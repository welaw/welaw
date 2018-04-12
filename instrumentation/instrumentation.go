package instrumentation

import (
	"github.com/go-kit/kit/metrics"
	"github.com/welaw/welaw/services"
)

type instrumentatingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           services.Service
}

func NewInstrumentatingMiddleware(rc metrics.Counter, rl metrics.Histogram, s services.Service) services.Service {
	return &instrumentatingMiddleware{
		requestCount:   rc,
		requestLatency: rl,
		next:           s,
	}
}
