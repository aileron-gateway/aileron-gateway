// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build example

package example_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func TestMultiServer(t *testing.T) {

	env := []string{}
	config := []string{"../../_example/multi-server/"}
	entrypoint := getEntrypointRunner(t, env, config)

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(5*time.Second, cancel)

	var resp1, resp2 *http.Response
	var err1, err2 error
	go func() {
		req1, _ := http.NewRequest(http.MethodGet, "http://localhost:8081/", nil)
		req1.Header.Set("Accept", "text/plain")
		resp1, err1 = http.DefaultTransport.RoundTrip(req1)

		req2, _ := http.NewRequest(http.MethodGet, "http://localhost:8082/", nil)
		req2.Header.Set("Accept", "text/plain")
		resp2, err2 = http.DefaultTransport.RoundTrip(req2)

		timer.Stop()
		cancel()
	}()

	if err := entrypoint.Run(ctx); err != nil {
		t.Error(err)
	}

	testutil.Diff(t, nil, err1)
	testutil.Diff(t, nil, err2)
	testutil.Diff(t, http.StatusOK, resp1.StatusCode)
	testutil.Diff(t, http.StatusOK, resp2.StatusCode)

}
