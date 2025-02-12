//go:build example

// + build example

package example_test

import (
	"context"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/status"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func runServer(t *testing.T, ctx context.Context) {

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		t.Error(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &server{})

	go func() {
		if err := grpcServer.Serve(lis); err != nil && err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	<-ctx.Done()

	grpcServer.GracefulStop()
	if err := lis.Close(); err != nil {
		t.Error(err)
	}
}

func TestProxyGrpc(t *testing.T) {

	targetDir := "./../.."
	changeDirectory(t, targetDir)

	env := []string{}
	config := []string{"./_example/proxy-grpc/config-http-http.yaml"}
	entrypoint := getEntrypointRunner(t, env, config)

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(5*time.Second, cancel)

	go runServer(t, ctx)
	time.Sleep(1 * time.Second)

	conn, err := grpc.NewClient("localhost:50000", grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	var in *pb.HelloReply
	go func() {
		client := pb.NewGreeterClient(conn)
		in, err = client.SayHello(ctx, &pb.HelloRequest{Name: "AILERON"})
		if err != nil {
			log.Printf("error getting feature: %v", err)
		} else {
			log.Printf("feature: %v", in)
		}
		timer.Stop()
		cancel()
	}()

	if err := entrypoint.Run(ctx); err != nil {
		t.Error(err)
	}

	testutil.Diff(t, codes.OK, status.Code(err))
	testutil.Diff(t, "Hello AILERON", in.GetMessage())
}
