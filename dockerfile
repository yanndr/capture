FROM golang:1.20 as builder

RUN apt-get update
RUN apt-get install -y libavcodec-dev libavformat-dev libavutil-dev libswscale-dev
COPY . /go/src
WORKDIR /go/src

RUN make build

FROM debian:stretch-slim
RUN apt-get update && apt-get install -y libavcodec-dev libavformat-dev libavutil-dev libswscale-dev
COPY --from=builder /go/src/CaptureService /bin/CaptureService
#ENTRYPOINT ["/bin/CaptureService"]