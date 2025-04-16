// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build (linux || darwin || freebsd) && cgo

package goplugin

import (
	"plugin"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "core/v1"
	kind       = "GoPlugin"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.GoPlugin{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.GoPluginSpec{
				SymbolName: "Plugin",
			},
		},
	},
}

type API struct {
	*api.BaseResource
}

// testPlugin is the plugin only used for testing.
var testPlugin any

// testOpen is the plugin open function only used for testing.
var testOpen func(path string) (*plugin.Plugin, error)

func (*API) Create(a api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.GoPlugin)

	lg := log.DefaultOr(c.Metadata.Logger)

	eh, err := utilhttp.ErrorHandler(a, c.Spec.ErrorHandler)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	p, err := plugin.Open(c.Spec.PluginPath)
	if testOpen != nil { // Skip error when testing.
		p, err = testOpen(c.Spec.PluginPath)
	}
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	s, err := p.Lookup(c.Spec.SymbolName)
	if testPlugin != nil { // Replace plugin for testing.
		s = testPlugin
		err = nil
	}
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	if i, ok := s.(InitializerA); ok {
		if err := i.Init(); err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
	}
	if i, ok := s.(InitializerB); ok {
		if err := i.Init(lg); err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
	}
	if i, ok := s.(InitializerC); ok {
		if err := i.Init(lg, eh); err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
	}

	return s, nil
}

type InitializerA interface {
	Init() error
}

type InitializerB interface {
	Init(log.Logger) error
}

type InitializerC interface {
	Init(log.Logger, core.ErrorHandler) error
}
