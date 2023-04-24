package job

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/segmentio/ksuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) TrackJobView(job *JobPost) error {
	stmt := `INSERT INTO job_event (event_type, job_id, created_at) VALUES ($1, $2, NOW())`
	_, err := r.db.Exec(stmt, jobEventPageView, job.ID)
	return err
}

func (r *Repository) TrackJobClickout(jobID int) error {
	stmt := `INSERT INTO job_event (event_type, job_id, created_at) VALUES ($1, $2, NOW())`
	_, err := r.db.Exec(stmt, jobEventClickout, jobID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetJobByExternalID(externalID string) (JobPost, error) {
	res := r.db.QueryRow(`SELECT id, job_title, company, salary_range, location, slug, external_id FROM job WHERE external_id = $1`, externalID)
	var job JobPost
	err := res.Scan(&job.ID, &job.JobTitle, &job.Company, &job.SalaryRange, &job.Location, &job.Slug, &job.ExternalID)
	if err != nil {
		return job, err
	}

	return job, nil
}

func (r *Repository) SaveDraft(job *JobRq) (int, error) {
	externalID, err := ksuid.NewRandom()
	if err != nil {
		return 0, err
	}
	expiration := r.PlanTypeAndDurationToExpirations()
	sqlStatement := `
			INSERT INTO job (
				job_title,
				job_category,
				company,
				location,
				salary_range,
				job_type,
				application_link,
				description,
				created_at,
				url_id,
				slug,
				company_icon_image_id,
				external_id,
				blog_eligibility_expired_at,
				company_page_eligibility_expired_at,
				front_page_eligibility_expired_at,
				newsletter_eligibility_expired_at,
				plan_expired_at,
				social_media_eligibility_expired_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19) RETURNING id`
	slugTitle := slug.Make(fmt.Sprintf("%s %s %d", job.JobTitle, job.Company, time.Now().UTC().Unix()))
	createdAt := time.Now().UTC().Unix()

	var lastInsertID int
	res := r.db.QueryRow(
		sqlStatement,
		job.JobTitle,                        // job_title
		job.JobCategory,                     // job_category
		job.Company,                         // company
		job.Location,                        // location
		job.SalaryRange,                     // salary_range
		job.JobType,                         // job_type
		job.ApplicationLink,                 // application_link
		job.Description,                     // description
		time.Unix(createdAt, 0),             // created_at
		createdAt,                           // url_id
		slugTitle,                           // slug
		job.CompanyIconID,                   // company_icon_image_id
		externalID,                          // external_id
		expiration.BlogEligibilityExpiredAt, // blog_eligibility_expired_at
		expiration.CompanyPageEligibilityExpiredAt, // company_page_eligibility_expired_at
		expiration.FrontPageEligibilityExpiredAt,   // front_page_eligibility_expired_at
		expiration.NewsletterEligibilityExpiredAt,  // newsletter_eligibility_expired_at
		expiration.PlanExpiredAt,                   // plan_expired_at
		expiration.SocialMediaEligibilityExpiredAt, // social_media_eligibility_expired_at
	)

	if err := res.Scan(&lastInsertID); err != nil {
		return 0, err
	}
	return int(lastInsertID), err
}

func (r *Repository) UpdateJob(job *JobRqUpdate, jobID int) error {
	salaryMinInt, err := strconv.Atoi(strings.TrimSpace(job.SalaryMin))
	if err != nil {
		return err
	}
	salaryMaxInt, err := strconv.Atoi(strings.TrimSpace(job.SalaryMax))
	if err != nil {
		return err
	}
	salaryRange := salaryToSalaryRangeString(salaryMinInt, salaryMaxInt, job.SalaryCurrency)
	_, err = r.db.Exec(
		`UPDATE job SET job_title = $1, company = $2, company_url = $3, salary_min = $4, salary_max = $5, salary_currency = $6, salary_range = $7, location = $8, description = $9, perks = $10, interview_process = $11, how_to_apply = $12, company_icon_image_id = $13 WHERE id = $14`,
		job.JobTitle,
		job.Company,
		job.CompanyURL,
		job.SalaryMin,
		job.SalaryMax,
		job.SalaryCurrency,
		salaryRange,
		job.Location,
		job.Description,
		job.Perks,
		job.InterviewProcess,
		job.HowToApply,
		job.CompanyIconID,
		jobID,
	)
	if err != nil {
		return err
	}
	return err
}

func (r *Repository) ApproveJob(jobID int) error {
	_, err := r.db.Exec(
		`UPDATE job SET approved_at = NOW() WHERE id = $1`,
		jobID,
	)
	if err != nil {
		return err
	}
	return err
}

func (r *Repository) DisapproveJob(jobID int) error {
	_, err := r.db.Exec(
		`UPDATE job SET approved_at = NULL WHERE id = $1`,
		jobID,
	)
	if err != nil {
		return err
	}
	return err
}

func (r *Repository) GetViewCountForJob(jobID int) (int, error) {
	var count int
	row := r.db.QueryRow(`select count(*) as c from job_event where job_event.event_type = 'page_view' and job_event.job_id = $1`, jobID)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}

func (r *Repository) GetJobByStripeSessionID(sessionID string) (JobPost, error) {
	res := r.db.QueryRow(`SELECT j.id, j.job_title, j.company, j.salary_range, j.location, j.slug, j.external_id, j.approved_at, j.blog_eligibility_expired_at, j.company_page_eligibility_expired_at, j.front_page_eligibility_expired_at, j.newsletter_eligibility_expired_at, j.plan_expired_at, j.social_media_eligibility_expired_at FROM purchase_event p LEFT JOIN job j ON p.job_id = j.id WHERE p.stripe_session_id = $1`, sessionID)
	var job JobPost
	var approvedAt sql.NullTime
	err := res.Scan(
		&job.ID,
		&job.JobTitle,
		&job.Company,
		&job.SalaryRange,
		&job.Location,
		&job.Slug,
		&job.ExternalID,
		&approvedAt,
		&job.BlogEligibilityExpiredAt,
		&job.CompanyPageEligibilityExpiredAt,
		&job.FrontPageEligibilityExpiredAt,
		&job.NewsletterEligibilityExpiredAt,
		&job.PlanExpiredAt,
		&job.SocialMediaEligibilityExpiredAt,
	)
	if err != nil {
		return job, err
	}
	if approvedAt.Valid {
		job.ApprovedAt = &approvedAt.Time
	}

	return job, nil
}

func (r *Repository) GetStatsForJob(jobID int) ([]JobStat, error) {
	var stats []JobStat
	rows, err := r.db.Query(`SELECT COUNT(*) FILTER (WHERE event_type = 'clickout') AS clickout, COUNT(*) FILTER (WHERE event_type = 'page_view') AS pageview, TO_CHAR(DATE_TRUNC('day', created_at), 'YYYY-MM-DD') FROM job_event WHERE job_id = $1 GROUP BY DATE_TRUNC('day', created_at) ORDER BY DATE_TRUNC('day', created_at) ASC`, jobID)
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

func (r *Repository) JobPostByCreatedAt() ([]*JobPost, error) {
	jobs := []*JobPost{}
	var rows *sql.Rows
	rows, err := r.db.Query(
		`SELECT id, job_title, company, salary_range, location, description, created_at, url_id, slug, company_icon_image_id, external_id, expired, last_week_clickouts
		FROM job
		WHERE approved_at IS NOT NULL
		ORDER BY created_at DESC`)
	if err != nil {
		return jobs, err
	}
	for rows.Next() {
		job := &JobPost{}
		var createdAt time.Time
		var companyIcon sql.NullString
		err = rows.Scan(&job.ID, &job.JobTitle, &job.Company, &job.SalaryRange, &job.Location, &job.JobDescription, &createdAt, &job.CreatedAt, &job.Slug, &companyIcon, &job.ExternalID, &job.Expired, &job.LastWeekClickouts)
		if companyIcon.Valid {
			job.CompanyIconID = companyIcon.String
		}
		job.TimeAgo = createdAt.UTC().Format("January 2006")
		if err != nil {
			return jobs, err
		}
		jobs = append(jobs, job)
	}
	err = rows.Err()
	if err != nil {
		return jobs, err
	}
	return jobs, nil
}

func (r *Repository) TopNJobsByLocation(location string, max int) ([]*JobPost, error) {
	jobs := []*JobPost{}
	var rows *sql.Rows
	rows, err := r.db.Query(
		`SELECT id, job_title, company, salary_range, location, description, created_at, url_id, slug, company_icon_image_id, external_id
		FROM job
		WHERE location ILIKE '%' || $1 || '%'
		AND approved_at IS NOT NULL
		ORDER BY created_at DESC LIMIT $2`, location, max)
	if err != nil {
		return jobs, err
	}
	for rows.Next() {
		job := &JobPost{}
		var createdAt time.Time
		var companyIcon sql.NullString
		err = rows.Scan(&job.ID, &job.JobTitle, &job.Company, &job.SalaryRange, &job.Location, &job.JobDescription, &createdAt, &job.CreatedAt, &job.Slug, &companyIcon, &job.ExternalID)
		if companyIcon.Valid {
			job.CompanyIconID = companyIcon.String
		}
		job.TimeAgo = createdAt.UTC().Format("January 2006")
		if err != nil {
			return jobs, err
		}
		jobs = append(jobs, job)
	}
	err = rows.Err()
	if err != nil {
		return jobs, err
	}
	return jobs, nil
}

func (r *Repository) JobPostBySlug(slug string) (*JobPost, error) {
	job := &JobPost{}
	row := r.db.QueryRow(
		`SELECT id, job_title, company, salary_range, location, description, created_at, url_id, slug, company_icon_image_id, external_id, expired, last_week_clickouts
		FROM job
		WHERE approved_at IS NOT NULL
		AND slug = $1`, slug)
	var createdAt time.Time
	var companyIcon sql.NullString
	err := row.Scan(&job.ID, &job.JobTitle, &job.Company, &job.SalaryRange, &job.Location, &job.JobDescription, &createdAt, &job.CreatedAt, &job.Slug, &companyIcon, &job.ExternalID, &job.Expired, &job.LastWeekClickouts)
	if companyIcon.Valid {
		job.CompanyIconID = companyIcon.String
	}
	if err != nil {
		return job, err
	}
	job.TimeAgo = createdAt.UTC().Format("January 2006")
	return job, nil
}

func (r *Repository) JobPostBySlugAdmin(slug string) (*JobPost, error) {
	job := &JobPost{}
	row := r.db.QueryRow(
		`SELECT id, job_title, company, salary_range, location, description, created_at, url_id, slug, company_icon_image_id, external_id
		FROM job
		WHERE slug = $1`, slug)
	var createdAt time.Time
	var companyIcon sql.NullString
	err := row.Scan(&job.ID, &job.JobTitle, &job.Company, &job.SalaryRange, &job.Location, &job.JobDescription, &createdAt, &job.CreatedAt, &job.Slug, &companyIcon, &job.ExternalID)
	if companyIcon.Valid {
		job.CompanyIconID = companyIcon.String
	}
	if err != nil {
		return job, err
	}
	job.TimeAgo = createdAt.UTC().Format("January 2006")
	return job, nil
}

func (r *Repository) JobPostByIDForEdit(jobID int) (*JobPostForEdit, error) {
	job := &JobPostForEdit{}
	row := r.db.QueryRow(
		`SELECT job_title, company, company_email, company_url, salary_min, salary_max, salary_currency, location, description, perks, interview_process, how_to_apply, created_at, slug, approved_at, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period, plan_type, plan_duration, blog_eligibility_expired_at, company_page_eligibility_expired_at, front_page_eligibility_expired_at, newsletter_eligibility_expired_at, plan_expired_at, social_media_eligibility_expired_at
		FROM job
		WHERE id = $1`, jobID)
	var perks, interview, companyURL, companyIconID sql.NullString
	err := row.Scan(
		&job.JobTitle,
		&job.Company,
		&job.CompanyEmail,
		&companyURL,
		&job.SalaryMin,
		&job.SalaryMax,
		&job.SalaryCurrency,
		&job.Location,
		&job.JobDescription,
		&perks,
		&interview,
		&job.HowToApply,
		&job.CreatedAt,
		&job.Slug,
		&job.ApprovedAt,
		&job.SalaryMin,
		&job.SalaryMax,
		&job.SalaryCurrency,
		&companyIconID,
		&job.ExternalID,
		&job.SalaryPeriod,
		&job.PlanType,
		&job.PlanDuration,
		&job.BlogEligibilityExpiredAt,
		&job.CompanyPageEligibilityExpiredAt,
		&job.FrontPageEligibilityExpiredAt,
		&job.NewsletterEligibilityExpiredAt,
		&job.PlanExpiredAt,
		&job.SocialMediaEligibilityExpiredAt,
	)
	if err != nil {
		return job, err
	}
	if companyIconID.Valid {
		job.CompanyIconID = companyIconID.String
	}
	if perks.Valid {
		job.Perks = perks.String
	}
	if interview.Valid {
		job.InterviewProcess = interview.String
	}
	if companyURL.Valid {
		job.CompanyURL = companyURL.String
	} else {
		job.CompanyURL = ""
	}
	return job, nil
}

func (r *Repository) JobPostByExternalIDForEdit(externalID string) (*JobPostForEdit, error) {
	job := &JobPostForEdit{}
	row := r.db.QueryRow(
		`SELECT id, job_title, company, company_email, company_url, salary_min, salary_max, salary_currency, location, description, perks, interview_process, how_to_apply, created_at, slug, approved_at, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period
		FROM job
		WHERE external_id = $1`, externalID)
	var perks, interview, companyURL, companyIconID sql.NullString
	err := row.Scan(&job.ID, &job.JobTitle, &job.Company, &job.CompanyEmail, &companyURL, &job.SalaryMin, &job.SalaryMax, &job.SalaryCurrency, &job.Location, &job.JobDescription, &perks, &interview, &job.HowToApply, &job.CreatedAt, &job.Slug, &job.ApprovedAt, &job.SalaryMin, &job.SalaryMax, &job.SalaryCurrency, &companyIconID, &job.ExternalID, &job.SalaryPeriod)
	if err != nil {
		return job, err
	}
	if companyIconID.Valid {
		job.CompanyIconID = companyIconID.String
	}
	if perks.Valid {
		job.Perks = perks.String
	}
	if interview.Valid {
		job.InterviewProcess = interview.String
	}
	if companyURL.Valid {
		job.CompanyURL = companyURL.String
	} else {
		job.CompanyURL = ""
	}
	return job, nil
}

func (r *Repository) JobPostByURLID(URLID int64) (*JobPost, error) {
	job := &JobPost{}
	row := r.db.QueryRow(
		`SELECT id, job_title, company, salary_range, location, description, created_at, url_id, slug, company_icon_image_id, external_id, expired, last_week_clickouts
		FROM job
		WHERE approved_at IS NOT NULL
		AND url_id = $1`, URLID)
	var createdAt time.Time
	var companyIcon sql.NullString
	err := row.Scan(&job.ID, &job.JobTitle, &job.Company, &job.SalaryRange, &job.Location, &job.JobDescription, &createdAt, &job.CreatedAt, &job.Slug, &companyIcon, &job.ExternalID, &job.Expired, &job.LastWeekClickouts)
	if err != nil {
		return job, err
	}
	if companyIcon.Valid {
		job.CompanyIconID = companyIcon.String
	}
	job.TimeAgo = createdAt.UTC().Format("January 2006")
	return job, nil
}

func (r *Repository) DeleteJobCascade(jobID int) error {
	if _, err := r.db.Exec(
		`DELETE FROM image WHERE id IN (SELECT company_icon_image_id FROM job WHERE id = $1)`,
		jobID,
	); err != nil {
		return err
	}
	if _, err := r.db.Exec(
		`DELETE FROM edit_token WHERE job_id = $1`,
		jobID,
	); err != nil {
		return err
	}
	if _, err := r.db.Exec(
		`DELETE FROM apply_token WHERE job_id = $1`,
		jobID,
	); err != nil {
		return err
	}
	if _, err := r.db.Exec(
		`DELETE FROM job_event WHERE job_id = $1`,
		jobID,
	); err != nil {
		return err
	}
	if _, err := r.db.Exec(
		`DELETE FROM purchase_event WHERE job_id = $1`,
		jobID,
	); err != nil {
		return err
	}
	if _, err := r.db.Exec(
		`DELETE FROM job WHERE id = $1`,
		jobID,
	); err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetPendingJobs() ([]*JobPost, error) {
	jobs := []*JobPost{}
	var rows *sql.Rows
	rows, err := r.db.Query(`
	SELECT id, job_title, company, company_url, salary_range, location, description, perks, interview_process, how_to_apply, created_at, url_id, slug, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period
		FROM job WHERE approved_at IS NULL`)
	if err == sql.ErrNoRows {
		return jobs, nil
	}
	if err != nil {
		return jobs, err
	}
	defer rows.Close()
	for rows.Next() {
		job := &JobPost{}
		var createdAt time.Time
		var companyIcon sql.NullString
		err = rows.Scan(&job.ID, &job.JobTitle, &job.Company, &job.SalaryRange, &job.Location, &job.JobDescription, &createdAt, &job.CreatedAt, &job.Slug, &companyIcon, &job.ExternalID)
		if companyIcon.Valid {
			job.CompanyIconID = companyIcon.String
		}
		job.TimeAgo = createdAt.UTC().Format("January 2006")
		if err != nil {
			return jobs, err
		}
		jobs = append(jobs, job)
	}
	err = rows.Err()
	if err != nil {
		return jobs, err
	}
	return jobs, nil
}

// GetCompanyJobs returns jobs for a given company
func (r *Repository) GetCompanyJobs(companyName string, limit int) ([]*JobPost, error) {
	jobs := []*JobPost{}
	var rows *sql.Rows
	rows, err := r.db.Query(`
	SELECT id, job_title, company, company_url, salary_range, location, description, perks, interview_process, how_to_apply, created_at, url_id, slug, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period, last_week_clickouts
		FROM job WHERE approved_at IS NOT NULL AND expired IS FALSE AND company = $1 ORDER BY created_at DESC, approved_at DESC LIMIT $2`, companyName, limit)
	if err != nil {
		return jobs, err
	}
	defer rows.Close()
	for rows.Next() {
		job := &JobPost{}
		var createdAt time.Time
		var companyIcon sql.NullString
		err = rows.Scan(&job.ID, &job.JobTitle, &job.Company, &job.SalaryRange, &job.Location, &job.JobDescription, &createdAt, &job.CreatedAt, &job.Slug, &companyIcon, &job.ExternalID, &job.LastWeekClickouts)
		if companyIcon.Valid {
			job.CompanyIconID = companyIcon.String
		}
		job.TimeAgo = createdAt.UTC().Format("January 2006")
		if err != nil {
			return jobs, err
		}
		jobs = append(jobs, job)
	}
	err = rows.Err()
	if err != nil {
		return jobs, err
	}
	return jobs, nil
}

// GetRelevantJobs returns pinned and most recent jobs for now
func (r *Repository) GetRelevantJobs(location string, jobID int, limit int) ([]*JobPost, error) {
	jobs := []*JobPost{}
	var rows *sql.Rows
	rows, err := r.db.Query(`
	SELECT id, job_title, company, company_url, salary_range, location, description, perks, interview_process, how_to_apply, created_at, url_id, slug, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period, last_week_clickouts
		FROM job WHERE approved_at IS NOT NULL AND id != $1 AND expired IS FALSE ORDER BY created_at DESC, approved_at DESC, word_similarity($2, location) LIMIT $3`, jobID, location, limit)
	if err != nil {
		return jobs, err
	}
	defer rows.Close()
	for rows.Next() {
		job := &JobPost{}
		var createdAt time.Time
		var companyIcon sql.NullString
		err = rows.Scan(&job.ID, &job.JobTitle, &job.Company, &job.SalaryRange, &job.Location, &job.JobDescription, &createdAt, &job.CreatedAt, &job.Slug, &companyIcon, &job.ExternalID, &job.LastWeekClickouts)
		if companyIcon.Valid {
			job.CompanyIconID = companyIcon.String
		}
		job.TimeAgo = createdAt.UTC().Format("January 2006")
		if err != nil {
			return jobs, err
		}
		jobs = append(jobs, job)
	}
	err = rows.Err()
	if err != nil {
		return jobs, err
	}
	return jobs, nil
}

func (r *Repository) GetPinnedJobs() ([]*JobPost, error) {
	jobs := []*JobPost{}
	var rows *sql.Rows
	rows, err := r.db.Query(`
	SELECT id, job_title, company, salary_range, location, description, perks, interview_process, how_to_apply, created_at, url_id, slug, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period, last_week_clickouts
		FROM job WHERE approved_at IS NOT NULL AND front_page_eligibility_expired_at > NOW() ORDER BY approved_at DESC`)
	if err != nil {
		return jobs, err
	}
	defer rows.Close()
	for rows.Next() {
		job := &JobPost{}
		var createdAt time.Time
		var perks, interview, companyIcon sql.NullString
		err = rows.Scan(&job.ID, &job.JobTitle, &job.Company, &job.SalaryRange, &job.Location, &job.JobDescription, &perks, &interview, &createdAt, &job.CreatedAt, &job.Slug, &companyIcon, &job.ExternalID, &job.LastWeekClickouts)
		if companyIcon.Valid {
			job.CompanyIconID = companyIcon.String
		}
		job.TimeAgo = createdAt.UTC().Format("January 2006")
		if err != nil {
			return jobs, err
		}
		jobs = append(jobs, job)
	}
	err = rows.Err()
	if err != nil {
		return jobs, err
	}
	return jobs, nil
}

func (r *Repository) JobsByQuery(location, tag string, pageId, salary int, currency string, jobsPerPage int, includePinnedJobs bool) ([]*JobPost, int, error) {
	jobs := []*JobPost{}
	var rows *sql.Rows
	offset := pageId*jobsPerPage - jobsPerPage
	// replace `|` with white space
	// remove double white spaces
	// join with `|` for ps query
	tag = strings.Join(strings.Fields(strings.ReplaceAll(tag, "|", " ")), "|")
	rows, err := getQueryForArgs(r.db, location, tag, salary, currency, offset, jobsPerPage, includePinnedJobs)
	if err != nil {
		return jobs, 0, err
	}
	defer rows.Close()
	var fullRowsCount int
	for rows.Next() {
		job := &JobPost{}
		var createdAt time.Time
		var companyIcon sql.NullString
		err = rows.Scan(
			&fullRowsCount,
			&job.ID,
			&job.JobTitle,
			&job.Company,
			&job.SalaryRange,
			&job.Location,
			&job.JobDescription,
			&createdAt,
			&job.CreatedAt,
			&job.Slug,
			&companyIcon,
			&job.ExternalID,
			&job.Expired,
			&job.LastWeekClickouts,
			&job.BlogEligibilityExpiredAt,
			&job.CompanyPageEligibilityExpiredAt,
			&job.FrontPageEligibilityExpiredAt,
			&job.NewsletterEligibilityExpiredAt,
			&job.PlanExpiredAt,
			&job.SocialMediaEligibilityExpiredAt,
		)
		if companyIcon.Valid {
			job.CompanyIconID = companyIcon.String
		}
		job.TimeAgo = createdAt.UTC().Format("January 2006")
		if err != nil {
			return jobs, fullRowsCount, err
		}
		jobs = append(jobs, job)
	}
	err = rows.Err()
	if err != nil {
		return jobs, fullRowsCount, err
	}
	return jobs, fullRowsCount, nil
}

func (r *Repository) TokenByJobID(jobID int) (string, error) {
	tokenRow := r.db.QueryRow(
		`SELECT token
		FROM edit_token
		WHERE job_id = $1`, jobID)
	var token string
	err := tokenRow.Scan(&token)
	return token, err
}

func (r *Repository) JobPostIDByToken(token string) (int, error) {
	row := r.db.QueryRow(
		`SELECT job_id
		FROM edit_token
		WHERE token = $1`, token)
	var jobID int
	err := row.Scan(&jobID)
	if err != nil {
		return 0, err
	}
	return jobID, nil
}

func (r *Repository) GetLastNJobs(max int, loc string) ([]*JobPost, error) {
	var jobs []*JobPost
	var rows *sql.Rows
	var err error
	if strings.TrimSpace(loc) == "" {
		rows, err = r.db.Query(`SELECT id, job_title, description, company, salary_range, location, slug, company_icon_image_id, external_id, approved_at FROM job WHERE approved_at IS NOT NULL ORDER BY approved_at DESC LIMIT $1`, max)
	} else {
		rows, err = r.db.Query(`
	SELECT id, job_title, description, company, salary_range, location, slug, company_icon_image_id, external_id, approved_at
	FROM
	job
	WHERE
	approved_at IS NOT NULL
	AND location ILIKE '%' || $1 || '%'
	ORDER BY approved_at DESC LIMIT $2`, loc, max)
	}
	if err != nil {
		return jobs, err
	}
	for rows.Next() {
		job := &JobPost{}
		var companyIcon sql.NullString
		err := rows.Scan(&job.ID, &job.JobTitle, &job.JobDescription, &job.Company, &job.SalaryRange, &job.Location, &job.Slug, &companyIcon, &job.ExternalID, &job.ApprovedAt)
		if companyIcon.Valid {
			job.CompanyIconID = companyIcon.String
		}
		if err != nil {
			return jobs, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *Repository) GetLastNJobsFromID(max, jobID int) ([]*JobPost, error) {
	var jobs []*JobPost
	var rows *sql.Rows
	rows, err := r.db.Query(`SELECT id, job_title, company, salary_range, location, slug, company_icon_image_id, external_id FROM job WHERE id > $1 AND approved_at IS NOT NULL LIMIT $2`, jobID, max)
	if err != nil {
		return jobs, err
	}
	for rows.Next() {
		job := &JobPost{}
		var companyIcon sql.NullString
		err := rows.Scan(&job.ID, &job.JobTitle, &job.Company, &job.SalaryRange, &job.Location, &job.Slug, &companyIcon, &job.ExternalID)
		if companyIcon.Valid {
			job.CompanyIconID = companyIcon.String
		}
		if err != nil {
			return jobs, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *Repository) MarkJobAsExpired(jobID int) error {
	_, err := r.db.Exec(`UPDATE job SET expired = true WHERE id = $1`, jobID)
	return err
}

func (r *Repository) NewJobsLastWeekOrMonth() (int, int, error) {
	var week, month int
	row := r.db.QueryRow(`select lastweek.c as week, lastmonth.c as month 
from 
(select count(*) as c, 1 as id from job  where approved_at >= (now() - '7 days'::interval)::date) as lastweek
left join 
(select count(*) as c, 1 as id from job  where approved_at >= (now() - '30 days'::interval)::date) as lastmonth on lastmonth.id = lastweek.id 
`)
	if err := row.Scan(&week, &month); err != nil {
		return week, month, err
	}
	return week, month, nil
}

func (r *Repository) GetJobApplyURLs() ([]JobApplyURL, error) {
	jobURLs := make([]JobApplyURL, 0)
	var rows *sql.Rows
	rows, err := r.db.Query(`SELECT id, how_to_apply FROM job WHERE approved_at IS NOT NULL AND expired = false`)
	if err != nil {
		return jobURLs, err
	}
	for rows.Next() {
		jobURL := JobApplyURL{}
		if err := rows.Scan(&jobURL.ID, &jobURL.URL); err != nil {
			return jobURLs, err
		}
		jobURLs = append(jobURLs, jobURL)
	}
	return jobURLs, nil
}

type JobExpirationEntity struct {
	NewsletterEligibilityExpiredAt  time.Time
	BlogEligibilityExpiredAt        time.Time
	SocialMediaEligibilityExpiredAt time.Time
	FrontPageEligibilityExpiredAt   time.Time
	CompanyPageEligibilityExpiredAt time.Time
	PlanExpiredAt                   time.Time
}

func (r *Repository) PlanTypeAndDurationToExpirations() JobExpirationEntity {
	expiration := time.Now().AddDate(0, 0, 365*10)
	return JobExpirationEntity{
		NewsletterEligibilityExpiredAt:  expiration,
		BlogEligibilityExpiredAt:        expiration,
		SocialMediaEligibilityExpiredAt: expiration,
		FrontPageEligibilityExpiredAt:   expiration,
		CompanyPageEligibilityExpiredAt: expiration,
		PlanExpiredAt:                   expiration,
	}
}

func (r *Repository) PlanTypeAndDurationToExpirationsFromExistingExpirations(expiration JobExpirationEntity, planType string, planDuration int) (JobExpirationEntity, error) {
	return JobExpirationEntity{
		NewsletterEligibilityExpiredAt:  expiration.NewsletterEligibilityExpiredAt.AddDate(0, 0, 30*planDuration),
		BlogEligibilityExpiredAt:        expiration.BlogEligibilityExpiredAt.AddDate(0, 0, 30*planDuration),
		SocialMediaEligibilityExpiredAt: expiration.SocialMediaEligibilityExpiredAt.AddDate(0, 0, 30*planDuration),
		FrontPageEligibilityExpiredAt:   expiration.FrontPageEligibilityExpiredAt.AddDate(0, 0, 30*planDuration),
		CompanyPageEligibilityExpiredAt: expiration.CompanyPageEligibilityExpiredAt.AddDate(0, 0, 30*planDuration),
		PlanExpiredAt:                   expiration.PlanExpiredAt.AddDate(0, 0, 30*planDuration),
	}, nil
}

func (r *Repository) UpdateJobPlan(jobID int, expiration JobExpirationEntity) error {
	_, err := r.db.Exec(
		`UPDATE job SET newsletter_eligibility_expired_at = $1, blog_eligibility_expired_at = $2, social_media_eligibility_expired_at = $3, front_page_eligibility_expired_at = $4, company_page_eligibility_expired_at = $5, plan_expired_at = $6, approved_at = NOW() WHERE id = $7`,
		expiration.NewsletterEligibilityExpiredAt,
		expiration.BlogEligibilityExpiredAt,
		expiration.SocialMediaEligibilityExpiredAt,
		expiration.FrontPageEligibilityExpiredAt,
		expiration.CompanyPageEligibilityExpiredAt,
		expiration.PlanExpiredAt,
		jobID)
	return err
}

func (r *Repository) GetClickoutCountForJob(jobID int) (int, error) {
	var count int
	row := r.db.QueryRow(`select count(*) as c from job_event where job_event.event_type = 'clickout' and job_event.job_id = $1`, jobID)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}

func (r *Repository) LastJobPosted() (time.Time, error) {
	row := r.db.QueryRow(`SELECT created_at FROM job WHERE approved_at IS NOT NULL ORDER BY created_at DESC LIMIT 1`)
	var last time.Time
	if err := row.Scan(&last); err != nil {
		return last, err
	}

	return last, nil
}

func (r *Repository) SaveTokenForJob(token string, jobID int) error {
	_, err := r.db.Exec(`INSERT INTO edit_token (token, job_id, created_at) VALUES ($1, $2, $3)`, token, jobID, time.Now().UTC())
	if err != nil {
		return err
	}
	return err
}

func (r *Repository) GetValue(key string) (string, error) {
	res := r.db.QueryRow(`SELECT value FROM meta WHERE key = $1`, key)
	var val string
	err := res.Scan(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *Repository) SetValue(key, val string) error {
	_, err := r.db.Exec(`UPDATE meta SET value = $1 WHERE key = $2`, val, key)
	return err
}

func (r *Repository) ApplyToJob(jobID int, cv []byte, email, token string) error {
	stmt := `INSERT INTO apply_token (token, job_id, created_at, email, cv) VALUES ($1, $2, NOW(), $3, $4)`
	_, err := r.db.Exec(stmt, token, jobID, email, cv)
	return err
}

func (r *Repository) ConfirmApplyToJob(token string) error {
	_, err := r.db.Exec(
		`UPDATE apply_token SET confirmed_at = NOW() WHERE token = $1`,
		token,
	)
	return err
}

func (r *Repository) CleanupExpiredApplyTokens() error {
	_, err := r.db.Exec(
		`DELETE FROM apply_token WHERE created_at < NOW() - INTERVAL '3 days' OR confirmed_at IS NOT NULL`,
	)
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

func getQueryForArgs(conn *sql.DB, location, tag string, salary int, currency string, offset, max int, includePinnedJobs bool) (*sql.Rows, error) {
	planTypeFilter := "AND front_page_eligibility_expired_at < NOW()"
	if includePinnedJobs {
		planTypeFilter = "AND 1=1"
	}
	if tag == "" && location == "" && salary == 0 {
		return conn.Query(`
		SELECT count(*) OVER() AS full_count, id, job_title, company, company_url, salary_range, location, description, perks, interview_process, how_to_apply, created_at, url_id, slug, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period, expired, last_week_clickouts, plan_type, plan_duration, blog_eligibility_expired_at, company_page_eligibility_expired_at, front_page_eligibility_expired_at, newsletter_eligibility_expired_at, plan_expired_at, social_media_eligibility_expired_at
		FROM job
		WHERE approved_at IS NOT NULL `+planTypeFilter+` ORDER BY created_at DESC LIMIT $2 OFFSET $1`, offset, max)
	}
	if tag == "" && location != "" && salary == 0 {
		return conn.Query(`
		SELECT count(*) OVER() AS full_count, id, job_title, company, company_url, salary_range, location, description, perks, interview_process, how_to_apply, created_at, url_id, slug, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period, expired, last_week_clickouts, plan_type, plan_duration, blog_eligibility_expired_at, company_page_eligibility_expired_at, front_page_eligibility_expired_at, newsletter_eligibility_expired_at, plan_expired_at, social_media_eligibility_expired_at 
		FROM job
		WHERE approved_at IS NOT NULL `+planTypeFilter+` AND location ILIKE '%' || $1 || '%'
		ORDER BY created_at DESC LIMIT $3 OFFSET $2`, location, offset, max)
	}
	if tag != "" && location == "" && salary == 0 {
		return conn.Query(`
	SELECT count(*) OVER() AS full_count, id, job_title, company, company_url, salary_range, location, description, perks, interview_process, how_to_apply, created_at, url_id, slug, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period, expired, last_week_clickouts, plan_type, plan_duration, blog_eligibility_expired_at, company_page_eligibility_expired_at, front_page_eligibility_expired_at, newsletter_eligibility_expired_at, plan_expired_at, social_media_eligibility_expired_at
	FROM
	(
		SELECT id, job_title, company, company_url, salary_range, location, description, perks, interview_process, how_to_apply, created_at, url_id, slug, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period, expired, last_week_clickouts, plan_type, plan_duration, blog_eligibility_expired_at, company_page_eligibility_expired_at, front_page_eligibility_expired_at, newsletter_eligibility_expired_at, plan_expired_at, social_media_eligibility_expired_at, to_tsvector(job_title) || to_tsvector(company) || to_tsvector(description) AS doc
		FROM job WHERE approved_at IS NOT NULL `+planTypeFilter+`
	) AS job_
	WHERE job_.doc @@ to_tsquery($1)
	ORDER BY ts_rank(job_.doc, to_tsquery($1)) DESC, created_at DESC LIMIT $3 OFFSET $2`, tag, offset, max)
	}
	if tag != "" && location != "" && salary == 0 {
		return conn.Query(`
	SELECT count(*) OVER() AS full_count, id, job_title, company, company_url, salary_range, location, description, perks, interview_process, how_to_apply, created_at, url_id, slug, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period, expired, last_week_clickouts, plan_type, plan_duration, blog_eligibility_expired_at, company_page_eligibility_expired_at, front_page_eligibility_expired_at, newsletter_eligibility_expired_at, plan_expired_at, social_media_eligibility_expired_at
	FROM
	(
		SELECT id, job_title, company, company_url, salary_range, location, description, perks, interview_process, how_to_apply, created_at, url_id, slug, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period, expired, last_week_clickouts, plan_type, plan_duration, blog_eligibility_expired_at, company_page_eligibility_expired_at, front_page_eligibility_expired_at, newsletter_eligibility_expired_at, plan_expired_at, social_media_eligibility_expired_at, to_tsvector(job_title) || to_tsvector(company) || to_tsvector(description) AS doc
		FROM job WHERE approved_at IS NOT NULL `+planTypeFilter+`
	) AS job_
	WHERE job_.doc @@ to_tsquery($1)
	AND location ILIKE '%' || $2 || '%'
	ORDER BY ts_rank(job_.doc, to_tsquery($1)) DESC, created_at DESC LIMIT $4 OFFSET $3`, tag, location, offset, max)
	}
	if tag == "" && location == "" && salary != 0 {
		return conn.Query(`
		SELECT count(*) OVER() AS full_count, id, job_title, company, company_url, salary_range, location, description, perks, interview_process, how_to_apply, created_at, url_id, slug, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period, expired, last_week_clickouts, plan_type, plan_duration, blog_eligibility_expired_at, company_page_eligibility_expired_at, front_page_eligibility_expired_at, newsletter_eligibility_expired_at, plan_expired_at, social_media_eligibility_expired_at
		FROM job FULL JOIN fx_rate ON fx_rate.base = job.salary_currency_iso AND fx_rate.target = $4
		WHERE approved_at IS NOT NULL `+planTypeFilter+` AND (COALESCE(fx_rate.value, 1)*job.salary_max) >= $3 ORDER BY created_at DESC LIMIT $2 OFFSET $1`, offset, max, salary, currency)
	}
	if tag == "" && location != "" && salary != 0 {
		return conn.Query(`
		SELECT count(*) OVER() AS full_count, id, job_title, company, company_url, salary_range, location, description, perks, interview_process, how_to_apply, created_at, url_id, slug, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period, expired, last_week_clickouts, plan_type, plan_duration, blog_eligibility_expired_at, company_page_eligibility_expired_at, front_page_eligibility_expired_at, newsletter_eligibility_expired_at, plan_expired_at, social_media_eligibility_expired_at 
		FROM job FULL JOIN fx_rate ON fx_rate.base = job.salary_currency_iso AND fx_rate.target = $5
		WHERE approved_at IS NOT NULL `+planTypeFilter+` AND location ILIKE '%' || $1 || '%' AND (COALESCE(fx_rate.value, 1)*job.salary_max) >= $4
		ORDER BY created_at DESC LIMIT $3 OFFSET $2`, location, offset, max, salary, currency)
	}
	if tag != "" && location == "" && salary != 0 {
		return conn.Query(`
	SELECT count(*) OVER() AS full_count, id, job_title, company, company_url, salary_range, location, description, perks, interview_process, how_to_apply, created_at, url_id, slug, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period, expired, last_week_clickouts, plan_type, plan_duration, blog_eligibility_expired_at, company_page_eligibility_expired_at, front_page_eligibility_expired_at, newsletter_eligibility_expired_at, plan_expired_at, social_media_eligibility_expired_at
	FROM
	(
		SELECT id, job_title, company, company_url, salary_range, location, description, perks, interview_process, how_to_apply, created_at, url_id, slug, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period, expired, last_week_clickouts, plan_type, plan_duration, blog_eligibility_expired_at, company_page_eligibility_expired_at, front_page_eligibility_expired_at, newsletter_eligibility_expired_at, plan_expired_at, social_media_eligibility_expired_at, to_tsvector(job_title) || to_tsvector(company) || to_tsvector(description) AS doc
		FROM job FULL JOIN fx_rate ON fx_rate.base = job.salary_currency_iso AND fx_rate.target = $5 WHERE approved_at IS NOT NULL `+planTypeFilter+` AND (COALESCE(fx_rate.value, 1)*job.salary_max) >= $4
	) AS job_
	WHERE job_.doc @@ to_tsquery($1)
	ORDER BY ts_rank(job_.doc, to_tsquery($1)) DESC, created_at DESC LIMIT $3 OFFSET $2`, tag, offset, max, salary, currency)
	}

	return conn.Query(`
	SELECT count(*) OVER() AS full_count, id, job_title, company, company_url, salary_range, location, description, perks, interview_process, how_to_apply, created_at, url_id, slug, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period, expired, last_week_clickouts, plan_type, plan_duration, blog_eligibility_expired_at, company_page_eligibility_expired_at, front_page_eligibility_expired_at, newsletter_eligibility_expired_at, plan_expired_at, social_media_eligibility_expired_at
	FROM
	(
		SELECT id, job_title, company, company_url, salary_range, location, description, perks, interview_process, how_to_apply, created_at, url_id, slug, salary_min, salary_max, salary_currency, company_icon_image_id, external_id, salary_period, expired, last_week_clickouts, plan_type, plan_duration, blog_eligibility_expired_at, company_page_eligibility_expired_at, front_page_eligibility_expired_at, newsletter_eligibility_expired_at, plan_expired_at, social_media_eligibility_expired_at, to_tsvector(job_title) || to_tsvector(company) || to_tsvector(description) AS doc
		FROM job FULL JOIN fx_rate ON fx_rate.base = job.salary_currency_iso AND fx_rate.target = $6 WHERE approved_at IS NOT NULL `+planTypeFilter+` AND (COALESCE(fx_rate.value, 1)*job.salary_max) >= $5
	) AS job_
	WHERE job_.doc @@ to_tsquery($1)
	AND location ILIKE '%' || $2 || '%'
	ORDER BY ts_rank(job_.doc, to_tsquery($1)) DESC, created_at DESC LIMIT $4 OFFSET $3`, tag, location, offset, max, salary, currency)
}
