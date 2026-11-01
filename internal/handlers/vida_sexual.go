package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/puppe1990/amarra-cais/pkg/cais/flash"
	"github.com/puppe1990/amarra-cais/pkg/cais/httpx"
	"github.com/puppe1990/amarra-cais/pkg/cais/session"

	"github.com/puppe1990/aura-ciclos/internal/models"
)

type sexAnalysis struct {
	Days            int
	ProtectedDays   int
	UnprotectedDays int
	OrgasmYes       int
	OrgasmNo        int
	Times           int
}

func summarizeSex(logs []models.DailyLog) sexAnalysis {
	var a sexAnalysis
	for _, l := range logs {
		if l.Sex == "" && l.SexTimes == 0 {
			continue
		}
		a.Days++
		a.Times += l.SexTimes
		switch l.Sex {
		case "protegida":
			a.ProtectedDays++
		case "desprotegida":
			a.UnprotectedDays++
		}
		switch l.Orgasm {
		case "sim":
			a.OrgasmYes++
		case "nao":
			a.OrgasmNo++
		}
	}
	return a
}

func parseSexKind(v string) string {
	switch strings.TrimSpace(v) {
	case "protegida", "desprotegida":
		return v
	default:
		return ""
	}
}

func parseOrgasm(v string) string {
	switch strings.TrimSpace(v) {
	case "sim", "nao":
		return v
	default:
		return ""
	}
}

func parseSexTimes(v string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(v))
	if n < 0 {
		return 0
	}
	if n > 20 {
		return 20
	}
	return n
}

func (h *TrackerHandler) VidaSexual(w http.ResponseWriter, r *http.Request) {
	data, ok := h.load(w, r)
	if !ok {
		return
	}
	day := h.today()
	if raw := strings.TrimSpace(r.URL.Query().Get("dia")); raw != "" {
		if parsed, err := time.ParseInLocation(time.DateOnly, raw, time.UTC); err == nil {
			day = parsed
		}
	}
	log, err := h.store.GetDailyLog(data.user.ID, day)
	if err != nil {
		log = models.DailyLog{LoggedOn: day, UserID: data.user.ID}
	}
	h.page(w, r, "vida_sexual", map[string]any{
		"Title":     "Vida sexual",
		"ActiveNav": "registro",
		"Name":      greetName(data.user, data.profile),
		"Anonymous": data.user.Anonymous || data.profile.Anonymous,
		"Day":       day.Format(time.DateOnly),
		"Log":       log,
		"Analysis":  summarizeSex(data.logs),
		"Analise":   r.URL.Query().Get("analise") == "1",
	})
}

func (h *TrackerHandler) VidaSexualPost(w http.ResponseWriter, r *http.Request) {
	id, ok := session.UserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if err := httpx.ParseFormOrJSON(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	day, err := time.ParseInLocation(time.DateOnly, r.FormValue("logged_on"), time.UTC)
	if err != nil {
		day = h.today()
	}
	log, err := h.store.GetDailyLog(id, day)
	if err != nil {
		log = models.DailyLog{UserID: id, LoggedOn: day}
	}
	log.UserID = id
	log.LoggedOn = day
	log.Sex = parseSexKind(r.FormValue("sex"))
	log.Orgasm = parseOrgasm(r.FormValue("orgasm"))
	log.SexTimes = parseSexTimes(r.FormValue("sex_times"))
	if err := h.store.UpsertDailyLog(log); err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	flash.Set(w, "notice", "Vida sexual salva.", h.cfg.CookieSecure())
	http.Redirect(w, r, "/registro/vida-sexual?dia="+day.Format(time.DateOnly), http.StatusSeeOther)
}
