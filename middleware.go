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

func (mw loggingMiddleware) Extract(ctx context.Context, request ExtractRequest) (ExtractResponse, error) {
	defer func() {
		mw.logger.Log("method", "Extract", "ExtractRequest", request.Name)
	}()
	return mw.next.Extract(ctx, request)
}

func (mw loggingMiddleware) AddOverlay(request OverlayImageRequest) (ExtractResponse, error) {
	defer func() {
		mw.logger.Log("method", "AddOverlay")
	}()
	return mw.next.AddOverlay(request)
}

type instrumentingMiddleware struct {
	extracts metrics.Counter
	overlays metrics.Counter
	next     Service
}

// InstrumentingMiddleware returns a service middleware that instruments
// the number of integers summed and characters concatenated over the lifetime of
// the service.
func InstrumentingMiddleware(extracts, overlays metrics.Counter) Middleware {
	return func(next Service) Service {
		return instrumentingMiddleware{
			extracts: extracts,
			overlays: overlays,
			next:     next,
		}
	}
}

func (mw instrumentingMiddleware) Extract(ctx context.Context, request ExtractRequest) (ExtractResponse, error) {
	mw.extracts.Add(float64(1))
	return mw.next.Extract(ctx, request)
}

func (mw instrumentingMiddleware) AddOverlay(request OverlayImageRequest) (ExtractResponse, error) {
	mw.overlays.Add(float64(1))
	return mw.next.AddOverlay(request)
}
