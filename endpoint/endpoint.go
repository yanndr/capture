package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/yanndr/capture"
)

type Set struct {
	ExtractEndpoint endpoint.Endpoint
}

func New(svc capture.VideoCaptureService, logger log.Logger) Set {
	extractEndpoint := MakeExtractEndpoint(svc)

	return Set{
		ExtractEndpoint: extractEndpoint,
	}
}

func MakeExtractEndpoint(s capture.VideoCaptureService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(capture.ExtractRequest)
		return s.Extract(ctx, req), nil
	}
}
