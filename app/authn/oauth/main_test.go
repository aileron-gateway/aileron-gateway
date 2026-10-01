// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package oauth

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/security"
	"github.com/golang-jwt/jwt/v5"
)

func newJWTHandler() *security.JWTHandler {
	spec := &v1.JWTHandlerSpec{
		PrivateKeys: []*v1.SigningKeySpec{
			{
				KeyID:     "test",
				Algorithm: v1.SigningKeyAlgorithm_HS256,
				KeyType:   v1.SigningKeyType_COMMON,
				KeyString: base64.StdEncoding.EncodeToString([]byte("password")),
			},
		},
		PublicKeys: []*v1.SigningKeySpec{
			{
				KeyID:     "test",
				Algorithm: v1.SigningKeyAlgorithm_HS256,
				KeyType:   v1.SigningKeyType_COMMON,
				KeyString: base64.StdEncoding.EncodeToString([]byte("password")),
			},
		},
	}
	jh, _ := security.NewJWTHandler(spec, nil)
	return jh
}

type testClient struct {
	err    error
	status int
	header http.Header
	body   io.Reader

	got *http.Request
}

func (c *testClient) RoundTrip(r *http.Request) (*http.Response, error) {
	c.got = r
	if c.err != nil {
		return nil, c.err
	}
	return &http.Response{
		StatusCode: c.status,
		Header:     c.header,
		Body:       io.NopCloser(c.body),
	}, nil
}

type testRequester struct {
	status int
	body   []byte
	err    error
}

func (r *testRequester) doRequest(ctx context.Context, endpoint string, query map[string]string) (int, []byte, error) {
	return r.status, r.body, r.err
}

type testRedeemer struct {
	status int
	resp   *TokenResponse
	err    core.HTTPError

	gotCtx    context.Context
	gotParams map[string]string
}

func (r *testRedeemer) redeemToken(ctx context.Context, params map[string]string) (int, *TokenResponse, core.HTTPError) {
	r.gotCtx = ctx
	r.gotParams = params
	return r.status, r.resp, r.err
}

type testIntrospector struct {
	status int
	claims jwt.MapClaims
	err    core.HTTPError

	gotCtx    context.Context
	gotParams map[string]string
}

func (i *testIntrospector) tokenIntrospection(ctx context.Context, params map[string]string) (int, jwt.MapClaims, core.HTTPError) {
	i.gotCtx = ctx
	i.gotParams = params
	return i.status, i.claims, i.err
}
