package cycle

import (
	"testing"

	"github.com/puppe1990/aura-ciclos/internal/testdata"
)

func TestPhysiology_ovulationPeaksLH(t *testing.T) {
	got := Physiology(PhaseOvulation)
	if got.Hormones != "LH & Estrogênio" {
		t.Errorf("Hormones = %q", got.Hormones)
	}
	if got.Vitality != "Energia máxima" {
		t.Errorf("Vitality = %q", got.Vitality)
	}
	if got.Badge != "Fase ovulatória • Janela fértil" {
		t.Errorf("Badge = %q", got.Badge)
	}
}

func TestPhysiology_lutealUsesProgesterone(t *testing.T) {
	got := Physiology(PhaseLuteal)
	if got.Hormones != "Progesterona" {
		t.Errorf("Hormones = %q", got.Hormones)
	}
	if got.Badge != "Fase lútea" {
		t.Errorf("Badge = %q", got.Badge)
	}
}

func TestConfidenceWords(t *testing.T) {
	if ConfidenceWords(4) != "Excelente (Nível 4 de 4)" {
		t.Errorf("level 4 = %q", ConfidenceWords(4))
	}
	if ConfidenceWords(1) != "Baixa (Nível 1 de 4)" {
		t.Errorf("level 1 = %q", ConfidenceWords(1))
	}
}

func TestChanceWords(t *testing.T) {
	if ChanceWords(87) != "Chance Alta" {
		t.Errorf("87 = %q", ChanceWords(87))
	}
	if ChanceWords(0) != "" {
		t.Errorf("0 = %q", ChanceWords(0))
	}
}

func TestPregnancySize_week14Lemon(t *testing.T) {
	if PregnancySize(14) != "Limão" {
		t.Errorf("week 14 = %q", PregnancySize(14))
	}
}

func TestFertileIndex_fourthOfSix(t *testing.T) {
	ov := testdata.Day(2026, 11, 1)
	start := ov.AddDate(0, 0, -5)
	end := ov.AddDate(0, 0, 1)
	n, total := FertileIndex(ov, start, end)
	if n != 6 || total != 7 {
		t.Errorf("index = %d/%d, want 6/7 (day 0 of -5..+1 is the 6th day)", n, total)
	}
}
