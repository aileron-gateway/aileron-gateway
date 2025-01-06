package io_test

import (
	"io/fs"
	"os"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/io"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestReadWriteTest(t *testing.T) {
	type condition struct {
		open   func(string, int, fs.FileMode) (*os.File, error)
		remove func(string) error
	}

	type action struct {
		err any // error or errorutil.Kind
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndOpenSuccess := tb.Condition("open success", "open function successfully exit")
	cndRemoveSuccess := tb.Condition("remove success", "remove function successfully exit")
	actCheckError := tb.Action("error", "check that there is an error")
	actCheckNoError := tb.Action("no error", "check that there is no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"open success",
			[]string{cndOpenSuccess, cndRemoveSuccess},
			[]string{actCheckNoError},
			&condition{
				open:   os.OpenFile,
				remove: os.Remove,
			},
			&action{
				err: nil,
			},
		),
		gen(
			"open error",
			[]string{cndRemoveSuccess},
			[]string{actCheckError},
			&condition{
				open: func(_ string, _ int, _ fs.FileMode) (*os.File, error) {
					return nil, os.ErrPermission
				},
				remove: os.Remove,
			},
			&action{
				err: os.ErrPermission,
			},
		),
		gen(
			"remove error",
			[]string{cndOpenSuccess},
			[]string{actCheckError},
			&condition{
				open: os.OpenFile,
				remove: func(s string) error {
					os.Remove(s)
					return os.ErrPermission
				},
			},
			&action{
				err: os.ErrPermission,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmpOpen := io.OpenFile
			io.OpenFile = tt.C().open
			defer func() {
				io.OpenFile = tmpOpen
			}()

			tmpRemove := io.Remove
			io.Remove = tt.C().remove
			defer func() {
				io.Remove = tmpRemove
			}()

			err := io.ReadWriteTest(os.TempDir())
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}
