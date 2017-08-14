package transport

import (
	"context"
	"io"

	"github.com/yanndr/capture"
	"github.com/yanndr/capture/endpoint"
	"github.com/yanndr/capture/pb"
)

type grpcServer struct {
	endpoints endpoint.Set
}

func NewGRPCServer(endpoints endpoint.Set) pb.VideoCaptureServer {

	return &grpcServer{
		endpoints: endpoints,
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

			rep, err := s.endpoints.ExtractEndpoint(context.Background(), request)
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
