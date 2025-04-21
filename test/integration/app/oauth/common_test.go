// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package oauth_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	apps "github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

func testOAuthAuthenticationHandler(t *testing.T, h apps.AuthenticationHandler) {

	t.Helper()

	// errPattern := regexp.MustCompile(core.ErrPrefix + `authentication failed`)

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
	r1.Header = http.Header{"Authorization": {"Bearer FooBar"}}
	w1 := httptest.NewRecorder()
	_, result, shouldReturn, err := h.ServeAuthn(w1, r1)
	testutil.Diff(t, apps.AuthContinue, result)
	testutil.Diff(t, false, shouldReturn)
	testutil.DiffError(t, nil, nil, err)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, ``, string(b1))

}

func TestMinimalWithoutMetadata(t *testing.T) {

	configs := []string{"./config-minimal-without-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
		Name:       "", // "default" will be used.
		Namespace:  "", // "default" will be used.
	}
	eh, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testOAuthAuthenticationHandler(t, eh)

}

func TestMinimalWithMetadata(t *testing.T) {

	configs := []string{"./config-minimal-with-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	eh, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testOAuthAuthenticationHandler(t, eh)

}

func TestEmptyName(t *testing.T) {

	configs := []string{"./config-empty-name.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
		Name:       "default",
		Namespace:  "testNamespace",
	}
	eh, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testOAuthAuthenticationHandler(t, eh)

}

func TestEmptyNamespace(t *testing.T) {

	configs := []string{"./config-empty-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
		Name:       "testName",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testOAuthAuthenticationHandler(t, eh)

}

func TestEmptyNameNamespace(t *testing.T) {

	configs := []string{"./config-empty-name-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
		Name:       "default",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testOAuthAuthenticationHandler(t, eh)

}

func TestEmptySpec(t *testing.T) {

	configs := []string{"./config-empty-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	errPattern := regexp.MustCompile(core.ErrPrefix + `failed to load configs. config-empty-spec.yaml`)
	testutil.DiffError(t, app.ErrAppMainLoadConfigs, errPattern, err)

}

func TestInvalidSpec(t *testing.T) {

	configs := []string{"./config-invalid-spec.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	errPattern := regexp.MustCompile(core.ErrPrefix + `failed to load configs. config-invalid-spec.yaml`)
	testutil.DiffError(t, app.ErrAppMainLoadConfigs, errPattern, err)

}
