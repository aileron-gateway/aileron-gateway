package httpproxy

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/aileron-gateway/aileron-gateway/util/resilience"
)

// upstream is the interface of a target of circuit breaker.
type upstream interface {
	resilience.Entry
	// url returns url string of this target.
	url() *url.URL
	// notify is the notification interface to this target.
	// This method should be called when notifying the result
	notify(int, error)
}

// noopUpstream is no-operation load balancer upstream.
// This implements upstream interface.
type noopUpstream struct {
	weight    int
	rawURL    string
	parsedURL *url.URL
}

func (t *noopUpstream) ID() string {
	return t.rawURL
}

func (t *noopUpstream) Weight() int {
	return t.weight
}

// active returns the availability of this upstream.
// Noop load balancer upstream always return true.
func (t *noopUpstream) Active() bool {
	return true
}

func (t *noopUpstream) Hint() int {
	return 0 // Currently not used.
}

// url returns the proxy url for this upstream.
func (t *noopUpstream) url() *url.URL {
	return t.parsedURL
}

// notify count proxy result for this upstream tu update active status.
// Noop load balancer upstream does not accept any result
// because it is considered to e always active.
func (t *noopUpstream) notify(_ int, _ error) {
}

// lbUpstream is a load balancer upstream with circuit breaker.
// This implements proxy.upstream interface.
type lbUpstream struct {
	circuitBreaker

	// upstream is the url of this upstream server.
	weight    int
	rawURL    string
	parsedURL *url.URL

	// passiveEnabled enables passive health checking.
	// Enabling this reflect the result of actual request.
	passiveEnabled bool

	// initialDelay is the delay duration to starting the active health checking.
	initialDelay time.Duration
	// interval is the interval duration to send health checking request
	// when active health checking.
	interval time.Duration

	// closer is the close channel.
	closer chan struct{}
}

func (t *lbUpstream) ID() string {
	return t.rawURL
}

func (t *lbUpstream) Weight() int {
	return t.weight
}

func (t *lbUpstream) Hint() int {
	return 0 // Currently not used.
}

// url returns the proxy url for this upstream.
func (t *lbUpstream) url() *url.URL {
	return t.parsedURL
}

// notify notifies the upstream health status.
// The given result will be discarded if
// the passive health check is not enabled.
// Use this for passive health check result.
func (t *lbUpstream) notify(status int, err error) {
	if !t.passiveEnabled {
		return
	}
	t.feedback(status, err)
}

// feedback feedbacks the upstream health status to this object.
// Use this for active health check result.
func (t *lbUpstream) feedback(status int, err error) {
	if err != nil {
		t.countFailure()
		return
	}
	if status >= 500 {
		t.countFailure()
	} else {
		t.countSuccess()
	}
}

// close breaks health checking loop.
// This method is used for the convenience of testing.
//
//lint:ignore U1000 func (*lbUpstream).close is unused
func (t *lbUpstream) close() {
	close(t.closer)
}

// activeCheck actively check server health status by HTTP.
// The second argument of req should be reusable.
func (t *lbUpstream) activeCheckHTTP(rt http.RoundTripper, req *http.Request) {
	if t.closer == nil {
		t.closer = make(chan struct{}, 1)
	}

	// Wait for initial delay seconds.
	time.Sleep(t.initialDelay)

	// The interval must be grater than zero.
	// Otherwise, the icker will be panic.
	if t.interval <= 0 {
		t.interval = time.Second
	}

	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	for {
		// Send a health checking request.
		res, err := rt.RoundTrip(req)
		if err != nil {
			t.feedback(0, err)
		} else {
			t.feedback(res.StatusCode, nil)
			_, _ = io.Copy(io.Discard, res.Body)
			res.Body.Close()
		}

		select {
		case <-t.closer:
			return
		case <-ticker.C:
		}
	}
}

// activeCheck actively check server health status by TCP or UDP.
// Checkout the link below for available network and address format.
// https://pkg.go.dev/net#Dial
func (t *lbUpstream) activeCheck(nw, addr string) {
	if t.closer == nil {
		t.closer = make(chan struct{}, 1)
	}

	// Wait for initial delay seconds.
	time.Sleep(t.initialDelay)

	// The interval must be grater than zero.
	// Otherwise, the ticker will panic.
	if t.interval <= 0 {
		t.interval = time.Second
	}

	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	// Start active health checking.
	for {
		// Send a health checking request.
		conn, err := net.Dial(nw, addr)
		if err != nil {
			t.feedback(0, err)
		} else {
			conn.Close()
			t.feedback(http.StatusOK, nil)
		}

		select {
		case <-t.closer:
			return
		case <-ticker.C:
		}
	}
}
