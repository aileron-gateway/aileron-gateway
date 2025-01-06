package network

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/http"
	"slices"
	"sync"
	"time"

	"golang.org/x/net/http/httpguts"
	"golang.org/x/net/http2"
)

// HTTP2AddrResolver is used for resolving
// ip addresses in http2ClientConn pool.
var HTTP2AddrResolver = net.DefaultResolver

// dnsResolver resolves host to IP addresses.
type dnsResolver struct {
	// mu protects addrs and lastResolved.
	mu sync.RWMutex
	// host is the host name.
	// This should be an CNAME, not ip addresses.
	host string
	// ips is the list of ip addressed
	// resolved with the host by DNS resolver.
	ips []string
	// current is the current position of the
	// addrs to use.
	// IP addresses are returned by round robin way.
	current int
	// stop stops the resolve.
	stop chan struct{}
}

// length returns the current number of ip addresses.
// Returned length can be 0.
// Call resolve method to immediately update addresses.
// length is safe for concurrent call.
func (r *dnsResolver) length() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.ips)
}

// next returns the next ip address that should be used.
// Addresses are returned by round robin like order.
// next is safe for concurrent call.
// An empty string "" is returned when there is no addresses to return.
func (r *dnsResolver) next() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.ips) == 0 {
		return ""
	}
	r.current += 1
	if r.current > len(r.ips)-1 {
		r.current = 0
	}
	return r.ips[r.current]
}

func (r *dnsResolver) stopResolveLoop() {
	if r.stop != nil {
		close(r.stop)
	}
}

// resolveEveryInterval resolve host every given interval.
// If the interval is zero or negative, this method
// do nothing and returns immediately.
// This method is intended to be run in a isolated goroutine.
func (r *dnsResolver) resolveEveryInterval(ctx context.Context, interval time.Duration) {
	if r.stop != nil || interval <= 0 {
		return // Resolve already working.
	}
	r.stop = make(chan struct{})

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-r.stop:
			return
		case <-ticker.C:
		}
		err := r.resolve(ctx)
		if err != nil && VerboseLogs {
			log.Printf("network: DNS resolved error %s\n", err.Error())
		}
	}
}

// resolve resolves host to ip addresses.
// resolve is safe for concurrent call.
func (r *dnsResolver) resolve(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second) // Lookup timeout,currently 2 seconds.
	defer cancel()

	ips, err := HTTP2AddrResolver.LookupIPAddr(ctx, r.host)
	if err != nil {
		return err // Do not replace the existing ips.
	}

	ipAddrs := make([]string, 0, len(ips))
	for _, ip := range ips {
		ipAddrs = append(ipAddrs, ip.IP.String())
	}
	slices.Sort(ipAddrs) // Sort addrs to keep the order.

	if VerboseLogs {
		log.Printf("network: DNS resolved addrs %s >> %+v\n", r.host, ipAddrs)
	}

	r.mu.Lock()
	r.ips = ipAddrs // Update ip list.
	r.mu.Unlock()
	return nil
}

// hostConns is the connection pool
// bounded to a single host.
type hostConns struct {
	*dnsResolver

	// addr is the target address.
	addr string
	// port number of the addr.
	port string

	// mu protects conns and connToIP.
	mu sync.RWMutex
	// conns is the IP address to connection mapping.
	// key is the IP address string obtained from the dnsResolver.
	conns map[string][]*http2.ClientConn
	// connToIP is the connection to IP address mapping.
	// This is for deleting dead connections.
	// Connections are not removed until the
	// markDead method was called.
	// So markDead must be called for dead connections
	// otherwise, results in memory leaks.
	// This map is for markDead.
	connToIP map[*http2.ClientConn]string

	tc *tls.Config
	t  *http2.Transport
}

// markDead marks the given connection as dead.
// markDead returns the remaining number of active connections.
func (c *hostConns) markDead(conn *http2.ClientConn) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	ip, ok := c.connToIP[conn]
	if !ok {
		return len(c.connToIP)
	}
	delete(c.connToIP, conn)

	conns, ok := c.conns[ip]
	if !ok {
		return len(c.connToIP)
	}
	for i := range conns {
		if conn == conns[i] {
			conns = slices.Delete(conns, i, i+1)
			c.conns[ip] = conns
			break
		}
	}
	if len(conns) == 0 {
		delete(c.conns, ip) // Unregister the ip.
	}

	return len(c.connToIP)
}

// newClientConn returns a new client connection.
func (c *hostConns) newClientConn(ctx context.Context, addr string) (*http2.ClientConn, error) {
	conn, err := c.t.DialTLSContext(ctx, "tcp", addr, c.tc) // Dialer should have timeout itself.
	if err != nil {
		return nil, err
	}
	return c.t.NewClientConn(conn)
}

func (c *hostConns) getClientConn(ctx context.Context) (*http2.ClientConn, error) {
	// Available ip addresses may be zero.
	num := c.dnsResolver.length()

	var errs []error
	for i := 0; i < num; i++ {
		ip := c.dnsResolver.next()

		c.mu.RLock()
		conns := c.conns[ip]
		c.mu.RUnlock()

		for _, conn := range conns {
			if conn.ReserveNewRequest() {
				c.mu.Lock()
				c.connToIP[conn] = ip
				c.mu.Unlock()
				return conn, nil
			}
		}

		conn, err := c.newClientConn(ctx, net.JoinHostPort(ip, c.port))
		if err != nil {
			errs = append(errs, err) // Keep error and try next ip.
			continue
		}
		if conn.ReserveNewRequest() {
			c.mu.Lock()
			c.conns[ip] = append(c.conns[ip], conn)
			c.connToIP[conn] = ip
			c.mu.Unlock()
			return conn, nil
		}
	}

	// Finally, could not obtain any connection.
	errs = append([]error{http2.ErrNoCachedConn}, errs...)
	return nil, errors.Join(errs...)
}

// http2ConnPool is the HTTP2 client connection pool.
// This is used as the connection pool for http2.Transport.
// Hosts can be resolved to multiple ips.
// This implements http2.ClientConnPool interface.
type http2ConnPool struct {
	t *http2.Transport
	// mu protects conns and connMap.
	mu sync.RWMutex
	// conns is the addr and hostConns mapping.
	// The keys addr = host:port
	conns map[string]*hostConns
	// connMap is the connection to hostConn mapping.
	// This map is for markDead.
	connMap map[*http2.ClientConn]*hostConns
	// resolveInterval is the minimum resolve interval
	// used in dnsResolvers.
	resolveInterval time.Duration
}

func (p *http2ConnPool) GetClientConn(req *http.Request, addr string) (conn *http2.ClientConn, err error) {
	defer func() {
		if err != nil {
			return
		}
		if req.Close || httpguts.HeaderValuesContainsToken(req.Header["Connection"], "close") {
			conn.SetDoNotReuse()
		}
		p.mu.Lock()
		p.connMap[conn] = p.conns[addr]
		p.mu.Unlock()
	}()

	p.mu.RLock()
	conns, ok := p.conns[addr]
	p.mu.RUnlock()
	if ok {
		return conns.getClientConn(req.Context())
	}

	host, port, _ := net.SplitHostPort(addr)
	cfg := new(tls.Config)
	if p.t.TLSClientConfig != nil {
		*cfg = *p.t.TLSClientConfig.Clone()
	}
	if !slices.Contains(cfg.NextProtos, http2.NextProtoTLS) {
		cfg.NextProtos = append([]string{http2.NextProtoTLS}, cfg.NextProtos...)
	}
	if cfg.ServerName == "" {
		cfg.ServerName = host
	}

	conns = &hostConns{
		dnsResolver: &dnsResolver{
			host: host,
		},
		addr:     addr,
		port:     port,
		conns:    map[string][]*http2.ClientConn{},
		connToIP: map[*http2.ClientConn]string{},
		tc:       cfg,
		t:        p.t,
	}

	if err := conns.dnsResolver.resolve(req.Context()); err != nil {
		return nil, err
	}

	conn, err = conns.getClientConn(req.Context())
	if err != nil {
		return nil, err
	}

	// Start resolve loon in a new goroutine.
	// This goroutine will be stopped when MarkDead was called.
	//nolint:contextcheck // Non-inherited new context, use function like `context.WithXXX` or `r.Context` instead
	go conns.dnsResolver.resolveEveryInterval(context.Background(), p.resolveInterval)

	p.mu.Lock()
	p.conns[addr] = conns
	p.mu.Unlock()
	return conn, nil
}

func (p *http2ConnPool) MarkDead(conn *http2.ClientConn) {
	if VerboseLogs {
		log.Printf("network: connection marked dead %p\n", conn)
	}

	p.mu.RLock()
	conns, ok := p.connMap[conn]
	p.mu.RUnlock()
	if !ok {
		return
	}

	active := conns.markDead(conn)

	p.mu.Lock()
	defer p.mu.Unlock()
	if active == 0 {
		conns.dnsResolver.stopResolveLoop() // Must stop lookup goroutine to avoid leak.
		delete(p.conns, conns.addr)
	}
	delete(p.connMap, conn)
}
