package main

import (
	"context"
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

	res, err := c.ExtractImage(context.Background(),
		&pb.VideoCaptureRequest{Path: "../IMG_3116.mp4",
			Width:  640,
			Height: 480,
			Time:   10000,
			 OverlayImage: &pb.OverlayImage{Path: "../forumtube.png"},
		})
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
