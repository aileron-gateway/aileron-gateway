// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package prommeter

import (
	"net/http"
	"strconv"

	"github.com/aileron-gateway/aileron-gateway/core"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	// HandlerBase is the base struct for
	// http.Handler type resource.
	// This provides Patterns() and Methods() methods
	// to fulfill the core.Handler interface.
	*utilhttp.HandlerBase

	http.Handler

	reg *prometheus.Registry

	// middleware metrics
	mAPICalls *prometheus.CounterVec

	// tripperware metrics
	tAPICalls *prometheus.CounterVec
}

// Registry return the prometheus registry.
// The registry can be used to register custom metrics.
func (m *metrics) Registry() *prometheus.Registry {
	return m.reg
}

func (m *metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := utilhttp.WrapWriter(w)
		w = ww

		defer func() {
			m.mAPICalls.With(prometheus.Labels{
				"host":   r.Host,
				"path":   r.URL.Path,
				"code":   strconv.Itoa(ww.StatusCode()),
				"method": r.Method,
			}).Inc()
		}()

		next.ServeHTTP(w, r)
	})
}

func (m *metrics) Tripperware(next http.RoundTripper) http.RoundTripper {
	return core.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		var status int
		defer func() {
			m.tAPICalls.With(prometheus.Labels{
				"host":   r.URL.Host,
				"path":   r.URL.Path,
				"code":   strconv.Itoa(status),
				"method": r.Method,
			}).Inc()
		}()

		resp, err := next.RoundTrip(r)
		if resp != nil {
			status = resp.StatusCode
		}

		return resp, err
	})
}
