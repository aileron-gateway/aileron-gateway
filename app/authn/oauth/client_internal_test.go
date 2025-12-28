// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package oauth

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/golang-jwt/jwt/v5"
)

func TestTokenIntrospectionClient_tokenIntrospection(t *testing.T) {
	type condition struct {
		client *tokenIntrospectionClient
		ctx    context.Context
		query  map[string]string
	}

	type action struct {
		status     int
		claims     jwt.MapClaims
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
		errStatus  int
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success",
			&condition{
				client: &tokenIntrospectionClient{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					requester: &testRequester{
						status: http.StatusOK,
						body:   []byte(`{"access_token":"test_token"}`),
					},
					provider: &provider{
						tokenEP: "http://test.com/token",
					},
				},
				ctx: context.Background(),
			},
			&action{
				status: http.StatusOK,
				claims: jwt.MapClaims{
					"access_token": "test_token",
				},
				err: nil,
			},
		),
		gen(
			"doRequest error",
			&condition{
				client: &tokenIntrospectionClient{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					requester: &testRequester{
						status: -1,
						err:    errors.New("test error"),
					},
					provider: &provider{
						tokenEP: "http://test.com/token",
					},
				},
				ctx: context.Background(),
			},
			&action{
				status:     -1,
				claims:     nil,
				err:        app.ErrAppAuthnIntrospection,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to token introspection.`),
				errStatus:  http.StatusInternalServerError,
			},
		),
		gen(
			"server error",
			&condition{
				client: &tokenIntrospectionClient{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					requester: &testRequester{
						status: http.StatusInternalServerError,
						body:   []byte(`{"error":"test_error"}`),
					},
					provider: &provider{
						tokenEP: "http://test.com/token",
					},
				},
				ctx: context.Background(),
			},
			&action{
				status:     http.StatusInternalServerError,
				claims:     nil,
				err:        app.ErrAppAuthnIntrospection,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to token introspection. status:500 body:`),
				errStatus:  http.StatusInternalServerError,
			},
		),
		gen(
			"authn error",
			&condition{
				client: &tokenIntrospectionClient{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					requester: &testRequester{
						status: http.StatusUnauthorized,
						body:   []byte(`{"error":"test_error"}`),
					},
					provider: &provider{
						tokenEP: "http://test.com/token",
					},
				},
				ctx: context.Background(),
			},
			&action{
				status:     http.StatusUnauthorized,
				claims:     nil,
				err:        app.ErrAppAuthnIntrospection,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to token introspection. status:401 body:`),
				errStatus:  http.StatusUnauthorized,
			},
		),
		gen(
			"unmarshal error",
			&condition{
				client: &tokenIntrospectionClient{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					requester: &testRequester{
						status: http.StatusOK,
						body:   []byte(`plain text`),
					},
					provider: &provider{
						tokenEP: "http://test.com/token",
					},
				},
				ctx: context.Background(),
			},
			&action{
				status:     http.StatusOK,
				claims:     nil,
				err:        app.ErrAppGenUnmarshal,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to unmarshal from=json to=MapClaims`),
				errStatus:  http.StatusInternalServerError,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			status, claims, err := tt.C.client.tokenIntrospection(tt.C.ctx, tt.C.query)
			testutil.Diff(t, tt.A.status, status)
			testutil.Diff(t, tt.A.claims, claims)
			if tt.A.err != nil {
				e := err.(*utilhttp.HTTPError)
				testutil.DiffError(t, tt.A.err, tt.A.errPattern, e.Unwrap())
				testutil.Diff(t, tt.A.errStatus, e.StatusCode())
			}
		})
	}
}

func TestRedeemTokenClient_redeemToken(t *testing.T) {
	type condition struct {
		client *redeemTokenClient
		ctx    context.Context
		query  map[string]string
	}

	type action struct {
		status     int
		resp       *TokenResponse
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
		errStatus  int
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success",
			&condition{
				client: &redeemTokenClient{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					requester: &testRequester{
						status: http.StatusOK,
						body:   []byte(`{"access_token":"test_token"}`),
					},
					provider: &provider{
						tokenEP: "http://test.com/token",
					},
				},
				ctx: context.Background(),
			},
			&action{
				status: http.StatusOK,
				resp: &TokenResponse{
					AccessToken: "test_token",
					StatusCode:  http.StatusOK,
					RawBody:     []byte(`{"access_token":"test_token"}`),
				},
				err: nil,
			},
		),
		gen(
			"doRequest error",
			&condition{
				client: &redeemTokenClient{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					requester: &testRequester{
						status: -1,
						err:    errors.New("test error"),
					},
					provider: &provider{
						tokenEP: "http://test.com/token",
					},
				},
				ctx: context.Background(),
			},
			&action{
				status:     -1,
				resp:       nil,
				err:        app.ErrAppAuthnRedeemToken,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to redeem token`),
				errStatus:  http.StatusInternalServerError,
			},
		),
		gen(
			"server error",
			&condition{
				client: &redeemTokenClient{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					requester: &testRequester{
						status: http.StatusInternalServerError,
						body:   []byte(`{"error":"test_error"}`),
					},
					provider: &provider{
						tokenEP: "http://test.com/token",
					},
				},
				ctx: context.Background(),
			},
			&action{
				status:     http.StatusInternalServerError,
				resp:       nil,
				err:        app.ErrAppAuthnRedeemToken,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to redeem token. status:500 body:`),
				errStatus:  http.StatusInternalServerError,
			},
		),
		gen(
			"authn error",
			&condition{
				client: &redeemTokenClient{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					requester: &testRequester{
						status: http.StatusUnauthorized,
						body:   []byte(`{"error":"test_error"}`),
					},
					provider: &provider{
						tokenEP: "http://test.com/token",
					},
				},
				ctx: context.Background(),
			},
			&action{
				status:     http.StatusUnauthorized,
				resp:       nil,
				err:        app.ErrAppAuthnRedeemToken,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to redeem token. status:401 body:`),
				errStatus:  http.StatusUnauthorized,
			},
		),
		gen(
			"unmarshal error",
			&condition{
				client: &redeemTokenClient{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					requester: &testRequester{
						status: http.StatusOK,
						body:   []byte(`plain text`),
					},
					provider: &provider{
						tokenEP: "http://test.com/token",
					},
				},
				ctx: context.Background(),
			},
			&action{
				status:     http.StatusOK,
				resp:       nil,
				err:        app.ErrAppGenUnmarshal,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to unmarshal from=json to=TokenResponseModel`),
				errStatus:  http.StatusInternalServerError,
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			status, resp, err := tt.C.client.redeemToken(tt.C.ctx, tt.C.query)
			testutil.Diff(t, tt.A.status, status)
			testutil.Diff(t, tt.A.resp, resp)
			if tt.A.err != nil {
				e := err.(*utilhttp.HTTPError)
				testutil.DiffError(t, tt.A.err, tt.A.errPattern, e.Unwrap())
				testutil.Diff(t, tt.A.errStatus, e.StatusCode())
			}
		})
	}
}

func TestClientRequester_doRequest(t *testing.T) {
	type condition struct {
		req   *clientRequester
		ctx   context.Context
		query map[string]string
	}

	type action struct {
		authHeader string
		bodyQuery  string

		status     int
		body       []byte
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success / basic auth",
			&condition{
				req: &clientRequester{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					rt: &testClient{
						status: http.StatusOK,
						body:   bytes.NewReader([]byte("test")),
						header: http.Header{"Content-Type": []string{"application/json"}},
					},
					client: &client{
						id:     "test-id",
						secret: "test-secret",
					},
				},
				ctx: context.Background(),
			},
			&action{
				authHeader: "Basic dGVzdC1pZDp0ZXN0LXNlY3JldA==",
				status:     http.StatusOK,
				body:       []byte("test"),
				err:        nil,
			},
		),
		gen(
			"success / form auth",
			&condition{
				req: &clientRequester{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					rt: &testClient{
						status: http.StatusOK,
						body:   bytes.NewReader([]byte("test")),
						header: http.Header{"Content-Type": []string{"application/json"}},
					},
					client: &client{
						id:     "test-id",
						secret: "test-secret",
					},
					clientAuthMethod: clientAuthForm,
				},
				ctx: context.Background(),
			},
			&action{
				bodyQuery: "client_id=test-id&client_secret=test-secret",
				status:    http.StatusOK,
				body:      []byte("test"),
				err:       nil,
			},
		),
		gen(
			"with query param",
			&condition{
				req: &clientRequester{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					rt: &testClient{
						status: http.StatusOK,
						body:   bytes.NewReader([]byte("test")),
						header: http.Header{"Content-Type": []string{"application/json"}},
					},
					client: &client{
						id:     "test-id",
						secret: "test-secret",
					},
				},
				ctx:   context.Background(),
				query: map[string]string{"foo": "bar"},
			},
			&action{
				authHeader: "Basic dGVzdC1pZDp0ZXN0LXNlY3JldA==",
				bodyQuery:  "foo=bar",
				status:     http.StatusOK,
				body:       []byte("test"),
				err:        nil,
			},
		),
		gen(
			"request create error",
			&condition{
				req: &clientRequester{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					rt: &testClient{
						status: http.StatusOK,
						body:   bytes.NewReader([]byte("test")),
						header: http.Header{"Content-Type": []string{"application/json"}},
					},
					client: &client{
						id:     "test-id",
						secret: "test-secret",
					},
				},
				ctx: nil,
			},
			&action{
				status:     -1,
				body:       nil,
				err:        app.ErrAppGenCreateRequest,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create http request. method=POST url=http://test.com body=`),
			},
		),
		gen(
			"round trip error",
			&condition{
				req: &clientRequester{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					rt: &testClient{
						err: errors.New("test error"),
					},
					client: &client{
						id:     "test-id",
						secret: "test-secret",
					},
				},
				ctx: context.Background(),
			},
			&action{
				authHeader: "Basic dGVzdC1pZDp0ZXN0LXNlY3JldA==",
				status:     -1,
				body:       nil,
				err:        app.ErrAppGenRoundTrip,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to round trip. method=POST url=http://test.com body=`),
			},
		),
		gen(
			"invalid content type",
			&condition{
				req: &clientRequester{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					rt: &testClient{
						status: http.StatusOK,
						body:   bytes.NewReader([]byte("test")),
						header: http.Header{"Content-Type": []string{"text/plain"}},
					},
					client: &client{
						id:     "test-id",
						secret: "test-secret",
					},
				},
				ctx: context.Background(),
			},
			&action{
				authHeader: "Basic dGVzdC1pZDp0ZXN0LXNlY3JldA==",
				status:     http.StatusOK,
				body:       nil,
				err:        app.ErrAppGenInvalidResponse,
				errPattern: regexp.MustCompile(core.ErrPrefix + `invalid response. method=POST url=http://test.com Content-Type=text/plain`),
			},
		),
		gen(
			"read body error",
			&condition{
				req: &clientRequester{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					rt: &testClient{
						status: http.StatusOK,
						body:   &testutil.ErrorReader{},
						header: http.Header{"Content-Type": []string{"application/json"}},
					},
					client: &client{
						id:     "test-id",
						secret: "test-secret",
					},
				},
				ctx: context.Background(),
			},
			&action{
				authHeader: "Basic dGVzdC1pZDp0ZXN0LXNlY3JldA==",
				status:     http.StatusOK,
				body:       nil,
				err:        app.ErrAppGenReadHTTPBody,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to read response body. read=`),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			status, body, err := tt.C.req.doRequest(tt.C.ctx, "http://test.com", tt.C.query)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)
			testutil.Diff(t, tt.A.status, status)
			testutil.Diff(t, tt.A.body, body)

			c := tt.C.req.rt.(*testClient)
			if c.got != nil {
				b, _ := io.ReadAll(c.got.Body)
				testutil.Diff(t, "application/json", c.got.Header.Get("Accept"))
				testutil.Diff(t, "application/x-www-form-urlencoded", c.got.Header.Get("Content-Type"))
				testutil.Diff(t, tt.A.authHeader, c.got.Header.Get("Authorization"))
				testutil.Diff(t, tt.A.bodyQuery, string(b))
			}
		})
	}
}

func TestMapValue(t *testing.T) {
	type condition struct {
		claims jwt.MapClaims
		key    string
		f      func(jwt.MapClaims, string) any
	}

	type action struct {
		result any
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"key exists for int",
			&condition{
				claims: jwt.MapClaims{
					"test": 123,
				},
				key: "test",
				f: func(c jwt.MapClaims, key string) any {
					return mapValue[int](c, key)
				},
			},
			&action{
				result: 123,
			},
		),
		gen(
			"key exists for map",
			&condition{
				claims: jwt.MapClaims{
					"test": map[string]string{"foo": "bar"},
				},
				key: "test",
				f: func(c jwt.MapClaims, key string) any {
					return mapValue[map[string]string](c, key)
				},
			},
			&action{
				result: map[string]string{"foo": "bar"},
			},
		),
		gen(
			"wrong value type",
			&condition{
				claims: jwt.MapClaims{
					"test": 123,
				},
				key: "test",
				f: func(c jwt.MapClaims, key string) any {
					return mapValue[string](c, key)
				},
			},
			&action{
				result: "",
			},
		),
		gen(
			"key not found for int",
			&condition{
				claims: jwt.MapClaims{},
				key:    "not-exist",
				f: func(c jwt.MapClaims, key string) any {
					return mapValue[int](c, key)
				},
			},
			&action{
				result: 0,
			},
		),
		gen(
			"key not found for map",
			&condition{
				claims: jwt.MapClaims{},
				key:    "not-exist",
				f: func(c jwt.MapClaims, key string) any {
					return mapValue[map[string]any](c, key)
				},
			},
			&action{
				result: map[string]any(nil),
			},
		),
		gen(
			"claim is nil",
			&condition{
				claims: nil,
				key:    "not-exist",
				f: func(c jwt.MapClaims, key string) any {
					return mapValue[map[string]any](c, key)
				},
			},
			&action{
				result: map[string]any(nil),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			result := tt.C.f(tt.C.claims, tt.C.key)
			testutil.Diff(t, tt.A.result, result)
		})
	}
}
