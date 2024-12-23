package utils

import (
	"time"

	"golang.org/x/time/rate"

	"github.com/reddio-com/reddio/config"
)

var (
	GetReceiptRateLimiter *rate.Limiter
)

func IniLimiter() {
	GetReceiptRateLimiter = GenGetReceiptRateLimiter()
}

func GenGetReceiptRateLimiter() *rate.Limiter {
	qps := config.GetGlobalConfig().RateLimitConfig.GetReceipt
	limiter := rate.NewLimiter(rate.Every(10*time.Millisecond), int(qps/100))
	return limiter
}
