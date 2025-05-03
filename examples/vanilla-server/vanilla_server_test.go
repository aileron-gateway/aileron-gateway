// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package main_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/util/register"
)

func TestSingle(t *testing.T) {
	entrypoint := getEntrypointRunner(t, "./config-single.yaml")

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

		timer.Stop() // Stop the timer
		cancel()     // and immediately stop the server.
	}()

	if err := entrypoint.Run(ctx); err != nil {
		t.Error(err)
	}
	testutil.Diff(t, nil, err1)
	testutil.Diff(t, nil, err2)
	testutil.Diff(t, http.StatusOK, resp1.StatusCode)
	testutil.Diff(t, http.StatusOK, resp2.StatusCode)
}

func getEntrypointRunner(t *testing.T, config ...string) core.Runner {
	t.Helper()
	svr := api.NewDefaultServeMux()
	f := api.NewFactoryAPI()
	register.RegisterAll(f)
	svr.Handle("core/", f)
	if err := app.LoadConfigFiles(svr, config); err != nil {
		t.Error(err)
	}
	req := &api.Request{
		Method: api.MethodGet,
		Key:    "core/v1/Entrypoint",
		Format: api.FormatProtoReference,
		Content: &kernel.Reference{
			APIVersion: "core/v1",
			Kind:       "Entrypoint",
			Namespace:  ".entrypoint",
			Name:       ".entrypoint",
		},
	}
	res, err := svr.Serve(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	entrypoint, ok := res.Content.(core.Runner)
	if !ok {
		t.Error(errors.New("failed to assert entrypoint to Runner interface"))
	}
	return entrypoint
}
