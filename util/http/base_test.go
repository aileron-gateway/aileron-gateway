// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package http_test

import (
	"net/http"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

func TestMethods(t *testing.T) {
	type condition struct {
		methods []v1.HTTPMethod
	}

	type action struct {
		methods []string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil", &condition{
				methods: nil,
			},
			&action{
				methods: nil,
			},
		),
		gen(
			"unknown", &condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_HTTPMethodUnknown},
			},
			&action{
				methods: []string{},
			},
		),
		gen(
			"GET", &condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_GET},
			},
			&action{
				methods: []string{http.MethodGet},
			},
		),
		gen(
			"HEAD", &condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_HEAD},
			},
			&action{
				methods: []string{http.MethodHead},
			},
		),
		gen(
			"POST", &condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_POST},
			},
			&action{
				methods: []string{http.MethodPost},
			},
		),
		gen(
			"PUT", &condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_PUT},
			},
			&action{
				methods: []string{http.MethodPut},
			},
		),
		gen(
			"PATCH", &condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_PATCH},
			},
			&action{
				methods: []string{http.MethodPatch},
			},
		),
		gen(
			"DELETE", &condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_DELETE},
			},
			&action{
				methods: []string{http.MethodDelete},
			},
		),
		gen(
			"CONNECT", &condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_CONNECT},
			},
			&action{
				methods: []string{http.MethodConnect},
			},
		),
		gen(
			"OPTIONS", &condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_OPTIONS},
			},
			&action{
				methods: []string{http.MethodOptions},
			},
		),
		gen(
			"Trace", &condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_TRACE},
			},
			&action{
				methods: []string{http.MethodTrace},
			},
		),
		gen(
			"GET/HEAD/POST", &condition{
				methods: []v1.HTTPMethod{
					v1.HTTPMethod_GET,
					v1.HTTPMethod_HEAD,
					v1.HTTPMethod_POST,
				},
			},
			&action{
				methods: []string{
					http.MethodGet,
					http.MethodHead,
					http.MethodPost,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			methods := utilhttp.Methods(tt.C.methods)
			testutil.Diff(t, tt.A.methods, methods)
		})
	}
}

func TestBaseHandler_Patterns(t *testing.T) {
	type condition struct {
		h *utilhttp.HandlerBase
	}

	type action struct {
		patterns []string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no pattern", &condition{
				h: &utilhttp.HandlerBase{
					AcceptPatterns: []string{},
				},
			},
			&action{
				patterns: []string{},
			},
		),
		gen(
			"empty pattern", &condition{
				h: &utilhttp.HandlerBase{
					AcceptPatterns: []string{""},
				},
			},
			&action{
				patterns: []string{""},
			},
		),
		gen(
			"non empty pattern", &condition{
				h: &utilhttp.HandlerBase{
					AcceptPatterns: []string{"/test"},
				},
			},
			&action{
				patterns: []string{"/test"},
			},
		),
		gen(
			"multiple patterns", &condition{
				h: &utilhttp.HandlerBase{
					AcceptPatterns: []string{"/test1", "/test2"},
				},
			},
			&action{
				patterns: []string{"/test1", "/test2"},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			p := tt.C.h.Patterns()
			testutil.Diff(t, tt.A.patterns, p)
		})
	}
}

func TestHandler_Methods(t *testing.T) {
	type condition struct {
		h *utilhttp.HandlerBase
	}

	type action struct {
		methods []string
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no methods", &condition{
				h: &utilhttp.HandlerBase{},
			},
			&action{
				methods: nil,
			},
		),
		gen(
			"no methods", &condition{
				h: &utilhttp.HandlerBase{
					AcceptMethods: []string{""},
				},
			},
			&action{
				methods: []string{""},
			},
		),
		gen(
			"one method", &condition{
				h: &utilhttp.HandlerBase{
					AcceptMethods: []string{
						http.MethodGet,
					},
				},
			},
			&action{
				methods: []string{
					http.MethodGet,
				},
			},
		),
		gen(
			"multiple methods", &condition{
				h: &utilhttp.HandlerBase{
					AcceptMethods: []string{
						http.MethodGet,
						http.MethodPost,
					},
				},
			},
			&action{
				methods: []string{
					http.MethodGet,
					http.MethodPost,
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			ms := tt.C.h.Methods()
			testutil.Diff(t, tt.A.methods, ms)
		})
	}
}
