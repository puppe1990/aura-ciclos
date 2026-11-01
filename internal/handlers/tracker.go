package handlers

import (
	"net/http"
	"time"

	"github.com/puppe1990/amarra-cais/pkg/amarra/view"
	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/httpx"
	"github.com/puppe1990/amarra-cais/pkg/cais/i18n"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"
	"github.com/puppe1990/amarra-cais/pkg/cais/session"

	"github.com/puppe1990/aura-ciclos/internal/cycle"
	"github.com/puppe1990/aura-ciclos/internal/models"
	"github.com/puppe1990/aura-ciclos/internal/store"
)

type TrackerHandler struct {
	views *view.Renderer
	store store.Store
	site  meta.Site
	cfg   cais.Config
	now   func() time.Time
}

func NewTrackerHandler(views *view.Renderer, s store.Store, site meta.Site, _ *i18n.Catalog, cfg cais.Config) *TrackerHandler {
	return &TrackerHandler{views: views, store: s, site: site, cfg: cfg}
}

func (h *TrackerHandler) today() time.Time {
	if h.now != nil {
		return h.now()
	}
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return time.Now().UTC()
	}
	return time.Now().In(loc)
}

type loaded struct {
	user    models.User
	profile models.Profile
	periods []models.Period
	logs    []models.DailyLog
	fc      cycle.Result
}

func (h *TrackerHandler) load(w http.ResponseWriter, r *http.Request) (loaded, bool) {
	id, ok := session.UserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return loaded{}, false
	}
	user, err := h.store.FindUser(id)
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return loaded{}, false
	}
	profile, err := h.store.GetProfile(id)
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return loaded{}, false
	}
	periods, err := h.store.ListPeriods(id)
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return loaded{}, false
	}
	since := h.today().AddDate(0, 0, -120)
	logs, err := h.store.ListDailyLogs(id, since)
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return loaded{}, false
	}
	fc := cycle.Forecast(cycle.Input{
		Today:             h.today(),
		Periods:           toCyclePeriods(periods),
		Observations:      toObservations(logs),
		Mode:              cycle.Mode(profile.Mode),
		PregnancyLMP:      profile.PregnancyLMP,
		TypicalCycleDays:  profile.TypicalCycleDays,
		TypicalPeriodDays: profile.TypicalPeriodDays,
	})
	return loaded{user: user, profile: profile, periods: periods, logs: logs, fc: fc}, true
}

func (h *TrackerHandler) page(w http.ResponseWriter, r *http.Request, name string, extra map[string]any) {
	writeView(w, r, h.views, h.cfg, "app", name, amarraData(r, h.site, extra), 0)
}

func greetName(user models.User, profile models.Profile) string {
	if profile.Anonymous || user.Anonymous {
		return "você"
	}
	if profile.DisplayName != "" {
		return profile.DisplayName
	}
	return "você"
}

func regularityLabel(fc cycle.Result) string {
	if len(fc.Lengths) == 0 {
		return "Ainda sem histórico"
	}
	if fc.Irregular {
		return "Ciclo irregular"
	}
	return "Ciclo regular"
}

func formatDayBR(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	return t.UTC().Format("02/01")
}

func formatOptionalDay(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.DateOnly)
}

func toCyclePeriods(in []models.Period) []cycle.Period {
	out := make([]cycle.Period, len(in))
	for i, p := range in {
		out[i] = cycle.Period{StartedOn: p.StartedOn, EndedOn: p.EndedOn, Flow: p.Flow}
	}
	return out
}

func toObservations(logs []models.DailyLog) []cycle.Observation {
	out := make([]cycle.Observation, 0, len(logs))
	for _, l := range logs {
		out = append(out, cycle.Observation{
			Day:           l.LoggedOn,
			BBTCelsius:    l.BBTCelsius,
			OvulationTest: l.OvulationTest,
		})
	}
	return out
}
