// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package static

import (
	"cmp"
	"net/http"
	"path"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "core/v1"
	kind       = "StaticFileHandler"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.StaticFileHandler{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.StaticFileHandlerSpec{
				RootDir: "./",
			},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.StaticFileHandler)
	var fs http.FileSystem = http.Dir(path.Clean(c.Spec.RootDir))
	if !c.Spec.EnableListing {
		fs = &fileOnlyDir{ // Protect from directory listing attack.
			fs: fs,
		}
	}

	h := http.FileServer(fs)
	if c.Spec.StripPrefix != "" {
		h = http.StripPrefix(c.Spec.StripPrefix, h)
	}

	return &handler{
		HandlerBase: &utilhttp.HandlerBase{
			AcceptPatterns: c.Spec.Patterns,
			AcceptMethods:  utilhttp.Methods(c.Spec.Methods),
		},
		Handler: h,
		eh:      utilhttp.GlobalErrorHandler(cmp.Or(c.Metadata.ErrorHandler, utilhttp.DefaultErrorHandlerName)),
		header:  c.Spec.Header,
	}, nil
}
