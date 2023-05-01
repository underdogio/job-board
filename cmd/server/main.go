package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/underdogio/job-board/internal/config"
	"github.com/underdogio/job-board/internal/email"
	"github.com/underdogio/job-board/internal/handler"
	"github.com/underdogio/job-board/internal/payment"
	"github.com/underdogio/job-board/internal/server"
	"github.com/underdogio/job-board/internal/template"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	boil.DebugMode = true
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %+v", err)
	}
	conn, err := server.GetDbConn(cfg.DatabaseURL)
	defer server.CloseDbConn(conn)
	if err != nil {
		log.Fatalf("unable to connect to postgres: %v", err)
	}
	emailClient, err := email.NewClient(
		cfg.Email2APIKey,
		cfg.SupportEmail,
		cfg.NoReplyEmail,
		cfg.SiteName,
		cfg.IsLocal,
	)
	if err != nil {
		log.Fatalf("unable to connect to sparkpost API: %v", err)
	}
	sessionStore := sessions.NewCookieStore(cfg.SessionKey)
	robotsTxtContent, err := os.ReadFile("./static/robots.txt")
	if err != nil {
		log.Fatalf("unable to read robots.txt placeholder file: %v", err)
	}
	securityTxtContent, err := os.ReadFile("./static/security.txt")
	if err != nil {
		log.Fatalf("unable to read security.txt placeholder file: %v", err)
	}

	paymentRepo := payment.NewRepository(cfg.StripeKey, cfg.SiteName, cfg.SiteHost)

	svr := server.NewServer(
		cfg,
		conn,
		mux.NewRouter(),
		template.NewTemplate(),
		emailClient,
		sessionStore,
	)

	svr.RegisterRoute("/sitemap.xml", handler.SitemapIndexHandler(svr), []string{"GET"})
	svr.RegisterRoute("/sitemap-{number}.xml", handler.SitemapHandler(svr), []string{"GET"})
	svr.RegisterRoute("/robots.txt", handler.RobotsTXTHandler(svr, robotsTxtContent), []string{"GET"})
	svr.RegisterRoute("/.well-known/security.txt", handler.WellKnownSecurityHandler(svr, securityTxtContent), []string{"GET"})

	svr.RegisterPathPrefix("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("./static/assets"))), []string{"GET"})

	svr.RegisterRoute("/about", handler.AboutPageHandler(svr), []string{"GET"})
	svr.RegisterRoute("/privacy-policy", handler.PrivacyPolicyPageHandler(svr), []string{"GET"})
	svr.RegisterRoute("/terms-of-service", handler.TermsOfServicePageHandler(svr), []string{"GET"})

	svr.RegisterRoute("/", handler.IndexPageHandler(svr), []string{"GET"})
	svr.RegisterRoute(
		fmt.Sprintf("/Companies-Using-%s", strings.Title(cfg.SiteJobCategory)),
		handler.CompaniesHandler(svr),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/Remote-Companies-Using-%s", strings.Title(cfg.SiteJobCategory)),
		handler.CompaniesForLocationHandler(svr, "Remote"),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/Companies-Using-%s-In-{location}", strings.Title(cfg.SiteJobCategory)),
		handler.CompaniesHandler(svr),
		[]string{"GET"},
	)

	// developers pages
	// svr.RegisterRoute(
	// 	fmt.Sprintf("/%s-Developers", strings.Title(cfg.SiteJobCategory)),
	// 	handler.DevelopersHandler(svr, devRepo),
	// 	[]string{"GET"},
	// )
	// svr.RegisterRoute(
	// 	fmt.Sprintf("/%s-Developers-In-{location}", strings.Title(cfg.SiteJobCategory)),
	// 	handler.DevelopersHandler(svr, devRepo),
	// 	[]string{"GET"},
	// )
	svr.RegisterRoute(
		fmt.Sprintf("/%s-{tag}-Developers", strings.Title(cfg.SiteJobCategory)),
		handler.DevelopersHandler(svr),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/%s-{tag}-Developers-In-{location}", strings.Title(cfg.SiteJobCategory)),
		handler.DevelopersHandler(svr),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		"/Submit-Developer-Profile",
		handler.PermanentRedirectHandler(svr, fmt.Sprintf("/Join-%s-Community", strings.Title(cfg.SiteJobCategory))),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		"/community",
		handler.PermanentRedirectHandler(svr, fmt.Sprintf("/Join-%s-Community", strings.Title(cfg.SiteJobCategory))),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/Join-%s-Community", strings.Title(cfg.SiteJobCategory)),
		handler.SubmitDeveloperProfileHandler(svr),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/Hire-From-%s-Community", strings.Title(cfg.SiteJobCategory)),
		handler.SubmitRecruiterProfileHandler(svr),
		[]string{"GET"},
	)
	svr.RegisterRoute("/x/srp", handler.SaveRecruiterProfileHandler(svr), []string{"POST"})
	svr.RegisterRoute("/x/sdp", handler.SaveDeveloperProfileHandler(svr), []string{"POST"})
	svr.RegisterRoute("/x/udp", handler.UpdateDeveloperProfileHandler(svr), []string{"POST"})
	svr.RegisterRoute("/x/ddp", handler.DeleteDeveloperProfileHandler(svr), []string{"POST"})
	svr.RegisterRoute("/x/smdp/{id}", handler.SendMessageDeveloperProfileHandler(svr), []string{"POST"})
	svr.RegisterRoute("/developer/{slug}", handler.ViewDeveloperProfileHandler(svr), []string{"GET"})
	svr.RegisterRoute("/x/auth/message/{id}", handler.DeliverMessageDeveloperProfileHandler(svr), []string{"GET"})

	// blog
	svr.RegisterRoute("/profile/home", handler.ProfileHomepageHandler(svr), []string{"GET"})
	svr.RegisterRoute("/profile/{id}/edit", handler.EditProfileHandler(svr), []string{"GET"})
	svr.RegisterRoute("/profile/blog/create", handler.CreateDraftBlogPostHandler(svr), []string{"GET"})
	svr.RegisterRoute("/profile/blog/list", handler.GetUserBlogPostsHandler(svr), []string{"GET"})
	svr.RegisterRoute("/profile/blog/{id}/edit", handler.EditBlogPostHandler(svr), []string{"GET"})
	svr.RegisterRoute("/x/profile/blog/create", handler.CreateBlogPostHandler(svr), []string{"POST"})
	svr.RegisterRoute("/x/profile/blog/{id}/publish", handler.PublishBlogPostHandler(svr), []string{"POST"})
	svr.RegisterRoute("/x/profile/blog/{id}/unpublish", handler.UnpublishBlogPostHandler(svr), []string{"POST"})
	svr.RegisterRoute("/x/profile/blog/{id}/update", handler.UpdateBlogPostHandler(svr), []string{"POST"})
	svr.RegisterRoute("/blog/{slug}", handler.GetBlogPostBySlugHandler(svr), []string{"GET"})
	svr.RegisterRoute("/blog", handler.GetAllPublishedBlogPostsHandler(svr), []string{"GET"})

	// tasks
	svr.RegisterRoute("/x/task/weekly-newsletter", handler.TriggerWeeklyNewsletter(svr), []string{"POST"})
	svr.RegisterRoute("/x/task/ads-manager", handler.TriggerAdsManager(svr), []string{"POST"})
	svr.RegisterRoute("/x/task/twitter-scheduler", handler.TriggerTwitterScheduler(svr), []string{"POST"})
	svr.RegisterRoute("/x/task/telegram-scheduler", handler.TriggerTelegramScheduler(svr), []string{"POST"})
	// svr.RegisterRoute("/x/task/company-update", handler.TriggerCompanyUpdate(svr), []string{"POST"})
	svr.RegisterRoute("/x/task/sitemap-update", handler.TriggerSitemapUpdate(svr), []string{"POST"})
	svr.RegisterRoute("/x/task/cloudflare-stats-export", handler.TriggerCloudflareStatsExport(svr), []string{"POST"})
	svr.RegisterRoute("/x/task/expired-jobs", handler.TriggerExpiredJobsTask(svr), []string{"POST"})
	svr.RegisterRoute("/x/task/update-last-week-clickouts", handler.TriggerUpdateLastWeekClickouts(svr), []string{"POST"})
	svr.RegisterRoute("/x/task/monthly-highlights", handler.TriggerMonthlyHighlights(svr), []string{"POST"})
	svr.RegisterRoute("/x/task/fx-rate-update", handler.TriggerFXRateUpdate(svr), []string{"POST"})

	// view newsletter
	svr.RegisterRoute("/newsletter", handler.ViewNewsletterPageHandler(svr), []string{"GET"})

	// view support
	svr.RegisterRoute("/support", handler.ViewSupportPageHandler(svr), []string{"GET"})

	// post a job succeeded
	svr.RegisterRoute("/x/j/p/1", handler.PostAJobSuccessPageHandler(svr), []string{"GET"})

	// post a job failed
	svr.RegisterRoute("/x/j/p/0", handler.PostAJobFailurePageHandler(svr), []string{"GET"})

	// newsletter member save
	svr.RegisterRoute("/x/email/subscribe", handler.AddEmailSubscriberHandler(svr), []string{"GET"})
	svr.RegisterRoute("/x/email/unsubscribe", handler.RemoveEmailSubscriberHandler(svr), []string{"GET"})
	svr.RegisterRoute("/x/email/confirm/{token}", handler.ConfirmEmailSubscriberHandler(svr), []string{"GET"})

	// apply for job
	svr.RegisterRoute("/x/a/e", handler.ApplyForJobPageHandler(svr), []string{"POST"})

	// apply to job confirmation
	svr.RegisterRoute("/apply/{token}", handler.ApplyToJobConfirmation(svr), []string{"GET"})

	// submit job post
	svr.RegisterRoute("/x/s", handler.SubmitJobPostPageHandler(svr, paymentRepo), []string{"POST"})

	// re-submit job post payment for upsell
	svr.RegisterRoute("/x/s/upsell", handler.SubmitJobPostPaymentUpsellPageHandler(svr, paymentRepo), []string{"POST"})

	// save media file
	svr.RegisterRoute("/x/s/m", handler.SaveMediaPageHandler(svr), []string{"POST"})

	// update media file
	svr.RegisterRoute("/x/s/m/{id}", handler.UpdateMediaPageHandler(svr), []string{"PUT"})

	// retrieve media file
	svr.RegisterRoute("/x/s/m/{id}", handler.RetrieveMediaPageHandler(svr), []string{"GET"})

	// retrieve meta image media file
	svr.RegisterRoute("/x/s/m/meta/{id}", handler.RetrieveMediaMetaPageHandler(svr), []string{"GET"})

	// stripe payment confirmation webhook
	svr.RegisterRoute("/x/stripe/checkout/completed", handler.StripePaymentConfirmationWebookHandler(svr), []string{"POST"})

	// send feedback message
	svr.RegisterRoute("/x/s/message", handler.SendFeedbackMessage(svr), []string{"POST"})

	// track job clickout
	svr.RegisterRoute("/x/j/c/{id}", handler.TrackJobClickoutPageHandler(svr), []string{"GET"})

	// track job clickout + redirect to job page
	svr.RegisterRoute("/x/r", handler.TrackJobClickoutAndRedirectToJobPage(svr), []string{"GET"})

	// autocomplete locations
	svr.RegisterRoute("/x/loc/autocomplete", handler.AutocompleteLocation(svr), []string{"GET"})

	// autocomplete skills
	svr.RegisterRoute("/x/skill/autocomplete", handler.AutocompleteSkill(svr), []string{"GET"})

	// view job by slug
	svr.RegisterRoute("/job/{slug}", handler.JobBySlugPageHandler(svr), []string{"GET"})

	// view company by slug
	svr.RegisterRoute("/company/{slug}", handler.CompanyBySlugPageHandler(svr), []string{"GET"})

	//
	// auth routes
	//

	// sign on page
	svr.RegisterRoute("/auth", handler.GetAuthPageHandler(svr), []string{"GET"})

	// sign on email link
	svr.RegisterRoute("/x/auth/link", handler.RequestTokenSignOn(svr), []string{"POST"})
	svr.RegisterRoute("/x/auth/{token}", handler.VerifyTokenSignOn(svr, cfg.AdminEmail), []string{"GET"})

	//
	// private routes
	// at the moment only protected by static token
	//

	// @private: update job by token
	svr.RegisterRoute("/x/u", handler.UpdateJobPageHandler(svr), []string{"POST"})

	// @private: view edit job by token
	svr.RegisterRoute("/edit/{token}", handler.EditJobViewPageHandler(svr), []string{"GET"})

	// @private: disapprove job by token
	svr.RegisterRoute("/x/d", handler.DisapproveJobPageHandler(svr), []string{"POST"})

	//
	// landing page routes
	//

	// Aliases
	svr.RegisterRoute(
		fmt.Sprintf("/%s-Jobs", cfg.SiteJobCategory),
		handler.PermanentRedirectHandler(svr, "/"),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		"/Remote-Jobs",
		handler.PermanentRedirectHandler(svr, fmt.Sprintf("/Remote-%s-Jobs", strings.Title(cfg.SiteJobCategory))),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		"/youtube",
		handler.PermanentExternalRedirectHandler(svr, fmt.Sprintf("https://www.youtube.com/c/%s", cfg.SiteYoutube)),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		"/telegram",
		handler.PermanentExternalRedirectHandler(svr, fmt.Sprintf("https://t.me/%s", cfg.SiteTelegramChannel)),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		"/twitter",
		handler.PermanentExternalRedirectHandler(svr, fmt.Sprintf("https://twitter.com/%s", cfg.SiteTwitter)), []string{"GET"})
	svr.RegisterRoute(
		"/linkedin",
		handler.PermanentExternalRedirectHandler(svr, fmt.Sprintf("https://www.linkedin.com/company/%s", cfg.SiteLinkedin)),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		"/github",
		handler.PermanentExternalRedirectHandler(svr, fmt.Sprintf("https://github.com/%s", cfg.SiteGithub)),
		[]string{"GET"},
	)

	// Remote Landing Page No Skill
	svr.RegisterRoute(
		fmt.Sprintf("/Remote-%s-Jobs", strings.Title(cfg.SiteJobCategory)),
		handler.LandingPageForLocationHandler(svr, "Remote"),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/Remote-%s-Jobs-Paying-{salary}-{currency}-year", strings.Title(cfg.SiteJobCategory)),
		handler.LandingPageForLocationHandler(svr, "Remote"),
		[]string{"GET"},
	)

	// Salary Only Landing Page
	svr.RegisterRoute(
		fmt.Sprintf("/%s-Jobs-Paying-{salary}-{currency}-year", strings.Title(cfg.SiteJobCategory)),
		handler.IndexPageHandler(svr),
		[]string{"GET"},
	)

	// Remote Landing Page Skill
	svr.RegisterRoute(
		fmt.Sprintf("/Remote-%s-{skill}-Jobs", strings.Title(cfg.SiteJobCategory)),
		handler.LandingPageForLocationAndSkillPlaceholderHandler(svr, "Remote"),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/Remote-%s-{skill}-Jobs-Paying-{salary}-{currency}-year", strings.Title(cfg.SiteJobCategory)),
		handler.LandingPageForLocationAndSkillPlaceholderHandler(svr, "Remote"),
		[]string{"GET"},
	)

	// Location Only Landing Page
	svr.RegisterRoute(
		fmt.Sprintf("/%s-Jobs-In-{location}-Paying-{salary}-{currency}-year", strings.Title(cfg.SiteJobCategory)),
		handler.LandingPageForLocationPlaceholderHandler(svr),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/%s-Jobs-in-{location}-Paying-{salary}-{currency}-year", strings.Title(cfg.SiteJobCategory)),
		handler.LandingPageForLocationPlaceholderHandler(svr),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/%s-Jobs-In-{location}", strings.Title(cfg.SiteJobCategory)),
		handler.LandingPageForLocationPlaceholderHandler(svr),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/%s-Jobs-in-{location}", strings.Title(cfg.SiteJobCategory)),
		handler.LandingPageForLocationPlaceholderHandler(svr),
		[]string{"GET"},
	)

	// Skill Only Landing Page
	svr.RegisterRoute(
		fmt.Sprintf("/%s-{skill}-Jobs", strings.Title(cfg.SiteJobCategory)),
		handler.LandingPageForSkillPlaceholderHandler(svr),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/%s-{skill}-Jobs-Paying-{salary}-{currency}-year", strings.Title(cfg.SiteJobCategory)),
		handler.LandingPageForSkillPlaceholderHandler(svr),
		[]string{"GET"},
	)

	// Skill And Location Landing Page
	svr.RegisterRoute(
		fmt.Sprintf("/%s-{skill}-Jobs-In-{location}-Paying-{salary}-{currency}-year", strings.Title(cfg.SiteJobCategory)),
		handler.LandingPageForSkillAndLocationPlaceholderHandler(svr),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/%s-{skill}-Jobs-in-{location}-Paying-{salary}-{currency}-year", strings.Title(cfg.SiteJobCategory)),
		handler.LandingPageForSkillAndLocationPlaceholderHandler(svr),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/%s-{skill}-Jobs-In-{location}", strings.Title(cfg.SiteJobCategory)),
		handler.LandingPageForSkillAndLocationPlaceholderHandler(svr),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/%s-{skill}-Jobs-in-{location}", strings.Title(cfg.SiteJobCategory)),
		handler.LandingPageForSkillAndLocationPlaceholderHandler(svr),
		[]string{"GET"},
	)

	// // Salary for location
	// svr.RegisterRoute(
	// 	fmt.Sprintf("/%s-Developer-Salary-{location}", strings.Title(cfg.SiteJobCategory)),
	// 	handler.SalaryLandingPageLocationPlaceholderHandler(svr, devRepo),
	// 	[]string{"GET"},
	// )
	// // Salary for remote
	// svr.RegisterRoute(
	// 	fmt.Sprintf("/Remote-%s-Developer-Salary", strings.Title(cfg.SiteJobCategory)),
	// 	handler.SalaryLandingPageLocationHandler(svr, devRepo, "Remote"),
	// 	[]string{"GET"},
	// )

	// hire developers pages
	svr.RegisterRoute(
		fmt.Sprintf("/Hire-%s-Developers", strings.Title(cfg.SiteJobCategory)),
		handler.PostAJobPageHandler(svr),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/hire-%s-developers", cfg.SiteJobCategory),
		handler.PermanentRedirectHandler(svr, fmt.Sprintf("Hire-%s-Developers", strings.Title(cfg.SiteJobCategory))),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		"/post-a-job",
		handler.PermanentRedirectHandler(svr, fmt.Sprintf("Hire-%s-Developers", strings.Title(cfg.SiteJobCategory))),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		"/ad",
		handler.PermanentRedirectHandler(svr, fmt.Sprintf("Hire-%s-Developers", strings.Title(cfg.SiteJobCategory))),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/Hire-Remote-%s-Developers", strings.Title(cfg.SiteJobCategory)),
		handler.PostAJobForLocationPageHandler(svr, "Remote"),
		[]string{"GET"},
	)
	svr.RegisterRoute(
		fmt.Sprintf("/Hire-%s-Developers-In-{location}", strings.Title(cfg.SiteJobCategory)),
		handler.PostAJobForLocationFromURLPageHandler(svr),
		[]string{"GET"},
	)

	// generic payment intent view
	svr.RegisterRoute("/x/payment", handler.ShowPaymentPage(svr), []string{"GET"})
	// generic payment intent processing
	// svr.RegisterRoute("/x/payment-intent", handler.GeneratePaymentIntent(svr, paymentRepo), []string{"POST"})

	// RSS feed
	svr.RegisterRoute("/rss", handler.ServeRSSFeed(svr), []string{"GET"})

	//
	// admin routes
	// protected by jwt auth
	//

	// @admin: submit job without payment view
	svr.RegisterRoute("/manage/new", handler.PostAJobWithoutPaymentPageHandler(svr), []string{"GET"})

	// @admin: list/search jobs as admin
	svr.RegisterRoute("/manage/list", handler.ListJobsAsAdminPageHandler(svr), []string{"GET"})

	// @admin: view job as admin (alias to manage/edit/{token})
	svr.RegisterRoute("/manage/job/{slug}", handler.ManageJobBySlugViewPageHandler(svr), []string{"GET"})

	// @admin: view manage job page
	svr.RegisterRoute("/manage/{token}", handler.ManageJobViewPageHandler(svr), []string{"GET"})

	// @admin: submit job without payment
	svr.RegisterRoute("/x/sp", handler.SubmitJobPostWithoutPaymentHandler(svr), []string{"POST"})

	// @admin: approve job
	svr.RegisterRoute("/x/a", handler.ApproveJobPageHandler(svr), []string{"POST"})

	// @admin: permanently delete job and all child resources (image, clickouts, edit token)
	svr.RegisterRoute("/x/j/d", handler.PermanentlyDeleteJobByToken(svr), []string{"POST"})

	log.Fatal(svr.Run())
}
