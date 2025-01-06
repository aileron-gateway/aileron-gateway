package network

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"os"
	"reflect"
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/pion/dtls/v3"
)

// removeSocketConn wraps net.PacketConn
// which especially intended to net.UnixConn.
// removeSocketConn removes a unix domain docket file
// when Close method was called.
type removeSocketConn struct {
	net.PacketConn
	socket string
}

func (c *removeSocketConn) Close() error {
	defer func() {
		if info, err := os.Stat(c.socket); err == nil && !info.IsDir() {
			os.Remove(c.socket)
		}
	}()
	return c.PacketConn.Close()
}

// removeSocketListener wraps net.Listener
// which especially intended to net.UnixListener.
// removeSocketListener removes a unix domain docket file
// when Close method was called.
type removeSocketListener struct {
	net.Listener
	socket string
}

func (l *removeSocketListener) Close() error {
	defer func() {
		if info, err := os.Stat(l.socket); err == nil && !info.IsDir() {
			os.Remove(l.socket)
		}
	}()
	return l.Listener.Close()
}

func RemoveSocketConn(conn net.PacketConn) net.PacketConn {
	if conn == nil || reflect.ValueOf(conn).IsNil() || conn.LocalAddr() == nil {
		return conn
	}
	addr, ok := conn.LocalAddr().(*net.UnixAddr)
	if !ok || addr.Name == "" || strings.HasPrefix(addr.Name, "@") {
		return conn // Not a path name socket.
	}
	return &removeSocketConn{
		PacketConn: conn,
		socket:     addr.Name,
	}
}

func RemoveSocketListener(ln net.Listener) net.Listener {
	if ln == nil || reflect.ValueOf(ln).IsNil() || ln.Addr() == nil {
		return ln
	}
	addr, ok := ln.Addr().(*net.UnixAddr)
	if !ok || addr.Name == "" || strings.HasPrefix(addr.Name, "@") {
		return ln // Not a path name socket.
	}
	return &removeSocketListener{
		Listener: ln,
		socket:   addr.Name,
	}
}

// PacketConnConfig is the config for IP and UDP connection.
// TLS over UDP, or DTLS is not supported.
type PacketConnConfig struct {
	// Network is the network type to listen packet.
	// The network must be "udp", "udp4", "udp6", "unixgram", or an IP transport.
	// The IP transports are "ip", "ip4", or "ip6" followed by a colon
	// and a literal protocol number or a protocol name, as in "ip:1" or "ip:icmp".
	// See the net.ListenPacket documents.
	// https://pkg.go.dev/net#ListenPacket
	Network string

	// Address is the local address to listen to.
	// The value must be a valid address for the given network type.
	// See the net.Dial document for valid value examples.
	// https://pkg.go.dev/net#Dial
	// For unix type networks, "/var/run/example.sock" or "@example"
	// address can be specified.
	Address string

	// SockOption is the socket option.
	SockOption *SockOption
}

// NewPacketConn returns a new net.PacketConn from the given config.
// The nil config is equivalent to the zero config which means
//   - "udp" for the Network
//   - "127.0.0.1:0" for the Address
//   - Default TLS config defined by TLSConfig
func NewPacketConn(c *PacketConnConfig) (net.PacketConn, error) {
	if c == nil {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePackConn,
			Description: ErrDscPackConn,
			Detail:      "nil spec was given to new packet conn.",
		}
	}

	var conn net.PacketConn
	var syscallConn interface {
		SyscallConn() (syscall.RawConn, error)
	}

	var err error
	switch c.Network {
	case "udp", "udp4", "udp6":
		addr, resolveErr := net.ResolveUDPAddr(c.Network, c.Address)
		cn, listenErr := net.ListenUDP(c.Network, addr)
		conn, syscallConn = cn, cn
		err = errors.Join(listenErr, resolveErr)
	case "unixgram":
		addr, resolveErr := net.ResolveUnixAddr(c.Network, c.Address)
		cn, listenErr := net.ListenUnixgram(c.Network, addr)
		conn, syscallConn = cn, cn
		err = errors.Join(listenErr, resolveErr)
		conn = RemoveSocketConn(conn) // Remove .socket file when closed if necessary.
	default:
		if strings.HasPrefix(c.Network, "ip:") || strings.HasPrefix(c.Network, "ip4:") || strings.HasPrefix(c.Network, "ip6:") {
			addr, resolveErr := net.ResolveIPAddr(c.Network, c.Address)
			cn, listenErr := net.ListenIP(c.Network, addr)
			conn, syscallConn = cn, cn
			err = errors.Join(listenErr, resolveErr)
		} else {
			err = errors.New("unsupported network `" + c.Network + "`")
		}
	}

	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypePackConn,
			Description: ErrDscPackConn,
			Detail:      "create new packet conn.",
		}).Wrap(err)
	}

	controlFunc := c.SockOption.ControlFunc(SockOptSO | SockOptIP | SockOptIPV6 | SockOptUDP)
	if controlFunc != nil && syscallConn != nil {
		conn, err := syscallConn.SyscallConn()
		if err != nil {
			return nil, (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypePackConn,
				Description: ErrDscPackConn,
				Detail:      "create new packet conn.",
			}).Wrap(err)
		}
		if err := controlFunc("", "", conn); err != nil {
			return nil, (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypePackConn,
				Description: ErrDscPackConn,
				Detail:      "create new packet conn.",
			}).Wrap(err)
		}
	}

	return conn, nil
}

type ListenConfig struct {
	// Network is the network type to listen.
	// Known networks are "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only),
	// "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "unix", and "unixpacket".
	// See the net.Dial document for valid values.
	// https://pkg.go.dev/net#Dial
	Network string
	// Address is the local address to listen to.
	// The value must be a valid TCP or Unix address.
	// If not set, the address will be automatically chosen.
	// See the net.Dial document for valid values.
	// https://pkg.go.dev/net#Dial
	Address string

	// TLSConfig is the TLS configuration.
	// Network must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket".
	// otherwise, networking may not be work properly.
	TLSConfig *tls.Config
	// DTLSConfig is the TLS configuration.
	// Network must be "udp", "udp4" or "udp6".
	// otherwise, networking may not be work properly.
	DTLSConfig *dtls.Config

	// ConnectionLimit is the maximum number of connection
	// that allowed to connect at a point of time.
	// Connections that exceeds this limit until
	// they dis-connect the connections or
	// until they are accepted.
	ConnectionLimit int

	// ReadDeadline apply read deadline for connection.
	// If zero, the deadline is not explicitly applied.
	ReadDeadline time.Duration
	// WriteDeadline apply read deadline for connection.
	// If zero, the deadline is not explicitly applied.
	WriteDeadline time.Duration

	// Networks is the blacklist/whitelist of
	// networks that can be connected.
	// By default, Networks works as a whitelist.
	// To use as a blacklist, set Blacklist to true.
	Networks  []string
	Blacklist bool

	// SockOption is the socket option.
	// SocketOption is not applied for DTLS.
	SockOption *SockOption
}

// NewListenerFromSpec returns a new net.Listener from the given spec.
// This function returns nil listener and nil error when a nil spec was
// given as an argument.
// This function internally call network.NewListener function.
func NewListenerFromSpec(spec *k.ListenConfig) (net.Listener, error) {
	if spec == nil {
		return nil, nil
	}

	tlsConfig, err := TLSConfig(spec.TLSConfig)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeListener,
			Description: ErrDscListener,
			Detail:      "create new listener.",
		}).Wrap(err)
	}

	config := &ListenConfig{
		TLSConfig:       tlsConfig,
		DTLSConfig:      nil, // TODO: make DTLSConfig configurable.
		Network:         spec.Network,
		Address:         spec.Addr,
		ConnectionLimit: int(spec.ConnectionLimit),
		Networks:        spec.Networks,
		Blacklist:       spec.Blacklist,
		ReadDeadline:    time.Duration(spec.ReadDeadline) * time.Millisecond,
		WriteDeadline:   time.Duration(spec.WriteDeadline) * time.Millisecond,
		SockOption:      SockOptionFromSpec(spec.SockOption),
	}
	return NewListener(config)
}

// NewListener returns a new net.Listener from the given config.
// This function returns nil if a nil config was given.
// For example
//   - "tcp", "127.0.0.1:80"
//   - "tcp4", "127.0.0.1:80"
//   - "tcp6", "127.0.0.1:80"
//   - "unix", "/var/run/example.sock"
func NewListener(c *ListenConfig) (net.Listener, error) {
	if c == nil {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeListener,
			Description: ErrDscListener,
			Detail:      "nil spec was given to new listener.",
		}
	}

	lc := &net.ListenConfig{
		Control:         c.SockOption.ControlFunc(SockOptSO | SockOptIP | SockOptIPV6 | SockOptTCP),
		KeepAliveConfig: net.KeepAliveConfig{}, // TODO: make KeepAliveConfig configurable.
	}

	var ln net.Listener
	var err error
	switch c.Network {
	case "tcp", "tcp4", "tcp6":
		ln, err = lc.Listen(context.Background(), c.Network, c.Address)
		if c.TLSConfig != nil {
			ln = tls.NewListener(ln, c.TLSConfig)
		}
	case "udp", "udp4", "udp6":
		addr, resolveErr := net.ResolveUDPAddr(c.Network, c.Address)
		listener, listenErr := dtls.Listen(c.Network, addr, c.DTLSConfig)
		ln = listener
		err = errors.Join(listenErr, resolveErr)
	case "unix", "unixpacket":
		ln, err = lc.Listen(context.Background(), c.Network, c.Address)
		if c.TLSConfig != nil {
			ln = tls.NewListener(ln, c.TLSConfig)
		}
		ln = RemoveSocketListener(ln) // Remove .socket file when closed if necessary.
	default:
		err = errors.New("kernel/network: unsupported network `" + c.Network + "`")
	}
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeListener,
			Description: ErrDscListener,
			Detail:      "create new listener.",
		}).Wrap(err)
	}

	// Apply deadlines for connections.
	ln = ListenerWithReadDeadline(ln, c.ReadDeadline)
	ln = ListenerWithWriteDeadline(ln, c.WriteDeadline)

	// Apply whitelist/blacklist to the listener.
	list, err := netContainers(c.Networks)
	if err != nil {
		ln.Close() // Make sure to close internal listener.
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeListener,
			Description: ErrDscListener,
			Detail:      "create new listener.",
		}).Wrap(err)
	}
	ln = ListenerWithSecure(ln, list, c.Blacklist)

	// Apply connection limit is the limit is >0.
	ln = ListenerWithLimit(ln, c.ConnectionLimit)

	return ln, nil
}

// ListenerWithSecure returns a net listener with ip address and/or port filter.
// The given list will be considered as a whitelist when the blacklist=false,
// and blacklist when the blacklist=true.
// This function do nothing when the given Listener was nil
// or the length of given containers was 0.
func ListenerWithSecure(inner net.Listener, list []Container, blacklist bool) net.Listener {
	if inner == nil || len(list) == 0 {
		return inner
	}
	return &secureListener{
		Listener:  inner,
		list:      slices.Clip(list),
		blacklist: blacklist,
	}
}

// secureListener protects connections from not allowed clients.
// This listener allows or dis-allows connections by whitelist or blacklist.
// Closed connections will be returned by Accept() method when a connection
// was not allowed.
type secureListener struct {
	net.Listener
	// list is the list of ip or ip/port containers.
	// This list is used as whitelist or blacklist.
	// This list is used as whitelist by default and
	// blacklist when the "blacklist" is set to true.
	list []Container
	// blacklist is the flag to use the list as blacklist.
	// In most cases, securing by blacklist is not recommended.
	blacklist bool
}

func (l *secureListener) Accept() (net.Conn, error) {
	// Wait for the next connection.
	// Accept block the process until the next connection obtained.
	// If the listener is closed, Accept() unblock and return an error.
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeListener,
			Description: ErrDscListener,
			Detail:      "accept connection.",
		}).Wrap(err)
	}

	ip, port := splitHostPort(c.RemoteAddr().String())
	if ip == nil {
		c.Close() // Connection not allowed for invalid network.
		return c, nil
	}

	// When blacklist=false, dis-allowed by default(allowed=false) and allowed(allowed=true) if contained.
	// When blacklist=true, allowed by default(allowed=true) and dis-allowed(allowed=false) if contained.
	allowed := l.blacklist
	for _, item := range l.list {
		if item.Contains(ip, port) {
			if l.blacklist {
				allowed = false
			} else {
				allowed = true
			}
			break
		}
	}

	if !allowed {
		c.Close() // Connection not allowed.
	}
	return c, nil
}

// ListenerWithReadDeadline returns a net listener which set deadlines to the connection.
// This function do nothing when the given Listener was nil or the given deadline is zero.
// To set no deadline, set read to a negative value.
func ListenerWithReadDeadline(inner net.Listener, read time.Duration) net.Listener {
	if inner == nil || read == 0 {
		return inner
	}
	return &readDeadlineListener{
		Listener: inner,
		read:     read,
	}
}

// readDeadlineListener is the net listener which
// applies write deadline for the accepted connection.
type readDeadlineListener struct {
	net.Listener
	// read is the duration until read deadline.
	read time.Duration
}

func (l *readDeadlineListener) Accept() (net.Conn, error) {
	// Wait for the next connection.
	// Accept block the process until the next connection obtained.
	// If the listener is closed, Accept() unblock and return an error.
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeListener,
			Description: ErrDscListener,
			Detail:      "accept connection.",
		}).Wrap(err)
	}
	if l.read <= 0 {
		_ = c.SetReadDeadline(time.Time{}) // No timeouts.
	} else {
		_ = c.SetReadDeadline(time.Now().Add(l.read))
	}
	return c, nil
}

// ListenerWithWriteDeadline returns a net listener which set deadlines to the connection.
// This function do nothing when the given Listener was nil or the given deadline is zero.
// To set no deadline, set write to a negative value.
func ListenerWithWriteDeadline(inner net.Listener, write time.Duration) net.Listener {
	if inner == nil || write == 0 {
		return inner
	}
	return &writeDeadlineListener{
		Listener: inner,
		write:    write,
	}
}

// writeDeadlineListener is the net listener which
// applies write deadline for the accepted connection.
type writeDeadlineListener struct {
	net.Listener
	// write is the duration until write deadline.
	write time.Duration
}

func (l *writeDeadlineListener) Accept() (net.Conn, error) {
	// Wait for the next connection.
	// Accept block the process until the next connection obtained.
	// If the listener is closed, Accept() unblock and return an error.
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeListener,
			Description: ErrDscListener,
			Detail:      "accept connection.",
		}).Wrap(err)
	}
	if l.write <= 0 {
		_ = c.SetWriteDeadline(time.Time{}) // No timeouts.
	} else {
		_ = c.SetWriteDeadline(time.Now().Add(l.write))
	}
	return c, nil
}

// ListenerWithLimit returns a net listener with connection limit.
// This function do nothing when the given Listener was nil
// or the given limit was less than or equal to 0.
func ListenerWithLimit(inner net.Listener, limit int) net.Listener {
	if inner == nil || limit <= 0 {
		return inner
	}
	return &limitListener{
		Listener: inner,
		sem:      make(chan struct{}, limit), // limit must be positive.
	}
}

// limitListener is the net listener which can limit the number of simultaneous connections.
// Note that the linux command "netstat" or "ss" like below
// does not show the right number of connections currently accepted.
//   - netstat -uant | grep ESTABLISHED | grep 8080 | wc
//   - ss -o state established "( dport = :8080 )" -np | wc
//
// This is described in https://github.com/golang/go/issues/36212#issuecomment-567838193
// Use "lsof" command instead. For example,
//   - lsof -i:8080 | grep aileron
type limitListener struct {
	net.Listener
	// sem is the semaphore variable.
	// A new semaphore(struct{}) will be sent to this sem
	// when a new connection was acceptable.
	// The sent struct{} was removed from this sem when the connection was released.
	// The length of the sem must be grater than or equal to 1.
	sem chan struct{}
}

func (l *limitListener) Accept() (net.Conn, error) {
	l.sem <- struct{}{} // Accept semaphore.
	// Wait and accept the next connection.
	// Accept blocks the process until the next connection was obtained.
	// If the listener is closed, Accept() unblock and return an error.
	c, err := l.Listener.Accept()
	if err != nil {
		<-l.sem // Make sure to release accepted semaphore.
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeListener,
			Description: ErrDscListener,
			Detail:      "accept connection.",
		}).Wrap(err)
	}
	return &limitListenerConn{
		Conn:    c,
		release: func() { <-l.sem },
	}, nil
}

// limitListenerConn is the connection generated by limitListener with release function.
// release function is called when the connection is disconnected.
// Any occupied resources such as semaphore should be released in the release function.
type limitListenerConn struct {
	net.Conn
	once    sync.Once
	release func()
}

// Close closes the connection and release resources.
func (l *limitListenerConn) Close() error {
	err := l.Conn.Close()
	l.once.Do(l.release)
	return err
}
