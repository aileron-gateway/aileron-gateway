// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package io

import (
	stdcmp "cmp"
	"compress/gzip"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// testDir is the path to the test data.
// This path can be changed by the environmental variable.
var testDir = stdcmp.Or(os.Getenv("TEST_DIR"), "../../test/")

func TestNewMatchFunc(t *testing.T) {
	type condition struct {
		base    string
		ext     string
		pattern *regexp.Regexp
	}

	type action struct {
		results map[string]string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil pattern",
			[]string{},
			[]string{},
			&condition{
				base:    "test",
				ext:     ".log",
				pattern: nil,
			},
			&action{
				results: map[string]string{
					"test.log":         " ", // Space will be returned instead of empty.
					"test..log":        ".",
					"test.foo.log":     ".foo",
					"test.foo.bar.log": ".foo.bar",
					"test-foo.log":     "-foo",
				},
			},
		),
		gen(
			"nil pattern/base not match",
			[]string{},
			[]string{},
			&condition{
				base:    "test",
				ext:     ".log",
				pattern: nil,
			},
			&action{
				results: map[string]string{
					"placeholder.log":         "",
					"placeholder..log":        "",
					"placeholder.foo.log":     "",
					"placeholder.foo.bar.log": "",
					"placeholder-foo.log":     "",
				},
			},
		),
		gen(
			"nil pattern/ext not match",
			[]string{},
			[]string{},
			&condition{
				base:    "test",
				ext:     ".log",
				pattern: nil,
			},
			&action{
				results: map[string]string{
					"test.txt":         "",
					"test..txt":        "",
					"test.foo.txt":     "",
					"test.foo.bar.txt": "",
					"test-foo.txt":     "",
				},
			},
		),
		gen(
			"nil pattern/empty base",
			[]string{},
			[]string{},
			&condition{
				base:    "",
				ext:     ".log",
				pattern: nil,
			},
			&action{
				results: map[string]string{
					"test.log":         "test", // Space will be returned instead of empty.
					"test..log":        "test.",
					"test.foo.log":     "test.foo",
					"test.foo.bar.log": "test.foo.bar",
					"test-foo.log":     "test-foo",
				},
			},
		),
		gen(
			"nil pattern/empty ext",
			[]string{},
			[]string{},
			&condition{
				base:    "test",
				ext:     "",
				pattern: nil,
			},
			&action{
				results: map[string]string{
					"test.log":         ".log", // Space will be returned instead of empty.
					"test..log":        "..log",
					"test.foo.log":     ".foo.log",
					"test.foo.bar.log": ".foo.bar.log",
					"test-foo.log":     "-foo.log",
				},
			},
		),
		gen(
			"match id",
			[]string{},
			[]string{},
			&condition{
				base:    "test",
				ext:     ".log",
				pattern: regexp.MustCompile(`\.[0-9]+`),
			},
			&action{
				results: map[string]string{
					"test.0.log":  ".0",
					"test.1.log":  ".1",
					"test.01.log": ".01",
				},
			},
		),
		gen(
			"not match id",
			[]string{},
			[]string{},
			&condition{
				base:    "test",
				ext:     ".log",
				pattern: regexp.MustCompile(`\.[0-9]+`),
			},
			&action{
				results: map[string]string{
					"test.log":     "",
					"test..log":    "",
					"test.-1.log":  "",
					"test.foo.log": "",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			matchFunc := newMatchFunc(tt.C().base, tt.C().ext, tt.C().pattern)
			for k, v := range tt.A().results {
				testutil.Diff(t, v, matchFunc(k))
			}
		})
	}
}

func TestNewParseFunc(t *testing.T) {
	type condition struct {
		layout string
		loc    *time.Location
	}

	type action struct {
		id      int
		results map[string]int64
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"invalid pattern",
			[]string{},
			[]string{},
			&condition{
				layout: "15-04-05",
			},
			&action{
				id: -1,
				results: map[string]int64{
					"":         -1,
					".":        -1,
					"..":       -1,
					"0":        -1,
					"1":        -1,
					"foo":      -1,
					".foo":     -1,
					".foo.bar": -1,
				},
			},
		),
		gen(
			"id 0",
			[]string{},
			[]string{},
			&condition{
				layout: "15-04-05",
			},
			&action{
				id: 0,
				results: map[string]int64{
					".0":   0,
					".00":  0,
					".000": 0,
				},
			},
		),
		gen(
			"id 1",
			[]string{},
			[]string{},
			&condition{
				layout: "15-04-05",
			},
			&action{
				id: 1,
				results: map[string]int64{
					".1":   0,
					".01":  0,
					".001": 0,
				},
			},
		),
		gen(
			"time stamp",
			[]string{},
			[]string{},
			&condition{
				layout: "2006-01-02_15-04-05",
				loc:    time.UTC,
			},
			&action{
				id: 1,
				results: map[string]int64{
					".1970-01-01_00-00-00.1": 0,
					".1970-01-01_00-00-01.1": 1,
					".1970-01-01_00-01-00.1": 60,
					".1970-01-01_01-00-00.1": 3600,
				},
			},
		),
		gen(
			"invalid time stamp",
			[]string{},
			[]string{},
			&condition{
				layout: "15-04-05",
			},
			&action{
				id: -1,
				results: map[string]int64{
					"..1":           -1,
					". .1":          -1,
					".foo.1":        -1,
					".00-00-.1":     -1,
					".00-00-00.foo": -1,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			parseFunc := newParseFunc(tt.C().layout, tt.C().loc)
			for k, v := range tt.A().results {
				created, id := parseFunc(k)
				testutil.Diff(t, v, created)
				testutil.Diff(t, tt.A().id, id)
			}
		})
	}
}

// func TestLogicalFileConfig_New(t *testing.T) {

// 	type condition struct {
// 		c      *LogicalFileConfig
// 		open   func(string, int, fs.FileMode) (*os.File, error)
// 		remove func(string) error
// 	}

// 	type action struct {
// 		wc         io.WriteCloser
// 		err        any // error or errorutil.Kind
// 		errPattern *regexp.Regexp
// 	}

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())
// 	table := tb.Build()

// 	// testDir := testDir+"ut/kernel/io/"

// 	gen := testutil.NewCase[*condition, *action]
// 	testCases := []*testutil.Case[*condition, *action]{
// 		gen(
// 			"zero config",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				c: &LogicalFileConfig{
// 					// SrcDir:   testDir,
// 					// FileName: "dummy.txt",
// 				},
// 				open: func(s string, i int, fm fs.FileMode) (*os.File, error) {
// 					return os.NewFile(uintptr(os.Stdout.Fd()), "/test"), nil // Stderr for dummy
// 				},
// 				remove: func(s string) error { return nil },
// 			},
// 			&action{
// 				wc: &LogicalFile{
// 					curFile:        os.Stderr,
// 					fileBase:       "application",
// 					fileExt:        ".log",
// 					fileArchiveExt: ".log",
// 					srcDir:         ".",
// 					dstDir:         ".",
// 					srcMatchFunc:   newMatchFunc("application", ".log", idOnlyPattern),
// 					parseFunc:      newParseFunc("", time.UTC), // "" will be UTC.
// 				},
// 				err: nil,
// 			},
// 		),
// 	}

// 	testutil.Register(table, testCases...)

// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {

// 			if tt.C().open != nil {
// 				tmpOpen := OpenFile
// 				OpenFile = tt.C().open
// 				defer func() {
// 					OpenFile = tmpOpen
// 				}()
// 			}
// 			if tt.C().remove != nil {
// 				tmpRemove := Remove
// 				Remove = tt.C().remove
// 				defer func() {
// 					Remove = tmpRemove
// 				}()
// 			}

// 			wc, err := tt.C().c.New()
// 			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)

// 			opts := []cmp.Option{
// 				cmp.AllowUnexported(LogicalFile{}, atomic.Int64{}, atomic.Bool{}),
// 				cmpopts.IgnoreUnexported(sync.RWMutex{}, os.File{}),
// 				cmp.Comparer(testutil.ComparePointer[func(string) string]), // matchFunc
// 				// cmp.Comparer(testutil.ComparePointer[func(string) (int64 int)]), // parseFunc
// 				cmpopts.IgnoreFields(LogicalFile{}, "parseFunc", "manageFunc"),
// 			}
// 			testutil.Diff(t, tt.A().wc, wc, opts...)

// 		})
// 	}
// }

func TestLocation(t *testing.T) {
	type condition struct {
		tz string
	}

	type action struct {
		loc *time.Location
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputUTC := tb.Condition("UTC", "input UTC as timezone")
	cndInputLocal := tb.Condition("Local", "input Local as timezone")
	cndInvalid := tb.Condition("invalid", "input invalid timezone")
	actCheckUTC := tb.Action("check UTC", "check UTC timezone returned")
	actCheckLocal := tb.Action("check Local", "check local timezone returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"UTC",
			[]string{cndInputUTC},
			[]string{actCheckUTC},
			&condition{
				tz: "UTC",
			},
			&action{
				loc: time.UTC,
			},
		),
		gen(
			"Local",
			[]string{cndInputLocal},
			[]string{actCheckLocal},
			&condition{
				tz: "Local",
			},
			&action{
				loc: time.Local,
			},
		),
		gen(
			"invalid zone",
			[]string{cndInvalid},
			[]string{actCheckLocal},
			&condition{
				tz: "INVALID",
			},
			&action{
				loc: time.Local,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			loc := location(tt.C().tz)

			opts := []cmp.Option{
				cmp.AllowUnexported(time.Location{}),
				cmpopts.IgnoreFields(time.Location{}, "cacheZone"),
				cmpopts.IgnoreFields(time.Location{}, "zone", "tx"),
			}
			testutil.Diff(t, tt.A().loc, loc, opts...)
		})
	}
}

func TestCompressFiles(t *testing.T) {
	type condition struct {
		open   func(string, int, fs.FileMode) (*os.File, error)
		remove func(string) error
		create []string
		srcDir string
		dstDir string
		level  int
	}

	type action struct {
		srcFile    string
		dstFile    string
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndSrcExists := tb.Condition("src exists", "source file exists")
	// actCheckError := tb.Action("error", "check that there is an error")
	table := tb.Build()

	testDir := testDir + "ut/kernel/io/"
	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"srcDir!=dstDir/src file exists",
			[]string{cndSrcExists},
			[]string{},
			&condition{
				create: []string{testDir + "gzipSrc/test.log"},
				srcDir: testDir + "gzipSrc",
				dstDir: testDir + "gzipDst",
				level:  gzip.BestSpeed,
			},
			&action{
				srcFile: "",
				dstFile: testDir + "gzipDst/test.log.gz",
				err:     nil,
			},
		),
		gen(
			"list file error",
			[]string{cndSrcExists},
			[]string{},
			&condition{
				create: []string{},
				srcDir: testDir + "not-exists",
				dstDir: testDir + "gzipSrc",
				level:  gzip.BestSpeed,
			},
			&action{
				srcFile: "",
				dstFile: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeFile,
					Description: ErrDscListFile,
				},
			},
		),
		gen(
			"srcDir==dstDir/no compression",
			[]string{cndSrcExists},
			[]string{},
			&condition{
				create: []string{},
				srcDir: testDir + "gzipSrc",
				dstDir: testDir + "gzipSrc",
				level:  gzip.NoCompression,
			},
			&action{
				srcFile: "",
				dstFile: "",
				err:     nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			if tt.C().open != nil {
				tmpOpen := OpenFile
				OpenFile = tt.C().open
				defer func() {
					OpenFile = tmpOpen
				}()
			}

			if tt.C().remove != nil {
				tmpRemove := Remove
				Remove = tt.C().remove
				defer func() {
					Remove = tmpRemove
				}()
			}

			for _, f := range tt.C().create {
				if err := os.WriteFile(f, []byte("test"), os.ModePerm); err != nil {
					t.Errorf("%#v\n", err)
					return
				}
			}
			defer func() {
				for _, f := range tt.C().create {
					os.Remove(f)
				}
				os.Remove(tt.A().srcFile)
				os.Remove(tt.A().dstFile)
			}()

			matchFunc := func(s string) string {
				if strings.HasSuffix(s, ".log") {
					return "matched"
				}
				return ""
			}
			err := compressFiles(tt.C().srcDir, tt.C().dstDir, tt.C().level, matchFunc)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())

			_, err = os.Stat(tt.A().srcFile)
			testutil.Diff(t, tt.A().srcFile != "", err == nil)
			_, err = os.Stat(tt.A().dstFile)
			testutil.Diff(t, tt.A().dstFile != "", err == nil)
		})
	}
}

func TestGzipFile(t *testing.T) {
	type condition struct {
		open     func(string, int, fs.FileMode) (*os.File, error)
		remove   func(string) error
		create   []string
		srcDir   string
		dstDir   string
		filename string
		level    int
	}

	type action struct {
		srcFile string
		dstFile string
		err     error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndSrcExists := tb.Condition("src exists", "source file exists")
	actCheckError := tb.Action("error", "check that there is an error")
	table := tb.Build()

	testDir := testDir + "ut/kernel/io/"
	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"srcDir!=dstDir/src exists",
			[]string{cndSrcExists},
			[]string{},
			&condition{
				create:   []string{testDir + "gzipSrc/test.log"},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipDst",
				filename: "test.log",
				level:    gzip.BestSpeed,
			},
			&action{
				srcFile: "",
				dstFile: testDir + "gzipDst/test.log.gz",
				err:     nil,
			},
		),
		gen(
			"srcDir!=dstDir/src exists/dst exists",
			[]string{cndSrcExists},
			[]string{},
			&condition{
				create:   []string{testDir + "gzipSrc/test.log", testDir + "gzipDst/test.log.gz"},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipDst",
				filename: "test.log",
				level:    gzip.BestSpeed,
			},
			&action{
				srcFile: "",
				dstFile: testDir + "gzipDst/test.log.gz",
				err:     nil,
			},
		),
		gen(
			"srcDir!=dstDir/src exists/no compression",
			[]string{cndSrcExists},
			[]string{},
			&condition{
				create:   []string{testDir + "gzipSrc/test.log"},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipDst",
				filename: "test.log",
				level:    gzip.NoCompression,
			},
			&action{
				srcFile: "",
				dstFile: testDir + "gzipDst/test.log",
				err:     nil,
			},
		),
		gen(
			"srcDir!=dstDir/src exists/dst exists/no compression",
			[]string{cndSrcExists},
			[]string{},
			&condition{
				create:   []string{testDir + "gzipSrc/test.log", testDir + "gzipDst/test.log"},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipDst",
				filename: "test.log",
				level:    gzip.NoCompression,
			},
			&action{
				srcFile: "",
				dstFile: testDir + "gzipDst/test.log",
				err:     nil,
			},
		),
		gen(
			"srcDir==dstDir/src exists",
			[]string{cndSrcExists},
			[]string{},
			&condition{
				create:   []string{testDir + "gzipSrc/test.log"},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipSrc",
				filename: "test.log",
				level:    gzip.BestSpeed,
			},
			&action{
				srcFile: "",
				dstFile: testDir + "gzipSrc/test.log.gz",
				err:     nil,
			},
		),
		gen(
			"srcDir==dstDir/src exists/dst exists",
			[]string{cndSrcExists},
			[]string{},
			&condition{
				create:   []string{testDir + "gzipSrc/test.log", testDir + "gzipSrc/test.log.gz"},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipSrc",
				filename: "test.log",
				level:    gzip.BestSpeed,
			},
			&action{
				srcFile: "",
				dstFile: testDir + "gzipSrc/test.log.gz",
				err:     nil,
			},
		),
		gen(
			"srcDir==dstDir/src exists/no compression",
			[]string{cndSrcExists},
			[]string{},
			&condition{
				create:   []string{testDir + "gzipSrc/test.log"},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipSrc",
				filename: "test.log",
				level:    gzip.NoCompression,
			},
			&action{
				srcFile: testDir + "gzipSrc/test.log",
				dstFile: testDir + "gzipSrc/test.log",
				err:     nil,
			},
		),
		gen(
			"srcDir!=dstDir/src not exists",
			[]string{},
			[]string{},
			&condition{
				create:   []string{},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipDst",
				filename: "test.log",
				level:    gzip.BestSpeed,
			},
			&action{
				srcFile: "",
				dstFile: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeFile,
					Description: ErrDscFileSys,
				},
			},
		),
		gen(
			"srcDir!=dstDir/src not exists/no compression",
			[]string{},
			[]string{},
			&condition{
				create:   []string{},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipDst",
				filename: "test.log",
				level:    gzip.NoCompression,
			},
			&action{
				srcFile: "",
				dstFile: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeFile,
					Description: ErrDscFileSys,
				},
			},
		),
		gen(
			"srcDir==dstDir/src file not exists",
			[]string{},
			[]string{},
			&condition{
				create:   []string{},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipSrc",
				filename: "test.log",
				level:    gzip.BestSpeed,
			},
			&action{
				srcFile: "",
				dstFile: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeFile,
					Description: ErrDscFileSys,
				},
			},
		),
		gen(
			"srcDir==dstDir/src file not exists/no compression",
			[]string{},
			[]string{},
			&condition{
				create:   []string{},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipSrc",
				filename: "test.log",
				level:    gzip.NoCompression,
			},
			&action{
				srcFile: "",
				dstFile: "",
				err:     nil,
			},
		),
		gen(
			"level -3",
			[]string{},
			[]string{},
			&condition{
				create:   []string{testDir + "gzipSrc/test.log"},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipDst",
				filename: "test.log",
				level:    -3,
			},
			&action{
				srcFile: "",
				dstFile: testDir + "gzipDst/test.log.gz",
				err:     nil,
			},
		),
		gen(
			"level -2",
			[]string{},
			[]string{},
			&condition{
				create:   []string{testDir + "gzipSrc/test.log"},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipDst",
				filename: "test.log",
				level:    -2,
			},
			&action{
				srcFile: "",
				dstFile: testDir + "gzipDst/test.log.gz",
				err:     nil,
			},
		),
		gen(
			"level -1",
			[]string{},
			[]string{},
			&condition{
				create:   []string{testDir + "gzipSrc/test.log"},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipDst",
				filename: "test.log",
				level:    -1,
			},
			&action{
				srcFile: "",
				dstFile: testDir + "gzipDst/test.log.gz",
				err:     nil,
			},
		),
		gen(
			"level -10",
			[]string{},
			[]string{},
			&condition{
				create:   []string{testDir + "gzipSrc/test.log"},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipDst",
				filename: "test.log",
				level:    10,
			},
			&action{
				srcFile: "",
				dstFile: testDir + "gzipDst/test.log.gz",
				err:     nil,
			},
		),
		gen(
			"failed to open source file",
			[]string{cndSrcExists},
			[]string{actCheckError},
			&condition{
				open: func(s string, i int, fm fs.FileMode) (*os.File, error) {
					return nil, os.ErrPermission
				},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipDst",
				filename: "test.log",
				level:    gzip.BestSpeed,
			},
			&action{
				srcFile: "",
				dstFile: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeFile,
					Description: ErrDscFileSys,
				},
			},
		),
		gen(
			"failed to open destination file",
			[]string{cndSrcExists},
			[]string{actCheckError},
			&condition{
				open: func(s string, i int, fm fs.FileMode) (*os.File, error) {
					if filepath.Base(s) == "test.log.gz" {
						return nil, os.ErrPermission
					}
					return os.OpenFile(s, i, fm)
				},
				create:   []string{testDir + "gzipSrc/test.log"},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipDst",
				filename: "test.log",
				level:    gzip.BestSpeed,
			},
			&action{
				srcFile: testDir + "gzipSrc/test.log",
				dstFile: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeFile,
					Description: ErrDscFileSys,
				},
			},
		),
		gen(
			"copy error",
			[]string{cndSrcExists},
			[]string{actCheckError},
			&condition{
				open: func(s string, i int, fm fs.FileMode) (*os.File, error) {
					f, err := os.OpenFile(s, i, fm)
					if filepath.Base(s) == "test.log" {
						f.Close()
						return f, err
					}
					return f, err
				},
				create:   []string{testDir + "gzipSrc/test.log"},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipDst",
				filename: "test.log",
				level:    gzip.BestSpeed,
			},
			&action{
				srcFile: testDir + "gzipSrc/test.log",
				dstFile: "",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeFile,
					Description: ErrDscFileSys,
				},
			},
		),
		gen(
			"remove error",
			[]string{cndSrcExists},
			[]string{actCheckError},
			&condition{
				remove: func(s string) error {
					return os.ErrPermission
				},
				create:   []string{testDir + "gzipSrc/test.log"},
				srcDir:   testDir + "gzipSrc",
				dstDir:   testDir + "gzipDst",
				filename: "test.log",
				level:    gzip.BestSpeed,
			},
			&action{
				srcFile: testDir + "gzipSrc/test.log",
				dstFile: testDir + "gzipDst/test.log.gz",
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeFile,
					Description: ErrDscFileSys,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			if tt.C().open != nil {
				tmpOpen := OpenFile
				OpenFile = tt.C().open
				defer func() {
					OpenFile = tmpOpen
				}()
			}

			if tt.C().remove != nil {
				tmpRemove := Remove
				Remove = tt.C().remove
				defer func() {
					Remove = tmpRemove
				}()
			}

			for _, f := range tt.C().create {
				f, err := os.Create(f)
				if err != nil {
					t.Errorf("%#v\n", err)
					return
				}
				f.Close()
			}
			defer func() {
				for _, f := range tt.C().create {
					os.Remove(f)
				}
				os.Remove(tt.A().srcFile)
				os.Remove(tt.A().dstFile)
			}()

			t.Log(tt.C().srcDir, tt.C().dstDir, tt.C().filename, tt.C().level)
			err := gzipFile(tt.C().srcDir, tt.C().dstDir, tt.C().filename, tt.C().level)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			t.Logf("%#v\n", err)

			_, err = os.Stat(tt.A().srcFile)
			testutil.Diff(t, tt.A().srcFile != "", err == nil)
			_, err = os.Stat(tt.A().dstFile)
			testutil.Diff(t, tt.A().dstFile != "", err == nil)
		})
	}
}
