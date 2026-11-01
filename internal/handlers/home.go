package handlers

import (
	"net/http"

	"github.com/puppe1990/amarra-cais/pkg/amarra/view"
	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/i18n"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"
	"github.com/puppe1990/amarra-cais/pkg/cais/session"
)

type HomeHandler struct {
	views *view.Renderer
	site  meta.Site
	cfg   cais.Config
}

func NewHomeHandler(views *view.Renderer, site meta.Site, _ *i18n.Catalog, cfg cais.Config) *HomeHandler {
	return &HomeHandler{views: views, site: site, cfg: cfg}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, ok := session.UserID(r); ok {
		http.Redirect(w, r, afterAuthPath, http.StatusSeeOther)
		return
	}
	view.Write(w, r, h.views, view.Page{
		Layout: "app",
		Name:   "home",
		Data: amarraData(r, h.site, map[string]any{
			"ActiveNav": "home",
		}),
	}, h.cfg)
}
