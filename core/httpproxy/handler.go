package httpproxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aileron-gateway/aileron-gateway/core"
	kio "github.com/aileron-gateway/aileron-gateway/kernel/io"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/quic-go/quic-go/http3"
)

// reverseProxy is the reverse proxy handler.
// This implements http.Handler interface.
type reverseProxy struct {
	// HandlerBase is the base struct for
	// http.Handler type resource.
	// This provides Patterns() and Methods() methods
	// to fulfill the core.Handler interface.
	*utilhttp.HandlerBase

	lg log.Logger
	eh core.ErrorHandler

	// lbs is the slice of load balancers.
	lbs []loadBalancer

	// rt is the round tripper to be used for proxy requests.
	rt http.RoundTripper
}

// findUpstream returns a proxy upstream.
// false is turned when no proxy upstream is available.
func (p *reverseProxy) findUpstream(r *http.Request) (upstream, *url.URL, core.HTTPError) {
	for _, lb := range p.lbs {
		if t, url, matched := lb.upstream(r); matched {
			if t == nil {
				// Matched but no upstream available.
				// That means all upstream servers are downed.
				err := core.ErrCoreProxyUnavailable.WithStack(nil, map[string]any{"path": r.URL.Path})
				return nil, nil, utilhttp.NewHTTPError(err, http.StatusBadGateway)
			}
			return t, url, nil
		}
	}
	// Upstream not found.
	err := core.ErrCoreProxyNoUpstream.WithoutStack(nil, map[string]any{"path": r.URL.Path})
	return nil, nil, utilhttp.NewHTTPError(err, http.StatusNotFound)
}

func (p *reverseProxy) logIfError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	attr := core.ErrCoreProxyNoRecovery.WithStack(err, nil)
	p.lg.Info(ctx, "proxy error", attr.Name(), attr.Map())
}

func (p *reverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hooks := utilhttp.PreProxyHookFromContext(r.Context())
	for _, hook := range hooks {
		if err := hook(r); err != nil {
			p.eh.ServeHTTPError(w, r, err)
			return
		}
	}

	upstream, upstreamURL, findErr := p.findUpstream(r)
	if findErr != nil {
		p.eh.ServeHTTPError(w, r, findErr)
		return
	}

	outReq := r.Clone(r.Context())
	outReq.Host = ""
	rewriteRequestURL(outReq.URL, upstreamURL) // Rewrite request url to upstream url.
	if outReq.Header == nil {
		outReq.Header = make(http.Header, 0)
	}
	if outReq.ContentLength == 0 {
		r.Body = nil // Issue 16036: nil Body for http.Transport retries
	}
	if outReq.Body != nil {
		defer outReq.Body.Close()
	}

	removeHopByHopHeaders(outReq.Header)
	handleHopByHopHeaders(r.Header, outReq.Header)
	setXForwardedHeaders(r, outReq.Header)
	copyHeader(outReq.Header, utilhttp.ProxyHeaderFromContext(r.Context()))

	outRes, err := p.rt.RoundTrip(outReq)
	if err != nil {
		// Notify the error to upstream object so the
		// upstream object can know the health status of the upstream server.
		per := proxyErrorResponse(err)
		upstream.notify(per.StatusCode(), nil)
		p.eh.ServeHTTPError(w, r, per)
		return
	}

	// Notify the status code to upstream object
	// so the upstream object identify the upstream server
	// is healthy or not.
	upstream.notify(outRes.StatusCode, nil)

	// Handle protocol switching.
	// This blocks bi-directional communication until it finished.
	if outRes.StatusCode == http.StatusSwitchingProtocols {
		if err := handleUpgradeResponse(w, outReq, outRes); err != nil {
			p.eh.ServeHTTPError(w, r, err)
		}
		return
	}

	removeHopByHopHeaders(outRes.Header)
	copyHeader(w.Header(), outRes.Header)

	if len(outRes.Trailer) > 0 {
		announcedTrailerKeys := make([]string, 0, len(outRes.Trailer))
		for k := range outRes.Trailer {
			announcedTrailerKeys = append(announcedTrailerKeys, k)
		}
		w.Header().Add("Trailer", strings.Join(announcedTrailerKeys, ", "))
	}

	// No responses are allowed to write something into the response writer after WriteHeader.
	// According to the *http.Response, response body is always non-nil.
	w.WriteHeader(outRes.StatusCode)

	// Copy response body to client.
	// CopyBuffer blocks until copy finished.
	if _, err = kio.CopyBuffer(withImmediateFlush(w, shouldFlushImmediately(outRes)), outRes.Body); err != nil {
		// We can't write anything to the response writer any more.
		// So, we only output log of the returned error.
		p.logIfError(r.Context(), outRes.Body.Close())
		// Client canceled request. Following error may happen.
		//  	- [context.Canceled]
		//  	- [net.OpError] <-- Both read and write.
		p.eh.ServeHTTPError(w, r, utilhttp.NewHTTPError(err, -1))
		return
	}

	// close now, instead of defer, to populate outRes.Trailer
	p.logIfError(r.Context(), outRes.Body.Close())

	if len(outRes.Trailer) > 0 {
		// Force chunking if we saw a response trailer.
		// This prevents net/http from calculating the length for short bodies and adding a Content-Length.
		p.logIfError(r.Context(), http.NewResponseController(w).Flush())
		copyTrailer(w.Header(), outRes.Trailer) // Copy trailer with http.TrailerPrefix.
	}
}

// proxyErrorResponse returns an appropriate
// http error from the given error .
// This function can panic when a nil error was given.
func proxyErrorResponse(err error) core.HTTPError {
	switch {
	case err == context.Canceled:
		return utilhttp.NewHTTPError(err, -1) // LoggingOnly
	case http3RequestCanceled(err):
		return utilhttp.NewHTTPError(err, -1) // LoggingOnly
	case errors.Is(err, context.DeadlineExceeded):
		// Timeout occurred while attempting to get response from proxy upstream.
		err = core.ErrCoreProxyTimeout.WithStack(err, nil)
		return utilhttp.NewHTTPError(err, http.StatusGatewayTimeout)
	default:
		err = core.ErrCoreProxyRoundtrip.WithStack(err, nil)
		return utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}
}

// http3RequestCanceled returns if the
// context was canceled by the client when using http3 server.
// This error can happen when the server is http3.
func http3RequestCanceled(err error) bool {
	v, ok := err.(*http3.Error)
	if !ok {
		return false
	}
	return v.ErrorCode == http3.ErrCodeRequestCanceled
}

// handleUpgradeResponse handles protocol upgrade.
// This method is called when http.StatusSwitchingProtocols was detected.
// Refer to the handleUpgradeResponse method in https://go.dev/src/net/http/httputil/reverseproxy.go
func handleUpgradeResponse(rw http.ResponseWriter, req *http.Request, res *http.Response) core.HTTPError {
	reqUpType := upgradeType(req.Header)
	resUpType := upgradeType(res.Header)
	if len(reqUpType) != len(resUpType) || !strings.EqualFold(reqUpType, resUpType) {
		reason := fmt.Sprintf("backend tried to switch protocol %q when %q was requested", resUpType, reqUpType)
		err := core.ErrCoreProxyProtocolSwitch.WithStack(nil, map[string]any{"reason": reason})
		return utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}

	backConn, ok := res.Body.(io.ReadWriteCloser)
	if !ok {
		reason := "internal error: 101 switching protocols response with non-writable body"
		err := core.ErrCoreProxyProtocolSwitch.WithStack(nil, map[string]any{"reason": reason})
		return utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}

	backConnCloseCh := make(chan bool)
	defer close(backConnCloseCh)
	go func() {
		// Ensure that the cancellation of a request closes the backend.
		// See issue https://golang.org/issue/35559.
		select {
		case <-req.Context().Done():
		case <-backConnCloseCh:
		}
		backConn.Close()
	}()

	conn, brw, err := http.NewResponseController(rw).Hijack()
	if err != nil {
		reason := "Hijack failed from type " + fmt.Sprintf("%T", rw)
		err = core.ErrCoreProxyProtocolSwitch.WithStack(err, map[string]any{"reason": reason})
		return utilhttp.NewHTTPError(err, http.StatusInternalServerError)
	}
	defer conn.Close()

	copyHeader(rw.Header(), res.Header)

	res.Header = rw.Header()
	res.Body = nil // so res.Write only writes the headers; we have res.Body in backConn above
	if err := res.Write(brw); err != nil {
		err = core.ErrCoreProxyProtocolSwitch.WithStack(err, map[string]any{"reason": "response write failed"})
		return utilhttp.NewHTTPError(err, -1) // LoggingOnly
	}
	if err := brw.Flush(); err != nil {
		err = core.ErrCoreProxyProtocolSwitch.WithStack(err, map[string]any{"reason": "response flush failed"})
		return utilhttp.NewHTTPError(err, -1) // LoggingOnly
	}

	errChan := make(chan error, 1)
	bid := &kio.BidirectionalReadWriter{
		Frontend: conn,
		Backend:  backConn,
	}

	// Stat the bi-directional communication and wait until it finished.
	go bid.CopyToBackend(errChan)
	go bid.CopyFromBackend(errChan)
	if err = <-errChan; err != nil {
		err = core.ErrCoreProxyBidirectionalCom.WithStack(err, nil)
		return utilhttp.NewHTTPError(err, -1) // LoggingOnly
	}

	return nil
}
