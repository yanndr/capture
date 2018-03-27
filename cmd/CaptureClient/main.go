package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/yanndr/capture/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address = "localhost:50051"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("no input file")
		return
	}

	var config struct {
		CertPath string
	}

	if err := envconfig.Process("CAPTURE", &config); err != nil {
		log.Fatalf("error getting the environment variables: %v", err)
		os.Exit(1)
	}

	var c pb.VideoCaptureClient
	creds, err := credentials.NewClientTLSFromFile(config.CertPath, "")
	if err != nil {
		log.Printf("error with certificate: %v", err)
		return
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Printf("did not connect: %v", err)
		return
	}

	log.Println("Connected")
	defer conn.Close()
	c = pb.NewVideoCaptureClient(conn)

	fi, err := os.Open(args[0])
	if err != nil {
		log.Fatal("Could not open file: ", err)
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

	r := bytes.NewReader(res.Data)

	str, err := c.AddOverlay(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		chunk := make([]byte, 64*1024)
		n, err := r.Read(chunk)
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
		str.Send(&pb.OverlayImageRequest{Original: chunk, Position: &pb.Position{X: -10, Y: -10}})
	}

	if len(args) > 1 {
		fr, err := os.Open(args[1])
		if err != nil {
			log.Fatal(err)
			return
		}
		defer fr.Close()

		for {
			chunk := make([]byte, 64*1024)
			n, err := fr.Read(chunk)
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
			str.Send(&pb.OverlayImageRequest{Overlay: chunk, Position: &pb.Position{X: -10, Y: -10}})
		}

		resp, err := str.CloseAndRecv()
		if err != nil {
			log.Println("failed to add overlay")
			log.Println(err)
			return
		}

		f, err := os.Create("../result.png")
		defer f.Close()

		if err != nil {
			log.Println(err)
			return
		}
		_, err = f.Write(resp.Data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
