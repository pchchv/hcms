package database

import (
	"database/sql"
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

func TestGeneratePassword_UsesCharset(t *testing.T) {
	pw, err := generatePassword(100)
	if err != nil {
		t.Fatalf("generatePassword: %v", err)
	}

	for _, ch := range pw {
		if !containsRune(passwordCharset, ch) {
			t.Errorf("password contains invalid character: %q", ch)
		}
	}
}

func containsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	t.Cleanup(func() { db.Close() })
	return db
}
