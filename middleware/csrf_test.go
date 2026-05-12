package middleware

import "testing"

func TestGenerateToken_Deterministic(t *testing.T) {
	t1 := GenerateToken("session-abc")
	t2 := GenerateToken("session-abc")
	if t1 != t2 {
		t.Error("GenerateToken should be deterministic for the same session ID")
	}
}

func TestGenerateToken_DifferentInputs(t *testing.T) {
	t1 := GenerateToken("session-aaa")
	t2 := GenerateToken("session-bbb")
	if t1 == t2 {
		t.Error("different session IDs should produce different tokens")
	}
}
