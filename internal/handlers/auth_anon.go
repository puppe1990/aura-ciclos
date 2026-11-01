package handlers

import (
	"errors"
	"net/http"

	"github.com/puppe1990/amarra-cais/pkg/cais/flash"
	"github.com/puppe1990/amarra-cais/pkg/cais/httpx"
	"github.com/puppe1990/amarra-cais/pkg/cais/session"

	"github.com/puppe1990/aura-ciclos/internal/anon"
	"github.com/puppe1990/aura-ciclos/internal/store"
)

func (h *AuthHandler) Anonimo(w http.ResponseWriter, r *http.Request) {
	if _, ok := session.UserID(r); ok {
		http.Redirect(w, r, afterAuthPath, http.StatusSeeOther)
		return
	}
	h.renderAuth(w, r, "anonimo", map[string]any{
		"Title": h.t(r, "auth.anon_title"),
	}, 0)
}

func (h *AuthHandler) AnonimoPost(w http.ResponseWriter, r *http.Request) {
	if err := httpx.ParseFormOrJSON(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var (
		code   string
		userID int64
		err    error
	)
	for attempt := 0; attempt < 5; attempt++ {
		code, err = anon.RandomRecoveryCode()
		if err != nil {
			httpx.ServerError(w, err, h.cfg)
			return
		}
		hash, hashErr := session.HashPassword(anon.NormalizeCode(code))
		if hashErr != nil {
			httpx.ServerError(w, hashErr, h.cfg)
			return
		}
		userID, err = h.store.CreateUser(anon.EmailFor(code), hash)
		if err == nil {
			break
		}
		if errors.Is(err, store.ErrEmailTaken) {
			continue
		}
		httpx.ServerError(w, err, h.cfg)
		return
	}
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	if err := h.store.MarkAnonymous(userID, true); err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	if err := session.SignIn(w, h.sessions, r, userID, session.CookieOptionsFromConfig(h.cfg)); err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	flash.Set(w, "notice", h.t(r, "auth.anon_created"), h.cfg.CookieSecure())
	http.Redirect(w, r, afterAuthPath, http.StatusSeeOther)
}
