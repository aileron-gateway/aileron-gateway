//go:build integration

package throttle_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

func TestThrottleMiddleware_MaxConnections(t *testing.T) {

	configs := []string{"./config-throttle-max-connections.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "ThrottleMiddleware",
		Name:       "default",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})
	h := m.Middleware(handler)

	// Run the test case 10 times to check reproducibility.
	for iter := 0; iter < 10; iter++ {

		// Send 5 requests (3 succeeded, 2 failed)
		var wg sync.WaitGroup
		success := atomic.Int32{}
		failure := atomic.Int32{}
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				r := httptest.NewRequest(http.MethodGet, "http://throttle.com/max-connection", nil)
				w := httptest.NewRecorder()
				h.ServeHTTP(w, r)
				b, _ := io.ReadAll(w.Result().Body)
				if w.Result().StatusCode == http.StatusOK {
					success.Add(1)
					testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
					testutil.Diff(t, "test", string(b))
				} else {
					failure.Add(1)
					testutil.Diff(t, http.StatusTooManyRequests, w.Result().StatusCode)
					testutil.Diff(t, `{"status":429,"statusText":"Too Many Requests"}`, string(b))
				}
			}()
		}
		wg.Wait()
		testutil.Diff(t, int32(3), success.Load())
		testutil.Diff(t, int32(2), failure.Load())

	}

}

func TestThrottleMiddleware_TokenBucket(t *testing.T) {

	configs := []string{"./config-throttle-token-bucket.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "ThrottleMiddleware",
		Name:       "default",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})
	h := m.Middleware(handler)

	// Send 10 requests (7 succeeded, 3 failed)
	var wg sync.WaitGroup
	success := atomic.Int32{}
	failure := atomic.Int32{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := httptest.NewRequest(http.MethodGet, "http://throttle.com/token-bucket", nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			b, _ := io.ReadAll(w.Result().Body)
			if w.Result().StatusCode == http.StatusOK {
				success.Add(1)
				testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
				testutil.Diff(t, "test", string(b))
			} else {
				failure.Add(1)
				testutil.Diff(t, http.StatusTooManyRequests, w.Result().StatusCode)
				testutil.Diff(t, `{"status":429,"statusText":"Too Many Requests"}`, string(b))
			}
		}()
		time.Sleep(15 * time.Millisecond)
	}
	wg.Wait()
	testutil.Diff(t, int32(7), success.Load())
	testutil.Diff(t, int32(3), failure.Load())

}

func TestThrottleMiddleware_FixedWindow(t *testing.T) {

	configs := []string{"./config-throttle-fixed-window.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "ThrottleMiddleware",
		Name:       "default",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})
	h := m.Middleware(handler)

	// Run the test case 10 times to check reproducibility.
	for iter := 0; iter < 10; iter++ {

		// Send 10 requests (5 succeeded, 5 failed)
		var wg sync.WaitGroup
		success := atomic.Int32{}
		failure := atomic.Int32{}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				r := httptest.NewRequest(http.MethodGet, "http://throttle.com/fixed-window", nil)
				w := httptest.NewRecorder()
				h.ServeHTTP(w, r)
				b, _ := io.ReadAll(w.Result().Body)
				if w.Result().StatusCode == http.StatusOK {
					success.Add(1)
					testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
					testutil.Diff(t, "test", string(b))
				} else {
					failure.Add(1)
					testutil.Diff(t, http.StatusTooManyRequests, w.Result().StatusCode)
					testutil.Diff(t, `{"status":429,"statusText":"Too Many Requests"}`, string(b))
				}
			}()
		}
		wg.Wait()
		testutil.Diff(t, int32(5), success.Load())
		testutil.Diff(t, int32(5), failure.Load())
		time.Sleep(150 * time.Millisecond)

	}

}

func TestThrottleMiddleware_LeakyBucket(t *testing.T) {

	configs := []string{"./config-throttle-leaky-bucket.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "ThrottleMiddleware",
		Name:       "default",
		Namespace:  "",
	}
	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})
	h := m.Middleware(handler)

	// Run the test case 10 times to check reproducibility.
	for iter := 0; iter < 10; iter++ {

		// Send 10 requests (5 succeeded, 5 failed)
		var wg sync.WaitGroup
		success := atomic.Int32{}
		failure := atomic.Int32{}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				r := httptest.NewRequest(http.MethodGet, "http://throttle.com/leaky-bucket", nil)
				w := httptest.NewRecorder()
				h.ServeHTTP(w, r)
				b, _ := io.ReadAll(w.Result().Body)
				if w.Result().StatusCode == http.StatusOK {
					success.Add(1)
					testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
					testutil.Diff(t, "test", string(b))
				} else {
					failure.Add(1)
					testutil.Diff(t, http.StatusTooManyRequests, w.Result().StatusCode)
					testutil.Diff(t, `{"status":429,"statusText":"Too Many Requests"}`, string(b))
				}
			}()
		}
		wg.Wait()

		testutil.Diff(t, int32(4), success.Load())
		testutil.Diff(t, int32(6), failure.Load())

	}

}
