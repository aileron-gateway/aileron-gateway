//go:build example
// +build example

package example

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func TestPrommetrics(t *testing.T) {

	env := []string{}
	config := []string{"../../_example/prometheus/"}
	entrypoint := getEntrypointRunner(t, env, config)

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(5*time.Second, cancel)

	var resp *http.Response
	var err error
	go func() {
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost:8080/metrics", nil)
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
