package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"path"
	"sort"

	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/prufen/prufen/cjail/proto"
)

type cjailService struct {
	pb.UnimplementedCJailServer

	imagesRefToDir map[string]string
}

func (s *cjailService) ListImages(ctx context.Context, req *pb.ListImagesRequest) (*pb.ListImagesResponse, error) {
	refs := []string{}
	for ref, _ := range s.imagesRefToDir {
		refs = append(refs, ref)
	}
	sort.Strings(refs)
	log.Printf("available image refs: %v", refs)

	resp := &pb.ListImagesResponse{
		ImageRef: refs,
	}

	return resp, nil
}

type options struct {
	ListenAddress         string `long:"listen-address"             env:"CJAIL_LISTEN_ADDRESS"             default:"localhost:8080" description:"gRPC listen address"`
	ImagesRefToDirJsonMap string `long:"images-ref-to-dir-json-map" env:"CJAIL_IMAGES_REF_TO_DIR_JSON_MAP" default:"{}"             description:"JSON map from image refs to directories with image files"`
}

func main() {
	log.Printf("starting CJail server...")

	var opts options

	parser := flags.NewNamedParser(path.Base(os.Args[0]), flags.Default)
	parser.AddGroup("Application Options", "", &opts)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	imagesRefToDir := map[string]string{}
	if err := json.Unmarshal([]byte(opts.ImagesRefToDirJsonMap), &imagesRefToDir); err != nil {
		log.Fatalf("Failed to parse images refs to directories JSON map: %v\nJSON map:\n%v", err, opts.ImagesRefToDirJsonMap)
	}

	service := &cjailService{
		imagesRefToDir: imagesRefToDir,
	}

	log.Printf("listening on %q", opts.ListenAddress)
	listener, err := net.Listen("tcp", opts.ListenAddress)
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCJailServer(grpcServer, service)
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
