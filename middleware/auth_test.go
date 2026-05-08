package middleware

import (
	"context"
	"testing"
)

func TestGetSession_NoValue(t *testing.T) {
	s := GetSession(context.Background())
	if s != nil {
		t.Error("expected nil session from empty context")
	}
}
