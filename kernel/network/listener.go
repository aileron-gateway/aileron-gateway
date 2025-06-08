// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package network

import (
	"cmp"
	"context"
	"crypto/tls"
	"errors"
	"net"
	"os"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-projects/go/znet"
	"github.com/aileron-projects/go/zsyscall"
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
	SockOption *zsyscall.SockOption
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

	controlFunc := c.SockOption.ControlFunc(zsyscall.SockOptSO | zsyscall.SockOptIP | zsyscall.SockOptIPV6 | zsyscall.SockOptUDP)
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
	// Address is the local address to listen to.
	// The value must be a valid TCP or Unix address.
	// If not set, the address will be automatically chosen.
	// Address can have network type prefix in the format of "<Network>://<Address>".
	// Known networks are "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only),
	// "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "unix", and "unixpacket".
	// See the net.Dial document for valid values.
	// https://pkg.go.dev/net#Dial
	Address string
	// TLSConfig is the TLS configuration.
	// Network must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket".
	// otherwise, networking may not be work properly.
	TLSConfig *tls.Config
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
	Networks []string
	// KeepAliveConfig is the keep-alive configuration.
	// This config is not used for DTLS.
	KeepAliveConfig net.KeepAliveConfig
	// SockOption is the socket option.
	// SocketOption is not applied for DTLS.
	SockOption *zsyscall.SockOption
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
		Address:         spec.Addr,
		ConnectionLimit: int(spec.ConnectionLimit),
		Networks:        spec.Networks,
		ReadDeadline:    time.Duration(spec.ReadDeadline) * time.Millisecond,
		WriteDeadline:   time.Duration(spec.WriteDeadline) * time.Millisecond,
		SockOption:      SockOptionFromSpec(spec.SockOption),
	}
	if spec.KeepAliveConfig != nil {
		kc := spec.KeepAliveConfig
		if kc.Disable {
			config.KeepAliveConfig = net.KeepAliveConfig{
				Enable: false,            // Use net.ListenConfig.KeepAlive rather than net.ListenConfig.KeepAliveConfig.
				Idle:   -1 * time.Second, // Negative value of net.ListenConfig.KeepAlive disables keep-alive.
			}
		} else {
			config.KeepAliveConfig = net.KeepAliveConfig{
				Enable:   true,
				Idle:     time.Second * time.Duration(kc.Idle),
				Interval: time.Second * time.Duration(kc.Interval),
				Count:    int(kc.Count),
			}
		}
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
		Control:         c.SockOption.ControlFunc(zsyscall.SockOptSO | zsyscall.SockOptIP | zsyscall.SockOptIPV6 | zsyscall.SockOptTCP),
		KeepAlive:       c.KeepAliveConfig.Idle,
		KeepAliveConfig: c.KeepAliveConfig,
	}

	var ln net.Listener
	var err error
	net, addr := znet.ParseNetAddr(c.Address)
	switch net {
	case "", "tcp", "tcp4", "tcp6":
		net = cmp.Or(net, "tcp") // Default tcp.
		ln, err = lc.Listen(context.Background(), net, addr)
		if c.TLSConfig != nil {
			ln = tls.NewListener(ln, c.TLSConfig)
		}
	case "unix", "unixpacket":
		ln, err = lc.Listen(context.Background(), net, addr)
		if c.TLSConfig != nil {
			ln = tls.NewListener(ln, c.TLSConfig)
		}
		ln = RemoveSocketListener(ln) // Remove .socket file when closed if necessary.
	default:
		err = errors.New("kernel/network: unknown address `" + c.Address + "`")
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
	if len(c.Networks) > 0 {
		wln, err := znet.NewWhiteListListener(ln, c.Networks...)
		if err != nil {
			ln.Close() // Make sure to close internal listener.
			return nil, (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeListener,
				Description: ErrDscListener,
				Detail:      "create new listener.",
			}).Wrap(err)
		}
		ln = wln
	}
	if c.ConnectionLimit > 0 { // Connection limit if >0.
		ln = znet.NewLimitListener(ln, c.ConnectionLimit)
	}
	return ln, nil
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

func SockOptionFromSpec(spec *kernel.SockOption) *zsyscall.SockOption {
	if spec == nil {
		return nil
	}
	return &zsyscall.SockOption{
		SO:   SockSOOptionFromSpec(spec.SOOption),
		IP:   SockIPOptionFromSpec(spec.IPOption),
		IPV6: SockIPV6OptionFromSpec(spec.IPV6Option),
		TCP:  SockTCPOptionFromSpec(spec.TCPOption),
		UDP:  SockUDPOptionFromSpec(spec.UDPOption),
	}
}

func SockSOOptionFromSpec(spec *kernel.SockSOOption) *zsyscall.SockSOOption {
	if spec == nil {
		return nil
	}
	return &zsyscall.SockSOOption{
		BindToDevice:       spec.BindToDevice,
		Debug:              spec.Debug,
		KeepAlive:          spec.KeepAlive,
		Linger:             spec.Linger,
		Mark:               int(spec.Mark),
		ReceiveBuffer:      int(spec.ReceiveBuffer),
		ReceiveBufferForce: int(spec.ReceiveBufferForce),
		ReceiveTimeout:     time.Duration(spec.ReceiveTimeout) * time.Millisecond,
		SendTimeout:        time.Duration(spec.SendTimeout) * time.Millisecond,
		ReuseAddr:          spec.ReuseAddr,
		ReusePort:          spec.ReusePort,
		SendBuffer:         int(spec.SendBuffer),
		SendBufferForce:    int(spec.SendBufferForce),
	}
}

func SockIPOptionFromSpec(spec *kernel.SockIPOption) *zsyscall.SockIPOption {
	if spec == nil {
		return nil
	}
	return &zsyscall.SockIPOption{
		BindAddressNoPort:   spec.BindAddressNoPort,
		FreeBind:            spec.FreeBind,
		LocalPortRangeUpper: uint16(spec.LocalPortRangeUpper), //nolint:gosec // G115: integer overflow conversion int32 -> uint16
		LocalPortRangeLower: uint16(spec.LocalPortRangeLower), //nolint:gosec // G115: integer overflow conversion int32 -> uint16
		Transparent:         spec.Transparent,
		TTL:                 int(spec.TTL),
	}
}

func SockIPV6OptionFromSpec(spec *kernel.SockIPV6Option) *zsyscall.SockIPV6Option {
	if spec == nil {
		return nil
	}
	return &zsyscall.SockIPV6Option{}
}

func SockTCPOptionFromSpec(spec *kernel.SockTCPOption) *zsyscall.SockTCPOption {
	if spec == nil {
		return nil
	}
	return &zsyscall.SockTCPOption{
		CORK:            spec.CORK,
		DeferAccept:     int(spec.DeferAccept),
		KeepCount:       int(spec.KeepCount),
		KeepIdle:        int(spec.KeepIdle),
		KeepInterval:    int(spec.KeepInterval),
		Linger2:         spec.Linger2,
		MaxSegment:      int(spec.MaxSegment),
		NoDelay:         spec.NoDelay,
		QuickAck:        spec.QuickAck,
		SynCount:        int(spec.SynCount),
		UserTimeout:     int(spec.UserTimeout),
		WindowClamp:     int(spec.WindowClamp),
		FastOpen:        spec.FastOpen,
		FastOpenConnect: spec.FastOpenConnect,
	}
}

func SockUDPOptionFromSpec(spec *kernel.SockUDPOption) *zsyscall.SockUDPOption {
	if spec == nil {
		return nil
	}
	return &zsyscall.SockUDPOption{
		CORK:    spec.CORK,
		Segment: int(spec.Segment),
		GRO:     spec.GRO,
	}
}
