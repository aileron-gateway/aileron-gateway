package network

import (
	"context"
	"net"
	"net/netip"

	quicgo "github.com/quic-go/quic-go"
	"golang.org/x/net/quic"
)

// Dialer is an interface for dialer, or network client.
// This interface is intended to be used for *net.Dialer and *tls.Dialer
// and *DTLSDialer.
type Dialer interface {
	// Dial connects to the address on the named network.
	// Known networks are "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only),
	// "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4"
	// (IPv4-only), "ip6" (IPv6-only), "unix", "unixgram" and "unixpacket".
	// See https://pkg.go.dev/net#Dialer for more information.
	Dial(network string, address string) (net.Conn, error)

	// DialContext connects to the address on the named network using
	// the provided context.
	// Known networks are "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only),
	// "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4"
	// (IPv4-only), "ip6" (IPv6-only), "unix", "unixgram" and "unixpacket".
	// See https://pkg.go.dev/net#Dialer for more information.
	//
	// The provided Context must be non-nil. If the context expires before
	// the connection is complete, an error is returned. Once successfully
	// connected, any expiration of the context will not affect the
	// connection.
	DialContext(ctx context.Context, network string, address string) (net.Conn, error)
}

// QUICDialer is the dialer of QUIC protocol.
// QUICDialer is related to the "golang.org/x/net/quic" package
// while QuicDialer is related to "github.com/quic-go/quic-go".
type QUICDialer interface {
	// Dial connects to the address on the named network.
	// Network must be a UDP type "udp", "udp4" (IPv4-only), "udp6" (IPv6-only).
	// See https://pkg.go.dev/net#ResolveUDPAddr.
	Dial(network string, address string) (*quic.Conn, error)

	// DialContext connects to the address on the named network using
	// the provided context.
	// Network must be a UDP type "udp", "udp4" (IPv4-only), "udp6" (IPv6-only).
	// See https://pkg.go.dev/net#ResolveUDPAddr.
	//
	// The provided Context must be non-nil. If the context expires before
	// the connection is complete, an error is returned. Once successfully
	// connected, any expiration of the context will not affect the
	// connection.
	DialContext(ctx context.Context, network string, address string) (*quic.Conn, error)
}

// QuicDialer is the dialer of QUIC protocol.
// QUICDialer is related to the "golang.org/x/net/quic" package
// while QuicDialer is related to "github.com/quic-go/quic-go".
type QuicDialer interface {
	// Dial connects to the address on the named network.
	// Network must be a UDP type "udp", "udp4" (IPv4-only), "udp6" (IPv6-only).
	// See https://pkg.go.dev/net#ResolveUDPAddr.
	Dial(network string, addr string) (quicgo.Connection, error)

	// DialContext connects to the address on the named network using
	// the provided context.
	// Network must be a UDP type "udp", "udp4" (IPv4-only), "udp6" (IPv6-only).
	// See https://pkg.go.dev/net#ResolveUDPAddr.
	//
	// The provided Context must be non-nil. If the context expires before
	// the connection is complete, an error is returned. Once successfully
	// connected, any expiration of the context will not affect the
	// connection.
	DialContext(ctx context.Context, network string, addr string) (quicgo.Connection, error)
}

// QUICListener is the listener of QUIC protocol.
// QUICListener is related to the "golang.org/x/net/quic" package
// while QuicListener is related to "github.com/quic-go/quic-go".
type QUICListener interface {
	Accept(context.Context) (*quic.Conn, error)
	Close(context.Context) error
	LocalAddr() netip.AddrPort
}

// QuicListener is the listener of QUIC protocol.
// QUICListener is related to the "golang.org/x/net/quic" package
// while QuicListener is related to "github.com/quic-go/quic-go".
type QuicListener interface {
	Accept(context.Context) (quicgo.Connection, error)
	Close() error
	Addr() net.Addr
}
