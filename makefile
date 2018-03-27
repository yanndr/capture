protos:
	protoc -I pb/ pb/capture.proto --go_out=plugins=grpc:pb

docker-service:
	docker build -t capture ./cmd/CaptureService/
