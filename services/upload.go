package services

// Uploader handles file uploads.
type Uploader struct {
	Dir     string
	MaxSize int64 // bytes, default 5 MB
}

// NewUploader creates a new Uploader for the
// given directory with a 5 MB max size.
func NewUploader(dir string) *Uploader {
	return &Uploader{
		Dir:     dir,
		MaxSize: 5 * 1024 * 1024, // 5 MB
	}
}
