package soaprest

// import (
// 	"encoding/xml"
// 	"errors"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"

// 	"log/slog"

// 	"github.com/aileron-gateway/aileron-gateway/app"
// 	"github.com/aileron-gateway/aileron-gateway/kernel/log"
// 	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
// 	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
// 	"github.com/google/go-cmp/cmp"
// 	"github.com/google/go-cmp/cmp/cmpopts"
// )

// // Define a ResponseWriter that fails on Write
// type erroringResponseWriter struct {
// 	*httptest.ResponseRecorder
// }

// func (w *erroringResponseWriter) Write(b []byte) (int, error) {
// 	return 0, errors.New("forced write error")
// }

// func TestSOAPErrorHandler_ServeHTTPError(t *testing.T) {
// 	type condition struct {
// 		eh  *soapErrorHandler
// 		err error
// 		req *http.Request

// 		loggingOnly         bool
// 		failWritingResponse bool
// 	}

// 	type action struct {
// 		status int
// 		header http.Header
// 		body   *soapFaultEnvelope
// 	}

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())
// 	table := tb.Build()

// 	debugLogger := log.NewJSONSLogger(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})

// 	gen := testutil.NewCase[*condition, *action]
// 	testCases := []*testutil.Case[*condition, *action]{
// 		gen(
// 			"InvalidSOAP11Request",
// 			nil,
// 			nil,
// 			&condition{
// 				eh: &soapErrorHandler{
// 					lg:          debugLogger,
// 					stackAlways: false,
// 				},
// 				err: utilhttp.NewHTTPError(
// 					app.ErrAppMiddleSOAPRESTVersionMismatch.WithoutStack(nil, nil),
// 					http.StatusForbidden,
// 				),
// 				req: httptest.NewRequest(http.MethodGet, "http://test.com/foo", nil),
// 			},
// 			&action{
// 				status: http.StatusForbidden,
// 				header: http.Header{
// 					"Content-Type":           []string{"text/xml; charset=utf-8"},
// 					"X-Content-Type-Options": []string{"nosniff"},
// 					"Vary":                   []string{"Accept"},
// 				},
// 				body: &soapFaultEnvelope{
// 					Body: &soapFaultBody{
// 						Fault: &soap11Fault{
// 							Faultcode:   faultCodeVersionMismatch,
// 							Faultstring: http.StatusText(http.StatusForbidden),
// 							Faultactor:  "test.com",
// 							Detail: &soap11FaultDetail{
// 								Message:    "Expected a SOAP 1.1 request, but received a request in a different format.",
// 								StatusCode: http.StatusForbidden,
// 							},
// 						},
// 					},
// 				},
// 			},
// 		),
// 		gen(
// 			"BadRequest",
// 			nil,
// 			nil,
// 			&condition{
// 				eh: &soapErrorHandler{
// 					lg:          debugLogger,
// 					stackAlways: false,
// 				},
// 				err: utilhttp.ErrBadRequest,
// 				req: httptest.NewRequest(http.MethodGet, "http://test.com/foo", nil),
// 			},
// 			&action{
// 				status: http.StatusBadRequest,
// 				header: http.Header{
// 					"Content-Type":           []string{"text/xml; charset=utf-8"},
// 					"X-Content-Type-Options": []string{"nosniff"},
// 					"Vary":                   []string{"Accept"},
// 				},
// 				body: &soapFaultEnvelope{
// 					Body: &soapFaultBody{
// 						Fault: &soap11Fault{
// 							Faultcode:   faultCodeClient,
// 							Faultstring: http.StatusText(http.StatusBadRequest),
// 							Faultactor:  "test.com",
// 							Detail: &soap11FaultDetail{
// 								Message:    "An error has occurred while processing the request from the client.",
// 								StatusCode: http.StatusBadRequest,
// 							},
// 						},
// 					},
// 				},
// 			},
// 		),
// 		gen(
// 			"InternalServerError",
// 			nil,
// 			nil,
// 			&condition{
// 				eh: &soapErrorHandler{
// 					lg:          debugLogger,
// 					stackAlways: false,
// 				},
// 				err: utilhttp.ErrInternalServerError,
// 				req: httptest.NewRequest(http.MethodGet, "http://test.com/foo", nil),
// 			},
// 			&action{
// 				status: http.StatusInternalServerError,
// 				header: http.Header{
// 					"Content-Type":           []string{"text/xml; charset=utf-8"},
// 					"X-Content-Type-Options": []string{"nosniff"},
// 					"Vary":                   []string{"Accept"},
// 				},
// 				body: &soapFaultEnvelope{
// 					Body: &soapFaultBody{
// 						Fault: &soap11Fault{
// 							Faultcode:   faultCodeServer,
// 							Faultstring: http.StatusText(http.StatusInternalServerError),
// 							Faultactor:  "test.com",
// 							Detail: &soap11FaultDetail{
// 								Message:    "An error has occurred on the upstream server.",
// 								StatusCode: http.StatusInternalServerError,
// 							},
// 						},
// 					},
// 				},
// 			},
// 		),
// 		gen(
// 			"Logging only status -1",
// 			nil,
// 			nil,
// 			&condition{
// 				eh: &soapErrorHandler{
// 					lg:          debugLogger,
// 					stackAlways: false,
// 				},
// 				err: utilhttp.NewHTTPError(errors.New("test error"), -1),
// 				req: httptest.NewRequest(http.MethodGet, "http://test.com/foo", nil),

// 				loggingOnly: true,
// 			},
// 			&action{
// 				status: http.StatusOK,
// 				header: nil,
// 				body:   nil,
// 			},
// 		),
// 		gen(
// 			"Failed to write response body",
// 			nil,
// 			nil,
// 			&condition{
// 				eh: &soapErrorHandler{
// 					lg:          debugLogger,
// 					stackAlways: false,
// 				},
// 				err: utilhttp.ErrInternalServerError,
// 				req: httptest.NewRequest(http.MethodGet, "http://test.com/foo", nil),

// 				loggingOnly:         true, // Since the write operation to the body will fail, the body check will not be performed.
// 				failWritingResponse: true,
// 			},
// 			&action{
// 				status: http.StatusInternalServerError,
// 				header: nil,
// 				body:   nil,
// 			},
// 		),
// 	}

// 	testutil.Register(table, testCases...)

// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {
// 			rr := httptest.NewRecorder()
// 			var w http.ResponseWriter

// 			if tt.C().failWritingResponse {
// 				w = &erroringResponseWriter{
// 					ResponseRecorder: rr,
// 				}
// 			} else {
// 				w = rr
// 			}

// 			tt.C().eh.ServeHTTPError(w, tt.C().req, tt.C().err)
// 			resp := rr.Result()
// 			defer resp.Body.Close()
// 			b, _ := io.ReadAll(resp.Body)

// 			testutil.Diff(t, tt.A().status, resp.StatusCode)

// 			for k, v := range tt.A().header {
// 				testutil.Diff(t, v, resp.Header[k])
// 			}

// 			if !tt.C().loggingOnly {
// 				var gotBody soapFaultEnvelope
// 				if err := xml.Unmarshal(b, &gotBody); err != nil {
// 					t.Fatalf("failed to unmarshal response body: %v", err)
// 				}

// 				opts := []cmp.Option{
// 					cmpopts.IgnoreFields(soapFaultEnvelope{}, "XMLName"),
// 					cmpopts.IgnoreFields(soapFaultBody{}),
// 					cmpopts.IgnoreFields(soap11Fault{}, "XMLName"),
// 					cmpopts.IgnoreFields(soap11FaultDetail{}, "XMLName"),
// 				}

// 				testutil.Diff(t, tt.A().body, &gotBody, opts...)
// 			}
// 		})
// 	}
// }
