package middleware

import (
	"net"
	"net/http"
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

func (rl *RateLimiter) getOrCreate(ip string) *bucket {
	rl.mu.RLock()
	b, ok := rl.buckets[ip]
	rl.mu.RUnlock()
	if ok {
		return b
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if b, ok = rl.buckets[ip]; ok {
		return b
	}

	b = &bucket{windowStart: time.Now()}
	rl.buckets[ip] = b
	return b
}

// realIP extracts the client IP, respecting X-Forwarded-For and X-Real-IP headers.
func realIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// take the first IP in the list
		if idx := len(xff); idx > 0 {
			for i, ch := range xff {
				if ch == ',' {
					xff = xff[:i]
					break
				}
			}
			return xff
		}
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}
