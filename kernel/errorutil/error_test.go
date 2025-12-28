// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package errorutil_test

import (
	"errors"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/errorutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNew(t *testing.T) {
	type condition struct {
		code  string
		kind  string
		msg   string
		stack []byte
		err   error
	}

	type action struct {
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"zero values",
			&condition{
				code:  "",
				kind:  "",
				msg:   "",
				stack: nil,
				err:   nil,
			},
			&action{},
		),
		gen(
			"non-zero values",
			&condition{
				code:  "code",
				kind:  "kind",
				msg:   "msg",
				stack: []byte("test"),
				err:   errors.New("test"),
			},
			&action{},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := errorutil.New(tt.C.code, tt.C.kind, tt.C.msg, tt.C.stack, tt.C.err)

			testutil.Diff(t, tt.C.code, err.Code())
			testutil.Diff(t, tt.C.kind, err.Kind())
			testutil.Diff(t, tt.C.msg, err.Error())
			testutil.Diff(t, string(tt.C.stack), err.StackTrace())
			testutil.Diff(t, tt.C.err, err.Unwrap(), cmpopts.EquateErrors())
		})
	}
}
