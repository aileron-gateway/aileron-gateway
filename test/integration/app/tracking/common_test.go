// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package tracking_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/uid"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

func check(t *testing.T, m core.Middleware, rID, tID bool) {
	t.Helper()

	var reqID, trcID string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.Diff(t, true, uid.IDFromContext(r.Context()) != "")
		h := utilhttp.ProxyHeaderFromContext(r.Context())
		reqID = h.Get("X-Request-ID")
		trcID = h.Get("X-Trace-ID")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	h := m.Middleware(handler)

	r1 := httptest.NewRequest(http.MethodGet, "http://tracking-test.com/test", nil)
	r1.Header.Set("X-Trace-Extract", "existing-trace-id")
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)

	t.Log(reqID, trcID)
	testutil.Diff(t, rID, reqID != "")
	testutil.Diff(t, tID, trcID != "")

	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "ok", string(b1))
}

func TestMinimalWithoutMetadata(t *testing.T) {
	configs := []string{"./config-minimal-without-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "TrackingMiddleware",
		Name:       "default",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m, false, false)
}

func TestMinimalWithMetadata(t *testing.T) {
	configs := []string{"./config-minimal-with-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "TrackingMiddleware",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m, false, false)
}

func TestEmptyName(t *testing.T) {
	configs := []string{"./config-empty-name.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "TrackingMiddleware",
		Name:       "default",
		Namespace:  "testNamespace",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m, false, false)
}

func TestEmptyNamespace(t *testing.T) {
	configs := []string{"./config-empty-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "TrackingMiddleware",
		Name:       "testName",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m, false, false)
}

func TestEmptyNameNamespace(t *testing.T) {
	configs := []string{"./config-empty-name-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "TrackingMiddleware",
		Name:       "default",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, m, false, false)
}

func TestEmptySpec(t *testing.T) {
	configs := []string{"./config-empty-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "TrackingMiddleware",
		Name:       "default",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	check(t, eh, false, false)
}

func TestInvalidSpec(t *testing.T) {
	configs := []string{"./config-invalid-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	errPattern := regexp.MustCompile(core.ErrPrefix + `failed to load configs. config-invalid-spec.yaml`)
	testutil.DiffError(t, app.ErrAppMainLoadConfigs, errPattern, err)
}
