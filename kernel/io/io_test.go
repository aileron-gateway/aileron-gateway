package io

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testReadWriter struct {
	*testReader
	*testWriter
}

type testWriter struct {
	err      error
	written  []byte
	readOnly int
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	if w.readOnly > 0 {
		w.written = append(w.written, p[:w.readOnly]...)
		return w.readOnly, w.err
	}
	w.written = append(w.written, p...)
	return len(p), w.err
}

type testReader struct {
	err  error
	data []byte
}

func (r *testReader) Read(p []byte) (n int, err error) {
	if len(p) >= len(r.data) {
		copy(p, r.data)
		return len(r.data), r.err
	}
	copy(p, r.data[:len(p)])
	return len(p), r.err
}

func TestCopyBuffer(t *testing.T) {
	type condition struct {
		dst *testWriter
		src *testReader
	}

	type action struct {
		written []byte
		err     any // error or errorutil.Kind
	}

	cndNoError := "no reader/writer error"
	cndReaderError := "reader error"
	cndWriterError := "writer error"
	actCheckWritten := "check written"
	actCheckError := "check error"
	actCheckNoError := "check no error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNoError, "no error occurs on read and write")
	tb.Condition(cndReaderError, "reader method returns an error on read")
	tb.Condition(cndWriterError, "writer method returns an error on write")
	tb.Action(actCheckWritten, "check the value written to a writer")
	tb.Action(actCheckError, "check that the expected non-nil error returned")
	tb.Action(actCheckNoError, "check that there is no error")
	table := tb.Build()

	testErrRead := errors.New("test read error")
	testErrWrite := errors.New("test write error")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no error",
			[]string{cndNoError},
			[]string{actCheckWritten, actCheckNoError},
			&condition{
				dst: &testWriter{
					err: nil,
				},
				src: &testReader{
					err:  io.EOF,
					data: []byte("test"),
				},
			},
			&action{
				written: []byte("test"),
				err:     nil,
			},
		),
		gen(
			"read error",
			[]string{cndReaderError},
			[]string{actCheckWritten, actCheckError},
			&condition{
				dst: &testWriter{
					err: nil,
				},
				src: &testReader{
					err:  testErrRead,
					data: []byte("test"),
				},
			},
			&action{
				written: []byte("test"),
				err:     testErrRead,
			},
		),
		gen(
			"write error",
			[]string{cndWriterError},
			[]string{actCheckWritten, actCheckError},
			&condition{
				dst: &testWriter{
					err: testErrWrite,
				},
				src: &testReader{
					err:  nil,
					data: []byte("test"),
				},
			},
			&action{
				written: []byte("test"),
				err:     testErrWrite,
			},
		),
		gen(
			"short write error",
			[]string{cndWriterError},
			[]string{actCheckWritten, actCheckError},
			&condition{
				dst: &testWriter{
					err:      nil,
					readOnly: 2, // read only 2 bytes.
				},
				src: &testReader{
					err:  nil,
					data: []byte("test"),
				},
			},
			&action{
				written: []byte("te"), // Only 2 bytes are read.
				err:     io.ErrShortWrite,
			},
		),
		gen(
			"context cancel error",
			[]string{cndReaderError},
			[]string{actCheckWritten, actCheckError},
			&condition{
				dst: &testWriter{
					err: nil,
				},
				src: &testReader{
					err:  context.Canceled,
					data: []byte("test"),
				},
			},
			&action{
				written: []byte("test"),
				err:     context.Canceled,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			n, err := CopyBuffer(tt.C().dst, tt.C().src)

			testutil.Diff(t, tt.A().written, tt.C().dst.written)
			testutil.Diff(t, len(tt.A().written), int(n))
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestBidirectionalReadWriter_copyFromBackend(t *testing.T) {
	type condition struct {
		frontend *testReadWriter
		backend  *testReadWriter
	}

	type action struct {
		written []byte
		err     any // error or errorutil.Kind
	}

	cndNoError := "no error"
	cndReadError := "read error"
	cndWriterError := "write error"
	actCheckWritten := "check written bytes"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNoError, "no error on read and write")
	tb.Condition(cndReadError, "an error was returned on read")
	tb.Condition(cndWriterError, "an error was returned on write")
	tb.Action(actCheckWritten, "check the written bytes")
	tb.Action(actCheckError, "check the returned error")
	table := tb.Build()

	testErrRead := errors.New("test read error")
	testErrWrite := errors.New("test write error")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no error",
			[]string{cndNoError},
			[]string{actCheckWritten},
			&condition{
				frontend: &testReadWriter{
					testWriter: &testWriter{},
				},
				backend: &testReadWriter{
					testReader: &testReader{
						err:  io.EOF,
						data: []byte("test"),
					},
				},
			},
			&action{
				written: []byte("test"),
				err:     nil,
			},
		),
		gen(
			"read error",
			[]string{cndReadError},
			[]string{actCheckWritten, actCheckError},
			&condition{
				frontend: &testReadWriter{
					testWriter: &testWriter{},
				},
				backend: &testReadWriter{
					testReader: &testReader{
						err:  testErrRead,
						data: []byte("test"),
					},
				},
			},
			&action{
				written: []byte("test"),
				err:     testErrRead,
			},
		),
		gen(
			"write error",
			[]string{cndWriterError},
			[]string{actCheckWritten, actCheckError},
			&condition{
				frontend: &testReadWriter{
					testWriter: &testWriter{
						err: testErrWrite,
					},
				},
				backend: &testReadWriter{
					testReader: &testReader{
						err:  io.EOF,
						data: []byte("test"),
					},
				},
			},
			&action{
				written: []byte("test"),
				err:     testErrWrite,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			errChan := make(chan error, 1)

			rw := &BidirectionalReadWriter{
				Frontend: tt.C().frontend,
				Backend:  tt.C().backend,
			}

			rw.CopyFromBackend(errChan)
			err := <-errChan

			testutil.Diff(t, tt.A().written, tt.C().frontend.written)
			testutil.Diff(t, int64(len(tt.A().written)), rw.WrittenToFront)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestBidirectionalReadWriter_copyToBackend(t *testing.T) {
	type condition struct {
		frontend *testReadWriter
		backend  *testReadWriter
	}

	type action struct {
		written []byte
		err     any // error or errorutil.Kind
	}

	cndNoError := "no error"
	cndReadError := "read error"
	cndWriterError := "write error"
	actCheckWritten := "check written bytes"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndNoError, "no error on read and write")
	tb.Condition(cndReadError, "an error was returned on read")
	tb.Condition(cndWriterError, "an error was returned on write")
	tb.Action(actCheckWritten, "check the written bytes")
	tb.Action(actCheckError, "check the returned error")
	table := tb.Build()

	testErrRead := errors.New("test read error")
	testErrWrite := errors.New("test write error")

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no error",
			[]string{cndNoError},
			[]string{actCheckWritten},
			&condition{
				frontend: &testReadWriter{
					testReader: &testReader{
						err:  io.EOF,
						data: []byte("test"),
					},
				},
				backend: &testReadWriter{
					testWriter: &testWriter{},
				},
			},
			&action{
				written: []byte("test"),
				err:     nil,
			},
		),
		gen(
			"read error",
			[]string{cndReadError},
			[]string{actCheckWritten, actCheckError},
			&condition{
				frontend: &testReadWriter{
					testReader: &testReader{
						err:  testErrRead,
						data: []byte("test"),
					},
				},
				backend: &testReadWriter{
					testWriter: &testWriter{},
				},
			},
			&action{
				written: []byte("test"),
				err:     testErrRead,
			},
		),
		gen(
			"write error",
			[]string{cndWriterError},
			[]string{actCheckWritten, actCheckError},
			&condition{
				frontend: &testReadWriter{
					testReader: &testReader{
						err:  io.EOF,
						data: []byte("test"),
					},
				},
				backend: &testReadWriter{
					testWriter: &testWriter{
						err: testErrWrite,
					},
				},
			},
			&action{
				written: []byte("test"),
				err:     testErrWrite,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			errChan := make(chan error, 1)

			rw := &BidirectionalReadWriter{
				Frontend: tt.C().frontend,
				Backend:  tt.C().backend,
			}

			rw.CopyToBackend(errChan)
			err := <-errChan

			testutil.Diff(t, tt.A().written, tt.C().backend.written)
			testutil.Diff(t, int64(len(tt.A().written)), rw.WrittenToBack)
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}
