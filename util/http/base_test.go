package http_test

import (
	"net/http"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
)

func TestMethods(t *testing.T) {
	type condition struct {
		methods []v1.HTTPMethod
	}

	type action struct {
		methods []string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputNil := tb.Condition("nil", "input nil")
	cndInputUnknown := tb.Condition("Unknown", "input unknown method")
	cndInputGet := tb.Condition("Get", "input Get method")
	cndInputHead := tb.Condition("Head", "input Head method")
	cndInputPost := tb.Condition("Post", "input Post method")
	cndInputPut := tb.Condition("Put", "input Put method")
	cndInputPatch := tb.Condition("Patch", "input Patch method")
	cndInputDelete := tb.Condition("Delete", "input Delete method")
	cndInputConnect := tb.Condition("Connect", "input Connect method")
	cndInputOptions := tb.Condition("Options", "input Options method")
	cndInputTrace := tb.Condition("Trace", "input Trace method")
	actCheckNil := tb.Action("nil", "check that nil is returned")
	actCheckGet := tb.Action("Get", "check that Get is included in the returned slice")
	actCheckHead := tb.Action("Head", "check that Head is included in the returned slice")
	actCheckPost := tb.Action("Post", "check that Post is included in the returned slice")
	actCheckPut := tb.Action("Put", "check that Put is included in the returned slice")
	actCheckPatch := tb.Action("Patch", "check that Patch is included in the returned slice")
	actCheckDelete := tb.Action("Delete", "check that Delete is included in the returned slice")
	actCheckConnect := tb.Action("Connect", "check that Connect is included in the returned slice")
	actCheckOptions := tb.Action("Options", "check that Options is included in the returned slice")
	actCheckTrace := tb.Action("Trace", "check that Trace is included in the returned slice")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil",
			[]string{cndInputNil},
			[]string{actCheckNil},
			&condition{
				methods: nil,
			},
			&action{
				methods: nil,
			},
		),
		gen(
			"unknown",
			[]string{cndInputUnknown},
			[]string{actCheckNil},
			&condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_HTTPMethodUnknown},
			},
			&action{
				methods: []string{},
			},
		),
		gen(
			"GET",
			[]string{cndInputGet},
			[]string{actCheckGet},
			&condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_GET},
			},
			&action{
				methods: []string{http.MethodGet},
			},
		),
		gen(
			"HEAD",
			[]string{cndInputHead},
			[]string{actCheckHead},
			&condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_HEAD},
			},
			&action{
				methods: []string{http.MethodHead},
			},
		),
		gen(
			"POST",
			[]string{cndInputPost},
			[]string{actCheckPost},
			&condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_POST},
			},
			&action{
				methods: []string{http.MethodPost},
			},
		),
		gen(
			"PUT",
			[]string{cndInputPut},
			[]string{actCheckPut},
			&condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_PUT},
			},
			&action{
				methods: []string{http.MethodPut},
			},
		),
		gen(
			"PATCH",
			[]string{cndInputPatch},
			[]string{actCheckPatch},
			&condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_PATCH},
			},
			&action{
				methods: []string{http.MethodPatch},
			},
		),
		gen(
			"DELETE",
			[]string{cndInputDelete},
			[]string{actCheckDelete},
			&condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_DELETE},
			},
			&action{
				methods: []string{http.MethodDelete},
			},
		),
		gen(
			"CONNECT",
			[]string{cndInputConnect},
			[]string{actCheckConnect},
			&condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_CONNECT},
			},
			&action{
				methods: []string{http.MethodConnect},
			},
		),
		gen(
			"OPTIONS",
			[]string{cndInputOptions},
			[]string{actCheckOptions},
			&condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_OPTIONS},
			},
			&action{
				methods: []string{http.MethodOptions},
			},
		),
		gen(
			"Trace",
			[]string{cndInputTrace},
			[]string{actCheckTrace},
			&condition{
				methods: []v1.HTTPMethod{v1.HTTPMethod_TRACE},
			},
			&action{
				methods: []string{http.MethodTrace},
			},
		),
		gen(
			"GET/HEAD/POST",
			[]string{cndInputGet, cndInputHead, cndInputPost},
			[]string{actCheckGet, actCheckHead, actCheckPost},
			&condition{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			methods := utilhttp.Methods(tt.C().methods)
			testutil.Diff(t, tt.A().methods, methods)
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	actCheckPattern := tb.Action("check pattern", "check the returned pattern")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no pattern",
			[]string{},
			[]string{actCheckPattern},
			&condition{
				h: &utilhttp.HandlerBase{
					AcceptPatterns: []string{},
				},
			},
			&action{
				patterns: []string{},
			},
		),
		gen(
			"empty pattern",
			[]string{},
			[]string{actCheckPattern},
			&condition{
				h: &utilhttp.HandlerBase{
					AcceptPatterns: []string{""},
				},
			},
			&action{
				patterns: []string{""},
			},
		),
		gen(
			"non empty pattern",
			[]string{},
			[]string{actCheckPattern},
			&condition{
				h: &utilhttp.HandlerBase{
					AcceptPatterns: []string{"/test"},
				},
			},
			&action{
				patterns: []string{"/test"},
			},
		),
		gen(
			"multiple patterns",
			[]string{},
			[]string{actCheckPattern},
			&condition{
				h: &utilhttp.HandlerBase{
					AcceptPatterns: []string{"/test1", "/test2"},
				},
			},
			&action{
				patterns: []string{"/test1", "/test2"},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			p := tt.C().h.Patterns()
			testutil.Diff(t, tt.A().patterns, p)
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

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	actCheckMethods := tb.Action("check methods", "check the returned methods")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no methods",
			[]string{},
			[]string{actCheckMethods},
			&condition{
				h: &utilhttp.HandlerBase{},
			},
			&action{
				methods: nil,
			},
		),
		gen(
			"no methods",
			[]string{},
			[]string{actCheckMethods},
			&condition{
				h: &utilhttp.HandlerBase{
					AcceptMethods: []string{""},
				},
			},
			&action{
				methods: []string{""},
			},
		),
		gen(
			"one method",
			[]string{},
			[]string{actCheckMethods},
			&condition{
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
			"multiple methods",
			[]string{},
			[]string{actCheckMethods},
			&condition{
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

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ms := tt.C().h.Methods()
			testutil.Diff(t, tt.A().methods, ms)
		})
	}
}
