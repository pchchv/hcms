package middleware

import (
	"sync"
	"sync/atomic"
	"time"
)

// bucket holds the request count and the start of the current window for one IP.
// mu protects only windowStart resets (rare); count is incremented atomically.
type bucket struct {
	mu          sync.Mutex
	windowStart time.Time
	count       atomic.Int64
}
