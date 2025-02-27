//go:build example
// +build example

package example_test

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
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

	ln, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		t.Error(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &server{})

	go func() {
		if err := grpcServer.Serve(ln); err != nil && err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	<-ctx.Done()

	grpcServer.GracefulStop()
}

func TestProxyGrpc(t *testing.T) {

	wd, _ := os.Getwd()
	defer changeDirectory(t, wd)
	changeDirectory(t, "./../../")

	env := []string{}
	config := []string{"./_example/proxy-grpc/config-http-http.yaml"}
	entrypoint := getEntrypointRunner(t, env, config)

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(5*time.Second, cancel)

	go runServer(t, ctx)
	time.Sleep(1 * time.Second)

	conn, err := grpc.NewClient("localhost:50000", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
