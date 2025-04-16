// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package api

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/aileron-gateway/aileron-gateway/kernel/encoder"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// KeyAccept is the key of the parameter.
// This key can be used for Get method for the factory API.
const KeyAccept = "Accept"

// Resource is the interface of API resources.
// Resources are used by being registered to a FactoryAPI.
type Resource interface {
	// Default returns a new instance of ProtoMessage with default values.
	Default() protoreflect.ProtoMessage
	// Mutate changes the given ProtoMessage if necessary
	// and return the changed ProtoMessage.
	// Given ProtoMessage has the same type as the one returned by the Default().
	// Input message is the merged with the default value and
	// user defined configurations.
	Mutate(protoreflect.ProtoMessage) protoreflect.ProtoMessage
	// Validate validates the given ProtoMessage
	// and return an error when it was invalid.
	// Given ProtoMessage has the same type as the one returned by the Default().
	Validate(protoreflect.ProtoMessage) error
	// Create creates a new instance of the resource.
	// Given ProtoMessage has the same type as the one returned by the Default().
	Create(API[*Request, *Response], protoreflect.ProtoMessage) (any, error)
	// Delete deletes the created resource.
	// Given ProtoMessage has the same type as the one returned by the Default().
	Delete(API[*Request, *Response], protoreflect.ProtoMessage, any) error
}

// BaseResource is the base struct for api.Resource interface.
// Embed this struct to avoid unnecessary method implementation
// to satisfy api.Resource interface.
// This struct does not implement Create method because the method is
// required by all resource implementations.
type BaseResource struct {
	DefaultProto protoreflect.ProtoMessage
}

func (b *BaseResource) Default() protoreflect.ProtoMessage {
	return proto.Clone(b.DefaultProto)
}

func (b *BaseResource) Validate(msg protoreflect.ProtoMessage) error {
	v, _ := protovalidate.New()
	if err := v.Validate(msg); err != nil {
		json, _ := encoder.MarshalProtoToJSON(msg, &protojson.MarshalOptions{Multiline: true, Indent: "  ", AllowPartial: true})
		return (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeFactory,
			Description: ErrDscProtoValidate,
			Detail:      reflect.TypeOf(msg).String() + string(addLineNumber(json)),
		}).Wrap(err)
	}
	return nil
}

func (b *BaseResource) Mutate(msg protoreflect.ProtoMessage) protoreflect.ProtoMessage {
	return msg
}

func (b *BaseResource) Delete(_ API[*Request, *Response], _ protoreflect.ProtoMessage, _ any) error {
	return nil
}

// addLineNumber adds line number to the input content.
func addLineNumber(in []byte) []byte {
	re := regexp.MustCompile("(?m)^")
	row := 0
	b := re.ReplaceAllFunc(in, func(b []byte) []byte { row += 1; return []byte(fmt.Sprintf("%04d|", row)) })
	line := []byte("-----1---------1---------1---------1---------")
	pre := append(append([]byte("\n\n"), line...), []byte("\n")...)
	post := append(append([]byte("\n"), line...), []byte("\n\n")...)
	return append(append(pre, b...), post...)
}

// NewFactoryAPI returns a new instance of FactoryAPI.
func NewFactoryAPI() *FactoryAPI {
	return &FactoryAPI{
		protoStore: map[string]protoreflect.ProtoMessage{},
		objStore:   map[string]any{},
		resources:  map[string]Resource{},
	}
}

// FactoryAPI is the API.
// This implements api.API[*api.Request, *api.Response] interface.
// Use api.NewFactoryAPI() to obtain an instance of this struct.
// Otherwise the instance will not be initialized properly.
type FactoryAPI struct {
	// protoStore stores proto messages.
	// The key will be IDs in the format of "APIGroup/APIVersion/Kind/Namespace/Name".
	protoStore map[string]protoreflect.ProtoMessage
	// objStore stores objects created by resources.
	// The key will be IDs in the format of "APIGroup/APIVersion/Kind/Namespace/Name".
	objStore map[string]any
	// resources stores resources.
	// The keys should be in the format of "APIGroup/APIVersion/Kind".
	resources map[string]Resource
}

// Register registers a resource to this API.
// nil resource will be ignored.
func (a *FactoryAPI) Register(key string, r Resource) error {
	if r == nil {
		return nil
	}
	if _, ok := a.resources[key]; ok {
		return &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeFactory,
			Description: ErrDscDuplicateKey,
			Detail:      "key=" + key,
		}
	}
	printDebug(debugLv2, "FactoryAPI:", "Register:", "key="+key)
	a.resources[key] = r
	return nil
}

func (a *FactoryAPI) Serve(ctx context.Context, req *Request) (*Response, error) {
	// Save this route in the context
	// so that the contained APIs can use parent APIs.
	ctx = ContextWithRoute(ctx, a)

	if req == nil {
		// Nil request is not allowed.
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeFactory,
			Description: ErrDscNil,
		}
	}

	r, ok := a.resources[req.Key]
	if !ok {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeFactory,
			Description: ErrDscNoAPI,
			Detail:      "key=" + req.Key,
		}
	}

	var content any
	switch req.Method {
	case MethodDelete:
		printDebug(debugLv2, "FactoryAPI:", "DELETE:", "key="+req.Key)
		if err := a.delete(ctx, req, r); err != nil {
			return nil, err // Return err as-is.
		}

	case MethodGet:
		printDebug(debugLv2, "FactoryAPI:", "GET:", "key="+req.Key)
		c, err := a.get(ctx, req, r)
		if err != nil {
			return nil, err // Return err as-is.
		}
		content = c

	case MethodPost:
		printDebug(debugLv2, "FactoryAPI:", "POST:", "key="+req.Key)
		if err := a.post(ctx, req, r); err != nil {
			return nil, err // Return err as-is.
		}

	default:
		printDebug(debugLv2, "FactoryAPI:", "UNDEFINED:", req.Method, "key="+req.Key)
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeFactory,
			Description: ErrDscNoMethod,
			Detail:      string(req.Method),
		}
	}

	return &Response{
		Params:  map[string]string{},
		Content: content,
	}, nil
}

func (a *FactoryAPI) delete(ctx context.Context, req *Request, r Resource) error {
	msg, err := ProtoMessage(req.Format, req.Content, r.Default(), nil)
	if err != nil {
		return err // Return err as-is.
	}
	id, err := ParseID(msg)
	if err != nil {
		return err // Return err as-is.
	}

	p := a.protoStore[id]
	o := a.objStore[id]

	root := RootAPIFromContext(ctx)
	if err := r.Delete(root, p, o); err != nil {
		return err // Return err as-is.
	}

	delete(a.protoStore, id)
	delete(a.objStore, id)
	return nil
}

func (a *FactoryAPI) get(ctx context.Context, req *Request, r Resource) (any, error) {
	msg, err := ProtoMessage(req.Format, req.Content, r.Default(), nil)
	if err != nil {
		return nil, err // Return err as-is.
	}
	id, err := ParseID(msg)
	if err != nil {
		return nil, err // Return err as-is.
	}

	if strings.HasSuffix(id, "/template/template") {
		msg = r.Mutate(r.Default())
	} else {
		p, ok := a.protoStore[id]
		if !ok {
			return nil, &er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeFactory,
				Description: ErrDscNoManifest,
				Detail:      "key=" + id,
			}
		}
		msg = p
	}

	opt := &protojson.MarshalOptions{
		Multiline:       true,
		Indent:          "  ",
		AllowPartial:    true,
		UseProtoNames:   false,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}

	switch Format(req.Params[KeyAccept]) {
	case FormatJSON:
		return encoder.MarshalProtoToJSON(msg, opt)
	case FormatYAML:
		return encoder.MarshalProtoToYAML(msg, opt)
	case FormatProtoMessage:
		return proto.Clone(msg), nil
	default:
		if v, ok := a.objStore[id]; ok {
			return v, nil
		}
		root := RootAPIFromContext(ctx)
		configJSON, _ := encoder.MarshalProtoToJSON(msg, nil)
		printDebug(debugLv1, "FactoryAPI:", "GET:", "Create Resource:", "key="+id, string(configJSON))
		content, err := r.Create(root, msg)
		if err != nil {
			return nil, err // Return err as-is.
		}
		a.objStore[id] = content
		return content, nil
	}
}

func (a *FactoryAPI) post(_ context.Context, req *Request, r Resource) error {
	msg, err := ProtoMessage(req.Format, req.Content, r.Default(), nil)
	if err != nil {
		return err // Return err as-is.
	}

	configJSON, _ := encoder.MarshalProtoToJSON(msg, nil)
	printDebug(debugLv2, "FactoryAPI:", "POST:", "Mutate Resource:", string(configJSON))
	msg = r.Mutate(msg)

	configJSON, _ = encoder.MarshalProtoToJSON(msg, nil)
	printDebug(debugLv2, "FactoryAPI:", "POST:", "Validate Resource:", string(configJSON))
	if err := r.Validate(msg); err != nil {
		return err // Return err as-is.
	}

	id, err := ParseID(msg)
	if err != nil {
		return err // Return err as-is.
	}

	if _, ok := a.protoStore[id]; ok {
		return &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeFactory,
			Description: ErrDscDuplicateKey,
			Detail:      "key=" + req.Key,
		}
	}

	// Clone message so it will not be changed.
	a.protoStore[id] = msg
	return nil
}
