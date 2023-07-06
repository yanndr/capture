package transport

import (
	"context"
	kitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/yanndr/capture"
	"github.com/yanndr/capture/endpoint"
	"github.com/yanndr/capture/pb"
	"io"
)

type grpcServer struct {
	extract kitendpoint.Endpoint
	overlay kitendpoint.Endpoint
	pb.UnimplementedVideoCaptureServer
}

func NewGRPCServer(endpoints endpoint.Set) pb.VideoCaptureServer {

	return &grpcServer{
		extract: endpoints.ExtractEndpoint,
		overlay: endpoints.AddOverlayEndpoint,
	}
}

func encodeGPRCExtractResponse(resp interface{}) *pb.VideoCaptureReply {
	res := resp.(capture.ExtractResponse)
	return &pb.VideoCaptureReply{Data: res.Data}
}

func (s *grpcServer) ExtractImage(stream pb.VideoCapture_ExtractImageServer) error {
	data := []byte{}
	var time int64
	var width, height int32
	var name string

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			request := capture.ExtractRequest{Video: data, Name: name, Height: height, Width: width, Time: time}

			rep, err := s.extract(context.Background(), request)
			if err != nil {
				return err
			}

			return stream.SendAndClose(encodeGPRCExtractResponse(rep))
		}
		if err != nil {
			return err
		}
		data = append(data, req.Video...)
		time = req.Time
		width = req.Width
		height = req.Height
		name = req.Name
	}
}

func (s *grpcServer) AddOverlay(stream pb.VideoCapture_AddOverlayServer) error {
	original := []byte{}
	logo := []byte{}
	var x, y int32

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			request := capture.OverlayImageRequest{Original: original, Overlay: logo, X: x, Y: y}

			rep, err := s.overlay(context.Background(), request)
			if err != nil {
				return err
			}

			return stream.SendAndClose(encodeGPRCExtractResponse(rep))
		}
		if err != nil {
			return err
		}
		original = append(original, req.Original...)
		logo = append(logo, req.Overlay...)
		x = req.Position.X
		y = req.Position.Y
	}
}
