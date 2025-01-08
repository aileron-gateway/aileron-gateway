package otelmeter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

func TestRegistry(t *testing.T) {
	type condition struct {
		mp *sdkmetric.MeterProvider
	}

	type action struct {
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil registry",
			[]string{},
			[]string{},
			&condition{
				mp: nil,
			},
			&action{},
		),
		gen(
			"non nil registry",
			[]string{},
			[]string{},
			&condition{
				mp: sdkmetric.NewMeterProvider(),
			},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			metrics := &otelMeter{
				mp: tt.C().mp,
			}

			mp := metrics.MeterProvider()

			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[*sdkmetric.MeterProvider]),
			}
			testutil.Diff(t, tt.C().mp, mp, opts...)
		})
	}
}

func TestMiddleware(t *testing.T) {
	type condition struct {
		numReq int
	}

	type action struct {
		value int64
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	CndSingleRequest := tb.Condition("send single request", "send single request")
	CndMultipleRequests := tb.Condition("send multiple requests", "send multiple requests")
	ActCheckCount := tb.Action("check count", "check that an expected value returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"single request",
			[]string{CndSingleRequest},
			[]string{ActCheckCount},
			&condition{
				numReq: 1,
			},
			&action{
				value: 1,
			},
		),
		gen(
			"multiple requests",
			[]string{CndMultipleRequests},
			[]string{ActCheckCount},
			&condition{
				numReq: 5,
			},
			&action{
				value: 5,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			reader := sdkmetric.NewManualReader()
			mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
			meter := mp.Meter("test-meter")
			mAPICall, _ := meter.Int64Counter(
				"http_requests_total",
				metric.WithDescription("Total number of received http requests"),
			)

			metrics := otelMeter{
				mAPICalls: mAPICall,
			}
			h := metrics.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
			resp := httptest.NewRecorder()
			for i := 0; i < tt.C().numReq; i++ {
				h.ServeHTTP(resp, req)
			}

			rm := metricdata.ResourceMetrics{}
			err := reader.Collect(context.Background(), &rm)
			testutil.Diff(t, nil, err)

			expect := metricdata.ScopeMetrics{
				Scope: instrumentation.Scope{Name: "test-meter"},
				Metrics: []metricdata.Metrics{{
					Name:        "http_requests_total",
					Description: "Total number of received http requests",
					Unit:        "",
					Data: metricdata.Sum[int64]{
						Temporality: metricdata.CumulativeTemporality,
						IsMonotonic: true,
						DataPoints: []metricdata.DataPoint[int64]{{
							Attributes: attribute.NewSet(
								attribute.Int("code", 200),
								attribute.String("host", "test.com"),
								attribute.String("method", "GET"),
								attribute.String("path", "/test"),
							),
						}},
					},
				}},
			}
			metricdatatest.AssertEqual(t, expect, rm.ScopeMetrics[0], metricdatatest.IgnoreTimestamp(), metricdatatest.IgnoreValue())

			for _, m := range rm.ScopeMetrics[0].Metrics {
				switch a := m.Data.(type) {
				case metricdata.Sum[int64]:
					for _, dp := range a.DataPoints {
						testutil.Diff(t, dp.Value, tt.A().value)
					}
				default:
					t.Fatalf("unexpected data type %v", a)
				}
			}
		})
	}
}

func TestTripperware(t *testing.T) {
	type condition struct {
		numReq int
	}

	type action struct {
		value int64
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"single request",
			[]string{},
			[]string{},
			&condition{
				numReq: 1,
			},
			&action{
				value: 1,
			},
		),
		gen(
			"multiple requests",
			[]string{},
			[]string{},
			&condition{
				numReq: 5,
			},
			&action{
				value: 5,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			reader := sdkmetric.NewManualReader()
			mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
			meter := mp.Meter("test-meter")
			tAPICall, _ := meter.Int64Counter(
				"http_client_requests_total",
				metric.WithDescription("Total number of sent http requests"),
			)

			metrics := otelMeter{
				tAPICalls: tAPICall,
			}
			h := metrics.Tripperware(core.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: 200}, nil
			}))

			req := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
			for i := 0; i < tt.C().numReq; i++ {
				h.RoundTrip(req)
			}

			rm := metricdata.ResourceMetrics{}
			err := reader.Collect(context.Background(), &rm)
			testutil.Diff(t, nil, err)

			expect := metricdata.ScopeMetrics{
				Scope: instrumentation.Scope{Name: "test-meter"},
				Metrics: []metricdata.Metrics{{
					Name:        "http_client_requests_total",
					Description: "Total number of sent http requests",
					Unit:        "",
					Data: metricdata.Sum[int64]{
						Temporality: metricdata.CumulativeTemporality,
						IsMonotonic: true,
						DataPoints: []metricdata.DataPoint[int64]{{
							Attributes: attribute.NewSet(
								attribute.Int("code", 200),
								attribute.String("host", "test.com"),
								attribute.String("method", "GET"),
								attribute.String("path", "/test"),
							),
						}},
					},
				}},
			}
			metricdatatest.AssertEqual(t, expect, rm.ScopeMetrics[0], metricdatatest.IgnoreTimestamp(), metricdatatest.IgnoreValue())

			for _, m := range rm.ScopeMetrics[0].Metrics {
				switch a := m.Data.(type) {
				case metricdata.Sum[int64]:
					for _, dp := range a.DataPoints {
						testutil.Diff(t, dp.Value, tt.A().value)
					}
				default:
					t.Fatalf("unexpected data type %v", a)
				}
			}
		})
	}
}

func TestFinalize(t *testing.T) {
	type condition struct {
		times int
	}

	type action struct {
		err        error
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"call Finalize once",
			[]string{},
			[]string{},
			&condition{
				times: 1,
			},
			&action{
				err: nil,
			},
		),
		gen(
			"call Finalize more than once",
			[]string{},
			[]string{},
			&condition{
				times: 5,
			},
			&action{
				err:        sdkmetric.ErrReaderShutdown,
				errPattern: regexp.MustCompile(sdkmetric.ErrReaderShutdown.Error()),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			mp := sdkmetric.NewMeterProvider(
				sdkmetric.WithReader(sdkmetric.NewManualReader()),
			)

			om := &otelMeter{
				mp: mp,
			}

			var err error

			for i := 0; i < tt.C().times; i++ {
				err = om.Finalize()
			}

			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err, cmpopts.EquateErrors())
		})
	}
}
