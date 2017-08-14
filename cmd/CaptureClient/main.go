package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/yanndr/capture/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address = "localhost:50051"
)

func main() {
	var c pb.VideoCaptureClient
	creds, err := credentials.NewClientTLSFromFile("../cert.pem", "")
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Printf("did not connect: %v", err)
		return
	}

	log.Println("Connected")
	defer conn.Close()
	c = pb.NewVideoCaptureClient(conn)
	log.Println(c)

	fi, err := os.Open("../IMG_3116.mp4")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer fi.Close()

	stream, err := c.ExtractImage(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		chunk := make([]byte, 64*1024)
		n, err := fi.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			return
		}
		if n < len(chunk) {
			chunk = chunk[:n]
		}
		stream.Send(&pb.VideoCaptureRequest{Name: fi.Name(), Video: chunk, Width: 640, Height: 480, Time: 10000})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Println("failed extract")
		log.Println(err)
		return
	}

	f, err := os.Create("../result.png")
	defer f.Close()

	if err != nil {
		log.Println(err)
		return
	}
	_, err = f.Write(res.Data)
	if err != nil {
		log.Println(err)
		return
	}
}
