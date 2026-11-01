package handlers

import (
	"net/http"

	"github.com/puppe1990/amarra-cais/pkg/amarra/view"
	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/i18n"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"

	"github.com/puppe1990/aura-ciclos/internal/store"
)

type DashboardHandler struct{}

func NewDashboardHandler(_ *view.Renderer, _ store.Store, _ meta.Site, _ *i18n.Catalog, _ cais.Config) *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, afterAuthPath, http.StatusSeeOther)
}
