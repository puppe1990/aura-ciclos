package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/puppe1990/amarra-cais/pkg/cais/flash"
	"github.com/puppe1990/amarra-cais/pkg/cais/httpx"
	"github.com/puppe1990/amarra-cais/pkg/cais/session"

	"github.com/puppe1990/aura-ciclos/internal/cycle"
	"github.com/puppe1990/aura-ciclos/internal/models"
	"github.com/puppe1990/aura-ciclos/internal/symptoms"
)

func (h *TrackerHandler) Registro(w http.ResponseWriter, r *http.Request) {
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
	selected := map[string]bool{}
	for _, id := range log.SymptomIDs {
		selected[id] = true
	}
	h.page(w, r, "registro", map[string]any{
		"Title":      "Registro",
		"ActiveNav":  "registro",
		"Name":       greetName(data.user, data.profile),
		"Anonymous":  data.user.Anonymous || data.profile.Anonymous,
		"Mode":       data.profile.Mode,
		"ModeLabel":  cycle.ModeLabel(cycle.Mode(data.profile.Mode)),
		"Forecast":   data.fc,
		"PhaseLabel": cycle.PhaseLabel(data.fc.Phase),
		"Day":        day.Format(time.DateOnly),
		"DayLabel":   formatDayPT(day),
		"Log":        log,
		"Groups":     symptoms.Grouped(),
		"Selected":   selected,
		"Trying":     data.profile.Mode == models.ModeTrying,
		"BBTDisplay": formatBBT(log.BBTCelsius),
	})
}

func formatBBT(v *float64) string {
	if v == nil {
		return ""
	}
	return strconv.FormatFloat(*v, 'f', 2, 64)
}

func (h *TrackerHandler) RegistroPost(w http.ResponseWriter, r *http.Request) {
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
	energy, _ := strconv.Atoi(r.FormValue("energy"))
	water, _ := strconv.Atoi(r.FormValue("water_glasses"))
	sleep, _ := strconv.ParseFloat(strings.ReplaceAll(r.FormValue("sleep_hours"), ",", "."), 64)
	exercise, _ := strconv.Atoi(r.FormValue("exercise_minutes"))
	var bbt *float64
	if raw := strings.TrimSpace(strings.ReplaceAll(r.FormValue("bbt_celsius"), ",", ".")); raw != "" {
		if v, err := strconv.ParseFloat(raw, 64); err == nil {
			bbt = &v
		}
	}
	ids := r.Form["symptoms"]
	valid := make([]string, 0, len(ids))
	for _, sid := range ids {
		if _, ok := symptoms.ByID(sid); ok {
			valid = append(valid, sid)
		}
	}
	log := models.DailyLog{
		UserID:          id,
		LoggedOn:        day,
		SymptomIDs:      valid,
		Mood:            r.FormValue("mood"),
		Energy:          energy,
		SleepHours:      sleep,
		WaterGlasses:    water,
		ExerciseMinutes: exercise,
		Sex:             r.FormValue("sex"),
		BBTCelsius:      bbt,
		OvulationTest:   r.FormValue("ovulation_test"),
		PregnancyTest:   r.FormValue("pregnancy_test"),
		Notes:           strings.TrimSpace(r.FormValue("notes")),
	}
	if err := h.store.UpsertDailyLog(log); err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	flash.Set(w, "notice", "Registro salvo. As previsões usam esses dados no próximo cálculo.", h.cfg.CookieSecure())
	http.Redirect(w, r, "/registro?dia="+day.Format(time.DateOnly), http.StatusSeeOther)
}

func formatDayPT(day time.Time) string {
	meses := []string{"janeiro", "fevereiro", "março", "abril", "maio", "junho", "julho", "agosto", "setembro", "outubro", "novembro", "dezembro"}
	return "Hoje, " + strconv.Itoa(day.Day()) + " de " + meses[day.Month()-1]
}
