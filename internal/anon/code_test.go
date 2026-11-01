package anon

import (
	"strings"
	"testing"

	"github.com/puppe1990/aura-ciclos/internal/testdata"
)

func TestNewRecoveryCode_formatAndEmail(t *testing.T) {
	f := testdata.New()
	code, err := NewRecoveryCode(f)
	if err != nil {
		t.Fatal(err)
	}
	if !ValidCode(code) {
		t.Fatalf("generated code %q failed ValidCode", code)
	}
	email := EmailFor(code)
	if !strings.HasPrefix(email, "anon.") || !strings.HasSuffix(email, "@aura.local") {
		t.Errorf("email = %q", email)
	}
	if DisplayEmail(email) {
		t.Error("anon email must not be shown in the UI")
	}
}

func TestNormalizeCode_stripsDashesAndCase(t *testing.T) {
	got := NormalizeCode("ab12-cd34-ef56-gh78")
	if got != "AB12CD34EF56GH78" {
		t.Errorf("NormalizeCode = %q", got)
	}
	if EmailFor("ab12-cd34-ef56-gh78") != EmailFor("AB12CD34EF56GH78") {
		t.Error("email must be stable across formatting")
	}
}

func TestValidCode_rejectsShortAndPIILooking(t *testing.T) {
	if ValidCode("short") {
		t.Error("short code must fail")
	}
	if ValidCode("user@example.com") {
		t.Error("email must fail as recovery code")
	}
}

func TestIsAnonEmail(t *testing.T) {
	if !IsAnonEmail("anon.abcd1234abcd1234@aura.local") {
		t.Error("expected anon email")
	}
	if IsAnonEmail("sofia@example.com") {
		t.Error("real email is not anon")
	}
}
