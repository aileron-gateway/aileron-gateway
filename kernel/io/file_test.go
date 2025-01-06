package io_test

import (
	stdcmp "cmp"
	"os"
	"path/filepath"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/io"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// testDir is the path to the test data.
// This path can be changed by the environmental variable.
var testDir = stdcmp.Or(os.Getenv("TEST_DIR"), "../../test/")

var testDataDir = testDir + "ut/kernel/io/testdir/"

func TestReadFiles(t *testing.T) {
	type condition struct {
		paths     []string
		recursive bool
	}

	type action struct {
		contents map[string][]byte
		err      error
	}

	cndFiles := "pass file paths"
	cndDirs := "pass directory paths"
	actCheckContent := "check returned contents"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndFiles, "pass file paths")
	tb.Condition(cndDirs, "pass dir paths")
	tb.Action(actCheckContent, "check that the returned content is the same as expected")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that the returned error was the one expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"path 2 file path",
			[]string{cndFiles},
			[]string{actCheckContent, actCheckNoError},
			&condition{
				paths: []string{
					filepath.Clean(testDataDir + "file1.txt"),
				},
			},
			&action{
				contents: map[string][]byte{
					filepath.Clean(testDataDir + "file1.txt"): []byte("test"),
				},
				err: nil,
			},
		),
		gen(
			"path 2 file paths",
			[]string{cndFiles},
			[]string{actCheckContent, actCheckNoError},
			&condition{
				paths: []string{
					filepath.Clean(testDataDir + "file1.txt"),
					filepath.Clean(testDataDir + "file2.txt"),
				},
			},
			&action{
				contents: map[string][]byte{
					filepath.Clean(testDataDir + "file1.txt"): []byte("test"),
					filepath.Clean(testDataDir + "file2.txt"): []byte("test"),
				},
				err: nil,
			},
		),
		gen(
			"pass 1 dir path",
			[]string{cndDirs},
			[]string{actCheckContent, actCheckNoError},
			&condition{
				paths: []string{
					testDataDir,
				},
			},
			&action{
				contents: map[string][]byte{
					filepath.Clean(testDataDir + "empty.txt"): []byte(""),
					filepath.Clean(testDataDir + "file1.txt"): []byte("test"),
					filepath.Clean(testDataDir + "file2.txt"): []byte("test"),
				},
				err: nil,
			},
		),
		gen(
			"pass 2 dir paths",
			[]string{cndDirs},
			[]string{actCheckContent, actCheckNoError},
			&condition{
				paths: []string{
					testDataDir,
					testDataDir + "subdir/",
				},
			},
			&action{
				contents: map[string][]byte{
					filepath.Clean(testDataDir + "empty.txt"):        []byte(""),
					filepath.Clean(testDataDir + "file1.txt"):        []byte("test"),
					filepath.Clean(testDataDir + "file2.txt"):        []byte("test"),
					filepath.Clean(testDataDir + "subdir/file3.txt"): []byte("test"),
				},
				err: nil,
			},
		),
		gen(
			"pass nothing",
			[]string{},
			[]string{actCheckContent, actCheckNoError},
			&condition{
				paths: []string{},
			},
			&action{
				contents: map[string][]byte{},
				err:      nil,
			},
		),
		gen(
			"pass empty string",
			[]string{},
			[]string{actCheckContent, actCheckNoError},
			&condition{
				paths: []string{
					"",
				},
			},
			&action{
				contents: map[string][]byte{},
				err:      nil,
			},
		),
		gen(
			"pass 1 dir path recursive",
			[]string{cndDirs},
			[]string{actCheckContent, actCheckNoError},
			&condition{
				paths: []string{
					testDataDir,
				},
				recursive: true,
			},
			&action{
				contents: map[string][]byte{
					filepath.Clean(testDataDir + "empty.txt"):        []byte(""),
					filepath.Clean(testDataDir + "file1.txt"):        []byte("test"),
					filepath.Clean(testDataDir + "file2.txt"):        []byte("test"),
					filepath.Clean(testDataDir + "subdir/file3.txt"): []byte("test"),
				},
				err: nil,
			},
		),
		gen(
			"pass dir path which does not exist",
			[]string{cndDirs},
			[]string{actCheckContent, actCheckError},
			&condition{
				paths: []string{
					"not_exist/",
				},
			},
			&action{
				contents: nil,
				err: &er.Error{
					Package:     io.ErrPkg,
					Type:        io.ErrTypeFile,
					Description: io.ErrDscListFile,
				},
			},
		),
		gen(
			"pass file path which does not exist",
			[]string{cndFiles},
			[]string{actCheckContent, actCheckError},
			&condition{
				paths: []string{
					"not_exist",
				},
			},
			&action{
				contents: nil,
				err: &er.Error{
					Package:     io.ErrPkg,
					Type:        io.ErrTypeFile,
					Description: io.ErrDscListFile,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			contents, err := io.ReadFiles(tt.C().recursive, tt.C().paths...)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().contents, contents)
		})
	}
}

func TestSplitMultiDoc(t *testing.T) {
	type condition struct {
		docs []byte
		sep  string
	}

	type action struct {
		contents [][]byte
	}

	cndEmptyDocs := "input empty docs as file path"
	cndMultipleDocs := "input empty docs as file path"
	cndDefaultSeparator := "use default separator"
	actCheckContent := "check content"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndEmptyDocs, "input empty docs")
	tb.Condition(cndMultipleDocs, "input docs with multiple blocks separated by the separator")
	tb.Condition(cndDefaultSeparator, "docs are separated by the default separator")
	tb.Action(actCheckContent, "check that the content is the same as expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"read docs with 1 block",
			[]string{cndDefaultSeparator},
			[]string{actCheckContent},
			&condition{
				docs: []byte("test"),
				sep:  "",
			},
			&action{
				contents: [][]byte{
					[]byte("test"),
				},
			},
		),
		gen(
			"read docs with 2 blocks",
			[]string{cndDefaultSeparator, cndMultipleDocs},
			[]string{actCheckContent},
			&condition{
				docs: []byte("test1\n---\ntest2"),
				sep:  "",
			},
			&action{
				contents: [][]byte{
					[]byte("test1"),
					[]byte("test2"),
				},
			},
		),
		gen(
			"read empty docs",
			[]string{cndDefaultSeparator, cndEmptyDocs},
			[]string{actCheckContent},
			&condition{
				docs: []byte(""),
				sep:  "",
			},
			&action{
				contents: nil,
			},
		),
		gen(
			"read docs double separator",
			[]string{cndDefaultSeparator, cndMultipleDocs},
			[]string{actCheckContent},
			&condition{
				docs: []byte("test1\n---\n---\ntest2"),
				sep:  "",
			},
			&action{
				contents: [][]byte{
					[]byte("test1"),
					[]byte("test2"),
				},
			},
		),
		gen(
			"read docs with only separator",
			[]string{cndDefaultSeparator, cndMultipleDocs},
			[]string{actCheckContent},
			&condition{
				docs: []byte("---\n"),
				sep:  "",
			},
			&action{
				contents: nil,
			},
		),
		gen(
			"read docs with non default separator",
			[]string{},
			[]string{actCheckContent},
			&condition{
				docs: []byte("test1\n***\ntest2"),
				sep:  "***\n",
			},
			&action{
				contents: [][]byte{
					[]byte("test1"),
					[]byte("test2"),
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			contents := io.SplitMultiDoc(tt.C().docs, tt.C().sep)
			testutil.Diff(t, tt.A().contents, contents)
		})
	}
}
