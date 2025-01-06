package api

import (
	"context"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"google.golang.org/protobuf/encoding/protojson"
)

// Creator is the interface that creates a new object from given configuration.
// Structs implement this interface can be registered to a ExtensionAPI.
type Creator interface {
	Create(API[*Request, *Response], Format, any) (any, error)
}

// NewExtensionAPI returns a new instance of the ExtensionAPI.
func NewExtensionAPI() *ExtensionAPI {
	return &ExtensionAPI{
		manifestStore: make(map[string]any, 2),
		formatStore:   make(map[string]Format, 2),
		objStore:      make(map[string]any, 2),
		creators:      make(map[string]Creator, 2),
	}
}

// ExtensionAPI is the API.
// This implements api.API[*api.Request, *api.Response] interface.
// Use api.NewExtensionAPI() to obtain an instance of this struct.
// Otherwise the instance will not be initialized properly.
type ExtensionAPI struct {
	// manifestStore stores manifests given.
	// Typically, the key will be IDs in the format of "APIGroup/APIVersion/Kind/Namespace/Name".
	manifestStore map[string]any
	// formatStore stores format of the manifest.
	// The key is the same as manifestStore.
	formatStore map[string]Format
	// objStore stores objects created by creators.
	// The key is the same as manifestStore.
	objStore map[string]any
	// creators stores instances of creator.
	// Typically, the key should be in the format of "APIGroup/APIVersion/Kind".
	creators map[string]Creator
}

// Register registers a creator to this API.
// nil creator will be ignored.
func (a *ExtensionAPI) Register(key string, c Creator) error {
	if c == nil {
		return nil
	}
	if _, ok := a.creators[key]; ok {
		return &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeExt,
			Description: ErrDscDuplicateKey,
			Detail:      "key=" + key,
		}
	}
	printDebug(debugLv2, "ExtensionAPI:", "Register:", "key="+key)
	a.creators[key] = c
	return nil
}

func (a *ExtensionAPI) Serve(ctx context.Context, req *Request) (*Response, error) {
	// Save this route in the context
	// so that the contained APIs can use parent APIs.
	ctx = ContextWithRoute(ctx, a)

	if req == nil {
		// Nil request is not allowed.
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeExt,
			Description: ErrDscNil,
		}
	}

	c, ok := a.creators[req.Key]
	if !ok {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeExt,
			Description: ErrDscNoAPI,
			Detail:      "key=" + req.Key,
		}
	}

	id := req.Key
	opt := &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}
	msg, err := ProtoMessage(req.Format, req.Content, &k.Template{}, opt)
	if err != nil {
		return nil, err // Return err as-is.
	}
	if msg != nil {
		if id, err = ParseID(msg); err != nil {
			return nil, err // Return err as-is.
		}
	}

	var content any
	switch req.Method {
	case MethodDelete: // Delete deletes stored object from this container.
		printDebug(debugLv2, "ExtensionAPI:", "DELETE:", "key="+id)
		delete(a.manifestStore, id)
		delete(a.formatStore, id)
		delete(a.objStore, id)

	case MethodGet: // Get returns stored object. Nil will be returned if not stored.
		printDebug(debugLv2, "ExtensionAPI:", "GET:", "key="+id)
		if v, ok := a.objStore[id]; ok {
			content = v // Object has already been created.
		} else {
			manifest, ok := a.manifestStore[id]
			if !ok {
				return nil, &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeExt,
					Description: ErrDscNoManifest,
					Detail:      "key=" + id,
				}
			}
			root := RootAPIFromContext(ctx)
			printDebug(debugLv1, "ExtensionAPI:", "GET:", "Create Resource:", "key="+id, manifest)
			content, err = c.Create(root, a.formatStore[id], manifest)
			if err != nil {
				return nil, err // Return err as-is.
			}
			a.objStore[id] = content
		}

	case MethodPost: // Post stores the given object in this container.
		printDebug(debugLv2, "ExtensionAPI:", "POST:", "key="+id)
		if _, ok := a.manifestStore[id]; ok {
			return nil, &er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeExt,
				Description: ErrDscDuplicateKey,
				Detail:      "key=" + id,
			}
		}
		a.manifestStore[id] = req.Content
		a.formatStore[id] = req.Format

	default:
		printDebug(debugLv2, "ExtensionAPI:", "UNDEFINED:", req.Method, "key="+id)
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeExt,
			Description: ErrDscNoMethod,
			Detail:      string(req.Method),
		}
	}

	return &Response{
		Content: content,
	}, nil
}
