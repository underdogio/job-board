package job

import (
	"time"

	"github.com/lib/pq"
)

const (
	jobEventPageView = "page_view"
	jobEventClickout = "clickout"

	SearchTypeJob    = "job"
	SearchTypeSalary = "salary"
)

const (
	JobAdBasic = iota
	JobAdSponsoredBackground
	JobAdSponsoredPinnedFor30Days
	JobAdSponsoredPinnedFor7Days
	JobAdWithCompanyLogo
	JobAdSponsoredPinnedFor60Days
	JobAdSponsoredPinnedFor90Days
)

type Job struct {
	CreatedAt                       int64
	JobTitle                        string
	Company                         string
	SalaryMin                       string
	SalaryMax                       string
	SalaryCurrency                  string
	SalaryPeriod                    string
	SalaryRange                     string
	Location                        string
	Description                     string
	Perks                           string
	InterviewProcess                string
	HowToApply                      string
	Email                           string
	Expired                         bool
	LastWeekClickouts               int
	PlanType                        string
	PlanDuration                    int
	NewsletterEligibilityExpiredAt  time.Time
	BlogEligibilityExpiredAt        time.Time
	SocialMediaEligibilityExpiredAt time.Time
	FrontPageEligibilityExpiredAt   time.Time
	CompanyPageEligibilityExpiredAt time.Time
	PlanExpiredAt                   time.Time
}

type JobRq struct {
	JobTitle        string `json:"job_title"`
	JobCategory     string `json:"job_category"`
	Company         string `json:"company_name"`
	Location        string `json:"job_location"`
	SalaryRange     string `json:"salary_range"`
	JobType         string `json:"job_type"`
	ApplicationLink string `json:"application_link"`
	Description     string `json:"job_description"`
	CompanyIconID   string `json:"company_icon_id,omitempty"`
	Email           string `json:"email"`
	StripeToken     string `json:"stripe_token,omitempty"`
}

type JobRqUpsell struct {
	Token           string `json:"token"`
	Email           string `json:"email"`
	StripeToken     string `json:"stripe_token,omitempty"`
	PlanType        string `json:"plan_type"`
	PlanDuration    int
	PlanDurationStr string `json:"plan_duration"`
}

type JobRqUpdate struct {
	JobTitle         string `json:"job_title"`
	Location         string `json:"job_location"`
	Company          string `json:"company_name"`
	CompanyURL       string `json:"company_url"`
	SalaryMin        string `json:"salary_min"`
	SalaryMax        string `json:"salary_max"`
	SalaryCurrency   string `json:"salary_currency"`
	Description      string `json:"job_description"`
	HowToApply       string `json:"how_to_apply"`
	Perks            string `json:"perks"`
	InterviewProcess string `json:"interview_process"`
	Email            string `json:"company_email"`
	Token            string `json:"token"`
	CompanyIconID    string `json:"company_icon_id,omitempty"`
	SalaryPeriod     string `json:"salary_period"`
}

type JobPost struct {
	ID                              int
	CreatedAt                       int64
	TimeAgo                         string
	JobTitle                        string
	JobCategory                     string
	Company                         string
	Location                        string
	SalaryRange                     string
	JobType                         string
	ApplicationLink                 string
	JobDescription                  string
	Slug                            string
	CompanyIconID                   string
	ExternalID                      string
	ApprovedAt                      *time.Time
	Expired                         bool
	LastWeekClickouts               int
	NewsletterEligibilityExpiredAt  time.Time
	BlogEligibilityExpiredAt        time.Time
	SocialMediaEligibilityExpiredAt time.Time
	FrontPageEligibilityExpiredAt   time.Time
	CompanyPageEligibilityExpiredAt time.Time
	PlanExpiredAt                   time.Time
}

type JobPostForEdit struct {
	ID                                                                        int
	JobTitle, Company, CompanyEmail, CompanyURL, Location                     string
	SalaryMin, SalaryMax                                                      int
	SalaryCurrency, JobDescription, Perks, InterviewProcess, HowToApply, Slug string
	CreatedAt                                                                 time.Time
	ApprovedAt                                                                pq.NullTime
	CompanyIconID                                                             string
	ExternalID                                                                string
	SalaryPeriod                                                              string
	PlanType                                                                  string
	PlanDuration                                                              int
	NewsletterEligibilityExpiredAt                                            time.Time
	BlogEligibilityExpiredAt                                                  time.Time
	SocialMediaEligibilityExpiredAt                                           time.Time
	FrontPageEligibilityExpiredAt                                             time.Time
	CompanyPageEligibilityExpiredAt                                           time.Time
	PlanExpiredAt                                                             time.Time
}

type JobStat struct {
	Date      string `json:"date"`
	Clickouts int    `json:"clickouts"`
	PageViews int    `json:"pageviews"`
}

type JobApplyURL struct {
	ID  int
	URL string
}

type Applicant struct {
	Cv    []byte
	Email string
}

type JobCategory struct {
	Key   string
	Label string
}

var JobCategories = []JobCategory{
	{
		Key:   "software-engineer",
		Label: "Software Engineer",
	},
	{
		Key:   "engineering-manager",
		Label: "Engineering Manager",
	},
	{
		Key:   "data-engineer",
		Label: "Data Engineer",
	},
	{
		Key:   "devops-engineer",
		Label: "DevOps Engineer",
	},
	{
		Key:   "security-engineer",
		Label: "Security Engineer",
	},
	{
		Key:   "qa-engineer",
		Label: "QA Engineer",
	},
	{
		Key:   "data-scientist",
		Label: "Data Scientist",
	},
	{
		Key:   "data-analyst",
		Label: "Data Analyst",
	},
	{
		Key:   "mobile-developer",
		Label: "Mobile Developer",
	},
	{
		Key:   "ui-ux-designer",
		Label: "UI/UX Designer",
	},
	{
		Key:   "design-manager",
		Label: "Design Manager",
	},
	{
		Key:   "ux-researcher",
		Label: "UX Researcher",
	},
	{
		Key:   "product-manager",
		Label: "Product Manager",
	},
	{
		Key:   "project-manager",
		Label: "Project Manager",
	},
	{
		Key:   "product-designer",
		Label: "Product Designer",
	},
	{
		Key:   "technical-project-manager",
		Label: "Technical Project Manager",
	},
	{
		Key:   "staff-engineer",
		Label: "Staff Engineer",
	},
	{
		Key:   "sales",
		Label: "Sales",
	},
	{
		Key:   "business-development",
		Label: "Business Development",
	},
	{
		Key:   "account-managers",
		Label: "Account Managers",
	},
	{
		Key:   "customer-success",
		Label: "Customer Success",
	},
	{
		Key:   "marketing",
		Label: "Marketing",
	},
	{
		Key:   "business-operations",
		Label: "Business Operations",
	},
	{
		Key:   "business-product-manager",
		Label: "Business Product Manager",
	},
	{
		Key:   "business-project-manager",
		Label: "Business Project Manager",
	},
	{
		Key:   "recruiter",
		Label: "Recruiter",
	},
	{
		Key:   "human-resources",
		Label: "Human Resources",
	},
	{
		Key:   "finance-accounting",
		Label: "Finance/Accounting",
	},
	{
		Key:   "legal",
		Label: "Legal",
	},
	{
		Key:   "content-writing",
		Label: "Content/Writing",
	},
	{
		Key:   "community-manager",
		Label: "Community Manager",
	},
}
