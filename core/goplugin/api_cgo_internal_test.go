// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build (linux || darwin || freebsd) && cgo

package goplugin

import (
	"io"
	"plugin"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type testPluginInitA struct {
	err error
}

func (p *testPluginInitA) Init() error {
	return p.err
}

type testPluginInitB struct {
	lg  log.Logger
	err error
}

func (p *testPluginInitB) Init(lg log.Logger) error {
	p.lg = lg
	return p.err
}

type testPluginInitC struct {
	lg  log.Logger
	eh  core.ErrorHandler
	err error
}

func (p *testPluginInitC) Init(lg log.Logger, eh core.ErrorHandler) error {
	p.lg = lg
	p.eh = eh
	return p.err
}

func TestCreate(t *testing.T) {
	type condition struct {
		manifest   protoreflect.ProtoMessage
		testPlugin any
		testOpen   func(string) (*plugin.Plugin, error)
	}
	type action struct {
		expect     any
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	testServer := api.NewContainerAPI()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default manifest",
			[]string{}, []string{},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create GoPlugin`),
			},
		),
		gen(
			"valid plugin/initA success",
			[]string{}, []string{},
			&condition{
				manifest:   Resource.Default(),
				testPlugin: &testPluginInitA{},
				testOpen:   func(p string) (*plugin.Plugin, error) { return &plugin.Plugin{}, nil },
			},
			&action{
				expect: &testPluginInitA{},
				err:    nil,
			},
		),
		gen(
			"valid plugin/initA failed",
			[]string{}, []string{},
			&condition{
				manifest: Resource.Default(),
				testPlugin: &testPluginInitA{
					err: io.EOF,
				},
				testOpen: func(p string) (*plugin.Plugin, error) { return &plugin.Plugin{}, nil },
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create GoPlugin`),
			},
		),
		gen(
			"valid plugin/initB success",
			[]string{}, []string{},
			&condition{
				manifest:   Resource.Default(),
				testPlugin: &testPluginInitB{},
				testOpen:   func(p string) (*plugin.Plugin, error) { return &plugin.Plugin{}, nil },
			},
			&action{
				expect: &testPluginInitB{
					lg: log.GlobalLogger(log.DefaultLoggerName),
				},
				err: nil,
			},
		),
		gen(
			"valid plugin/initB failed",
			[]string{}, []string{},
			&condition{
				manifest: Resource.Default(),
				testPlugin: &testPluginInitB{
					err: io.EOF,
				},
				testOpen: func(p string) (*plugin.Plugin, error) { return &plugin.Plugin{}, nil },
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create GoPlugin`),
			},
		),
		gen(
			"valid plugin/initC success",
			[]string{}, []string{},
			&condition{
				manifest:   Resource.Default(),
				testPlugin: &testPluginInitC{},
				testOpen:   func(p string) (*plugin.Plugin, error) { return &plugin.Plugin{}, nil },
			},
			&action{
				expect: &testPluginInitC{
					lg: log.GlobalLogger(log.DefaultLoggerName),
					eh: utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
				},
				err: nil,
			},
		),
		gen(
			"valid plugin/initC failed",
			[]string{}, []string{},
			&condition{
				manifest: Resource.Default(),
				testPlugin: &testPluginInitC{
					err: io.EOF,
				},
				testOpen: func(p string) (*plugin.Plugin, error) { return &plugin.Plugin{}, nil },
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create GoPlugin`),
			},
		),
		gen(
			"lookup error",
			[]string{}, []string{},
			&condition{
				manifest: Resource.Default(),
				testOpen: func(p string) (*plugin.Plugin, error) { return &plugin.Plugin{}, nil },
			},
			&action{
				expect:     nil,
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create GoPlugin`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmpPlugin := testPlugin
			testPlugin = tt.C().testPlugin
			defer func() {
				testPlugin = tmpPlugin
			}()
			tmpOpen := testOpen
			testOpen = tt.C().testOpen
			defer func() {
				testOpen = tmpOpen
			}()

			got, err := Resource.Create(testServer, tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[log.Logger]),
				cmp.AllowUnexported(utilhttp.DefaultErrorHandler{}),
				cmp.AllowUnexported(testPluginInitA{}, testPluginInitB{}, testPluginInitC{}),
				cmpopts.EquateErrors(),
			}

			testutil.Diff(t, tt.A().expect, got, opts...)
		})
	}
}
