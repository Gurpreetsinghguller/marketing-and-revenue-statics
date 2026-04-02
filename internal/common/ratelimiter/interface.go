package ratelimiter

import (
	"time"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/config"
)

type RateLimiter interface {
	Allow() bool
}

// we can pass the whole config or the ratelimiter config struct in this Factory
func RateLimiterFactory() RateLimiter {
	rlcfg := config.GetConfig().RateLimit
	switch rlcfg.RateLimiterType {
	case "fixedwindow":
		return NewFixedWindowRateLimiter(rlcfg.MaxRequests, time.Duration(rlcfg.WindowSeconds)*time.Second)
	default:
		return nil
	}
}
