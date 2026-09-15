// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package compression

import (
	"slices"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"google.golang.org/protobuf/proto"
)

const (
	apiVersion = "app/v1"
	kind       = "CompressionMiddleware"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.CompressionMiddleware{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.CompressionMiddlewareSpec{
				MinimumSize: 1 << 10, // 1024 bytes.
			},
		},
	},
}

type API struct {
	*api.BaseResource
}

// Mutate changes configured values.
// The values of the msg which is given as the argument is the merged message of default values and user defined values.
// Changes for the fields of msg in this function make the final values which will be the input for validate and create function.
// Default values for "repeated" or "oneof" fields can also be applied in this function if necessary.
// Please check msg!=nil and asserting the mgs does not panic even they won't from the view of overall architecture of the gateway.
func (*API) Mutate(msg proto.Message) proto.Message {
	c := msg.(*v1.CompressionMiddleware)

	if len(c.Spec.TargetMIMEs) == 0 {
		// Apply default target mimes. Checkout the links below for references.
		// https://www.iana.org/assignments/media-types/media-types.xhtml
		// https://developers.cloudflare.com/speed/optimization/content/brotli/content-compression/
		// https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/ServingCompressedFiles.html
		c.Spec.TargetMIMEs = []string{
			"application/json", "application/manifest+json", "application/graphql+json", // json
			"text/html", "text/richtext", "text/plain", "text/css", // text
			"text/xml", "application/xml", "application/xhtml+xml", "image/svg+xml", // xml
			"application/javascript", "text/javascript", "text/js", // javascript
		}
	}

	return c
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg proto.Message) (any, error) {
	c := msg.(*v1.CompressionMiddleware)
	return &compression{
		mimes:       slices.Clip(c.Spec.TargetMIMEs),
		minimumSize: int64(c.Spec.MinimumSize),

		gzipDisabled: c.Spec.GzipLevel == 0,                                      // Disable=0
		gwPool:       newGzipWriterPool(restrictBetween(c.Spec.GzipLevel, 1, 9)), // BestSpeed=1, BestCompression=9

		brotliDisabled: c.Spec.BrotliLevel == 0,                                           // Disable=0
		bwPool:         newBrotliWriterPool(restrictBetween(c.Spec.BrotliLevel-1, 0, 12)), // BestSpeed=0, BestCompression=11
	}, nil
}

// restrictBetween restricts the given target int value to the value between "min" and "max".
func restrictBetween(target int32, min, max int) int {
	t := int(target)
	if t < min {
		return min
	}
	if t > max {
		return max
	}
	return t
}
