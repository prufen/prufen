package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"sort"
	"syscall"

	"github.com/jessevdk/go-flags"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "github.com/prufen/prufen/cjail/proto"
)

type cjailService struct {
	pb.UnimplementedCJailServer

	imagesRefToDir map[string]string
	nsjailPath     string
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

func (s *cjailService) executeImpl(ctx context.Context, imageDir string, execFile string, execArgs []string) (*pb.ExecuteResponse, error) {
	log.Printf("using base image from %q", imageDir)

	c := []string{
		s.nsjailPath,
		"-Mo",
		"--chroot", imageDir,
		"-v",
		// /proc/ is needed for AddressSanitizer.
		//"--disable_proc",
		// TODO: setup CWD
		//"--cwd", "/tmp",
		// Write nsjail log to the create pipe.
		// We will pass created pipe file descriptor to nsjail os.exec, from os.exec docs:
		// "entry i becomes file descriptor 3+i",
		// i.e. single provided file descriptor will have number 3 in nsjail (but may be different in the current process).
		"--log_fd", "3",
		// TODO: should be overridable in the request.
		"--time_limit", "60", // seconds
		// Disallow mount and pivot_root syscalls --- they were explicitly enabled in K8s AppArmor profile to make mount namespaces work.
		// Disallow part of tracing. ptrace is needed for LeakSanitizer.
		// Disallow network related syscalls.
		"--seccomp_string", "POLICY jail { ERRNO(1) { mount, pivot_root, process_vm_readv, process_vm_writev, bind, listen, accept, accept4, connect, getsockopt, setsockopt, getsockname, getpeername } } USE jail DEFAULT ALLOW",
		// Set file size to 50 MB.
		"--rlimit_fsize", "50",
		// Use unlimited address space size to make AddressSanitizer/Valgrind work.
		"--rlimit_as", "inf",
		"--", execFile,
	}
	for _, arg := range execArgs {
		c = append(c, arg)
	}

	cmd := exec.CommandContext(ctx, c[0], c[1:]...)

	// Setup pipe for nsjail log.
	// TODO: be sure that it will get closed.
	logr, logw, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("os.Pipe(): %v", err)
	}
	cmd.ExtraFiles = []*os.File{logw}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("cmd.StderrPipe() failed for %#v: %v", c, err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("cmd.StdoutPipe() failed for %#v: %v", c, err)
	}

	eg, _ := errgroup.WithContext(ctx)
	// Wait for nsjail separately, as we need termination of process to get all logs.
	egNsjail, _ := errgroup.WithContext(ctx)

	// Read nsjail logs.
	var debugBytes []byte
	egNsjail.Go(func() error {
		defer logr.Close()

		r := bufio.NewReader(logr)
		b := make([]byte, 4096)
		for {
			n, err := r.Read(b)
			if err == io.EOF {
				break
			}
			debugBytes = append(debugBytes, b[:n]...)
			log.Printf("nsjail output:\n%s", string(b[:n]))
		}
		log.Printf("nsjail output ended")
		return nil
	})

	// Read command stderr.
	// TODO: Limit sizes of stdout and stderr.
	var stderrBytes []byte
	eg.Go(func() error {
		defer stderr.Close()

		r := bufio.NewReader(stderr)
		b := make([]byte, 4096)
		for {
			n, err := r.Read(b)
			if err == io.EOF {
				break
			}
			log.Printf("stderr output:\n%s", string(b[:n]))
			stderrBytes = append(stderrBytes, b[:n]...)
		}

		log.Printf("stderr output ended")
		return nil
	})

	// Read command stdout.
	// TODO: Limit sizes of stdout and stderr.
	var stdoutBytes []byte
	eg.Go(func() error {
		defer stdout.Close()

		r := bufio.NewReader(stdout)
		b := make([]byte, 4096)
		for {
			n, err := r.Read(b)
			if err == io.EOF {
				break
			}
			log.Printf("stdout output:\n%s", string(b[:n]))
			stdoutBytes = append(stdoutBytes, b[:n]...)
		}

		log.Printf("stdout output ended")
		return nil
	})

	log.Printf("running %#v", c)
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start %#v: %v", c, err)
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("reading stdout/stderr goroutines failed: %v", err)
	}

	exitCode := 0
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			exitCode = exiterr.Sys().(syscall.WaitStatus).ExitStatus()
			log.Printf("%#v completed with the exit code %d", c, exitCode)
		} else {
			// TODO: handle waiting of egNsjail.
			return nil, fmt.Errorf("failed to execute %#v: %v", c, err)
		}
	} else {
		log.Printf("successfully completed %#v", c)
	}

	// TODO: handle error.
	logw.Close()
	if err := egNsjail.Wait(); err != nil {
		return nil, fmt.Errorf("reading nsjail output failed: %v", err)
	}

	resp := &pb.ExecuteResponse{
		ExitCode:    int64(exitCode),
		Stdout:      stdoutBytes,
		Stderr:      stderrBytes,
		DebugOutput: debugBytes,
	}

	return resp, nil
}

func (s *cjailService) Execute(ctx context.Context, req *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
	log.Printf("execute(%v)", req)

	if req.GetBaseImageRef() == "" {
		return nil, status.Errorf(codes.NotFound, "base image ref not set")
	}

	imageDir, found := s.imagesRefToDir[req.GetBaseImageRef()]
	if !found {
		return nil, status.Errorf(codes.NotFound, "base image ref %q not found", req.GetBaseImageRef())
	}

	resp, err := s.executeImpl(ctx, imageDir, req.GetFile(), req.GetArgs())
	if err != nil {
		log.Print(err)
		status.Errorf(codes.Internal, "%v", err)
	}

	return resp, nil
}

type options struct {
	ListenAddress         string `long:"listen-address"             env:"CJAIL_LISTEN_ADDRESS"             default:"localhost:8080" description:"gRPC listen address"`
	ImagesRefToDirJsonMap string `long:"images-ref-to-dir-json-map" env:"CJAIL_IMAGES_REF_TO_DIR_JSON_MAP" default:"{}"             description:"JSON map from image refs to directories with image files"`
	NsjailExec            string `long:"nsjail-exec"                env:"CJAIL_NSJAIL_EXEC"                default:"nsjail"         description:"path to nsjail executable"`
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
		log.Fatalf("failed to parse images refs to directories JSON map: %v\nJSON map:\n%v", err, opts.ImagesRefToDirJsonMap)
	}

	nsjailPath, err := exec.LookPath(opts.NsjailExec)
	if err != nil {
		log.Fatalf("nsjail executable not found %q: %v", opts.NsjailExec, err)
	}

	service := &cjailService{
		imagesRefToDir: imagesRefToDir,
		nsjailPath:     nsjailPath,
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
