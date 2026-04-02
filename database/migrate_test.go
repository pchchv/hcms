package database

import (
	"testing"

	_ "modernc.org/sqlite"
)

func TestGeneratePassword_Length(t *testing.T) {
	for _, n := range []int{6, 12, 24} {
		if pw, err := generatePassword(n); err != nil {
			t.Fatalf("generatePassword(%d): %v", n, err)
		} else if len([]rune(pw)) != n {
			t.Errorf("generatePassword(%d) returned length %d", n, len([]rune(pw)))
		}
	}
}
