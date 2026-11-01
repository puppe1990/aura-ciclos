package handlers

import (
	"testing"

	"github.com/puppe1990/aura-ciclos/internal/models"
	"github.com/puppe1990/aura-ciclos/internal/testdata"
)

func TestBBTSparkline_threePoints(t *testing.T) {
	a, b, c := 36.40, 36.45, 36.65
	got := bbtSparkline([]models.DailyLog{
		{LoggedOn: testdata.Day(2026, 10, 30), BBTCelsius: &a},
		{LoggedOn: testdata.Day(2026, 10, 31), BBTCelsius: &b},
		{LoggedOn: testdata.Day(2026, 11, 1), BBTCelsius: &c},
	})
	if got.Empty {
		t.Fatal("sparkline should not be empty")
	}
	if got.Latest != "36,65 °C" {
		t.Errorf("Latest = %q", got.Latest)
	}
	if got.Path == "" || got.Area == "" {
		t.Errorf("missing path")
	}
	if got.Delta != "(+0.2°)" {
		t.Errorf("Delta = %q", got.Delta)
	}
}

func TestAvatarLetter(t *testing.T) {
	if avatarLetter("Sofia") != "S" {
		t.Errorf("Sofia = %q", avatarLetter("Sofia"))
	}
	if avatarLetter("você") != "A" {
		t.Errorf("anon = %q", avatarLetter("você"))
	}
}

func TestWaterLiters(t *testing.T) {
	if waterLiters(8) != "2,0" {
		t.Errorf("8 glasses = %q", waterLiters(8))
	}
}
