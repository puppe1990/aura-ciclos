package maya

import (
	"strings"
	"testing"

	"github.com/puppe1990/aura-ciclos/internal/cycle"
	"github.com/puppe1990/aura-ciclos/internal/testdata"
)

func TestReply_emptyFallsBackToInvite(t *testing.T) {
	got := Reply("", Context{})
	if !strings.Contains(strings.ToLower(got), "pergunta") {
		t.Errorf("empty reply = %q", got)
	}
}

func TestReply_colicasMentionsHeatAndRest(t *testing.T) {
	got := Reply("estou com muita cólica hoje", Context{Phase: cycle.PhaseMenstruation})
	low := strings.ToLower(got)
	if !strings.Contains(low, "cólica") && !strings.Contains(low, "colica") {
		t.Errorf("missing cramp mention: %s", got)
	}
	if !strings.Contains(low, "calor") && !strings.Contains(low, "descanso") {
		t.Errorf("missing comfort advice: %s", got)
	}
}

func TestReply_atrasoMentionsStressAndPregnancy(t *testing.T) {
	got := Reply("minha menstruação atrasou 5 dias", Context{})
	low := strings.ToLower(got)
	if !strings.Contains(low, "atras") {
		t.Errorf("missing late period copy: %s", got)
	}
	if !strings.Contains(low, "estresse") && !strings.Contains(low, "gravidez") {
		t.Errorf("missing causes: %s", got)
	}
}

func TestReply_ovulationExplainsWindow(t *testing.T) {
	got := Reply("quando eu ovulo?", Context{
		Phase:     cycle.PhaseFertile,
		CycleDay:  12,
		FertileOn: true,
		ChancePct: 55,
	})
	if !strings.Contains(strings.ToLower(got), "fértil") && !strings.Contains(strings.ToLower(got), "fertil") {
		t.Errorf("missing fertile window: %s", got)
	}
}

func TestReply_tpmInLuteal(t *testing.T) {
	got := Reply("tpm está pesada", Context{Phase: cycle.PhaseLuteal})
	if !strings.Contains(strings.ToLower(got), "tpm") {
		t.Errorf("missing TPM: %s", got)
	}
}

func TestReply_unknownStillHelpful(t *testing.T) {
	q := testdata.New().Sentence(6)
	got := Reply(q, Context{})
	if len(got) < 40 {
		t.Errorf("fallback too short: %q", got)
	}
}

func TestReply_disclaimerAlwaysPresent(t *testing.T) {
	got := Reply("dor de cabeça", Context{})
	if !strings.Contains(got, "não substitui") {
		t.Errorf("missing clinical disclaimer: %s", got)
	}
}
