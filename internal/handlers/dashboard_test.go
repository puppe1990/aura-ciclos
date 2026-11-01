package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/flash"
	"github.com/puppe1990/amarra-cais/pkg/cais/session"

	appi18n "github.com/puppe1990/aura-ciclos/internal/i18n"
)

func TestDashboardHandler_redirectsToCiclo(t *testing.T) {
	h := NewDashboardHandler(setupTestViews(t), setupTestStore(t), testSite(), appi18n.DefaultCatalog(), cais.Config{})

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want 303", rr.Code)
	}
	if rr.Header().Get("Location") != afterAuthPath {
		t.Errorf("Location = %q, want %s", rr.Header().Get("Location"), afterAuthPath)
	}
}

func TestDashboardHandler_includesFlashStillRedirects(t *testing.T) {
	h := NewDashboardHandler(setupTestViews(t), setupTestStore(t), testSite(), appi18n.DefaultCatalog(), cais.Config{})

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	req = flash.WithMessage(req, flash.Message{Kind: "notice", Message: "Welcome back!"})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rr.Code)
	}
}

func TestDashboardHandler_signedInRedirectsToCiclo(t *testing.T) {
	h := NewDashboardHandler(setupTestViews(t), setupTestStore(t), testSite(), appi18n.DefaultCatalog(), cais.Config{})

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	req = session.WithUserID(req, 1)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rr.Code)
	}
	if rr.Header().Get("Location") != afterAuthPath {
		t.Errorf("Location = %q, want %s", rr.Header().Get("Location"), afterAuthPath)
	}
}
