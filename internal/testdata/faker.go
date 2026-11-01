// Package testdata builds fake domain values for tests and dev seeds.
//
// Every value derives from one seeded gofakeit Faker, so a failing test
// reproduces with identical data on every run.
package testdata

import (
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

// Seed keeps generated fixtures identical across runs.
const Seed = uint64(20261101)

func New() *gofakeit.Faker {
	return gofakeit.New(Seed)
}

func Day(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func EndOf(start time.Time, periodDays int) *time.Time {
	end := start.AddDate(0, 0, periodDays-1)
	return &end
}

func Flow(f *gofakeit.Faker) string {
	return f.RandomString([]string{"leve", "medio", "intenso"})
}

// RegularStarts returns n period start dates ending at lastStart, each cycleLen
// days apart, oldest first.
func RegularStarts(lastStart time.Time, n, cycleLen int) []time.Time {
	out := make([]time.Time, n)
	for i := 0; i < n; i++ {
		out[i] = lastStart.AddDate(0, 0, -cycleLen*(n-1-i))
	}
	return out
}

func Email(f *gofakeit.Faker) string {
	return strings.ToLower(f.FirstName()) + "@example.com"
}

func DisplayName(f *gofakeit.Faker) string {
	return f.FirstName()
}

func SymptomIDs(f *gofakeit.Faker, n int) []string {
	pool := []string{"colicas", "dor_cabeca", "acne", "inchaco", "ansiedade", "fadiga", "libido_alta", "nauseas"}
	out := make([]string, 0, n)
	seen := map[string]bool{}
	for len(out) < n {
		id := f.RandomString(pool)
		if seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	return out
}

func Celsius(f *gofakeit.Faker, min, max float64) float64 {
	return f.Float64Range(min, max)
}
