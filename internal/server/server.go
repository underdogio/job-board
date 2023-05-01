package server

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	stdtemplate "html/template"

	"github.com/dustin/go-humanize"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/underdogio/job-board/internal/company"
	"github.com/underdogio/job-board/internal/config"
	"github.com/underdogio/job-board/internal/database"
	"github.com/underdogio/job-board/internal/email"
	"github.com/underdogio/job-board/internal/job"
	"github.com/underdogio/job-board/internal/middleware"
	"github.com/underdogio/job-board/internal/repositories"
	"github.com/underdogio/job-board/internal/template"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/allegro/bigcache/v3"
)

const (
	CacheKeyPinnedJobs       = "pinnedJobs"
	CacheKeyNewJobsLastWeek  = "newJobsLastWeek"
	CacheKeyNewJobsLastMonth = "newJobsLastMonth"
)

type Server struct {
	cfg          config.Config
	Conn         *sql.DB
	router       *mux.Router
	tmpl         *template.Template
	emailClient  email.Client
	SessionStore *sessions.CookieStore
	bigCache     *bigcache.BigCache
	emailRe      *regexp.Regexp
}

func NewServer(
	cfg config.Config,
	conn *sql.DB,
	r *mux.Router,
	t *template.Template,
	emailClient email.Client,
	sessionStore *sessions.CookieStore,
) Server {
	bigCache, err := bigcache.NewBigCache(bigcache.DefaultConfig(12 * time.Hour))
	boil.SetDB(conn)
	svr := Server{
		cfg:          cfg,
		Conn:         conn,
		router:       r,
		tmpl:         t,
		emailClient:  emailClient,
		SessionStore: sessionStore,
		bigCache:     bigCache,
		emailRe:      regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"),
	}
	if err != nil {
		svr.Log(err, "unable to initialise big cache")
	}

	return svr
}

func (s Server) RegisterRoute(path string, handler func(w http.ResponseWriter, r *http.Request), methods []string) {
	s.router.HandleFunc(path, handler).Methods(methods...)
}

func (s Server) RegisterPathPrefix(path string, handler http.Handler, methods []string) {
	s.router.PathPrefix(path).Handler(handler).Methods(methods...)
}

func (s Server) StringToHTML(str string) stdtemplate.HTML {
	return s.tmpl.StringToHTML(str)
}

func (s Server) JSEscapeString(str string) string {
	return s.tmpl.JSEscapeString(str)
}

func (s Server) MarkdownToHTML(str string) stdtemplate.HTML {
	return s.tmpl.MarkdownToHTML(str)
}

func (s Server) GetConfig() config.Config {
	return s.cfg
}

/*
func (s Server) RenderSalaryForLocation(ctx context.Context, w http.ResponseWriter, r *http.Request,  location string) {
	loc, err := database.FindSeoLocationG(ctx, location)
	complimentaryRemote := false
	if err != nil {
		complimentaryRemote = true
		loc.Name = "Remote"
		loc.Currency = "$"
	}
	set, err := database.GetSalaryDataForLocationAndCurrency(s.Conn, loc.Name, loc.Currency)
	if err != nil {
		s.Log(err, fmt.Sprintf("unable to retrieve salary stats for location %s and currency %s, err: %#v", location, loc.Currency, err))
		s.JSON(w, http.StatusInternalServerError, map[string]string{"status": "error"})
		return
	}
	trendSet, err := database.GetSalaryTrendsForLocationAndCurrency(s.Conn, loc.Name, loc.Currency)
	if err != nil {
		s.Log(err, fmt.Sprintf("unable to retrieve salary trends for location %s and currency %s, err: %#v", location, loc.Currency, err))
		s.JSON(w, http.StatusInternalServerError, map[string]string{"status": "error"})
		return
	}
	if len(set) < 1 {
		complimentaryRemote = true
		set, err = database.GetSalaryDataForLocationAndCurrency(s.Conn, "Remote", "$")
		if err != nil {
			s.Log(err, fmt.Sprintf("unable to retrieve salary stats for location %s and currency %s, err: %#v", location, loc.Currency, err))
			s.JSON(w, http.StatusInternalServerError, map[string]string{"status": "error"})
			return
		}
		trendSet, err = database.GetSalaryTrendsForLocationAndCurrency(s.Conn, "Remote", "$")
		if err != nil {
			s.Log(err, fmt.Sprintf("unable to retrieve salary stats for location %s and currency %s, err: %#v", location, loc.Currency, err))
			s.JSON(w, http.StatusInternalServerError, map[string]string{"status": "error"})
			return
		}
	}
	jsonRes, err := json.Marshal(set)
	if err != nil {
		s.Log(err, fmt.Sprintf("unable to marshal data set %v, err: %#v", set, err))
		s.JSON(w, http.StatusInternalServerError, map[string]string{"status": "error"})
		return
	}
	jsonTrendRes, err := json.Marshal(trendSet)
	if err != nil {
		s.Log(err, fmt.Sprintf("unable to marshal data set trneds %v, err: %#v", trendSet, err))
		s.JSON(w, http.StatusInternalServerError, map[string]string{"status": "error"})
		return
	}
	var sampleMin, sampleMax stats.Sample
	for _, x := range set {
		sampleMin.Xs = append(sampleMin.Xs, float64(x.Min))
		sampleMax.Xs = append(sampleMax.Xs, float64(x.Max))
	}
	min, _ := sampleMin.Bounds()
	_, max := sampleMax.Bounds()
	min = min - 30000
	max = max + 30000
	if min < 0 {
		min = 0
	}
	ua := r.Header.Get("user-agent")
	ref := r.Header.Get("referer")
	ips := strings.Split(r.Header.Get("x-forwarded-for"), ", ")
	if len(ips) > 0 && strings.Contains(ref, s.GetConfig().SiteHost) {
		hashedIP := sha256.Sum256([]byte(ips[0]))
		go func() {
			if err := database.TrackSearchEvent(s.Conn, ua, hex.EncodeToString(hashedIP[:]), location, "", len(set), job.SearchTypeSalary); err != nil {
				fmt.Printf("err while saving loc: %s\n", err)
			}
		}()
	}
	jobPosts, err := repositories.TopNJobsByLocation(loc.Name, 3)
	if err != nil {
		s.Log(err, "repositories.TopNJobsByLocation")
	}
	if len(jobPosts) == 0 {
		jobPosts, err = repositories.TopNJobsByLocation("Remote", 3)
		if err != nil {
			s.Log(err, "repositories.TopNJobsByLocation")
		}
	}
	lastJobPosted, err := repositories.LastJobPosted()
	if err != nil {
		s.Log(err, "could not retrieve last job posted at")
		lastJobPosted = time.Now().AddDate(0, 0, -1)
	}

	emailSubscribersCount, err := database.CountEmailSubscribers(s.Conn)
	if err != nil {
		s.Log(err, "database.CountEmailSubscribers")
	}
	topDevelopers, err := repositories.GetTopDevelopers(10)
	if err != nil {
		s.Log(err, "unable to get top developers")
	}
	topDeveloperSkills, err := repositories.GetTopDeveloperSkills(7)
	if err != nil {
		s.Log(err, "unable to get top developer skills")
	}
	lastDevUpdatedAt, err := repositories.GetLastDevUpdatedAt()
	if err != nil {
		s.Log(err, "unable to retrieve last developer joined at")
	}
	topDeveloperNames := make([]string, 0, len(topDevelopers))
	for _, d := range topDevelopers {
		topDeveloperNames = append(topDeveloperNames, strings.Split(d.Name, " ")[0])
	}
	messagesSentLastMonth, err := repositories.GetDeveloperMessagesSentLastMonth()
	if err != nil {
		s.Log(err, "GetDeveloperMessagesSentLastMonth")
	}
	devsRegisteredLastMonth, err := repositories.GetDevelopersRegisteredLastMonth()
	if err != nil {
		s.Log(err, "GetDevelopersRegisteredLastMonth")
	}
	devPageViewsLastMonth, err := repositories.GetDeveloperProfilePageViewsLastMonth()
	if err != nil {
		s.Log(err, "GetDeveloperProfilePageViewsLastMonth")
	}

	s.Render(r, w, http.StatusOK, "salary-explorer.html", map[string]interface{}{
		"Location":                           strings.ReplaceAll(location, "-", " "),
		"LocationURLEncoded":                 url.PathEscape(strings.ReplaceAll(location, "-", " ")),
		"Currency":                           loc.Currency,
		"DataSet":                            string(jsonRes),
		"DataSetTrends":                      string(jsonTrendRes),
		"TextCompanies":                      textifyCompanies(loc.Name, jobPosts, jobPosts),
		"TextJobTitles":                      textifyJobTitles(jobPosts),
		"P10Max":                             humanize.Comma(int64(math.Round(sampleMax.Quantile(0.1)))),
		"P10Min":                             humanize.Comma(int64(math.Round(sampleMin.Quantile(0.1)))),
		"P50Max":                             humanize.Comma(int64(math.Round(sampleMax.Quantile(0.5)))),
		"P50Min":                             humanize.Comma(int64(math.Round(sampleMin.Quantile(0.5)))),
		"P90Max":                             humanize.Comma(int64(math.Round(sampleMax.Quantile(0.9)))),
		"P90Min":                             humanize.Comma(int64(math.Round(sampleMin.Quantile(0.9)))),
		"MeanMin":                            humanize.Comma(int64(math.Round(sampleMin.Mean()))),
		"MeanMax":                            humanize.Comma(int64(math.Round(sampleMax.Mean()))),
		"StdDevMin":                          humanize.Comma(int64(math.Round(sampleMin.StdDev()))),
		"StdDevMax":                          humanize.Comma(int64(math.Round(sampleMax.StdDev()))),
		"Count":                              len(set),
		"Country":                            loc.Country,
		"Region":                             loc.Region,
		"Population":                         loc.Population,
		"Min":                                int64(math.Round(min)),
		"Max":                                int64(math.Round(max)),
		"ComplimentaryRemote":                complimentaryRemote,
		"LastJobPostedAt":                    lastJobPosted.Format(time.RFC3339),
		"LastJobPostedAtHumanized":           humanize.Time(lastJobPosted),
		"MonthAndYear":                       time.Now().UTC().Format("January 2006"),
		"EmailSubscribersCount":              humanize.Comma(int64(emailSubscribersCount)),
		"TopDevelopers":                      topDevelopers,
		"TopDeveloperNames":                  textifyGeneric(topDeveloperNames),
		"TopDeveloperSkills":                 textifyGeneric(topDeveloperSkills),
		"DeveloperMessagesSentLastMonth":     messagesSentLastMonth,
		"DevelopersRegisteredLastMonth":      devsRegisteredLastMonth,
		"DeveloperProfilePageViewsLastMonth": devPageViewsLastMonth,
		"LastDevCreatedAt":                   lastDevUpdatedAt.Format(time.RFC3339),
		"LastDevCreatedAtHumanized":          humanize.Time(lastDevUpdatedAt),
	})
}
*/

var hasBot = regexp.MustCompile(`(?i)(googlebot|bingbot|slurp|baiduspider|duckduckbot|yandexbot|sogou|exabot|facebookexternalhit|facebot|ia_archiver|linkedinbot|python-urllib|python-requests|go-http-client|msnbot|ahrefs)`)

func (s Server) RenderPageForLocationAndTag(w http.ResponseWriter, r *http.Request, location, tag, page, salary, currency, htmlView string) {
	var validSalary bool
	for _, band := range s.GetConfig().AvailableSalaryBands {
		if fmt.Sprintf("%d", band) == salary {
			validSalary = true
			break
		}
	}
	var validCurrency bool
	for _, availableCurrency := range s.GetConfig().AvailableCurrencies {
		if availableCurrency == currency {
			validCurrency = true
			break
		}
	}
	if (salary != "" && !validSalary) || (currency != "" && !validCurrency) {
		s.Redirect(w, r, http.StatusMovedPermanently, "/")
		return
	}
	showPage := true
	if page == "" {
		page = "1"
		showPage = false
	}
	tag = strings.TrimSpace(tag)
	location = strings.TrimSpace(location)
	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
	if err != nil {
		s.Log(err, "unable to compile regex (this should never happen)")
	}
	tag = reg.ReplaceAllString(tag, "")
	location = reg.ReplaceAllString(location, "")
	pageID, err := strconv.Atoi(page)
	if err != nil {
		pageID = 1
		showPage = false
	}
	isLandingPage := tag == "" && location == "" && page == "1" && salary == ""
	var newJobsLastWeek, newJobsLastMonth int
	newJobsLastWeekCached, okWeek := s.CacheGet(CacheKeyNewJobsLastWeek)
	newJobsLastMonthCached, okMonth := s.CacheGet(CacheKeyNewJobsLastMonth)
	if !okMonth || !okWeek {
		// load and cache last jobs count
		var err error
		newJobsLastWeek, newJobsLastMonth, err = repositories.NewJobsLastWeekOrMonth(r.Context())
		if err != nil {
			s.Log(err, "unable to retrieve new jobs last week last month")
		}
		buf := &bytes.Buffer{}
		enc := gob.NewEncoder(buf)
		if err := enc.Encode(newJobsLastWeek); err != nil {
			s.Log(err, "unable to encode new jobs last week")
		}
		if err := s.CacheSet(CacheKeyNewJobsLastWeek, buf.Bytes()); err != nil {
			s.Log(err, "unable to cache set new jobs lat week")
		}
		buf.Reset()
		if err := enc.Encode(newJobsLastMonth); err != nil {
			s.Log(err, "unable to encode new jobs last month")
		}
		if err := s.CacheSet(CacheKeyNewJobsLastMonth, buf.Bytes()); err != nil {
			s.Log(err, "unable to cache set new jobs lat month")
		}
	} else {
		dec := gob.NewDecoder(bytes.NewReader(newJobsLastWeekCached))
		if err := dec.Decode(&newJobsLastWeek); err != nil {
			s.Log(err, "unable to decode cached new jobs last week")
		}
		dec = gob.NewDecoder(bytes.NewReader(newJobsLastMonthCached))
		if err := dec.Decode(&newJobsLastMonth); err != nil {
			s.Log(err, "unable to decode cached new jobs last month")
		}
	}
	var pinnedJobs []*database.Job
	// only load pinned jobs for main landing page
	if isLandingPage {
		pinnedJobsCached, ok := s.CacheGet(CacheKeyPinnedJobs)
		if !ok {
			// load and cache jobs
			pinnedJobs, err = repositories.GetPinnedJobs(r.Context())
			if err != nil {
				s.Log(err, "unable to get pinned jobs")
			}
			for i, j := range pinnedJobs {
				pinnedJobs[i].Description = string(s.tmpl.MarkdownToHTML(j.Description))
			}
			buf := &bytes.Buffer{}
			enc := gob.NewEncoder(buf)
			if err := enc.Encode(pinnedJobs); err != nil {
				s.Log(err, "unable to encode pinned jobs")
			}
			if err := s.CacheSet(CacheKeyPinnedJobs, buf.Bytes()); err != nil {
				s.Log(err, "unable to set pinnedJobs cache")
			}
		} else {
			// pinned jobs are cached
			dec := gob.NewDecoder(bytes.NewReader(pinnedJobsCached))
			if err := dec.Decode(&pinnedJobs); err != nil {
				s.Log(err, "unable to decode pinned jobs")
			}
		}
	}
	jobsForPage, totalJobCount, err := repositories.JobsByQuery(r.Context(), location, tag, pageID, salary, currency, s.cfg.JobsPerPage, !isLandingPage)
	if err != nil {
		s.Log(err, "unable to get jobs by query")
		s.JSON(w, http.StatusInternalServerError, "Oops! An internal error has occurred")
		return
	}
	var complementaryRemote bool
	if len(jobsForPage) == 0 {
		complementaryRemote = true
		jobsForPage, totalJobCount, err = repositories.JobsByQuery(r.Context(), "Remote", tag, pageID, salary, currency, s.cfg.JobsPerPage, !isLandingPage)
		if len(jobsForPage) == 0 {
			jobsForPage, totalJobCount, err = repositories.JobsByQuery(r.Context(), "Remote", "", pageID, salary, currency, s.cfg.JobsPerPage, !isLandingPage)
		}
	}
	if err != nil {
		s.Log(err, "unable to retrieve jobs by query")
		s.JSON(w, http.StatusInternalServerError, "Oops! An internal error has occurred")
		return
	}
	pages := []int{}
	pageLinksPerPage := 8
	pageLinkShift := ((pageLinksPerPage / 2) + 1)
	firstPage := 1
	if pageID-pageLinkShift > 0 {
		firstPage = pageID - pageLinkShift
	}
	for i, j := firstPage, 1; i <= totalJobCount/s.cfg.JobsPerPage+1 && j <= pageLinksPerPage; i, j = i+1, j+1 {
		pages = append(pages, i)
	}

	var locFromDB = &database.SeoLocation{
		Name:     "Remote",
		Currency: "$",
	}
	if location != "" && !strings.EqualFold(location, "remote") {
		locFromDB, err = database.SeoLocations(database.SeoLocationWhere.Name.EQ(location)).OneG(r.Context())
		if err != nil {
			locFromDB.Name = "Remote"
			locFromDB.Currency = "$"
		}
	}
	var minSalary int64 = 1<<63 - 1
	var maxSalary int64 = 0
	for i, j := range jobsForPage {
		jobsForPage[i].Description = string(s.tmpl.MarkdownToHTML(j.Description))
	}

	ua := r.Header.Get("user-agent")
	ref := r.Header.Get("referer")
	ips := strings.Split(r.Header.Get("x-forwarded-for"), ", ")
	if len(ips) > 0 && strings.Contains(ref, s.GetConfig().SiteHost) {
		hashedIP := sha256.Sum256([]byte(ips[0]))
		go func() {
			if hasBot.MatchString(ua) {
				return
			}
			loc := strings.TrimSpace(location)
			trimmerTag := strings.TrimSpace(tag)
			if loc == "" && tag == "" {
				return
			}
			searchEvent := database.SearchEvent{
				SessionID: hex.EncodeToString(hashedIP[:]),
				Location:  null.StringFrom(loc),
				Tag:       null.StringFrom(trimmerTag),
				Results:   len(jobsForPage),
				Type:      null.StringFrom(job.SearchTypeJob),
			}
			if err := searchEvent.InsertG(r.Context(), boil.Infer()); err != nil {
				fmt.Printf("err while saving event: %s\n", err)
			}
		}()
	}

	locationWithCountry := strings.Title(location)
	relatedLocations := make([]string, 0)
	if locFromDB.Name != "Remote" {
		locationWithCountry = fmt.Sprintf("%s", locFromDB.Name)
		if locFromDB.Country.Valid && locFromDB.Country.String != "" {
			locationWithCountry = fmt.Sprintf("%s, %s", locFromDB.Name, locFromDB.Country.String)
		}
		if locFromDB.Region.Valid && locFromDB.Region.String != "" {
			locationWithCountry = fmt.Sprintf("%s, %s, %s", locFromDB.Name, locFromDB.Region.String, locFromDB.Country.String)
		}
		if err := database.SeoLocations(
			qm.Select(database.SeoLocationColumns.Name),
			database.SeoLocationWhere.Country.EQ(locFromDB.Country),
			database.SeoLocationWhere.Name.NEQ(locFromDB.Name),
			qm.OrderBy(database.SeoLocationColumns.Population+" DESC"),
			qm.Limit(6),
		).BindG(r.Context(), relatedLocations); err != nil {
			s.Log(err, fmt.Sprintf("unable to get random locations for country %s", locFromDB.Country.String))
		}
	}
	if currency == "" {
		currency = "USD"
	}
	lastJobPosted, err := repositories.LastJobPosted(r.Context())
	if err != nil {
		s.Log(err, "could not retrieve last job posted at")
		lastJobPosted = time.Now().AddDate(0, 0, -1)
	}

	emailSubscribersCount, err := database.EmailSubscribers(
		database.EmailSubscriberWhere.ConfirmedAt.IsNotNull(),
	).CountG(r.Context())
	if err != nil {
		s.Log(err, "database.CountEmailSubscribers")
	}

	s.Render(r, w, http.StatusOK, htmlView, map[string]interface{}{
		"Jobs":                      jobsForPage,
		"PinnedJobs":                pinnedJobs,
		"JobsMinusOne":              len(jobsForPage) - 1,
		"LocationFilter":            strings.Title(location),
		"LocationFilterWithCountry": locationWithCountry,
		"LocationFilterURLEnc":      url.PathEscape(strings.Title(location)),
		"TagFilter":                 tag,
		"SalaryFilter":              salary,
		"CurrencyFilter":            currency,
		"AvailableCurrencies":       s.GetConfig().AvailableCurrencies,
		"AvailableSalaryBands":      s.GetConfig().AvailableSalaryBands,
		"TagFilterURLEnc":           url.PathEscape(tag),
		"CurrentPage":               pageID,
		"ShowPage":                  showPage,
		"PageSize":                  s.cfg.JobsPerPage,
		"PageIndexes":               pages,
		"TotalJobCount":             totalJobCount,
		"TextJobCount":              textifyJobCount(totalJobCount),
		"TextCompanies":             textifyCompanies(location, pinnedJobs, jobsForPage),
		"TextJobTitles":             textifyJobTitles(jobsForPage),
		"LastJobPostedAt":           lastJobPosted.Format(time.RFC3339),
		"LastJobPostedAtHumanized":  humanize.Time(lastJobPosted),
		"HasSalaryInfo":             maxSalary > 0,
		"MinSalary":                 fmt.Sprintf("%s%s", locFromDB.Currency, humanize.Comma(minSalary)),
		"MaxSalary":                 fmt.Sprintf("%s%s", locFromDB.Currency, humanize.Comma(maxSalary)),
		"LocationFromDB":            locFromDB.Name,
		"CountryFromDB":             locFromDB.Country,
		"RegionFromDB":              locFromDB.Region,
		"PopulationFromDB":          locFromDB.Population,
		"LocationEmojiFromDB":       locFromDB.Emoji,
		"RelatedLocations":          relatedLocations,
		"ComplementaryRemote":       complementaryRemote,
		"MonthAndYear":              time.Now().UTC().Format("January 2006"),
		"NewJobsLastWeek":           newJobsLastWeek,
		"NewJobsLastMonth":          newJobsLastMonth,
		"EmailSubscribersCount":     humanize.Comma(int64(emailSubscribersCount)),
	})
}

func textifyJobCount(n int) string {
	if n <= 50 {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%d+", (n/50)*50)
}

func textifyCompanies(location string, pinnedJobs, jobs []*database.Job) string {
	if len(pinnedJobs) > 2 && location == "" {
		jobs = pinnedJobs
	}
	switch {
	case len(jobs) == 1:
		return jobs[0].Company
	case len(jobs) == 2:
		return fmt.Sprintf("%s and %s", jobs[0].Company, jobs[1].Company)
	case len(jobs) > 2:
		return fmt.Sprintf("%s, %s and %s", jobs[0].Company, jobs[1].Company, jobs[2].Company)
	}

	return ""
}

func textifyGeneric(items []string) string {
	switch {
	case len(items) == 1:
		return items[0]
	case len(items) == 2:
		return fmt.Sprintf("%s and %s", items[0], items[1])
	case len(items) > 2:
		return fmt.Sprintf("%s and %s", strings.Join(items[:len(items)-1], ", "), items[len(items)-1])
	}

	return ""
}

func textifyCompanyNames(companies []*database.Company, max int) string {
	switch {
	case len(companies) == 1:
		return companies[0].Name
	case len(companies) == 2:
		return fmt.Sprintf("%s and %s", companies[0].Name, companies[1].Name)
	case len(companies) > 2:
		names := make([]string, 0, len(companies))
		if max >= len(companies)-1 {
			max = len(companies) - 1
		}
		for i := 0; i < max; i++ {
			names = append(names, companies[i].Name)
		}
		return fmt.Sprintf("%s and many others", strings.Join(names, ", "))
	}

	return ""
}

func textifyJobTitles(jobs []*database.Job) string {
	switch {
	case len(jobs) == 1:
		return jobs[0].JobTitle
	case len(jobs) == 2:
		return fmt.Sprintf("%s and %s", jobs[0].JobTitle, jobs[1].JobTitle)
	case len(jobs) > 2:
		return fmt.Sprintf("%s, %s and %s", jobs[0].JobTitle, jobs[1].JobTitle, jobs[2].JobTitle)
	}

	return ""
}

func (s Server) RenderPageForProfileRegistration(w http.ResponseWriter, r *http.Request, htmlView string) {
	topDevelopers, err := repositories.GetTopDevelopers(r.Context(), 10)
	if err != nil {
		s.Log(err, "unable to get top developers")
	}
	topDeveloperSkills, err := repositories.GetTopDeveloperSkills(r.Context(), 7)
	if err != nil {
		s.Log(err, "unable to get top developer skills")
	}
	lastDevUpdatedAt, err := repositories.GetLastDevUpdatedAt(r.Context())
	if err != nil {
		s.Log(err, "unable to retrieve last developer joined at")
	}
	topDeveloperNames := make([]string, 0, len(topDevelopers))
	for _, d := range topDevelopers {
		topDeveloperNames = append(topDeveloperNames, strings.Split(d.Name, " ")[0])
	}
	messagesSentLastMonth, err := repositories.GetDeveloperMessagesSentLastMonth(r.Context())
	if err != nil {
		s.Log(err, "GetDeveloperMessagesSentLastMonth")
	}
	devsRegisteredLastMonth, err := repositories.GetDevelopersRegisteredLastMonth(r.Context())
	if err != nil {
		s.Log(err, "GetDevelopersRegisteredLastMonth")
	}
	devPageViewsLastMonth, err := repositories.GetDeveloperProfilePageViewsLastMonth(r.Context())
	if err != nil {
		s.Log(err, "GetDeveloperProfilePageViewsLastMonth")
	}

	s.Render(r, w, http.StatusOK, htmlView, map[string]interface{}{
		"TopDevelopers":                      topDevelopers,
		"TopDeveloperNames":                  textifyGeneric(topDeveloperNames),
		"TopDeveloperSkills":                 textifyGeneric(topDeveloperSkills),
		"DeveloperMessagesSentLastMonth":     messagesSentLastMonth,
		"DevelopersRegisteredLastMonth":      devsRegisteredLastMonth,
		"DeveloperProfilePageViewsLastMonth": devPageViewsLastMonth,
		"MonthAndYear":                       time.Now().UTC().Format("January 2006"),
		"LastDevCreatedAt":                   lastDevUpdatedAt.Format(time.RFC3339),
		"LastDevCreatedAtHumanized":          humanize.Time(lastDevUpdatedAt),
	})
}

type DeveloperView struct {
	Developer          *database.DeveloperProfile
	CreatedAtHumanized string
	UpdatedAtHumanized string
	SkillsArray        []string
}

func (s Server) RenderPageForDevelopers(w http.ResponseWriter, r *http.Request, location, tag, page, htmlView string) {
	showPage := true
	if page == "" {
		page = "1"
		showPage = false
	}
	location = strings.TrimSpace(location)
	tag = strings.TrimSpace(tag)
	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
	if err != nil {
		s.Log(err, "unable to compile regex (this should never happen)")
	}
	location = reg.ReplaceAllString(location, "")
	tag = reg.ReplaceAllString(tag, "")
	pageID, err := strconv.Atoi(page)
	if err != nil {
		pageID = 1
		showPage = false
	}
	var complementaryRemote bool
	locSearch := location
	if strings.EqualFold(location, "remote") {
		locSearch = ""
	}
	developers, totalDevelopersCount, err := repositories.DevelopersByLocationAndTag(r.Context(), locSearch, tag, pageID, s.cfg.DevelopersPerPage)
	if err != nil {
		s.Log(err, "unable to get developers by location and tag")
		s.JSON(w, http.StatusInternalServerError, "Oops! An internal error has occurred")
		return
	}
	if len(developers) == 0 {
		complementaryRemote = true
		developers, totalDevelopersCount, err = repositories.DevelopersByLocationAndTag(r.Context(), "", "", pageID, s.cfg.DevelopersPerPage)
	}
	pages := []int{}
	pageLinksPerPage := 8
	pageLinkShift := ((pageLinksPerPage / 2) + 1)
	firstPage := 1
	if pageID-pageLinkShift > 0 {
		firstPage = pageID - pageLinkShift
	}
	for i, j := firstPage, 1; i <= totalDevelopersCount/s.cfg.DevelopersPerPage+1 && j <= pageLinksPerPage; i, j = i+1, j+1 {
		pages = append(pages, i)
	}
	developersForPage := make([]*DeveloperView, len(developers))
	for i, j := range developers {
		developersForPage[i] = &DeveloperView{
			Developer:          j,
			CreatedAtHumanized: humanize.Time(j.CreatedAt.UTC()),
			UpdatedAtHumanized: j.UpdatedAt.Time.UTC().Format("January 2006"),
			SkillsArray:        strings.Split(j.Skills, ","),
		}
	}
	ref := r.Header.Get("referer")
	ips := strings.Split(r.Header.Get("x-forwarded-for"), ", ")
	if len(ips) > 0 && strings.Contains(ref, s.GetConfig().SiteHost) {
		hashedIP := sha256.Sum256([]byte(ips[0]))
		go func() {
			searchEvent := &database.SearchEvent{
				SessionID: hex.EncodeToString(hashedIP[:]),
				Location:  null.StringFrom(location),
				Tag:       null.StringFrom(tag),
				Results:   len(developersForPage),
				Type:      null.StringFrom(repositories.SearchTypeDeveloper),
			}
			if err := searchEvent.InsertG(r.Context(), boil.Infer()); err != nil {
				fmt.Printf("err while saving event: %s\n", err)
			}
		}()
	}
	loc, err := database.SeoLocations(database.SeoLocationWhere.Name.EQ(location)).OneG(r.Context())
	if err != nil {
		loc.Name = "Remote"
		loc.Currency = "$"
	}
	topDevelopers, err := repositories.GetTopDevelopers(r.Context(), 10)
	if err != nil {
		s.Log(err, "unable to get top developer names")
	}
	topDeveloperSkills, err := repositories.GetTopDeveloperSkills(r.Context(), 7)
	if err != nil {
		s.Log(err, "unable to get top developer skills")
	}
	topDeveloperNames := make([]string, 0, len(topDevelopers))
	for _, d := range topDevelopers {
		topDeveloperNames = append(topDeveloperNames, strings.Split(d.Name, " ")[0])
	}

	emailSubscribersCount, err := database.EmailSubscribers(
		database.EmailSubscriberWhere.ConfirmedAt.IsNotNull(),
	).CountG(r.Context())
	if err != nil {
		s.Log(err, "database.CountEmailSubscribers")
	}
	lastDevUpdatedAt, err := repositories.GetLastDevUpdatedAt(r.Context())
	if err != nil {
		s.Log(err, "unable to retrieve last developer joined at")
	}
	messagesSentLastMonth, err := repositories.GetDeveloperMessagesSentLastMonth(r.Context())
	if err != nil {
		s.Log(err, "GetDeveloperMessagesSentLastMonth")
	}
	devsRegisteredLastMonth, err := repositories.GetDevelopersRegisteredLastMonth(r.Context())
	if err != nil {
		s.Log(err, "GetDevelopersRegisteredLastMonth")
	}
	devPageViewsLastMonth, err := repositories.GetDeveloperProfilePageViewsLastMonth(r.Context())
	if err != nil {
		s.Log(err, "GetDeveloperProfilePageViewsLastMonth")
	}

	s.Render(r, w, http.StatusOK, htmlView, map[string]interface{}{
		"Developers":                         developersForPage,
		"TopDeveloperSkills":                 textifyGeneric(topDeveloperSkills),
		"DevelopersMinusOne":                 len(developersForPage) - 1,
		"LocationFilter":                     strings.Title(location),
		"LocationURLEncoded":                 url.PathEscape(strings.ReplaceAll(location, "-", " ")),
		"TextCount":                          textifyJobCount(totalDevelopersCount),
		"TagFilter":                          tag,
		"TagFilterURLEncoded":                url.PathEscape(tag),
		"CurrentPage":                        pageID,
		"ShowPage":                           showPage,
		"PageSize":                           s.cfg.DevelopersPerPage,
		"Country":                            loc.Country,
		"Region":                             loc.Region,
		"PageIndexes":                        pages,
		"TotalDevelopersCount":               totalDevelopersCount,
		"ComplementaryRemote":                complementaryRemote,
		"MonthAndYear":                       time.Now().UTC().Format("January 2006"),
		"EmailSubscribersCount":              humanize.Comma(int64(emailSubscribersCount)),
		"DevelopersBannerLink":               s.GetConfig().DevelopersBannerLink,
		"DevelopersBannerText":               s.GetConfig().DevelopersBannerText,
		"TopDevelopers":                      topDevelopers,
		"TopDeveloperNames":                  textifyGeneric(topDeveloperNames),
		"DeveloperMessagesSentLastMonth":     messagesSentLastMonth,
		"DevelopersRegisteredLastMonth":      devsRegisteredLastMonth,
		"DeveloperProfilePageViewsLastMonth": devPageViewsLastMonth,
		"LastDevCreatedAt":                   lastDevUpdatedAt.Format(time.RFC3339),
		"LastDevCreatedAtHumanized":          humanize.Time(lastDevUpdatedAt),
	})

}

func (s Server) RenderPageForCompanies(w http.ResponseWriter, r *http.Request, location, page, htmlView string) {
	showPage := true
	if page == "" {
		page = "1"
		showPage = false
	}
	location = strings.TrimSpace(location)
	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
	if err != nil {
		s.Log(err, "unable to compile regex (this should never happen)")
	}
	location = reg.ReplaceAllString(location, "")
	pageID, err := strconv.Atoi(page)
	if err != nil {
		pageID = 1
		showPage = false
	}
	var complementaryRemote bool
	companiesForPage, totalCompaniesCount, err := repositories.CompaniesByQuery(r.Context(), location, pageID, s.cfg.CompaniesPerPage)
	if err != nil {
		s.Log(err, "unable to get companies by query")
		s.JSON(w, http.StatusInternalServerError, "Oops! An internal error has occurred")
		return
	}
	if len(companiesForPage) == 0 {
		complementaryRemote = true
		companiesForPage, totalCompaniesCount, err = repositories.CompaniesByQuery(r.Context(), "Remote", pageID, s.cfg.CompaniesPerPage)
	}
	loc, err := database.SeoLocations(database.SeoLocationWhere.Name.EQ(location)).OneG(r.Context())
	if err != nil {
		loc = &database.SeoLocation{}
		loc.Name = "Remote"
		loc.Currency = "$"
	}
	pageID64 := pageID
	pages := []int{}
	pageLinksPerPage := 8
	pageLinkShift := (pageLinksPerPage / 2) + 1
	firstPage := 1
	if pageID64-pageLinkShift > 0 {
		firstPage = pageID64 - pageLinkShift
	}
	for i, j := firstPage, 1; i <= totalCompaniesCount/s.cfg.CompaniesPerPage+1 && j <= pageLinksPerPage; i, j = i+1, j+1 {
		pages = append(pages, i)
	}

	ref := r.Header.Get("referer")
	ips := strings.Split(r.Header.Get("x-forwarded-for"), ", ")
	if len(ips) > 0 && strings.Contains(ref, s.GetConfig().SiteHost) {
		hashedIP := sha256.Sum256([]byte(ips[0]))
		go func() {
			trackEvent := &database.SearchEvent{
				SessionID: hex.EncodeToString(hashedIP[:]),
				Location:  null.StringFrom(location),
				Tag:       null.StringFrom(""),
				Results:   len(companiesForPage),
				Type:      null.StringFrom(company.SearchTypeCompany),
			}
			if err := trackEvent.InsertG(r.Context(), boil.Infer()); err != nil {
				fmt.Printf("err while saving event: %s\n", err)
			}
		}()
	}
	jobPosts, err := repositories.TopNJobsByLocation(r.Context(), loc.Name, 3)
	if err != nil {
		s.Log(err, "database.TopNJobsByLocation")
	}
	if len(jobPosts) == 0 {
		jobPosts, err = repositories.TopNJobsByLocation(r.Context(), "Remote", 3)
		if err != nil {
			s.Log(err, "database.TopNJobsByLocation")
		}
	}

	var lastJobPostedAt, lastJobPostedAtHumanized string
	if len(jobPosts) > 0 {
		lastJobPostedAt = jobPosts[0].CreatedAt.Format(time.RFC3339)
		lastJobPostedAtHumanized = humanize.Time(jobPosts[0].CreatedAt)
	}

	emailSubscribersCount, err := database.EmailSubscribers(
		database.EmailSubscriberWhere.ConfirmedAt.IsNotNull(),
	).CountG(r.Context())
	if err != nil {
		s.Log(err, "database.CountEmailSubscribers")
	}
	topDevelopers, err := repositories.GetTopDevelopers(r.Context(), 10)
	if err != nil {
		s.Log(err, "unable to get top developers")
	}
	topDeveloperSkills, err := repositories.GetTopDeveloperSkills(r.Context(), 7)
	if err != nil {
		s.Log(err, "unable to get top developer skills")
	}
	lastDevUpdatedAt, err := repositories.GetLastDevUpdatedAt(r.Context())
	if err != nil {
		s.Log(err, "unable to retrieve last developer joined at")
	}
	topDeveloperNames := make([]string, 0, len(topDevelopers))
	for _, d := range topDevelopers {
		topDeveloperNames = append(topDeveloperNames, strings.Split(d.Name, " ")[0])
	}
	messagesSentLastMonth, err := repositories.GetDeveloperMessagesSentLastMonth(r.Context())
	if err != nil {
		s.Log(err, "GetDeveloperMessagesSentLastMonth")
	}
	devsRegisteredLastMonth, err := repositories.GetDevelopersRegisteredLastMonth(r.Context())
	if err != nil {
		s.Log(err, "GetDevelopersRegisteredLastMonth")
	}
	devPageViewsLastMonth, err := repositories.GetDeveloperProfilePageViewsLastMonth(r.Context())
	if err != nil {
		s.Log(err, "GetDeveloperProfilePageViewsLastMonth")
	}

	s.Render(r, w, http.StatusOK, htmlView, map[string]interface{}{
		"Companies":                          companiesForPage,
		"CompaniesMinusOne":                  len(companiesForPage) - 1,
		"LocationFilter":                     strings.Title(location),
		"LocationURLEncoded":                 url.PathEscape(strings.ReplaceAll(location, "-", " ")),
		"TextCompanies":                      textifyCompanies(loc.Name, jobPosts, jobPosts),
		"TextJobTitles":                      textifyJobTitles(jobPosts),
		"TextJobCount":                       textifyJobCount(totalCompaniesCount),
		"CurrentPage":                        pageID,
		"ShowPage":                           showPage,
		"PageSize":                           s.cfg.CompaniesPerPage,
		"PageIndexes":                        pages,
		"TotalCompaniesCount":                totalCompaniesCount,
		"ComplementaryRemote":                complementaryRemote,
		"MonthAndYear":                       time.Now().UTC().Format("January 2006"),
		"Country":                            loc.Country,
		"Region":                             loc.Region,
		"Population":                         loc.Population,
		"LastJobPostedAt":                    lastJobPostedAt,
		"LastJobPostedAtHumanized":           lastJobPostedAtHumanized,
		"EmailSubscribersCount":              humanize.Comma(int64(emailSubscribersCount)),
		"TopDevelopers":                      topDevelopers,
		"TopDeveloperNames":                  textifyGeneric(topDeveloperNames),
		"TopDeveloperSkills":                 textifyGeneric(topDeveloperSkills),
		"DeveloperMessagesSentLastMonth":     messagesSentLastMonth,
		"DevelopersRegisteredLastMonth":      devsRegisteredLastMonth,
		"DeveloperProfilePageViewsLastMonth": devPageViewsLastMonth,
		"LastDevCreatedAt":                   lastDevUpdatedAt.Format(time.RFC3339),
		"LastDevCreatedAtHumanized":          humanize.Time(lastDevUpdatedAt),
	})
}

func (s Server) RenderPageForLocationAndTagAdmin(r *http.Request, w http.ResponseWriter, location, tag, page, salary, currency, htmlView string) {
	showPage := true
	if page == "" {
		page = "1"
		showPage = false
	}
	salaryInt, err := strconv.Atoi(salary)
	if err != nil {
		salaryInt = 0
	}
	salaryInt = int(salaryInt)
	tag = strings.TrimSpace(tag)
	location = strings.TrimSpace(location)
	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
	if err != nil {
		s.Log(err, "unable to compile regex (this should never happen)")
	}
	tag = reg.ReplaceAllString(tag, "")
	location = reg.ReplaceAllString(location, "")
	pageID, err := strconv.Atoi(page)
	if err != nil {
		pageID = 1
		showPage = false
	}
	var pinnedJobs []*database.Job
	pinnedJobs, err = repositories.GetPinnedJobs(r.Context())
	if err != nil {
		s.Log(err, "unable to get pinned jobs")
	}
	var pendingJobs []*database.Job
	pendingJobs, err = repositories.GetPendingJobs(r.Context())
	if err != nil {
		s.Log(err, "unable to get pending jobs")
	}
	for i, j := range pendingJobs {
		pendingJobs[i].SalaryRange = j.SalaryRange
	}
	jobsForPage, totalJobCount, err := repositories.JobsByQuery(r.Context(), location, tag, pageID, salary, currency, s.cfg.JobsPerPage, false)
	if err != nil {
		s.Log(err, "unable to get jobs by query")
		s.JSON(w, http.StatusInternalServerError, "Oops! An internal error has occurred")
		return
	}
	var complementaryRemote bool
	if len(jobsForPage) == 0 {
		complementaryRemote = true
		jobsForPage, totalJobCount, err = repositories.JobsByQuery(r.Context(), "Remote", tag, pageID, salary, currency, s.cfg.JobsPerPage, false)
		if len(jobsForPage) == 0 {
			jobsForPage, totalJobCount, err = repositories.JobsByQuery(r.Context(), "Remote", "", pageID, salary, currency, s.cfg.JobsPerPage, false)
		}
	}
	if err != nil {
		s.Log(err, "unable to retrieve jobs by query")
		s.JSON(w, http.StatusInternalServerError, "Oops! An internal error has occurred")
		return
	}
	pages := []int{}
	pageLinksPerPage := 8
	pageLinkShift := ((pageLinksPerPage / 2) + 1)
	firstPage := 1
	if pageID-pageLinkShift > 0 {
		firstPage = pageID - pageLinkShift
	}
	for i, j := firstPage, 1; i <= totalJobCount/s.cfg.JobsPerPage+1 && j <= pageLinksPerPage; i, j = i+1, j+1 {
		pages = append(pages, i)
	}
	for i, j := range jobsForPage {
		jobsForPage[i].Description = string(s.tmpl.MarkdownToHTML(j.Description))
	}
	for i, j := range pinnedJobs {
		pinnedJobs[i].Description = string(s.tmpl.MarkdownToHTML(j.Description))
	}

	s.Render(r, w, http.StatusOK, htmlView, map[string]interface{}{
		"Jobs":                jobsForPage,
		"PinnedJobs":          pinnedJobs,
		"PendingJobs":         pendingJobs,
		"JobsMinusOne":        len(jobsForPage) - 1,
		"LocationFilter":      location,
		"TagFilter":           tag,
		"CurrentPage":         pageID,
		"ShowPage":            showPage,
		"PageSize":            s.cfg.JobsPerPage,
		"PageIndexes":         pages,
		"TotalJobCount":       totalJobCount,
		"ComplementaryRemote": complementaryRemote,
		"MonthAndYear":        time.Now().UTC().Format("January 2006"),
	})
}

func (s Server) RenderPostAJobForLocation(w http.ResponseWriter, r *http.Request, location string) {
	var defaultJobPageviewsLast30Days = 15000
	var defaultJobApplicantsLast30Days = int64(1000)
	var defaultPageviewsLast30Days = int64(4000)

	pageviewsLast30Days := 0
	jobPageviewsLast30Days := int64(0)
	jobApplicantsLast30Days := int64(0)
	err := database.CloudflareBrowserStats(
		qm.Select("SUM("+database.CloudflareBrowserStatColumns.PageViews+") AS pageviews"),
		database.CloudflareBrowserStatWhere.Date.GT(time.Now().AddDate(0, 0, -30)),
		qm.And(database.CloudflareBrowserStatColumns.UaBrowserFamily+" NOT ILIKE ?", "%bot%"),
	).BindG(r.Context(), &pageviewsLast30Days)
	if err != nil {
		s.Log(err, "could not retrieve pageviews for last 30 days")
	}
	if pageviewsLast30Days < defaultJobPageviewsLast30Days {
		pageviewsLast30Days = defaultJobPageviewsLast30Days
	}
	jobPageviewsLast30Days, err = database.JobEvents(
		database.JobEventWhere.EventType.EQ("page_view"),
		database.JobEventWhere.CreatedAt.GT(time.Now().AddDate(0, 0, -30)),
	).CountG(r.Context())
	if err != nil {
		s.Log(err, "could not retrieve pageviews for last 30 days")
	}
	if err != nil {
		s.Log(err, "could not retrieve job pageviews for last 30 days")
	}
	if jobPageviewsLast30Days < defaultPageviewsLast30Days {
		jobPageviewsLast30Days = defaultPageviewsLast30Days
	}
	jobApplicantsLast30Days, err = database.JobEvents(
		database.JobEventWhere.EventType.EQ("clickout"),
		database.JobEventWhere.CreatedAt.GT(time.Now().AddDate(0, 0, -30)),
	).CountG(r.Context())
	if err != nil {
		s.Log(err, "could not retrieve job clickouts for last 30 days")
	}
	if jobApplicantsLast30Days < defaultJobApplicantsLast30Days {
		jobApplicantsLast30Days = defaultJobApplicantsLast30Days
	}
	featuredCompanies, err := repositories.FeaturedCompaniesPostAJob(r.Context())
	if err != nil {
		s.Log(err, "could not retrieve featured companies for post a job page")
	}
	lastJobPosted, err := repositories.LastJobPosted(r.Context())
	if err != nil {
		s.Log(err, "could not retrieve last job posted at")
		lastJobPosted = time.Now().AddDate(0, 0, -1)
	}
	newJobsLastWeek, newJobsLastMonth, err := repositories.NewJobsLastWeekOrMonth(r.Context())
	if err != nil {
		s.Log(err, "unable to retrieve new jobs last week last month")
		newJobsLastWeek = 1
	}
	s.Render(r, w, http.StatusOK, "post-a-job.html", map[string]interface{}{
		"Location":                 location,
		"PageviewsLastMonth":       humanize.Comma(int64(pageviewsLast30Days)),
		"JobPageviewsLastMonth":    humanize.Comma(int64(jobPageviewsLast30Days)),
		"JobApplicantsLastMonth":   humanize.Comma(int64(jobApplicantsLast30Days)),
		"FeaturedCompanies":        featuredCompanies,
		"FeaturedCompaniesNames":   textifyCompanyNames(featuredCompanies, 10),
		"LastJobPostedAtHumanized": humanize.Time(lastJobPosted),
		"LastJobPostedAt":          lastJobPosted.Format(time.RFC3339),
		"NewJobsLastWeek":          newJobsLastWeek,
		"NewJobsLastMonth":         newJobsLastMonth,
		"StripePublishableKey":     s.GetConfig().StripePublishableKey,
		"JobCategories":            job.JobCategories,
	})
}

func (s Server) Render(r *http.Request, w http.ResponseWriter, status int, htmlView string, data interface{}) error {
	dataMap := make(map[string]interface{}, 0)
	if data != nil {
		dataMap = data.(map[string]interface{})
	}
	profile, _ := middleware.GetUserFromJWT(r, s.SessionStore, s.GetJWTSigningKey())
	dataMap["LoggedUser"] = profile
	if profile != nil {
		dataMap["IsUserRecruiter"] = profile.IsRecruiter
		dataMap["IsUserDeveloper"] = profile.IsDeveloper
		dataMap["IsUserAdmin"] = profile.IsAdmin
	}
	dataMap["SiteName"] = s.GetConfig().SiteName
	dataMap["SiteJobCategory"] = strings.Title(strings.ToLower(s.GetConfig().SiteJobCategory))
	dataMap["SiteJobCategoryURLEncoded"] = strings.ReplaceAll(strings.Title(strings.ToLower(s.GetConfig().SiteJobCategory)), " ", "-")
	dataMap["SupportEmail"] = s.GetConfig().SupportEmail
	dataMap["SiteHost"] = s.GetConfig().SiteHost
	dataMap["SiteTwitter"] = s.GetConfig().SiteTwitter
	dataMap["SiteGithub"] = s.GetConfig().SiteGithub
	dataMap["SiteLinkedin"] = s.GetConfig().SiteLinkedin
	dataMap["SiteYoutube"] = s.GetConfig().SiteYoutube
	dataMap["SiteTelegramChannel"] = s.GetConfig().SiteTelegramChannel
	dataMap["PrimaryColor"] = s.GetConfig().PrimaryColor
	dataMap["SecondaryColor"] = s.GetConfig().SecondaryColor
	dataMap["SiteLogoImageID"] = s.GetConfig().SiteLogoImageID
	dataMap["PlanPriceID"] = s.GetConfig().PlanPriceID
	dataMap["PlanPrice"] = "299"
	dataMap["MonthAndYear"] = time.Now().UTC().Format("January 2006")

	return s.tmpl.Render(w, status, htmlView, dataMap)
}

func (s Server) XML(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(status)
	w.Write(data)
}

func (s Server) JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s Server) TEXT(w http.ResponseWriter, status int, text string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(text))
}

func (s Server) MEDIA(w http.ResponseWriter, status int, media []byte, mediaType string) {
	w.Header().Set("Content-Type", mediaType)
	w.Header().Set("Cache-Control", "max-age=31536000")
	w.WriteHeader(status)
	w.Write(media)
}

func (s Server) Log(err error, msg string) {
	log.Printf("%s: %+v", msg, err)
}

func (s Server) GetEmail() email.Client {
	return s.emailClient
}

func (s Server) Redirect(w http.ResponseWriter, r *http.Request, status int, dst string) {
	http.Redirect(w, r, dst, status)
}

func (s Server) Run() error {
	addr := fmt.Sprintf(":%s", s.cfg.Port)
	if s.cfg.Env == "dev" {
		log.Printf("local env http://localhost:%s", s.cfg.Port)
		addr = fmt.Sprintf("0.0.0.0:%s", s.cfg.Port)
	}
	return http.ListenAndServe(
		addr,
		middleware.HTTPSMiddleware(
			middleware.GzipMiddleware(
				middleware.LoggingMiddleware(middleware.HeadersMiddleware(s.router, s.cfg.Env)),
			),
			s.cfg.Env,
		),
	)
}

func (s Server) GetJWTSigningKey() []byte {
	return s.cfg.JwtSigningKey
}

func (s Server) CacheGet(key string) ([]byte, bool) {
	out, err := s.bigCache.Get(key)
	if err != nil {
		return []byte{}, false
	}
	return out, true
}

func (s Server) CacheSet(key string, val []byte) error {
	return s.bigCache.Set(key, val)
}

func (s Server) CacheDelete(key string) error {
	return s.bigCache.Delete(key)
}

func (s Server) SeenSince(r *http.Request, timeAgo time.Duration) bool {
	ipAddrs := strings.Split(r.Header.Get("x-forwarded-for"), ", ")
	if len(ipAddrs) == 0 {
		return false
	}
	lastSeen, err := s.bigCache.Get(ipAddrs[0])
	if err == bigcache.ErrEntryNotFound {
		s.bigCache.Set(ipAddrs[0], []byte(time.Now().Format(time.RFC3339)))
		return false
	}
	if err != nil {
		return false
	}
	lastSeenTime, err := time.Parse(time.RFC3339, string(lastSeen))
	if err != nil {
		s.bigCache.Set(ipAddrs[0], []byte(time.Now().Format(time.RFC3339)))
		return false
	}
	if !lastSeenTime.After(time.Now().Add(-timeAgo)) {
		s.bigCache.Set(ipAddrs[0], []byte(time.Now().Format(time.RFC3339)))
		return false
	}

	return true
}

func (s Server) IsEmail(val string) bool {
	return s.emailRe.MatchString(val)
}

// GetDbConn tries to establish a connection to postgres and return the connection handler
func GetDbConn(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}

// CloseDbConn closes db conn
func CloseDbConn(conn *sql.DB) {
	conn.Close()
}
