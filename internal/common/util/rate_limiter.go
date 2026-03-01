package util

import (
	"sync"
	"time"
)

type rateLimitEntry struct {
	count     int
	resetTime time.Time
}

// FixedWindowRateLimiter is an in-memory fixed-window rate limiter.
type FixedWindowRateLimiter struct {
	maxRequests int
	window      time.Duration
	entries     map[string]*rateLimitEntry
	mu          sync.Mutex
}

func NewFixedWindowRateLimiter(maxRequests int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		maxRequests: maxRequests,
		window:      window,
		entries:     make(map[string]*rateLimitEntry),
	}
}

// Allow reports whether a request for identifier should be allowed.
// When blocked, retryAfter is the number of seconds until reset.
func (r *FixedWindowRateLimiter) Allow(identifier string, now time.Time) (allowed bool, retryAfter int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry, exists := r.entries[identifier]
	if !exists || now.After(entry.resetTime) {
		r.entries[identifier] = &rateLimitEntry{
			count:     1,
			resetTime: now.Add(r.window),
		}
		return true, 0
	}

	if entry.count >= r.maxRequests {
		retry := int(entry.resetTime.Sub(now).Seconds())
		if retry < 1 {
			retry = 1
		}
		return false, retry
	}

	entry.count++
	return true, 0
}
