package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/puppe1990/amarra-cais/pkg/cais/flash"
	"github.com/puppe1990/amarra-cais/pkg/cais/httpx"
	"github.com/puppe1990/amarra-cais/pkg/cais/session"

	"github.com/puppe1990/aura-ciclos/internal/cycle"
	"github.com/puppe1990/aura-ciclos/internal/models"
)

func (h *TrackerHandler) Ciclo(w http.ResponseWriter, r *http.Request) {
	data, ok := h.load(w, r)
	if !ok {
		return
	}
	rems, err := h.store.ListReminders(data.user.ID)
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	phys := cycle.Physiology(data.fc.Phase)
	fertileN, fertileTotal := cycle.FertileIndex(h.today(), data.fc.FertileStart, data.fc.FertileEnd)
	h.page(w, r, "ciclo", map[string]any{
		"Title":           "Ciclo",
		"ActiveNav":       "ciclo",
		"Name":            greetName(data.user, data.profile),
		"Anonymous":       data.user.Anonymous || data.profile.Anonymous,
		"Mode":            data.profile.Mode,
		"ModeLabel":       cycle.ModeLabel(cycle.Mode(data.profile.Mode)),
		"Forecast":        data.fc,
		"PhaseLabel":      cycle.PhaseLabel(data.fc.Phase),
		"PhaseBadge":      phys.Badge,
		"RegularLabel":    regularityLabel(data.fc),
		"Week":            weekStrip(h.today(), data.fc, true),
		"Ring":            ringView(data.fc),
		"Reminders":       rems,
		"PregnancyNote":   cycle.PregnancyNote(data.fc.PregnancyWeek),
		"PregnancySize":   cycle.PregnancySize(data.fc.PregnancyWeek),
		"Today":           h.today().Format(time.DateOnly),
		"NextPeriodDate":  formatDayBR(data.fc.NextPeriodOn),
		"FertileStart":    formatDayBR(data.fc.FertileStart),
		"FertileEnd":      formatDayBR(data.fc.FertileEnd),
		"OvulationDate":   formatDayBR(data.fc.OvulationOn),
		"DueDate":         formatDayBR(data.fc.DueOn),
		"Trying":          data.profile.Mode == models.ModeTrying,
		"Pregnant":        data.profile.Mode == models.ModePregnancy,
		"Physiology":      phys,
		"ChanceLabel":     cycle.ChanceWords(data.fc.ConceptionChance),
		"ConfidenceWords": cycle.ConfidenceWords(data.fc.Confidence),
		"IntervalRange":   cycle.IntervalRange(data.fc.Lengths),
		"FertileN":        fertileN,
		"FertileTotal":    fertileTotal,
		"Spark":           bbtSparkline(data.logs),
		"PeakFertile":     data.fc.ConceptionChance >= 75,
	})
}

func (h *TrackerHandler) PeriodoPost(w http.ResponseWriter, r *http.Request) {
	id, ok := session.UserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if err := httpx.ParseFormOrJSON(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	started, err := time.ParseInLocation(time.DateOnly, r.FormValue("started_on"), time.UTC)
	if err != nil {
		flash.Set(w, "alert", "Informe a data de início da menstruação.", h.cfg.CookieSecure())
		http.Redirect(w, r, "/ciclo", http.StatusSeeOther)
		return
	}
	flow := r.FormValue("flow")
	if flow == "" {
		flow = "medio"
	}
	p := models.Period{UserID: id, StartedOn: started, Flow: flow}
	if endedRaw := strings.TrimSpace(r.FormValue("ended_on")); endedRaw != "" {
		ended, err := time.ParseInLocation(time.DateOnly, endedRaw, time.UTC)
		if err == nil {
			p.EndedOn = &ended
		}
	}
	if err := h.store.UpsertPeriod(p); err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	flash.Set(w, "notice", "Menstruação registrada. As previsões foram recalculadas.", h.cfg.CookieSecure())
	http.Redirect(w, r, "/ciclo", http.StatusSeeOther)
}

func (h *TrackerHandler) ModoPost(w http.ResponseWriter, r *http.Request) {
	id, ok := session.UserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if err := httpx.ParseFormOrJSON(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p, err := h.store.GetProfile(id)
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	switch r.FormValue("mode") {
	case string(models.ModeTrying):
		p.Mode = models.ModeTrying
	case string(models.ModePregnancy):
		p.Mode = models.ModePregnancy
	default:
		p.Mode = models.ModeCycle
	}
	if raw := strings.TrimSpace(r.FormValue("pregnancy_lmp")); raw != "" {
		if lmp, err := time.ParseInLocation(time.DateOnly, raw, time.UTC); err == nil {
			p.PregnancyLMP = &lmp
		}
	}
	if n, err := strconv.Atoi(r.FormValue("typical_cycle_days")); err == nil && n >= 15 && n <= 60 {
		p.TypicalCycleDays = n
	}
	if n, err := strconv.Atoi(r.FormValue("typical_period_days")); err == nil && n >= 1 && n <= 15 {
		p.TypicalPeriodDays = n
	}
	if err := h.store.SaveProfile(p); err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	flash.Set(w, "notice", "Modo atualizado.", h.cfg.CookieSecure())
	http.Redirect(w, r, "/relatorios", http.StatusSeeOther)
}

func (h *TrackerHandler) LembretePost(w http.ResponseWriter, r *http.Request) {
	id, ok := session.UserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if err := httpx.ParseFormOrJSON(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	kind := r.FormValue("kind")
	if kind != "pill" && kind != "water" && kind != "period" {
		kind = "pill"
	}
	hour, _ := strconv.Atoi(r.FormValue("hour"))
	minute, _ := strconv.Atoi(r.FormValue("minute"))
	if hour < 0 || hour > 23 {
		hour = 8
	}
	if minute < 0 || minute > 59 {
		minute = 0
	}
	if err := h.store.UpsertReminder(models.Reminder{
		UserID:  id,
		Kind:    kind,
		Hour:    hour,
		Minute:  minute,
		Enabled: true,
	}); err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	flash.Set(w, "notice", "Lembrete salvo. Ele aparece na tela do ciclo neste aparelho.", h.cfg.CookieSecure())
	http.Redirect(w, r, "/apoio", http.StatusSeeOther)
}

type weekChip struct {
	Date       string
	Weekday    string
	Day        int
	Today      bool
	Phase      cycle.Phase
	PhaseClass string
}

func weekStrip(anchor time.Time, fc cycle.Result, labelToday bool) []weekChip {
	weekdays := []string{"DOM", "SEG", "TER", "QUA", "QUI", "SEX", "SÁB"}
	out := make([]weekChip, 0, 7)
	start := anchor.AddDate(0, 0, -3)
	for i := 0; i < 7; i++ {
		d := start.AddDate(0, 0, i)
		dayFC := fc
		if !sameUTCDay(d, anchor) {
			dayFC.Phase = phaseOn(d, fc)
		}
		chip := weekChip{
			Date:       d.Format(time.DateOnly),
			Weekday:    weekdays[d.Weekday()],
			Day:        d.Day(),
			Today:      sameUTCDay(d, anchor),
			Phase:      dayFC.Phase,
			PhaseClass: phaseDot(dayFC.Phase),
		}
		if chip.Today && labelToday {
			chip.Weekday = "HOJE"
		}
		out = append(out, chip)
	}
	return out
}

func phaseOn(day time.Time, fc cycle.Result) cycle.Phase {
	if fc.Phase == cycle.PhasePregnancy {
		return cycle.PhasePregnancy
	}
	if !day.Before(fc.FertileStart) && !day.After(fc.FertileEnd) {
		if sameUTCDay(day, fc.OvulationOn) {
			return cycle.PhaseOvulation
		}
		return cycle.PhaseFertile
	}
	if !fc.NextPeriodOn.IsZero() {
		start := fc.NextPeriodOn.AddDate(0, 0, -fc.CycleLength)
		cycleDay := int(day.Sub(start).Hours()/24) + 1
		if cycleDay >= 1 && cycleDay <= fc.PeriodLength {
			return cycle.PhaseMenstruation
		}
		if day.Before(fc.OvulationOn) {
			return cycle.PhaseFollicular
		}
	}
	return cycle.PhaseLuteal
}

func phaseDot(p cycle.Phase) string {
	switch p {
	case cycle.PhaseMenstruation, cycle.PhaseOvulation:
		return "bg-primary"
	case cycle.PhaseFertile:
		return "bg-tertiary-container"
	case cycle.PhaseLuteal:
		return "bg-secondary"
	default:
		return "bg-outline-variant"
	}
}

func sameUTCDay(a, b time.Time) bool {
	aa := time.Date(a.Year(), a.Month(), a.Day(), 0, 0, 0, 0, time.UTC)
	bb := time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, time.UTC)
	return aa.Equal(bb)
}

type ringViewModel struct {
	PeriodDash, FollicularDash, FertileDash, LutealDash float64
	PeriodOff, FollicularOff, FertileOff, LutealOff     float64
	PinX, PinY                                          float64
}

func ringView(fc cycle.Result) ringViewModel {
	const circ = 615.75
	length := fc.CycleLength
	if length < 1 {
		length = 28
	}
	period := float64(fc.PeriodLength) / float64(length) * circ
	ovDay := float64(length - cycle.LutealDays)
	if ovDay < 1 {
		ovDay = 14
	}
	fertileDays := 7.0
	follicularDays := ovDay - 5 - float64(fc.PeriodLength)
	if follicularDays < 0 {
		follicularDays = 0
	}
	lutealDays := float64(length) - float64(fc.PeriodLength) - follicularDays - fertileDays
	if lutealDays < 0 {
		lutealDays = 0
	}
	follicular := follicularDays / float64(length) * circ
	fertile := fertileDays / float64(length) * circ
	luteal := lutealDays / float64(length) * circ
	// SVG circle starts at 3 o'clock; the template rotates -90° so day 0 sits at 12 o'clock.
	theta := 2 * math.Pi * float64(fc.CycleDay) / float64(length)
	return ringViewModel{
		PeriodDash:     period,
		FollicularDash: follicular,
		FertileDash:    fertile,
		LutealDash:     luteal,
		PeriodOff:      0,
		FollicularOff:  -period,
		FertileOff:     -(period + follicular),
		LutealOff:      -(period + follicular + fertile),
		PinX:           120 + 98*math.Cos(theta),
		PinY:           120 + 98*math.Sin(theta),
	}
}
