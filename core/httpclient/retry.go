package httpclient

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/util/resilience"
)

const (
	ErrPkg           = "core/httpclient"
	ErrTypeRetry     = "retry"
	ErrDescRetryFail = "sending request failed after retry."
)

// retry applies retry to HTTP requests.
// Applying this tripperware will read entire body of the request and keep it in the memory.
// retry does not check the methods or headers to determine if the request is
// retry-able or not even the request contains any idempotency keys.
// See isReplayable() of http.Request below.
// https://cs.opensource.google/go/go/+/refs/tags/go1.23.1:src/net/http/request.go;l=1543
// This implements core.Tripperware interface.
type retry struct {
	// maxRetry is maximum retry count.
	// Initial requests is not included.
	maxRetry int

	// waiter determined wait time to the next request.
	waiter resilience.Waiter

	// maxContentLength is the maximum content length of the requests that can be retried.
	// Because retrying the request keep the entire body on memory, this value should not be
	// set too large.
	maxContentLength int64

	// retryStatus is the list of HTTP status codes to be retried.
	// When this is not set, only networking errors that cannot get response status codes are retried.
	retryStatus []int
}

func (t *retry) Tripperware(next http.RoundTripper) http.RoundTripper {
	return core.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		if t.maxRetry <= 0 {
			return next.RoundTrip(r)
		}

		// Skip to retry when Content-Length is unknown or exceeds maxContentLength.
		if r.ContentLength < 0 || r.ContentLength > t.maxContentLength {
			return next.RoundTrip(r)
		}

		var errs []string // Accumulate all errors occurred while repeating retries.
		var ticker *time.Ticker

		// make the request retry-able.
		// i.e. Make it possible to read the request body multiple times.
		newReq, err := setupRewindBody(r)
		if err != nil {
			return nil, (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeRetry,
				Description: ErrDescRetryFail,
			}).Wrap(err)
		}
		r = newReq

		// Initialize ticker with dummy duration.
		ticker = time.NewTicker(time.Second)
		defer ticker.Stop()

	loop:
		for i := 0; i <= t.maxRetry; i++ {
			res, err := next.RoundTrip(r)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return nil, err // Client may canceled the request.
				}
				info := strconv.Itoa(i+1) + "th request failed at unix millis " +
					strconv.FormatInt(time.Now().UnixMilli(), 10) + ". caused by " + err.Error()
				errs = append(errs, info)
			} else {
				if !slices.Contains(t.retryStatus, res.StatusCode) {
					return res, nil
				}
				// The returned status code indicates that the request should be retried.
				info := strconv.Itoa(i+1) + "th request failed at unix millis " +
					strconv.FormatInt(time.Now().UnixMilli(), 10) + " with status code " + strconv.Itoa(res.StatusCode)
				errs = append(errs, info)
			}

			ticker.Reset(t.waiter.Wait(i + 1))
			select {
			case <-r.Context().Done():
				break loop // Request should be timeout or canceled wile waiting.
			case <-ticker.C:
			}

			r, err = rewindBody(r) // Reset the request body to read it from the beginning.
			if err != nil {
				info := "failed to rewind body after " + strconv.Itoa(i+1) + "th request. caused by " + err.Error()
				errs = append(errs, info)
				break
			}
		}

		// Request was not completed successfully.
		// Return all the accumulated errors.
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeRetry,
			Description: ErrDescRetryFail,
			Detail:      strings.Join(errs, "; "),
		}).Wrap(err)
	})
}

// readTrackingBody tracks the read/write operation for the HTTP request body.
// Referencing https://go.dev/src/net/http/transport.go
type readTrackingBody struct {
	io.ReadCloser
	didRead  bool // didRead is the flag to know if the Read() method was called at least once.
	didClose bool // didClose is the flag to know if the Close() method was called at least once.
}

func (r *readTrackingBody) Read(data []byte) (int, error) {
	r.didRead = true
	return r.ReadCloser.Read(data)
}

func (r *readTrackingBody) Close() error {
	r.didClose = true
	return r.ReadCloser.Close()
}

// setupRewindBody returns a new request with a custom body wrapper
// that can report whether the body needs rewinding.
// The given request MUST NOT be nil, otherwise panics.
// Referencing https://go.dev/src/net/http/transport.go
func setupRewindBody(req *http.Request) (*http.Request, error) {
	if req.Body == nil || req.Body == http.NoBody {
		return req, nil
	}

	// Make the body reusable.
	if req.GetBody == nil {
		buf, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewReader(buf))
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(buf)), nil
		}
	}

	newReq := *req
	newReq.Body = &readTrackingBody{ReadCloser: req.Body}

	return &newReq, nil
}

// rewindBody rewinds the request body.
// The given request MUST NOT be nil, otherwise panics.
// Referencing https://go.dev/src/net/http/transport.go
func rewindBody(req *http.Request) (rewound *http.Request, err error) {
	if req.Body == nil || req.Body == http.NoBody {
		return req, nil // nothing to rewind.
	}

	if rtb, ok := req.Body.(*readTrackingBody); ok && !rtb.didRead && !rtb.didClose {
		return req, nil // nothing to rewind.
	}

	if req.GetBody == nil {
		// Cannot rewind.
		return nil, errors.New("cannot rewind because the GetBody is nil")
	}

	body, err := req.GetBody()
	if err != nil {
		return nil, err
	}

	newReq := *req
	newReq.Body = &readTrackingBody{ReadCloser: body}

	return &newReq, nil
}
