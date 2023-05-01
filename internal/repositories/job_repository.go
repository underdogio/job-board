package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/segmentio/ksuid"
	"github.com/underdogio/job-board/internal/database"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const (
	jobEventPageView = "page_view"
	jobEventClickout = "clickout"

	SearchTypeJob    = "job"
	SearchTypeSalary = "salary"
)

type JobStat struct {
	Date      string `json:"date"`
	Clickouts int    `json:"clickouts"`
	PageViews int    `json:"pageviews"`
}

func TrackJobView(ctx context.Context, jobID string) error {
	event := &database.JobEvent{
		EventType: jobEventPageView,
		JobID:     jobID,
		CreatedAt: time.Now().UTC(),
	}
	return event.InsertG(context.Background(), boil.Infer())
}

func TrackJobClickout(ctx context.Context, jobID string) error {
	event := &database.JobEvent{
		EventType: jobEventClickout,
		JobID:     jobID,
		CreatedAt: time.Now().UTC(),
	}
	return event.InsertG(context.Background(), boil.Infer())
}

func GetJobByExternalID(ctx context.Context, externalID string) (*database.Job, error) {
	return database.Jobs(
		database.JobWhere.ExternalID.EQ(externalID),
	).OneG(ctx)
}

func SaveDraft(ctx context.Context, job *database.Job) error {
	id, err := ksuid.NewRandom()
	if err != nil {
		return err
	}
	externalID, err := ksuid.NewRandom()
	if err != nil {
		return err
	}
	job.ID = id.String()
	job.ExternalID = externalID.String()
	job.Slug = null.StringFrom(slug.Make(fmt.Sprintf("%s %s %d", job.JobTitle, job.Company, time.Now().UTC().Unix())))
	job.CreatedAt = time.Now().UTC()

	ApplyPlanTypeAndDurationToExpirations(job)
	return job.InsertG(ctx, boil.Infer())
}

func UpdateJob(ctx context.Context, job *database.Job) error {
	_, err := job.UpdateG(ctx, boil.Infer())
	return err
}

func ApproveJob(ctx context.Context, jobID string) error {
	_, err := database.Jobs(
		database.JobWhere.ID.EQ(jobID),
	).UpdateAllG(ctx, database.M{
		database.JobColumns.ApprovedAt: time.Now().UTC(),
	})
	return err
}

func DisapproveJob(ctx context.Context, jobID string) error {
	_, err := database.Jobs(
		database.JobWhere.ID.EQ(jobID),
	).UpdateAllG(ctx, database.M{
		database.JobColumns.ApprovedAt: nil,
	})
	return err
}

func GetViewCountForJob(ctx context.Context, job *database.Job) (int64, error) {
	return job.JobEvents(
		database.JobEventWhere.EventType.EQ(jobEventPageView),
	).CountG(ctx)
}

func GetJobByStripeSessionID(ctx context.Context, sessionID string) (*database.Job, error) {
	return database.Jobs(
		qm.Load(database.JobRels.PurchaseEvents, database.PurchaseEventWhere.StripeSessionID.EQ(sessionID)),
	).OneG(ctx)
}

func GetStatsForJob(ctx context.Context, jobID int) ([]JobStat, error) {
	var stats []JobStat
	database.NewQuery(
		qm.Select("COUNT(*) FILTER (WHERE event_type = 'clickout') AS clickout, COUNT(*) FILTER (WHERE event_type = 'page_view') AS pageview, TO_CHAR(DATE_TRUNC('day', created_at), 'YYYY-MM-DD')"),
		qm.From("job_event"),
		qm.Where("job_id = ?", jobID),
		qm.GroupBy("DATE_TRUNC('day', created_at)"),
		qm.OrderBy("DATE_TRUNC('day', created_at) ASC"),
	).BindG(ctx, stats)

	return stats, nil
}

func JobPostByCreatedAt(ctx context.Context) ([]*database.Job, error) {
	return database.Jobs(
		database.JobWhere.ApprovedAt.IsNotNull(),
		qm.OrderBy("created_at DESC"),
	).AllG(ctx)
}

func TopNJobsByLocation(ctx context.Context, location string, max int) ([]*database.Job, error) {
	return database.Jobs(
		qm.Where("location ILIKE ?", `%`+location+`%`),
		database.JobWhere.ApprovedAt.IsNotNull(),
		qm.OrderBy(database.JobColumns.CreatedAt+" DESC"),
		qm.Limit(max),
	).AllG(ctx)
}

func JobPostBySlug(ctx context.Context, slug string) (*database.Job, error) {
	return database.Jobs(
		database.JobWhere.ApprovedAt.IsNotNull(),
		database.JobWhere.Slug.EQ(null.StringFrom(slug)),
	).OneG(ctx)
}

func JobPostBySlugAdmin(ctx context.Context, slug string) (*database.Job, error) {
	return database.Jobs(
		database.JobWhere.Slug.EQ(null.StringFrom(slug)),
	).OneG(ctx)
}

func JobPostByID(ctx context.Context, jobID string) (*database.Job, error) {
	return database.FindJobG(ctx, jobID)
}

func JobPostByExternalID(ctx context.Context, externalID string) (*database.Job, error) {
	return database.Jobs(
		database.JobWhere.ExternalID.EQ(externalID),
	).OneG(ctx)
}

func JobPostByURLID(ctx context.Context, URLID int) (*database.Job, error) {
	return database.Jobs(
		database.JobWhere.ApprovedAt.IsNotNull(),
		database.JobWhere.URLID.EQ(URLID),
	).OneG(ctx)
}

func DeleteJobCascade(ctx context.Context, job *database.Job) error {

	if _, err := database.Images(
		database.ImageWhere.ID.EQ(job.CompanyIconImageID.String),
	).DeleteAllG(ctx); err != nil {
		return err
	}

	if _, err := job.R.EditTokens.DeleteAllG(ctx); err != nil {
		return err
	}
	if _, err := job.R.ApplyTokens.DeleteAllG(ctx); err != nil {
		return err
	}
	if _, err := job.R.JobEvents.DeleteAllG(ctx); err != nil {
		return err
	}
	if _, err := job.R.PurchaseEvents.DeleteAllG(ctx); err != nil {
		return err
	}
	if _, err := job.DeleteG(ctx); err != nil {
		return err
	}
	return nil
}

func GetPendingJobs(ctx context.Context) ([]*database.Job, error) {
	return database.Jobs(
		database.JobWhere.ApprovedAt.IsNull(),
	).AllG(ctx)
}

// GetCompanyJobs returns jobs for a given company
func GetCompanyJobs(ctx context.Context, companyName string, limit int) ([]*database.Job, error) {
	return database.Jobs(
		database.JobWhere.ApprovedAt.IsNotNull(),
		database.JobWhere.Expired.EQ(null.BoolFrom(false)),
		database.JobWhere.Company.EQ(companyName),
		qm.OrderBy("created_at DESC, approved_at DESC"),
		qm.Limit(limit),
	).AllG(ctx)
}

// GetRelevantJobs returns pinned and most recent jobs for now
func GetRelevantJobs(ctx context.Context, location, jobID string, limit int) ([]*database.Job, error) {
	return database.Jobs(
		database.JobWhere.ApprovedAt.IsNotNull(),
		database.JobWhere.ID.NEQ(jobID),
		database.JobWhere.Expired.EQ(null.BoolFrom(false)),
		qm.OrderBy(database.JobColumns.CreatedAt+" DESC"),
		qm.OrderBy(database.JobColumns.ApprovedAt+" DESC"),
		qm.OrderBy("word_similarity(?, location)", location),
		qm.Limit(limit),
	).AllG(ctx)
}

func GetPinnedJobs(ctx context.Context) ([]*database.Job, error) {
	return database.Jobs(
		database.JobWhere.ApprovedAt.IsNotNull(),
		database.JobWhere.FrontPageEligibilityExpiredAt.GT(null.TimeFrom(time.Now())),
		qm.OrderBy(database.JobColumns.ApprovedAt+" DESC"),
	).AllG(ctx)
}

func JobsByQuery(ctx context.Context, location, tag string, pageId int, salary, currency string, jobsPerPage int, includePinnedJobs bool) ([]*database.Job, int, error) {
	offset := pageId*jobsPerPage - jobsPerPage
	// replace `|` with white space
	// remove double white spaces
	// join with `|` for ps query
	tag = strings.Join(strings.Fields(strings.ReplaceAll(tag, "|", " ")), "|")
	queries := getQueryForArgs(location, tag, salary, currency, offset, jobsPerPage, includePinnedJobs)
	fmt.Println(queries)
	query := database.Jobs(queries...)
	jobs, err := query.AllG(ctx)
	if err != nil {
		return nil, 0, err
	}
	fullRowsCount, err := database.Jobs(queries...).CountG(ctx)
	if err != nil {
		return nil, 0, err
	}

	return jobs, int(fullRowsCount), nil
}

func TokenByJobID(ctx context.Context, jobID string) (*database.EditToken, error) {
	return database.EditTokens(
		database.EditTokenWhere.JobID.EQ(jobID),
	).OneG(ctx)
}

func JobPostByToken(ctx context.Context, token string) (*database.Job, error) {
	return database.Jobs(
		qm.InnerJoin(database.JobRels.EditTokens+" ON "+database.EditTokenTableColumns.JobID+" = "+database.JobTableColumns.ID),
		database.EditTokenWhere.Token.EQ(token),
	).OneG(ctx)
}

func GetLastNJobs(ctx context.Context, max int, loc string) ([]*database.Job, error) {
	if strings.TrimSpace(loc) == "" {
		return database.Jobs(
			database.JobWhere.ApprovedAt.IsNotNull(),
			qm.OrderBy(database.JobColumns.ApprovedAt+" DESC"),
			qm.Limit(max),
		).AllG(ctx)
	} else {
		return database.Jobs(
			database.JobWhere.ApprovedAt.IsNotNull(),
			qm.Where(database.JobColumns.Location+" ILIKE ?", `%`+loc+`%`),
			qm.OrderBy(database.JobColumns.ApprovedAt+" DESC"),
			qm.Limit(max),
		).AllG(ctx)
	}
}

func GetLastNJobsFromID(ctx context.Context, max int, jobID string) ([]*database.Job, error) {
	return database.Jobs(
		database.JobWhere.ID.GT(jobID),
		database.JobWhere.ApprovedAt.IsNotNull(),
		qm.Limit(max),
	).AllG(ctx)
}

func MarkJobAsExpired(ctx context.Context, job *database.Job) error {
	job.Expired = null.BoolFrom(true)
	_, err := job.UpdateG(ctx, boil.Whitelist(database.JobColumns.Expired))
	return err
}

func NewJobsLastWeekOrMonth(ctx context.Context) (int, int, error) {
	week, err := database.Jobs(
		database.JobWhere.ApprovedAt.GTE(null.TimeFrom(time.Now().AddDate(0, 0, -7))),
	).CountG(ctx)
	if err != nil {
		return 0, 0, err
	}
	month, err := database.Jobs(
		database.JobWhere.ApprovedAt.GTE(null.TimeFrom(time.Now().AddDate(0, -1, 0))),
	).CountG(ctx)
	if err != nil {
		return 0, 0, err
	}
	return int(week), int(month), nil
}

func GetJobApplyURLs(ctx context.Context) ([]*database.Job, error) {
	return database.Jobs(
		database.JobWhere.ApprovedAt.IsNotNull(),
		database.JobWhere.Expired.EQ(null.BoolFrom(false)),
	).AllG(ctx)
}

type JobExpirationEntity struct {
	NewsletterEligibilityExpiredAt  time.Time
	BlogEligibilityExpiredAt        time.Time
	SocialMediaEligibilityExpiredAt time.Time
	FrontPageEligibilityExpiredAt   time.Time
	CompanyPageEligibilityExpiredAt time.Time
	PlanExpiredAt                   time.Time
}

func ApplyPlanTypeAndDurationToExpirations(job *database.Job) {
	expiration := null.TimeFrom(time.Now().AddDate(0, 0, 365*10))
	job.NewsletterEligibilityExpiredAt = expiration
	job.BlogEligibilityExpiredAt = expiration
	job.SocialMediaEligibilityExpiredAt = expiration
	job.FrontPageEligibilityExpiredAt = expiration
	job.CompanyPageEligibilityExpiredAt = expiration
	job.PlanExpiredAt = expiration
}

func PlanTypeAndDurationToExpirationsFromExistingExpirations(expiration JobExpirationEntity, planType string, planDuration int) (JobExpirationEntity, error) {
	return JobExpirationEntity{
		NewsletterEligibilityExpiredAt:  expiration.NewsletterEligibilityExpiredAt.AddDate(0, 0, 30*planDuration),
		BlogEligibilityExpiredAt:        expiration.BlogEligibilityExpiredAt.AddDate(0, 0, 30*planDuration),
		SocialMediaEligibilityExpiredAt: expiration.SocialMediaEligibilityExpiredAt.AddDate(0, 0, 30*planDuration),
		FrontPageEligibilityExpiredAt:   expiration.FrontPageEligibilityExpiredAt.AddDate(0, 0, 30*planDuration),
		CompanyPageEligibilityExpiredAt: expiration.CompanyPageEligibilityExpiredAt.AddDate(0, 0, 30*planDuration),
		PlanExpiredAt:                   expiration.PlanExpiredAt.AddDate(0, 0, 30*planDuration),
	}, nil
}

func GetClickoutCountForJob(ctx context.Context, job *database.Job) (int64, error) {
	return job.JobEvents(
		database.JobEventWhere.EventType.EQ("clickout"),
	).CountG(ctx)
}

func LastJobPosted(ctx context.Context) (time.Time, error) {
	lastJob, err := database.Jobs(
		database.JobWhere.ApprovedAt.IsNotNull(),
		qm.OrderBy("created_at DESC"),
	).OneG(ctx)
	if err != nil {
		return time.Time{}, err
	}
	return lastJob.CreatedAt, nil
}

func SaveTokenForJob(ctx context.Context, token, jobID string) error {
	tokenObj := &database.EditToken{
		Token:     token,
		JobID:     jobID,
		CreatedAt: time.Now().UTC(),
	}
	return tokenObj.InsertG(ctx, boil.Infer())
}

func GetValue(ctx context.Context, key string) (*database.Metum, error) {
	return database.Meta(
		database.MetumWhere.Key.EQ(key),
	).OneG(ctx)
}

func SetValue(ctx context.Context, key, val string) error {
	_, err := database.Meta(
		database.MetumWhere.Key.EQ(key),
	).UpdateAllG(ctx, database.M{
		database.MetumColumns.Value: val,
	})
	return err
}

func ApplyToJob(ctx context.Context, jobID string, cv []byte, email, token string) error {
	applyToken := &database.ApplyToken{
		Token:     token,
		JobID:     jobID,
		Email:     email,
		CV:        cv,
		CreatedAt: time.Now().UTC(),
	}
	return applyToken.InsertG(ctx, boil.Infer())
}

func ConfirmApplyToJob(ctx context.Context, token string) error {
	_, err := database.ApplyTokens(
		database.ApplyTokenWhere.Token.EQ(token),
	).UpdateAllG(ctx, database.M{
		database.ApplyTokenColumns.ConfirmedAt: time.Now().UTC(),
	})

	return err
}

func CleanupExpiredApplyTokens(ctx context.Context) error {
	_, err := database.ApplyTokens(
		database.ApplyTokenWhere.CreatedAt.LT(time.Now().UTC().AddDate(0, 0, -3)),
		qm.Or2(database.ApplyTokenWhere.ConfirmedAt.IsNull()),
	).DeleteAllG(ctx)
	return err
}

func salaryToSalaryRangeString(salaryMin, salaryMax int, currency string) string {
	salaryMinStr := fmt.Sprintf("%d", salaryMin)
	salaryMaxStr := fmt.Sprintf("%d", salaryMax)
	if currency != "₹" {
		if salaryMin > 1000 {
			salaryMinStr = fmt.Sprintf("%dk", salaryMin/1000)
		}
		if salaryMax > 1000 {
			salaryMaxStr = fmt.Sprintf("%dk", salaryMax/1000)
		}
	} else {
		if salaryMin > 100000 {
			salaryMinStr = fmt.Sprintf("%dL", salaryMin/100000)
		}
		if salaryMax > 100000 {
			salaryMaxStr = fmt.Sprintf("%dL", salaryMax/100000)
		}
	}

	return fmt.Sprintf("%s%s - %s%s", currency, salaryMinStr, currency, salaryMaxStr)
}

func getQueryForArgs(location, tag string, salary, currency string, offset, max int, includePinnedJobs bool) []qm.QueryMod {
	query := []qm.QueryMod{}
	if tag != "" {
		query = append(
			query,
			qm.Select("to_tsvector(job_title || ' ' || company || ' ' || description) as vector, ts_rank_cd(vector, query, 32 /* rank/(rank+1) */ ) as rank, *"),
			qm.From(fmt.Sprintf("to_tsquery('%s') query", tag)),
			qm.Where("query @@ vector"),
			qm.OrderBy("rank DESC"),
		)
	}
	if !includePinnedJobs {
		query = append(query, database.JobWhere.FrontPageEligibilityExpiredAt.LT(null.TimeFrom(time.Now())))
	}
	if location != "" {
		query = append(query, qm.Where("location ILIKE '%' || ? || '%'", location))
	}
	if salary != "" {
		query = append(query, database.JobWhere.SalaryRange.EQ(salary))
	}
	query = append(
		query,
		database.JobWhere.ApprovedAt.IsNotNull(),
		qm.OrderBy("created_at DESC"),
		qm.GroupBy("id, created_at"),
		qm.Limit(max),
		qm.Offset(offset),
	)
	return query
}
