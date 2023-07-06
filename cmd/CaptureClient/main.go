package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
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

	outputPtr := flag.String("o", "result.png", "output of the result")

	flag.Parse()

	if len(args) == 0 {
		log.Fatal("no input file")
		return
	}

	var config struct {
		CertPath string
	}

	if err := envconfig.Process("CAPTURE", &config); err != nil {
		log.Fatalf("error getting the environment variables: %v", err)
	}

	var c pb.VideoCaptureClient

	var opts []grpc.DialOption

	if config.CertPath != "" {
		creds, err := credentials.NewClientTLSFromFile(config.CertPath, "")
		if err != nil {
			log.Printf("error with certificate: %v", err)
			return
		}
		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	} else {
		opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	}

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

	var data []byte
	if len(args) > 1 {
		r := bytes.NewReader(res.Data)
		if data, err = addOverlay(c, r, args[1]); err != nil {
			log.Fatal(err)
		}

	} else {
		data = res.Data
	}

	f, err := os.Create(*outputPtr)
	defer f.Close()

	if err != nil {
		log.Println(err)
		return
	}
	_, err = f.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func addOverlay(c pb.VideoCaptureClient, image io.Reader, overlay string) ([]byte, error) {
	str, err := c.AddOverlay(context.Background())
	if err != nil {
		return nil, err
	}

	for {
		chunk := make([]byte, 64*1024)
		n, err := image.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		if n < len(chunk) {
			chunk = chunk[:n]
		}
		str.Send(&pb.OverlayImageRequest{Original: chunk, Position: &pb.Position{X: -10, Y: -10}})
	}

	fr, err := os.Open(overlay)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer fr.Close()

	for {
		chunk := make([]byte, 64*1024)
		n, err := fr.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if n < len(chunk) {
			chunk = chunk[:n]
		}
		str.Send(&pb.OverlayImageRequest{Overlay: chunk, Position: &pb.Position{X: -10, Y: -10}})
	}

	resp, err := str.CloseAndRecv()
	if err != nil {
		return nil, fmt.Errorf("error on close and received: %v", err)
	}

	return resp.Data, nil
}
