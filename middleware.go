package capture

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"golang.org/x/net/context"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) Extract(ctx context.Context, request ExtractRequest) ExtractResponse {
	defer func() {
		mw.logger.Log("method", "Extract", "ExtractRequest", request)
	}()
	return mw.next.Extract(ctx, request)
}

type instrumentingMiddleware struct {
	extracts metrics.Counter
	next     Service
}

// InstrumentingMiddleware returns a service middleware that instruments
// the number of integers summed and characters concatenated over the lifetime of
// the service.
func InstrumentingMiddleware(extracts metrics.Counter) Middleware {
	return func(next Service) Service {
		return instrumentingMiddleware{
			extracts: extracts,
			next:     next,
		}
	}
}

func (mw instrumentingMiddleware) Extract(ctx context.Context, request ExtractRequest) ExtractResponse {
	mw.extracts.Add(float64(1))
	return mw.next.Extract(ctx, request)
}
