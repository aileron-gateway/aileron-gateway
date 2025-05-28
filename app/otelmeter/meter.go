// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package otelmeter

import (
	"context"
	"net/http"

	"github.com/aileron-gateway/aileron-gateway/core"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type otelMeter struct {
	mp *sdkmetric.MeterProvider
	// middleware metrics
	mAPICalls metric.Int64Counter
	// tripperware metrics
	tAPICalls metric.Int64Counter
}

// MeterProvider return the opentelemetry metric provider.
// The registry can be used to register custom metrics.
func (m *otelMeter) MeterProvider() *sdkmetric.MeterProvider {
	return m.mp
}

func (m *otelMeter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := utilhttp.WrapWriter(w)
		w = ww

		defer func(ctx context.Context) {
			m.mAPICalls.Add(ctx, 1,
				metric.WithAttributes(
					attribute.String("host", r.Host),
					attribute.String("path", r.URL.Path),
					attribute.Int("code", ww.StatusCode()),
					attribute.String("method", r.Method),
				),
			)
		}(r.Context())

		next.ServeHTTP(w, r)
	})
}

func (m *otelMeter) Tripperware(next http.RoundTripper) http.RoundTripper {
	return core.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		var status int
		defer func() {
			m.tAPICalls.Add(r.Context(), 1,
				metric.WithAttributes(
					attribute.String("host", r.Host),
					attribute.String("path", r.URL.Path),
					attribute.Int("code", status),
					attribute.String("method", r.Method),
				),
			)
		}()

		resp, err := next.RoundTrip(r)
		if resp != nil {
			status = resp.StatusCode
		}

		return resp, err
	})
}

// Finalize ensures m.mp.Shutdown is called before the application halts.
// Shutdown flushes metrics data remaining in the buffer.
// This implements the core.Finalizer interface.
func (m *otelMeter) Finalize() error {
	return m.mp.Shutdown(context.Background())
}
