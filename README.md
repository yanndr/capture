# capture
GRPC service to capture an image from a video

## About
I needed to extract an image from a video for a personal website project. I wanted to keep the code in an external library; so I created a service that exposed the function by GRPC. 
Why GRPC? Because I wanted to learn about it, [Go-kit](https://github.com/go-kit/kit) and [Promotheus](https://github.com/prometheus/client_golang)

## Installation
This service depends on [screengen](github.com/opennota/screengen), which requires having [ffmpeg](https://ffmpeg.org/) installed if you want to run it locally.

```
go get github.com/yanndr/capture

docker build -t capture ./cmd/CaptureService/

```


## Usage 

Run the Service with Docker

Without a certificate:
```
docker run --rm -it -p  50051:50051 -p 8080:8081  capture
```
With a certificate
```
docker run --rm -it -p  50051:50051 -p 8080:8081 -v /PATHTOCERT/cert -e CAPTURE_CERTPATH='/cert/cert.pem' -e CAPTURE_KEYPATH='cert/key.pem' capture
```

Then the Service will expose 2 functions by GRPC:
rpc ExtractImage (stream VideoCaptureRequest) returns (VideoCaptureReply) {}
rpc AddOverlay (stream OverlayImageRequest) returns (VideoCaptureReply) {}

metrics will be available at http://localhost:8080/metrics

There is a client example in cmd/CaptureClient:

```
go run cmd/CaptureClient/main.go pathtoVideofile 