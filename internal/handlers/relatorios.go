package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/puppe1990/aura-ciclos/internal/cycle"
	"github.com/puppe1990/aura-ciclos/internal/models"
	"github.com/puppe1990/aura-ciclos/internal/symptoms"
)

func (h *TrackerHandler) Relatorios(w http.ResponseWriter, r *http.Request) {
	data, ok := h.load(w, r)
	if !ok {
		return
	}
	counts := map[string]int{}
	for _, log := range data.logs {
		for _, id := range log.SymptomIDs {
			counts[id]++
		}
	}
	type symptomStat struct {
		Label string
		Count int
	}
	var top []symptomStat
	for id, n := range counts {
		s, ok := symptoms.ByID(id)
		if !ok {
			continue
		}
		top = append(top, symptomStat{Label: s.Label, Count: n})
	}
	sort.Slice(top, func(i, j int) bool { return top[i].Count > top[j].Count })
	if len(top) > 6 {
		top = top[:6]
	}
	h.page(w, r, "relatorios", map[string]any{
		"Title":         "Relatórios",
		"ActiveNav":     "relatorios",
		"Name":          greetName(data.user, data.profile),
		"Anonymous":     data.user.Anonymous || data.profile.Anonymous,
		"Mode":          data.profile.Mode,
		"ModeLabel":     cycle.ModeLabel(cycle.Mode(data.profile.Mode)),
		"Forecast":      data.fc,
		"PhaseLabel":    cycle.PhaseLabel(data.fc.Phase),
		"RegularLabel":  regularityLabel(data.fc),
		"StdDevLabel":   fmt.Sprintf("%.1f dias", data.fc.StdDev),
		"Lengths":       data.fc.Lengths,
		"TopSymptoms":   top,
		"Profile":       data.profile,
		"PregnancyNote": cycle.PregnancyNote(data.fc.PregnancyWeek),
		"LMP":           formatOptionalDay(data.profile.PregnancyLMP),
		"Trying":        data.profile.Mode == models.ModeTrying,
		"Pregnant":      data.profile.Mode == models.ModePregnancy,
	})
}
