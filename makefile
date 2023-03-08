VERSION=`cat cmd/CaptureService/version`

LDFLAGS=-ldflags "-X main.version=${VERSION}"

default:
	go build cmd/CaptureService/main.go

build:
	go build ${LDFLAGS} -o CaptureService cmd/CaptureService/main.go
	
protos:
	protoc -I pb/ pb/capture.proto --go-grpc_out=pb

docker-service:
	docker build -t ydruffin/capture:latest ./cmd/CaptureService/
	docker build -t ydruffin/capture:${VERSION} ./cmd/CaptureService/
