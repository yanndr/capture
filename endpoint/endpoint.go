package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/yanndr/capture"
)

type Set struct {
	ExtractEndpoint    endpoint.Endpoint
	AddOverlayEndpoint endpoint.Endpoint
}

func New(svc capture.Service, logger log.Logger, duration metrics.Histogram) Set {
	extractEndpoint := MakeExtractEndpoint(svc)
	extractEndpoint = LoggingMiddleware(log.With(logger, "method", "Extract"))(extractEndpoint)
	extractEndpoint = InstrumentingMiddleware(duration.With("method", "Extract"))(extractEndpoint)

	overlayEndpoint := MakeAddOverlayEndpoint(svc)
	overlayEndpoint = LoggingMiddleware(log.With(logger, "method", "AddOverlay"))(extractEndpoint)
	overlayEndpoint = InstrumentingMiddleware(duration.With("method", "AddOverlay"))(extractEndpoint)

	return Set{
		ExtractEndpoint:    extractEndpoint,
		AddOverlayEndpoint: overlayEndpoint,
	}
}

func MakeExtractEndpoint(s capture.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(capture.ExtractRequest)
		return s.Extract(ctx, req)
	}
}

func MakeAddOverlayEndpoint(s capture.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(capture.OverlayImageRequest)
		return s.AddOverlay(req)
	}
}
