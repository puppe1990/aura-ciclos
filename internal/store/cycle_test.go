package store

import (
	"testing"

	"github.com/puppe1990/aura-ciclos/internal/cycle"
	"github.com/puppe1990/aura-ciclos/internal/models"
	"github.com/puppe1990/aura-ciclos/internal/testdata"
)

func TestStore_profileDefaultsAndModeChange(t *testing.T) {
	s := newTestStore(t)
	f := testdata.New()
	id, err := s.CreateUser(testdata.Email(f), "hash")
	if err != nil {
		t.Fatal(err)
	}
	p, err := s.GetProfile(id)
	if err != nil {
		t.Fatal(err)
	}
	if p.Mode != models.ModeCycle || p.TypicalCycleDays != 28 {
		t.Fatalf("default profile = %+v", p)
	}
	p.Mode = models.ModeTrying
	p.DisplayName = testdata.DisplayName(f)
	if err := s.SaveProfile(p); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetProfile(id)
	if err != nil {
		t.Fatal(err)
	}
	if got.Mode != models.ModeTrying || got.DisplayName != p.DisplayName {
		t.Fatalf("saved profile = %+v", got)
	}
}

func TestStore_periodsRoundTripAndForecastInput(t *testing.T) {
	s := newTestStore(t)
	f := testdata.New()
	id, err := s.CreateUser(testdata.Email(f), "hash")
	if err != nil {
		t.Fatal(err)
	}
	last := testdata.Day(2026, 10, 19)
	for _, start := range testdata.RegularStarts(last, 4, 28) {
		p := models.Period{
			UserID:    id,
			StartedOn: start,
			EndedOn:   testdata.EndOf(start, 5),
			Flow:      testdata.Flow(f),
		}
		if err := s.UpsertPeriod(p); err != nil {
			t.Fatal(err)
		}
	}
	got, err := s.ListPeriods(id)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 4 {
		t.Fatalf("len(periods) = %d", len(got))
	}
	fc := cycle.Forecast(cycle.Input{
		Today:   testdata.Day(2026, 11, 1),
		Periods: toCyclePeriods(got),
	})
	if fc.CycleDay != 14 {
		t.Errorf("CycleDay = %d, want 14", fc.CycleDay)
	}
}

func TestStore_dailyLogSymptomsAndBBT(t *testing.T) {
	s := newTestStore(t)
	f := testdata.New()
	id, err := s.CreateUser(testdata.Email(f), "hash")
	if err != nil {
		t.Fatal(err)
	}
	temp := testdata.Celsius(f, 36.2, 36.9)
	log := models.DailyLog{
		UserID:        id,
		LoggedOn:      testdata.Day(2026, 11, 1),
		SymptomIDs:    testdata.SymptomIDs(f, 3),
		Mood:          "ansiedade",
		Energy:        2,
		SleepHours:    6.5,
		WaterGlasses:  4,
		BBTCelsius:    &temp,
		OvulationTest: "negativo",
		Notes:         f.Sentence(5),
	}
	if err := s.UpsertDailyLog(log); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetDailyLog(id, log.LoggedOn)
	if err != nil {
		t.Fatal(err)
	}
	if got.Mood != "ansiedade" || len(got.SymptomIDs) != 3 {
		t.Fatalf("log = %+v", got)
	}
	if got.BBTCelsius == nil || *got.BBTCelsius != temp {
		t.Fatalf("bbt = %v want %v", got.BBTCelsius, temp)
	}
	log.Mood = "calma"
	if err := s.UpsertDailyLog(log); err != nil {
		t.Fatal(err)
	}
	got, err = s.GetDailyLog(id, log.LoggedOn)
	if err != nil {
		t.Fatal(err)
	}
	if got.Mood != "calma" {
		t.Errorf("upsert did not replace mood: %q", got.Mood)
	}
}

func TestStore_remindersAndMaya(t *testing.T) {
	s := newTestStore(t)
	f := testdata.New()
	id, err := s.CreateUser(testdata.Email(f), "hash")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.UpsertReminder(models.Reminder{UserID: id, Kind: "pill", Hour: 8, Minute: 0, Enabled: true}); err != nil {
		t.Fatal(err)
	}
	rems, err := s.ListReminders(id)
	if err != nil {
		t.Fatal(err)
	}
	if len(rems) != 1 || rems[0].Kind != "pill" {
		t.Fatalf("reminders = %+v", rems)
	}
	if err := s.AddMayaMessage(id, "user", "estou com cólica"); err != nil {
		t.Fatal(err)
	}
	if err := s.AddMayaMessage(id, "maya", "Calor local ajuda."); err != nil {
		t.Fatal(err)
	}
	msgs, err := s.ListMayaMessages(id, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 2 {
		t.Fatalf("messages = %d", len(msgs))
	}
}

func TestStore_anonUserFlag(t *testing.T) {
	s := newTestStore(t)
	id, err := s.CreateUser("anon.abcd1234abcd1234@aura.local", "hash")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.MarkAnonymous(id, true); err != nil {
		t.Fatal(err)
	}
	u, err := s.FindUser(id)
	if err != nil {
		t.Fatal(err)
	}
	if !u.Anonymous {
		t.Fatal("user should be anonymous")
	}
}

func toCyclePeriods(in []models.Period) []cycle.Period {
	out := make([]cycle.Period, len(in))
	for i, p := range in {
		out[i] = cycle.Period{StartedOn: p.StartedOn, EndedOn: p.EndedOn, Flow: p.Flow}
	}
	return out
}

func TestStore_migrationsIncludeCycleTables(t *testing.T) {
	s := newTestStore(t)
	for _, table := range []string{"profiles", "periods", "daily_logs", "reminders", "maya_messages"} {
		var name string
		if err := s.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name = ?", table).Scan(&name); err != nil {
			t.Fatalf("%s table not found: %v", table, err)
		}
	}
}
