//go:build example && !windows
// +build example,!windows

package example_test

import (
	"context"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func TestUnixSocket(t *testing.T) {

	wd, _ := os.Getwd()
	defer changeDirectory(t, wd)
	changeDirectory(t, "./../../")

	env := []string{}
	config := []string{"_example/listen-unix-socket/"}
	entrypoint := getEntrypointRunner(t, env, config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dial := func(network, addr string) (net.Conn, error) {
		return net.Dial("unix", "@gateway")
	}

	transport := &http.Transport{
		Dial: dial,
	}

	var resp *http.Response
	var err error
	go func() {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		req.Header.Set("Accept", "text/plain")
		resp, err = transport.RoundTrip(req)
		cancel()
	}()

	if err := entrypoint.Run(ctx); err != nil {
		t.Error(err)
	}

	testutil.Diff(t, nil, err)
	testutil.Diff(t, http.StatusOK, resp.StatusCode)
}
