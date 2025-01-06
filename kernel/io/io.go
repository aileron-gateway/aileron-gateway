package io

import (
	"io"
	"sync"
)

// pool is the buffer pool of []byte.
//
//	buf := *pool.Get().(*[]byte)
//	defer pool.Put(&buf)
var pool = sync.Pool{
	New: func() any {
		b := make([]byte, 4*1024)
		return &b
	},
}

// CopyBuffer returns any write errors or non-EOF read errors,
// and the amount of bytes written.
// CopyBuffer does not report io.EOF.
// dst and src must not be nil.
// The implementation is almost the same as io.CopyBuffer
// except for the minor performance improvement.
//   - https://pkg.go.dev/io#CopyBuffer
func CopyBuffer(dst io.Writer, src io.Reader) (int64, error) {
	buf := *pool.Get().(*[]byte)
	defer pool.Put(&buf)

	var written int64
	for {
		nRead, readErr := src.Read(buf)

		// First, write all read bytes.
		if nRead > 0 {
			nWrite, writeErr := dst.Write(buf[:nRead])
			written += int64(nWrite)
			if writeErr != nil {
				return written, writeErr
			}
			if nRead != nWrite {
				return written, io.ErrShortWrite
			}
		}

		if readErr != nil {
			if readErr == io.EOF {
				return written, nil
			}
			return written, readErr
		}
	}
}

// BidirectionalReadWriter copies data
// from frontend to backend and vice versa.
type BidirectionalReadWriter struct {
	Frontend       io.ReadWriter
	Backend        io.ReadWriter
	WrittenToFront int64
	WrittenToBack  int64
}

// CopyFromBackend copies data from backend to frontend.
// This method returns after copy was completed.
// An err will be sent to errChan if any.
func (b *BidirectionalReadWriter) CopyFromBackend(errChan chan<- error) {
	n, err := CopyBuffer(b.Frontend, b.Backend)
	b.WrittenToFront += n
	errChan <- err
}

// CopyToBackend copies data from frontend to backend.
// This method returns after copy was completed.
// An err will be sent to errChan if any.
func (b *BidirectionalReadWriter) CopyToBackend(errChan chan<- error) {
	n, err := CopyBuffer(b.Backend, b.Frontend)
	b.WrittenToBack += n
	errChan <- err
}
