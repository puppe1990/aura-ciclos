package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/session"

	appi18n "github.com/puppe1990/aura-ciclos/internal/i18n"
	"github.com/puppe1990/aura-ciclos/internal/models"
	"github.com/puppe1990/aura-ciclos/internal/store"
	"github.com/puppe1990/aura-ciclos/internal/testdata"
)

func fixedNow() time.Time {
	return testdata.Day(2026, 11, 1)
}

func signedInTracker(t *testing.T) (*TrackerHandler, store.Store, models.User) {
	t.Helper()
	s := setupTestStore(t)
	f := testdata.New()
	hash, err := session.HashPassword("password123")
	if err != nil {
		t.Fatal(err)
	}
	id, err := s.CreateUser(testdata.Email(f), hash)
	if err != nil {
		t.Fatal(err)
	}
	last := testdata.Day(2026, 10, 19)
	for _, start := range testdata.RegularStarts(last, 7, 28) {
		if err := s.UpsertPeriod(models.Period{
			UserID:    id,
			StartedOn: start,
			EndedOn:   testdata.EndOf(start, 5),
			Flow:      testdata.Flow(f),
		}); err != nil {
			t.Fatal(err)
		}
	}
	u, err := s.FindUser(id)
	if err != nil {
		t.Fatal(err)
	}
	p, err := s.GetProfile(id)
	if err != nil {
		t.Fatal(err)
	}
	p.DisplayName = testdata.DisplayName(f)
	if err := s.SaveProfile(p); err != nil {
		t.Fatal(err)
	}
	h := NewTrackerHandler(setupTestViews(t), s, testSite(), appi18n.NewCatalog("pt"), cais.Config{})
	h.now = fixedNow
	return h, s, u
}

func requestWithUser(method, path string, userID int64, body url.Values) *http.Request {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, strings.NewReader(body.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	return session.WithUserID(req, userID)
}

func TestCiclo_rendersDay14AndFertileWindow(t *testing.T) {
	h, _, u := signedInTracker(t)
	req := requestWithUser(http.MethodGet, "/ciclo", u.ID, nil)
	rr := httptest.NewRecorder()
	h.Ciclo(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	if !strings.Contains(body, "Dia 14") {
		t.Errorf("missing cycle day, got: %s", body)
	}
	if !strings.Contains(body, "Janela") && !strings.Contains(body, "fértil") && !strings.Contains(body, "Fértil") && !strings.Contains(body, "Ovul") {
		t.Errorf("missing fertile/ovulation copy, got: %s", body)
	}
	if !strings.Contains(body, `href="/ciclo"`) {
		t.Errorf("missing ciclo nav")
	}
}

func TestCiclo_rendersStitchHero(t *testing.T) {
	h, _, u := signedInTracker(t)
	req := requestWithUser(http.MethodGet, "/ciclo", u.ID, nil)
	rr := httptest.NewRecorder()
	h.Ciclo(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Ritmo Biológico Ativo",
		"Olá,",
		"Registrar Sintomas de Hoje",
		"IA Adaptativa Aura",
		"Fisiologia",
		"LH &amp; Estrogênio",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("missing Stitch ciclo copy %q", want)
		}
	}
}

func TestRegistro_savesSymptomsFromFaker(t *testing.T) {
	h, s, u := signedInTracker(t)
	f := testdata.New()
	ids := testdata.SymptomIDs(f, 3)
	form := url.Values{}
	form.Set("logged_on", "2026-11-01")
	form.Set("mood", "ansiedade")
	form.Set("energy", "2")
	form.Set("sleep_hours", "7")
	form.Set("water_glasses", "5")
	for _, id := range ids {
		form.Add("symptoms", id)
	}
	req := requestWithUser(http.MethodPost, "/registro", u.ID, form)
	rr := httptest.NewRecorder()
	h.RegistroPost(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
	got, err := s.GetDailyLog(u.ID, testdata.Day(2026, 11, 1))
	if err != nil {
		t.Fatal(err)
	}
	if got.Mood != "ansiedade" || len(got.SymptomIDs) != 3 {
		t.Fatalf("saved log = %+v", got)
	}
}

func TestRegistro_rendersCatalog(t *testing.T) {
	h, _, u := signedInTracker(t)
	req := requestWithUser(http.MethodGet, "/registro", u.ID, nil)
	rr := httptest.NewRecorder()
	h.Registro(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Cólicas") {
		t.Errorf("missing symptom catalog, got: %s", rr.Body.String())
	}
}

func TestRegistro_rendersStitchDiary(t *testing.T) {
	h, _, u := signedInTracker(t)
	req := requestWithUser(http.MethodGet, "/registro", u.ID, nil)
	rr := httptest.NewRecorder()
	h.Registro(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Fluxo Menstrual",
		"Humor &amp; Energia",
		"Estilo de Vida",
		"Salvar Registro do Dia",
		"Bateria Corporal",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("missing Stitch registro copy %q", want)
		}
	}
}

func TestRelatorios_showsRegularity(t *testing.T) {
	h, _, u := signedInTracker(t)
	req := requestWithUser(http.MethodGet, "/relatorios", u.ID, nil)
	rr := httptest.NewRecorder()
	h.Relatorios(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "Regular") && !strings.Contains(body, "regular") {
		t.Errorf("missing regularity, got: %s", body)
	}
}

func TestRelatorios_rendersStitchModes(t *testing.T) {
	h, _, u := signedInTracker(t)
	req := requestWithUser(http.MethodGet, "/relatorios", u.ID, nil)
	rr := httptest.NewRecorder()
	h.Relatorios(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Selecione seu Modo",
		"Acompanhar Ciclo",
		"Crescer a Família",
		"Acompanhar Gravidez",
		"Relatórios Clínicos",
		"Consulta Médica Facilitada",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("missing Stitch relatórios copy %q", want)
		}
	}
	if strings.Contains(body, "Criptografia de ponta a ponta") || strings.Contains(body, "E2EE") {
		t.Error("relatórios must not claim E2EE")
	}
}

func TestApoio_mayaAnswersColica(t *testing.T) {
	h, s, u := signedInTracker(t)
	form := url.Values{"question": {"estou com cólica forte"}}
	req := requestWithUser(http.MethodPost, "/apoio/maya", u.ID, form)
	rr := httptest.NewRecorder()
	h.MayaPost(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
	msgs, err := s.ListMayaMessages(u.ID, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) < 2 {
		t.Fatalf("messages = %d, want user+maya", len(msgs))
	}
}

func TestApoio_listsArticles(t *testing.T) {
	h, _, u := signedInTracker(t)
	req := requestWithUser(http.MethodGet, "/apoio", u.ID, nil)
	rr := httptest.NewRecorder()
	h.Apoio(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "janela fértil") && !strings.Contains(rr.Body.String(), "Janela fértil") && !strings.Contains(rr.Body.String(), "O que é a janela") {
		t.Errorf("missing articles, got: %s", rr.Body.String())
	}
}

func TestApoio_rendersStitchSupport(t *testing.T) {
	h, _, u := signedInTracker(t)
	req := requestWithUser(http.MethodGet, "/apoio", u.ID, nil)
	rr := httptest.NewRecorder()
	h.Apoio(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Seu Refúgio Seguro",
		"Conversar com a Maya",
		"Central de Lembretes",
		"Biblioteca Clínica",
		"SQLite",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("missing Stitch apoio copy %q", want)
		}
	}
	if strings.Contains(body, "Criptografia Ponta a Ponta") || strings.Contains(body, "(E2EE)") {
		t.Error("apoio must not claim E2EE")
	}
}

func TestModo_switchesToPregnancy(t *testing.T) {
	h, s, u := signedInTracker(t)
	form := url.Values{"mode": {"pregnancy"}, "pregnancy_lmp": {"2026-08-09"}}
	req := requestWithUser(http.MethodPost, "/modo", u.ID, form)
	rr := httptest.NewRecorder()
	h.ModoPost(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
	p, err := s.GetProfile(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if p.Mode != models.ModePregnancy {
		t.Errorf("mode = %q", p.Mode)
	}
}

func TestPeriodoPost_startsPeriodToday(t *testing.T) {
	h, s, u := signedInTracker(t)
	form := url.Values{"started_on": {"2026-11-01"}, "flow": {"medio"}}
	req := requestWithUser(http.MethodPost, "/ciclo/periodo", u.ID, form)
	rr := httptest.NewRecorder()
	h.PeriodoPost(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d", rr.Code)
	}
	periods, err := s.ListPeriods(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, p := range periods {
		if p.StartedOn.Equal(testdata.Day(2026, 11, 1)) {
			found = true
		}
	}
	if !found {
		t.Fatal("period was not stored")
	}
}

func TestReminderPost_savesPill(t *testing.T) {
	h, s, u := signedInTracker(t)
	form := url.Values{"kind": {"pill"}, "hour": {"8"}, "minute": {"0"}}
	req := requestWithUser(http.MethodPost, "/lembretes", u.ID, form)
	rr := httptest.NewRecorder()
	h.LembretePost(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
	rems, err := s.ListReminders(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(rems) != 1 || rems[0].Kind != "pill" {
		t.Fatalf("reminders = %+v", rems)
	}
}
