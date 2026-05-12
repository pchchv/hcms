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

func TestVerify_Valid(t *testing.T) {
	sessionID := "test-session-id"
	token := GenerateToken(sessionID)
	if !Verify(sessionID, token) {
		t.Error("Verify should return true for valid token")
	}
}

func TestVerify_Invalid(t *testing.T) {
	if Verify("session-id", "wrongtoken") {
		t.Error("Verify should return false for wrong token")
	}
}

func TestVerify_EmptyToken(t *testing.T) {
	if Verify("session-id", "") {
		t.Error("Verify should return false for empty token")
	}
}

func TestVerify_WrongSession(t *testing.T) {
	token := GenerateToken("correct-session")
	if Verify("wrong-session", token) {
		t.Error("Verify should return false for wrong session")
	}
}
