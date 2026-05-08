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

// RateLimiter implements IP-based rate limiting.
type RateLimiter struct {
	mu      sync.RWMutex
	buckets map[string]*bucket
}

// NewRateLimiter creates a new RateLimiter.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{buckets: make(map[string]*bucket)}
}
