// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.3
// source: pb/capture.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// VideoCaptureClient is the client API for VideoCapture service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VideoCaptureClient interface {
	ExtractImage(ctx context.Context, opts ...grpc.CallOption) (VideoCapture_ExtractImageClient, error)
	AddOverlay(ctx context.Context, opts ...grpc.CallOption) (VideoCapture_AddOverlayClient, error)
}

type videoCaptureClient struct {
	cc grpc.ClientConnInterface
}

func NewVideoCaptureClient(cc grpc.ClientConnInterface) VideoCaptureClient {
	return &videoCaptureClient{cc}
}

func (c *videoCaptureClient) ExtractImage(ctx context.Context, opts ...grpc.CallOption) (VideoCapture_ExtractImageClient, error) {
	stream, err := c.cc.NewStream(ctx, &VideoCapture_ServiceDesc.Streams[0], "/pb.VideoCapture/ExtractImage", opts...)
	if err != nil {
		return nil, err
	}
	x := &videoCaptureExtractImageClient{stream}
	return x, nil
}

type VideoCapture_ExtractImageClient interface {
	Send(*VideoCaptureRequest) error
	CloseAndRecv() (*VideoCaptureReply, error)
	grpc.ClientStream
}

type videoCaptureExtractImageClient struct {
	grpc.ClientStream
}

func (x *videoCaptureExtractImageClient) Send(m *VideoCaptureRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *videoCaptureExtractImageClient) CloseAndRecv() (*VideoCaptureReply, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(VideoCaptureReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *videoCaptureClient) AddOverlay(ctx context.Context, opts ...grpc.CallOption) (VideoCapture_AddOverlayClient, error) {
	stream, err := c.cc.NewStream(ctx, &VideoCapture_ServiceDesc.Streams[1], "/pb.VideoCapture/AddOverlay", opts...)
	if err != nil {
		return nil, err
	}
	x := &videoCaptureAddOverlayClient{stream}
	return x, nil
}

type VideoCapture_AddOverlayClient interface {
	Send(*OverlayImageRequest) error
	CloseAndRecv() (*VideoCaptureReply, error)
	grpc.ClientStream
}

type videoCaptureAddOverlayClient struct {
	grpc.ClientStream
}

func (x *videoCaptureAddOverlayClient) Send(m *OverlayImageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *videoCaptureAddOverlayClient) CloseAndRecv() (*VideoCaptureReply, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(VideoCaptureReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// VideoCaptureServer is the server API for VideoCapture service.
// All implementations must embed UnimplementedVideoCaptureServer
// for forward compatibility
type VideoCaptureServer interface {
	ExtractImage(VideoCapture_ExtractImageServer) error
	AddOverlay(VideoCapture_AddOverlayServer) error
	mustEmbedUnimplementedVideoCaptureServer()
}

// UnimplementedVideoCaptureServer must be embedded to have forward compatible implementations.
type UnimplementedVideoCaptureServer struct {
}

func (UnimplementedVideoCaptureServer) ExtractImage(VideoCapture_ExtractImageServer) error {
	return status.Errorf(codes.Unimplemented, "method ExtractImage not implemented")
}
func (UnimplementedVideoCaptureServer) AddOverlay(VideoCapture_AddOverlayServer) error {
	return status.Errorf(codes.Unimplemented, "method AddOverlay not implemented")
}
func (UnimplementedVideoCaptureServer) mustEmbedUnimplementedVideoCaptureServer() {}

// UnsafeVideoCaptureServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VideoCaptureServer will
// result in compilation errors.
type UnsafeVideoCaptureServer interface {
	mustEmbedUnimplementedVideoCaptureServer()
}

func RegisterVideoCaptureServer(s grpc.ServiceRegistrar, srv VideoCaptureServer) {
	s.RegisterService(&VideoCapture_ServiceDesc, srv)
}

func _VideoCapture_ExtractImage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(VideoCaptureServer).ExtractImage(&videoCaptureExtractImageServer{stream})
}

type VideoCapture_ExtractImageServer interface {
	SendAndClose(*VideoCaptureReply) error
	Recv() (*VideoCaptureRequest, error)
	grpc.ServerStream
}

type videoCaptureExtractImageServer struct {
	grpc.ServerStream
}

func (x *videoCaptureExtractImageServer) SendAndClose(m *VideoCaptureReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *videoCaptureExtractImageServer) Recv() (*VideoCaptureRequest, error) {
	m := new(VideoCaptureRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _VideoCapture_AddOverlay_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(VideoCaptureServer).AddOverlay(&videoCaptureAddOverlayServer{stream})
}

type VideoCapture_AddOverlayServer interface {
	SendAndClose(*VideoCaptureReply) error
	Recv() (*OverlayImageRequest, error)
	grpc.ServerStream
}

type videoCaptureAddOverlayServer struct {
	grpc.ServerStream
}

func (x *videoCaptureAddOverlayServer) SendAndClose(m *VideoCaptureReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *videoCaptureAddOverlayServer) Recv() (*OverlayImageRequest, error) {
	m := new(OverlayImageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// VideoCapture_ServiceDesc is the grpc.ServiceDesc for VideoCapture service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VideoCapture_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.VideoCapture",
	HandlerType: (*VideoCaptureServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ExtractImage",
			Handler:       _VideoCapture_ExtractImage_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "AddOverlay",
			Handler:       _VideoCapture_AddOverlay_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "pb/capture.proto",
}
