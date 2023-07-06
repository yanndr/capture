VERSION=$(shell cat cmd/CaptureService/version)
LDFLAGS=-ldflags "-X main.version=${VERSION}"

default:
	go build cmd/CaptureService/main.go

build:
	go build ${LDFLAGS} -o CaptureService cmd/CaptureService/main.go
	
protos:
	protoc  --go_out=./ --go-grpc_out=./ pb/capture.proto

docker-service:
	docker build -t ydruffin/capture:latest .
	docker build -t ydruffin/capture:${VERSION} .
