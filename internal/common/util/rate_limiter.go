package util

import (
	"sync"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/ratelimiter"
)

type RateLimiterWrapper struct {
	limitermu      sync.Mutex
	limiterEntries map[string]ratelimiter.RateLimiter
}

func NewRateLimiterWrapper() *RateLimiterWrapper {
	return &RateLimiterWrapper{
		limitermu:      sync.Mutex{},
		limiterEntries: make(map[string]ratelimiter.RateLimiter),
	}
}

func (rw *RateLimiterWrapper) Allow(identifer string) bool {
	rw.limitermu.Lock()
	defer rw.limitermu.Unlock()

	limiter, exist := rw.limiterEntries[identifer]
	if !exist {
		limiter = ratelimiter.RateLimiterFactory()
		rw.limiterEntries[identifer] = limiter
	}
	return limiter.Allow()
}

// -----------------------------------Advance version of above wrapper
// Wrapper that provides per-identifier limiting using a factory
type IdentifierRateLimiter struct {
	factory  func() ratelimiter.RateLimiter // Creates new limiter for each identifier
	limiters sync.Map                       // map[identifier]RateLimiter
}

func NewIdentifierRateLimiter(factory func() ratelimiter.RateLimiter) *IdentifierRateLimiter {
	return &IdentifierRateLimiter{
		factory: factory,
	}
}

// This is a proxy method that calls the Allow method of actual limiter
func (l *IdentifierRateLimiter) Allow(identifier string) bool {
	// Get or create limiter for this identifier
	limiter, loaded := l.limiters.LoadOrStore(identifier, l.factory())
	if loaded {
		// log here that existing limiter is being used for this identifier
	}

	// Type assert to RateLimiter
	rl := limiter.(ratelimiter.RateLimiter)
	return rl.Allow()
}
