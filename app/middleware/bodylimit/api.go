// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package bodylimit

import (
	"cmp"
	"os"
	"path/filepath"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	httputil "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "BodyLimitMiddleware"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.BodyLimitMiddleware{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.BodyLimitMiddlewareSpec{
				MaxSize:  1 << 22, // 4MiB
				MemLimit: 0,
				TempPath: os.TempDir(),
			},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.BodyLimitMiddleware)

	eh, err := httputil.ErrorHandler(a, c.Spec.ErrorHandler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	return &bodyLimit{
		eh:       eh,
		maxSize:  c.Spec.MaxSize,
		memLimit: cmp.Or(c.Spec.MemLimit, c.Spec.MaxSize),
		tempPath: filepath.Clean(c.Spec.TempPath) + "/",
	}, nil
}
