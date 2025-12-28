// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpclient

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/network"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/quic-go/quic-go/http3"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// testDir is the path to the test data.
var testDir = "../../test/"

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		expect     any
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				expect: network.DefaultHTTPTransport,
				err:    nil,
			},
		),
		gen(
			"create http",
			&condition{
				manifest: &v1.HTTPClient{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPClientSpec{
						Transports: &v1.HTTPClientSpec_HTTPTransportConfig{
							HTTPTransportConfig: &k.HTTPTransportConfig{},
						},
					},
				},
			},
			&action{
				expect: &http.Transport{
					Proxy:               http.ProxyFromEnvironment,
					MaxIdleConnsPerHost: 1024,
					TLSNextProto:        make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
				},
				err: nil,
			},
		),
		gen(
			"create http2",
			&condition{
				manifest: &v1.HTTPClient{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPClientSpec{
						Transports: &v1.HTTPClientSpec_HTTP2TransportConfig{
							HTTP2TransportConfig: &k.HTTP2TransportConfig{},
						},
					},
				},
			},
			&action{
				expect: &http2.Transport{},
				err:    nil,
			},
		),
		gen(
			"create http3",
			&condition{
				manifest: &v1.HTTPClient{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPClientSpec{
						Transports: &v1.HTTPClientSpec_HTTP3TransportConfig{
							HTTP3TransportConfig: &k.HTTP3TransportConfig{},
						},
					},
				},
			},
			&action{
				expect: &http3.Transport{},
				err:    nil,
			},
		),
		gen(
			"create with retry",
			&condition{
				manifest: &v1.HTTPClient{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPClientSpec{
						RetryConfig: &v1.RetryConfig{
							MaxRetry:         1,
							RetryStatusCodes: []int32{500, 503},
						},
					},
				},
			},
			&action{
				expect: (&retry{}).Tripperware(nil), // Only core.RoundTripperFunc type will be checked. We need to improve the tests.
				err:    nil,
			},
		),
		gen(
			"fail to get http transport config",
			&condition{
				manifest: &v1.HTTPClient{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPClientSpec{
						Transports: &v1.HTTPClientSpec_HTTPTransportConfig{
							HTTPTransportConfig: &k.HTTPTransportConfig{
								TLSConfig: &k.TLSConfig{RootCAs: []string{testDir + "ut/not-exist/"}},
							},
						},
					},
				},
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPClient`),
			},
		),
		gen(
			"fail to get http2 transport config",
			&condition{
				manifest: &v1.HTTPClient{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPClientSpec{
						Transports: &v1.HTTPClientSpec_HTTP2TransportConfig{
							HTTP2TransportConfig: &k.HTTP2TransportConfig{
								TLSConfig: &k.TLSConfig{RootCAs: []string{testDir + "ut/not-exist/"}},
							},
						},
					},
				},
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPClient`),
			},
		),
		gen(
			"fail to get http3 transport config",
			&condition{
				manifest: &v1.HTTPClient{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPClientSpec{
						Transports: &v1.HTTPClientSpec_HTTP3TransportConfig{
							HTTP3TransportConfig: &k.HTTP3TransportConfig{
								TLSConfig: &k.TLSConfig{RootCAs: []string{testDir + "ut/not-exist/"}},
							},
						},
					},
				},
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPClient`),
			},
		),
		gen(
			"fail to get tripperware",
			&condition{
				manifest: &v1.HTTPClient{
					Metadata: &k.Metadata{},
					Spec: &v1.HTTPClientSpec{
						Tripperwares: []*k.Reference{{APIVersion: "wrong"}},
					},
				},
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create HTTPClient`),
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			got, err := Resource.Create(api.NewContainerAPI(), tt.C.manifest)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(http.Transport{}, http2.Transport{}, http3.Transport{}),
				cmp.Comparer(testutil.ComparePointer[func(*http.Request) (*url.URL, error)]),                   // Proxy
				cmp.Comparer(testutil.ComparePointer[func(context.Context, string, string) (net.Conn, error)]), // DialContext
				cmpopts.IgnoreFields(http.Transport{}, "DialContext"),
				cmpopts.IgnoreFields(http2.Transport{}, "DialTLSContext"),
				cmp.Comparer(func(x, y core.RoundTripperFunc) bool { return true }), // Retry tripperware.
			}
			testutil.Diff(t, tt.A.expect, got, opts...)
		})
	}
}
