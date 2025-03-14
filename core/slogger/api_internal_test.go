package slogger

import (
	"bufio"
	"context"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/cron"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	kio "github.com/aileron-gateway/aileron-gateway/kernel/io"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
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
		logPart    string
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndDefaultManifest := tb.Condition("default manifest", "input default manifest")
	tb.Condition(cndDefaultManifest, "input default manifest")
	cndErrorReference := tb.Condition("error reference", "input an error reference to an object")
	actCheckError := tb.Action("check the returned error", "check that the returned error is the one expected")
	actCheckNoError := tb.Action("check no error", "check that there is no error returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
			[]string{cndDefaultManifest},
			[]string{actCheckNoError},
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				logPart: `"msg":"test"`,
				err:     nil,
			},
		),
		gen(
			"output discard",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				manifest: &v1.SLogger{
					Spec: &v1.SLoggerSpec{
						LogOutput: &v1.LogOutputSpec{
							OutputTarget: v1.OutputTarget_Discard,
						},
					},
				},
			},
			&action{
				logPart: `"msg":"test"`,
				err:     nil,
			},
		),
		gen(
			"output stderr",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				manifest: &v1.SLogger{
					Spec: &v1.SLoggerSpec{
						LogOutput: &v1.LogOutputSpec{
							OutputTarget: v1.OutputTarget_Stderr,
						},
					},
				},
			},
			&action{
				logPart: `"msg":"test"`,
				err:     nil,
			},
		),
		gen(
			"output stdout",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				manifest: &v1.SLogger{
					Spec: &v1.SLoggerSpec{
						LogOutput: &v1.LogOutputSpec{
							OutputTarget: v1.OutputTarget_Stdout,
						},
					},
				},
			},
			&action{
				logPart: `"msg":"test"`,
				err:     nil,
			},
		),
		gen(
			"output unknown",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				manifest: &v1.SLogger{
					Spec: &v1.SLoggerSpec{
						LogOutput: &v1.LogOutputSpec{
							OutputTarget: v1.OutputTarget(9999),
						},
					},
				},
			},
			&action{
				logPart: `"msg":"test"`,
				err:     nil,
			},
		),
		gen(
			"create with replacer",
			[]string{cndDefaultManifest},
			[]string{actCheckNoError},
			&condition{
				manifest: &v1.SLogger{
					Spec: &v1.SLoggerSpec{
						LogOutput:    &v1.LogOutputSpec{},
						Unstructured: true,
						FieldReplacers: []*v1.FieldReplacerSpec{
							{Field: "foo.bar.baz"},
						},
					},
				},
			},
			&action{
				logPart: `msg=test`,
				err:     nil,
			},
		),
		gen(
			"create replacer failed",
			[]string{cndDefaultManifest},
			[]string{actCheckNoError},
			&condition{
				manifest: &v1.SLogger{
					Spec: &v1.SLoggerSpec{
						LogOutput:    &v1.LogOutputSpec{},
						Unstructured: true,
						FieldReplacers: []*v1.FieldReplacerSpec{
							{
								Field: "foo",
								Replacer: &kernel.ReplacerSpec{
									Replacers: &kernel.ReplacerSpec_Regexp{
										Regexp: &kernel.RegexpReplacer{Pattern: "[0-9"},
									},
								},
							},
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create SLogger`),
			},
		),
		gen(
			"create unstructured logger",
			[]string{cndDefaultManifest},
			[]string{actCheckNoError},
			&condition{
				manifest: &v1.SLogger{
					Spec: &v1.SLoggerSpec{
						LogOutput:    &v1.LogOutputSpec{},
						Unstructured: true,
					},
				},
			},
			&action{
				logPart: `msg=test`,
				err:     nil,
			},
		),
		gen(
			"create file logger",
			[]string{cndDefaultManifest},
			[]string{actCheckNoError},
			&condition{
				manifest: &v1.SLogger{
					Spec: &v1.SLoggerSpec{
						LogOutput: &v1.LogOutputSpec{
							OutputTarget: v1.OutputTarget_File,
							LogDir:       os.TempDir(),
						},
					},
				},
			},
			&action{
				logPart: ``,
				err:     nil,
			},
		),
		gen(
			"invalid timezone format",
			[]string{cndErrorReference},
			[]string{actCheckError, actCheckNoError},
			&condition{
				manifest: &v1.SLogger{
					Spec: &v1.SLoggerSpec{
						LogOutput: &v1.LogOutputSpec{
							OutputTarget: v1.OutputTarget_File,
							TimeZone:     "INVALID",
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create SLogger`),
			},
		),
		gen(
			"invalid log output spec",
			[]string{cndErrorReference},
			[]string{actCheckError, actCheckNoError},
			&condition{
				manifest: &v1.SLogger{
					Spec: &v1.SLoggerSpec{
						LogOutput: &v1.LogOutputSpec{
							OutputTarget: v1.OutputTarget_File,
							LogDir:       "invalid dir \x00\n",
						},
					},
				},
			},
			&action{
				err:        core.ErrCoreGenCreateObject,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create SLogger`),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmpDiscard, tmpStdout, tmpStderr := io.Discard, os.Stdout, os.Stderr
			defer func() {
				io.Discard, os.Stdout, os.Stderr = tmpDiscard, tmpStdout, tmpStderr
			}()
			r, w, _ := os.Pipe()
			io.Discard, os.Stdout, os.Stderr = w, w, w

			server := api.NewContainerAPI()
			got, err := Resource.Create(server, tt.C().manifest)
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
			if tt.A().err != nil {
				testutil.Diff(t, nil, got)
				return
			}

			lg := got.(log.Logger)
			lg.Info(context.Background(), "test")
			w.Close()

			line, _ := io.ReadAll(r)
			testutil.Diff(t, true, strings.Contains(string(line), tt.A().logPart))
		})
	}
}

type testCloser struct {
	called bool
	err    error
}

func (c *testCloser) Close() error {
	c.called = true
	return c.err
}

func TestFinalizableLogger_Finalize(t *testing.T) {
	type condition struct {
		fl *finalizableLogger
	}

	type action struct {
		err any // error or errorutil.Kind
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndCloseError := tb.Condition("close error", "input an error closer")
	actCheckError := tb.Action("check error", "check that the returned error is the one expected")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil close",
			[]string{},
			[]string{},
			&condition{
				fl: &finalizableLogger{
					closer: nil,
				},
			},
			&action{
				err: nil,
			},
		),
		gen(
			"no error closer",
			[]string{},
			[]string{},
			&condition{
				fl: &finalizableLogger{
					closer: &testCloser{},
				},
			},
			&action{
				err: nil,
			},
		),
		gen(
			"error closer",
			[]string{cndCloseError},
			[]string{actCheckError},
			&condition{
				fl: &finalizableLogger{
					closer: &testCloser{
						err: io.EOF, // Dummy error.
					},
				},
			},
			&action{
				err: io.EOF, // Dummy error
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			err := tt.C().fl.Finalize()
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestNewFileWriter(t *testing.T) {
	type condition struct {
		spec *v1.LogOutputSpec
	}

	type action struct {
		w   *kio.LogicalFile
		err error
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	testConfig := &kio.LogicalFileConfig{
		SrcDir: os.TempDir(),
		DstDir: os.TempDir(),
	}
	testFW, _ := testConfig.New()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"file",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.LogOutputSpec{
					LogDir:    os.TempDir(),
					BackupDir: os.TempDir(),
				},
			},
			&action{
				w:   testFW,
				err: nil,
			},
		),
		gen(
			"use cron",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.LogOutputSpec{
					LogDir:    os.TempDir(),
					BackupDir: os.TempDir(),
					Cron:      "* * * * *",
				},
			},
			&action{
				w:   testFW,
				err: nil,
			},
		),
		gen(
			"invalid cron",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.LogOutputSpec{
					LogDir:    os.TempDir(),
					BackupDir: os.TempDir(),
					Cron:      "INVALID",
				},
			},
			&action{
				w: nil,
				err: &er.Error{
					Package:     cron.ErrPkg,
					Type:        cron.ErrTypeParse,
					Description: cron.ErrDscParse,
				},
			},
		),
		gen(
			"log dir not exist",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.LogOutputSpec{
					LogDir:    "this dir is not exist \x00\n",
					BackupDir: os.TempDir(),
				},
			},
			&action{
				w: nil,
				err: &er.Error{
					Package:     kio.ErrPkg,
					Type:        kio.ErrTypeFile,
					Description: kio.ErrDscFileSys,
				},
			},
		),
		gen(
			"backup dir not exist",
			[]string{},
			[]string{},
			&condition{
				spec: &v1.LogOutputSpec{
					OutputTarget: v1.OutputTarget_File,
					LogDir:       os.TempDir(),
					BackupDir:    "this dir is not exist \x00\n",
				},
			},
			&action{
				w: nil,
				err: &er.Error{
					Package:     kio.ErrPkg,
					Type:        kio.ErrTypeFile,
					Description: kio.ErrDscFileSys,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			w, err := newFileWriter(tt.C().spec, time.UTC)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())

			opts := []cmp.Option{
				cmp.AllowUnexported(sync.Mutex{}, sync.RWMutex{}, atomic.Int32{}, atomic.Int64{}, atomic.Bool{}),
				cmp.Comparer(testutil.ComparePointer[*os.File]),
				// Options for bufferef writer.
				cmp.AllowUnexported(bufio.Writer{}),
				// Options for file writer.
				cmp.AllowUnexported(kio.LogicalFile{}),
				cmpopts.IgnoreFields(kio.LogicalFile{}, "curFile", "mu"),
				cmp.Comparer(testutil.ComparePointer[func() error]),              // manageFunc in LogicalFile.
				cmp.Comparer(testutil.ComparePointer[func(string) string]),       // matchFunc in LogicalFile.
				cmp.Comparer(testutil.ComparePointer[func(string) (int64, int)]), // parseFunc in LogicalFile.
			}

			testutil.Diff(t, tt.A().w, w, opts...)
		})
	}
}
