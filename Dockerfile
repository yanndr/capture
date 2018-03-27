FROM golang:1.10 as builder


RUN apt-get update
RUN apt-get install -y libavcodec-dev libavformat-dev libavutil-dev libswscale-dev
RUN mkdir /go/src/app
COPY ./cmd/CaptureService /go/src/app
RUN go get ./...
RUN go install app

FROM debian:jessie-slim
RUN apt-get update && apt-get install -y libavcodec-dev libavformat-dev libavutil-dev libswscale-dev
COPY --from=builder /go/bin/app /bin/app
ENTRYPOINT ["/bin/app"]