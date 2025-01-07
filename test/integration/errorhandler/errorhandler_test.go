//go:build integration
// +build integration

package errorhandler_test

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/errorutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

var (
	ErrFooFoo = errorutil.NewKind("E1001", "FooFailure", "failed to run foo caused by file not exist")
	ErrBarBar = errorutil.NewKind("E2001", "BarFailure", "failed to run bar caused by file permission error")
)

func TestDefault(t *testing.T) {

	configs := []string{
		testDataDir + "config-default.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	tmp := os.Stdout
	defer func() { os.Stdout = tmp }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Replace default logger with info level.
	tmpLg := log.GlobalLogger(log.DefaultLoggerName)
	log.SetGlobalLogger(log.DefaultLoggerName,
		log.NewJSONSLogger(w, &slog.HandlerOptions{Level: slog.LevelInfo}))
	defer log.SetGlobalLogger(log.DefaultLoggerName, tmpLg)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ErrorHandler",
		Name:       "default",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.ErrorHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	w1 := httptest.NewRecorder()
	eh.ServeHTTPError(w1, req, errors.New("test"))
	testutil.Diff(t, http.StatusInternalServerError, w1.Result().StatusCode)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, `{"status":500,"statusText":"Internal Server Error"}`, string(b1))

	w2 := httptest.NewRecorder()
	eh.ServeHTTPError(w2, req, utilhttp.NewHTTPError(nil, http.StatusUnauthorized))
	testutil.Diff(t, http.StatusUnauthorized, w2.Result().StatusCode)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, `{"status":401,"statusText":"Unauthorized"}`, string(b2))

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	t.Log(buf.String())
	testutil.Diff(t, true, strings.Contains(buf.String(), `goroutine`)) // Stack trace exists.
	testutil.Diff(t, true, strings.Contains(buf.String(), `"msg":"serve http error. status=500"`))
	testutil.Diff(t, false, strings.Contains(buf.String(), `"msg":"serve http error. status=401"`))

}

func TestStackAlways(t *testing.T) {

	configs := []string{
		testDataDir + "config-stack-always.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	tmp := os.Stdout
	defer func() { os.Stdout = tmp }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Replace default logger with debug level.
	tmpLg := log.GlobalLogger(log.DefaultLoggerName)
	log.SetGlobalLogger(log.DefaultLoggerName,
		log.NewJSONSLogger(w, &slog.HandlerOptions{Level: slog.LevelDebug}))
	defer log.SetGlobalLogger(log.DefaultLoggerName, tmpLg)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ErrorHandler",
		Name:       "default",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.ErrorHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	w1 := httptest.NewRecorder()
	eh.ServeHTTPError(w1, req, errors.New("test"))
	testutil.Diff(t, http.StatusInternalServerError, w1.Result().StatusCode)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, `{"status":500,"statusText":"Internal Server Error"}`, string(b1))

	w2 := httptest.NewRecorder()
	eh.ServeHTTPError(w2, req, utilhttp.NewHTTPError(nil, http.StatusUnauthorized))
	testutil.Diff(t, http.StatusUnauthorized, w2.Result().StatusCode)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, `{"status":401,"statusText":"Unauthorized"}`, string(b2))

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	t.Log(buf.String())
	testutil.Diff(t, true, strings.Contains(buf.String(), `goroutine`)) // Stack trace exists.
	testutil.Diff(t, true, strings.Contains(buf.String(), `"msg":"serve http error. status=500"`))
	testutil.Diff(t, true, strings.Contains(buf.String(), `"msg":"serve http error. status=401"`))

}

func TestErrorMessageCode(t *testing.T) {

	configs := []string{
		testDataDir + "config-error-message-code.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ErrorHandler",
		Name:       "default",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.ErrorHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	errFoo := utilhttp.NewHTTPError(ErrFooFoo.WithoutStack(nil, nil), http.StatusUnauthorized)

	r1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r1.Header.Set("Accept", "application/json")
	w1 := httptest.NewRecorder()
	eh.ServeHTTPError(w1, r1, errFoo)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, "bob", w1.Result().Header.Get("alice"))
	testutil.Diff(t, `{"foo":"bar"}`, string(b1))

	r2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r2.Header.Set("Accept", "text/plain")
	w2 := httptest.NewRecorder()
	eh.ServeHTTPError(w2, r2, errFoo)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, "bob", w2.Result().Header.Get("alice"))
	testutil.Diff(t, `foo=bar`, string(b2))

	// First content application/json will be used.
	r3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r3.Header.Set("Accept", "invalid/type")
	w3 := httptest.NewRecorder()
	eh.ServeHTTPError(w3, r3, errFoo)
	testutil.Diff(t, http.StatusOK, w3.Result().StatusCode)
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, "bob", w3.Result().Header.Get("alice"))
	testutil.Diff(t, `{"foo":"bar"}`, string(b3))

	errBar := utilhttp.NewHTTPError(ErrBarBar.WithoutStack(nil, nil), http.StatusUnauthorized)

	r4 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r4.Header.Set("Accept", "application/json")
	w4 := httptest.NewRecorder()
	eh.ServeHTTPError(w4, r4, errBar)
	testutil.Diff(t, http.StatusUnauthorized, w4.Result().StatusCode)
	b4, _ := io.ReadAll(w4.Result().Body)
	testutil.Diff(t, `{"status":401,"statusText":"Unauthorized"}`, string(b4))

	r5 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r5.Header.Set("Accept", "text/plain")
	w5 := httptest.NewRecorder()
	eh.ServeHTTPError(w5, r5, errBar)
	testutil.Diff(t, http.StatusUnauthorized, w5.Result().StatusCode)
	b5, _ := io.ReadAll(w5.Result().Body)
	testutil.Diff(t, "status: 401\nstatusText: Unauthorized\n", string(b5))

}

func TestErrorMessageKind(t *testing.T) {

	configs := []string{
		testDataDir + "config-error-message-kind.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ErrorHandler",
		Name:       "default",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.ErrorHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	errFoo := utilhttp.NewHTTPError(ErrFooFoo.WithoutStack(nil, nil), http.StatusUnauthorized)

	r1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r1.Header.Set("Accept", "application/json")
	w1 := httptest.NewRecorder()
	eh.ServeHTTPError(w1, r1, errFoo)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "bob", w1.Result().Header.Get("alice"))
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, `{"foo":"bar"}`, string(b1))

	r2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r2.Header.Set("Accept", "text/plain")
	w2 := httptest.NewRecorder()
	eh.ServeHTTPError(w2, r2, errFoo)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, "bob", w2.Result().Header.Get("alice"))
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, `foo=bar`, string(b2))

	// First content application/json will be used.
	r3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r3.Header.Set("Accept", "invalid/type")
	w3 := httptest.NewRecorder()
	eh.ServeHTTPError(w3, r3, errFoo)
	testutil.Diff(t, http.StatusOK, w3.Result().StatusCode)
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, `{"foo":"bar"}`, string(b3))

	errBar := utilhttp.NewHTTPError(ErrBarBar.WithoutStack(nil, nil), http.StatusUnauthorized)

	r4 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r4.Header.Set("Accept", "application/json")
	w4 := httptest.NewRecorder()
	eh.ServeHTTPError(w4, r4, errBar)
	testutil.Diff(t, http.StatusUnauthorized, w4.Result().StatusCode)
	b4, _ := io.ReadAll(w4.Result().Body)
	testutil.Diff(t, `{"status":401,"statusText":"Unauthorized"}`, string(b4))

	r5 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r5.Header.Set("Accept", "text/plain")
	w5 := httptest.NewRecorder()
	eh.ServeHTTPError(w5, r5, errBar)
	testutil.Diff(t, http.StatusUnauthorized, w5.Result().StatusCode)
	b5, _ := io.ReadAll(w5.Result().Body)
	testutil.Diff(t, "status: 401\nstatusText: Unauthorized\n", string(b5))

}

func TestErrorMessageMsg(t *testing.T) {

	configs := []string{
		testDataDir + "config-error-message-msg.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "ErrorHandler",
		Name:       "default",
		Namespace:  "",
	}
	eh, err := api.ReferTypedObject[core.ErrorHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	errFoo := utilhttp.NewHTTPError(ErrFooFoo.WithoutStack(nil, nil), http.StatusUnauthorized)

	r1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r1.Header.Set("Accept", "application/json")
	w1 := httptest.NewRecorder()
	eh.ServeHTTPError(w1, r1, errFoo)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
	testutil.Diff(t, "bob", w1.Result().Header.Get("alice"))
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, `{"foo":"bar"}`, string(b1))

	r2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r2.Header.Set("Accept", "text/plain")
	w2 := httptest.NewRecorder()
	eh.ServeHTTPError(w2, r2, errFoo)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode)
	testutil.Diff(t, "bob", w2.Result().Header.Get("alice"))
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, `foo=bar`, string(b2))

	// First content application/json will be used.
	r3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r3.Header.Set("Accept", "invalid/type")
	w3 := httptest.NewRecorder()
	eh.ServeHTTPError(w3, r3, errFoo)
	testutil.Diff(t, http.StatusOK, w3.Result().StatusCode)
	testutil.Diff(t, "bob", w3.Result().Header.Get("alice"))
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, `{"foo":"bar"}`, string(b3))

	errBar := utilhttp.NewHTTPError(ErrBarBar.WithoutStack(nil, nil), http.StatusUnauthorized)

	r4 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r4.Header.Set("Accept", "application/json")
	w4 := httptest.NewRecorder()
	eh.ServeHTTPError(w4, r4, errBar)
	testutil.Diff(t, http.StatusUnauthorized, w4.Result().StatusCode)
	b4, _ := io.ReadAll(w4.Result().Body)
	testutil.Diff(t, `{"status":401,"statusText":"Unauthorized"}`, string(b4))

	r5 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r5.Header.Set("Accept", "text/plain")
	w5 := httptest.NewRecorder()
	eh.ServeHTTPError(w5, r5, errBar)
	testutil.Diff(t, http.StatusUnauthorized, w5.Result().StatusCode)
	b5, _ := io.ReadAll(w5.Result().Body)
	testutil.Diff(t, "status: 401\nstatusText: Unauthorized\n", string(b5))

}
