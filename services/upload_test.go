package services

import (
	"bytes"
	"mime/multipart"
	"net/textproto"
	"testing"
)

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

func makeFakeMultipartFile(t *testing.T, content []byte, filename, contentType string) *multipart.FileHeader {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="image"; filename="`+filename+`"`)
	h.Set("Content-Type", contentType)
	part, err := w.CreatePart(h)
	if err != nil {
		t.Fatalf("create part: %v", err)
	}

	part.Write(content)
	w.Close()

	reader := multipart.NewReader(&buf, w.Boundary())
	form, err := reader.ReadForm(10 << 20)
	if err != nil {
		t.Fatalf("read form: %v", err)
	}

	files := form.File["image"]
	if len(files) == 0 {
		t.Fatal("no files in form")
	}

	return files[0]
}
