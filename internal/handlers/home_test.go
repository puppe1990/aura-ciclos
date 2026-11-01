package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/session"

	appi18n "github.com/puppe1990/aura-ciclos/internal/i18n"
)

func newHomeHandler(t *testing.T) *HomeHandler {
	t.Helper()
	return NewHomeHandler(setupTestViews(t), testSite(), appi18n.DefaultCatalog(), cais.Config{})
}

func TestHomeHandler_Returns200(t *testing.T) {
	h := newHomeHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestHomeHandler_RendersHTML(t *testing.T) {
	h := newHomeHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, `id="amarra-main"`) {
		t.Errorf("body missing #amarra-main, got: %s", body)
	}
	if !strings.Contains(body, "Aura Ciclos") {
		t.Errorf("body missing heading, got: %s", body)
	}
}

func TestHomeHandler_stitchLanding(t *testing.T) {
	h := newHomeHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	body := rr.Body.String()
	if strings.Contains(body, "01 —") || strings.Contains(body, "01-") {
		t.Error("guest home still has Amarra rails numbered aside")
	}
	if !strings.Contains(body, `viewBox="0 0 240 240"`) {
		t.Error("guest home missing Stitch cycle ring preview")
	}
	if !strings.Contains(body, `data-testid="amarra-ready"`) {
		t.Error("missing amarra-ready")
	}
	if !strings.Contains(body, "Active biological rhythm") && !strings.Contains(body, "Ritmo Biológico") {
		t.Error("missing Stitch kicker")
	}
}

func TestHomeHandler_ContentType(t *testing.T) {
	h := newHomeHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}
}

// #208: a visitor without a session must not be offered app chrome that bounces
// (Dashboard) or signs out a session that does not exist.
func TestHomeHandler_anonymousHidesAuthChrome(t *testing.T) {
	h := newHomeHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	body := rr.Body.String()
	if strings.Contains(body, `action="/logout"`) {
		t.Errorf("anonymous page must not offer sign-out (#208), got: %s", body)
	}
	if strings.Contains(body, `href="/dashboard"`) {
		t.Errorf("anonymous page must not link to /dashboard (#208), got: %s", body)
	}
	if strings.Contains(body, `href="/ciclo"`) {
		t.Errorf("anonymous page must not link to /ciclo, got: %s", body)
	}
}

func TestHomeHandler_loggedInRedirectsToCiclo(t *testing.T) {
	h := newHomeHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
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
