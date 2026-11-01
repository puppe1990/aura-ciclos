package handlers

import (
	"net/http"
	"strings"

	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/httpx"
	"github.com/puppe1990/amarra-cais/pkg/cais/session"

	"github.com/puppe1990/aura-ciclos/internal/articles"
	"github.com/puppe1990/aura-ciclos/internal/cycle"
	"github.com/puppe1990/aura-ciclos/internal/maya"
)

func (h *TrackerHandler) Apoio(w http.ResponseWriter, r *http.Request) {
	data, ok := h.load(w, r)
	if !ok {
		return
	}
	msgs, err := h.store.ListMayaMessages(data.user.ID, 20)
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	rems, err := h.store.ListReminders(data.user.ID)
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	suggestion := maya.Reply("como estou hoje", maya.Context{
		Phase:     data.fc.Phase,
		CycleDay:  data.fc.CycleDay,
		FertileOn: data.fc.ConceptionChance > 0,
		ChancePct: data.fc.ConceptionChance,
	})
	h.page(w, r, "apoio", map[string]any{
		"Title":          "Apoio & IA",
		"ActiveNav":      "apoio",
		"Name":           greetName(data.user, data.profile),
		"Anonymous":      data.user.Anonymous || data.profile.Anonymous,
		"Mode":           data.profile.Mode,
		"ModeLabel":      cycle.ModeLabel(cycle.Mode(data.profile.Mode)),
		"Forecast":       data.fc,
		"PhaseLabel":     cycle.PhaseLabel(data.fc.Phase),
		"PhaseBadge":     cycle.Physiology(data.fc.Phase).Badge,
		"Messages":       msgs,
		"Articles":       articles.All(),
		"Reminders":      rems,
		"Suggestion":     suggestion,
		"NextPeriodDate": formatDayBR(data.fc.NextPeriodOn),
	})
}

func (h *TrackerHandler) MayaPost(w http.ResponseWriter, r *http.Request) {
	id, ok := session.UserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if err := httpx.ParseFormOrJSON(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	q := strings.TrimSpace(r.FormValue("question"))
	if q == "" {
		http.Redirect(w, r, "/apoio", http.StatusSeeOther)
		return
	}
	data, ok := h.load(w, r)
	if !ok {
		return
	}
	reply := maya.Reply(q, maya.Context{
		Phase:     data.fc.Phase,
		CycleDay:  data.fc.CycleDay,
		FertileOn: data.fc.ConceptionChance > 0,
		ChancePct: data.fc.ConceptionChance,
	})
	if err := h.store.AddMayaMessage(id, "user", q); err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	if err := h.store.AddMayaMessage(id, "maya", reply); err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	http.Redirect(w, r, "/apoio", http.StatusSeeOther)
}

func (h *TrackerHandler) Artigo(w http.ResponseWriter, r *http.Request, slug string) {
	data, ok := h.load(w, r)
	if !ok {
		return
	}
	art, found := articles.BySlug(slug)
	if !found {
		http.NotFound(w, r)
		return
	}
	h.page(w, r, "artigo", map[string]any{
		"Title":      art.Title,
		"ActiveNav":  "apoio",
		"Name":       greetName(data.user, data.profile),
		"Anonymous":  data.user.Anonymous || data.profile.Anonymous,
		"ModeLabel":  cycle.ModeLabel(cycle.Mode(data.profile.Mode)),
		"Article":    art,
		"Paragraphs": strings.Split(art.Body, "\n\n"),
	})
}

func ArtigoRoute(h *TrackerHandler) http.HandlerFunc {
	return cais.StringParam("slug", h.Artigo)
}
