package core_test

import (
	"fmt"
	"net/http"

	"github.com/aileron-gateway/aileron-gateway/core"
)

type foo struct {
}

func (f *foo) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("foo", "bar")
		next.ServeHTTP(w, r)
	})
}

func ExampleMiddleware() {
	middleware := &foo{}
	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	r, _ := http.NewRequest(http.MethodGet, "http://ecample.com", nil)
	handler.ServeHTTP(nil, r) // Omit ResponseWriter in this example.

	fmt.Println(r.Header.Get("foo"))
	// Output:
	// bar
}

func addHeader(r *http.Request) (*http.Response, error) {
	r.Header.Add("foo", "bar")
	return nil, nil // Return nil for example.
	// return http.DefaultTransport.RoundTrip(r)
}

func ExampleRoundTripperFunc() {
	var rt http.RoundTripper = core.RoundTripperFunc(addHeader)

	r, _ := http.NewRequest(http.MethodGet, "http://ecample.com", nil)
	_, _ = rt.RoundTrip(r)

	fmt.Println(r.Header.Get("foo"))
	// Output:
	// bar
}
