package cycle

import (
	"math"
	"testing"
	"time"

	"github.com/puppe1990/aura-ciclos/internal/testdata"
)

func TestForecast_noPeriods_usesDefault28DayCycle(t *testing.T) {
	today := testdata.Day(2026, 11, 1)
	got := Forecast(Input{Today: today})

	if got.CycleLength != 28 {
		t.Errorf("CycleLength = %d, want 28", got.CycleLength)
	}
	if got.PeriodLength != 5 {
		t.Errorf("PeriodLength = %d, want 5", got.PeriodLength)
	}
	if got.CycleDay != 1 {
		t.Errorf("CycleDay = %d, want 1", got.CycleDay)
	}
	if got.Phase != PhaseMenstruation {
		t.Errorf("Phase = %q, want %q", got.Phase, PhaseMenstruation)
	}
	if !got.OvulationOn.Equal(testdata.Day(2026, 11, 14)) {
		t.Errorf("OvulationOn = %s, want 2026-11-14", got.OvulationOn.Format(time.DateOnly))
	}
	if !got.NextPeriodOn.Equal(testdata.Day(2026, 11, 29)) {
		t.Errorf("NextPeriodOn = %s, want 2026-11-29", got.NextPeriodOn.Format(time.DateOnly))
	}
	if got.Confidence != 1 || got.ConfidencePct != 25 {
		t.Errorf("confidence = %d/%d%%, want 1/25", got.Confidence, got.ConfidencePct)
	}
	if got.Irregular {
		t.Error("empty history must not be marked irregular")
	}
}

func periodsFromStarts(starts []time.Time, periodDays int, f interface{ RandomString([]string) string }) []Period {
	out := make([]Period, len(starts))
	for i, start := range starts {
		flow := "medio"
		if f != nil {
			flow = testdata.Flow(testdata.New())
		}
		out[i] = Period{StartedOn: start, EndedOn: testdata.EndOf(start, periodDays), Flow: flow}
	}
	return out
}

func TestForecast_sixRegularCycles_ovulationOnDay14(t *testing.T) {
	today := testdata.Day(2026, 11, 1)
	last := testdata.Day(2026, 10, 19)
	got := Forecast(Input{Today: today, Periods: periodsFromStarts(testdata.RegularStarts(last, 7, 28), 5, testdata.New())})

	if got.CycleDay != 14 {
		t.Errorf("CycleDay = %d, want 14", got.CycleDay)
	}
	if got.CycleLength != 28 {
		t.Errorf("CycleLength = %d, want 28", got.CycleLength)
	}
	if got.Phase != PhaseOvulation {
		t.Errorf("Phase = %q, want ovulation", got.Phase)
	}
	if !got.OvulationOn.Equal(today) {
		t.Errorf("OvulationOn = %s, want today", got.OvulationOn.Format(time.DateOnly))
	}
	if !got.FertileStart.Equal(testdata.Day(2026, 10, 27)) {
		t.Errorf("FertileStart = %s, want 2026-10-27", got.FertileStart.Format(time.DateOnly))
	}
	if !got.FertileEnd.Equal(testdata.Day(2026, 11, 2)) {
		t.Errorf("FertileEnd = %s, want 2026-11-02", got.FertileEnd.Format(time.DateOnly))
	}
	if got.ConceptionChance != 87 {
		t.Errorf("ConceptionChance = %d, want 87", got.ConceptionChance)
	}
	if got.Irregular {
		t.Error("regular 28-day history must not be irregular")
	}
	if got.Confidence != 4 || got.ConfidencePct != 98 {
		t.Errorf("confidence = %d/%d%%, want 4/98", got.Confidence, got.ConfidencePct)
	}
	if !got.NextPeriodOn.Equal(testdata.Day(2026, 11, 16)) {
		t.Errorf("NextPeriodOn = %s, want 2026-11-16", got.NextPeriodOn.Format(time.DateOnly))
	}
	if got.DaysUntilPeriod != 15 {
		t.Errorf("DaysUntilPeriod = %d, want 15", got.DaysUntilPeriod)
	}
}

func TestForecast_irregularStddevMarksIrregular(t *testing.T) {
	today := testdata.Day(2026, 11, 1)
	starts := []time.Time{
		testdata.Day(2026, 4, 1),
		testdata.Day(2026, 4, 25),
		testdata.Day(2026, 5, 30),
		testdata.Day(2026, 6, 21),
		testdata.Day(2026, 7, 31),
		testdata.Day(2026, 8, 21),
		testdata.Day(2026, 9, 28),
	}
	periods := make([]Period, len(starts))
	for i, start := range starts {
		periods[i] = Period{StartedOn: start, EndedOn: testdata.EndOf(start, 5)}
	}
	got := Forecast(Input{Today: today, Periods: periods})
	if !got.Irregular {
		t.Fatalf("Irregular = false, stddev=%.2f", got.StdDev)
	}
	if got.StdDev < 4 {
		t.Errorf("StdDev = %.2f, want >= 4", got.StdDev)
	}
	if got.ConfidencePct > 60 {
		t.Errorf("irregular confidence %d%% must cap at 60", got.ConfidencePct)
	}
}

func TestForecast_usesMeanOfLastSixIntervals(t *testing.T) {
	today := testdata.Day(2026, 11, 1)
	// Eight starts: two old 40-day gaps, then six 28-day gaps ending Oct 19.
	starts := []time.Time{
		testdata.Day(2026, 3, 25),
		testdata.Day(2026, 5, 4), // 40-day gap, dropped when keeping last six
		testdata.Day(2026, 6, 1),
		testdata.Day(2026, 6, 29),
		testdata.Day(2026, 7, 27),
		testdata.Day(2026, 8, 24),
		testdata.Day(2026, 9, 21),
		testdata.Day(2026, 10, 19),
	}
	periods := make([]Period, len(starts))
	for i, start := range starts {
		periods[i] = Period{StartedOn: start, EndedOn: testdata.EndOf(start, 5)}
	}
	got := Forecast(Input{Today: today, Periods: periods})
	if got.CycleLength != 28 {
		t.Errorf("CycleLength = %d, want 28 (ignore the two old 40-day gaps)", got.CycleLength)
	}
}

func TestForecast_fertileWindowSpansOvulationMinus5ToPlus1(t *testing.T) {
	today := testdata.Day(2026, 10, 27)
	last := testdata.Day(2026, 10, 19)
	got := Forecast(Input{Today: today, Periods: periodsFromStarts(testdata.RegularStarts(last, 3, 28), 5, testdata.New())})
	if got.Phase != PhaseFertile {
		t.Errorf("Phase = %q, want fertile (day 9, ov-5)", got.Phase)
	}
	span := int(got.FertileEnd.Sub(got.FertileStart).Hours()/24) + 1
	if span != 7 {
		t.Errorf("fertile span = %d days, want 7", span)
	}
}

func TestForecast_bbtShiftMovesOvulation(t *testing.T) {
	f := testdata.New()
	today := testdata.Day(2026, 10, 27)
	last := testdata.Day(2026, 10, 19)
	obs := make([]Observation, 0, 9)
	for i := 0; i < 6; i++ {
		d := last.AddDate(0, 0, i)
		temp := 36.40
		obs = append(obs, Observation{Day: d, BBTCelsius: &temp})
	}
	for i := 6; i < 9; i++ {
		d := last.AddDate(0, 0, i)
		temp := 36.70
		obs = append(obs, Observation{Day: d, BBTCelsius: &temp})
	}
	_ = f
	got := Forecast(Input{
		Today:        today,
		Periods:      periodsFromStarts(testdata.RegularStarts(last, 3, 28), 5, testdata.New()),
		Observations: obs,
	})
	want := testdata.Day(2026, 10, 24)
	if !got.OvulationOn.Equal(want) {
		t.Errorf("OvulationOn = %s, want %s (day before BBT shift)", got.OvulationOn.Format(time.DateOnly), want.Format(time.DateOnly))
	}
}

func TestForecast_positiveOvulationTestWinsOverCalendar(t *testing.T) {
	last := testdata.Day(2026, 10, 19)
	today := testdata.Day(2026, 10, 28)
	testDay := testdata.Day(2026, 10, 26)
	got := Forecast(Input{
		Today:   today,
		Periods: periodsFromStarts(testdata.RegularStarts(last, 3, 28), 5, testdata.New()),
		Observations: []Observation{{
			Day:           testDay,
			OvulationTest: "positivo",
		}},
	})
	if !got.OvulationOn.Equal(testDay) {
		t.Errorf("OvulationOn = %s, want ovulation-test day %s", got.OvulationOn.Format(time.DateOnly), testDay.Format(time.DateOnly))
	}
}

func TestForecast_pregnancyMode_weekFromLMPAndDueIn280Days(t *testing.T) {
	lmp := testdata.Day(2026, 8, 9)
	today := testdata.Day(2026, 11, 1)
	got := Forecast(Input{Today: today, Mode: ModePregnancy, PregnancyLMP: &lmp})
	if got.Phase != PhasePregnancy {
		t.Errorf("Phase = %q, want pregnancy", got.Phase)
	}
	if got.PregnancyWeek != 13 {
		t.Errorf("PregnancyWeek = %d, want 13", got.PregnancyWeek)
	}
	wantDue := testdata.Day(2027, 5, 16)
	if !got.DueOn.Equal(wantDue) {
		t.Errorf("DueOn = %s, want %s", got.DueOn.Format(time.DateOnly), wantDue.Format(time.DateOnly))
	}
}

func TestForecast_typicalCycleOverrideWithSinglePeriod(t *testing.T) {
	last := testdata.Day(2026, 10, 2)
	today := testdata.Day(2026, 11, 1)
	got := Forecast(Input{
		Today:             today,
		Periods:           []Period{{StartedOn: last, EndedOn: testdata.EndOf(last, 4)}},
		TypicalCycleDays:  30,
		TypicalPeriodDays: 4,
	})
	if got.CycleLength != 30 {
		t.Errorf("CycleLength = %d, want typical 30 (one start, no interval yet)", got.CycleLength)
	}
	if got.PeriodLength != 4 {
		t.Errorf("PeriodLength = %d, want 4", got.PeriodLength)
	}
	if got.CycleDay != 31 {
		t.Errorf("CycleDay = %d, want 31", got.CycleDay)
	}
}

func TestConceptionChance_peakOnOvulationDay(t *testing.T) {
	ov := testdata.Day(2026, 11, 1)
	cases := []struct {
		offset int
		want   int
	}{
		{-6, 0},
		{-5, 10},
		{-4, 20},
		{-3, 33},
		{-2, 55},
		{-1, 75},
		{0, 87},
		{1, 25},
		{2, 0},
	}
	for _, tc := range cases {
		got := ConceptionChance(ov.AddDate(0, 0, tc.offset), ov)
		if got != tc.want {
			t.Errorf("offset %d: chance = %d, want %d", tc.offset, got, tc.want)
		}
	}
}

func TestForecast_openPeriodIsMenstruation(t *testing.T) {
	today := testdata.Day(2026, 11, 1)
	got := Forecast(Input{
		Today:   today,
		Periods: []Period{{StartedOn: today}},
	})
	if got.Phase != PhaseMenstruation {
		t.Errorf("Phase = %q, want menstruation", got.Phase)
	}
	if got.CycleDay != 1 {
		t.Errorf("CycleDay = %d, want 1", got.CycleDay)
	}
}

func TestStdDev_sampleOfKnownLengths(t *testing.T) {
	got := sampleStdDev([]float64{24, 35, 22, 40, 21, 38})
	if math.Abs(got-8.6023) > 0.01 {
		t.Errorf("stddev = %.4f, want ~8.6023", got)
	}
}
