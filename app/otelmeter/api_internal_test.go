// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package otelmeter

import (
	"context"
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
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
			},
		),
		gen(
			"HTTPExporter",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryMeter{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.OpenTelemetryMeterSpec{
						Exporters: &v1.OpenTelemetryMeterSpec_HTTPExporterSpec{
							HTTPExporterSpec: &v1.HTTPMetricsExporterSpec{},
						},
						PeriodicReader: &v1.PeriodicReaderSpec{},
					},
				},
			},
			&action{
				err: nil,
			},
		),
		gen(
			"StdoutExporter",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryMeter{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.OpenTelemetryMeterSpec{
						Exporters: &v1.OpenTelemetryMeterSpec_StdoutExporterSpec{
							StdoutExporterSpec: &v1.StdoutMetricsExporterSpec{},
						},
						PeriodicReader: &v1.PeriodicReaderSpec{},
					},
				},
			},
			&action{
				err: nil,
			},
		),
		gen(
			"fail to create http exporter",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryMeter{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.OpenTelemetryMeterSpec{
						Exporters: &v1.OpenTelemetryMeterSpec_HTTPExporterSpec{
							HTTPExporterSpec: &v1.HTTPMetricsExporterSpec{
								TLSConfig: &k.TLSConfig{
									RootCAs: []string{
										"notExistCA",
									},
								},
							},
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create OpenTelemetryMeter`),
			},
		),
		gen(
			"fail to create grpc exporter",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryMeter{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.OpenTelemetryMeterSpec{
						Exporters: &v1.OpenTelemetryMeterSpec_GRPCExporterSpec{
							GRPCExporterSpec: &v1.GRPCMetricsExporterSpec{
								TLSConfig: &k.TLSConfig{
									RootCAs: []string{
										"notExistCA",
									},
								},
							},
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create OpenTelemetryMeter`),
			},
		),
		gen(
			"create grpc exporter with insecure",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryMeter{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.OpenTelemetryMeterSpec{
						Exporters: &v1.OpenTelemetryMeterSpec_GRPCExporterSpec{
							GRPCExporterSpec: &v1.GRPCMetricsExporterSpec{
								Insecure: true,
							},
						},
						PeriodicReader: &v1.PeriodicReaderSpec{},
					},
				},
			},
			&action{
				err: nil,
			},
		),
		gen(
			"create grpc exporter with appendOption",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.OpenTelemetryMeter{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "default",
						Name:      "default",
					},
					Spec: &v1.OpenTelemetryMeterSpec{
						Exporters: &v1.OpenTelemetryMeterSpec_GRPCExporterSpec{
							GRPCExporterSpec: &v1.GRPCMetricsExporterSpec{
								EndpointURL: "http://testURL",
							},
						},
						PeriodicReader: &v1.PeriodicReaderSpec{},
					},
				},
			},
			&action{
				err: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			server := api.NewContainerAPI()
			postTestResource(server, nil)

			_, err := Resource.Create(server, tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
		})
	}
}

func postTestResource(server api.API[*api.Request, *api.Response], res any) {
	ref := &k.Reference{
		APIVersion: "container/v1",
		Kind:       "Container",
		Namespace:  "externalOptions",
		Name:       "externalOptions",
	}
	req := &api.Request{
		Method:  api.MethodPost,
		Key:     ref.APIVersion + "/" + ref.Kind + "/" + ref.Namespace + "/" + ref.Name,
		Content: res,
	}
	if _, err := server.Serve(context.Background(), req); err != nil {
		panic(err)
	}
}
