// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package opa

import (
	"context"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/open-policy-agent/opa/bundle"
	"github.com/open-policy-agent/opa/loader"
)

// loadBundle load OPA bundle from local file path
// or remote HTTP server.
func loadBundle(path string, rt http.RoundTripper, loader loader.FileLoader) (*bundle.Bundle, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, path, nil)
		if err != nil {
			return nil, (&er.Error{
				Package:     "authz/opa",
				Type:        "load bundle",
				Description: "failed create request.",
				Detail:      path,
			}).Wrap(err)
		}

		resp, err := rt.RoundTrip(req)
		if err != nil {
			return nil, (&er.Error{
				Package:     "authz/opa",
				Type:        "load bundle",
				Description: "failed to get bundle from endpoint.",
				Detail:      path,
			}).Wrap(err)
		}
		defer func() {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}()

		if resp.StatusCode != http.StatusOK {
			return nil, (&er.Error{
				Package:     "authz/opa",
				Type:        "load bundle",
				Description: "bundle endpoint returned non 200 OK status.",
				Detail:      "got " + strconv.Itoa(resp.StatusCode),
			}).Wrap(err)
		}

		f, err := os.CreateTemp(os.TempDir(), "*.tar.gz")
		if err != nil {
			return nil, (&er.Error{
				Package:     "authz/opa",
				Type:        "load bundle",
				Description: "failed to read bundle from response body.",
				Detail:      path,
			}).Wrap(err)
		}

		// Remove temp tarball.
		defer os.Remove(f.Name())
		path = f.Name()

		_, err = f.ReadFrom(resp.Body)
		if err != nil {
			f.Close()
			return nil, (&er.Error{
				Package:     "authz/opa",
				Type:        "load bundle",
				Description: "failed to read bundle from response body.",
				Detail:      path,
			}).Wrap(err)
		}
		f.Close()
	}

	b, err := loader.AsBundle(path)
	if err != nil {
		return nil, (&er.Error{
			Package:     "authz/opa",
			Type:        "load bundle",
			Description: "failed to load bundle.",
			Detail:      path,
		}).Wrap(err)
	}

	return b, nil
}

func verificationConfig(spec *v1.BundleVerificationSpec) (*bundle.VerificationConfig, error) {
	if spec == nil {
		return nil, nil
	}
	vc := &bundle.VerificationConfig{
		PublicKeys: make(map[string]*bundle.KeyConfig, len(spec.VerificationKeys)),
		KeyID:      spec.KeyID,
		Scope:      spec.Scope,
		Exclude:    spec.Excludes,
	}
	for _, vk := range spec.VerificationKeys {
		key, err := os.ReadFile(vk.KeyFile)
		if err != nil {
			return nil, err
		}
		vc.PublicKeys[vk.KeyID] = &bundle.KeyConfig{
			Key:       string(key),
			Algorithm: vk.Algorithm,
			Scope:     vk.Scope,
		}
	}
	return vc, nil
}
