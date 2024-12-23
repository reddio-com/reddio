package utils

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"github.com/reddio-com/reddio/config"
)

var (
	GetReceiptRateLimiter *rate.Limiter
)

func IniLimiter() {
	GetReceiptRateLimiter = GenGetReceiptRateLimiter()
	logrus.Infof("GetReceipt Limit %v qps", config.GetGlobalConfig().RateLimitConfig.GetReceipt)
}

func GenGetReceiptRateLimiter() *rate.Limiter {
	qps := config.GetGlobalConfig().RateLimitConfig.GetReceipt
	if qps < 1 {
		return nil
	}
	limiter := rate.NewLimiter(rate.Limit(qps), 1)
	return limiter
}
