package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/session"

	"github.com/puppe1990/aura-ciclos/internal/anon"
	appi18n "github.com/puppe1990/aura-ciclos/internal/i18n"
	"github.com/puppe1990/aura-ciclos/internal/store"
)

func TestAuth_AnonimoPost_createsRecoveryCode(t *testing.T) {
	s, err := store.NewSQLiteStore(":memory:", "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	h := NewAuthHandler(setupTestViews(t), s, testSite(), s.Sessions(), cais.Config{}, appi18n.NewCatalog("pt"))

	req := httptest.NewRequest(http.MethodPost, "/anonimo", nil)
	rr := httptest.NewRecorder()
	h.AnonimoPost(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	if !strings.Contains(body, "código") && !strings.Contains(body, "Código") {
		t.Errorf("missing recovery code copy: %s", body)
	}
}

func TestAuth_LoginPost_recoveryCode(t *testing.T) {
	s, err := store.NewSQLiteStore(":memory:", "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	h := NewAuthHandler(setupTestViews(t), s, testSite(), s.Sessions(), cais.Config{}, appi18n.NewCatalog("pt"))

	code := "AB12-CD34-EF56-GH78"
	hash, err := session.HashPassword(anon.NormalizeCode(code))
	if err != nil {
		t.Fatal(err)
	}
	id, err := s.CreateUser(anon.EmailFor(code), hash)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.MarkAnonymous(id, true); err != nil {
		t.Fatal(err)
	}

	form := url.Values{"recovery_code": {code}}
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	h.LoginPost(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
	if rr.Header().Get("Location") != afterAuthPath {
		t.Errorf("Location = %q, want %s", rr.Header().Get("Location"), afterAuthPath)
	}
}
