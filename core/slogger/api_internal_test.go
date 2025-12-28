// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package slogger

import (
	"bufio"
	"context"
	"io"
	"io/fs"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-projects/go/zlog"
	"github.com/aileron-projects/go/ztime/zcron"
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default manifest",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			tmpDiscard, tmpStdout, tmpStderr := io.Discard, os.Stdout, os.Stderr
			defer func() {
				io.Discard, os.Stdout, os.Stderr = tmpDiscard, tmpStdout, tmpStderr
			}()
			r, w, _ := os.Pipe()
			io.Discard, os.Stdout, os.Stderr = w, w, w

			server := api.NewContainerAPI()
			got, err := Resource.Create(server, tt.C.manifest)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)
			if tt.A.err != nil {
				testutil.Diff(t, nil, got)
				return
			}

			lg := got.(log.Logger)
			lg.Info(context.Background(), "test")
			w.Close()

			line, _ := io.ReadAll(r)
			testutil.Diff(t, true, strings.Contains(string(line), tt.A.logPart))
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

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil close",
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

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.C.fl.Finalize()
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
		})
	}
}

func TestNewFileWriter(t *testing.T) {
	type condition struct {
		spec *v1.LogOutputSpec
	}

	type action struct {
		w   *zlog.LogicalFile
		err error
	}

	testConfig := &zlog.LogicalFileConfig{
		Manager: &zlog.FileManagerConfig{
			SrcDir:  os.TempDir(),
			DstDir:  os.TempDir(),
			Pattern: "TestNewFileWriter.%i.log",
		},
		FileName: "TestNewFileWriter.log",
	}
	testFW, _ := zlog.NewLogicalFile(testConfig)

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"file",
			&condition{
				spec: &v1.LogOutputSpec{
					LogDir:      os.TempDir(),
					BackupDir:   os.TempDir(),
					LogFileName: "TestNewFileWriter.log",
				},
			},
			&action{
				w:   testFW,
				err: nil,
			},
		),
		gen(
			"use cron",
			&condition{
				spec: &v1.LogOutputSpec{
					LogDir:      os.TempDir(),
					BackupDir:   os.TempDir(),
					LogFileName: "TestNewFileWriter.log",
					Cron:        "* * * * *",
				},
			},
			&action{
				w:   testFW,
				err: nil,
			},
		),
		gen(
			"invalid cron",
			&condition{
				spec: &v1.LogOutputSpec{
					LogDir:    os.TempDir(),
					BackupDir: os.TempDir(),
					Cron:      "INVALID",
				},
			},
			&action{
				w:   nil,
				err: &zcron.ParseError{What: "number of fields"},
			},
		),
		gen(
			"log dir not exist",
			&condition{
				spec: &v1.LogOutputSpec{
					LogDir:    "this dir is not exist \x00\n",
					BackupDir: os.TempDir(),
				},
			},
			&action{
				w:   nil,
				err: &fs.PathError{Op: "mkdir", Path: "this dir is not exist \x00\n"},
			},
		),
		gen(
			"backup dir not exist",
			&condition{
				spec: &v1.LogOutputSpec{
					OutputTarget: v1.OutputTarget_File,
					LogDir:       os.TempDir(),
					BackupDir:    "this dir is not exist \x00\n",
				},
			},
			&action{
				w:   nil,
				err: &fs.PathError{Op: "mkdir", Path: "this dir is not exist \x00\n"},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			w, err := newFileWriter(tt.C.spec)
			t.Logf("%#v\n", err)
			testutil.DiffError(t, tt.A.err, nil, err, cmpopts.IgnoreFields(fs.PathError{}, "Err"))

			opts := []cmp.Option{
				cmp.AllowUnexported(sync.RWMutex{}, atomic.Int32{}, atomic.Int64{}, atomic.Bool{}),
				cmpopts.IgnoreUnexported(sync.Mutex{}, os.File{}),
				cmp.AllowUnexported(zlog.LogicalFile{}, zlog.FileManager{}),
				cmp.AllowUnexported(bufio.Writer{}),
			}
			testutil.Diff(t, tt.A.w, w, opts...)
		})
	}
}
