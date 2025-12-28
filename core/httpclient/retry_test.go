// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpclient

import (
	"bytes"
	stdcmp "cmp"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testBody struct {
	io.ReadCloser
	readErr  error
	closeErr error
}

func (b *testBody) Read(p []byte) (n int, err error) {
	if b.readErr != nil {
		return 0, b.readErr
	}
	return len(p), nil
}

func (b *testBody) Close() error {
	return b.closeErr
}

func TestRetry_Tripperware(t *testing.T) {
	type condition struct {
		retry         *retry
		contentLength int
		status        []int
		cancelCtx     bool // Cancel context in round tripper.
		doneCtx       bool // Simulate client's cancel
		body          io.ReadCloser
		bodyFunc      func() (io.ReadCloser, error)
	}

	type action struct {
		called int
		status int
		err    error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"0 retry/without error",
			&condition{
				retry: &retry{
					maxRetry: 0,
				},
				status: []int{http.StatusNotFound, http.StatusOK},
				body:   io.NopCloser(bytes.NewReader([]byte("test body"))),
			},
			&action{
				called: 1,
				status: http.StatusNotFound,
			},
		),
		gen(
			"0 retry/with error",
			&condition{
				retry: &retry{
					maxRetry: 0,
				},
				status: []int{0, http.StatusOK}, // 0 will raise an error.
				body:   io.NopCloser(bytes.NewReader([]byte("test body"))),
			},
			&action{
				called: 1,
				err:    io.EOF,
			},
		),
		gen(
			"negative content length/without error",
			&condition{
				retry: &retry{
					maxRetry: 1,
				},
				contentLength: -1, // Overwrite the requests's content length.
				status:        []int{http.StatusNotFound, http.StatusOK},
				body:          io.NopCloser(bytes.NewReader([]byte("test body"))),
			},
			&action{
				called: 1,
				status: http.StatusNotFound,
			},
		),
		gen(
			"negative content length/with error",
			&condition{
				retry: &retry{
					maxRetry: 1,
				},
				contentLength: -1,                      // Overwrite the requests's content length.
				status:        []int{0, http.StatusOK}, // 0 will raise an error.
				body:          io.NopCloser(bytes.NewReader([]byte("test body"))),
			},
			&action{
				called: 1,
				err:    io.EOF,
			},
		),
		gen(
			"large content length/without error",
			&condition{
				retry: &retry{
					maxRetry:         1,
					maxContentLength: 100,
				},
				contentLength: 999, // Overwrite the requests's content length.
				status:        []int{http.StatusNotFound, http.StatusOK},
				body:          io.NopCloser(bytes.NewReader([]byte("test body"))),
			},
			&action{
				called: 1,
				status: http.StatusNotFound,
			},
		),
		gen(
			"large content length/with error",
			&condition{
				retry: &retry{
					maxRetry:         1,
					maxContentLength: 100,
				},
				contentLength: 999,                     // Overwrite the requests's content length.
				status:        []int{0, http.StatusOK}, // 0 will raise an error.
				body:          io.NopCloser(bytes.NewReader([]byte("test body"))),
			},
			&action{
				called: 1,
				err:    io.EOF,
			},
		),
		gen(
			"1 retry/without error",
			&condition{
				retry: &retry{
					maxRetry:         1,
					maxContentLength: 100,
				},
				status: []int{0, http.StatusFound},
				body:   io.NopCloser(bytes.NewReader([]byte("test body"))),
			},
			&action{
				called: 2,
				status: http.StatusFound,
			},
		),
		gen(
			"1 retry/with error",
			&condition{
				retry: &retry{
					maxRetry:         1,
					maxContentLength: 100,
				},
				status: []int{0, 0}, // 0 will raise an error.
				body:   io.NopCloser(bytes.NewReader([]byte("test body"))),
			},
			&action{
				called: 2,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeRetry,
					Description: ErrDescRetryFail,
				},
			},
		),
		gen(
			"2 retry/without error",
			&condition{
				retry: &retry{
					maxRetry:         2,
					maxContentLength: 100,
				},
				status: []int{0, 0, http.StatusFound},
				body:   io.NopCloser(bytes.NewReader([]byte("test body"))),
			},
			&action{
				called: 3,
				status: http.StatusFound,
			},
		),
		gen(
			"2 retry/with error",
			&condition{
				retry: &retry{
					maxRetry:         2,
					maxContentLength: 100,
				},
				status: []int{0, 0, 0}, // 0 will raise an error.
				body:   io.NopCloser(bytes.NewReader([]byte("test body"))),
			},
			&action{
				called: 3,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeRetry,
					Description: ErrDescRetryFail,
				},
			},
		),
		gen(
			"retry http status/without error",
			&condition{
				retry: &retry{
					maxRetry:         3,
					maxContentLength: 100,
					retryStatus:      []int{http.StatusInternalServerError, http.StatusGatewayTimeout},
				},
				status: []int{0, http.StatusInternalServerError, http.StatusGatewayTimeout, http.StatusFound},
				body:   io.NopCloser(bytes.NewReader([]byte("test body"))),
			},
			&action{
				called: 4,
				status: http.StatusFound,
			},
		),
		gen(
			"retry http status/with error",
			&condition{
				retry: &retry{
					maxRetry:         3,
					maxContentLength: 100,
					retryStatus:      []int{http.StatusInternalServerError, http.StatusGatewayTimeout},
				},
				status: []int{0, http.StatusInternalServerError, http.StatusGatewayTimeout, http.StatusInternalServerError},
				body:   io.NopCloser(bytes.NewReader([]byte("test body"))),
			},
			&action{
				called: 4,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeRetry,
					Description: ErrDescRetryFail,
				},
			},
		),
		gen(
			"setupRewindBody error",
			&condition{
				retry: &retry{
					maxRetry:         1,
					maxContentLength: 100,
				},
				status: []int{0, http.StatusFound},
				body:   io.NopCloser(&testutil.ErrorReader{}),
			},
			&action{
				called: 0,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeRetry,
					Description: ErrDescRetryFail,
				},
			},
		),
		gen(
			"rewindBody error",
			&condition{
				retry: &retry{
					maxRetry:         1,
					maxContentLength: 100,
				},
				status: []int{0, http.StatusFound},
				body:   io.NopCloser(bytes.NewReader([]byte("test body"))),
				bodyFunc: func() (io.ReadCloser, error) {
					return io.NopCloser(bytes.NewReader([]byte("test body"))), io.ErrUnexpectedEOF
				},
			},
			&action{
				called: 1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeRetry,
					Description: ErrDescRetryFail,
				},
			},
		),
		gen(
			"context cancel error",
			&condition{
				retry: &retry{
					maxRetry:         1,
					maxContentLength: 100,
				},
				status:    []int{0, http.StatusFound},
				cancelCtx: true,
				body:      io.NopCloser(bytes.NewReader([]byte("test body"))),
			},
			&action{
				called: 1,
				err:    context.Canceled,
			},
		),
		gen(
			"context done",
			&condition{
				retry: &retry{
					maxRetry:         1,
					maxContentLength: 100,
				},
				status:  []int{0, http.StatusFound},
				body:    io.NopCloser(bytes.NewReader([]byte("test body"))),
				doneCtx: true,
			},
			&action{
				called: 1,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeRetry,
					Description: ErrDescRetryFail,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			// Test request.
			ctx, cancel := context.WithCancel(context.Background())
			req := &http.Request{
				ContentLength: int64(stdcmp.Or(tt.C.contentLength, 9)),
				Body:          tt.C.body,
				GetBody:       tt.C.bodyFunc,
			}
			req = req.WithContext(ctx)

			// Test roundtripper
			called := 0
			rt := tt.C.retry.Tripperware(core.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
				called += 1

				var body []byte
				var err error
				if r.Body != nil {
					body, err = io.ReadAll(r.Body)
				}
				testutil.Diff(t, nil, err)
				testutil.Diff(t, "test body", string(body))

				if tt.C.doneCtx {
					cancel() // Simulate client's cancel while waiting.
				}
				if tt.C.cancelCtx {
					return nil, context.Canceled
				}

				status := 0
				if len(tt.C.status) >= called {
					status = tt.C.status[called-1]
				}
				if status == 0 {
					return nil, io.EOF
				}

				return &http.Response{
					StatusCode: status,
				}, nil
			}))

			res, err := rt.RoundTrip(req)

			testutil.Diff(t, tt.A.called, called)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			if tt.A.status > 0 {
				testutil.Diff(t, tt.A.status, res.StatusCode)
			}
		})
	}
}

func TestSetupRewindBody(t *testing.T) {
	type condition struct {
		req   *http.Request
		read  int
		close bool
	}

	type action struct {
		req  *http.Request
		body string
		err  error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil body",
			&condition{
				req: &http.Request{
					Body: nil,
				},
			},
			&action{
				req: &http.Request{
					Body: nil,
				},
			},
		),
		gen(
			"NoBody body",
			&condition{
				req: &http.Request{
					Body: http.NoBody,
				},
			},
			&action{
				req: &http.Request{
					Body: http.NoBody,
				},
			},
		),
		gen(
			"trackingBody/not read/not closed",
			&condition{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewReader([]byte("test body"))),
				},
				read:  0,
				close: false,
			},
			&action{
				req: &http.Request{
					Body: &readTrackingBody{
						ReadCloser: io.NopCloser(bytes.NewReader([]byte("test body"))),
					},
					GetBody: func() (io.ReadCloser, error) {
						return io.NopCloser(bytes.NewReader([]byte("test body"))), nil
					},
				},
				body: "test body",
			},
		),
		gen(
			"trackingBody/read/not closed",
			&condition{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewReader([]byte("test body"))),
				},
				read:  4,
				close: false,
			},
			&action{
				req: &http.Request{
					Body: &readTrackingBody{
						ReadCloser: io.NopCloser(bytes.NewReader([]byte("test body"))),
					},
					GetBody: func() (io.ReadCloser, error) {
						return io.NopCloser(bytes.NewReader([]byte("test body"))), nil
					},
				},
				body: "test body",
			},
		),
		gen(
			"trackingBody/not read/closed",
			&condition{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewReader([]byte("test body"))),
				},
				read:  0,
				close: true,
			},
			&action{
				req: &http.Request{
					Body: &readTrackingBody{
						ReadCloser: io.NopCloser(bytes.NewReader([]byte("test body"))),
					},
					GetBody: func() (io.ReadCloser, error) {
						return io.NopCloser(bytes.NewReader([]byte("test body"))), nil
					},
				},
				body: "test body",
			},
		),
		gen(
			"trackingBody/read/closed",
			&condition{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewReader([]byte("test body"))),
				},
				read:  4,
				close: true,
			},
			&action{
				req: &http.Request{
					Body: &readTrackingBody{
						ReadCloser: io.NopCloser(bytes.NewReader([]byte("test body"))),
					},
					GetBody: func() (io.ReadCloser, error) {
						return io.NopCloser(bytes.NewReader([]byte("test body"))), nil
					},
				},
				body: "test body",
			},
		),
		gen(
			"with GetBody",
			&condition{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewReader([]byte("test body"))),
					GetBody: func() (io.ReadCloser, error) {
						return io.NopCloser(bytes.NewReader([]byte("rewound body"))), nil
					},
				},
				read:  4,
				close: true,
			},
			&action{
				req: &http.Request{
					Body: &readTrackingBody{
						ReadCloser: io.NopCloser(bytes.NewReader([]byte("test body"))),
					},
					GetBody: func() (io.ReadCloser, error) {
						return io.NopCloser(bytes.NewReader([]byte("rewound body"))), nil
					},
				},
				body: "rewound body",
			},
		),
		gen(
			"with error body",
			&condition{
				req: &http.Request{
					Body:    &testBody{readErr: io.ErrUnexpectedEOF},
					GetBody: nil,
				},
			},
			&action{
				req: nil,
				err: io.ErrUnexpectedEOF,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			newReq, err := setupRewindBody(tt.C.req)

			if tt.A.err != nil {
				testutil.Diff(t, tt.A.err.Error(), err.Error())
				testutil.Diff(t, (*http.Request)(nil), newReq)
				return
			}
			testutil.Diff(t, nil, err)

			opts := []cmp.Option{
				cmp.AllowUnexported(readTrackingBody{}),
				cmp.AllowUnexported(bytes.Reader{}),
				cmp.AllowUnexported(http.Request{}),
				cmp.Comparer(func(x, y func() (io.ReadCloser, error)) bool {
					if x == nil || y == nil {
						return x == nil && y == nil
					}
					xrc, _ := x()
					yrc, _ := y()
					if xrc == nil || yrc == nil {
						return xrc == yrc
					}
					bx, _ := io.ReadAll(xrc)
					by, _ := io.ReadAll(yrc)
					return bytes.Equal(bx, by)
				}),
			}
			testutil.Diff(t, tt.A.req, newReq, opts...)

			if tt.C.read > 0 {
				b := make([]byte, tt.C.read)
				newReq.Body.Read(b)
			}
			if tt.C.close {
				newReq.Body.Close()
			}

			if newReq.Body != nil {
				newReq, _ = rewindBody(newReq)
				b, _ := io.ReadAll(newReq.Body)
				testutil.Diff(t, tt.A.body, string(b))
			}
		})
	}
}

func TestRewindBody(t *testing.T) {
	type condition struct {
		body    io.ReadCloser
		getBody func() (io.ReadCloser, error)
	}

	type action struct {
		body io.ReadCloser
		err  error
	}

	testGetBody := func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader([]byte("rewound body"))), nil
	}
	testErrGetBody := func() (io.ReadCloser, error) {
		return nil, io.ErrUnexpectedEOF
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil body",
			&condition{
				body: nil,
			},
			&action{
				body: nil,
			},
		),
		gen(
			"NoBody body",
			&condition{
				body: http.NoBody,
			},
			&action{
				body: http.NoBody,
			},
		),
		gen(
			"trackingBody/not read/not closed",
			&condition{
				body: &readTrackingBody{
					ReadCloser: io.NopCloser(bytes.NewReader([]byte("test body"))),
					didRead:    false,
					didClose:   false,
				},
			},
			&action{
				body: &readTrackingBody{
					ReadCloser: io.NopCloser(bytes.NewReader([]byte("test body"))),
				},
			},
		),
		gen(
			"trackingBody/not read/closed",
			&condition{
				body: &readTrackingBody{
					ReadCloser: io.NopCloser(bytes.NewReader([]byte("test body"))),
					didRead:    false,
					didClose:   true,
				},
				getBody: testGetBody,
			},
			&action{
				body: &readTrackingBody{
					ReadCloser: io.NopCloser(bytes.NewReader([]byte("rewound body"))),
				},
			},
		),
		gen(
			"trackingBody/read/not closed",
			&condition{
				body: &readTrackingBody{
					ReadCloser: io.NopCloser(bytes.NewReader([]byte("test body"))),
					didRead:    true,
					didClose:   false,
				},
				getBody: testGetBody,
			},
			&action{
				body: &readTrackingBody{
					ReadCloser: io.NopCloser(bytes.NewReader([]byte("rewound body"))),
				},
			},
		),
		gen(
			"trackingBody/read/closed",
			&condition{
				body: &readTrackingBody{
					ReadCloser: io.NopCloser(bytes.NewReader([]byte("test body"))),
					didRead:    true,
					didClose:   true,
				},
				getBody: testGetBody,
			},
			&action{
				body: &readTrackingBody{
					ReadCloser: io.NopCloser(bytes.NewReader([]byte("rewound body"))),
				},
			},
		),
		gen(
			"cannot get body",
			&condition{
				body:    io.NopCloser(bytes.NewReader([]byte("test body"))),
				getBody: nil,
			},
			&action{
				body: nil,
				err:  errors.New("cannot rewind because the GetBody is nil"),
			},
		),
		gen(
			"get body returns error",
			&condition{
				body:    io.NopCloser(bytes.NewReader([]byte("test body"))),
				getBody: testErrGetBody,
			},
			&action{
				body: nil,
				err:  io.ErrUnexpectedEOF,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			req := &http.Request{
				Body:    tt.C.body,
				GetBody: tt.C.getBody,
			}

			newReq, err := rewindBody(req)
			if tt.A.err != nil {
				testutil.Diff(t, tt.A.err.Error(), err.Error())
				return
			}
			testutil.Diff(t, nil, err)

			opts := []cmp.Option{
				cmp.AllowUnexported(readTrackingBody{}),
				cmp.AllowUnexported(bytes.Reader{}),
			}
			fmt.Printf("%#v\n", newReq)
			testutil.Diff(t, tt.A.body, newReq.Body, opts...)
		})
	}
}
