package compression

import (
	"slices"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"google.golang.org/protobuf/reflect/protoreflect"
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
				BrotliLevel: 4,
				GzipLevel:   6,
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
func (*API) Mutate(msg protoreflect.ProtoMessage) protoreflect.ProtoMessage {
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

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.CompressionMiddleware)

	gzipLevel := restrictBetween(int(c.Spec.GzipLevel), 1, 9)      // BestSpeed=1, BestCompression=9.
	brotliLevel := restrictBetween(int(c.Spec.BrotliLevel), 0, 11) // BestSpeed=0, BestCompression=11.

	return &compression{
		mimes:       slices.Clip(c.Spec.TargetMIMEs),
		minimumSize: int64(c.Spec.MinimumSize),
		gwPool:      newGzipWriterPool(gzipLevel),
		bwPool:      newBrotliWriterPool(brotliLevel),
	}, nil
}

// restrictBetween restricts the given target int value to the value between "min" and "max".
func restrictBetween(target, min, max int) int {
	if target < min {
		return min
	}
	if target > max {
		return max
	}
	return target
}
