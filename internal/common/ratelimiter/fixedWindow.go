package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowRateLImiter struct {
	mu          sync.Mutex
	limit       int
	windowSize  time.Duration
	windowStart time.Time
	count       int
}

var _ RateLimiter = (*FixedWindowRateLImiter)(nil)

func NewFixedWindowRateLimiter(limit int, wSize time.Duration) RateLimiter {
	return &FixedWindowRateLImiter{
		limit:       limit,
		windowSize:  wSize,
		windowStart: time.Now(),
		count:       0,
		mu:          sync.Mutex{},
	}
}

func (fw *FixedWindowRateLImiter) Allow() bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()
	if now.Sub(fw.windowStart) >= fw.windowSize {
		fw.windowStart = now
		fw.count = 0
	}

	if fw.count < fw.limit {
		fw.count++
		return true
	}
	// retry after logic
	// if now.Sub(fw.windowStart) < fw.windowSize {
	// 	retryAfter := fw.windowSize - now.Sub(fw.windowStart)
	// }
	return false
}
