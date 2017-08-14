package transport

import (
	"context"
	"io"

	"github.com/yanndr/capture"
	"github.com/yanndr/capture/pb"
)

type grpcServer struct {
	// extract grpctransport.Handler
	service capture.Service
}

func NewGRPCServer(svc capture.Service) pb.VideoCaptureServer {

	return &grpcServer{
		service: svc,
	}
}

func decodeGPRCExtractRequest(req *pb.VideoCaptureRequest) capture.ExtractRequest {
	return capture.ExtractRequest{Video: req.Video, Height: req.Height, Width: req.Width, Time: req.Time}
}

func encodeGPRCExtractResponse(resp capture.ExtractResponse) *pb.VideoCaptureReply {
	return &pb.VideoCaptureReply{Data: resp.Data}
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
			rep, err := s.service.Extract(context.Background(), request)
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
