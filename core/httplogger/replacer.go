package httplogger

import (
	"net/textproto"
	"slices"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// headerReplacers returns replacers for HTTP request/response headers.
// The returned map will be key value pairs of header names and replacers.
// The header name is formatted by the textproto.CanonicalMIMEHeaderKey.
func headerReplacers(specs []*v1.LogHeaderSpec) (map[string][]stringReplFunc, bool, error) {
	allHeader := false
	result := make(map[string][]stringReplFunc, len(specs))

	for _, spec := range specs {
		if spec == nil || spec.Name == "" {
			continue
		}
		if spec.Name == "*" {
			allHeader = true
			continue
		}

		fs, err := txtutil.NewStringReplacers(spec.Replacers...)
		if err != nil {
			return nil, false, err
		}

		// Canonical header name.
		// "foo-bar" will be converted to "Foo-Bar".
		name := textproto.CanonicalMIMEHeaderKey(spec.Name)
		result[name] = slices.Clip(append(result[name], stringReplacerToFunc(fs)...))
	}

	return result, allHeader, nil
}

// bodyReplacers returns replacers for HTTP request/response bodies.
// The returned map will be key value pairs of MIME type and replacers.
func bodyReplacers(specs []*v1.LogBodySpec) (map[string][]bytesReplFunc, error) {
	result := make(map[string][]bytesReplFunc, len(specs))

	for _, spec := range specs {
		if spec == nil || spec.Mime == "" {
			continue
		}

		fs, err := txtutil.NewBytesReplacers(spec.Replacers...)
		if err != nil {
			return nil, err
		}
		if len(spec.JSONFields) > 0 {
			jr := &jsonFieldReplacer{
				fields:    spec.JSONFields,
				replacers: bytesReplacerToFunc(fs),
			}
			result[spec.Mime] = append(result[spec.Mime], jr.Replace)
		} else {
			result[spec.Mime] = slices.Clip(append(result[spec.Mime], bytesReplacerToFunc(fs)...))
		}
	}

	return result, nil
}

type jsonFieldReplacer struct {
	fields    []string
	replacers []bytesReplFunc
}

func (r *jsonFieldReplacer) Replace(data []byte) []byte {
	values := gjson.GetManyBytes(data, r.fields...)
	for i, f := range r.fields {
		if !values[i].Exists() {
			continue
		}
		for j := 0; j < len(r.replacers); j++ {
			data, _ = sjson.SetRawBytes(data, f, r.replacers[j]([]byte(values[i].Raw)))
		}
	}
	return data
}
