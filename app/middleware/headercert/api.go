// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package headercert

import (
	"cmp"
	"crypto/x509"
	"errors"
	"os"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "app/v1"
	kind       = "HeaderCertMiddleware"
	Key        = apiVersion + "/" + kind
)

var ErrAddCert = errors.New("headercert: failed to add root certificate to CertPool")

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.HeaderCertMiddleware{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.HeaderCertMiddlewareSpec{
				RootCAs:            []string{},
				CertHeader:         "X-SSL-Client-Cert",
				FingerprintpHeader: "",
			},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.HeaderCertMiddleware)
	_ = log.DefaultOr(c.Metadata.Logger)
	eh := utilhttp.GlobalErrorHandler(cmp.Or(c.Metadata.ErrorHandler, utilhttp.DefaultErrorHandlerName))

	pool, err := loadRootCert(c.Spec.RootCAs)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	return &headerCert{
		eh: eh,
		opts: x509.VerifyOptions{
			Roots: pool,
		},
		certHeader: c.Spec.CertHeader,
		fpHeader:   c.Spec.FingerprintpHeader,
	}, nil
}

func loadRootCert(rootCAs []string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	// Read the root certificate specified in the local file
	for _, c := range rootCAs {
		pem, err := os.ReadFile(c)
		if err != nil {
			return nil, err
		}
		// Add the root certificate to CertPool
		if !pool.AppendCertsFromPEM(pem) {
			return nil, ErrAddCert
		}
	}
	return pool, nil
}
