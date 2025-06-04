// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package api

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/encoder"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-projects/go/zos"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	debugLv1 = 1
	debugLv2 = 2
	debugLv3 = 3
)

var (
	// DebugAPICall if true, trace api call.
	DebugLv3 = strings.Contains(os.Getenv("GODEBUG"), "gatewayapi=3")
	// DebugResource if true, print all manifest before calling all resource api
	DebugLv2 = DebugLv3 || strings.Contains(os.Getenv("GODEBUG"), "gatewayapi=2")
	// DebugManifest if true, print all manifest before calling create methods.
	DebugLv1 = DebugLv2 || strings.Contains(os.Getenv("GODEBUG"), "gatewayapi=1")
)

func printDebug(level int, args ...any) {
	if !DebugLv3 && level >= debugLv3 {
		return
	}
	if !DebugLv2 && level >= debugLv2 {
		return
	}
	if !DebugLv1 && level >= debugLv1 {
		return
	}
	log.Println(args...)
}

// Method is the type of method that
// all built in APIs can accept.
type Method string

const (
	MethodDelete Method = "DELETE" // Delete operation for APIs.
	MethodGet    Method = "GET"    // Get operation for APIs.
	MethodPost   Method = "POST"   // Post operation for APIs.
)

// Format is the type of data format that
// all built-in APIs can accept.
type Format string

const (
	FormatJSON           Format = "JSON"           // JSON []byte
	FormatYAML           Format = "YAML"           // YAML []byte
	FormatProtoMessage   Format = "ProtoMessage"   // protoreflect.ProtoMessage
	FormatProtoReference Format = "ProtoReference" // protoreflect.ProtoMessage only for kernel.Reference
)

// Unmarshal un-marshals the in to into with this format.
// JSON and YAML is supported.
func (f Format) Unmarshal(in any, into any) error {
	switch f {
	case FormatJSON:
		in := in.([]byte)
		return encoder.UnmarshalJSON(in, into)
	case FormatYAML:
		in := in.([]byte)
		return encoder.UnmarshalYAML(in, into)
	}
	return &er.Error{
		Package:     ErrPkg,
		Type:        ErrTypeUtil,
		Description: ErrDscFormatSupport,
		Detail:      string(f),
	}
}

// Request is the default API request.
// This struct does not implement any method and just transfer data with its fields.
// Configuration rules depend on API implementations.
type Request struct {
	Params  map[string]string
	Content any
	Method  Method
	Key     string
	Format  Format
}

// Response is the default API response.
// This struct does not implement any method and just transfer data with its fields.
// Configuration rules depend on API implementations.
type Response struct {
	Params  map[string]string
	Content any
}

// NewDefaultServeMux returns a new instance of DefaultServeMux
// which is the multiplexer for api.API[*api.Request, *api.Response].
func NewDefaultServeMux() *DefaultServeMux {
	return &DefaultServeMux{
		apis: map[string]API[*Request, *Response]{},
	}
}

// DefaultServeMux is the default API multiplexer.
// This implements api.ServeMux[string, *api.Request, *api.Response] interface.
type DefaultServeMux struct {
	apis map[string]API[*Request, *Response]
	keys []string
}

func (m *DefaultServeMux) Serve(ctx context.Context, req *Request) (*Response, error) {
	// Save this route in the context
	// so that the contained APIs can use parent APIs.
	ctx = ContextWithRoute(ctx, m)

	if req == nil {
		// Nil request is not allowed.
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeUtil,
			Description: ErrDscNil,
		}
	}

	// Find API route with prefix matching.
	// Note that the keys are sorted descending order.
	for _, k := range m.keys {
		if strings.HasPrefix(req.Key, k) {
			return m.apis[k].Serve(ctx, req)
		}
	}

	return nil, &er.Error{
		Package:     ErrPkg,
		Type:        ErrTypeUtil,
		Description: ErrDscNoAPI,
		Detail:      "key=" + req.Key,
	}
}

func (m *DefaultServeMux) Handle(key string, a API[*Request, *Response]) error {
	if a == nil {
		return nil // Ignore nil API.
	}

	if _, ok := m.apis[key]; ok {
		return &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeUtil,
			Description: ErrDscDuplicateKey,
			Detail:      "key=" + key,
		}
	}

	m.apis[key] = a
	m.keys = append(m.keys, key)
	sort.Sort(sort.Reverse(sort.StringSlice(m.keys)))
	return nil
}

// routeKey is the key type to save API routes
// into a context.
type routeKey struct{}

var (
	APIRouteContextKey = &routeKey{}
)

// ContextWithRoute returns a new context with API route information.
// This function panics when invalid type of data is stored with api.APIRouteContextKey
// in the given context.
func ContextWithRoute(ctx context.Context, a API[*Request, *Response]) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if a == nil {
		return ctx
	}
	v := ctx.Value(APIRouteContextKey)
	if v == nil {
		route := []API[*Request, *Response]{a}
		return context.WithValue(ctx, APIRouteContextKey, route)
	}
	routes := v.([]API[*Request, *Response])
	routes = append(routes, a)
	return context.WithValue(ctx, APIRouteContextKey, routes)
}

// RootAPIFromContext returns a root API of API routes.
// This function panics when invalid type of data is stored with api.APIRouteContextKey
// in the given context.
func RootAPIFromContext(ctx context.Context) API[*Request, *Response] {
	if ctx == nil {
		return nil
	}
	routes, ok := ctx.Value(APIRouteContextKey).([]API[*Request, *Response])
	if !ok || len(routes) == 0 {
		return nil
	}
	return routes[0]
}

// ProtoMessage returns protoreflect.ProtoMessage by parsing the given content.
// Typically, the defaultMsg is the configuration for API resources with default value
// in protobuf message type.
// defaultMsg will be ignored when the format is api.FormatProtoReference.
func ProtoMessage(format Format, content any, defaultMsg protoreflect.ProtoMessage, opt *protojson.UnmarshalOptions) (protoreflect.ProtoMessage, error) {
	msg := defaultMsg
	var err error
	switch format {
	case FormatJSON:
		src := proto.Clone(defaultMsg)
		b, ok := content.([]byte)
		if !ok {
			return nil, &er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeUtil,
				Description: ErrDscAssert,
				Detail:      fmt.Sprintf("convert from %T to []byte", content),
			}
		}
		b, err = zos.EnvSubst2(b)
		if err != nil {
			break
		}
		err = encoder.UnmarshalProtoFromJSON(b, src, opt)
		proto.Merge(msg, src)
	case FormatYAML:
		src := proto.Clone(defaultMsg)
		b, ok := content.([]byte)
		if !ok {
			return nil, &er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeUtil,
				Description: ErrDscAssert,
				Detail:      fmt.Sprintf("convert from %T to []byte", content),
			}
		}
		b, err = zos.EnvSubst2(b)
		if err != nil {
			break
		}
		err = encoder.UnmarshalProtoFromYAML(b, src, opt)
		proto.Merge(msg, src)
	case FormatProtoMessage:
		src, ok := content.(protoreflect.ProtoMessage)
		if !ok {
			return nil, &er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeUtil,
				Description: ErrDscAssert,
				Detail:      fmt.Sprintf("convert from %T to ProtoMessage", content),
			}
		}
		proto.Merge(msg, src)
	case FormatProtoReference:
		src, ok := content.(protoreflect.ProtoMessage)
		if !ok {
			return nil, &er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeUtil,
				Description: ErrDscAssert,
				Detail:      fmt.Sprintf("convert from %T to ProtoMessage", content),
			}
		}
		msg = src
	default:
		return nil, nil
	}
	return msg, err
}

// ParseID returns the identifier of an object in hte apiVersion/kind/namespace/name format.
// This function can parse an identifier from following format.
// Name and namespace in the Metadata field is prior to the outer values.
//
//	&struct {
//		APIVersion string `json:"apiVersion"`
//		Kind       string `json:"kind"`
//		Name       string `json:"name"`
//		Namespace  string `json:"namespace"`
//		Metadata   *struct {
//			Name      string `json:"name"`
//			Namespace string `json:"namespace"`
//		} `json:"metadata"`
//	}
func ParseID(msg protoreflect.ProtoMessage) (string, error) {
	into := &struct {
		Metadata *struct {
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
		} `json:"metadata"`
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Name       string `json:"name"`
		Namespace  string `json:"namespace"`
	}{}

	b, err := encoder.MarshalProtoToJSON(msg, nil)
	if err != nil {
		return "", err // Return err as-is. Even may be no error.
	}

	if err := encoder.UnmarshalJSON(b, into); err != nil {
		return "", err // Return err as-is.
	}

	if into.Metadata != nil {
		into.Name = into.Metadata.Name
		into.Namespace = into.Metadata.Namespace
	}
	into.Name = cmp.Or(into.Name, "default")           // Set "default" if empty.
	into.Namespace = cmp.Or(into.Namespace, "default") // Set "default" if empty.

	id := into.APIVersion + "/" + into.Kind + "/" + into.Namespace + "/" + into.Name
	return id, nil
}

// ReferObject refer and get an object from the API searching by the given reference information.
// This returns nil when the nil reference is given.
// This returns an error when the Get method of the API returns any error.
// This function panics when nil API is given by the first argument.
func ReferObject(a API[*Request, *Response], ref *k.Reference) (any, error) {
	if ref == nil {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeUtil,
			Description: ErrDscNil,
			Detail:      "cannot reference resource by nil",
		}
	}
	req := &Request{
		Method:  MethodGet,
		Key:     ref.APIVersion + "/" + ref.Kind,
		Format:  FormatProtoReference,
		Content: ref,
	}
	res, err := a.Serve(context.Background(), req)
	if err != nil {
		return nil, err // Return err as-is.
	}
	return res.Content, nil
}

// ReferTypedObject refer and get an object from the given reference information.
// Note that a zero value of type T will be returned when there is an error.
// This function panics when nil API is given by the first argument.
func ReferTypedObject[T any](a API[*Request, *Response], ref *k.Reference) (T, error) {
	var t T
	obj, err := ReferObject(a, ref)
	if err != nil {
		return t, err
	}
	typed, ok := obj.(T)
	if !ok {
		key := strings.Join([]string{ref.APIVersion, ref.Kind, ref.Namespace, ref.Name}, "/")
		typ := strings.TrimPrefix(reflect.TypeOf(new(T)).String(), "*")
		return t, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeUtil,
			Description: ErrDscAssert,
			Detail:      fmt.Sprintf("from %T to %s. may be %s is not defined?", obj, typ, key),
		}
	}
	return typed, nil
}

// ReferTypedObjects returns objects referenced by the given refs.
// Note that a nil slice will be returned when there is an error.
// This function panics when nil API is given by the first argument.
func ReferTypedObjects[T any](a API[*Request, *Response], refs ...*k.Reference) ([]T, error) {
	ts := make([]T, 0, len(refs))
	for _, ref := range refs {
		typed, err := ReferTypedObject[T](a, ref)
		if err != nil {
			return nil, err // Return err as-is.
		}
		ts = append(ts, typed)
	}
	return ts, nil
}
