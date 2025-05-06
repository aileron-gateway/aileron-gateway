// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package soaprest

import (
	"encoding/xml"
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/app/middleware/soaprest/zxml"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		err        any
		errPattern *regexp.Regexp
		expect     any
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{},
			[]string{},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				err: nil,
				expect: &soapREST{
					converter: &zxml.JSONConverter{
						EncodeDecoder: &zxml.Simple{
							TextKey:      "$",
							AttrPrefix:   "@",
							NamespaceSep: ":",
							TrimSpace:    false,
							PreferShort:  false,
						},
						Header: xml.Header,
					},
				},
			},
		),
		gen(
			"create with modified simple converter",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.SOAPRESTMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Name:      "default",
						Namespace: "default",
					},
					Spec: &v1.SOAPRESTMiddlewareSpec{
						Rules: &v1.SOAPRESTMiddlewareSpec_Simple{
							Simple: &v1.SimpleSpec{
								TextKey:      "$test",
								AttrPrefix:   "@test",
								NamespaceSep: ":test",
								TrimSpace:    true,
								PreferShort:  true,
							},
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &soapREST{
					converter: &zxml.JSONConverter{
						EncodeDecoder: &zxml.Simple{
							TextKey:      "$test",
							AttrPrefix:   "@test",
							NamespaceSep: ":test",
							TrimSpace:    true,
							PreferShort:  true,
						},
						Header: xml.Header,
					},
				},
			},
		),
		gen(
			"create with modified rayfish converter",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.SOAPRESTMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Name:      "default",
						Namespace: "default",
					},
					Spec: &v1.SOAPRESTMiddlewareSpec{
						Rules: &v1.SOAPRESTMiddlewareSpec_Rayfish{
							Rayfish: &v1.RayfishSpec{
								NameKey:      "#testName",
								TextKey:      "#testText",
								ChildrenKey:  "#testChildren",
								AttrPrefix:   "@test",
								NamespaceSep: ":test",
								TrimSpace:    true,
							},
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &soapREST{
					converter: &zxml.JSONConverter{
						EncodeDecoder: &zxml.RayFish{
							NameKey:      "#testName",
							TextKey:      "#testText",
							ChildrenKey:  "#testChildren",
							AttrPrefix:   "@test",
							NamespaceSep: ":test",
							TrimSpace:    true,
						},
						Header: xml.Header,
					},
				},
			},
		),
		gen(
			"create with modified badgerfish converter",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.SOAPRESTMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Name:      "default",
						Namespace: "default",
					},
					Spec: &v1.SOAPRESTMiddlewareSpec{
						Rules: &v1.SOAPRESTMiddlewareSpec_Badgerfish{
							Badgerfish: &v1.BadgerfishSpec{
								TextKey:      "$test",
								AttrPrefix:   "@test",
								NamespaceSep: ":test",
								TrimSpace:    true,
							},
						},
					},
				},
			},
			&action{
				err: nil,
				expect: &soapREST{
					converter: &zxml.JSONConverter{
						EncodeDecoder: &zxml.BadgerFish{
							TextKey:      "$test",
							AttrPrefix:   "@test",
							NamespaceSep: ":test",
							TrimSpace:    true,
						},
						Header: xml.Header,
					},
				},
			},
		),
		gen(
			"fail to get ErrorHandler",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.SOAPRESTMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.SOAPRESTMiddlewareSpec{
						ErrorHandler: &k.Reference{
							Name: "notExist",
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create SOAPRESTMiddleware`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			server := api.NewContainerAPI()
			a := &API{}
			got, err := a.Create(server, tt.C().manifest)
			opts := []cmp.Option{
				cmp.AllowUnexported(soapREST{}, zxml.JSONConverter{}, zxml.Simple{}),
				cmpopts.IgnoreFields(soapREST{}, "eh"),
				cmpopts.IgnoreFields(zxml.JSONConverter{}, "jsonEncoderOpts", "jsonDecoderOpts", "xmlEncoderOpts", "xmlDecoderOpts"),
				cmpopts.IgnoreFields(zxml.Simple{}, "emptyVal"),
				cmpopts.IgnoreFields(zxml.RayFish{}, "emptyVal"),
				cmpopts.IgnoreFields(zxml.BadgerFish{}, "emptyVal"),
			}
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
			testutil.Diff(t, tt.A().expect, got, opts...)
		})
	}
}
