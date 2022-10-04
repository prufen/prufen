package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	v8 "rogchap.com/v8go"

	pb "github.com/prufen/prufen/jail/js/proto"
)

type jsjailService struct {
	pb.UnimplementedJSJailServer
}

func (s *jsjailService) Execute(ctx context.Context, req *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
	v8ctx := v8.NewContext()
	defer v8ctx.Close()

	resp := &pb.ExecuteResponse{}
	result, err := v8ctx.RunScript(req.GetScript(), "script.js")
	if result != nil {
		resp.Result = fmt.Sprintf("%v", result)
	}
	if err != nil {
		resp.Error = fmt.Sprintf("%v", err)
	}

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
	pb.RegisterJSJailServer(grpcServer, &jsjailService{})
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
