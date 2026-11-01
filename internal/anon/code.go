package anon

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"unicode"

	"github.com/brianvoe/gofakeit/v7"
)

const (
	codeLen     = 16
	emailHost   = "@aura.local"
	emailPrefix = "anon."
	alphabet    = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
)

func NewRecoveryCode(f *gofakeit.Faker) (string, error) {
	if f == nil {
		return "", fmt.Errorf("recovery code: faker is required")
	}
	raw := make([]byte, codeLen)
	for i := range raw {
		raw[i] = alphabet[f.IntN(len(alphabet))]
	}
	return FormatCode(string(raw)), nil
}

func RandomRecoveryCode() (string, error) {
	raw := make([]byte, codeLen)
	max := big.NewInt(int64(len(alphabet)))
	for i := range raw {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("recovery code: %w", err)
		}
		raw[i] = alphabet[n.Int64()]
	}
	return FormatCode(string(raw)), nil
}

func FormatCode(raw string) string {
	n := NormalizeCode(raw)
	if len(n) != codeLen {
		return n
	}
	return n[0:4] + "-" + n[4:8] + "-" + n[8:12] + "-" + n[12:16]
}

func NormalizeCode(code string) string {
	var b strings.Builder
	for _, r := range strings.ToUpper(code) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func ValidCode(code string) bool {
	n := NormalizeCode(code)
	if len(n) != codeLen {
		return false
	}
	if strings.Contains(code, "@") {
		return false
	}
	for _, r := range n {
		if !strings.ContainsRune(alphabet+"IO01", r) {
			return false
		}
	}
	return true
}

func EmailFor(code string) string {
	return emailPrefix + strings.ToLower(NormalizeCode(code)) + emailHost
}

func IsAnonEmail(email string) bool {
	e := strings.ToLower(strings.TrimSpace(email))
	return strings.HasPrefix(e, emailPrefix) && strings.HasSuffix(e, emailHost)
}

func DisplayEmail(email string) bool {
	return !IsAnonEmail(email)
}
