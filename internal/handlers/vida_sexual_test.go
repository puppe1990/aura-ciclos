package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/puppe1990/aura-ciclos/internal/models"
	"github.com/puppe1990/aura-ciclos/internal/testdata"
)

func TestVidaSexual_rendersFloScreen(t *testing.T) {
	h, _, u := signedInTracker(t)
	req := requestWithUser(http.MethodGet, "/registro/vida-sexual", u.ID, nil)
	rr := httptest.NewRecorder()
	h.VidaSexual(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Vida sexual",
		"Preservativo relações",
		"Protegido",
		"Desprotegido",
		"Orgasmo feminino",
		">Sim<",
		">Não<",
		"vezes",
		"Salvar",
		"Análise",
		`name="sex"`,
		`name="orgasm"`,
		`name="sex_times"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("missing %q", want)
		}
	}
	if strings.Contains(body, `name="recovery_code"`) {
		t.Error("vida sexual must not mention recovery code")
	}
}

func TestVidaSexualPost_savesWithoutWipingSymptoms(t *testing.T) {
	h, s, u := signedInTracker(t)
	day := testdata.Day(2026, 11, 1)
	if err := s.UpsertDailyLog(models.DailyLog{
		UserID:       u.ID,
		LoggedOn:     day,
		SymptomIDs:   []string{"colicas", "acne"},
		Mood:         "ansiedade",
		WaterGlasses: 4,
	}); err != nil {
		t.Fatal(err)
	}

	form := url.Values{
		"logged_on": {"2026-11-01"},
		"sex":       {"desprotegida"},
		"orgasm":    {"nao"},
		"sex_times": {"2"},
	}
	req := requestWithUser(http.MethodPost, "/registro/vida-sexual", u.ID, form)
	rr := httptest.NewRecorder()
	h.VidaSexualPost(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
	if loc := rr.Header().Get("Location"); !strings.Contains(loc, "/registro/vida-sexual") {
		t.Errorf("Location = %q", loc)
	}
	got, err := s.GetDailyLog(u.ID, day)
	if err != nil {
		t.Fatal(err)
	}
	if got.Sex != "desprotegida" || got.Orgasm != "nao" || got.SexTimes != 2 {
		t.Fatalf("saved = sex=%q orgasm=%q times=%d", got.Sex, got.Orgasm, got.SexTimes)
	}
	if got.Mood != "ansiedade" || got.WaterGlasses != 4 || len(got.SymptomIDs) != 2 {
		t.Fatalf("wiped other fields: %+v", got)
	}
}

func TestRegistroPost_keepsSexDetails(t *testing.T) {
	h, s, u := signedInTracker(t)
	day := testdata.Day(2026, 11, 1)
	if err := s.UpsertDailyLog(models.DailyLog{
		UserID:     u.ID,
		LoggedOn:   day,
		Sex:        "protegida",
		SexTimes:   4,
		Orgasm:     "sim",
		SymptomIDs: []string{"colicas"},
	}); err != nil {
		t.Fatal(err)
	}
	form := url.Values{
		"logged_on":     {"2026-11-01"},
		"mood":          {"calma"},
		"energy":        {"3"},
		"sleep_hours":   {"8"},
		"water_glasses": {"6"},
		"symptoms":      {"colicas"},
	}
	req := requestWithUser(http.MethodPost, "/registro", u.ID, form)
	rr := httptest.NewRecorder()
	h.RegistroPost(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d", rr.Code)
	}
	got, err := s.GetDailyLog(u.ID, day)
	if err != nil {
		t.Fatal(err)
	}
	if got.Sex != "protegida" || got.SexTimes != 4 || got.Orgasm != "sim" {
		t.Fatalf("registro post wiped sex details: sex=%q times=%d orgasm=%q", got.Sex, got.SexTimes, got.Orgasm)
	}
	if got.Mood != "calma" {
		t.Errorf("mood = %q", got.Mood)
	}
}

func TestRegistro_linksToVidaSexual(t *testing.T) {
	h, _, u := signedInTracker(t)
	req := requestWithUser(http.MethodGet, "/registro", u.ID, nil)
	rr := httptest.NewRecorder()
	h.Registro(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "/registro/vida-sexual") {
		t.Error("registro must link to /registro/vida-sexual")
	}
	if !strings.Contains(body, "Vida sexual") {
		t.Error("registro must show Vida sexual")
	}
}

func TestVidaSexual_analiseCountsLogs(t *testing.T) {
	h, s, u := signedInTracker(t)
	if err := s.UpsertDailyLog(models.DailyLog{
		UserID: u.ID, LoggedOn: testdata.Day(2026, 10, 28),
		Sex: "protegida", SexTimes: 1, Orgasm: "sim",
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.UpsertDailyLog(models.DailyLog{
		UserID: u.ID, LoggedOn: testdata.Day(2026, 10, 30),
		Sex: "desprotegida", SexTimes: 2, Orgasm: "nao",
	}); err != nil {
		t.Fatal(err)
	}
	req := requestWithUser(http.MethodGet, "/registro/vida-sexual?analise=1", u.ID, nil)
	rr := httptest.NewRecorder()
	h.VidaSexual(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, want := range []string{"Análise", "2 dias", "Protegido", "Desprotegido", "3 vezes"} {
		if !strings.Contains(body, want) {
			t.Errorf("missing analysis copy %q in %s", want, body)
		}
	}
}

func TestSummarizeSex_countsActivity(t *testing.T) {
	got := summarizeSex([]models.DailyLog{
		{Sex: "protegida", SexTimes: 1, Orgasm: "sim"},
		{Mood: "calma"},
		{Sex: "desprotegida", SexTimes: 2, Orgasm: "nao"},
		{SexTimes: 1},
	})
	if got.Days != 3 || got.ProtectedDays != 1 || got.UnprotectedDays != 1 {
		t.Fatalf("days=%d prot=%d unprot=%d", got.Days, got.ProtectedDays, got.UnprotectedDays)
	}
	if got.OrgasmYes != 1 || got.OrgasmNo != 1 || got.Times != 4 {
		t.Fatalf("orgasm yes=%d no=%d times=%d", got.OrgasmYes, got.OrgasmNo, got.Times)
	}
}
