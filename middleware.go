package capture

import (
	"github.com/go-kit/kit/log"
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
