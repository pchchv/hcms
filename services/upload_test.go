package services

import "testing"

func TestNewUploader(t *testing.T) {
	dir := t.TempDir()
	u := NewUploader(dir)
	if u.Dir != dir {
		t.Errorf("expected dir %q, got %q", dir, u.Dir)
	}

	if u.MaxSize != 5*1024*1024 {
		t.Errorf("expected MaxSize=5MB, got %d", u.MaxSize)
	}
}
