// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package tracking

import (
	"crypto/rand"
	"errors"
	"regexp"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
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
				manifest: &v1.TrackingMiddleware{
					APIVersion: apiVersion,
					Kind:       kind,
					Metadata: &k.Metadata{
						Namespace: "",
						Name:      "default",
					},
					Spec: &v1.TrackingMiddlewareSpec{},
				},
			},
			&action{
				err: nil,
			},
		),
		gen(
			"fail to get errorhandler",
			[]string{},
			[]string{},
			&condition{
				manifest: &v1.TrackingMiddleware{
					Spec: &v1.TrackingMiddlewareSpec{
						ErrorHandler: &k.Reference{
							APIVersion: "wrong",
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create TrackingMiddleware`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			server := api.NewContainerAPI()

			a := &API{}
			_, err := a.Create(server, tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
		})
	}
}

// errorReaderMock implements io.Reader, a mock that always returns an error
type errorReaderMock struct{}

func (*errorReaderMock) Read(p []byte) (n int, err error) {
	return 0, errors.New("error reading random data")
}

func TestNewReqIDFunc(t *testing.T) {
	type condition struct {
		encodingType k.EncodingType
		mockError    bool
	}

	type action struct {
		expectedPattern *regexp.Regexp
		err             any
		errPattern      *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"base64 encoding",
			[]string{},
			[]string{},
			&condition{
				encodingType: k.EncodingType_Base64,
				mockError:    false,
			},
			&action{
				expectedPattern: regexp.MustCompile(`^[A-Za-z0-9+/]+={0,2}$`),
				err:             nil,
			},
		),
		gen(
			"error generating ID due to reader error",
			[]string{},
			[]string{},
			&condition{
				encodingType: k.EncodingType_Base64,
				mockError:    true,
			},
			&action{
				err:        &er.Error{Package: "uid", Type: "hosted id", Description: "failed to generate a new uid."},
				errPattern: regexp.MustCompile("error reading random data"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			if tt.C().mockError {
				// Replace rand.Reader with a mock that returns an error
				originalReader := rand.Reader
				rand.Reader = &errorReaderMock{}
				defer func() { rand.Reader = originalReader }()
			}

			fn := newReqIDFunc(tt.C().encodingType)
			id, err := fn()
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
			if tt.A().expectedPattern != nil {
				testutil.Diff(t, true, tt.A().expectedPattern.MatchString(id), nil)
			}
		})
	}
}

func TestNewTraceID(t *testing.T) {
	type condition struct {
		requestID string
	}

	type action struct {
		expectedTraceID string
		err             any
		errPattern      *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"valid request ID",
			[]string{},
			[]string{},
			&condition{
				requestID: "mock-request-id",
			},
			&action{
				expectedTraceID: "mock-request-id",
				err:             nil,
			},
		),
		gen(
			"empty request ID",
			[]string{},
			[]string{},
			&condition{
				requestID: "",
			},
			&action{
				expectedTraceID: "",
				err:             nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			traceID, err := newTraceID(tt.C().requestID)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
			testutil.Diff(t, tt.A().expectedTraceID, traceID)
		})
	}
}
