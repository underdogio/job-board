package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/underdogio/job-board/internal/server"
)

func IndexPageHandler(svr server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		location := r.URL.Query().Get("l")
		tag := r.URL.Query().Get("t")
		page := r.URL.Query().Get("p")

		var dst string
		if location != "" && tag != "" {
			dst = fmt.Sprintf("/%s-%s-Jobs-In-%s", strings.Title(svr.GetConfig().SiteJobCategory), tag, location)
		} else if location != "" {
			dst = fmt.Sprintf("/%s-Jobs-In-%s", strings.Title(svr.GetConfig().SiteJobCategory), location)
		} else if tag != "" {
			dst = fmt.Sprintf("/%s-%s-Jobs", strings.Title(svr.GetConfig().SiteJobCategory), tag)
		}
		if dst != "" && page != "" {
			dst += fmt.Sprintf("?p=%s", page)
		}
		if dst != "" {
			svr.Redirect(w, r, http.StatusMovedPermanently, dst)
			return
		}
		vars := mux.Vars(r)
		salary := vars["salary"]
		currency := vars["currency"]
		location = vars["location"]
		tag = vars["tag"]
		var validSalary bool
		for _, band := range svr.GetConfig().AvailableSalaryBands {
			if fmt.Sprintf("%d", band) == salary {
				validSalary = true
				break
			}
		}
		dst = "/"
		if location != "" && tag != "" {
			dst = fmt.Sprintf("/%s-%s-Jobs-In-%s", strings.Title(svr.GetConfig().SiteJobCategory), tag, location)
		} else if location != "" {
			dst = fmt.Sprintf("/%s-Jobs-In-%s", strings.Title(svr.GetConfig().SiteJobCategory), location)
		} else if tag != "" {
			dst = fmt.Sprintf("/%s-%s-Jobs", strings.Title(svr.GetConfig().SiteJobCategory), tag)
		}
		if page != "" {
			dst += fmt.Sprintf("?p=%s", page)
		}
		if (salary != "" && !validSalary) || (currency != "" && currency != "USD") {
			svr.Redirect(w, r, http.StatusMovedPermanently, dst)
			return
		}

		svr.RenderPageForLocationAndTag(w, r, "", "", page, salary, currency, "landing.html")
	}
}
