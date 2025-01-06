package httpclient

import (
	stdcmp "cmp"
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/network"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/quic-go/quic-go/http3"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// testDir is the path to the test data.
// This path can be changed by the environmental variable.
var testDir = stdcmp.Or(os.Getenv("TEST_DIR"), "../../test/")

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		expect     any
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndDefaultManifest := tb.Condition("default manifest", "input default manifest")
	cndHTTP := tb.Condition("http config", "input and expect http transport")
	cndHTTP2 := tb.Condition("http2 config", "input and expect http2 transport")
	cndHTTP3 := tb.Condition("http3 config", "input and expect http3 transport")
	cndValid := tb.Condition("valid config", "input valid config")
	cndErrorReference := tb.Condition("input error reference", "input an error reference to an object")
	actCheckError := tb.Action("check the returned error", "check that the returned error is the one expected")
	actCheckNoError := tb.Action("check no error", "check that there is no error returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{cndDefaultManifest, cndValid},
			[]string{actCheckNoError},
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
			[]string{cndHTTP, cndValid},
			[]string{actCheckNoError},
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
					ProxyConnectHeader:  http.Header{},
					MaxIdleConnsPerHost: 1024,
					TLSNextProto:        make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
				},
				err: nil,
			},
		),
		gen(
			"create http2",
			[]string{cndHTTP3, cndValid},
			[]string{actCheckNoError},
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
			[]string{cndHTTP3, cndValid},
			[]string{actCheckNoError},
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
			[]string{cndHTTP, cndValid},
			[]string{actCheckNoError},
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
			[]string{cndHTTP2},
			[]string{actCheckError},
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
			[]string{cndHTTP2},
			[]string{actCheckError},
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
			[]string{cndHTTP3},
			[]string{actCheckError},
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
			[]string{cndErrorReference},
			[]string{actCheckError},
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			got, err := Resource.Create(api.NewContainerAPI(), tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(http.Transport{}, http2.Transport{}, http3.Transport{}),
				cmp.Comparer(testutil.ComparePointer[func(*http.Request) (*url.URL, error)]),                   // Proxy
				cmp.Comparer(testutil.ComparePointer[func(context.Context, string, string) (net.Conn, error)]), // DialContext
				cmpopts.IgnoreFields(http.Transport{}, "DialContext"),
				cmpopts.IgnoreFields(http2.Transport{}, "DialTLSContext"),
				cmp.Comparer(func(x, y core.RoundTripperFunc) bool { return true }), // Retry tripperware.
			}
			testutil.Diff(t, tt.A().expect, got, opts...)
		})
	}
}
