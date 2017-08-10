package transport

import (
	"context"

	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/yanndr/capture"
	"github.com/yanndr/capture/endpoint"
	"github.com/yanndr/capture/pb"

	oldcontext "golang.org/x/net/context"
)

type grpcServer struct {
	extract grpctransport.Handler
}

func NewGRPCServer(endpoints endpoint.Set, logger log.Logger) pb.VideoCaptureServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}

	return &grpcServer{
		extract: grpctransport.NewServer(
			endpoints.ExtractEndpoint,
			decodeGPRCExtractRequest,
			decodeGPRCExtractRequest,
		),
	}
}

func decodeGPRCExtractRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.VideoCaptureRequest)
	return capture.ExtractRequest{Path: req.Path}, nil
}

func (s *grpcServer) ExtractImage(ctx oldcontext.Context, req *pb.VideoCaptureRequest) (*pb.VideoCaptureReply, error) {
	_, rep, err := s.extract.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.VideoCaptureReply), nil
}
