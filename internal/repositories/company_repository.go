package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/underdogio/job-board/internal/database"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const (
	companyEventPageView = "company_page_view"
)

// smart group by to map lower/upper case to same map entry with many entries and pickup the upper case one
// smart group by to find typos
func InferCompaniesFromJobs(ctx context.Context, since time.Time) ([]*database.Company, error) {
	var companies []*database.Company
	err := queries.Raw(
		`SELECT   trim(from company), 
         max(location)                  AS locations, 
         max(company_icon_image_id)     AS company_icon_id, 
         max(created_at)                AS last_job_created_at, 
         count(id)                      AS job_count, 
         count(approved_at IS NOT NULL) AS live_jobs_count,
		 max(company_page_eligibility_expired_at) AS company_page_eligibility_expired_at
FROM     job 
WHERE    company_icon_image_id IS NOT NULL 
AND      created_at > $1
AND      approved_at IS NOT NULL
GROUP BY trim(FROM company) 
ORDER BY trim(FROM company)`,
		since,
	).BindG(ctx, &companies)

	if err == sql.ErrNoRows || err == nil {
		return companies, nil
	}
	return companies, err
}

func SaveCompany(ctx context.Context, c *database.Company) error {
	if c.CompanyPageEligibilityExpiredAt.Time.Before(time.Now()) {
		c.CompanyPageEligibilityExpiredAt = null.TimeFrom(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	}
	return c.UpsertG(
		ctx,
		true,
		[]string{database.CompanyColumns.Name},
		boil.Whitelist(
			database.CompanyColumns.Locations,
			database.CompanyColumns.IconImageID,
			database.CompanyColumns.LastJobCreatedAt,
			database.CompanyColumns.TotalJobCount,
			database.CompanyColumns.ActiveJobCount,
			database.CompanyColumns.Slug,
			database.CompanyColumns.CompanyPageEligibilityExpiredAt,
		),
		boil.Infer(),
	)
}

func TrackCompanyView(ctx context.Context, company *database.Company) error {
	event := database.CompanyEvent{
		EventType: companyEventPageView,
		CompanyID: company.ID,
	}
	return event.InsertG(ctx, boil.Infer())
}

func CompanyBySlug(ctx context.Context, slug string) (*database.Company, error) {
	return database.Companies(
		database.CompanyWhere.Slug.EQ(slug),
	).OneG(ctx)
}

func CompaniesByQuery(ctx context.Context, location string, pageID, companiesPerPage int) ([]*database.Company, int, error) {
	offset := pageID*companiesPerPage - companiesPerPage
	mods := []qm.QueryMod{
		qm.Offset(offset),
		qm.Limit(companiesPerPage),
	}
	if location != "" {
		mods = append(mods, qm.Where(database.CompanyColumns.Locations+" ILIKE ?", "%"+location+"%"))
	}
	count, err := database.Companies(mods...).CountG(ctx)
	if err != nil {
		return []*database.Company{}, 0, err
	}
	companies, err := database.Companies(mods...).AllG(ctx)
	if err != nil {
		return companies, 0, err
	}
	return companies, int(count), nil
}

func FeaturedCompaniesPostAJob(ctx context.Context) ([]*database.Company, error) {
	return database.Companies(
		database.CompanyWhere.FeaturedPostAJob.EQ(null.BoolFrom(true)),
		qm.Limit(15),
	).AllG(ctx)
}

func GetCompanySlugs(ctx context.Context) ([]string, error) {
	slugs := make([]string, 0)
	err := database.Companies(
		qm.Select(database.CompanyColumns.Slug),
		database.CompanyWhere.Description.IsNotNull(),
	).BindG(ctx, &slugs)
	return slugs, err
}

func CompanyExists(ctx context.Context, company string) (bool, error) {
	return database.Jobs(
		qm.Where(database.JobColumns.Company+" ILIKE ?", "%"+company+"%"),
	).ExistsG(ctx)
}

func DeleteStaleImages(ctx context.Context, logoID string) error {
	_, err := database.Images(
		qm.Where(
			database.ImageColumns.ID+" NOT IN (SELECT company_icon_image_id FROM job WHERE company_icon_image_id IS NOT NULL)",
		),
		qm.Where(
			database.ImageColumns.ID+" NOT IN (SELECT icon_image_id FROM company)",
		),
		qm.Where(
			database.ImageColumns.ID+" NOT IN (SELECT image_id FROM developer_profile)",
		),
		database.ImageWhere.ID.NEQ(logoID),
	).DeleteAllG(ctx)
	return err
}
