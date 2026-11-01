package symptoms

import (
	"testing"

	"github.com/puppe1990/aura-ciclos/internal/testdata"
)

func TestCatalog_hasAtLeast70SymptomsWithUniqueIDs(t *testing.T) {
	all := All()
	if len(all) < 70 {
		t.Fatalf("len(All) = %d, want >= 70", len(all))
	}
	seen := map[string]bool{}
	for _, s := range all {
		if s.ID == "" || s.Label == "" || s.Category == "" {
			t.Fatalf("incomplete symptom: %+v", s)
		}
		if seen[s.ID] {
			t.Fatalf("duplicate id %q", s.ID)
		}
		seen[s.ID] = true
	}
}

func TestCatalog_coversRequiredCategories(t *testing.T) {
	need := []string{CategoryFluxo, CategoryFisicos, CategoryHumor, CategoryEstilo, CategorySexual}
	counts := map[string]int{}
	for _, s := range All() {
		counts[s.Category]++
	}
	for _, cat := range need {
		if counts[cat] == 0 {
			t.Errorf("missing category %q", cat)
		}
	}
}

func TestByID_findsKnownAndRejectsUnknown(t *testing.T) {
	got, ok := ByID("colicas")
	if !ok || got.Label != "Cólicas" {
		t.Fatalf("ByID(colicas) = %+v ok=%v", got, ok)
	}
	if _, ok := ByID("nao-existe"); ok {
		t.Fatal("unknown id must miss")
	}
}

func TestFilter_keepsFakerPicks(t *testing.T) {
	f := testdata.New()
	ids := testdata.SymptomIDs(f, 4)
	got := Filter(ids)
	if len(got) != 4 {
		t.Fatalf("Filter(%v) len = %d", ids, len(got))
	}
}

func TestGrouped_ordersCategories(t *testing.T) {
	groups := Grouped()
	if len(groups) < 4 {
		t.Fatalf("Grouped len = %d", len(groups))
	}
	if groups[0].Category != CategoryFluxo {
		t.Errorf("first group = %q, want fluxo", groups[0].Category)
	}
}
