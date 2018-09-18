VERSION=`cat cmd/CaptureService/version`

default:
	go build cmd/CaptureService/main.go
	
protos:
	protoc -I pb/ pb/capture.proto --go_out=plugins=grpc:pb

docker-service:
	docker build -t ydruffin/capture:latest ./cmd/CaptureService/
	docker build -t ydruffin/capture:${VERSION} ./cmd/CaptureService/
