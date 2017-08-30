FROM golang:1.8

WORKDIR /go/src/app
COPY . .
RUN apt-get update
RUN apt-get install -y libavcodec-dev libavformat-dev libavutil-dev libswscale-dev
RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

CMD ["go-wrapper", "run"] # ["app"]