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
		service    = capture.VideoCaptureService{}
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
		// log.Fatal(err)
	}
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	server := grpc.NewServer(opts...)
	pb.RegisterVideoCaptureServer(server, grpcServer)

	if err := server.Serve(grpcListener); err != nil {
		logger.Log("transport", "gRPC", "during", "listen", "err", err)
		os.Exit(1)
	}

	// lis, err := net.Listen("tcp", port)
	// if err != nil {
	// 	// log.Fatalf("failed to listen: %v", err)
	// }

	// creds, err := credentials.NewServerTLSFromFile("../cert.pem", "../key.pem")
	// if err != nil {
	// 	// log.Fatal(err)
	// }
	// opts := []grpc.ServerOption{grpc.Creds(creds)}

	// s := grpc.NewServer(opts...)
	// pb.RegisterVideoCaptureServer(s, &capture.VideoCaptureService{})
	// // Register reflection service on gRPC server.
	// reflection.Register(s)
	// if err := s.Serve(lis); err != nil {
	// 	// log.Fatalf("failed to serve: %v", err)
	// }
}
