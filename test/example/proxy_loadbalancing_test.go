//go:build example
// +build example

package example_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func runServers(t *testing.T, ctx context.Context) {
	addrs := []string{
		":8001",
		":8002",
		":8003",
		":8004",
		":8005",
	}

	servers := make([]*http.Server, len(addrs))

	for i, addr := range addrs {
		go func(i int, a string) {
			mux := &http.ServeMux{}
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello! %s from %s\n", r.RemoteAddr, a)
			})
			svr := &http.Server{
				Addr:    addr,
				Handler: mux,
			}

			servers[i] = svr

			log.Println("Server listens at", a)
			if err := svr.ListenAndServe(); err != nil {
				t.Error(err)
			}
		}(i, addr)
	}

	time.Sleep(time.Second * 2)

	<-ctx.Done()

	for _, svr := range servers {
		if err := svr.Shutdown(context.Background()); err != nil {
			t.Error(err)
		}
	}
}

func TestProxyLoadBalancing(t *testing.T) {

	targetDir := "./../.."
	changeDirectory(t, targetDir)

	env := []string{}
	config := []string{"./_example/proxy-loadbalancing/config-direct-hash.yaml"}
	entrypoint := getEntrypointRunner(t, env, config)

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(5*time.Second, cancel)

	go runServers(t, ctx)

	time.Sleep(1 * time.Second)

	var resp *http.Response
	var err error

	go func() {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
		resp, err = http.DefaultTransport.RoundTrip(req)
		timer.Stop()
		cancel()
	}()

	if err := entrypoint.Run(ctx); err != nil {
		t.Error(err)
	}

	testutil.Diff(t, nil, err)
	testutil.Diff(t, http.StatusOK, resp.StatusCode)
}
