// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package http

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/errorutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testErrorHandler struct {
	core.ErrorHandler
	id string
}

func TestSetGlobalErrorHandler(t *testing.T) {
	type condition struct {
		setHandler bool
		eh         core.ErrorHandler
		name       string
	}

	type action struct {
		expect core.ErrorHandler
	}

	CndSetNil := "set nil"
	CndSetNonNil := "set non-nil"
	CndDefaultName := "default name"
	ActCheckReplaced := "check replaced"
	ActCheckStored := "check stored"
	ActCheckNil := "check nil"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndSetNil, "set nil as error handler")
	tb.Condition(CndSetNonNil, "set non-nil error handler as input")
	tb.Condition(CndDefaultName, "set error handler by default name")
	tb.Action(ActCheckReplaced, "check that error handler is replaced")
	tb.Action(ActCheckStored, "check that the error handler is stored")
	tb.Action(ActCheckNil, "check that the returned value is nil")
	table := tb.Build()

	testEH := &testErrorHandler{
		ErrorHandler: nil,
		id:           "test",
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil by default name",
			[]string{CndSetNil, CndDefaultName},
			[]string{ActCheckReplaced},
			&condition{
				setHandler: true,
				eh:         nil,
				name:       DefaultErrorHandlerName,
			},
			&action{
				expect: &DefaultErrorHandler{LG: log.GlobalLogger(log.DefaultLoggerName)},
			},
		),
		gen(
			"nil by not default name",
			[]string{CndSetNil},
			[]string{ActCheckNil},
			&condition{
				setHandler: true,
				eh:         nil,
				name:       "test",
			},
			&action{
				expect: nil,
			},
		),
		gen(
			"nil by empty name",
			[]string{CndSetNil},
			[]string{ActCheckNil},
			&condition{
				setHandler: true,
				eh:         nil,
				name:       "",
			},
			&action{
				expect: nil,
			},
		),
		gen(
			"non-nil by default name",
			[]string{CndSetNonNil, CndDefaultName},
			[]string{ActCheckReplaced},
			&condition{
				setHandler: true,
				eh:         testEH,
				name:       DefaultErrorHandlerName,
			},
			&action{
				expect: testEH,
			},
		),
		gen(
			"non-nil by not default name",
			[]string{CndSetNonNil},
			[]string{},
			&condition{
				setHandler: true,
				eh:         testEH,
				name:       "test",
			},
			&action{
				expect: testEH,
			},
		),
		gen(
			"non-nil by empty name",
			[]string{CndSetNonNil},
			[]string{},
			&condition{
				setHandler: true,
				eh:         testEH,
				name:       "",
			},
			&action{
				expect: testEH,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := GlobalErrorHandler(DefaultErrorHandlerName)
			defer func() {
				SetGlobalErrorHandler(tt.C().name, nil)
				SetGlobalErrorHandler(DefaultErrorHandlerName, tmp)
			}()

			if tt.C().setHandler {
				SetGlobalErrorHandler(tt.C().name, tt.C().eh)
			}

			eh := GlobalErrorHandler(tt.C().name)

			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[log.Logger]),
				cmp.AllowUnexported(DefaultErrorHandler{}),
			}

			if v, ok := tt.A().expect.(*testErrorHandler); ok {
				testutil.Diff(t, v.id, eh.(*testErrorHandler).id)
			} else {
				testutil.Diff(t, tt.A().expect, eh, opts...)
			}
		})
	}
}

func TestGlobalLogger(t *testing.T) {
	type condition struct {
		name string
	}

	type action struct {
		expect core.ErrorHandler
	}

	CndLoggerExist := "error handler exists"
	CndLoggerNotExist := "error handler not exists"
	CndDefaultName := "default name"
	ActCheckNonNil := "check non-nil"
	ActCheckNil := "check nil"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndLoggerExist, "get error handler which exists in the global error handler holder")
	tb.Condition(CndLoggerNotExist, "get error handler which does not exist in the global error handler holder")
	tb.Condition(CndDefaultName, "set error handler by default name")
	tb.Action(ActCheckNonNil, "check that the returned value is non-nil")
	tb.Action(ActCheckNil, "check that the returned value is nil")
	table := tb.Build()

	testEH := &testErrorHandler{
		ErrorHandler: nil,
		id:           "test",
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"default name",
			[]string{CndLoggerExist, CndDefaultName},
			[]string{ActCheckNonNil},
			&condition{
				name: DefaultErrorHandlerName,
			},
			&action{
				expect: &DefaultErrorHandler{LG: log.GlobalLogger(log.DefaultLoggerName)},
			},
		),
		gen(
			"not default name",
			[]string{CndLoggerExist},
			[]string{ActCheckNonNil},
			&condition{
				name: "test_error_handler",
			},
			&action{
				expect: testEH,
			},
		),
		gen(
			"not-nil error handler",
			[]string{CndLoggerNotExist},
			[]string{ActCheckNil},
			&condition{
				name: "not_exist_handler_name",
			},
			&action{
				expect: nil,
			},
		),
		gen(
			"not-nil error handler by empty name",
			[]string{CndLoggerNotExist},
			[]string{ActCheckNil},
			&condition{
				name: "",
			},
			&action{
				expect: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			SetGlobalErrorHandler("test_error_handler", testEH)

			eh := GlobalErrorHandler(tt.C().name)

			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[log.Logger]),
				cmp.AllowUnexported(DefaultErrorHandler{}),
			}

			if v, ok := tt.A().expect.(*testErrorHandler); ok {
				testutil.Diff(t, v.id, eh.(*testErrorHandler).id)
			} else {
				testutil.Diff(t, tt.A().expect, eh, opts...)
			}
		})
	}
}

func TestErrorHandler(t *testing.T) {
	type condition struct {
		ref *k.Reference
	}

	type action struct {
		eh  core.ErrorHandler
		err error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNilReference := tb.Condition("nil reference", "input nil as reference")
	cndNonNilReference := tb.Condition("non-nil reference", "input non-nil as reference")
	actCheckHandler := tb.Action("check handler", "check that the returned handler is non-nil and is the expected handler type")
	actCheckNoError := tb.Action("no error", "check that no error was returned")
	actCheckError := tb.Action("error", "check that an error was returned")
	table := tb.Build()

	noopEH := &testErrorHandler{id: "noop"}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"non-nil reference",
			[]string{cndNonNilReference},
			[]string{actCheckHandler, actCheckNoError},
			&condition{
				ref: testResourceRef("noop"),
			},
			&action{
				eh:  noopEH,
				err: nil,
			},
		),
		gen(
			"nil reference",
			[]string{cndNilReference},
			[]string{actCheckHandler, actCheckNoError},
			&condition{
				ref: nil,
			},
			&action{
				eh:  &DefaultErrorHandler{LG: log.GlobalLogger(log.DefaultLoggerName)},
				err: nil,
			},
		),
		gen(
			"not exists",
			[]string{cndNonNilReference},
			[]string{actCheckError},
			&condition{
				ref: testResourceRef("this is not exist"),
			},
			&action{
				eh: nil,
				err: &er.Error{
					Package:     api.ErrPkg,
					Type:        api.ErrTypeUtil,
					Description: api.ErrDscAssert,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// Prepare an api for test.
			a := api.NewContainerAPI()
			postTestResource(a, "noop", noopEH)

			eh, err := ErrorHandler(a, tt.C().ref)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())

			opts := []cmp.Option{
				cmp.Comparer(testutil.ComparePointer[log.Logger]),
				cmp.AllowUnexported(DefaultErrorHandler{}),
				cmp.AllowUnexported(testErrorHandler{}),
			}
			testutil.Diff(t, tt.A().eh, eh, opts...)
		})
	}
}

func TestNewErrorMessage(t *testing.T) {
	type condition struct {
		spec *v1.ErrorMessageSpec
	}

	type action struct {
		em  *ErrorMessage
		err error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndNonNilReference := tb.Condition("non-nil reference", "input non-nil as reference")
	actCheckNil := tb.Action("check nil", "check nil was returned for message")
	actCheckMessage := tb.Action("check message", "check the returned message is not nil and have the expected values")
	actCheckNoError := tb.Action("no error", "check that no error was returned")
	actCheckError := tb.Action("error", "check that an error was returned")
	table := tb.Build()

	tpl, _ := txtutil.NewTemplate(txtutil.TplText, "test")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil spec",
			[]string{},
			[]string{actCheckNil, actCheckNoError},
			&condition{
				spec: nil,
			},
			&action{
				em:  nil,
				err: nil,
			},
		),
		gen(
			"no mime contents",
			[]string{cndNonNilReference},
			[]string{actCheckNil, actCheckNoError},
			&condition{
				spec: &v1.ErrorMessageSpec{},
			},
			&action{
				em:  nil,
				err: nil,
			},
		),
		gen(
			"successful",
			[]string{cndNonNilReference},
			[]string{actCheckMessage, actCheckNoError},
			&condition{
				spec: &v1.ErrorMessageSpec{
					Codes:          []string{"E0001"},
					Kinds:          []string{"ErrTest"},
					Messages:       []string{".*"},
					HeaderTemplate: map[string]string{"Code": "{{code}}"},
					MIMEContents: []*v1.MIMEContentSpec{
						{
							TemplateType: k.TemplateType_Text,
							Template:     "test",
							MIMEType:     "text/plain",
						},
					},
				},
			},
			&action{
				em: &ErrorMessage{
					codes: []string{"E0001"},
					kinds: []string{"ErrTest"},
					messages: []*regexp.Regexp{
						regexp.MustCompile(`.*`),
					},
					headerTpl: map[string]*txtutil.FastTemplate{
						"Code": txtutil.NewFastTemplate("{{code}}", "{{", "}}"),
					},
					contents: []*MIMEContent{
						{
							Template:   tpl,
							MIMEType:   "text/plain",
							StatusCode: http.StatusOK,
							Header:     http.Header{},
						},
					},
				},
				err: nil,
			},
		),
		gen(
			"message compile error",
			[]string{cndNonNilReference},
			[]string{actCheckNil, actCheckError},
			&condition{
				spec: &v1.ErrorMessageSpec{
					Messages:     []string{"[0-9a-"},
					MIMEContents: []*v1.MIMEContentSpec{nil},
				},
			},
			&action{
				em: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeErrHandler,
					Description: ErrDscRegexp,
				},
			},
		),
		gen(
			"mime content create error",
			[]string{cndNonNilReference},
			[]string{actCheckNil, actCheckError},
			&condition{
				spec: &v1.ErrorMessageSpec{
					Codes: []string{"E0001"},
					MIMEContents: []*v1.MIMEContentSpec{
						{
							MIMEType: "",
						},
					},
				},
			},
			&action{
				em: nil,
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeMime,
					Description: ErrDscParseMime,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			em, err := NewErrorMessage(tt.C().spec)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())

			opts := []cmp.Option{
				cmp.AllowUnexported(ErrorMessage{}),
				cmp.AllowUnexported(MIMEContent{}),
				cmp.AllowUnexported(regexp.Regexp{}),
				cmp.AllowUnexported(txtutil.FastTemplate{}),
				cmpopts.IgnoreInterfaces(struct{ txtutil.Template }{}),
			}
			testutil.Diff(t, tt.A().em, em, opts...)
		})
	}
}

func TestErrorMessage_Match(t *testing.T) {
	type condition struct {
		em   *ErrorMessage
		code string
		kind string
		msg  string
	}

	type action struct {
		matched bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndCode := tb.Condition("input code", "input non-empty code")
	cndKind := tb.Condition("input kind", "input non-empty kind")
	cndMsg := tb.Condition("input message", "input non-empty message")
	cndExactMatch := tb.Condition("exact match", "expect exact match")
	cndPathMatch := tb.Condition("path match", "expect path match")
	actCheckMatched := tb.Action("check matched", "check that the input matched to the message")
	actCheckNotMatched := tb.Action("check not matched", "check that the input did not match to the message")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"code exact match",
			[]string{cndCode, cndExactMatch},
			[]string{actCheckMatched},
			&condition{
				em: &ErrorMessage{
					codes: []string{"E0001"},
				},
				code: "E0001",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"code exact match / multiple value",
			[]string{cndCode, cndExactMatch},
			[]string{actCheckMatched},
			&condition{
				em: &ErrorMessage{
					codes: []string{"E0002", "E0001"},
				},
				code: "E0001",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"code path match",
			[]string{cndCode, cndPathMatch},
			[]string{actCheckMatched},
			&condition{
				em: &ErrorMessage{
					codes: []string{"E000*"},
				},
				code: "E0001",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"code path match / multiple value",
			[]string{cndCode, cndPathMatch},
			[]string{actCheckMatched},
			&condition{
				em: &ErrorMessage{
					codes: []string{"E0002", "E000*"},
				},
				code: "E0001",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"kind exact match",
			[]string{cndKind, cndExactMatch},
			[]string{actCheckMatched},
			&condition{
				em: &ErrorMessage{
					kinds: []string{"ErrTest"},
				},
				kind: "ErrTest",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"kind exact match / multiple value",
			[]string{cndKind, cndExactMatch},
			[]string{actCheckMatched},
			&condition{
				em: &ErrorMessage{
					kinds: []string{"ErrDummy", "ErrTest"},
				},
				kind: "ErrTest",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"kind path match",
			[]string{cndKind, cndPathMatch},
			[]string{actCheckMatched},
			&condition{
				em: &ErrorMessage{
					kinds: []string{"ErrTe*"},
				},
				kind: "ErrTest",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"kind path match / multiple value",
			[]string{cndKind, cndPathMatch},
			[]string{actCheckMatched},
			&condition{
				em: &ErrorMessage{
					kinds: []string{"ErrDum*", "ErrTe*"},
				},
				kind: "ErrTest",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"message match",
			[]string{cndMsg},
			[]string{actCheckMatched},
			&condition{
				em: &ErrorMessage{
					messages: []*regexp.Regexp{
						regexp.MustCompile(`test error`),
					},
				},
				msg: "This is a test error message.",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"message match / multiple value",
			[]string{cndMsg},
			[]string{actCheckMatched},
			&condition{
				em: &ErrorMessage{
					messages: []*regexp.Regexp{
						regexp.MustCompile(`not match`),
						regexp.MustCompile(`test error`),
					},
				},
				msg: "This is a test error message.",
			},
			&action{
				matched: true,
			},
		),
		gen(
			"not matched",
			[]string{cndCode, cndKind, cndMsg},
			[]string{actCheckNotMatched},
			&condition{
				em:   &ErrorMessage{},
				code: "E0001",
				kind: "ErrTest",
				msg:  "This is a test error message.",
			},
			&action{
				matched: false,
			},
		),
		gen(
			"not matched",
			[]string{cndCode, cndKind, cndMsg},
			[]string{actCheckNotMatched},
			&condition{
				em: &ErrorMessage{
					codes: []string{"E0002"},
					kinds: []string{"ErrDum*"},
					messages: []*regexp.Regexp{
						regexp.MustCompile(`not match`),
					},
				},
				code: "E0001",
				kind: "ErrTest",
				msg:  "This is a test error message.",
			},
			&action{
				matched: false,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			matched := tt.C().em.Match(tt.C().code, tt.C().kind, []byte(tt.C().msg))
			testutil.Diff(t, tt.A().matched, matched)
		})
	}
}

func TestErrorMessage_Content(t *testing.T) {
	type condition struct {
		em     *ErrorMessage
		accept string
	}

	type action struct {
		content *MIMEContent
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndMatchFirst := tb.Condition("accept first", "input accept to match the first content")
	cndMatchSecond := tb.Condition("accept second", "input accept to match the second content")
	actCheckMatched := tb.Action("check matched", "check that the input matched to a message")
	table := tb.Build()

	tpl1, _ := txtutil.NewTemplate(txtutil.TplText, "test1")
	tpl2, _ := txtutil.NewTemplate(txtutil.TplText, "test2")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"match to first",
			[]string{cndMatchFirst},
			[]string{actCheckMatched},
			&condition{
				em: &ErrorMessage{
					contents: []*MIMEContent{
						{
							Template: tpl1,
							MIMEType: "application/json",
						},
						{
							Template: tpl2,
							MIMEType: "text/plain",
						},
					},
				},
				accept: "application/json; charset=utf-8",
			},
			&action{
				content: &MIMEContent{
					Template: tpl1,
					MIMEType: "application/json",
				},
			},
		),
		gen(
			"match to second",
			[]string{cndMatchSecond},
			[]string{actCheckMatched},
			&condition{
				em: &ErrorMessage{
					contents: []*MIMEContent{
						{
							Template: tpl1,
							MIMEType: "application/json",
						},
						{
							Template: tpl2,
							MIMEType: "text/plain",
						},
					},
				},
				accept: "text/plain; charset=utf-8",
			},
			&action{
				content: &MIMEContent{
					Template: tpl2,
					MIMEType: "text/plain",
				},
			},
		),
		gen(
			"complex accept",
			[]string{cndMatchSecond},
			[]string{actCheckMatched},
			&condition{
				em: &ErrorMessage{
					contents: []*MIMEContent{
						{
							Template: tpl1,
							MIMEType: "application/json",
						},
						{
							Template: tpl2,
							MIMEType: "text/plain",
						},
					},
				},
				accept: "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,application/json,*/*;q=0.8",
			},
			&action{
				content: &MIMEContent{
					Template: tpl1,
					MIMEType: "application/json",
				},
			},
		),
		gen(
			"match to nothing",
			[]string{},
			[]string{},
			&condition{
				em: &ErrorMessage{
					contents: []*MIMEContent{
						{
							Template: tpl1,
							MIMEType: "application/json",
						},
						{
							Template: tpl2,
							MIMEType: "text/plain",
						},
					},
				},
				accept: "application/xml; charset=utf-8",
			},
			&action{
				content: &MIMEContent{
					Template: tpl1,
					MIMEType: "application/json",
				},
			},
		),
		gen(
			"no content",
			[]string{},
			[]string{},
			&condition{
				em: &ErrorMessage{
					contents: []*MIMEContent{},
				},
				accept: "application/json",
			},
			&action{
				content: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			content := tt.C().em.Content(tt.C().accept)

			opts := []cmp.Option{
				cmp.AllowUnexported(MIMEContent{}),
				cmpopts.IgnoreInterfaces(struct{ txtutil.Template }{}),
			}
			testutil.Diff(t, tt.A().content, content, opts...)
		})
	}
}

func TestDefaultErrorHandler_ServeHTTPError(t *testing.T) {
	type condition struct {
		eh  *DefaultErrorHandler
		err error
	}

	type action struct {
		status int
		header http.Header
		body   *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndMatchMessage := tb.Condition("match message", "request matches to a message")
	actCheckStatus := tb.Action("check stats", "check the response status code")
	actCheckOverwritten := tb.Action("check overwritten", "check the overwritten response body")
	table := tb.Build()

	testErrKind := errorutil.NewKind("E0001", "ErrTest", "This is a test error kind")
	testErr := testErrKind.WithoutStack(nil, nil)
	tpl, _ := txtutil.NewTemplate(txtutil.TplGoText, "{{.code}}.{{.kind}}")

	debugLogger := log.NewJSONSLogger(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no content/nil error",
			[]string{},
			[]string{actCheckStatus},
			&condition{
				eh: &DefaultErrorHandler{
					LG:   debugLogger,
					Msgs: []*ErrorMessage{},
				},
				err: nil,
			},
			&action{
				status: http.StatusInternalServerError,
				body:   regexp.MustCompile(`Internal Server Error`),
			},
		),
		gen(
			"1 content/nil error",
			[]string{},
			[]string{actCheckStatus},
			&condition{
				eh: &DefaultErrorHandler{
					LG: debugLogger,
					Msgs: []*ErrorMessage{
						{
							codes:    []string{"*"},
							contents: []*MIMEContent{},
						},
					},
				},
				err: nil,
			},
			&action{
				status: http.StatusInternalServerError,
				body:   regexp.MustCompile(`Internal Server Error`),
			},
		),
		gen(
			"primitive error",
			[]string{},
			[]string{actCheckStatus},
			&condition{
				eh: &DefaultErrorHandler{
					LG: debugLogger,
					Msgs: []*ErrorMessage{
						{
							codes:    []string{"*"},
							contents: []*MIMEContent{},
						},
					},
				},
				err: errors.New("test error"),
			},
			&action{
				status: http.StatusInternalServerError,
				body:   regexp.MustCompile(`Internal Server Error`),
			},
		),
		gen(
			"errorutil error",
			[]string{},
			[]string{actCheckStatus},
			&condition{
				eh: &DefaultErrorHandler{
					LG: debugLogger,
					Msgs: []*ErrorMessage{
						{
							codes:    []string{"*"},
							contents: []*MIMEContent{},
						},
					},
				},
				err: testErr,
			},
			&action{
				status: http.StatusInternalServerError,
				body:   regexp.MustCompile(`Internal Server Error`),
			},
		),
		gen(
			"match / http error including primitive error",
			[]string{cndMatchMessage},
			[]string{actCheckStatus, actCheckOverwritten},
			&condition{
				eh: &DefaultErrorHandler{
					LG: debugLogger,
					Msgs: []*ErrorMessage{
						{
							codes: []string{"E0001"},
							contents: []*MIMEContent{
								{
									StatusCode: http.StatusForbidden,
									Header:     http.Header{"Foo": []string{"bar"}},
									Template:   tpl,
									MIMEType:   "application/json",
								},
							},
						},
					},
				},
				err: NewHTTPError(errors.New("test error"), http.StatusInternalServerError),
			},
			&action{
				status: http.StatusInternalServerError,
				body:   regexp.MustCompile(`Internal Server Error`),
			},
		),
		gen(
			"match / http error including errorutil error",
			[]string{cndMatchMessage},
			[]string{actCheckStatus, actCheckOverwritten},
			&condition{
				eh: &DefaultErrorHandler{
					LG: debugLogger,
					Msgs: []*ErrorMessage{
						{
							codes: []string{"E0001"},
							contents: []*MIMEContent{
								{
									StatusCode: http.StatusForbidden,
									Header:     http.Header{"Foo": []string{"bar"}},
									Template:   tpl,
									MIMEType:   "application/json",
								},
							},
						},
					},
				},
				err: NewHTTPError(testErr, http.StatusInternalServerError),
			},
			&action{
				status: http.StatusForbidden,
				header: http.Header{"Foo": []string{"bar"}},
				body:   regexp.MustCompile(`^E0001.ErrTest$`),
			},
		),
		gen(
			"400 error / http error",
			[]string{cndMatchMessage},
			[]string{actCheckStatus, actCheckOverwritten},
			&condition{
				eh: &DefaultErrorHandler{
					LG: debugLogger,
					Msgs: []*ErrorMessage{
						{
							codes:    []string{"*"},
							contents: []*MIMEContent{},
						},
					},
				},
				err: NewHTTPError(errors.New("test error"), http.StatusBadRequest),
			},
			&action{
				status: http.StatusBadRequest,
				body:   regexp.MustCompile(`Bad Request`),
			},
		),
		gen(
			"400 error / http error with errorutil error",
			[]string{cndMatchMessage},
			[]string{actCheckStatus, actCheckOverwritten},
			&condition{
				eh: &DefaultErrorHandler{
					LG: debugLogger,
					Msgs: []*ErrorMessage{
						{
							codes: []string{"E0001"},
							contents: []*MIMEContent{
								{
									StatusCode: http.StatusBadRequest,
									Header:     http.Header{"Foo": []string{"bar"}},
									Template:   tpl,
									MIMEType:   "application/json",
								},
							},
						},
					},
				},
				err: NewHTTPError(testErr, http.StatusInternalServerError),
			},
			&action{
				status: http.StatusBadRequest,
				header: http.Header{"Foo": []string{"bar"}},
				body:   regexp.MustCompile(`^E0001.ErrTest$`),
			},
		),
		gen(
			"template body when 400 error + stack always",
			[]string{cndMatchMessage},
			[]string{actCheckStatus, actCheckOverwritten},
			&condition{
				eh: &DefaultErrorHandler{
					LG:          debugLogger,
					StackAlways: true,
					Msgs: []*ErrorMessage{
						{
							codes: []string{"E0001"},
							contents: []*MIMEContent{
								{
									StatusCode: http.StatusBadRequest,
									Header:     http.Header{"Foo": []string{"bar"}},
									Template:   tpl,
									MIMEType:   "application/json",
								},
							},
						},
					},
				},
				err: NewHTTPError(testErr, http.StatusBadRequest),
			},
			&action{
				status: http.StatusBadRequest,
				header: http.Header{"Foo": []string{"bar"}},
				body:   regexp.MustCompile(`^E0001.ErrTest$`),
			},
		),
		gen(
			"template header",
			[]string{cndMatchMessage},
			[]string{actCheckStatus, actCheckOverwritten},
			&condition{
				eh: &DefaultErrorHandler{
					LG: debugLogger,
					Msgs: []*ErrorMessage{
						{
							codes: []string{"E0001"},
							headerTpl: map[string]*txtutil.FastTemplate{
								"Status": txtutil.NewFastTemplate("{{status}}", "{{", "}}"),
							},
							contents: []*MIMEContent{
								{
									StatusCode: http.StatusBadRequest,
									Header:     http.Header{"Foo": []string{"bar"}},
									Template:   tpl,
									MIMEType:   "application/json",
								},
							},
						},
					},
				},
				err: NewHTTPError(testErr, http.StatusBadRequest),
			},
			&action{
				status: http.StatusBadRequest,
				header: http.Header{"Foo": {"bar"}, "Status": {"400"}},
				body:   regexp.MustCompile(`^E0001.ErrTest$`),
			},
		),
		gen(
			"logging only status -1",
			[]string{cndMatchMessage},
			[]string{actCheckStatus, actCheckOverwritten},
			&condition{
				eh: &DefaultErrorHandler{
					LG:          debugLogger,
					StackAlways: true,
					Msgs:        []*ErrorMessage{},
				},
				err: NewHTTPError(testErr, -1), // LoggingOnly
			},
			&action{
				status: 200, // LoggingOnly does not write into response. It becomes 200 by the specification of response writer.
				body:   regexp.MustCompile(`^$`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "http://test.com/foo", nil)
			r.Header.Set("Accept", "text/html,application/xml,application/json;q=0.9,*/*;q=0.8")

			tt.C().eh.ServeHTTPError(w, r, tt.C().err)

			resp := w.Result()
			defer resp.Body.Close()
			b, _ := io.ReadAll(resp.Body)

			t.Log(string(b) + "\n")
			testutil.Diff(t, tt.A().status, resp.StatusCode)
			testutil.Diff(t, true, tt.A().body.Match(b))
			for k, v := range tt.A().header {
				testutil.Diff(t, v, resp.Header[k])
			}
		})
	}
}

func postTestResource(server api.API[*api.Request, *api.Response], name string, res any) {
	ref := testResourceRef(name)
	req := &api.Request{
		Method:  api.MethodPost,
		Key:     ref.APIVersion + "/" + ref.Kind + "/" + ref.Namespace + "/" + ref.Name,
		Content: res,
	}
	if _, err := server.Serve(context.Background(), req); err != nil {
		panic(err)
	}
}

func testResourceRef(name string) *k.Reference {
	return &k.Reference{
		APIVersion: "core/v1",
		Kind:       "Container",
		Namespace:  "test",
		Name:       name,
	}
}
