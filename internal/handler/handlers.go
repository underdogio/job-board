package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bot-api/telegram"
	jwt "github.com/dgrijalva/jwt-go"
	humanize "github.com/dustin/go-humanize"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/machinebox/graphql"
	"github.com/microcosm-cc/bluemonday"
	"github.com/nfnt/resize"
	"github.com/segmentio/ksuid"
	"github.com/snabb/sitemap"
	"github.com/underdogio/job-board/internal/blog"
	"github.com/underdogio/job-board/internal/database"
	"github.com/underdogio/job-board/internal/developer"
	"github.com/underdogio/job-board/internal/email"
	"github.com/underdogio/job-board/internal/imagemeta"
	"github.com/underdogio/job-board/internal/job"
	"github.com/underdogio/job-board/internal/middleware"
	"github.com/underdogio/job-board/internal/payment"
	"github.com/underdogio/job-board/internal/repositories"
	"github.com/underdogio/job-board/internal/seo"
	"github.com/underdogio/job-board/internal/server"
	"github.com/underdogio/job-board/internal/user"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const (
	AuthStepVerifyDeveloperProfile = "1mCQFVDZTAx9VQa1lprjr0aLgoP"
	AuthStepLoginDeveloperProfile  = "1mEvrSr2G4e4iGeucwolKW6o64d"
)

func GetAuthPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svr.Render(r, w, http.StatusOK, "auth.html", nil)
	}
}

func CompaniesHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		location := vars["location"]
		page := r.URL.Query().Get("p")
		svr.RenderPageForCompanies(w, r, location, page, "companies.html")
	}
}

func DevelopersHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		location := vars["location"]
		tag := vars["tag"]
		page := r.URL.Query().Get("p")
		svr.RenderPageForDevelopers(w, r, location, tag, page, "developers.html")
	}
}

func SubmitDeveloperProfileHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svr.RenderPageForProfileRegistration(w, r, "submit-developer-profile.html")
	}
}

func SubmitRecruiterProfileHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svr.RenderPageForProfileRegistration(w, r, "submit-recruiter-profile.html")
	}
}

func SaveRecruiterProfileHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &struct {
			Fullname   string `json:"fullname"`
			CompanyURL string `json:"company_url"`
			Email      string `json:"email"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			svr.JSON(w, http.StatusBadRequest, "request is invalid")
			return
		}
		if !svr.IsEmail(req.Email) {
			svr.JSON(w, http.StatusBadRequest, "email is invalid")
			return
		}
		for _, e := range []string{"gmail.com", "outlook.com", "live.com", "yahoo.com", "icloud.com"} {
			if strings.Contains(req.Email, e) {
				svr.JSON(w, http.StatusBadRequest, "email must be a valid company email")
				return
			}
		}
		req.Fullname = strings.Title(strings.ToLower(bluemonday.StrictPolicy().Sanitize(req.Fullname)))
		existingRec, err := database.RecruiterProfiles(
			database.RecruiterProfileWhere.Email.EQ(req.Email),
		).OneG(r.Context())
		if err != nil {
			svr.Log(err, "unable to create profile")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		if existingRec.Email == req.Email {
			svr.JSON(w, http.StatusBadRequest, "recruiter profile with this email already exists")
			return
		}
		k, err := ksuid.NewRandom()
		if err != nil {
			svr.Log(err, "unable to generate token")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		t := time.Now().UTC()
		rec := database.RecruiterProfile{
			ID:         k.String(),
			Name:       null.StringFrom(req.Fullname),
			CompanyURL: req.CompanyURL,
			CreatedAt:  t,
			UpdatedAt:  null.TimeFrom(t),
			Email:      strings.ToLower(req.Email),
		}
		err = repositories.SaveTokenSignOn(r.Context(), strings.ToLower(req.Email), k.String(), user.UserTypeRecruiter)
		if err != nil {
			svr.Log(err, "unable to save sign on token")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		err = rec.InsertG(r.Context(), boil.Infer())
		if err != nil {
			svr.Log(err, "unable to save recruiter profile")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		err = svr.GetEmail().SendHTMLEmail(
			email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
			email.Address{Email: req.Email},
			email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
			fmt.Sprintf("Verify Your Recruiter Profile on %s", svr.GetConfig().SiteName),
			fmt.Sprintf(
				"Verify Your Recruiter Profile on %s https://%s/x/auth/%s",
				svr.GetConfig().SiteName,
				svr.GetConfig().SiteHost,
				k.String(),
			),
		)
		if err != nil {
			svr.Log(err, "unable to send email while submitting recruiter profile")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		svr.JSON(w, http.StatusOK, nil)
	}
}

func SaveDeveloperProfileHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &struct {
			Fullname           string   `json:"fullname"`
			HourlyRate         string   `json:"hourly_rate"`
			LinkedinURL        string   `json:"linkedin_url"`
			CurrentLocation    string   `json:"current_location"`
			GithubURL          *string  `json:"github_url,omitempty"`
			TwitterURL         *string  `json:"twitter_url,omitempty"`
			Bio                string   `json:"bio"`
			Tags               string   `json:"tags"`
			ProfileImageID     string   `json:"profile_image_id"`
			Email              string   `json:"email"`
			SearchStatus       string   `json:"search_status"`
			RoleLevel          string   `json:"role_level"`
			RoleTypes          []string `json:"role_types"`
			DetectedLocationID string   `json:"detected_location_id"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			svr.JSON(w, http.StatusBadRequest, "request is invalid")
			return
		}
		if !svr.IsEmail(req.Email) {
			svr.JSON(w, http.StatusBadRequest, "email is invalid")
			return
		}
		linkedinRe := regexp.MustCompile(`^https:\/\/(?:[a-z]{2,3}\.)?linkedin\.com\/.*$`)
		if !linkedinRe.MatchString(req.LinkedinURL) {
			svr.JSON(w, http.StatusBadRequest, "linkedin url is invalid")
			return
		}
		req.Bio = bluemonday.StrictPolicy().Sanitize(req.Bio)
		req.Fullname = strings.Title(strings.ToLower(bluemonday.StrictPolicy().Sanitize(req.Fullname)))
		req.CurrentLocation = strings.Title(strings.ToLower(bluemonday.StrictPolicy().Sanitize(req.CurrentLocation)))
		req.Tags = bluemonday.StrictPolicy().Sanitize(req.Tags)
		if len(strings.Split(req.Tags, ",")) > 10 {
			svr.JSON(w, http.StatusBadRequest, "too many skills")
			return
		}
		if _, ok := developer.ValidSearchStatus[req.SearchStatus]; !ok {
			svr.JSON(w, http.StatusBadRequest, "invalid search status")
			return
		}
		if _, ok := developer.ValidRoleLevels[req.RoleLevel]; !ok {
			svr.JSON(w, http.StatusBadRequest, "invalid role level")
			return
		}
		for _, v := range req.RoleTypes {
			if _, ok := developer.ValidRoleTypes[v]; !ok {
				svr.JSON(w, http.StatusBadRequest, "invalid role type")
				return
			}
		}
		existingDev, err := repositories.DeveloperProfileByEmail(r.Context(), req.Email)
		if err != nil {
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		if existingDev.Email == req.Email {
			svr.JSON(w, http.StatusBadRequest, "developer profile with this email already exists")
			return
		}
		k, err := ksuid.NewRandom()
		if err != nil {
			svr.Log(err, "unable to generate token")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		t := time.Now().UTC()
		// detectedLocationID := &req.DetectedLocationID
		// if req.DetectedLocationID == "" {
		// 	svr.Log(err, "detected location should be set")
		// 	svr.JSON(w, http.StatusBadRequest, "detected_location_id should be set")
		// 	return
		// }
		if req.HourlyRate == "" || req.HourlyRate == "0" {
			svr.JSON(w, http.StatusBadRequest, "Please specify hourly rate")
			return
		}
		hourlyRate, err := strconv.ParseInt(req.HourlyRate, 10, 64)
		if err != nil {
			svr.Log(err, "unable to parse string to int")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		if hourlyRate > 1000 && hourlyRate < 0 {
			svr.Log(err, "Hourly rate cannot be more than 1000 or less than 0")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		dev := &database.DeveloperProfile{
			ID:          k.String(),
			Name:        req.Fullname,
			Location:    req.CurrentLocation,
			HourlyRate:  int(hourlyRate),
			LinkedinURL: req.LinkedinURL,
			GithubURL:   null.StringFromPtr(req.GithubURL),
			TwitterURL:  null.StringFromPtr(req.TwitterURL),
			Bio:         req.Bio,
			Available:   true,
			CreatedAt:   t,
			UpdatedAt:   null.TimeFrom(t),
			Email:       strings.ToLower(req.Email),
			ImageID:     req.ProfileImageID,
			Skills:      req.Tags,
			// SearchStatus:       req.SearchStatus,
			// RoleTypes:          req.RoleTypes,
			// RoleLevel:          req.RoleLevel,
			// DetectedLocationID: detectedLocationID,
		}
		err = repositories.SaveTokenSignOn(r.Context(), strings.ToLower(req.Email), k.String(), user.UserTypeDeveloper)
		if err != nil {
			svr.Log(err, "unable to save sign on token")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		err = repositories.SaveDeveloperProfile(r.Context(), dev)
		if err != nil {
			svr.Log(err, "unable to save developer profile")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		emailSubscriber := database.EmailSubscriber{
			Email: req.Email,
			Token: k.String(),
		}
		err = emailSubscriber.InsertG(r.Context(), boil.Infer())
		if err != nil {
			svr.Log(err, "unable to add email subscriber to db")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		err = svr.GetEmail().SendHTMLEmail(
			email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
			email.Address{Email: req.Email},
			email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
			fmt.Sprintf("Verify Your Developer Profile on %s", svr.GetConfig().SiteName),
			fmt.Sprintf(
				"Verify Your Developer Profile on %s https://%s/x/auth/%s",
				svr.GetConfig().SiteName,
				svr.GetConfig().SiteHost,
				k.String(),
			),
		)
		if err != nil {
			svr.Log(err, "unable to send email while submitting developer profile")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		svr.JSON(w, http.StatusOK, nil)
	}
}

func TriggerFXRateUpdate(svr server.Server) http.HandlerFunc {
	return middleware.MachineAuthenticatedMiddleware(
		svr.GetConfig().MachineToken,
		func(w http.ResponseWriter, r *http.Request) {
			go func() {
				log.Println("going through list of available currencies")
				for _, base := range svr.GetConfig().AvailableCurrencies {
					req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.currencyapi.com/v3/latest?apikey=%s&base_currency=%s", svr.GetConfig().FXAPIKey, base), nil)
					if err != nil {
						svr.Log(err, "http.NewRequest")
						continue
					}
					res, err := http.DefaultClient.Do(req)
					if err != nil {
						svr.Log(err, "http.DefaultClient.Do")
						continue
					}
					var ratesResponse struct {
						Rates map[string]struct {
							Code  string  `json:"code"`
							Value float64 `json:"value"`
						} `json:"data"`
					}
					defer res.Body.Close()
					if err := json.NewDecoder(res.Body).Decode(&ratesResponse); err != nil {
						svr.Log(err, "json.NewDecoder(res.Body).Decode(ratesResponse)")
						continue
					}
					log.Printf("rate response for currency %s: %#v", base, ratesResponse)
					for _, target := range svr.GetConfig().AvailableCurrencies {
						if target == base {
							continue
						}
						cur, ok := ratesResponse.Rates[target]
						if !ok {
							svr.Log(errors.New("could not find target currency"), fmt.Sprintf("could not find target currency %s for base %s", target, base))
							continue
						}
						log.Println("updating fx rate pair ", base, target, cur.Code, cur.Value)
						fx := database.FXRate{
							Base:      base,
							UpdatedAt: time.Now(),
							Target:    target,
							Value:     cur.Value,
						}
						if err := fx.InsertG(r.Context(), boil.Infer()); err != nil {
							svr.Log(err, "database.AddFxRate")
							continue
						}
					}
				}
			}()
		},
	)
}

func TriggerSitemapUpdate(svr server.Server) http.HandlerFunc {
	return middleware.MachineAuthenticatedMiddleware(
		svr.GetConfig().MachineToken,
		func(w http.ResponseWriter, r *http.Request) {
			go func() {
				ctx := r.Context()
				tx, err := svr.Conn.BeginTx(r.Context(), nil)
				if err != nil {
					svr.Log(err, "svr.Conn.BeginTx")
					return
				}
				queries.Raw(`INSERT INTO seo_skill select distinct company from job on conflict do nothing`).ExecContext(r.Context(), tx)
				landingPages, err := seo.GenerateSearchSeoLandingPages(ctx, tx, svr.GetConfig().SiteJobCategory)
				if err != nil {
					svr.Log(err, "seo.GenerateSearchSEOLandingPages")
					return
				}
				fmt.Println("generating post a job landing page")
				postAJobLandingPages, err := seo.GeneratePostAJobSeoLandingPages(ctx, tx, svr.GetConfig().SiteJobCategory)
				if err != nil {
					svr.Log(err, "seo.GeneratePostAJobSEOLandingPages")
					return
				}
				fmt.Println("generating salary landing page")
				salaryLandingPages, err := seo.GenerateSalarySeoLandingPages(ctx, tx, svr.GetConfig().SiteJobCategory)
				if err != nil {
					svr.Log(err, "seo.GenerateSalarySEOLandingPages")
					return
				}
				fmt.Println("generating companies landing page")
				companyLandingPages, err := seo.GenerateCompaniesLandingPages(ctx, tx, svr.GetConfig().SiteJobCategory)
				if err != nil {
					svr.Log(err, "seo.GenerateCompaniesLandingPages")
					return
				}
				fmt.Println("generating dev skill landing pages")
				developerSkillsPages, err := seo.GenerateDevelopersSkillLandingPages(ctx, svr.GetConfig().SiteJobCategory)
				if err != nil {
					svr.Log(err, "seo.GenerateDevelopersSkillLandingPages")
					return
				}
				fmt.Println("generating dev profile landing pages")
				developerProfilePages, err := seo.GenerateDevelopersProfileLandingPages(ctx)
				if err != nil {
					svr.Log(err, "seo.GenerateDevelopersProfileLandingPages")
					return
				}
				fmt.Println("generating company profile landing page")
				companyProfilePages, err := seo.GenerateCompanyProfileLandingPages(ctx)
				if err != nil {
					svr.Log(err, "seo.GenerateDevelopersProfileLandingPages")
					return
				}
				fmt.Println("generating dev location pages")
				developerLocationPages, err := seo.GenerateDevelopersLocationPages(ctx, svr.GetConfig().SiteJobCategory)
				if err != nil {
					svr.Log(err, "seo.GenerateDevelopersLocationPages")
					return
				}
				fmt.Println("generating blog pages")
				blogPosts, err := seo.BlogPages(ctx)
				if err != nil {
					svr.Log(err, "seo.BlogPages")
					return
				}
				fmt.Println("generating static pages")
				pages := seo.StaticPages(svr.GetConfig().SiteJobCategory)
				jobPosts, err := repositories.JobPostByCreatedAt(ctx)
				if err != nil {
					svr.Log(err, "database.JobPostByCreatedAt")
					return
				}
				n := null.TimeFrom(time.Now().UTC())

				_, err = database.Sitemaps().DeleteAll(ctx, tx)
				if err != nil {
					svr.Log(err, "database.Sitemaps().DeleteAll")
					return
				}
				for _, j := range jobPosts {
					fmt.Println("job post page generating...")
					sitemap := database.Sitemap{
						Loc:        fmt.Sprintf(`https://%s/job/%s`, svr.GetConfig().SiteHost, j.Slug.String),
						Changefreq: null.StringFrom("weekly"),
						Lastmod:    null.TimeFrom(j.CreatedAt),
					}
					if err := sitemap.Insert(ctx, tx, boil.Infer()); err != nil {
						svr.Log(err, fmt.Sprintf("database.SaveSitemapEntry: %s", j.Slug.String))
					}
				}

				for _, b := range blogPosts {
					fmt.Println("blog post page generating...")
					sitemap := database.Sitemap{
						Loc:        fmt.Sprintf(`https://%s/blog/%s`, svr.GetConfig().SiteHost, b.Path),
						Changefreq: null.StringFrom("weekly"),
						Lastmod:    n,
					}
					if err := sitemap.Insert(ctx, tx, boil.Infer()); err != nil {
						svr.Log(err, fmt.Sprintf("database.SaveSitemapEntry: %s", b.Path))
					}
				}

				for _, p := range pages {
					fmt.Println("static page generating...")
					sitemap := database.Sitemap{
						Loc:        fmt.Sprintf(`https://%s/%s`, svr.GetConfig().SiteHost, p),
						Changefreq: null.StringFrom("weekly"),
						Lastmod:    n,
					}
					if err := sitemap.Insert(ctx, tx, boil.Infer()); err != nil {
						svr.Log(err, fmt.Sprintf("database.SaveSitemapEntry: %s", p))
					}
				}

				for i, p := range postAJobLandingPages {
					fmt.Println("post a job landing page generating...", i, len(postAJobLandingPages))
					sitemap := database.Sitemap{
						Loc:        fmt.Sprintf(`https://%s/%s`, svr.GetConfig().SiteHost, p),
						Changefreq: null.StringFrom("weekly"),
						Lastmod:    n,
					}
					if err := sitemap.Insert(ctx, tx, boil.Infer()); err != nil {
						svr.Log(err, fmt.Sprintf("database.SaveSitemapEntry: %s", p))
					}
				}

				for i, p := range salaryLandingPages {
					fmt.Println("salary landing page generating...", i, len(salaryLandingPages))
					sitemap := database.Sitemap{
						Loc:        fmt.Sprintf(`https://%s/%s`, svr.GetConfig().SiteHost, p),
						Changefreq: null.StringFrom("weekly"),
						Lastmod:    n,
					}
					if err := sitemap.Insert(ctx, tx, boil.Infer()); err != nil {
						svr.Log(err, fmt.Sprintf("database.SaveSitemapEntry: %s", p))
					}
				}

				for i, p := range landingPages {
					fmt.Println("landing page generating...", i, len(landingPages))
					sitemap := database.Sitemap{
						Loc:        fmt.Sprintf(`https://%s/%s`, svr.GetConfig().SiteHost, p.URI),
						Changefreq: null.StringFrom("weekly"),
						Lastmod:    n,
					}
					if err := sitemap.Insert(ctx, tx, boil.Infer()); err != nil {
						svr.Log(err, fmt.Sprintf("database.SaveSitemapEntry: %s", p.URI))
					}
				}

				for _, p := range companyLandingPages {
					sitemap := database.Sitemap{
						Loc:        fmt.Sprintf(`https://%s/%s`, svr.GetConfig().SiteHost, p),
						Changefreq: null.StringFrom("weekly"),
						Lastmod:    n,
					}
					if err := sitemap.Insert(ctx, tx, boil.Infer()); err != nil {
						svr.Log(err, fmt.Sprintf("database.SaveSitemapEntry: %s", p))
					}
				}

				for _, p := range developerSkillsPages {
					sitemap := database.Sitemap{
						Loc:        fmt.Sprintf(`https://%s/%s`, svr.GetConfig().SiteHost, p),
						Changefreq: null.StringFrom("weekly"),
						Lastmod:    n,
					}
					if err := sitemap.Insert(ctx, tx, boil.Infer()); err != nil {
						svr.Log(err, fmt.Sprintf("database.SaveSitemapEntry: %s", p))
					}
				}

				for _, p := range developerProfilePages {
					sitemap := database.Sitemap{
						Loc:        fmt.Sprintf(`https://%s/%s`, svr.GetConfig().SiteHost, p),
						Changefreq: null.StringFrom("weekly"),
						Lastmod:    n,
					}
					if err := sitemap.Insert(ctx, tx, boil.Infer()); err != nil {
						svr.Log(err, fmt.Sprintf("database.SaveSitemapEntry: %s", p))
					}
				}
				for _, p := range companyProfilePages {
					sitemap := database.Sitemap{
						Loc:        fmt.Sprintf(`https://%s/%s`, svr.GetConfig().SiteHost, p),
						Changefreq: null.StringFrom("weekly"),
						Lastmod:    n,
					}
					if err := sitemap.Insert(ctx, tx, boil.Infer()); err != nil {
						svr.Log(err, fmt.Sprintf("database.SaveSitemapEntry: %s", p))
					}
				}

				for _, p := range developerLocationPages {
					sitemap := database.Sitemap{
						Loc:        fmt.Sprintf(`https://%s/%s`, svr.GetConfig().SiteHost, p),
						Changefreq: null.StringFrom("weekly"),
						Lastmod:    n,
					}
					if err := sitemap.Insert(ctx, tx, boil.Infer()); err != nil {
						svr.Log(err, fmt.Sprintf("database.SaveSitemapEntry: %s", p))
					}
				}
			}()
		})
}

func TriggerExpiredJobsTask(svr server.Server) http.HandlerFunc {
	return middleware.MachineAuthenticatedMiddleware(
		svr.GetConfig().MachineToken,
		func(w http.ResponseWriter, r *http.Request) {
			go func() {
				jobs, err := repositories.GetJobApplyURLs(r.Context())
				if err != nil {
					svr.Log(err, "unable to get job apply URL for cleanup")
					return
				}
				for _, job := range jobs {
					if svr.IsEmail(job.ApplicationLink) {
						continue
					}
					res, err := http.Get(job.ApplicationLink)
					if err != nil {
						svr.Log(err, fmt.Sprintf("error while checking expired apply URL for job %s %s", job.ID, job.ApplicationLink))
						continue
					}
					if res.StatusCode == http.StatusNotFound {
						fmt.Printf("found expired job %s URL %s returned 404\n", job.ID, job.ApplicationLink)
						if err := repositories.MarkJobAsExpired(r.Context(), job); err != nil {
							svr.Log(err, fmt.Sprintf("unable to mark job %s %s as expired", job.ID, job.ApplicationLink))
						}
					}
				}
			}()
			svr.JSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
		},
	)
}

func TriggerUpdateLastWeekClickouts(svr server.Server) http.HandlerFunc {
	return middleware.MachineAuthenticatedMiddleware(
		svr.GetConfig().MachineToken,
		func(w http.ResponseWriter, r *http.Request) {
			go func() {
				err := UpdateLastWeekClickouts(r.Context(), svr.Conn)
				if err != nil {
					svr.Log(err, "unable to update last week clickouts")
					return
				}
			}()
			svr.JSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
		},
	)
}

func UpdateLastWeekClickouts(ctx context.Context, conn *sql.DB) error {
	_, err := conn.Exec(`WITH cte AS (SELECT job_id, count(*) AS clickouts FROM job_event WHERE event_type = 'clickout' AND created_at > CURRENT_DATE - 7 GROUP BY job_id)
	UPDATE job SET last_week_clickouts = cte.clickouts FROM cte WHERE cte.job_id = id`)
	return err
}

func TriggerCloudflareStatsExport(svr server.Server) http.HandlerFunc {
	return middleware.MachineAuthenticatedMiddleware(
		svr.GetConfig().MachineToken,
		func(w http.ResponseWriter, r *http.Request) {
			go func() {
				ctx := r.Context()
				client := graphql.NewClient(svr.GetConfig().CloudflareAPIEndpoint)
				req := graphql.NewRequest(
					`query {
  viewer {
    zones(filter: {zoneTag: $zoneTag}) {
      httpRequests1dGroups(orderBy: [date_ASC]  filter: { date_gt: $fromDate } limit: 10000) {
        dimensions {
          date
        }
      sum {
        pageViews
        requests
        bytes
        cachedBytes
        threats
        countryMap {
          clientCountryName
          requests
          threats
        }
	browserMap {
          uaBrowserFamily
          pageViews
        }
        responseStatusMap {
          edgeResponseStatus
          requests
        }
      }
        uniq {
          uniques
        }
    }
  }
}
}`,
				)
				var err error
				var daysAgo int
				daysAgoStr := r.URL.Query().Get("days_ago")
				daysAgo, err = strconv.Atoi(daysAgoStr)
				if err != nil {
					daysAgo = 3
				}
				req.Var("zoneTag", svr.GetConfig().CloudflareZoneTag)
				req.Var("fromDate", time.Now().UTC().AddDate(0, 0, -daysAgo).Format("2006-01-02"))
				req.Header.Set("Cache-Control", "no-cache")
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", svr.GetConfig().CloudflareAPIToken))
				type cloudFlareStatsResponse struct {
					Viewer struct {
						Zones []struct {
							HttpRequests1dGroups []struct {
								Dimensions struct {
									Date string `json:"date"`
								} `json:"dimensions"`
								Sum struct {
									Bytes       int64 `json:"bytes"`
									CachedBytes int64 `json:"cachedBytes"`
									CountryMap  []struct {
										ClientCountryName string `json:"clientCountryName"`
										Requests          int64  `json:"requests"`
										Threats           int64  `json:"threats"`
									} `json:"countryMap"`
									BrowserMap []struct {
										UABrowserFamily string `json:"uaBrowserFamily"`
										PageViews       int64  `json:"pageViews"`
									} `json:"browserMap"`
									PageViews         int64 `json:"pageViews"`
									Requests          int64 `json:"requests"`
									ResponseStatusMap []struct {
										EdgeResponseStatus int   `json:"edgeResponseStatus"`
										Requests           int64 `json:"requests"`
									} `json:"responseStatusMap"`
									Threats int64 `json:"threats"`
								} `json:"sum"`
								Uniq struct {
									Uniques int64 `json:"uniques"`
								} `json:"uniq"`
							} `json:"httpRequests1dGroups"`
						} `json:"zones"`
					} `json:"viewer"`
				}
				var res cloudFlareStatsResponse
				if err := client.Run(context.Background(), req, &res); err != nil {
					svr.Log(err, "unable to complete graphql request to cloudflare")
					return
				}
				stat := database.CloudflareStat{}
				statusCodeStat := database.CloudflareStatusCodeStat{}
				countryStat := database.CloudflareCountryStat{}
				browserStat := database.CloudflareBrowserStat{}
				if len(res.Viewer.Zones) < 1 {
					svr.Log(errors.New("got empty response from cloudflare APIs"), "expecting 1 zone got none")
					return
				}
				log.Printf("retrieved %d cloudflare stat entries\n", len(res.Viewer.Zones[0].HttpRequests1dGroups))
				for _, d := range res.Viewer.Zones[0].HttpRequests1dGroups {
					stat.Date, err = time.Parse("2006-01-02", d.Dimensions.Date)
					if err != nil {
						svr.Log(err, "unable to parse date from cloudflare stat")
						return
					}
					stat.Bytes = d.Sum.Bytes
					stat.CachedBytes = d.Sum.CachedBytes
					stat.PageViews = d.Sum.PageViews
					stat.Requests = d.Sum.Requests
					stat.Threats = d.Sum.Threats
					stat.Uniques = d.Uniq.Uniques
					if err := stat.InsertG(ctx, boil.Infer()); err != nil {
						svr.Log(err, "database.SaveCloudflareStat")
						return
					}
					// status code stat
					for _, v := range d.Sum.ResponseStatusMap {
						statusCodeStat.Date = stat.Date
						statusCodeStat.StatusCode = v.EdgeResponseStatus
						statusCodeStat.Requests = v.Requests
						if err := statusCodeStat.InsertG(ctx, boil.Infer()); err != nil {
							svr.Log(err, "database.SaveCloudflareStatusCodeStat")
							return
						}
					}
					// country stat
					for _, v := range d.Sum.CountryMap {
						countryStat.Date = stat.Date
						countryStat.CountryCode = v.ClientCountryName
						countryStat.Requests = v.Requests
						countryStat.Threats = v.Threats
						if err := countryStat.InsertG(ctx, boil.Infer()); err != nil {
							svr.Log(err, "database.SaveCloudflareCountryStat")
							return
						}
					}
					// browser stat
					for _, v := range d.Sum.BrowserMap {
						browserStat.Date = stat.Date
						browserStat.PageViews = v.PageViews
						browserStat.UaBrowserFamily = v.UABrowserFamily
						if err := browserStat.InsertG(ctx, boil.Infer()); err != nil {
							svr.Log(err, "database.SaveCloudflareBrowserStat")
							return
						}
					}
				}
				log.Println("done exporting cloudflare stats")
			}()
			svr.JSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
		},
	)
}

func TriggerWeeklyNewsletter(svr server.Server) http.HandlerFunc {
	return middleware.MachineAuthenticatedMiddleware(
		svr.GetConfig().MachineToken,
		func(w http.ResponseWriter, r *http.Request) {
			go func() {
				lastJobID, err := repositories.GetValue(r.Context(), "last_sent_job_id_weekly")
				if err != nil {
					svr.Log(err, "unable to retrieve last newsletter weekly job id")
					return
				}
				jobPosts, err := repositories.GetLastNJobsFromID(r.Context(), svr.GetConfig().NewsletterJobsToSend, lastJobID.Value)
				if err != nil {
					svr.Log(err, "unable to retrieve last newsletter weekly job id")
					return
				}
				if len(jobPosts) < 1 {
					log.Printf("found 0 new jobs for weekly newsletter. quitting")
					return
				}
				fmt.Printf("found %d/%d jobs for weekly newsletter\n", len(jobPosts), svr.GetConfig().NewsletterJobsToSend)
				subscribers, err := database.EmailSubscribers(
					database.EmailSubscriberWhere.ConfirmedAt.IsNotNull(),
				).AllG(r.Context())
				if err != nil {
					svr.Log(err, "unable to retrieve subscribers")
					return
				}
				var jobsHTMLArr []string
				for _, j := range jobPosts {
					jobsHTMLArr = append(jobsHTMLArr, `<p><b>Job Title:</b> `+j.JobTitle+`<br /><b>Company:</b> `+j.Company+`<br /><b>Location:</b> `+j.Location+`<br /><b>Salary:</b> `+j.SalaryRange+`<br /><b>Detail:</b> <a href="https://`+svr.GetConfig().SiteHost+`/job/`+j.Slug.String+`">https://`+svr.GetConfig().SiteHost+`/job/`+j.Slug.String+`</a></p>`)
					lastJobID.Value = j.ID
				}
				jobsHTML := strings.Join(jobsHTMLArr, " ")
				campaignContentHTML := `<p>Here's a list of the newest ` + fmt.Sprintf("%d", len(jobPosts)) + ` ` + svr.GetConfig().SiteJobCategory + ` jobs this week on ` + svr.GetConfig().SiteName + `</p>
` + jobsHTML + `
	<p>Check out more jobs at <a title="` + svr.GetConfig().SiteName + `" href="https://` + svr.GetConfig().SiteHost + `">https://` + svr.GetConfig().SiteHost + `</a></p>
	<p>Get companies apply to you, join the ` + svr.GetConfig().SiteJobCategory + ` Developer Community <a title="` + svr.GetConfig().SiteName + ` Community" href="https://` + svr.GetConfig().SiteHost + `/Join-` + strings.Title(svr.GetConfig().SiteJobCategory) + `-Community">https://` + svr.GetConfig().SiteHost + `/Join-` + strings.Title(svr.GetConfig().SiteJobCategory) + `-Community</a></p>
	<p>` + svr.GetConfig().SiteName + `</p>
	<hr />`
				unsubscribeLink := `
	<h6><strong> ` + svr.GetConfig().SiteName + `</strong> | London, United Kingdom<br />This email was sent to <strong>%s</strong> | <a href="https://` + svr.GetConfig().SiteHost + `/x/email/unsubscribe?token=%s">Unsubscribe</a></h6>`

				for _, s := range subscribers {
					err = svr.GetEmail().SendHTMLEmail(
						email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
						email.Address{Email: s.Email},
						email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
						fmt.Sprintf("Go Jobs This Week (%d New)", len(jobPosts)),
						campaignContentHTML+fmt.Sprintf(unsubscribeLink, s.Email, s.Token),
					)
					if err != nil {
						svr.Log(err, fmt.Sprintf("unable to send email for newsletter email %s", s.Email))
						continue
					}
				}
				_, err = lastJobID.UpdateG(r.Context(), boil.Whitelist("value"))
				if err != nil {
					svr.Log(err, "unable to save last weekly newsletter job id to db")
					return
				}
			}()
			svr.JSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
		},
	)
}

func TriggerTelegramScheduler(svr server.Server) http.HandlerFunc {
	return middleware.MachineAuthenticatedMiddleware(
		svr.GetConfig().MachineToken,
		func(w http.ResponseWriter, r *http.Request) {
			go func() {
				ctx := r.Context()
				lastTelegramJobID, err := repositories.GetValue(ctx, "last_telegram_job_id")
				if err != nil {
					svr.Log(err, "unable to retrieve last telegram job id")
					return
				}
				jobPosts, err := repositories.GetLastNJobsFromID(r.Context(), svr.GetConfig().TwitterJobsToPost, lastTelegramJobID.Value)
				if err != nil {
					svr.Log(err, "unable to retrieve last telegram job id")
					return
				}
				log.Printf("found %d/%d jobs to post on telegram\n", len(jobPosts), svr.GetConfig().TwitterJobsToPost)
				if len(jobPosts) == 0 {
					return
				}
				api := telegram.New(svr.GetConfig().TelegramAPIToken)
				bgCtx := context.Background()
				for _, j := range jobPosts {
					_, err := api.SendMessage(bgCtx, telegram.NewMessage(svr.GetConfig().TelegramChannelID, fmt.Sprintf("%s with %s - %s | %s\n\n#%s #%sjobs\n\nhttps://%s/job/%s", j.JobTitle, j.Company, j.Location, j.SalaryRange, svr.GetConfig().SiteJobCategory, svr.GetConfig().SiteJobCategory, svr.GetConfig().SiteHost, j.Slug.String)))
					if err != nil {
						svr.Log(err, "unable to post on telegram")
						continue
					}
					lastTelegramJobID.Value = j.ID
				}
				_, err = lastTelegramJobID.UpdateG(bgCtx, boil.Whitelist("value"))
				if err != nil {
					svr.Log(err, fmt.Sprintf("unable to save last telegram job id to db as %s", lastTelegramJobID.Value))
					return
				}
				log.Printf("updated last telegram job id to %s\n", lastTelegramJobID.Value)
				log.Printf("posted last %d jobs to telegram", len(jobPosts))
			}()
			svr.JSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
		},
	)
}

func TriggerMonthlyHighlights(svr server.Server) http.HandlerFunc {
	return middleware.MachineAuthenticatedMiddleware(
		svr.GetConfig().MachineToken,
		func(w http.ResponseWriter, r *http.Request) {
			go func() {
				pageviewsLast30Days, err := GetWebsitePageViewsLast30Days(svr.Conn)
				if err != nil {
					svr.Log(err, "could not retrieve pageviews for last 30 days")
					return
				}
				jobPageviewsLast30Days, err := GetJobPageViewsLast30Days(svr.Conn)
				if err != nil {
					svr.Log(err, "could not retrieve job pageviews for last 30 days")
					return
				}
				jobApplicantsLast30Days, err := GetJobClickoutsLast30Days(svr.Conn)
				if err != nil {
					svr.Log(err, "could not retrieve job clickouts for last 30 days")
					return
				}
				_, newJobsLastMonth, err := repositories.NewJobsLastWeekOrMonth(r.Context())
				if err != nil {
					svr.Log(err, "unable to retrieve new jobs last week last month")
					return
				}
				pageviewsLast30DaysText := humanize.Comma(int64(pageviewsLast30Days))
				jobPageviewsLast30DaysText := humanize.Comma(int64(jobPageviewsLast30Days))
				jobApplicantsLast30DaysText := humanize.Comma(int64(jobApplicantsLast30Days))
				newJobsLastMonthText := humanize.Comma(int64(newJobsLastMonth))
				highlights := fmt.Sprintf(`This months highlight ✨

📣 %s new jobs posted last month
✉️  %s applicants last month
🌎 %s pageviews last month
💼 %s jobs viewed last month
`, newJobsLastMonthText, jobApplicantsLast30DaysText, pageviewsLast30DaysText, jobPageviewsLast30DaysText)
				err = svr.GetEmail().SendHTMLEmail(
					email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
					email.Address{Email: svr.GetEmail().DefaultAdminAddress()},
					email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
					fmt.Sprintf("%s Monthly Highlights", svr.GetConfig().SiteName),
					highlights,
				)
				if err != nil {
					svr.Log(err, "unable to send monthtly highlights email")
					return
				}
			}()
		},
	)
}

func GetWebsitePageViewsLast30Days(conn *sql.DB) (int, error) {
	var c int
	row := conn.QueryRow(`SELECT SUM(page_views) AS c FROM cloudflare_browser_stats WHERE date > CURRENT_DATE - 30 AND ua_browser_family NOT ILIKE '%bot%'`)
	if err := row.Scan(&c); err != nil {
		return 100000, nil
	}

	return c, nil
}

func GetJobPageViewsLast30Days(conn *sql.DB) (int, error) {
	var c int
	row := conn.QueryRow(`SELECT COUNT(*) AS c FROM job_event WHERE event_type = 'page_view' AND created_at > CURRENT_DATE - 30`)
	if err := row.Scan(&c); err != nil {
		return 100000, nil
	}

	return c, nil
}

func GetJobClickoutsLast30Days(conn *sql.DB) (int, error) {
	var c int
	row := conn.QueryRow(`SELECT COUNT(*) AS c FROM job_event WHERE event_type = 'clickout' AND created_at > CURRENT_DATE - 30`)
	if err := row.Scan(&c); err != nil {
		return 100000, nil
	}

	return c, nil
}

func TriggerTwitterScheduler(svr server.Server) http.HandlerFunc {
	return middleware.MachineAuthenticatedMiddleware(
		svr.GetConfig().MachineToken,
		func(w http.ResponseWriter, r *http.Request) {
			go func() {
				lastTwittedJobID, err := repositories.GetValue(r.Context(), "last_twitted_job_id")
				if err != nil {
					svr.Log(err, "unable to retrieve last twitter job id")
					return
				}
				jobPosts, err := repositories.GetLastNJobsFromID(r.Context(), svr.GetConfig().TwitterJobsToPost, lastTwittedJobID.Value)
				if err != nil {
					svr.Log(err, "unable to retrieve last n jobs from id")
					return
				}
				log.Printf("found %d/%d jobs to post on twitter\n", len(jobPosts), svr.GetConfig().TwitterJobsToPost)
				if len(jobPosts) == 0 {
					return
				}
				lastJobID := lastTwittedJobID.Value
				api := anaconda.NewTwitterApiWithCredentials(svr.GetConfig().TwitterAccessToken, svr.GetConfig().TwitterAccessTokenSecret, svr.GetConfig().TwitterClientKey, svr.GetConfig().TwitterClientSecret)
				for _, j := range jobPosts {
					_, err := api.PostTweet(fmt.Sprintf("%s with %s - %s | %s\n\n#%s #%sjobs\n\nhttps://%s/job/%s", j.JobTitle, j.Company, j.Location, j.SalaryRange, svr.GetConfig().SiteJobCategory, svr.GetConfig().SiteJobCategory, svr.GetConfig().SiteHost, j.Slug.String), url.Values{})
					if err != nil {
						svr.Log(err, "unable to post tweet")
						continue
					}
					lastJobID = j.ID
				}
				err = repositories.SetValue(r.Context(), "last_twitted_job_id", lastJobID)
				if err != nil {
					svr.Log(err, fmt.Sprintf("unable to save last twitter job id to db as %s", lastJobID))
					return
				}
				log.Printf("updated last twitted job id to %s\n", lastJobID)
				log.Printf("posted last %d jobs to twitter", len(jobPosts))
			}()
			svr.JSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
		},
	)
}

/*
func TriggerCompanyUpdate(svr server.Server) http.HandlerFunc {
	return middleware.MachineAuthenticatedMiddleware(
		svr.GetConfig().MachineToken,
		func(w http.ResponseWriter, r *http.Request) {
			go func() {
				since := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
				cs, err := repositories.InferCompaniesFromJobs(since)
				if err != nil {
					svr.Log(err, "unable to infer companies from jobs")
					return
				}
				log.Printf("inferred %d companies...\n", len(cs))
				for _, c := range cs {
					res, err := http.Get(c.URL)
					if err != nil {
						svr.Log(err, fmt.Sprintf("http.Get(%s): unable to get url", c.URL))
						continue
					}
					defer res.Body.Close()
					if res.StatusCode != http.StatusOK {
						svr.Log(errors.New("non 200 status code"), fmt.Sprintf("GET %s: status code error: %d %s", c.URL, res.StatusCode, res.Status))
						continue
					}

					doc, err := goquery.NewDocumentFromReader(res.Body)
					if err != nil {
						svr.Log(err, "goquery.NewDocumentFromReader")
						continue
					}
					description := doc.Find("title").Text()
					twitter := ""
					doc.Find("meta").Each(func(i int, s *goquery.Selection) {
						if name, _ := s.Attr("name"); strings.EqualFold(name, "description") {
							var ok bool
							desc, ok := s.Attr("content")
							if !ok {
								log.Println("unable to retrieve content for description tag for companyURL ", c.URL)
								return
							}
							if desc != "" {
								description = desc
							}
							log.Printf("description: %s\n", description)
						}
						if name, _ := s.Attr("name"); strings.EqualFold(name, "twitter:site") {
							var ok bool
							twtr, ok := s.Attr("content")
							if !ok {
								log.Println("unable to retrieve content for twitter:site")
								return
							}
							if twtr != "" {
								twitter = "https://twitter.com/" + strings.Trim(twtr, "@")
							}
							log.Printf("twitter: %s\n", twitter)
						}
					})
					github := ""
					linkedin := ""
					doc.Find("a").Each(func(i int, s *goquery.Selection) {
						if href, ok := s.Attr("href"); ok && strings.Contains(href, "github.com/") {
							github = href
							log.Printf("github: %s\n", github)
						}
						if href, ok := s.Attr("href"); ok && strings.Contains(href, "linkedin.com/") {
							linkedin = href
							log.Printf("linkedin: %s\n", linkedin)
						}
						if twitter == "" {
							if href, ok := s.Attr("href"); ok && strings.Contains(href, "twitter.com/") {
								twitter = href
								log.Printf("twitter: %s\n", twitter)
							}
						}
					})
					if description != "" {
						c.Description = &description
					}
					if twitter != "" {
						c.Twitter = &twitter
					}
					if github != "" {
						c.Github = &github
					}
					if linkedin != "" {
						c.Linkedin = &linkedin
					}
					companyID, err := ksuid.NewRandom()
					if err != nil {
						svr.Log(err, "ksuid.NewRandom: companyID")
						continue
					}
					newIconID, err := ksuid.NewRandom()
					if err != nil {
						svr.Log(err, "ksuid.NewRandom: newIconID")
						continue
					}
					if err := database.DuplicateImage(svr.Conn, c.IconImageID, newIconID.String()); err != nil {
						svr.Log(err, "database.DuplicateImage")
						continue
					}
					c.ID = companyID.String()
					c.Slug = slug.Make(c.Name)
					c.IconImageID = newIconID.String()
					if err := repositories.SaveCompany(c); err != nil {
						svr.Log(err, "repositories.SaveCompany")
						continue
					}
					log.Println(c.Name)
				}
				if err := repositories.DeleteStaleImages(svr.GetConfig().SiteLogoImageID); err != nil {
					svr.Log(err, "repositories.DeleteStaleImages")
					return
				}
			}()
			svr.JSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
		},
	)
}
*/

func TriggerAdsManager(svr server.Server) http.HandlerFunc {
	return middleware.MachineAuthenticatedMiddleware(
		svr.GetConfig().MachineToken,
		func(w http.ResponseWriter, r *http.Request) {
			// TODO: add column to jobs and send emails
			// get jobs if plan_expired_at is less than now and last_email_sent is between 2 weeks ago and now
			svr.JSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
		},
	)
}

func UpdateDeveloperProfileHandler(svr server.Server) http.HandlerFunc {
	return middleware.UserAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			req := &struct {
				ID                 string   `json:"id"`
				Fullname           string   `json:"fullname"`
				HourlyRate         string   `json:"hourly_rate"`
				LinkedinURL        string   `json:"linkedin_url"`
				Bio                string   `json:"bio"`
				CurrentLocation    string   `json:"current_location"`
				Skills             string   `json:"skills"`
				ImageID            string   `json:"profile_image_id"`
				Email              string   `json:"email"`
				SearchStatus       string   `json:"search_status"`
				RoleLevel          string   `json:"role_level"`
				RoleTypes          []string `json:"role_types"`
				DetectedLocationID string   `json:"detected_location_id"`
			}{}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				svr.Log(errors.New("invalid search status"), "invalid search status")
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			if !svr.IsEmail(req.Email) {
				svr.Log(errors.New("invalid search status"), "invalid search status")
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			linkedinRe := regexp.MustCompile(`^https:\/\/(?:[a-z]{2,3}\.)?linkedin\.com\/.*$`)
			if !linkedinRe.MatchString(req.LinkedinURL) {
				svr.Log(errors.New("invalid search status"), "invalid search status")
				svr.JSON(w, http.StatusBadRequest, "linkedin url is invalid")
				return
			}
			if _, ok := developer.ValidSearchStatus[req.SearchStatus]; !ok {
				svr.Log(errors.New("invalid search status"), "invalid search status")
				svr.JSON(w, http.StatusBadRequest, "invalid search status")
				return
			}
			if _, ok := developer.ValidRoleLevels[req.RoleLevel]; !ok {
				svr.Log(errors.New("invalid role level"), "invalid role level")
				svr.JSON(w, http.StatusBadRequest, "invalid role level")
				return
			}
			for _, v := range req.RoleTypes {
				if _, ok := developer.ValidRoleTypes[v]; !ok {
					svr.Log(errors.New("invalid role type"), "invalid role type")
					svr.JSON(w, http.StatusBadRequest, "invalid role type")
					return
				}
			}
			if req.HourlyRate == "" || req.HourlyRate == "0" {
				svr.JSON(w, http.StatusBadRequest, "please specify hourly rate")
				return
			}
			req.Bio = bluemonday.StrictPolicy().Sanitize(req.Bio)
			req.Fullname = strings.Title(strings.ToLower(bluemonday.StrictPolicy().Sanitize(req.Fullname)))
			req.CurrentLocation = strings.Title(strings.ToLower(bluemonday.StrictPolicy().Sanitize(req.CurrentLocation)))
			req.Skills = bluemonday.StrictPolicy().Sanitize(req.Skills)
			if len(strings.Split(req.Skills, ",")) > 10 {
				svr.JSON(w, http.StatusBadRequest, "too many skills")
				return
			}
			profile, err := middleware.GetUserFromJWT(r, svr.SessionStore, svr.GetJWTSigningKey())
			if err != nil {
				svr.Log(err, "unable to get email from JWT")
				svr.JSON(w, http.StatusForbidden, nil)
				return
			}
			if req.Email != profile.Email && !profile.IsAdmin {
				svr.JSON(w, http.StatusForbidden, nil)
				return
			}
			t := time.Now().UTC()
			avail := true
			if req.SearchStatus == developer.SearchStatusNotAvailable {
				avail = false
			}
			hourlyRate, err := strconv.ParseInt(req.HourlyRate, 10, 64)
			if err != nil {
				svr.Log(err, "unable to parse string to int")
				svr.JSON(w, http.StatusInternalServerError, nil)
				return
			}

			dev := &database.DeveloperProfile{
				ID:          req.ID,
				Name:        req.Fullname,
				Location:    req.CurrentLocation,
				HourlyRate:  int(hourlyRate),
				LinkedinURL: req.LinkedinURL,
				Bio:         req.Bio,
				Email:       req.Email,
				Available:   avail,
				UpdatedAt:   null.TimeFrom(t),
				Skills:      req.Skills,
				ImageID:     req.ImageID,
				// SearchStatus: req.SearchStatus,
				// RoleLevel:    req.RoleLevel,
			}
			err = repositories.UpdateDeveloperProfile(r.Context(), dev)
			if err != nil {
				svr.Log(err, "unable to update developer profile")
				svr.JSON(w, http.StatusInternalServerError, nil)
				return
			}
			svr.JSON(w, http.StatusOK, nil)
		},
	)
}

func DeleteDeveloperProfileHandler(svr server.Server) http.HandlerFunc {
	return middleware.UserAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			req := &struct {
				ID      string `json:"id"`
				ImageID string `json:"image_id"`
				Email   string `json:"email"`
			}{}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			if !svr.IsEmail(req.Email) {
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			profile, err := middleware.GetUserFromJWT(r, svr.SessionStore, svr.GetJWTSigningKey())
			if err != nil {
				svr.Log(err, "unable to get email from JWT")
				svr.JSON(w, http.StatusForbidden, nil)
				return
			}
			if profile.Email != req.Email && !profile.IsAdmin {
				svr.JSON(w, http.StatusForbidden, nil)
				return
			}
			err = repositories.DeleteDeveloperProfile(r.Context(), req.ID, req.Email)
			if err != nil {
				svr.Log(err, "unable to delete developer profile")
				svr.JSON(w, http.StatusInternalServerError, nil)
				return
			}
			if _, imageErr := database.Images(
				database.ImageWhere.ID.EQ(req.ImageID),
			).DeleteAllG(r.Context()); imageErr != nil {
				svr.Log(err, "unable to delete developer profile image id "+req.ImageID)
				svr.JSON(w, http.StatusInternalServerError, nil)
				return
			}
			if _, userErr := database.Users(
				database.UserWhere.Email.EQ(req.Email),
			).DeleteAllG(r.Context()); userErr != nil {
				svr.Log(err, "unable to delete user by email "+req.Email)
				svr.JSON(w, http.StatusInternalServerError, nil)
			}
			svr.JSON(w, http.StatusOK, nil)
		},
	)
}

func ConfirmEmailSubscriberHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		token := vars["token"]
		_, err := database.EmailSubscribers(
			database.EmailSubscriberWhere.Token.EQ(token),
		).UpdateAllG(
			r.Context(),
			database.M{
				database.EmailSubscriberColumns.ConfirmedAt: time.Now(),
			},
		)
		if err != nil {
			svr.Log(err, "unable to confirm subscriber using token "+token)
			svr.TEXT(w, http.StatusInternalServerError, "There was an error with your request. Please try again later.")
			return
		}
		svr.TEXT(w, http.StatusOK, "Your email subscription has been confirmed successfully.")
	}
}

func RemoveEmailSubscriberHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := database.EmailSubscribers(
			database.EmailSubscriberWhere.Token.EQ(r.URL.Query().Get("token")),
		).DeleteAllG(r.Context())
		if err != nil {
			svr.Log(err, "unable to add email subscriber to db")
			svr.TEXT(w, http.StatusInternalServerError, "")
			return
		}
		svr.TEXT(w, http.StatusOK, "Your email has been successfully removed.")
	}
}

func AddEmailSubscriberHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		emailStr := strings.ToLower(r.URL.Query().Get("email"))
		if !svr.IsEmail(emailStr) {
			svr.Log(errors.New("invalid email"), "request email is not a valid email")
			svr.JSON(w, http.StatusBadRequest, "invalid email provided")
			return
		}
		k, err := ksuid.NewRandom()
		if err != nil {
			svr.Log(err, "unable to generate email subscriber token")
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		subscriber := &database.EmailSubscriber{
			Email:     emailStr,
			Token:     k.String(),
			CreatedAt: time.Now(),
		}
		err = subscriber.InsertG(r.Context(), boil.Infer())
		if err != nil {
			svr.Log(err, "unable to add email subscriber to db")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		err = svr.GetEmail().SendHTMLEmail(
			email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
			email.Address{Email: emailStr},
			email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
			fmt.Sprintf("Confirm Your Email Subscription on %s", svr.GetConfig().SiteName),
			fmt.Sprintf(
				"Please click on the link below to confirm your subscription to receive weekly emails from %s\n\n%s\n\nIf this was not requested by you, please ignore this email.",
				svr.GetConfig().SiteName,
				fmt.Sprintf("https://%s/x/email/confirm/%s", svr.GetConfig().SiteHost, k.String()),
			),
		)
		if err != nil {
			svr.Log(err, "unable to send email while submitting message")
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		svr.JSON(w, http.StatusOK, nil)
	}
}

func SendMessageDeveloperProfileHandler(svr server.Server) http.HandlerFunc {
	return middleware.UserAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			profileID := vars["id"]
			sender, err := middleware.GetUserFromJWT(r, svr.SessionStore, svr.GetJWTSigningKey())
			if err != nil {
				svr.Log(err, "unable to get email from JWT")
				svr.JSON(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			req := &struct {
				Content string `json:"content"`
				Email   string `json:"email"`
			}{}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				reqData, ioErr := ioutil.ReadAll(r.Body)
				if ioErr != nil {
					svr.Log(ioErr, "unable to read request body data for developer profile message")
				}
				svr.Log(err, fmt.Sprintf("unable to decode request body from developer profile message %+v", string(reqData)))
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			if !svr.IsEmail(req.Email) {
				svr.Log(errors.New("invalid email"), "request email is not a valid email")
				svr.JSON(w, http.StatusBadRequest, "invalid email provided")
				return
			}
			dev, err := repositories.DeveloperProfileByID(r.Context(), profileID)
			if err != nil {
				svr.Log(err, "unable to find developer profile by id "+profileID)
				svr.JSON(w, http.StatusInternalServerError, nil)
				return
			}
			k, err := ksuid.NewRandom()
			if err != nil {
				svr.Log(err, "unable to generate message ID")
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			message := &database.DeveloperProfileMessage{
				ID:        k.String(),
				Email:     req.Email,
				Content:   req.Content,
				ProfileID: dev.ID,
				SenderID:  sender.UserID,
			}
			err = message.InsertG(r.Context(), boil.Infer())
			if err != nil {
				svr.Log(err, "unable to send message to developer profile")
				svr.JSON(w, http.StatusInternalServerError, nil)
				return
			}
			if err := repositories.TrackDeveloperProfileMessageSent(r.Context(), dev); err != nil {
				svr.Log(err, "unable to track message sent to developer profile")
			}
			err = svr.GetEmail().SendHTMLEmail(
				email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
				email.Address{Email: dev.Email},
				email.Address{Email: message.Email},
				fmt.Sprintf("New Message from %s", svr.GetConfig().SiteName),
				fmt.Sprintf(
					"You received a new message from %s: \n\nMessage: %s\n\nFrom: %s",
					svr.GetConfig().SiteName,
					message.Content,
					message.Email,
				),
			)
			if err != nil {
				svr.Log(err, "unable to send email to developer profile")
				svr.JSON(w, http.StatusBadRequest, "There was a problem while sending the email")
				return
			}
			if err := repositories.MarkDeveloperMessageAsSent(r.Context(), message.ID); err != nil {
				svr.Log(err, "unable to mark developer message as sent "+message.ID)
			}
			svr.JSON(w, http.StatusOK, nil)
		})
}

func AutocompleteLocation(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prefix := r.URL.Query().Get("k")
		locs, err := database.SeoLocations(
			qm.Where("name ILIKE ?", fmt.Sprintf("%s%%", prefix)),
			qm.OrderBy(database.SeoLocationColumns.Population+" DESC"),
			qm.Limit(5),
		).AllG(r.Context())
		if err != nil {
			svr.Log(err, "unable to retrieve locations by prefix")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		svr.JSON(w, http.StatusOK, locs)
	}
}

func AutocompleteSkill(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prefix := r.URL.Query().Get("k")
		skills, err := database.SeoSkills(
			qm.Where("name ILIKE ?", fmt.Sprintf("%s%%", prefix)),
			qm.OrderBy(database.SeoSkillColumns.Name+" ASC"),
			qm.Limit(5),
		).AllG(r.Context())
		if err != nil {
			svr.Log(err, "unable to retrieve skills by prefix")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		svr.JSON(w, http.StatusOK, skills)
	}
}

func DeliverMessageDeveloperProfileHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		messageID := vars["id"]
		message, err := repositories.MessageForDeliveryByID(r.Context(), messageID)
		if err != nil {
			svr.JSON(w, http.StatusBadRequest, "Your link may be invalid or expired")
			return
		}
		err = svr.GetEmail().SendHTMLEmail(
			email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
			email.Address{Email: message.R.Profile.Email},
			email.Address{Email: message.Email},
			fmt.Sprintf("New Message from %s", svr.GetConfig().SiteName),
			fmt.Sprintf(
				"You received a new message from %s: \n\nMessage: %s\n\nFrom: %s",
				svr.GetConfig().SiteName,
				message.Content,
				message.Email,
			),
		)
		if err != nil {
			svr.Log(err, "unable to send email to developer profile")
			svr.JSON(w, http.StatusBadRequest, "There was a problem while sending the email")
			return
		}
		if err := repositories.MarkDeveloperMessageAsSent(r.Context(), messageID); err != nil {
			svr.Log(err, "unable to mark developer message as sent "+messageID)
		}
		svr.JSON(w, http.StatusOK, "Message Sent Successfully")
	}
}

func EditProfileHandler(svr server.Server) http.HandlerFunc {
	return middleware.UserAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			profileID := vars["id"]
			profile, err := middleware.GetUserFromJWT(r, svr.SessionStore, svr.GetJWTSigningKey())
			if err != nil {
				svr.Log(err, "unable to get email from JWT")
				http.Redirect(w, r, "/auth", http.StatusUnauthorized)
				return
			}
			// todo: allow admin to edit any profile type
			// todo: check that only owners can edit their own profiles
			switch profile.Type {
			case user.UserTypeDeveloper:
				dev, err := repositories.DeveloperProfileByID(r.Context(), profileID)
				if err != nil {
					svr.Log(err, "unable to find developer profile")
					http.Redirect(w, r, "/auth", http.StatusUnauthorized)
					return
				}
				if dev.Email != profile.Email && !profile.IsAdmin {
					http.Redirect(w, r, "/auth", http.StatusUnauthorized)
					return
				}
				svr.Render(r, w, http.StatusOK, "edit-developer-profile.html", map[string]interface{}{
					"DeveloperProfile": dev,
				})
			case user.UserTypeRecruiter:
				rec, err := repositories.RecruiterProfileByID(r.Context(), profileID)
				if err != nil {
					svr.Log(err, "unable to find recruiter profile")
					http.Redirect(w, r, "/auth", http.StatusUnauthorized)
					return
				}
				svr.Render(r, w, http.StatusOK, "edit-recruiter-profile.html", map[string]interface{}{
					"RecruiterProfile": rec,
				})
			case user.UserTypeAdmin:
				svr.Log(err, "admin does not have profile to edit yet")
				http.Redirect(w, r, "/auth", http.StatusUnauthorized)
				return
			}
		},
	)
}

func ViewDeveloperProfileHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		profileSlug := vars["slug"]
		dev, err := repositories.DeveloperProfileBySlug(r.Context(), profileSlug)
		if err != nil {
			svr.Log(err, "unable to find developer profile by slug "+profileSlug)
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		if err := repositories.TrackDeveloperProfileView(r.Context(), dev); err != nil {
			svr.Log(err, "unable to track developer profile view")
		}

		svr.Render(r, w, http.StatusOK, "view-developer-profile.html", map[string]interface{}{
			"DeveloperProfile": dev,
			"UpdateAt":         dev.UpdatedAt.Time.UTC().Format("January 2006"),
			"Skills":           strings.Split(dev.Skills, ","),
			"MonthAndYear":     time.Now().UTC().Format("January 2006"),
		})
	}
}

func CompaniesForLocationHandler(svr server.Server, loc string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("p")
		svr.RenderPageForCompanies(w, r, loc, page, "companies.html")
	}
}

func PermanentRedirectHandler(svr server.Server, dst string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svr.Redirect(w, r, http.StatusMovedPermanently, fmt.Sprintf("https://%s/%s", svr.GetConfig().SiteHost, dst))
	}
}

func PermanentExternalRedirectHandler(svr server.Server, dst string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svr.Redirect(w, r, http.StatusMovedPermanently, dst)
	}
}

func PostAJobPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svr.RenderPostAJobForLocation(w, r, "")
	}
}

func ShowPaymentPage(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			svr.JSON(w, http.StatusBadRequest, "invalid email")
		}
		svr.Render(r, w, http.StatusOK, "payment.html", map[string]interface{}{
			"CurrencySymbol":       "$",
			"StripePublishableKey": svr.GetConfig().StripePublishableKey,
			"Email":                email,
			"Amount":               299,
			"AmountPence":          29900,
		})
	}
}

func PostAJobWithoutPaymentPageHandler(svr server.Server) http.HandlerFunc {
	return middleware.AdminAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			svr.Render(r, w, http.StatusOK, "post-a-job-without-payment.html", nil)
		},
	)
}

func SendFeedbackMessage(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &struct {
			Email   string `json:"email"`
			Message string `json:"message"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		if !svr.IsEmail(req.Email) {
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		if svr.SeenSince(r, time.Duration(1*time.Hour)) {
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		err := svr.
			GetEmail().
			SendHTMLEmail(
				email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
				email.Address{Email: svr.GetEmail().DefaultAdminAddress()},
				email.Address{Email: req.Email},
				"New Feedback Message",
				fmt.Sprintf("From: %s\nMessage: %s", req.Email, req.Message),
			)
		if err != nil {
			svr.Log(err, "unable to send email for feedback message")
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		svr.JSON(w, http.StatusOK, nil)
	}
}

func RequestTokenSignOn(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &struct {
			Email string `json:"email"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		if !svr.IsEmail(req.Email) {
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		u, err := repositories.GetUser(r.Context(), req.Email)
		if err != nil {
			svr.JSON(w, http.StatusNotFound, nil)
			return
		}
		k, err := ksuid.NewRandom()
		if err != nil {
			svr.Log(err, "unable to generate token")
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		err = repositories.SaveTokenSignOn(r.Context(), req.Email, k.String(), u.UserType)
		if err != nil {
			svr.Log(err, "unable to save sign on token")
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		token := k.String()
		err = svr.GetEmail().SendHTMLEmail(
			email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
			email.Address{Email: req.Email},
			email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
			fmt.Sprintf("Sign On on %s", svr.GetConfig().SiteName),
			fmt.Sprintf("Sign On on %s https://%s/x/auth/%s", svr.GetConfig().SiteName, svr.GetConfig().SiteHost, token))
		if err != nil {
			svr.Log(err, "unable to send email while applying to job")
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		svr.JSON(w, http.StatusOK, nil)
	}
}

func VerifyTokenSignOn(svr server.Server, adminEmail string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		token := vars["token"]
		u, _, err := repositories.GetOrCreateUserFromToken(r.Context(), token)
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to validate signon token %s", token))
			svr.TEXT(w, http.StatusBadRequest, "Invalid or expired token")
			return
		}
		fmt.Println("verify")
		sess, err := svr.SessionStore.Get(r, "____gc")
		if err != nil {
			svr.TEXT(w, http.StatusInternalServerError, "Invalid or expired token")
			svr.Log(err, "unable to get session cookie from request")
			return
		}
		stdClaims := &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour).UTC().Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
			Issuer:    fmt.Sprintf("https://%s", svr.GetConfig().SiteHost),
		}
		claims := middleware.UserJWT{
			UserID:         u.ID,
			Email:          u.Email,
			IsAdmin:        u.UserType == user.UserTypeAdmin,
			IsRecruiter:    u.UserType == user.UserTypeRecruiter,
			IsDeveloper:    u.UserType == user.UserTypeDeveloper,
			CreatedAt:      u.CreatedAt.Time,
			Type:           u.UserType,
			StandardClaims: *stdClaims,
		}
		tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := tkn.SignedString(svr.GetJWTSigningKey())
		if err != nil {
			svr.Log(err, "unable to sign jwt")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		sess.Values["jwt"] = ss
		err = sess.Save(r, w)
		if err != nil {
			svr.Log(err, "unable to save jwt into session cookie")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		fmt.Println("got step user type", u.UserType)
		switch u.UserType {
		case user.UserTypeDeveloper:
			dev, err := repositories.DeveloperProfileByEmail(r.Context(), u.Email)
			if err != nil {
				svr.Log(err, "unable to find developer profile by email")
				svr.JSON(w, http.StatusNotFound, "unable to find developer profile by email")
				return
			}
			if !dev.UpdatedAt.Time.After(dev.CreatedAt) {
				if activateDevProfileErr := repositories.ActivateDeveloperProfile(r.Context(), u.Email); activateDevProfileErr != nil {
					svr.Log(err, "unable to activate developer profile")
					svr.JSON(w, http.StatusInternalServerError, nil)
					return
				}
			}
			if _, err := database.EmailSubscribers(
				database.EmailSubscriberWhere.Token.EQ(token),
			).UpdateAllG(r.Context(), database.M{
				database.EmailSubscriberColumns.ConfirmedAt: time.Now(),
			}); err != nil {
				svr.Log(err, "unable to confirm subscriber using token "+token)
			}
			svr.Redirect(w, r, http.StatusMovedPermanently, "/profile/home")
			return
		case user.UserTypeRecruiter:
			rec, err := repositories.RecruiterProfileByEmail(r.Context(), u.Email)
			if err != nil {
				svr.Log(err, "unable to find recruiter profile by email")
				svr.JSON(w, http.StatusNotFound, "unable to find recruiter profile by email")
				return
			}
			if !rec.UpdatedAt.Time.After(rec.CreatedAt) {
				if activateRecProfileErr := repositories.ActivateRecruiterProfile(r.Context(), u.Email); activateRecProfileErr != nil {
					svr.Log(err, "unable to activate recruiter profile")
					svr.JSON(w, http.StatusInternalServerError, nil)
					return
				}
			}
			svr.Redirect(w, r, http.StatusMovedPermanently, "/profile/home")
			return
		case user.UserTypeAdmin:
			svr.Redirect(w, r, http.StatusMovedPermanently, "/profile/home")
			return
		}
		svr.Log(errors.New("unable to complete token verification flow"), fmt.Sprintf("email %s token %s and user type %s", u.Email, token, u.UserType))
		svr.Redirect(w, r, http.StatusMovedPermanently, "/")
	}
}

func ListJobsAsAdminPageHandler(svr server.Server) http.HandlerFunc {
	return middleware.AdminAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			loc := r.URL.Query().Get("l")
			skill := r.URL.Query().Get("s")
			page := r.URL.Query().Get("p")
			salary := ""
			currency := "USD"
			svr.RenderPageForLocationAndTagAdmin(r, w, loc, skill, page, salary, currency, "list-jobs-admin.html")
		},
	)
}

func PostAJobForLocationPageHandler(svr server.Server, location string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svr.RenderPostAJobForLocation(w, r, location)
	}
}

func PostAJobForLocationFromURLPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		location := vars["location"]
		location = strings.ReplaceAll(location, "-", " ")
		reg, err := regexp.Compile(`[^a-zA-Z0-9\\s]+`)
		if err != nil {
			log.Fatal(err)
		}
		location = reg.ReplaceAllString(location, "")
		svr.RenderPostAJobForLocation(w, r, location)
	}
}

func JobBySlugPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slug := vars["slug"]
		location := vars["l"]
		jobPost, err := repositories.JobPostBySlug(r.Context(), slug)
		if err != nil || jobPost == nil {
			svr.JSON(w, http.StatusNotFound, fmt.Sprintf("Job %s/job/%s not found", svr.GetConfig().SiteHost, slug))
			return
		}
		if err := repositories.TrackJobView(r.Context(), jobPost.ID); err != nil {
			svr.Log(err, fmt.Sprintf("unable to track job view for %s: %v", slug, err))
		}
		jobLocations := strings.Split(jobPost.Location, "/")
		var isQuickApply bool

		relevantJobs, err := repositories.GetRelevantJobs(r.Context(), jobPost.Location, jobPost.ID, 3)
		if err != nil {
			svr.Log(err, "unable to get relevant jobs")
		}
		for i, j := range relevantJobs {
			relevantJobs[i].Description = string(svr.MarkdownToHTML(j.Description))
			relevantJobs[i].SalaryRange = j.SalaryRange
		}
		svr.Render(r, w, http.StatusOK, "job.html", map[string]interface{}{
			"Job":                   jobPost,
			"JobURIEncoded":         url.QueryEscape(jobPost.Slug.String),
			"IsQuickApply":          isQuickApply,
			"HTMLJobDescription":    svr.MarkdownToHTML(jobPost.Description),
			"LocationFilter":        location,
			"ExternalJobId":         jobPost.ExternalID,
			"MonthAndYear":          jobPost.CreatedAt.UTC().Format("January 2006"),
			"GoogleJobCreatedAt":    jobPost.CreatedAt.Format(time.RFC3339),
			"GoogleJobValidThrough": jobPost.CreatedAt.AddDate(0, 5, 0),
			"GoogleJobLocation":     jobLocations[0],
			"GoogleJobDescription":  strconv.Quote(strings.ReplaceAll(string(svr.MarkdownToHTML(jobPost.Description)), "\n", "")),
			"RelevantJobs":          relevantJobs,
		})
	}
}

func CompanyBySlugPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slug := vars["slug"]
		company, err := repositories.CompanyBySlug(r.Context(), slug)
		if err != nil || company == nil {
			svr.JSON(w, http.StatusNotFound, fmt.Sprintf("Company %s/job/%s not found", svr.GetConfig().SiteHost, slug))
			return
		}
		if err := repositories.TrackCompanyView(r.Context(), company); err != nil {
			svr.Log(err, fmt.Sprintf("unable to track company view for %s: %v", slug, err))
		}
		companyJobs, err := repositories.GetCompanyJobs(r.Context(), company.Name, 3)
		if err != nil {
			svr.Log(err, "unable to get company jobs")
		}
		for i, j := range companyJobs {
			companyJobs[i].Description = string(svr.MarkdownToHTML(j.Description))
			companyJobs[i].SalaryRange = j.SalaryRange
		}
		if err := svr.Render(r, w, http.StatusOK, "company.html", map[string]interface{}{
			"Company":      company,
			"MonthAndYear": time.Now().UTC().Format("January 2006"),
			"CompanyJobs":  companyJobs,
		}); err != nil {
			svr.Log(err, "unable to render template")
		}
	}
}

func LandingPageForLocationHandler(svr server.Server, location string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		salary := vars["salary"]
		currency := vars["currency"]
		page := r.URL.Query().Get("p")
		svr.RenderPageForLocationAndTag(w, r, location, "", page, salary, currency, "landing.html")
	}
}

func LandingPageForLocationAndSkillPlaceholderHandler(svr server.Server, location string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		salary := vars["salary"]
		currency := vars["currency"]
		skill := strings.ReplaceAll(vars["skill"], "-", " ")
		page := r.URL.Query().Get("p")
		svr.RenderPageForLocationAndTag(w, r, location, skill, page, salary, currency, "landing.html")
	}
}

func LandingPageForLocationPlaceholderHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		salary := vars["salary"]
		currency := vars["currency"]
		loc := strings.ReplaceAll(vars["location"], "-", " ")
		page := r.URL.Query().Get("p")
		svr.RenderPageForLocationAndTag(w, r, loc, "", page, salary, currency, "landing.html")
	}
}

func LandingPageForSkillPlaceholderHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		salary := vars["salary"]
		currency := vars["currency"]
		skill := strings.ReplaceAll(vars["skill"], "-", " ")
		page := r.URL.Query().Get("p")
		svr.RenderPageForLocationAndTag(w, r, "", skill, page, salary, currency, "landing.html")
	}
}

func LandingPageForSkillAndLocationPlaceholderHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		salary := vars["salary"]
		currency := vars["currency"]
		loc := strings.ReplaceAll(vars["location"], "-", " ")
		skill := strings.ReplaceAll(vars["skill"], "-", " ")
		page := r.URL.Query().Get("p")
		svr.RenderPageForLocationAndTag(w, r, loc, skill, page, salary, currency, "landing.html")
	}
}

func ServeRSSFeed(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobPosts, err := repositories.GetLastNJobs(r.Context(), 20, r.URL.Query().Get("l"))
		if err != nil {
			svr.Log(err, "unable to retrieve jobs for RSS Feed")
			svr.XML(w, http.StatusInternalServerError, []byte{})
			return
		}
		now := time.Now()
		feed := &feeds.Feed{
			Title:       fmt.Sprintf("%s Jobs", svr.GetConfig().SiteName),
			Link:        &feeds.Link{Href: fmt.Sprintf("https://%s", svr.GetConfig().SiteHost)},
			Description: fmt.Sprintf("%s Jobs RSS Feed", svr.GetConfig().SiteName),
			Author:      &feeds.Author{Name: svr.GetConfig().SiteName, Email: svr.GetConfig().SupportEmail},
			Created:     now,
		}

		for _, j := range jobPosts {
			if j.CompanyIconImageID.Valid && j.CompanyIconImageID.String != "" {
				feed.Items = append(feed.Items, &feeds.Item{
					Title:       fmt.Sprintf("%s with %s - %s", j.JobTitle, j.Company, j.Location),
					Link:        &feeds.Link{Href: fmt.Sprintf("https://%s/job/%s", svr.GetConfig().SiteHost, j.Slug.String)},
					Description: string(svr.MarkdownToHTML(j.Description + "\n\n**Salary Range:** " + j.SalaryRange)),
					Author:      &feeds.Author{Name: svr.GetConfig().SiteName, Email: svr.GetConfig().SupportEmail},
					Enclosure:   &feeds.Enclosure{Length: "not implemented", Type: "image", Url: fmt.Sprintf("https://%s/x/s/m/%s", svr.GetConfig().SiteHost, j.CompanyIconImageID.String)},
					Created:     j.ApprovedAt.Time,
				})
			} else {
				feed.Items = append(feed.Items, &feeds.Item{
					Title:       fmt.Sprintf("%s with %s - %s", j.JobTitle, j.Company, j.Location),
					Link:        &feeds.Link{Href: fmt.Sprintf("https://%s/job/%s", svr.GetConfig().SiteHost, j.Slug.String)},
					Description: string(svr.MarkdownToHTML(j.Description + "\n\n**Salary Range:** " + j.SalaryRange)),
					Author:      &feeds.Author{Name: svr.GetConfig().SiteName, Email: svr.GetConfig().SupportEmail},
					Created:     j.ApprovedAt.Time,
				})
			}
		}
		rssFeed, err := feed.ToRss()
		if err != nil {
			svr.Log(err, "unable to convert rss feed to xml")
			svr.XML(w, http.StatusInternalServerError, []byte{})
			return
		}
		svr.XML(w, http.StatusOK, []byte(rssFeed))
	}
}

func StripePaymentConfirmationWebookHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		const MaxBodyBytes = int64(65536)
		req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			svr.Log(err, "error reading request body from stripe")
			svr.JSON(w, http.StatusServiceUnavailable, nil)
			return
		}

		stripeSig := req.Header.Get("Stripe-Signature")
		sess, err := payment.HandleCheckoutSessionComplete(body, svr.GetConfig().StripeEndpointSecret, stripeSig)
		if err != nil {
			svr.Log(err, "error while handling checkout session complete")
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		if sess != nil {
			affectedRows, err := database.PurchaseEvents(
				database.PurchaseEventWhere.StripeSessionID.EQ(sess.ID),
				database.PurchaseEventWhere.CompletedAt.IsNull(),
			).UpdateAllG(req.Context(), database.M{
				database.PurchaseEventColumns.CompletedAt: time.Now(),
			})
			if err != nil {
				svr.Log(err, "error while saving successful payment")
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			if affectedRows != 1 {
				svr.Log(errors.New("invalid number of rows affected when saving payment"), fmt.Sprintf("got %d expected 1", affectedRows))
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			jobPost, err := repositories.GetJobByStripeSessionID(req.Context(), sess.ID)
			if err != nil {
				svr.Log(errors.New("unable to find job by stripe session id"), fmt.Sprintf("session id %s", sess.ID))
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			purchaseEvent, err := database.PurchaseEvents(
				database.PurchaseEventWhere.StripeSessionID.EQ(sess.ID),
			).OneG(req.Context())
			if err != nil {
				svr.Log(errors.New("unable to find purchase event by stripe session id"), fmt.Sprintf("session id %s", sess.ID))
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			jobToken, err := repositories.TokenByJobID(req.Context(), jobPost.ID)
			if err != nil {
				svr.Log(errors.New("unable to find token for job id"), fmt.Sprintf("session id %s job id %s", sess.ID, jobPost.ID))
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}

			repositories.ApplyPlanTypeAndDurationToExpirations(jobPost)
			if err := repositories.UpdateJob(req.Context(), jobPost); err != nil {
				svr.Log(errors.New("unable to update job to new ad type"), fmt.Sprintf("unable to update job id %s for session id %s", jobPost.ID, sess.ID))
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			err = svr.GetEmail().SendHTMLEmail(
				email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().SupportSenderAddress()},
				email.Address{Email: purchaseEvent.Email},
				email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().SupportSenderAddress()},
				fmt.Sprintf("Your Job Ad is live on %s", svr.GetConfig().SiteName),
				fmt.Sprintf("Your Job Ad has been approved and it's now live. You can edit the Job Ad at any time and check page views and clickouts by following this link https://%s/edit/%s", svr.GetConfig().SiteHost, jobToken.Token))
			if err != nil {
				svr.Log(err, "unable to send email while upgrading job ad")
			}
			if err := svr.CacheDelete(server.CacheKeyPinnedJobs); err != nil {
				svr.Log(err, "unable to cleanup cache after approving job")
			}
			svr.JSON(w, http.StatusOK, nil)
			return
		}

		svr.JSON(w, http.StatusOK, nil)
	}
}

func SitemapIndexHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		index := sitemap.NewSitemapIndex()
		entries, err := GetSitemapIndex(r.Context(), svr.GetConfig().SiteHost)
		if err != nil {
			svr.Log(err, "database.GetSitemapIndex")
			svr.TEXT(w, http.StatusInternalServerError, "unable to fetch sitemap")
			return
		}
		for _, e := range entries {
			index.Add(&sitemap.URL{
				Loc:     e.Loc,
				LastMod: e.Lastmod.Ptr(),
			})
		}
		buf := new(bytes.Buffer)
		if _, err := index.WriteTo(buf); err != nil {
			svr.Log(err, "sitemapIndex.WriteTo")
			svr.TEXT(w, http.StatusInternalServerError, "unable to save sitemap index")
			return
		}
		svr.XML(w, http.StatusOK, buf.Bytes())
	}
}

func GetSitemapLastMod(ctx context.Context) (null.Time, error) {
	sitemap, err := database.Sitemaps(
		qm.OrderBy(database.SitemapColumns.Lastmod+" DESC"),
		qm.Limit(1),
	).OneG(ctx)
	if err != nil {
		return null.Time{}, err
	}
	return sitemap.Lastmod, nil
}

const SitemapSize = 1000

func GetSitemapEntryCount(ctx context.Context) (int64, error) {
	return database.Sitemaps().CountG(ctx)
}

func GetSitemapIndex(ctx context.Context, siteHost string) ([]database.Sitemap, error) {
	entries := make([]database.Sitemap, 0, 20)
	count, err := GetSitemapEntryCount(ctx)
	if err != nil {
		return entries, err
	}
	lastMod, err := GetSitemapLastMod(ctx)
	if err != nil {
		return entries, err
	}
	slots := math.Ceil(float64(count) / float64(SitemapSize))
	for i := 1; i <= int(slots); i++ {
		entries = append(entries, database.Sitemap{
			Loc:     fmt.Sprintf("https://%s/sitemap-%d.xml", siteHost, i),
			Lastmod: lastMod,
		})
	}

	return entries, nil
}

func GetSitemapNo(ctx context.Context, n int) (database.SitemapSlice, error) {
	return database.Sitemaps(
		qm.Limit(SitemapSize),
		qm.Offset((n-1)*SitemapSize),
	).AllG(ctx)
}

func SitemapHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sitemapNo := vars["number"]
		number, err := strconv.Atoi(sitemapNo)
		if err != nil || number < 1 {
			svr.Log(err, fmt.Sprintf("unable to parse sitemap number %s", sitemapNo))
			svr.TEXT(w, http.StatusBadRequest, "invalid sitemap number")
			return
		}
		entries, err := GetSitemapNo(r.Context(), number)
		if err != nil {
			svr.Log(err, fmt.Sprintf("database.GetSitemapNo %d", number))
			svr.TEXT(w, http.StatusInternalServerError, "unable to fetch sitemap")
			return
		}
		sitemapFile := sitemap.New()
		for _, e := range entries {
			sitemapFile.Add(&sitemap.URL{
				Loc:        e.Loc,
				LastMod:    e.Lastmod.Ptr(),
				ChangeFreq: sitemap.ChangeFreq(e.Changefreq.String),
			})
		}
		buf := new(bytes.Buffer)
		if _, err := sitemapFile.WriteTo(buf); err != nil {
			svr.Log(err, fmt.Sprintf("sitemapFile.WriteTo %d", number))
			svr.TEXT(w, http.StatusInternalServerError, "unable to save sitemap file")
			return
		}
		svr.XML(w, http.StatusOK, buf.Bytes())
	}
}

func RobotsTXTHandler(svr server.Server, robotsTxtContent []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svr.TEXT(w, http.StatusOK, strings.ReplaceAll(string(robotsTxtContent), "__host_placeholder__", svr.GetConfig().SiteHost))
	}
}

func WellKnownSecurityHandler(svr server.Server, securityTxtContent []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentWithHost := strings.ReplaceAll(string(securityTxtContent), "__host_placeholder__", svr.GetConfig().SiteHost)
		svr.TEXT(w, http.StatusOK, strings.ReplaceAll(contentWithHost, "__support_email_placeholder__", svr.GetConfig().SupportEmail))
	}
}

func AboutPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svr.Render(r, w, http.StatusOK, "about.html", nil)
	}
}

func PrivacyPolicyPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svr.Render(r, w, http.StatusOK, "privacy-policy.html", nil)
	}
}

func TermsOfServicePageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svr.Render(r, w, http.StatusOK, "terms-of-service.html", nil)
	}
}

// func SalaryLandingPageLocationPlaceholderHandler(svr server.Server) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		location := strings.ReplaceAll(vars["location"], "-", " ")
// 		svr.RenderSalaryForLocation(r.Context(), w, r, devRepo, location)
// 	}
// }

// func SalaryLandingPageLocationHandler(svr server.Server, location string) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		svr.RenderSalaryForLocation(r.Context(), w, r, devRepo, location)
// 	}
// }

func ViewNewsletterPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svr.RenderPageForLocationAndTag(w, r, "", "", "", "", "", "newsletter.html")
	}
}

func ViewCommunityNewsletterPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svr.RenderPageForLocationAndTag(w, r, "", "", "", "", "", "news.html")
	}
}

func DisableDirListing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ViewSupportPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svr.RenderPageForLocationAndTag(w, r, "", "", "", "", "", "support.html")
	}
}

var allowedMediaTypes = []string{"image/png", "image/jpeg", "image/jpg"}

func PostAJobSuccessPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svr.Render(r, w, http.StatusOK, "post-a-job-success.html", nil)
	}
}

func PostAJobFailurePageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svr.Render(r, w, http.StatusOK, "post-a-job-error.html", nil)
	}
}

func ApplyForJobPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// limits upload form size to 5mb
		maxPdfSize := 5 * 1024 * 1024
		ctx := r.Context()
		r.Body = http.MaxBytesReader(w, r.Body, int64(maxPdfSize))
		cv, header, err := r.FormFile("cv")
		if err != nil {
			svr.Log(err, "unable to read cv file")
			svr.JSON(w, http.StatusRequestEntityTooLarge, nil)
			return
		}
		defer cv.Close()
		fileBytes, err := ioutil.ReadAll(cv)
		if err != nil {
			svr.Log(err, "unable to read cv file content")
			svr.JSON(w, http.StatusRequestEntityTooLarge, nil)
			return
		}
		contentType := http.DetectContentType(fileBytes)
		if contentType != "application/pdf" {
			svr.Log(errors.New("PDF file is not application/pdf"), fmt.Sprintf("PDF file is not application/pdf got %s", contentType))
			svr.JSON(w, http.StatusUnsupportedMediaType, nil)
			return
		}
		if header.Size > int64(maxPdfSize) {
			svr.Log(errors.New("PDF file is too large"), fmt.Sprintf("PDF file too large: %d > %d", header.Size, maxPdfSize))
			svr.JSON(w, http.StatusRequestEntityTooLarge, nil)
			return
		}
		externalID := r.FormValue("job-id")
		emailAddr := r.FormValue("email")
		jobPost, err := repositories.JobPostByExternalID(ctx, externalID)
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to retrieve job by externalId %s, %v", externalID, err))
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		k, err := ksuid.NewRandom()
		if err != nil {
			svr.Log(err, "unable to generate token")
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		randomToken, err := k.Value()
		if err != nil {
			svr.Log(err, "unable to get token value")
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		randomTokenStr, ok := randomToken.(string)
		if !ok {
			svr.Log(err, "unable to assert token value as string")
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		profile, _ := middleware.GetUserFromJWT(r, svr.SessionStore, svr.GetJWTSigningKey())
		// user is not logged in
		// standard flow to confirm application
		if profile == nil {
			err = repositories.ApplyToJob(r.Context(), jobPost.ID, fileBytes, emailAddr, randomTokenStr)
			if err != nil {
				svr.Log(err, "unable to apply for job while saving to db")
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			err = svr.GetEmail().SendHTMLEmail(
				email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
				email.Address{Email: emailAddr},
				email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
				fmt.Sprintf("Confirm your job application with %s", jobPost.Company),
				fmt.Sprintf(
					"Thanks for applying for the position %s with %s - %s.<br />Please confirm your application now by following this link https://%s/apply/%s",
					jobPost.JobTitle,
					jobPost.Company,
					jobPost.Location,
					svr.GetConfig().SiteHost,
					randomTokenStr,
				),
			)
			if err != nil {
				svr.Log(err, "unable to send email while applying to job")
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			if r.FormValue("notify-jobs") == "true" {
				k, err := ksuid.NewRandom()
				if err != nil {
					svr.Log(err, "unable to generate email subscriber token")
					svr.JSON(w, http.StatusBadRequest, nil)
					return
				}
				emailEntry := database.EmailSubscriber{
					Email:     emailAddr,
					Token:     k.String(),
					CreatedAt: time.Now(),
				}
				err = emailEntry.InsertG(ctx, boil.Infer())
				if err != nil {
					svr.Log(err, "unable to add email subscriber to db")
					svr.JSON(w, http.StatusInternalServerError, nil)
					return
				}
				err = svr.GetEmail().SendHTMLEmail(
					email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
					email.Address{Email: emailAddr},
					email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
					fmt.Sprintf("Confirm Your Email Subscription on %s", svr.GetConfig().SiteName),
					fmt.Sprintf(
						"Please click on the link below to confirm your subscription to receive weekly emails from %s\n\n%s\n\nIf this was not requested by you, please ignore this email.",
						svr.GetConfig().SiteName,
						fmt.Sprintf("https://%s/x/email/confirm/%s", svr.GetConfig().SiteHost, k.String()),
					),
				)
				if err != nil {
					svr.Log(err, "unable to send email while submitting message")
					svr.JSON(w, http.StatusBadRequest, nil)
					return
				}
			}
			svr.JSON(w, http.StatusOK, nil)
			return
		}
		if profile.Email != emailAddr {
			svr.JSON(w, http.StatusBadRequest, "Please use the same email address you have registered on your profile.")
			return
		}
		err = repositories.ApplyToJob(r.Context(), jobPost.ID, fileBytes, emailAddr, randomTokenStr)
		if err != nil {
			svr.Log(err, "unable to apply for job while saving to db")
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		token, err := database.ApplyTokens(
			qm.Load(database.ApplyTokenRels.Job),
			database.ApplyTokenWhere.Token.EQ(randomTokenStr),
		).OneG(ctx)
		if err != nil {
			svr.Render(r, w, http.StatusBadRequest, "apply-message.html", map[string]interface{}{
				"Title":       "Invalid Job Application",
				"Description": "Oops, seems like the application you are trying to complete is no longer valid. Your application request may be expired or simply the company may not be longer accepting applications.",
			})
			return
		}
		err = svr.GetEmail().SendEmailWithPDFAttachment(
			email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
			email.Address{Email: token.R.Job.SubscriberEmail},
			email.Address{Email: token.Email},
			fmt.Sprintf("New Applicant from %s", svr.GetConfig().SiteName),
			fmt.Sprintf(
				"Hi, there is a new applicant for your position on %s: %s with %s - %s (https://%s/job/%s). Applicant's Email: %s. Please find applicant's CV attached below",
				svr.GetConfig().SiteName,
				token.R.Job.JobTitle,
				token.R.Job.Company,
				token.R.Job.Location,
				svr.GetConfig().SiteHost,
				token.R.Job.Slug.String,
				token.Email,
			),
			token.CV,
			"cv.pdf",
		)
		if err != nil {
			svr.Log(err, "unable to send email while applying to job")
			svr.Render(r, w, http.StatusBadRequest, "apply-message.html", map[string]interface{}{
				"Title":       "Job Application Failure",
				"Description": fmt.Sprintf("Oops, there was a problem while completing yuor application. Please try again later. If the problem persists, please contact %s", svr.GetConfig().SupportEmail),
			})
			return
		}
		err = repositories.ConfirmApplyToJob(r.Context(), randomTokenStr)
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to update apply_token with successfull application for token %s", randomTokenStr))
			svr.Render(r, w, http.StatusBadRequest, "apply-message.html", map[string]interface{}{
				"Title":       "Job Application Failure",
				"Description": fmt.Sprintf("Oops, there was a problem while completing yuor application. Please try again later. If the problem persists, please contact %s", svr.GetConfig().SupportEmail),
			})
			return
		}
		svr.Render(r, w, http.StatusOK, "apply-message.html", map[string]interface{}{
			"Title": "Job Application Successfull",
			"Description": svr.StringToHTML(
				fmt.Sprintf(
					"Thank you for applying for <b>%s with %s - %s</b><br /><a href=\"https://%s/job/%s\">https://%s/job/%s</a>. <br /><br />Your CV has been forwarded to company HR. <br />Consider joining our Golang Cafe Developer community where companies can apply to you",
					token.R.Job.JobTitle,
					token.R.Job.Company,
					token.R.Job.Location,
					svr.GetConfig().SiteHost,
					token.R.Job.Slug.String,
					svr.GetConfig().SiteHost,
					token.R.Job.Slug.String,
				),
			),
		})
	}
}

func ApplyToJobConfirmation(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		pToken := vars["token"]
		token, err := database.ApplyTokens(
			qm.Load(database.ApplyTokenRels.Job),
			database.ApplyTokenWhere.Token.EQ(pToken),
		).OneG(r.Context())
		if err != nil {
			svr.Render(r, w, http.StatusBadRequest, "apply-message.html", map[string]interface{}{
				"Title":       "Invalid Job Application",
				"Description": "Oops, seems like the application you are trying to complete is no longer valid. Your application request may be expired or simply the company may not be longer accepting applications.",
			})
			return
		}
		err = svr.GetEmail().SendEmailWithPDFAttachment(
			email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
			email.Address{Email: token.R.Job.SubscriberEmail},
			email.Address{Email: token.Email},
			fmt.Sprintf("New Applicant from %s", svr.GetConfig().SiteName),
			fmt.Sprintf(
				"Hi, there is a new applicant for your position on %s: %s with %s - %s (https://%s/job/%s). Applicant's Email: %s. Please find applicant's CV attached below",
				svr.GetConfig().SiteName,
				token.R.Job.JobTitle,
				token.R.Job.Company,
				token.R.Job.Location,
				svr.GetConfig().SiteHost,
				token.R.Job.Slug.String,
				token.Email,
			),
			token.CV,
			"cv.pdf",
		)
		if err != nil {
			svr.Log(err, "unable to send email while applying to job")
			svr.Render(r, w, http.StatusBadRequest, "apply-message.html", map[string]interface{}{
				"Title":       "Job Application Failure",
				"Description": fmt.Sprintf("Oops, there was a problem while completing yuor application. Please try again later. If the problem persists, please contact %s", svr.GetConfig().SupportEmail),
			})
			return
		}
		token.ConfirmedAt = null.TimeFrom(time.Now())
		_, err = token.UpdateG(r.Context(), boil.Whitelist(database.ApplyTokenColumns.ConfirmedAt))
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to update apply_token with successfull application for token %s", token.Token))
			svr.Render(r, w, http.StatusBadRequest, "apply-message.html", map[string]interface{}{
				"Title":       "Job Application Failure",
				"Description": fmt.Sprintf("Oops, there was a problem while completing yuor application. Please try again later. If the problem persists, please contact %s", svr.GetConfig().SupportEmail),
			})
			return
		}
		svr.Render(r, w, http.StatusOK, "apply-message.html", map[string]interface{}{
			"Title": "Job Application Successfull",
			"Description": svr.StringToHTML(
				fmt.Sprintf(
					"Thank you for applying for <b>%s with %s - %s</b><br /><a href=\"https://%s/job/%s\">https://%s/job/%s</a>. <br /><br />Your CV has been forwarded to company HR. <br />Consider joining our Golang Cafe Developer community where companies can apply to you",
					token.R.Job.JobTitle,
					token.R.Job.Company,
					token.R.Job.Location,
					svr.GetConfig().SiteHost,
					token.R.Job.Slug.String,
					svr.GetConfig().SiteHost,
					token.R.Job.Slug.String,
				),
			),
		})
	}
}

func SubmitJobPostWithoutPaymentHandler(svr server.Server) http.HandlerFunc {
	return middleware.AdminAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			decoder := json.NewDecoder(r.Body)
			jobRq := &job.JobRq{}
			if err := decoder.Decode(&jobRq); err != nil {
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			job := &database.Job{
				JobTitle:           jobRq.JobTitle,
				JobCategory:        jobRq.JobCategory,
				Company:            jobRq.Company,
				Location:           jobRq.Location,
				SalaryRange:        jobRq.SalaryRange,
				JobType:            jobRq.JobType,
				ApplicationLink:    jobRq.ApplicationLink,
				SubscriberEmail:    jobRq.SubscriberEmail,
				Description:        jobRq.Description,
				CompanyIconImageID: null.StringFrom(jobRq.CompanyIconID),
			}
			err := repositories.SaveDraft(r.Context(), job)
			if err != nil {
				svr.Log(err, fmt.Sprintf("unable to save job request: %#v", jobRq))
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			k, err := ksuid.NewRandom()
			if err != nil {
				svr.Log(err, "unable to generate unique token")
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			randomToken, err := k.Value()
			if err != nil {
				svr.Log(err, "unable to get token value")
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			randomTokenStr, ok := randomToken.(string)
			if !ok {
				svr.Log(err, "unbale to assert token value as string")
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			err = repositories.SaveTokenForJob(r.Context(), randomTokenStr, job.ID)
			if err != nil {
				svr.Log(err, "unable to generate token")
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			svr.JSON(w, http.StatusOK, map[string]interface{}{"token": randomTokenStr})
		},
	)
}

func SubmitJobPostPaymentUpsellPageHandler(svr server.Server, paymentRepo *payment.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		jobRq := &job.JobRqUpsell{}
		if err := decoder.Decode(&jobRq); err != nil {
			svr.Log(err, "unable to decode request")
			svr.JSON(w, http.StatusBadRequest, err.Error())
			return
		}
		planDuration, err := strconv.Atoi(jobRq.PlanDurationStr)
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to convert duration to int %s", jobRq.PlanDurationStr))
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		jobRq.PlanDuration = planDuration
		editToken, err := database.EditTokens(
			database.EditTokenWhere.Token.EQ(jobRq.Token),
		).OneG(r.Context())
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to find job by token %s", jobRq.Token))
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		planId := svr.GetConfig().PlanPriceID
		sess, err := paymentRepo.CreateSession(planId)
		if err != nil {
			svr.Log(err, "unable to create subscription")
		}

		err = svr.GetEmail().SendHTMLEmail(
			email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
			email.Address{Email: svr.GetEmail().DefaultAdminAddress()},
			email.Address{Email: jobRq.Email},
			fmt.Sprintf("New Upgrade on %s", svr.GetConfig().SiteName),
			fmt.Sprintf(
				"Hey! There is a new ad upgrade on %s. Please check https://%s/manage/%s",
				svr.GetConfig().SiteName,
				svr.GetConfig().SiteHost,
				jobRq.Token,
			),
		)
		if err != nil {
			svr.Log(err, "unable to send email to admin while upgrading job ad")
		}
		purchaseEvent := &database.PurchaseEvent{
			StripeSessionID: sess.ID,
			Email:           jobRq.Email,
			PlanID:          planId,
			JobID:           editToken.JobID,
			CreatedAt:       time.Now(),
		}

		err = purchaseEvent.InsertG(r.Context(), boil.Infer())
		if err != nil {
			svr.Log(err, "unable to save payment initiated event")
		}
		svr.JSON(w, http.StatusOK, map[string]string{"s_id": sess.ID})

	}
}

func SubmitJobPostPageHandler(svr server.Server, paymentRepo *payment.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		jobRq := &job.JobRq{}
		if err := decoder.Decode(&jobRq); err != nil {
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		jobID, err := ksuid.NewRandom()
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to generate unique job id: %#v", jobRq))
			svr.JSON(w, http.StatusBadRequest, err.Error())
			return
		}
		expiration := null.TimeFrom(time.Now().AddDate(0, 0, 365*10))
		jobDraft := &database.Job{
			ID:                              jobID.String(),
			JobTitle:                        jobRq.JobTitle,
			JobCategory:                     jobRq.JobCategory,
			Company:                         jobRq.Company,
			Location:                        jobRq.Location,
			SalaryRange:                     jobRq.SalaryRange,
			JobType:                         jobRq.JobType,
			ApplicationLink:                 jobRq.ApplicationLink,
			Description:                     jobRq.Description,
			CreatedAt:                       time.Now().UTC(),
			URLID:                           int(time.Now().UTC().Unix()),
			Slug:                            null.StringFrom(slug.Make(fmt.Sprintf("%s %s %d", jobRq.JobTitle, jobRq.Company, time.Now().UTC().Unix()))),
			CompanyIconImageID:              null.StringFrom(jobRq.CompanyIconID),
			BlogEligibilityExpiredAt:        expiration,
			CompanyPageEligibilityExpiredAt: expiration,
			FrontPageEligibilityExpiredAt:   expiration,
			NewsletterEligibilityExpiredAt:  expiration,
			PlanExpiredAt:                   expiration,
			SocialMediaEligibilityExpiredAt: expiration,
		}
		err = jobDraft.InsertG(r.Context(), boil.Infer())
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to save job request: %#v", jobRq))
			svr.JSON(w, http.StatusBadRequest, err.Error())
			return
		}

		k, err := ksuid.NewRandom()
		if err != nil {
			svr.Log(err, "unable to generate unique token")
			svr.JSON(w, http.StatusBadRequest, err.Error())
			return
		}
		randomToken, err := k.Value()
		if err != nil {
			svr.Log(err, "unable to get token value")
			svr.JSON(w, http.StatusBadRequest, err.Error())
			return
		}
		randomTokenStr, ok := randomToken.(string)
		if !ok {
			svr.Log(err, "unbale to assert token value as string")
			svr.JSON(w, http.StatusBadRequest, "unbale to assert token value as string")
			return
		}
		err = jobDraft.AddEditTokensG(r.Context(), true, &database.EditToken{
			Token:     randomTokenStr,
			CreatedAt: time.Now().UTC(),
		})
		if err != nil {
			svr.Log(err, "unable to generate token")
			svr.JSON(w, http.StatusBadRequest, err.Error())
			return
		}
		planPriceID := svr.GetConfig().PlanPriceID
		sess, err := paymentRepo.CreateSession(planPriceID)
		if err != nil {
			svr.Log(err, "unable to create payment session")
		}
		err = svr.GetEmail().SendHTMLEmail(
			email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().NoReplySenderAddress()},
			email.Address{Email: svr.GetEmail().DefaultAdminAddress()},
			email.Address{Email: jobRq.Email},
			fmt.Sprintf("New Job Ad on %s", svr.GetConfig().SiteName),
			fmt.Sprintf(
				"Hey! There is a new Ad on %s. Please approve https://%s/manage/%s",
				svr.GetConfig().SiteName,
				svr.GetConfig().SiteHost,
				randomTokenStr,
			),
		)
		if err != nil {
			svr.Log(err, "unable to send email to admin while posting job ad")
		}
		if sess.ID != "" {
			purchaseEvent := &database.PurchaseEvent{
				StripeSessionID: sess.ID,
				Email:           jobRq.Email,
				PlanID:          planPriceID,
				JobID:           jobID.String(),
				CreatedAt:       time.Now(),
			}
			err = purchaseEvent.InsertG(r.Context(), boil.Infer())
			if err != nil {
				svr.Log(err, "unable to save payment initiated event")
			}
			svr.JSON(w, http.StatusOK, map[string]string{"s_id": sess.ID})
			return
		}
		svr.JSON(w, http.StatusOK, nil)
	}
}

func RetrieveMediaPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		mediaID := vars["id"]
		media, err := database.FindImageG(r.Context(), mediaID)
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to retrieve media by ID: '%s'", mediaID))
			svr.MEDIA(w, http.StatusNotFound, []byte{}, "")
			return
		}
		height := r.URL.Query().Get("h")
		width := r.URL.Query().Get("w")
		if height == "" && width == "" {
			svr.MEDIA(w, http.StatusOK, media.Bytes, media.MediaType)
			return
		}
		he, err := strconv.Atoi(height)
		if err != nil {
			svr.MEDIA(w, http.StatusOK, media.Bytes, media.MediaType)
			return
		}
		wi, err := strconv.Atoi(width)
		if err != nil {
			svr.MEDIA(w, http.StatusOK, media.Bytes, media.MediaType)
			return
		}
		contentTypeInvalid := true
		for _, allowedMedia := range allowedMediaTypes {
			if allowedMedia == media.MediaType {
				contentTypeInvalid = false
			}
		}
		if contentTypeInvalid {
			svr.Log(errors.New("invalid media content type"), fmt.Sprintf("media file %s is not one of the allowed media types: %+v", media.MediaType, allowedMediaTypes))
			svr.JSON(w, http.StatusUnsupportedMediaType, nil)
			return
		}
		decImage, _, err := image.Decode(bytes.NewReader(media.Bytes))
		if err != nil {
			svr.Log(err, "unable to decode image from bytes")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		m := resize.Resize(uint(wi), uint(he), decImage, resize.Lanczos3)
		resizeImageBuf := new(bytes.Buffer)
		switch media.MediaType {
		case "image/jpg", "image/jpeg":
			if err := jpeg.Encode(resizeImageBuf, m, nil); err != nil {
				svr.Log(err, "unable to encode resizeImage into jpeg")
				svr.JSON(w, http.StatusInternalServerError, nil)
				return
			}
		case "image/png":
			if err := png.Encode(resizeImageBuf, m); err != nil {
				svr.Log(err, "unable to encode resizeImage into png")
				svr.JSON(w, http.StatusInternalServerError, nil)
				return
			}
		default:
			svr.MEDIA(w, http.StatusOK, media.Bytes, media.MediaType)
			return
		}
		svr.MEDIA(w, http.StatusOK, resizeImageBuf.Bytes(), media.MediaType)
	}
}

func RetrieveMediaMetaPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		jobID := vars["id"]
		job, err := repositories.GetJobByExternalID(r.Context(), jobID)
		if err != nil {
			svr.Log(err, "unable to retrieve job by external ID")
			svr.MEDIA(w, http.StatusNotFound, []byte{}, "image/png")
			return
		}
		media, err := imagemeta.GenerateImageForJob(job)
		if err != nil {
			svr.Log(err, "unable to generate media for job ID")
			svr.MEDIA(w, http.StatusNotFound, []byte{}, "image/png")
			return
		}
		mediaBytes, err := ioutil.ReadAll(media)
		if err != nil {
			svr.Log(err, "unable to generate media for job ID")
			svr.MEDIA(w, http.StatusNotFound, mediaBytes, "image/png")
			return
		}
		svr.MEDIA(w, http.StatusOK, mediaBytes, "image/png")
	}
}

func UpdateMediaPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var x, y, wi, he int
		var err error
		x, err = strconv.Atoi(r.URL.Query().Get("x"))
		if err != nil {
			x = 0
		}
		y, err = strconv.Atoi(r.URL.Query().Get("y"))
		if err != nil {
			y = 0
		}
		wi, err = strconv.Atoi(r.URL.Query().Get("w"))
		if err != nil {
			wi = 0
		}
		he, err = strconv.Atoi(r.URL.Query().Get("h"))
		if err != nil {
			he = 0
		}
		vars := mux.Vars(r)
		mediaID := vars["id"]
		// limits upload form size to 5mb
		maxMediaFileSize := 5 * 1024 * 1024
		r.Body = http.MaxBytesReader(w, r.Body, int64(maxMediaFileSize))
		imageFile, header, err := r.FormFile("image")
		if err != nil {
			svr.Log(err, "unable to read media file")
			svr.JSON(w, http.StatusRequestEntityTooLarge, nil)
			return
		}
		defer imageFile.Close()
		fileBytes, err := ioutil.ReadAll(imageFile)
		if err != nil {
			svr.Log(err, "unable to read imageFile file content")
			svr.JSON(w, http.StatusRequestEntityTooLarge, nil)
			return
		}
		contentType := http.DetectContentType(fileBytes)
		contentTypeInvalid := true
		for _, allowedMedia := range allowedMediaTypes {
			if allowedMedia == contentType {
				contentTypeInvalid = false
			}
		}
		if contentTypeInvalid {
			svr.Log(errors.New("invalid media content type"), fmt.Sprintf("media file %s is not one of the allowed media types: %+v", contentType, allowedMediaTypes))
			svr.JSON(w, http.StatusUnsupportedMediaType, nil)
			return
		}
		if header.Size > int64(maxMediaFileSize) {
			svr.Log(errors.New("media file is too large"), fmt.Sprintf("media file too large: %d > %d", header.Size, maxMediaFileSize))
			svr.JSON(w, http.StatusRequestEntityTooLarge, nil)
			return
		}
		decImage, _, err := image.Decode(bytes.NewReader(fileBytes))
		if err != nil {
			svr.Log(err, "unable to decode image from bytes")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		min := decImage.Bounds().Dy()
		if decImage.Bounds().Dx() < min {
			min = decImage.Bounds().Dx()
		}
		if he == 0 || wi == 0 || he != wi {
			he = min
			wi = min
		}
		cutImage := decImage.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(image.Rect(x, y, wi+x, he+y))
		cutImageBytes := new(bytes.Buffer)
		switch contentType {
		case "image/jpg", "image/jpeg":
			if err := jpeg.Encode(cutImageBytes, cutImage, nil); err != nil {
				svr.Log(err, "unable to encode cutImage into jpeg")
				svr.JSON(w, http.StatusInternalServerError, nil)
				return
			}
		case "image/png":
			if err := png.Encode(cutImageBytes, cutImage); err != nil {
				svr.Log(err, "unable to encode cutImage into png")
				svr.JSON(w, http.StatusInternalServerError, nil)
				return
			}
		default:
			svr.Log(errors.New("content type not supported for encoding"), fmt.Sprintf("content type %s not supported for encoding", contentType))
			svr.JSON(w, http.StatusInternalServerError, nil)
		}
		_, err = database.Images(
			database.ImageWhere.ID.EQ(mediaID),
		).UpdateAllG(r.Context(), database.M{
			database.ImageColumns.Bytes:     cutImageBytes.Bytes(),
			database.ImageColumns.MediaType: contentType,
		})
		if err != nil {
			svr.Log(err, "unable to update media image to db")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		svr.JSON(w, http.StatusOK, nil)
	}
}

func SaveMediaPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var x, y, wi, he int
		var err error
		x, err = strconv.Atoi(r.URL.Query().Get("x"))
		if err != nil {
			x = 0
		}
		y, err = strconv.Atoi(r.URL.Query().Get("y"))
		if err != nil {
			y = 0
		}
		wi, err = strconv.Atoi(r.URL.Query().Get("w"))
		if err != nil {
			wi = 0
		}
		he, err = strconv.Atoi(r.URL.Query().Get("h"))
		if err != nil {
			he = 0
		}
		// limits upload form size to 5mb
		maxMediaFileSize := 5 * 1024 * 1024
		allowedMediaTypes := []string{"image/png", "image/jpeg", "image/jpg"}
		r.Body = http.MaxBytesReader(w, r.Body, int64(maxMediaFileSize))
		cv, header, err := r.FormFile("image")
		if err != nil {
			svr.Log(err, "unable to read media file")
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		defer cv.Close()
		fileBytes, err := ioutil.ReadAll(cv)
		if err != nil {
			svr.Log(err, "unable to read cv file content")
			svr.JSON(w, http.StatusRequestEntityTooLarge, nil)
			return
		}
		contentType := http.DetectContentType(fileBytes)
		contentTypeInvalid := true
		for _, allowedMedia := range allowedMediaTypes {
			if allowedMedia == contentType {
				contentTypeInvalid = false
			}
		}
		if contentTypeInvalid {
			svr.Log(errors.New("invalid media content type"), fmt.Sprintf("media file %s is not one of the allowed media types: %+v", contentType, allowedMediaTypes))
			svr.JSON(w, http.StatusUnsupportedMediaType, nil)
			return
		}
		if header.Size > int64(maxMediaFileSize) {
			svr.Log(errors.New("media file is too large"), fmt.Sprintf("media file too large: %d > %d", header.Size, maxMediaFileSize))
			svr.JSON(w, http.StatusRequestEntityTooLarge, nil)
			return
		}
		decImage, _, err := image.Decode(bytes.NewReader(fileBytes))
		if err != nil {
			svr.Log(err, "unable to decode image from bytes")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		min := decImage.Bounds().Dy()
		if decImage.Bounds().Dx() < min {
			min = decImage.Bounds().Dx()
		}
		if he == 0 || wi == 0 || wi != he {
			he = min
			wi = min
		}
		cutImage := decImage.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(image.Rect(x, y, wi, he))
		cutImageBytes := new(bytes.Buffer)
		switch contentType {
		case "image/jpg", "image/jpeg":
			if err := jpeg.Encode(cutImageBytes, cutImage, nil); err != nil {
				svr.Log(err, "unable to encode cutImage into jpeg")
				svr.JSON(w, http.StatusInternalServerError, nil)
				return
			}
		case "image/png":
			if err := png.Encode(cutImageBytes, cutImage); err != nil {
				svr.Log(err, "unable to encode cutImage into png")
				svr.JSON(w, http.StatusInternalServerError, nil)
				return
			}
		default:
			svr.Log(errors.New("content type not supported for encoding"), fmt.Sprintf("content type %s not supported for encoding", contentType))
			svr.JSON(w, http.StatusInternalServerError, nil)
		}
		id, err := ksuid.NewRandom()
		if err != nil {
			svr.Log(err, "unable to generate ksuid for media")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		media := database.Image{
			ID:        id.String(),
			Bytes:     cutImageBytes.Bytes(),
			MediaType: contentType,
		}
		err = media.InsertG(r.Context(), boil.Infer())
		if err != nil {
			svr.Log(err, "unable to save media image to db")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		svr.JSON(w, http.StatusOK, map[string]interface{}{"id": id})
	}
}

func UpdateJobPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		jobRq := &job.JobRqUpdate{}
		if err := decoder.Decode(&jobRq); err != nil {
			svr.Log(err, fmt.Sprintf("unable to parse job request for update: %#v", jobRq))
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		job, err := repositories.JobPostByToken(r.Context(), jobRq.Token)
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to find job post ID by token: %s", jobRq.Token))
			svr.JSON(w, http.StatusNotFound, nil)
			return
		}

		job.JobTitle = jobRq.JobTitle
		job.JobCategory = jobRq.JobCategory
		job.Company = jobRq.Company
		job.Location = jobRq.Location
		job.SalaryRange = jobRq.SalaryRange
		job.Description = jobRq.Description
		job.CompanyIconImageID = null.StringFrom(jobRq.CompanyIconID)

		_, err = job.UpdateG(r.Context(), boil.Whitelist(
			database.JobColumns.JobTitle,
			database.JobColumns.JobCategory,
			database.JobColumns.Company,
			database.JobColumns.Location,
			database.JobColumns.SalaryRange,
			database.JobColumns.Description,
			database.JobColumns.CompanyIconImageID,
		))
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to update job request: %#v", jobRq))
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		if err := svr.CacheDelete(server.CacheKeyPinnedJobs); err != nil {
			svr.Log(err, "unable to cleanup cache after approving job")
		}
		svr.JSON(w, http.StatusOK, nil)
	}
}
func PermanentlyDeleteJobByToken(svr server.Server) http.HandlerFunc {
	return middleware.AdminAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			decoder := json.NewDecoder(r.Body)
			jobRq := &job.JobRqUpdate{}
			if err := decoder.Decode(&jobRq); err != nil {
				svr.Log(err, fmt.Sprintf("unable to parse job request for delete: %#v", jobRq))
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			job, err := repositories.JobPostByToken(r.Context(), jobRq.Token)
			if err != nil {
				svr.Log(err, fmt.Sprintf("unable to find job post ID by token: %s", jobRq.Token))
				svr.JSON(w, http.StatusNotFound, nil)
				return
			}
			err = repositories.DeleteJobCascade(r.Context(), job)
			if err != nil {
				svr.Log(err, fmt.Sprintf("unable to permanently delete job: %#v", jobRq))
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			svr.JSON(w, http.StatusOK, nil)
		},
	)
}

func ApproveJobPageHandler(svr server.Server) http.HandlerFunc {
	return middleware.AdminAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			decoder := json.NewDecoder(r.Body)
			jobRq := &job.JobRqUpdate{}
			if err := decoder.Decode(&jobRq); err != nil {
				svr.Log(err, fmt.Sprintf("unable to parse job request for update: %#v", jobRq))
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			job, err := repositories.JobPostByToken(r.Context(), jobRq.Token)
			if err != nil {
				svr.Log(err, fmt.Sprintf("unable to find job post ID by token: %s", jobRq.Token))
				svr.JSON(w, http.StatusNotFound, nil)
				return
			}
			err = repositories.ApproveJob(r.Context(), job.ID)
			if err != nil {
				svr.Log(err, fmt.Sprintf("unable to update job request: %#v", jobRq))
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			err = svr.GetEmail().SendHTMLEmail(
				email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().SupportSenderAddress()},
				email.Address{Email: job.SubscriberEmail},
				email.Address{Name: svr.GetEmail().DefaultSenderName(), Email: svr.GetEmail().SupportSenderAddress()},
				fmt.Sprintf("Your Job Ad on %s", svr.GetConfig().SiteName),
				fmt.Sprintf("Thanks for using %s,\n\nYour Job Ad has been approved and it's currently live on %s: https://%s.\n\nYou can track your Ad performance and renew your Ad via this edit link: https://%s/edit/%s\n.\n\nI am always available to answer any questions you may have,\n\nBest,\n\n%s\n%s", svr.GetConfig().SiteName, svr.GetConfig().SiteName, svr.GetConfig().SiteHost, svr.GetConfig().SiteHost, jobRq.Token, svr.GetConfig().SiteName, svr.GetConfig().AdminEmail),
			)
			if err != nil {
				svr.Log(err, "unable to send email while approving job ad")
			}
			if err := svr.CacheDelete(server.CacheKeyPinnedJobs); err != nil {
				svr.Log(err, "unable to cleanup cache after approving job")
			}
			svr.JSON(w, http.StatusOK, nil)
		},
	)
}

func DisapproveJobPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		jobRq := &job.JobRqUpdate{}
		if err := decoder.Decode(&jobRq); err != nil {
			svr.Log(err, fmt.Sprintf("unable to parse job request for update: %#v", jobRq))
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		job, err := repositories.JobPostByToken(r.Context(), jobRq.Token)
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to find job post ID by token: %s", jobRq.Token))
			svr.JSON(w, http.StatusNotFound, nil)
			return
		}
		err = repositories.DisapproveJob(r.Context(), job.ID)
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to update job request: %#v", jobRq))
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		svr.JSON(w, http.StatusOK, nil)
	}
}

func TrackJobClickoutPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		externalID := vars["id"]
		if externalID == "" {
			svr.Log(errors.New("got empty id for tracking job"), "got empty externalID for tracking")
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		jobPost, err := repositories.GetJobByExternalID(r.Context(), externalID)
		if err != nil {
			svr.Log(err, "unable to get JobID from externalID")
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		if err := repositories.TrackJobClickout(r.Context(), jobPost.ID); err != nil {
			svr.Log(err, fmt.Sprintf("unable to save job clickout for job id %s. %v", jobPost.ID, err))
			svr.JSON(w, http.StatusOK, nil)
			return
		}
		svr.JSON(w, http.StatusOK, nil)
	}
}

func TrackJobClickoutAndRedirectToJobPage(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		externalID := r.URL.Query().Get("j")
		if externalID == "" {
			svr.Log(errors.New("TrackJobClickoutAndRedirectToJobPage: got empty id for tracking job"), "got empty externalID for tracking")
			svr.JSON(w, http.StatusBadRequest, nil)
			return
		}
		reg, _ := regexp.Compile("[^a-zA-Z0-9 ]+")
		jobPost, err := repositories.GetJobByExternalID(r.Context(), reg.ReplaceAllString(externalID, ""))
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to get HowToApply from externalID %s", externalID))
			svr.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		if err := repositories.TrackJobClickout(r.Context(), jobPost.ID); err != nil {
			svr.Log(err, fmt.Sprintf("unable to save job clickout for job id %s. %v", jobPost.ID, err))
			svr.JSON(w, http.StatusOK, nil)
			return
		}
		svr.Redirect(w, r, http.StatusTemporaryRedirect, jobPost.ApplicationLink)
	}
}

func EditJobViewPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		token := vars["token"]
		isCallback := r.URL.Query().Get("callback")
		paymentSuccess := r.URL.Query().Get("payment")
		editToken, err := database.EditTokens(
			qm.Load(
				database.EditTokenRels.Job,
				qm.Load(database.JobRels.JobEvents),
				qm.Load(database.JobRels.PurchaseEvents),
			),
			database.EditTokenWhere.Token.EQ(token),
		).OneG(r.Context())
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to find job post ID by token: %s", token))
			svr.JSON(w, http.StatusNotFound, nil)
			return
		}
		clickoutCount, err := editToken.R.Job.JobEvents(
			database.JobEventWhere.EventType.EQ("clickout"),
		).CountG(r.Context())
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to retrieve job clickout count for job id %s", editToken.JobID))
		}
		viewCount, err := editToken.R.Job.JobEvents(
			database.JobEventWhere.EventType.EQ("page_view"),
		).CountG(r.Context())
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to retrieve job view count for job id %s", editToken.JobID))
		}
		conversionRate := ""
		if clickoutCount > 0 && viewCount > 0 {
			conversionRate = fmt.Sprintf("%.2f", float64(float64(clickoutCount)/float64(viewCount)*100))
		}
		stats, err := GetStatsForJob(svr.Conn, editToken.JobID)
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to retrieve stats for job id %s", editToken.JobID))
		}
		statsSet, err := json.Marshal(stats)
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to marshal stats for job id %s", editToken.JobID))
		}
		svr.Render(r, w, http.StatusOK, "edit.html", map[string]interface{}{
			"Job":                   editToken.R.Job,
			"Stats":                 string(statsSet),
			"Purchases":             editToken.R.Job.R.PurchaseEvents,
			"JobDescriptionEscaped": svr.JSEscapeString(editToken.R.Job.Description),
			"Token":                 token,
			"ViewCount":             viewCount,
			"ClickoutCount":         clickoutCount,
			"ConversionRate":        conversionRate,
			"IsCallback":            isCallback,
			"PaymentSuccess":        paymentSuccess,
			"StripePublishableKey":  svr.GetConfig().StripePublishableKey,
		})
	}
}

type JobStat struct {
	Date      string `json:"date"`
	Clickouts int    `json:"clickouts"`
	PageViews int    `json:"pageviews"`
}

func GetStatsForJob(db *sql.DB, jobID string) ([]JobStat, error) {
	var stats []JobStat
	rows, err := db.Query(`SELECT COUNT(*) FILTER (WHERE event_type = 'clickout') AS clickout, COUNT(*) FILTER (WHERE event_type = 'page_view') AS pageview, TO_CHAR(DATE_TRUNC('day', created_at), 'YYYY-MM-DD') FROM job_event WHERE job_id = $1 GROUP BY DATE_TRUNC('day', created_at) ORDER BY DATE_TRUNC('day', created_at) ASC`, jobID)
	if err == sql.ErrNoRows {
		return stats, nil
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var s JobStat
		if err := rows.Scan(&s.Clickouts, &s.PageViews, &s.Date); err != nil {
			return stats, err
		}
		stats = append(stats, s)
	}

	return stats, nil
}

func ManageJobBySlugViewPageHandler(svr server.Server) http.HandlerFunc {
	return middleware.AdminAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			slug := vars["slug"]
			jobPost, err := repositories.JobPostBySlugAdmin(r.Context(), slug)
			if err != nil {
				svr.JSON(w, http.StatusNotFound, nil)
				return
			}
			jobPostToken, err := repositories.TokenByJobID(r.Context(), jobPost.ID)
			if err != nil {
				svr.JSON(w, http.StatusNotFound, fmt.Sprintf("Job for %s/manage/job/%s not found", svr.GetConfig().SiteHost, slug))
				return
			}
			svr.Redirect(w, r, http.StatusMovedPermanently, fmt.Sprintf("/manage/%s", jobPostToken.Token))
		},
	)
}

func ManageJobViewPageHandler(svr server.Server) http.HandlerFunc {
	return middleware.AdminAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			token := vars["token"]
			job, err := repositories.JobPostByToken(r.Context(), token)
			if err != nil {
				svr.Log(err, fmt.Sprintf("unable to find job post ID by token: %s", token))
				svr.JSON(w, http.StatusNotFound, nil)
				return
			}
			clickoutCount, err := repositories.GetClickoutCountForJob(r.Context(), job)
			if err != nil {
				svr.Log(err, fmt.Sprintf("unable to retrieve job clickout count for job id %s", job.ID))
			}
			viewCount, err := repositories.GetViewCountForJob(r.Context(), job)
			if err != nil {
				svr.Log(err, fmt.Sprintf("unable to retrieve job view count for job id %s", job.ID))
			}
			conversionRate := ""
			if clickoutCount > 0 && viewCount > 0 {
				conversionRate = fmt.Sprintf("%.2f", float64(float64(clickoutCount)/float64(viewCount)*100))
			}
			svr.Render(r, w, http.StatusOK, "manage.html", map[string]interface{}{
				"Job":                   job,
				"JobDescriptionEscaped": svr.JSEscapeString(job.Description),
				"Token":                 token,
				"ViewCount":             viewCount,
				"ClickoutCount":         clickoutCount,
				"ConversionRate":        conversionRate,
			})
		},
	)
}

func GetBlogPostBySlugHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slug := vars["slug"]
		bp, err := repositories.GetBlogPostBySlug(r.Context(), slug)
		if err != nil {
			svr.Log(err, fmt.Sprintf("unable to retrieve blog post: Slug=%s", slug))
			svr.TEXT(w, http.StatusNotFound, "Could not retrieve blogpost. Please try again later.")
			return
		}
		svr.Render(r, w, http.StatusOK, "view-blogpost.html", map[string]interface{}{
			"BlogPost":         bp,
			"BlogPostTextHTML": svr.MarkdownToHTML(bp.Text),
		})
	}
}

func EditBlogPostHandler(svr server.Server) http.HandlerFunc {
	return middleware.AdminAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			id := vars["id"]
			profile, err := middleware.GetUserFromJWT(r, svr.SessionStore, svr.GetJWTSigningKey())
			if err != nil {
				svr.Log(err, "unable to retrieve user from JWT")
				svr.JSON(w, http.StatusForbidden, nil)
				return
			}
			bp, err := repositories.GetBlogPostByIDAndAuthor(r.Context(), id, profile.UserID)
			if err != nil {
				svr.Log(err, fmt.Sprintf("unable to retrieve blog post: ID=%s authorID=%s", id, profile.UserID))
				svr.TEXT(w, http.StatusNotFound, "Could not retrieve blogpost. Please try again later.")
				return
			}
			svr.Render(r, w, http.StatusOK, "edit-blogpost.html", map[string]interface{}{
				"BlogPost":    bp,
				"IsPublished": bp.PublishedAt.Valid,
			})
		},
	)
}

func CreateBlogPostHandler(svr server.Server) http.HandlerFunc {
	return middleware.AdminAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			decoder := json.NewDecoder(r.Body)
			blogRq := &blog.CreateRq{}
			if err := decoder.Decode(&blogRq); err != nil {
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			k, err := ksuid.NewRandom()
			if err != nil {
				svr.Log(err, "unable to generate unique blog post id")
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			blogPostID, err := k.Value()
			if err != nil {
				svr.Log(err, "unable to get blog post id value")
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			blogPostIDStr, ok := blogPostID.(string)
			if !ok {
				svr.Log(err, "unbale to assert blog post id value as string")
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			profile, err := middleware.GetUserFromJWT(r, svr.SessionStore, svr.GetJWTSigningKey())
			if err != nil {
				svr.Log(err, "unable to retrieve user from JWT")
				svr.JSON(w, http.StatusForbidden, nil)
				return
			}
			bp := &database.BlogPost{
				ID:          blogPostIDStr,
				Title:       blogRq.Title,
				Slug:        slug.Make(blogRq.Title),
				Description: blogRq.Description,
				Tags:        blogRq.Tags,
				Text:        blogRq.Text,
				CreatedBy:   profile.UserID,
			}
			if err := repositories.CreateBlogPost(r.Context(), bp); err != nil {
				svr.Log(err, fmt.Sprintf("unable to create blog post: ID=%s authorID=%s", blogPostIDStr, profile.UserID))
				svr.JSON(w, http.StatusInternalServerError, map[string]interface{}{"error": "could not create blog post. Please try again later." + err.Error()})
				return
			}
			svr.JSON(w, http.StatusOK, map[string]interface{}{"id": bp.ID})
		},
	)
}

func CreateDraftBlogPostHandler(svr server.Server) http.HandlerFunc {
	return middleware.AdminAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			svr.Render(r, w, http.StatusOK, "create-blogpost.html", map[string]interface{}{})
		},
	)
}

func UpdateBlogPostHandler(svr server.Server) http.HandlerFunc {
	return middleware.UserAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			decoder := json.NewDecoder(r.Body)
			bpRq := &blog.UpdateRq{}
			if err := decoder.Decode(&bpRq); err != nil {
				svr.JSON(w, http.StatusBadRequest, nil)
				return
			}
			profile, err := middleware.GetUserFromJWT(r, svr.SessionStore, svr.GetJWTSigningKey())
			if err != nil {
				svr.Log(err, "unable to retrieve user from JWT")
				svr.JSON(w, http.StatusForbidden, nil)
				return
			}
			bp := &database.BlogPost{
				ID:          bpRq.ID,
				Title:       bpRq.Title,
				Description: bpRq.Description,
				Tags:        bpRq.Tags,
				Text:        bpRq.Text,
				CreatedBy:   profile.UserID,
			}
			if err := repositories.UpdateBlogPost(r.Context(), bp); err != nil {
				svr.Log(err, fmt.Sprintf("unable to update blog post: ID=%s authorID=%s", bp.ID, profile.UserID))
				svr.JSON(w, http.StatusNotFound, map[string]interface{}{"error": "could not update blog post. Please try again later"})
				return
			}
			svr.JSON(w, http.StatusOK, nil)
		},
	)
}

func PublishBlogPostHandler(svr server.Server) http.HandlerFunc {
	return middleware.UserAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			id := vars["id"]
			profile, err := middleware.GetUserFromJWT(r, svr.SessionStore, svr.GetJWTSigningKey())
			if err != nil {
				svr.Log(err, "unable to get email from JWT")
				svr.JSON(w, http.StatusForbidden, nil)
				return
			}
			bp, err := repositories.GetBlogPostByIDAndAuthor(r.Context(), id, profile.UserID)
			if err != nil {
				svr.Log(err, fmt.Sprintf("unable to unpublish blog post: ID=%s authorID=%s", id, profile.UserID))
				svr.JSON(w, http.StatusInternalServerError, "Could not unpublish blogpost. Please try again later.")
				return
			}
			if err = repositories.PublishBlogPost(r.Context(), bp); err != nil {
				svr.Log(err, fmt.Sprintf("unable to unpublish blog post: ID=%s authorID=%s", id, profile.UserID))
				svr.JSON(w, http.StatusInternalServerError, "Could not unpublish blogpost. Please try again later.")
				return
			}
			svr.JSON(w, http.StatusOK, nil)
		},
	)
}

func UnpublishBlogPostHandler(svr server.Server) http.HandlerFunc {
	return middleware.UserAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			id := vars["id"]
			profile, err := middleware.GetUserFromJWT(r, svr.SessionStore, svr.GetJWTSigningKey())
			if err != nil {
				svr.Log(err, "unable to get email from JWT")
				svr.JSON(w, http.StatusForbidden, nil)
				return
			}
			bp, err := repositories.GetBlogPostByIDAndAuthor(r.Context(), id, profile.UserID)
			if err != nil {
				svr.Log(err, fmt.Sprintf("unable to unpublish blog post: ID=%s authorID=%s", id, profile.UserID))
				svr.JSON(w, http.StatusInternalServerError, "Could not unpublish blogpost. Please try again later.")
				return
			}
			if err := repositories.UnpublishBlogPost(r.Context(), bp); err != nil {
				svr.Log(err, fmt.Sprintf("unable to unpublish blog post: ID=%s authorID=%s", id, profile.UserID))
				svr.JSON(w, http.StatusInternalServerError, "Could not unpublish blogpost. Please try again later.")
				return
			}
			svr.JSON(w, http.StatusOK, nil)
		},
	)
}

func GetAllPublishedBlogPostsHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all, err := repositories.GetAllPublishedBlogPost(r.Context())
		if err != nil {
			svr.Log(err, "unable to retrieve blogposts")
			svr.TEXT(w, http.StatusNotFound, "could not retrieve blog posts. Please try again later")
			return
		}
		fmt.Println("returning all blogposts", len(all))
		svr.Render(r, w, http.StatusOK, "list-blogposts.html", map[string]interface{}{
			"BlogPosts": all,
		})
	}
}

func GetUserBlogPostsHandler(svr server.Server) http.HandlerFunc {
	return middleware.UserAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			profile, err := middleware.GetUserFromJWT(r, svr.SessionStore, svr.GetJWTSigningKey())
			if err != nil {
				svr.Log(err, "unable to get email from JWT")
				svr.JSON(w, http.StatusForbidden, nil)
				return
			}
			all, err := repositories.GetBlogPostByCreatedBy(r.Context(), profile.UserID)
			if err != nil {
				svr.Log(err, "unable to retrieve blogposts")
				svr.TEXT(w, http.StatusNotFound, "could not retrieve blog posts. Please try again later")
				return
			}
			fmt.Println("returning all blogposts", len(all))
			svr.Render(r, w, http.StatusOK, "user-blogposts.html", map[string]interface{}{
				"BlogPosts": all,
			})
		},
	)
}

func ProfileHomepageHandler(svr server.Server) http.HandlerFunc {
	return middleware.UserAuthenticatedMiddleware(
		svr.SessionStore,
		svr.GetJWTSigningKey(),
		func(w http.ResponseWriter, r *http.Request) {
			profile, err := middleware.GetUserFromJWT(r, svr.SessionStore, svr.GetJWTSigningKey())
			if err != nil {
				svr.Log(err, "unable to get email from JWT")
				svr.JSON(w, http.StatusForbidden, nil)
				return
			}
			switch profile.Type {
			case user.UserTypeDeveloper:
				dev, err := repositories.DeveloperProfileByEmail(r.Context(), profile.Email)
				if err != nil {
					svr.Log(err, "unable to find developer profile")
					svr.JSON(w, http.StatusNotFound, nil)
					return
				}
				svr.Render(r, w, http.StatusOK, "profile-home.html", map[string]interface{}{
					"IsAdmin":       profile.IsAdmin,
					"UserID":        profile.UserID,
					"UserEmail":     profile.Email,
					"UserCreatedAt": profile.CreatedAt,
					"ProfileID":     dev.ID,
					"UserType":      profile.Type,
					"Developer":     dev,
				})
			case user.UserTypeRecruiter:
				rec, err := repositories.RecruiterProfileByEmail(r.Context(), profile.Email)
				if err != nil {
					svr.Log(err, "unable to find recruiter profile")
					svr.JSON(w, http.StatusNotFound, nil)
					return
				}
				svr.Render(r, w, http.StatusOK, "profile-home.html", map[string]interface{}{
					"IsAdmin":       profile.IsAdmin,
					"UserID":        profile.UserID,
					"UserEmail":     profile.Email,
					"UserCreatedAt": profile.CreatedAt,
					"ProfileID":     rec.ID,
					"UserType":      profile.Type,
					"Recruiter":     rec,
				})
			case user.UserTypeAdmin:
				dev, err := repositories.DeveloperProfileByEmail(r.Context(), profile.Email)
				if err != nil {
					svr.Log(err, "unable to find developer profile")
					svr.JSON(w, http.StatusNotFound, nil)
					return
				}
				svr.Render(r, w, http.StatusOK, "profile-home.html", map[string]interface{}{
					"IsAdmin":       profile.IsAdmin,
					"UserID":        profile.UserID,
					"UserEmail":     profile.Email,
					"UserCreatedAt": profile.CreatedAt,
					"ProfileID":     dev.ID,
					"UserType":      profile.Type,
					"Developer":     dev,
				})
			}
		},
	)
}
