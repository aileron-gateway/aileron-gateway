package network

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"regexp"
	"strings"
	"syscall"
	"time"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/pion/dtls/v3"
)

// NetworkType is the type that the net.Dialer can accept.
// See https://pkg.go.dev/net#Dial.
func NetworkType(t k.NetworkType) string {
	nw := map[k.NetworkType]string{
		k.NetworkType_HTTP:       "tcp", // HTTP is an alias for "tcp".
		k.NetworkType_TCP:        "tcp",
		k.NetworkType_TCP4:       "tcp4",
		k.NetworkType_TCP6:       "tcp6",
		k.NetworkType_UDP:        "udp",
		k.NetworkType_UDP4:       "udp4",
		k.NetworkType_UDP6:       "udp6",
		k.NetworkType_IP:         "ip",
		k.NetworkType_IP4:        "ip4",
		k.NetworkType_IP6:        "ip6",
		k.NetworkType_Unix:       "unix",
		k.NetworkType_UnixGram:   "unixgram",
		k.NetworkType_UnixPacket: "unixpacket",
	}
	return nw[t]
}

// replaceTargetDialer is a dialer that
// replaces dial target to a fixed value.
// This is intended to be used for
// HTTP communication over Unis Sockets.
// This implements Dialer interface.
type replaceTargetDialer struct {
	Dialer
	// Following slices are converting maps.
	// They all have the same length.
	// (fromNet[i], fromAddr[i]) > (toNet[i], toAddr[i])
	fromNet  []string
	fromAddr []string
	toNet    []string
	toAddr   []string
}

func (d *replaceTargetDialer) Dial(network string, address string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, address)
}

func (d *replaceTargetDialer) DialContext(ctx context.Context, network string, address string) (net.Conn, error) {
	for i := 0; i < len(d.fromAddr); i++ {
		if address == d.fromAddr[i] && network == d.fromNet[i] {
			network = d.toNet[i]
			address = d.toAddr[i]
		}
	}
	return d.Dialer.DialContext(ctx, network, address)
}

func newReplaceTargetDialer(replaces []string, dialer Dialer) (Dialer, error) {
	if len(replaces) == 0 {
		return dialer, nil
	}

	pattern := regexp.MustCompile(`^\(([\w:\-]+)\|(.+)\)[ >]*\(([\w:\-]+)\|(.+)\)$`)
	d := &replaceTargetDialer{
		Dialer:   dialer,
		fromNet:  make([]string, 0, len(replaces)),
		fromAddr: make([]string, 0, len(replaces)),
		toNet:    make([]string, 0, len(replaces)),
		toAddr:   make([]string, 0, len(replaces)),
	}

	uniqueCheck := map[string]struct{}{}
	for _, r := range replaces {
		r = strings.ReplaceAll(r, " ", "") // Remove white spaces.
		matches := pattern.FindStringSubmatch(r)
		if len(matches) != 5 {
			return nil, errors.New("invalid replace pattern. Pattern must be `" + pattern.String() + "`")
		}
		_, fromErr := resolveAddr(matches[1], matches[2]) // Check if the target is valid.
		_, toErr := resolveAddr(matches[3], matches[4])   // Check if the target is valid.
		if err := errors.Join(fromErr, toErr); err != nil {
			return nil, err
		}
		d.fromNet = append(d.fromNet, matches[1])
		d.fromAddr = append(d.fromAddr, matches[2])
		d.toNet = append(d.toNet, matches[3])
		d.toAddr = append(d.toAddr, matches[4])
		if _, ok := uniqueCheck[matches[1]+matches[2]]; ok {
			return nil, errors.New("replace source duplication. duplicate (" + matches[1] + "|" + matches[2] + ")")
		}
		uniqueCheck[matches[1]+matches[2]] = struct{}{}
	}

	return d, nil
}

func resolveAddr(network, address string) (net.Addr, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		return net.ResolveTCPAddr(network, address)
	case "udp", "udp4", "udp6":
		return net.ResolveUDPAddr(network, address)
	case "unix", "unixgram", "unixpacket":
		return net.ResolveUnixAddr(network, address)
	case "":
		return nil, nil
	default:
		if strings.HasPrefix(network, "ip:") || strings.HasPrefix(network, "ip4:") || strings.HasPrefix(network, "ip6:") {
			return net.ResolveIPAddr(network, address)
		} else {
			return nil, errors.New("invalid network `" + network + "`")
		}
	}
}

// NewDialerFromSpec returns a new net.Dialer from the given spec.
// This function returns an error if nil spec was given by the argument.
// This function internally calls NewDialer.
func NewDialerFromSpec(spec *k.DialConfig) (Dialer, error) {
	if spec == nil {
		return nil, nil
	}

	tlsConfig, err := TLSConfig(spec.TLSConfig)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDialer,
			Description: ErrDscDialer,
			Detail:      "create new dialer.",
		}).Wrap(err)
	}

	config := &DialConfig{
		TLSConfig:      tlsConfig,
		LocalNetwork:   spec.LocalNetwork,
		LocalAddress:   spec.LocalAddress,
		ReplaceTargets: spec.ReplaceTargets,
		Timeout:        time.Duration(spec.Timeout) * time.Millisecond,
		FallbackDelay:  time.Duration(spec.FallbackDelay) * time.Millisecond,
		SockOption:     SockOptionFromSpec(spec.SockOption),
	}
	return NewDialer(config)
}

// NewDialer returns a new Dialer from the given config.
// This function creates *net.Dialer when TLSConfig is not set
// and creates *tls.Dialer if TLSConfig is set.
// This function returns an error if nil config was given by an argument.
func NewDialer(c *DialConfig) (Dialer, error) {
	if c == nil {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDialer,
			Description: ErrDscDialer,
			Detail:      "nil spec was given to new dialer.",
		}
	}

	var localAddr net.Addr
	if c.LocalNetwork != "" && c.LocalAddress != "" {
		addr, err := resolveAddr(c.LocalNetwork, c.LocalAddress)
		if err != nil {
			return nil, (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeDialer,
				Description: ErrDscDialer,
				Detail:      "create new dialer.",
			}).Wrap(err)
		}
		localAddr = addr
	}

	var d Dialer
	if c.DTLSConfig != nil {
		d = &dtlsDialer{
			tlsConfig:    c.DTLSConfig,
			localNetwork: c.LocalNetwork,
			localAddress: c.LocalAddress,
			control:      c.SockOption.ControlFunc(SockOptSO | SockOptIP | SockOptIPV6 | SockOptTCP | SockOptUDP),
		}
	} else {
		dd := &net.Dialer{
			LocalAddr:     localAddr,
			FallbackDelay: c.FallbackDelay,
			Timeout:       c.Timeout,
			Control:       c.SockOption.ControlFunc(SockOptSO | SockOptIP | SockOptIPV6 | SockOptTCP | SockOptUDP),
		}
		d = dd
		if c.TLSConfig != nil {
			d = &tls.Dialer{
				NetDialer: dd,
				Config:    c.TLSConfig,
			}
		}
	}

	var err error
	d, err = newReplaceTargetDialer(c.ReplaceTargets, d)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDialer,
			Description: ErrDscDialer,
			Detail:      "create new dialer.",
		}).Wrap(err)
	}

	return d, nil
}

// DialConfig is the config for net.Dialer.
// Supported networks is the same as net.Dialer.
type DialConfig struct {
	// TLSConfig is the tls configuration used when dialing
	// to an address.
	// net.Dialer will be used when not set,
	// and tls.Dialer will be used when set.
	TLSConfig *tls.Config
	// TLSConfig is the TLS configuration to use for new connections.
	// A nil configuration is equivalent to the zero configuration.
	DTLSConfig *dtls.Config

	// LocalNetwork is the local network type to listen.
	// https://pkg.go.dev/net#Dialer
	LocalNetwork string
	// LocalAddress is the local address to listen to.
	LocalAddress string

	// ReplaceTargets is the dial target replaces.
	// The format must be
	// `(<FromNetwork>|<FromAddress>) (<ToNetwork>|<ToAddress>)`
	// For example,
	// `(tcp|example.com:80) (unix|/var/run/example.sock)` or
	// `(tcp|example.com:80) (unix|@example)`.
	// Supported values of networks and addresses
	// follow the specification of net.Dialer.
	// See https://pkg.go.dev/net#Dialer
	// As described in the net.Dialer comments,
	// TCP and UDP must have a port number.
	ReplaceTargets []string

	// Timeout is the Timeout field of net.Dialer.
	// See https://pkg.go.dev/net#Dialer
	Timeout time.Duration
	// FallbackDelay is the FallbackDelay field of net.Dialer.
	// FallbackDelay is not used for DTLS.
	// See https://pkg.go.dev/net#Dialer
	FallbackDelay time.Duration
	// SockOption is the socket option.
	SockOption *SockOption
}

// dtlsDialer is the dialer for DTLS.
// This implements Dialer interface.
type dtlsDialer struct {
	tlsConfig    *dtls.Config
	localNetwork string
	localAddress string
	control      ControlFunc
}

// Dial connects to the address with DTL.
// Networks have to be one of the udp types.
// Those are "udp", "udp4" (IPv4-only), "udp6" (IPv6-only) or "unixgram".
// For example:
//   - Dial("udp", "[2001:db8::1]:8443")
//   - Dial("udp", ":8443")
func (d *dtlsDialer) Dial(network string, address string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, address)
}

// DialContext connects to the address with the provided context.
// The provided Context must be non-nil. If the context expires before
// the connection is complete, an error is returned. Once successfully
// connected, any expiration of the context will not affect the connection.
func (d *dtlsDialer) DialContext(ctx context.Context, network string, address string) (net.Conn, error) {
	var conn net.PacketConn
	var syscallConn interface {
		SyscallConn() (syscall.RawConn, error)
	}
	var addr net.Addr
	var err error

	// For DTLS, Network must be one of "udp", "udp4", "udp6".
	switch network {
	case "udp", "udp4", "udp6":
		rAddr, resolveErr1 := net.ResolveUDPAddr(network, address)
		lAddr, resolveErr2 := net.ResolveUDPAddr(d.localNetwork, d.localAddress)
		udpConn, listenErr := net.ListenUDP(network, lAddr)
		err = errors.Join(listenErr, resolveErr1, resolveErr2)
		syscallConn = udpConn
		conn = udpConn
		addr = rAddr
	default:
		err = errors.New("network not supported by DTLS `" + network + "`")
	}
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDialer,
			Description: ErrDscDialer,
			Detail:      "dial DTLS.",
		}).Wrap(err)
	}

	if d.control != nil && syscallConn != nil {
		conn, err := syscallConn.SyscallConn()
		if err != nil {
			return nil, (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeDialer,
				Description: ErrDscDialer,
				Detail:      "dial DTLS.",
			}).Wrap(err)
		}
		if err := d.control("", "", conn); err != nil {
			return nil, (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeDialer,
				Description: ErrDscDialer,
				Detail:      "dial DTLS.",
			}).Wrap(err)
		}
	}

	dtlsConn, err := dtls.Client(conn, addr, d.tlsConfig)
	if err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeDialer,
			Description: ErrDscDialer,
			Detail:      "create QUIC dialer.",
		}).Wrap(err)
	}
	return dtlsConn, nil
}
