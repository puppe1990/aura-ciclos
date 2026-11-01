package cycle

import (
	"math"
	"sort"
	"time"
)

const (
	DefaultCycleDays  = 28
	DefaultPeriodDays = 5
	LutealDays        = 14
	IrregularStdDev   = 4.0
	maxIntervals      = 6
	bbtRiseC          = 0.2
	pregnancyLength   = 280
)

type Mode string

const (
	ModeCycle     Mode = "cycle"
	ModeTrying    Mode = "trying"
	ModePregnancy Mode = "pregnancy"
)

type Phase string

const (
	PhaseMenstruation Phase = "menstruacao"
	PhaseFollicular   Phase = "folicular"
	PhaseFertile      Phase = "fertil"
	PhaseOvulation    Phase = "ovulacao"
	PhaseLuteal       Phase = "lutea"
	PhasePregnancy    Phase = "gravidez"
)

type Period struct {
	StartedOn time.Time
	EndedOn   *time.Time
	Flow      string
}

type Observation struct {
	Day           time.Time
	BBTCelsius    *float64
	OvulationTest string
}

type Input struct {
	Today             time.Time
	Periods           []Period
	Observations      []Observation
	Mode              Mode
	PregnancyLMP      *time.Time
	TypicalCycleDays  int
	TypicalPeriodDays int
}

type Result struct {
	CycleDay         int
	CycleLength      int
	PeriodLength     int
	Phase            Phase
	OvulationOn      time.Time
	FertileStart     time.Time
	FertileEnd       time.Time
	NextPeriodOn     time.Time
	DaysUntilPeriod  int
	Irregular        bool
	StdDev           float64
	Confidence       int
	ConfidencePct    int
	ConceptionChance int
	PregnancyWeek    int
	DueOn            time.Time
	DaysUntilDue     int
	Lengths          []int
	FollicularDays   int
	LutealDays       int
}

func Forecast(in Input) Result {
	today := truncate(in.Today)
	if in.Mode == ModePregnancy && in.PregnancyLMP != nil {
		return pregnancyForecast(today, truncate(*in.PregnancyLMP))
	}

	cycleLen := pickDefault(in.TypicalCycleDays, DefaultCycleDays)
	periodLen := pickDefault(in.TypicalPeriodDays, DefaultPeriodDays)
	starts := sortedStarts(in.Periods)
	lengths := startIntervals(starts)
	if n := len(lengths); n > maxIntervals {
		lengths = lengths[n-maxIntervals:]
	}

	stddev := 0.0
	irregular := false
	if len(lengths) > 0 {
		cycleLen = int(math.Round(mean(floats(lengths))))
		if cycleLen < 15 {
			cycleLen = 15
		}
		stddev = sampleStdDev(floats(lengths))
		irregular = stddev >= IrregularStdDev
	}
	if recorded := meanPeriodDays(in.Periods); recorded > 0 {
		periodLen = recorded
	}

	anchor := today
	if len(starts) > 0 {
		anchor = starts[len(starts)-1]
	}
	cycleDay := daysBetween(anchor, today) + 1
	if cycleDay < 1 {
		cycleDay = 1
	}

	ovulationOn := calendarOvulation(anchor, cycleLen)
	if shifted, ok := bbtOvulation(in.Observations); ok {
		ovulationOn = shifted
	}
	if tested, ok := testOvulation(in.Observations); ok {
		ovulationOn = tested
	}

	fertileStart := ovulationOn.AddDate(0, 0, -5)
	fertileEnd := ovulationOn.AddDate(0, 0, 1)
	nextPeriod := anchor.AddDate(0, 0, cycleLen)
	until := daysBetween(today, nextPeriod)
	if until < 0 {
		until = 0
	}

	conf, pct := confidence(len(lengths), irregular)
	phase := classifyPhase(today, cycleDay, periodLen, ovulationOn, fertileStart, fertileEnd)
	return Result{
		CycleDay:         cycleDay,
		CycleLength:      cycleLen,
		PeriodLength:     periodLen,
		Phase:            phase,
		OvulationOn:      ovulationOn,
		FertileStart:     fertileStart,
		FertileEnd:       fertileEnd,
		NextPeriodOn:     nextPeriod,
		DaysUntilPeriod:  until,
		Irregular:        irregular,
		StdDev:           stddev,
		Confidence:       conf,
		ConfidencePct:    pct,
		ConceptionChance: ConceptionChance(today, ovulationOn),
		Lengths:          lengths,
		FollicularDays:   cycleLen - LutealDays,
		LutealDays:       LutealDays,
	}
}

func pregnancyForecast(today, lmp time.Time) Result {
	week := daysBetween(lmp, today)/7 + 1
	if week < 1 {
		week = 1
	}
	due := lmp.AddDate(0, 0, pregnancyLength)
	until := daysBetween(today, due)
	if until < 0 {
		until = 0
	}
	return Result{
		Phase:         PhasePregnancy,
		PregnancyWeek: week,
		DueOn:         due,
		DaysUntilDue:  until,
		CycleLength:   DefaultCycleDays,
		PeriodLength:  DefaultPeriodDays,
		Confidence:    4,
		ConfidencePct: 98,
	}
}

func ConceptionChance(today, ovulation time.Time) int {
	offset := daysBetween(ovulation, today)
	switch offset {
	case -5:
		return 10
	case -4:
		return 20
	case -3:
		return 33
	case -2:
		return 55
	case -1:
		return 75
	case 0:
		return 87
	case 1:
		return 25
	default:
		return 0
	}
}

func classifyPhase(today time.Time, cycleDay, periodLen int, ovulation, fertileStart, fertileEnd time.Time) Phase {
	if !today.Before(truncate(today)) && cycleDay <= periodLen {
		return PhaseMenstruation
	}
	if sameDay(today, ovulation) {
		return PhaseOvulation
	}
	if !today.Before(fertileStart) && !today.After(fertileEnd) {
		return PhaseFertile
	}
	if today.Before(ovulation) {
		return PhaseFollicular
	}
	return PhaseLuteal
}

func calendarOvulation(anchor time.Time, cycleLen int) time.Time {
	ovDay := cycleLen - LutealDays
	if ovDay < 1 {
		ovDay = 1
	}
	return anchor.AddDate(0, 0, ovDay-1)
}

func bbtOvulation(obs []Observation) (time.Time, bool) {
	type sample struct {
		day  time.Time
		temp float64
	}
	var samples []sample
	for _, o := range obs {
		if o.BBTCelsius == nil {
			continue
		}
		samples = append(samples, sample{day: truncate(o.Day), temp: *o.BBTCelsius})
	}
	sort.Slice(samples, func(i, j int) bool { return samples[i].day.Before(samples[j].day) })
	if len(samples) < 9 {
		return time.Time{}, false
	}
	for i := 6; i+2 < len(samples); i++ {
		cover := 0.0
		for j := i - 6; j < i; j++ {
			cover += samples[j].temp
		}
		cover /= 6
		if samples[i].temp < cover+bbtRiseC || samples[i+1].temp < cover+bbtRiseC || samples[i+2].temp < cover+bbtRiseC {
			continue
		}
		return samples[i].day.AddDate(0, 0, -1), true
	}
	return time.Time{}, false
}

func testOvulation(obs []Observation) (time.Time, bool) {
	var latest time.Time
	found := false
	for _, o := range obs {
		if o.OvulationTest != "positivo" {
			continue
		}
		day := truncate(o.Day)
		if !found || day.After(latest) {
			latest = day
			found = true
		}
	}
	return latest, found
}

func confidence(intervals int, irregular bool) (int, int) {
	switch {
	case intervals == 0:
		return 1, 25
	case intervals == 1:
		return 2, 50
	case intervals <= 3:
		if irregular {
			return 2, 45
		}
		return 3, 75
	case intervals <= 5:
		if irregular {
			return 3, 60
		}
		return 4, 90
	default:
		if irregular {
			return 3, 60
		}
		return 4, 98
	}
}

func sortedStarts(periods []Period) []time.Time {
	starts := make([]time.Time, 0, len(periods))
	for _, p := range periods {
		starts = append(starts, truncate(p.StartedOn))
	}
	sort.Slice(starts, func(i, j int) bool { return starts[i].Before(starts[j]) })
	return starts
}

func startIntervals(starts []time.Time) []int {
	if len(starts) < 2 {
		return nil
	}
	out := make([]int, 0, len(starts)-1)
	for i := 1; i < len(starts); i++ {
		out = append(out, daysBetween(starts[i-1], starts[i]))
	}
	return out
}

func meanPeriodDays(periods []Period) int {
	sum, n := 0, 0
	for _, p := range periods {
		if p.EndedOn == nil {
			continue
		}
		d := daysBetween(truncate(p.StartedOn), truncate(*p.EndedOn)) + 1
		if d > 0 {
			sum += d
			n++
		}
	}
	if n == 0 {
		return 0
	}
	return int(math.Round(float64(sum) / float64(n)))
}

func sampleStdDev(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	m := mean(values)
	sum := 0.0
	for _, v := range values {
		d := v - m
		sum += d * d
	}
	return math.Sqrt(sum / float64(len(values)-1))
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func floats(values []int) []float64 {
	out := make([]float64, len(values))
	for i, v := range values {
		out[i] = float64(v)
	}
	return out
}

func daysBetween(from, to time.Time) int {
	a := truncate(from)
	b := truncate(to)
	return int(b.Sub(a).Hours() / 24)
}

func truncate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func sameDay(a, b time.Time) bool {
	return truncate(a).Equal(truncate(b))
}

func pickDefault(v, fallback int) int {
	if v > 0 {
		return v
	}
	return fallback
}
