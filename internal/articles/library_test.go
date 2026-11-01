package articles

import "testing"

func TestLibrary_hasReviewedPortugueseArticles(t *testing.T) {
	all := All()
	if len(all) < 6 {
		t.Fatalf("len(All) = %d, want >= 6", len(all))
	}
	seen := map[string]bool{}
	for _, a := range all {
		if a.Slug == "" || a.Title == "" || a.Summary == "" || len(a.Body) < 200 {
			t.Fatalf("thin article: %+v", a)
		}
		if a.ReviewedBy == "" {
			t.Fatalf("missing reviewer on %s", a.Slug)
		}
		if seen[a.Slug] {
			t.Fatalf("duplicate slug %s", a.Slug)
		}
		seen[a.Slug] = true
	}
}

func TestBySlug_findsAndMisses(t *testing.T) {
	first := All()[0]
	got, ok := BySlug(first.Slug)
	if !ok || got.Title != first.Title {
		t.Fatalf("BySlug(%s) = %+v ok=%v", first.Slug, got, ok)
	}
	if _, ok := BySlug("nao-existe"); ok {
		t.Fatal("unknown slug must miss")
	}
}
