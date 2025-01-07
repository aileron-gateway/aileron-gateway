//go:build integration
// +build integration

package httpproxy_test

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/net/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/route_guide/routeguide"
)

func TestStream_SSE(t *testing.T) {

	configs := []string{
		testDataDir + "config-stream.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	svr := &http.Server{
		Addr: ":10001",
		Handler: &returnChunkedHandler{
			header: http.Header{
				"Content-Type":      []string{"text/event-stream"},
				"Transfer-Encoding": []string{"identity"},
			},
			interval: 10 * time.Millisecond,
			bodies:   []string{"1", "2", "3", "4", "5"},
		},
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	r := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
	w := httptest.NewRecorder()
	ww := &wrappedResponseWriter{ResponseRecorder: w}
	handler.ServeHTTP(ww, r)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	testutil.Diff(t, []string{"1", "2", "3", "4", "5"}, ww.bodies)
	t.Logf("%#v\n", ww.milliTimes)

}

func TestStream_ChunkedResponse(t *testing.T) {

	configs := []string{
		testDataDir + "config-stream.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	svr := &http.Server{
		Addr: ":10001",
		Handler: &returnChunkedHandler{
			header: http.Header{
				"Content-Type":      []string{"text/plain"},
				"Transfer-Encoding": []string{"chunked"},
			},
			interval: 10 * time.Millisecond,
			bodies:   []string{"1", "2", "3", "4", "5"},
		},
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	r := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
	w := httptest.NewRecorder()
	ww := &wrappedResponseWriter{ResponseRecorder: w}
	handler.ServeHTTP(ww, r)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	testutil.Diff(t, []string{"1", "2", "3", "4", "5"}, ww.bodies)
	t.Logf("%#v\n", ww.milliTimes)

}

func TestStream_ChunkedRequest(t *testing.T) {

	configs := []string{
		testDataDir + "config-stream.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	ch := &receiveChunkedHandler{}
	svr := &http.Server{
		Addr:    ":10001",
		Handler: ch,
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	rd, wr := io.Pipe()
	r := httptest.NewRequest(http.MethodGet, "http://test.com/test", rd)
	r.Header.Set("Transfer-Encoding", "chunked")
	w := httptest.NewRecorder()
	go func() {
		for _, b := range []string{"1", "2", "3", "4", "5"} {
			wr.Write([]byte(b))
			time.Sleep(10 * time.Millisecond)
		}
		wr.Close()
	}()
	handler.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Result().Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	testutil.Diff(t, "ok", string(b))
	testutil.Diff(t, []string{"1", "2", "3", "4", "5"}, ch.bodies)
	t.Logf("%#v\n", ch.milliTimes)

}

func TestStream_WebSocket(t *testing.T) {

	configs := []string{
		testDataDir + "config-stream-websocket.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "HTTPServer",
		Name:       "default",
		Namespace:  "",
	}
	runner, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	wsh := &webSocketHandler{}
	svr := &http.Server{
		Addr:    ":10001",
		Handler: websocket.Handler(wsh.ServeWebSocket),
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { runner.Run(ctx) }()
	time.Sleep(time.Second) // Wait a little until the server starts.

	ws, err := websocket.Dial("ws://localhost:11000/test", "", "http://localhost:11000")
	testutil.DiffError(t, nil, nil, err)
	defer ws.Close()

	for _, msg := range []string{"1", "2", "3", "4", "5"} {
		websocket.Message.Send(ws, msg)
		if err != nil {
			panic(err)
		}
		rcv := ""
		err := websocket.Message.Receive(ws, &rcv)
		if err != nil {
			panic(err)
		}
		testutil.Diff(t, msg+" ok", rcv)
		time.Sleep(10 * time.Millisecond)
	}
	testutil.Diff(t, []string{"1", "2", "3", "4", "5"}, wsh.bodies)
	t.Logf("%#v\n", wsh.milliTimes)

}

func TestStream_SendOctetStream(t *testing.T) {

	configs := []string{
		testDataDir + "config-stream.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	bh := &receiveBinaryHandler{}
	svr := &http.Server{
		Addr:    ":10001",
		Handler: bh,
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	rd, wr := io.Pipe()
	r := httptest.NewRequest(http.MethodGet, "http://test.com/test", rd)
	r.Header.Set("Content-Type", "application/octet-stream")
	w := httptest.NewRecorder()
	md5Hash := []byte{}
	go func() {
		hash := md5.New()
		hashWriter := io.MultiWriter(hash, wr)
		buf := make([]byte, 1000)
		for i := 0; i < 10_000; i++ {
			rand.Read(buf)
			hashWriter.Write(buf)
		}
		wr.Close()
		md5Hash = hash.Sum(nil)
	}()
	handler.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Result().Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	testutil.Diff(t, "ok", string(b))
	testutil.Diff(t, int64(10_000*1000), bh.size)
	testutil.Diff(t, hex.EncodeToString(md5Hash), hex.EncodeToString(bh.md5Hash))
	t.Log(hex.EncodeToString(md5Hash), hex.EncodeToString(bh.md5Hash))

}

func TestStream_ReceiveOctetStream(t *testing.T) {

	configs := []string{
		testDataDir + "config-stream.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	bh := &returnBinaryHandler{}
	svr := &http.Server{
		Addr:    ":10001",
		Handler: bh,
	}
	go func() { svr.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer svr.Close()

	r := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	hash := md5.New()
	size, _ := io.Copy(hash, w.Body)
	md5Hash := hash.Sum(nil)

	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	testutil.Diff(t, size, bh.size)
	testutil.Diff(t, hex.EncodeToString(md5Hash), hex.EncodeToString(bh.md5Hash))
	t.Log(hex.EncodeToString(md5Hash), hex.EncodeToString(bh.md5Hash))

}

func TestStream_Grpc(t *testing.T) {

	// This gRPC test leverages the resource of gRPC official contents.
	// https://grpc.io/docs/languages/go/basics/#bidirectional-streaming-rpc
	// https://github.com/grpc/grpc-go/tree/master/examples/route_guide

	configs := []string{
		testDataDir + "config-stream-grpc.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// Both server and client use TLS.
	mux := &http.ServeMux{}
	mux.Handle("/", handler)
	proxy := &http.Server{
		Addr:    ":14444",
		Handler: mux,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	go func() { proxy.ListenAndServeTLS(testDataDir+"testdata/cert.pem", testDataDir+"testdata/key.pem") }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer proxy.Close()

	svr := newGrpcServer("localhost:15555")
	defer svr.Stop()

	pem, _ := os.ReadFile(testDataDir + "testdata/cert.pem")
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{RootCAs: pool})),
	}
	conn, client := newGrpcClient("localhost:14444", opts...)
	defer conn.Close()

	// A simple RPC
	// https://grpc.io/docs/languages/go/basics/#simple-rpc
	feature := printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})
	testutil.Diff(t, int32(409146138), feature.Location.Latitude)
	testutil.Diff(t, int32(-746188906), feature.Location.Longitude)

	// A server-side streaming RPC
	// https://grpc.io/docs/languages/go/basics/#server-side-streaming-rpc
	features := printFeatures(client, &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	})
	testutil.Diff(t, int32(408122808), features[1].Location.Latitude)
	testutil.Diff(t, int32(-743999179), features[1].Location.Longitude)

	// A client-side streaming RPC
	// https://grpc.io/docs/languages/go/basics/#client-side-streaming-rpc
	summery := runRecordRoute(client)
	testutil.Diff(t, true, summery.PointCount > 0)
	testutil.Diff(t, true, summery.Distance > 0)

	// A bidirectional streaming RPC
	// https://grpc.io/docs/languages/go/basics/#bidirectional-streaming-rpc
	notes := runRouteChat(client)
	testutil.Diff(t, int32(0), notes[2].Location.Latitude)
	testutil.Diff(t, int32(3), notes[2].Location.Longitude)

}

func TestStream_GrpcNonTLSServer(t *testing.T) {

	// This gRPC test leverages the resource of gRPC official contents.
	// https://grpc.io/docs/languages/go/basics/#bidirectional-streaming-rpc
	// https://github.com/grpc/grpc-go/tree/master/examples/route_guide

	configs := []string{
		testDataDir + "config-stream-grpc.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// Proxy server must use H2c (Non TLS) to
	// accept insecure gRPC (Non TLS).
	mux := &http.ServeMux{}
	mux.Handle("/", handler)
	proxy := &http.Server{
		Addr:    ":14444",
		Handler: h2c.NewHandler(handler, &http2.Server{}),
	}

	go func() { proxy.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer proxy.Close()

	svr := newGrpcServer("localhost:15555")
	defer svr.Stop()

	opts := []grpc.DialOption{
		// Because the server supports H2c,
		// we can use insecure.NewCredentials().
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, client := newGrpcClient("localhost:14444", opts...)
	defer conn.Close()

	// A server-side streaming RPC
	// https://grpc.io/docs/languages/go/basics/#server-side-streaming-rpc
	features := printFeatures(client, &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	})
	testutil.Diff(t, int32(408122808), features[1].Location.Latitude)
	testutil.Diff(t, int32(-743999179), features[1].Location.Longitude)

	// A client-side streaming RPC
	// https://grpc.io/docs/languages/go/basics/#client-side-streaming-rpc
	summery := runRecordRoute(client)
	testutil.Diff(t, true, summery.PointCount > 0)
	testutil.Diff(t, true, summery.Distance > 0)

	// A bidirectional streaming RPC
	// https://grpc.io/docs/languages/go/basics/#bidirectional-streaming-rpc
	notes := runRouteChat(client)
	testutil.Diff(t, int32(0), notes[2].Location.Latitude)
	testutil.Diff(t, int32(3), notes[2].Location.Longitude)

}

func TestStream_GrpcNonTLSClient(t *testing.T) {

	// This gRPC test leverages the resource of gRPC official contents.
	// https://grpc.io/docs/languages/go/basics/#bidirectional-streaming-rpc
	// https://github.com/grpc/grpc-go/tree/master/examples/route_guide

	configs := []string{
		testDataDir + "config-stream-grpc.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ReverseProxyHandler",
		Name:       "default",
		Namespace:  "",
	}
	handler, err := api.ReferTypedObject[http.Handler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// Proxy server does not support H2c (Non TLS) and working as a TLS server.
	// So, insecure gRPC (Non TLS) client need to skip
	// X509 certificate verification.
	mux := &http.ServeMux{}
	mux.Handle("/", handler)
	proxy := &http.Server{
		Addr:    ":14444",
		Handler: handler,
	}
	go func() { proxy.ListenAndServeTLS(testDataDir+"testdata/cert.pem", testDataDir+"testdata/key.pem") }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer proxy.Close()

	svr := newGrpcServer("localhost:15555")
	defer svr.Stop()

	pem, _ := os.ReadFile(testDataDir + "testdata/cert.pem")
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)
	opts := []grpc.DialOption{
		// grpc.WithTransportCredentials(insecure.NewCredentials()), // This will fail. Because the server requires TLS.
		// Instead of using insecure client, we skip and ignore insecure certification.
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
	}
	conn, client := newGrpcClient("localhost:14444", opts...)
	defer conn.Close()

	// A server-side streaming RPC
	// https://grpc.io/docs/languages/go/basics/#server-side-streaming-rpc
	features := printFeatures(client, &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	})
	testutil.Diff(t, int32(408122808), features[1].Location.Latitude)
	testutil.Diff(t, int32(-743999179), features[1].Location.Longitude)

	// A client-side streaming RPC
	// https://grpc.io/docs/languages/go/basics/#client-side-streaming-rpc
	summery := runRecordRoute(client)
	testutil.Diff(t, true, summery.PointCount > 0)
	testutil.Diff(t, true, summery.Distance > 0)

	// A bidirectional streaming RPC
	// https://grpc.io/docs/languages/go/basics/#bidirectional-streaming-rpc
	notes := runRouteChat(client)
	testutil.Diff(t, int32(0), notes[2].Location.Latitude)
	testutil.Diff(t, int32(3), notes[2].Location.Longitude)

}
