package network

import (
	"errors"
	"syscall"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
)

const (
	SockOptSO = 1 << iota
	SockOptIP
	SockOptIPV6
	SockOptTCP
	SockOptUDP
)

var (
	// setsockoptInt is the re-mapped function of syscall.SetsockoptInt for testing.
	setsockoptInt = syscall.SetsockoptInt
	// setsockoptLinger is the re-mapped function of syscall.SetsockoptLinger for testing.
	setsockoptLinger = syscall.SetsockoptLinger
	// setsockoptTimeval is the re-mapped function of syscall.SetsockoptTimeval for testing.
	setsockoptTimeval = syscall.SetsockoptTimeval
)

// Controller is a function type that controls socket.
type Controller func(fd uintptr) error

// ControlFunc is the function type that handle RawConn.
type ControlFunc func(string, string, syscall.RawConn) error

// SockOption is aggregated socket options.
//   - https://man7.org/linux/man-pages/man7/socket.7.html
//   - https://man7.org/linux/man-pages/man7/ip.7.html
//   - https://man7.org/linux/man-pages/man7/ipv6.7.html
//   - https://man7.org/linux/man-pages/man7/tcp.7.html
//   - https://man7.org/linux/man-pages/man7/udp.7.html
type SockOption struct {
	SO   *SockSOOption
	IP   *SockIPOption
	IPV6 *SockIPV6Option
	TCP  *SockTCPOption
	UDP  *SockUDPOption
}

func (o *SockOption) ControlFunc(opts int) ControlFunc {
	if o == nil {
		return nil
	}
	cs := o.Controllers(opts)
	if len(cs) == 0 {
		return nil
	}
	return func(network, address string, conn syscall.RawConn) error {
		var errs []error
		err := conn.Control(func(fd uintptr) {
			for _, control := range cs {
				errs = append(errs, control(fd))
			}
		})
		return errors.Join(append([]error{err}, errs...)...)
	}
}

func (o *SockOption) Controllers(opts int) []Controller {
	var cs []Controller
	if o.SO != nil && opts&SockOptSO != 0 {
		cs = append(cs, o.SO.Controllers()...)
	}
	if o.IP != nil && opts&SockOptIP != 0 {
		cs = append(cs, o.IP.Controllers()...)
	}
	if o.IPV6 != nil && opts&SockOptIPV6 != 0 {
		cs = append(cs, o.IPV6.Controllers()...)
	}
	if o.TCP != nil && opts&SockOptTCP != 0 {
		cs = append(cs, o.TCP.Controllers()...)
	}
	if o.UDP != nil && opts&SockOptUDP != 0 {
		cs = append(cs, o.UDP.Controllers()...)
	}
	return cs
}

func SockOptionFromSpec(spec *k.SockOption) *SockOption {
	if spec == nil {
		return nil
	}
	return &SockOption{
		SO:   SockSOOptionFromSpec(spec.SOOption),
		IP:   SockIPOptionFromSpec(spec.IPOption),
		IPV6: SockIPV6OptionFromSpec(spec.IPV6Option),
		TCP:  SockTCPOptionFromSpec(spec.TCPOption),
		UDP:  SockUDPOptionFromSpec(spec.UDPOption),
	}
}

// SockSOOption is socket options for SOL_SOCKET level.
//   - https://man7.org/linux/man-pages/man7/socket.7.html
//   - https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/sys_socket.h.html
type SockSOOption struct {
	BindToDevice       string  // SO_BINDTODEVICE
	Debug              bool    // SO_DEBUG
	IncomingCPU        bool    // SO_INCOMING_CPU (gettable since Linux 3.19, settable since Linux 4.4)
	KeepAlive          bool    // SO_KEEPALIVE
	Linger             int32   // SO_LINGER
	Mark               int     // SO_MARK (since Linux 2.6.25)
	ReceiveBuffer      int     // SO_RCVBUF
	ReceiveBufferForce int     // SO_RCVBUFFORCE (since Linux 2.6.14)
	ReceiveTimeout     float64 // SO_RCVTIMEO
	SendTimeout        float64 // SO_SNDTIMEO
	ReuseAddr          bool    // SO_REUSEADDR
	ReusePort          bool    // SO_REUSEPORT (since Linux 3.9)
	SendBuffer         int     // SO_SNDBUF
	SendBufferForce    int     // SO_SNDBUFFORCE (since Linux 2.6.14)
}

func SockSOOptionFromSpec(spec *k.SockSOOption) *SockSOOption {
	if spec == nil {
		return nil
	}
	return &SockSOOption{
		BindToDevice:       spec.BindToDevice,
		Debug:              spec.Debug,
		IncomingCPU:        spec.IncomingCPU,
		KeepAlive:          spec.KeepAlive,
		Linger:             spec.Linger,
		Mark:               int(spec.Mark),
		ReceiveBuffer:      int(spec.ReceiveBuffer),
		ReceiveBufferForce: int(spec.ReceiveBufferForce),
		ReceiveTimeout:     spec.ReceiveTimeout,
		SendTimeout:        spec.SendTimeout,
		ReuseAddr:          spec.ReuseAddr,
		ReusePort:          spec.ReusePort,
		SendBuffer:         int(spec.SendBuffer),
		SendBufferForce:    int(spec.SendBufferForce),
	}
}

// SockIPOption is socket options for IPPROTO_IP level.
//   - https://man7.org/linux/man-pages/man7/ip.7.html
type SockIPOption struct {
	BindAddressNoPort   bool // IP_BIND_ADDRESS_NO_PORT (since Linux 4.2)
	FreeBind            bool // IP_FREEBIND (since Linux 2.4)
	LocalPortRangeUpper int  // IP_LOCAL_PORT_RANGE (since Linux 6.3)
	LocalPortRangeLower int  // IP_LOCAL_PORT_RANGE (since Linux 6.3)
	Transparent         bool // IP_TRANSPARENT (since Linux 2.6.24)
	TTL                 int  // IP_TTL (since Linux 1.0)
}

func SockIPOptionFromSpec(spec *k.SockIPOption) *SockIPOption {
	if spec == nil {
		return nil
	}
	return &SockIPOption{
		BindAddressNoPort:   spec.BindAddressNoPort,
		FreeBind:            spec.FreeBind,
		LocalPortRangeUpper: int(spec.LocalPortRangeUpper),
		LocalPortRangeLower: int(spec.LocalPortRangeLower),
		Transparent:         spec.Transparent,
		TTL:                 int(spec.TTL),
	}
}

// SockIPV6Option is socket options for IPPROTO_IPV6 level.
//   - https://man7.org/linux/man-pages/man7/ipv6.7.html
type SockIPV6Option struct {
}

func SockIPV6OptionFromSpec(spec *k.SockIPV6Option) *SockIPV6Option {
	if spec == nil {
		return nil
	}
	return &SockIPV6Option{}
}

// SockTCPOption is socket options for IPPROTO_TCP level.
//   - https://man7.org/linux/man-pages/man7/tcp.7.html
type SockTCPOption struct {
	CORK            bool  // TCP_CORK (since Linux 2.2)
	DeferAccept     int   // TCP_DEFER_ACCEPT (since Linux 2.4)
	KeepCount       int   // TCP_KEEPCNT (since Linux 2.4)
	KeepIdle        int   // TCP_KEEPIDLE (since Linux 2.4)
	KeepInterval    int   // TCP_KEEPINTVL (since Linux 2.4)
	Linger2         int32 // TCP_LINGER2 (since Linux 2.4)
	MaxSegment      int   // TCP_MAXSEG
	NoDelay         bool  // TCP_NODELAY
	QuickAck        bool  // TCP_QUICKACK (since Linux 2.4.4)
	SynCount        int   // TCP_SYNCNT (since Linux 2.4)
	UserTimeout     int   // TCP_USER_TIMEOUT (since Linux 2.6.37)
	WindowClamp     int   // TCP_WINDOW_CLAMP (since Linux 2.4)
	FastOpen        bool  // TCP_FASTOPEN (since Linux 3.6)
	FastOpenConnect bool  // TCP_FASTOPEN_CONNECT (since Linux 4.11)
}

func SockTCPOptionFromSpec(spec *k.SockTCPOption) *SockTCPOption {
	if spec == nil {
		return nil
	}
	return &SockTCPOption{
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

// SockUDPOption is socket options for IPPROTO_UDP level.
//   - https://man7.org/linux/man-pages/man7/udp.7.html
type SockUDPOption struct {
	CORK    bool // UDP_CORK (since Linux 2.5.44)
	Segment int  // UDP_SEGMENT (since Linux 4.18)
	GRO     bool // UDP_GRO (since Linux 5.0)
}

func SockUDPOptionFromSpec(spec *k.SockUDPOption) *SockUDPOption {
	if spec == nil {
		return nil
	}
	return &SockUDPOption{
		CORK:    spec.CORK,
		Segment: int(spec.Segment),
		GRO:     spec.GRO,
	}
}

func appendIfNotNil(arr []Controller, target Controller) []Controller {
	if target == nil {
		return arr
	}
	return append(arr, target)
}
