package api

import (
	"context"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
)

// NewContainerAPI returns a new instance of ContainerAPI.
func NewContainerAPI() *ContainerAPI {
	return &ContainerAPI{
		objStore: map[string]any{},
	}
}

// ContainerAPI is the API that stores any object.
// This implements api.API[*api.Request, *api.Response] interface.
// Use api.NewContainerAPI() to obtain an instance of this struct.
// Otherwise the instance will not be initialized properly.
// Allowed operations are described below and see the examples for usages.
//   - Post: Store objects inside the container.
//   - Get: Get the stored object from the container. Nil content will be returned when no object was found.
//   - Delete: Delete objects from the container.
type ContainerAPI struct {
	// objStore stores objects given by clients.
	// Typically, the key will be IDs in the format of "APIGroup/APIVersion/Kind/Namespace/Name".
	objStore map[string]any
}

func (a *ContainerAPI) Serve(ctx context.Context, req *Request) (*Response, error) {
	if req == nil {
		// Nil request is not allowed.
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeContainer,
			Description: ErrDscNil,
		}
	}

	id := req.Key
	msg, err := ProtoMessage(req.Format, req.Content, &k.Reference{}, nil)
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
		printDebug(debugLv2, "ContainerAPI:", "DELETE:", "key="+id)
		delete(a.objStore, id)

	case MethodGet: // Get returns stored object. Nil will be returned if no object found.
		printDebug(debugLv2, "ContainerAPI:", "GET:", "key="+id)
		content = a.objStore[id]

	case MethodPost: // Post stores the given object in this container.
		printDebug(debugLv2, "ContainerAPI:", "POST:", "key="+id)
		if _, ok := a.objStore[id]; ok {
			return nil, &er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeContainer,
				Description: ErrDscDuplicateKey,
				Detail:      "key=" + req.Key,
			}
		}
		a.objStore[id] = req.Content

	default:
		printDebug(debugLv2, "ContainerAPI:", "UNDEFINED:", req.Method, "key="+id)
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeContainer,
			Description: ErrDscNoMethod,
			Detail:      string(req.Method),
		}
	}

	return &Response{
		Content: content,
	}, nil
}
