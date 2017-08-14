package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/yanndr/capture"
	"github.com/yanndr/capture/pb"
)

type Set struct {
	ExtractEndpoint endpoint.Endpoint
}

func New(svc capture.Service, logger log.Logger, duration metrics.Histogram) Set {
	extractEndpoint := MakeExtractEndpoint(svc)
	extractEndpoint = LoggingMiddleware(log.With(logger, "method", "Extract"))(extractEndpoint)
	extractEndpoint = InstrumentingMiddleware(duration.With("method", "Extract"))(extractEndpoint)

	return Set{
		ExtractEndpoint: extractEndpoint,
	}
}

func MakeExtractEndpoint(s capture.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(capture.ExtractRequest)
		return s.Extract(ctx, req), nil
	}
}

func (s *grpcServer) ExtractImage(stream pb.VideoCapture_ExtractImageServer) error
