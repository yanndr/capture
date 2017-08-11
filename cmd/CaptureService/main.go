package main

import (
	"net"
	"os"

	"github.com/yanndr/capture"

	"github.com/yanndr/capture/endpoint"
	"github.com/yanndr/capture/pb"
	"github.com/yanndr/capture/transport"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port = ":50051"
)

func main() {

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var (
		service    = capture.NewService(logger)
		endpoints  = endpoint.New(service, logger)
		grpcServer = transport.NewGRPCServer(endpoints, logger)
	)

	grpcListener, err := net.Listen("tcp", port)
	if err != nil {
		logger.Log("transport", "gRPC", "during", "listen", "err", err)
		os.Exit(1)
	}

	creds, err := credentials.NewServerTLSFromFile("../cert.pem", "../key.pem")
	if err != nil {
		logger.Log("transport", "gRPC", "during", "credential", "err", err)
		os.Exit(1)
	}
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	server := grpc.NewServer(opts...)
	pb.RegisterVideoCaptureServer(server, grpcServer)

	if err := server.Serve(grpcListener); err != nil {
		logger.Log("transport", "gRPC", "during", "serve", "err", err)
		os.Exit(1)
	}
}
