// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package compression

import (
	"compress/gzip"
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/andybalholm/brotli"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestMutate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		manifest protoreflect.ProtoMessage
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"apply default values",
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				manifest: &v1.CompressionMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CompressionMiddlewareSpec{
						BrotliLevel: 4,
						GzipLevel:   6,
						MinimumSize: 1 << 10,
						TargetMIMEs: []string{
							"application/json", "application/manifest+json", "application/graphql+json", // json
							"text/html", "text/richtext", "text/plain", "text/css", // text
							"text/xml", "application/xml", "application/xhtml+xml", "image/svg+xml", // xml
							"application/javascript", "text/javascript", "text/js", // javascript
						},
					},
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			msg := Resource.Mutate(tt.C.manifest)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(v1.CompressionMiddleware{}, v1.CompressionMiddlewareSpec{}),
				cmpopts.IgnoreUnexported(k.Metadata{}, k.Status{}),
			}
			testutil.Diff(t, tt.A.manifest, msg, opts...)
		})
	}
}

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		expectedMIMEs []string
		err           any
		errPattern    *regexp.Regexp
		gzipLevel     int
		brotliLevel   int
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			&condition{
				manifest: &v1.CompressionMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CompressionMiddlewareSpec{
						BrotliLevel: 4,
						GzipLevel:   6,
						TargetMIMEs: nil,
					},
				},
			},
			&action{
				err:           nil,
				expectedMIMEs: nil,
				gzipLevel:     6,
				brotliLevel:   4,
			},
		),
		gen(
			"invalid gzip level max",
			&condition{
				manifest: &v1.CompressionMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CompressionMiddlewareSpec{
						GzipLevel: 20, // Invalid level
					},
				},
			},
			&action{
				err:           nil,
				expectedMIMEs: nil,
				gzipLevel:     9, // Adjusted to max valid level
			},
		),
		gen(
			"invalid brotli level min",
			&condition{
				manifest: &v1.CompressionMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.CompressionMiddlewareSpec{
						BrotliLevel: -1, // Invalid level
					},
				},
			},
			&action{
				err:           nil,
				expectedMIMEs: nil,
				brotliLevel:   0, // Adjusted to min valid level
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			server := api.NewContainerAPI()

			a := &API{}
			comp, err := a.Create(server, tt.C.manifest)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)

			if err == nil {
				// check MIME types
				compression := comp.(*compression)
				testutil.Diff(t, tt.A.expectedMIMEs, compression.mimes)

				gwPool := compression.gwPool.Get().(*gzip.Writer)
				testutil.Diff(t, false, gwPool == nil)

				bwPool := compression.bwPool.Get().(*brotli.Writer)
				testutil.Diff(t, false, bwPool == nil)
			}
		})
	}
}
