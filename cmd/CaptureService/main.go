package main

import (
	"net"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yanndr/capture"
	"github.com/yanndr/capture/pb"
	"github.com/yanndr/capture/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	gRPCPort  = ":50051"
	debugAddr = ":8080"
)

func main() {

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var extracts metrics.Counter
	{
		// Business-level metrics.
		extracts = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "capture",
			Subsystem: "captureSvx",
			Name:      "extracts_summed",
			Help:      "Total count of extracts done via the Extract method.",
		}, []string{})
	}

	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	var (
		service    = capture.NewService(logger, extracts)
		grpcServer = transport.NewGRPCServer(service)
	)

	debugListener, err := net.Listen("tcp", debugAddr)
	if err != nil {
		logger.Log("transport", "debug/HTTP", "during", "Listen", "err", err)
		os.Exit(1)
	}

	go func() {
		if err := http.Serve(debugListener, http.DefaultServeMux); err != nil {
			logger.Log("transport", "HTTP", "during", "serve", "err", err)
		}
	}()

	grpcListener, err := net.Listen("tcp", gRPCPort)
	if err != nil {
		logger.Log("transport", "gRPC", "during", "listen", "err", err)
		os.Exit(1)
	}

	creds, err := credentials.NewServerTLSFromFile("../../../cert.pem", "../../../key.pem")
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
