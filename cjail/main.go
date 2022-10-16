package main

import (
	"context"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/prufen/prufen/cjail/proto"
)

type cjailService struct {
	pb.UnimplementedCJailServer
}

func (s *cjailService) ListImages(ctx context.Context, req *pb.ListImagesRequest) (*pb.ListImagesResponse, error) {
	resp := &pb.ListImagesResponse{}

	// TODO: return the list of the available images.

	log.Printf("available images: %+v", resp)
	return resp, nil
}

func main() {
	log.Printf("starting server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCJailServer(grpcServer, &cjailService{})
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
