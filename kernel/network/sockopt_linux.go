//go:build linux

package network

import (
	"cmp"
	"encoding/binary"
	"runtime"
	"syscall"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
)

// TODO: Remove unnecessary options. Make options configurable.
const (
	_ = syscall.SO_DONTROUTE                     //= 0x5
	_ = syscall.SO_RCVLOWAT                      //= 0x12
	_ = syscall.SO_SECURITY_AUTHENTICATION       //= 0x16
	_ = syscall.SO_SECURITY_ENCRYPTION_NETWORK   //= 0x18
	_ = syscall.SO_SECURITY_ENCRYPTION_TRANSPORT //= 0x17
	_ = syscall.SO_SNDLOWAT                      //= 0x13

	_ = syscall.IPPROTO_AH       //= 0x33
	_ = syscall.IPPROTO_COMP     //= 0x6c
	_ = syscall.IPPROTO_DCCP     //= 0x21
	_ = syscall.IPPROTO_DSTOPTS  //= 0x3c
	_ = syscall.IPPROTO_EGP      //= 0x8
	_ = syscall.IPPROTO_ENCAP    //= 0x62
	_ = syscall.IPPROTO_ESP      //= 0x32
	_ = syscall.IPPROTO_FRAGMENT //= 0x2c
	_ = syscall.IPPROTO_GRE      //= 0x2f
	_ = syscall.IPPROTO_HOPOPTS  //= 0x0
	_ = syscall.IPPROTO_ICMP     //= 0x1
	_ = syscall.IPPROTO_ICMPV6   //= 0x3a
	_ = syscall.IPPROTO_IDP      //= 0x16
	_ = syscall.IPPROTO_IGMP     //= 0x2
	_ = syscall.IPPROTO_IP       //= 0x0
	_ = syscall.IPPROTO_IPIP     //= 0x4
	_ = syscall.IPPROTO_IPV6     //= 0x29
	_ = syscall.IPPROTO_MTP      //= 0x5c
	_ = syscall.IPPROTO_NONE     //= 0x3b
	_ = syscall.IPPROTO_PIM      //= 0x67
	_ = syscall.IPPROTO_PUP      //= 0xc
	_ = syscall.IPPROTO_RAW      //= 0xff
	_ = syscall.IPPROTO_ROUTING  //= 0x2b
	_ = syscall.IPPROTO_RSVP     //= 0x2e
	_ = syscall.IPPROTO_SCTP     //= 0x84
	_ = syscall.IPPROTO_TCP      //= 0x6
	_ = syscall.IPPROTO_TP       //= 0x1d
	_ = syscall.IPPROTO_UDP      //= 0x11
	_ = syscall.IPPROTO_UDPLITE  //= 0x88

	_ = syscall.IPV6_2292DSTOPTS     //= 0x4
	_ = syscall.IPV6_2292HOPLIMIT    //= 0x8
	_ = syscall.IPV6_2292HOPOPTS     //= 0x3
	_ = syscall.IPV6_2292PKTINFO     //= 0x2
	_ = syscall.IPV6_2292PKTOPTIONS  //= 0x6
	_ = syscall.IPV6_2292RTHDR       //= 0x5
	_ = syscall.IPV6_ADDRFORM        //= 0x1
	_ = syscall.IPV6_ADD_MEMBERSHIP  //= 0x14
	_ = syscall.IPV6_AUTHHDR         //= 0xa
	_ = syscall.IPV6_CHECKSUM        //= 0x7
	_ = syscall.IPV6_DROP_MEMBERSHIP //= 0x15
	_ = syscall.IPV6_DSTOPTS         //= 0x3b
	_ = syscall.IPV6_HOPLIMIT        //= 0x34
	_ = syscall.IPV6_HOPOPTS         //= 0x36
	_ = syscall.IPV6_IPSEC_POLICY    //= 0x22
	_ = syscall.IPV6_JOIN_ANYCAST    //= 0x1b
	_ = syscall.IPV6_JOIN_GROUP      //= 0x14
	_ = syscall.IPV6_LEAVE_ANYCAST   //= 0x1c
	_ = syscall.IPV6_LEAVE_GROUP     //= 0x15
	_ = syscall.IPV6_MTU             //= 0x18
	_ = syscall.IPV6_MTU_DISCOVER    //= 0x17
	_ = syscall.IPV6_MULTICAST_HOPS  //= 0x12
	_ = syscall.IPV6_MULTICAST_IF    //= 0x11
	_ = syscall.IPV6_MULTICAST_LOOP  //= 0x13
	_ = syscall.IPV6_NEXTHOP         //= 0x9
	_ = syscall.IPV6_PKTINFO         //= 0x32
	_ = syscall.IPV6_PMTUDISC_DO     //= 0x2
	_ = syscall.IPV6_PMTUDISC_DONT   //= 0x0
	_ = syscall.IPV6_PMTUDISC_PROBE  //= 0x3
	_ = syscall.IPV6_PMTUDISC_WANT   //= 0x1
	_ = syscall.IPV6_RECVDSTOPTS     //= 0x3a
	_ = syscall.IPV6_RECVERR         //= 0x19
	_ = syscall.IPV6_RECVHOPLIMIT    //= 0x33
	_ = syscall.IPV6_RECVHOPOPTS     //= 0x35
	_ = syscall.IPV6_RECVPKTINFO     //= 0x31
	_ = syscall.IPV6_RECVRTHDR       //= 0x38
	_ = syscall.IPV6_RECVTCLASS      //= 0x42
	_ = syscall.IPV6_ROUTER_ALERT    //= 0x16
	_ = syscall.IPV6_RTHDR           //= 0x39
	_ = syscall.IPV6_RTHDRDSTOPTS    //= 0x37
	_ = syscall.IPV6_RTHDR_LOOSE     //= 0x0
	_ = syscall.IPV6_RTHDR_STRICT    //= 0x1
	_ = syscall.IPV6_RTHDR_TYPE_0    //= 0x0
	_ = syscall.IPV6_RXDSTOPTS       //= 0x3b
	_ = syscall.IPV6_RXHOPOPTS       //= 0x36
	_ = syscall.IPV6_TCLASS          //= 0x43
	_ = syscall.IPV6_UNICAST_HOPS    //= 0x10
	_ = syscall.IPV6_V6ONLY          //= 0x1a
	_ = syscall.IPV6_XFRM_POLICY     //= 0x23

	_ = syscall.IP_ADD_MEMBERSHIP         //= 0x23
	_ = syscall.IP_ADD_SOURCE_MEMBERSHIP  //= 0x27
	_ = syscall.IP_BLOCK_SOURCE           //= 0x26
	_ = syscall.IP_DEFAULT_MULTICAST_LOOP //= 0x1
	_ = syscall.IP_DEFAULT_MULTICAST_TTL  //= 0x1
	_ = syscall.IP_DF                     //= 0x4000
	_ = syscall.IP_DROP_MEMBERSHIP        //= 0x24
	_ = syscall.IP_DROP_SOURCE_MEMBERSHIP //= 0x28
	_ = syscall.IP_FREEBIND               //= 0xf
	_ = syscall.IP_HDRINCL                //= 0x3
	_ = syscall.IP_IPSEC_POLICY           //= 0x10
	_ = syscall.IP_MAXPACKET              //= 0xffff
	_ = syscall.IP_MAX_MEMBERSHIPS        //= 0x14
	_ = syscall.IP_MF                     //= 0x2000
	_ = syscall.IP_MINTTL                 //= 0x15
	_ = syscall.IP_MSFILTER               //= 0x29
	_ = syscall.IP_MSS                    //= 0x240
	_ = syscall.IP_MTU                    //= 0xe
	_ = syscall.IP_MTU_DISCOVER           //= 0xa
	_ = syscall.IP_MULTICAST_IF           //= 0x20
	_ = syscall.IP_MULTICAST_LOOP         //= 0x22
	_ = syscall.IP_MULTICAST_TTL          //= 0x21
	_ = syscall.IP_OFFMASK                //= 0x1fff
	_ = syscall.IP_OPTIONS                //= 0x4
	_ = syscall.IP_ORIGDSTADDR            //= 0x14
	_ = syscall.IP_PASSSEC                //= 0x12
	_ = syscall.IP_PKTINFO                //= 0x8
	_ = syscall.IP_PKTOPTIONS             //= 0x9
	_ = syscall.IP_PMTUDISC               //= 0xa
	_ = syscall.IP_PMTUDISC_DO            //= 0x2
	_ = syscall.IP_PMTUDISC_DONT          //= 0x0
	_ = syscall.IP_PMTUDISC_PROBE         //= 0x3
	_ = syscall.IP_PMTUDISC_WANT          //= 0x1
	_ = syscall.IP_RECVERR                //= 0xb
	_ = syscall.IP_RECVOPTS               //= 0x6
	_ = syscall.IP_RECVORIGDSTADDR        //= 0x14
	_ = syscall.IP_RECVRETOPTS            //= 0x7
	_ = syscall.IP_RECVTOS                //= 0xd
	_ = syscall.IP_RECVTTL                //= 0xc
	_ = syscall.IP_RETOPTS                //= 0x7
	_ = syscall.IP_RF                     //= 0x8000
	_ = syscall.IP_ROUTER_ALERT           //= 0x5
	_ = syscall.IP_TOS                    //= 0x1
	_ = syscall.IP_TRANSPARENT            //= 0x13
	_ = syscall.IP_TTL                    //= 0x2
	_ = syscall.IP_UNBLOCK_SOURCE         //= 0x25
	_ = syscall.IP_XFRM_POLICY            //= 0x11

	_ = syscall.SO_NO_CHECK  //= 0xb
	_ = syscall.SO_PEERNAME  //= 0x1c
	_ = syscall.SO_PRIORITY  //= 0xc
	_ = syscall.SO_TIMESTAMP //= 0x1d

	_ = syscall.TCP_CONGESTION       //= 0xd
	_ = syscall.TCP_DEFER_ACCEPT     //= 0x9
	_ = syscall.TCP_INFO             //= 0xb
	_ = syscall.TCP_KEEPCNT          //= 0x6
	_ = syscall.TCP_KEEPIDLE         //= 0x4
	_ = syscall.TCP_KEEPINTVL        //= 0x5
	_ = syscall.TCP_LINGER2          //= 0x8
	_ = syscall.TCP_MAXSEG           //= 0x2
	_ = syscall.TCP_MAXWIN           //= 0xffff
	_ = syscall.TCP_MAX_WINSHIFT     //= 0xe
	_ = syscall.TCP_MD5SIG           //= 0xe
	_ = syscall.TCP_MD5SIG_MAXKEYLEN //= 0x50
	_ = syscall.TCP_MSS              //= 0x200
	_ = syscall.TCP_NODELAY          //= 0x1
	_ = syscall.TCP_QUICKACK         //= 0xc
	_ = syscall.TCP_SYNCNT           //= 0x7
	_ = syscall.TCP_WINDOW_CLAMP     //= 0xa
)

// Socket options.
// To avoid package dependencies, values are copied from golang.org/x/sys/unix
// TODO: Remove unnecessary options. Make options configurable.
const (
	SO_BINDTOIFINDEX    = 0x3e
	SO_BUSY_POLL        = 0x2e
	SO_BUSY_POLL_BUDGET = 0x46
	SO_COOKIE           = 0x39
	SO_INCOMING_CPU     = 0x31
	SO_PREFER_BUSY_POLL = 0x45
	SO_RCVMARK          = 0x4b
	SO_RESERVE_MEM      = 0x49
	SO_REUSEPORT        = 0xf
	SO_TXREHASH         = 0x4a
	SO_TXTIME           = 0x3d
	SO_WIFI_STATUS      = 0x29
	SO_ZEROCOPY         = 0x3c

	IPPROTO_BEETPH   = 0x5e
	IPPROTO_ETHERNET = 0x8f
	IPPROTO_L2TP     = 0x73
	IPPROTO_MH       = 0x87
	IPPROTO_MPLS     = 0x89
	IPPROTO_MPTCP    = 0x106

	IPV6_ADDR_PREFERENCES     = 0x48
	IPV6_AUTOFLOWLABEL        = 0x46
	IPV6_DONTFRAG             = 0x3e
	IPV6_FLOW                 = 0x11
	IPV6_FREEBIND             = 0x4e
	IPV6_HDRINCL              = 0x24
	IPV6_MINHOPCOUNT          = 0x49
	IPV6_MULTICAST_ALL        = 0x1d
	IPV6_ORIGDSTADDR          = 0x4a
	IPV6_PATHMTU              = 0x3d
	IPV6_PMTUDISC_INTERFACE   = 0x4
	IPV6_PMTUDISC_OMIT        = 0x5
	IPV6_RECVERR_RFC4884      = 0x1f
	IPV6_RECVFRAGSIZE         = 0x4d
	IPV6_RECVORIGDSTADDR      = 0x4a
	IPV6_RECVPATHMTU          = 0x3c
	IPV6_ROUTER_ALERT_ISOLATE = 0x1e
	IPV6_TRANSPARENT          = 0x4b
	IPV6_UNICAST_IF           = 0x4c
	IPV6_USER_FLOW            = 0xe

	IP_BIND_ADDRESS_NO_PORT = 0x18
	IP_CHECKSUM             = 0x17
	IP_LOCAL_PORT_RANGE     = 0x33
	IP_MULTICAST_ALL        = 0x31
	IP_NODEFRAG             = 0x16
	IP_PMTUDISC_INTERFACE   = 0x4
	IP_PMTUDISC_OMIT        = 0x5
	IP_PROTOCOL             = 0x34
	IP_RECVERR_RFC4884      = 0x1a
	IP_RECVFRAGSIZE         = 0x19
	IP_UNICAST_IF           = 0x32
	IP_USER_FLOW            = 0xd

	TCPOPT_EOL            = 0x0
	TCPOPT_MAXSEG         = 0x2
	TCPOPT_NOP            = 0x1
	TCPOPT_SACK           = 0x5
	TCPOPT_SACK_PERMITTED = 0x4
	TCPOPT_TIMESTAMP      = 0x8
	TCPOPT_TSTAMP_HDR     = 0x101080a
	TCPOPT_WINDOW         = 0x3

	TCP_CC_INFO              = 0x1a
	TCP_CM_INQ               = 0x24
	TCP_COOKIE_IN_ALWAYS     = 0x1
	TCP_COOKIE_MAX           = 0x10
	TCP_COOKIE_MIN           = 0x8
	TCP_COOKIE_OUT_NEVER     = 0x2
	TCP_COOKIE_PAIR_SIZE     = 0x20
	TCP_COOKIE_TRANSACTIONS  = 0xf
	TCP_FASTOPEN             = 0x17
	TCP_FASTOPEN_CONNECT     = 0x1e
	TCP_FASTOPEN_KEY         = 0x21
	TCP_FASTOPEN_NO_COOKIE   = 0x22
	TCP_INQ                  = 0x24
	TCP_MD5SIG_EXT           = 0x20
	TCP_MD5SIG_FLAG_IFINDEX  = 0x2
	TCP_MD5SIG_FLAG_PREFIX   = 0x1
	TCP_MD5SIG_MAXKEYLEN     = 0x50
	TCP_MSS_DEFAULT          = 0x218
	TCP_MSS_DESIRED          = 0x4c4
	TCP_NOTSENT_LOWAT        = 0x19
	TCP_QUEUE_SEQ            = 0x15
	TCP_QUICKACK             = 0xc
	TCP_REPAIR               = 0x13
	TCP_REPAIR_OFF           = 0x0
	TCP_REPAIR_OFF_NO_WP     = -0x1
	TCP_REPAIR_ON            = 0x1
	TCP_REPAIR_OPTIONS       = 0x16
	TCP_REPAIR_QUEUE         = 0x14
	TCP_REPAIR_WINDOW        = 0x1d
	TCP_SAVED_SYN            = 0x1c
	TCP_SAVE_SYN             = 0x1b
	TCP_SYNCNT               = 0x7
	TCP_S_DATA_IN            = 0x4
	TCP_S_DATA_OUT           = 0x8
	TCP_THIN_DUPACK          = 0x11
	TCP_THIN_LINEAR_TIMEOUTS = 0x10
	TCP_TIMESTAMP            = 0x18
	TCP_TX_DELAY             = 0x25
	TCP_ULP                  = 0x1f
	TCP_USER_TIMEOUT         = 0x12
	TCP_V4_FLOW              = 0x1
	TCP_V6_FLOW              = 0x5
	TCP_ZEROCOPY_RECEIVE     = 0x23

	UDP_CORK                   = 0x1
	UDP_ENCAP                  = 0x64
	UDP_ENCAP_ESPINUDP         = 0x2
	UDP_ENCAP_ESPINUDP_NON_IKE = 0x1
	UDP_ENCAP_GTP0             = 0x4
	UDP_ENCAP_GTP1U            = 0x5
	UDP_ENCAP_L2TPINUDP        = 0x3
	UDP_GRO                    = 0x68
	UDP_SEGMENT                = 0x67
)

// setsockoptString is the re-mapped function of syscall.SetsockoptString for testing.
var setsockoptString = syscall.SetsockoptString

func (c *SockSOOption) Controllers() []Controller {
	var controllers []Controller
	controllers = appendIfNotNil(controllers, soBindToDevice(c.BindToDevice))
	controllers = appendIfNotNil(controllers, soDebug(c.Debug))
	controllers = appendIfNotNil(controllers, soIncomingCPU(c.IncomingCPU))
	controllers = appendIfNotNil(controllers, soKeepAlive(c.KeepAlive))
	controllers = appendIfNotNil(controllers, soLinger(c.Linger))
	controllers = appendIfNotNil(controllers, soMark(c.Mark))
	controllers = appendIfNotNil(controllers, soRcvbuf(c.ReceiveBuffer))
	controllers = appendIfNotNil(controllers, soRcvbufForce(c.ReceiveBufferForce))
	controllers = appendIfNotNil(controllers, soSndtimeo(c.SendTimeout))
	controllers = appendIfNotNil(controllers, soRcvtimeo(c.ReceiveTimeout))
	controllers = appendIfNotNil(controllers, soReuseaddr(c.ReuseAddr))
	controllers = appendIfNotNil(controllers, soReuseport(c.ReusePort))
	controllers = appendIfNotNil(controllers, soSndbuf(c.SendBuffer))
	controllers = appendIfNotNil(controllers, soSndbufForce(c.SendBufferForce))
	return controllers
}

func soBindToDevice(value string) Controller {
	if value == "" {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptString(int(fd), syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, value); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "SOL_SOCKET.SO_BINDTODEVICE",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func soDebug(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_DEBUG, 1); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "SOL_SOCKET.SO_DEBUG",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func soIncomingCPU(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, SO_INCOMING_CPU, 1); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "SOL_SOCKET.SO_INCOMING_CPU",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func soKeepAlive(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, 1); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "SOL_SOCKET.SO_KEEPALIVE",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func soLinger(value int32) Controller {
	if value <= 0 {
		return nil
	}
	l := &syscall.Linger{
		Onoff:  1,
		Linger: value,
	}
	return func(fd uintptr) error {
		if err := setsockoptLinger(int(fd), syscall.SOL_SOCKET, syscall.SO_LINGER, l); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "SOL_SOCKET.SO_LINGER",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func soMark(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, value); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "SOL_SOCKET.SO_MARK",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func soRcvbuf(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF, value); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "SOL_SOCKET.SO_RCVBUF",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func soRcvbufForce(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUFFORCE, value); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "SOL_SOCKET.SO_RCVBUFFORCE",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func soSndbuf(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_SNDBUF, value); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "SOL_SOCKET.SO_SNDBUF",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func soSndbufForce(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_SNDBUFFORCE, value); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "SOL_SOCKET.SO_SNDBUFFORCE",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func soSndtimeo(value float64) Controller {
	if value <= 0 {
		return nil
	}
	tv := syscall.NsecToTimeval(int64(1_000_000_000 * value))
	return func(fd uintptr) error {
		if err := setsockoptTimeval(int(fd), syscall.SOL_SOCKET, syscall.SO_SNDTIMEO, &tv); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "SOL_SOCKET.SO_SNDTIMEO",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func soRcvtimeo(value float64) Controller {
	if value <= 0 {
		return nil
	}
	tv := syscall.NsecToTimeval(int64(1_000_000_000 * value))
	return func(fd uintptr) error {
		if err := setsockoptTimeval(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "SOL_SOCKET.SO_RCVTIMEO",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func soReuseaddr(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "SOL_SOCKET.SO_REUSEADDR",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func soReuseport(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, SO_REUSEPORT, 1); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "SOL_SOCKET.SO_REUSEPORT",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func (c *SockIPOption) Controllers() []Controller {
	var controllers []Controller
	controllers = appendIfNotNil(controllers, ipBindAddressNoPort(c.BindAddressNoPort))
	controllers = appendIfNotNil(controllers, ipFreeBind(c.FreeBind))
	controllers = appendIfNotNil(controllers, ipLocalPortRange(c.LocalPortRangeUpper, c.LocalPortRangeLower))
	controllers = appendIfNotNil(controllers, ipTransparent(c.Transparent))
	controllers = appendIfNotNil(controllers, ipTTL(c.TTL))
	return controllers
}

func ipBindAddressNoPort(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_IP, IP_BIND_ADDRESS_NO_PORT, 1); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_IP.IP_BIND_ADDRESS_NO_PORT",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func ipFreeBind(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_FREEBIND, 1); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_IP.IP_FREEBIND",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func ipLocalPortRange(upper, lower int) Controller {
	upper = cmp.Or(upper, 0)
	lower = cmp.Or(lower, 0)
	if upper <= 0 && lower <= 0 {
		return nil
	}
	b := make([]byte, 0, 4)
	b = binary.NativeEndian.AppendUint16(b, uint16(upper)) //nolint:gosec // G115: integer overflow conversion int -> uint16
	b = binary.NativeEndian.AppendUint16(b, uint16(lower)) //nolint:gosec // G115: integer overflow conversion int -> uint16
	pr := binary.BigEndian.Uint32(b)
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_IP, IP_LOCAL_PORT_RANGE, int(pr)); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_IP.IP_LOCAL_PORT_RANGE",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func ipTransparent(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_TRANSPARENT, 1); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_IP.IP_TRANSPARENT",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func ipTTL(ttl int) Controller {
	if ttl <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_TTL, ttl); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_IP.IP_TTL",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func (c *SockIPV6Option) Controllers() []Controller {
	return nil
}

func (c *SockTCPOption) Controllers() []Controller {
	var controllers []Controller
	controllers = appendIfNotNil(controllers, tcpCORK(c.CORK))
	controllers = appendIfNotNil(controllers, tcpDeferAccept(c.DeferAccept))
	controllers = appendIfNotNil(controllers, tcpKeepCount(c.KeepCount))
	controllers = appendIfNotNil(controllers, tcpKeepIdle(c.KeepIdle))
	controllers = appendIfNotNil(controllers, tcpKeepInterval(c.KeepInterval))
	controllers = appendIfNotNil(controllers, tcpLinger2(c.Linger2))
	controllers = appendIfNotNil(controllers, tcpMaxSegment(c.MaxSegment))
	controllers = appendIfNotNil(controllers, tcpNoDelay(c.NoDelay))
	controllers = appendIfNotNil(controllers, tcpQuickAck(c.QuickAck))
	controllers = appendIfNotNil(controllers, tcpSynCount(c.SynCount))
	controllers = appendIfNotNil(controllers, tcpUserTimeout(c.UserTimeout))
	controllers = appendIfNotNil(controllers, tcpWindowClamp(c.WindowClamp))
	controllers = appendIfNotNil(controllers, tcpFastOpen(c.FastOpen))
	controllers = appendIfNotNil(controllers, tcpFastOpenConnect(c.FastOpenConnect))
	return controllers
}

func tcpCORK(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_CORK, 1); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_TCP.TCP_CORK",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func tcpDeferAccept(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_DEFER_ACCEPT, value); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_TCP.TCP_DEFER_ACCEPT",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func tcpKeepCount(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_KEEPCNT, value); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_TCP.TCP_KEEPCNT",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func tcpKeepIdle(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_KEEPIDLE, value); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_TCP.TCP_KEEPIDLE",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}
func tcpKeepInterval(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_KEEPINTVL, value); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_TCP.TCP_KEEPINTVL",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func tcpLinger2(value int32) Controller {
	if value <= 0 {
		return nil
	}
	l := &syscall.Linger{
		Onoff:  1,
		Linger: value,
	}
	return func(fd uintptr) error {
		if err := setsockoptLinger(int(fd), syscall.IPPROTO_TCP, syscall.TCP_LINGER2, l); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_TCP.TCP_LINGER2",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func tcpMaxSegment(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_MAXSEG, value); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_TCP.TCP_MAXSEG",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func tcpNoDelay(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_TCP.TCP_NODELAY",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func tcpQuickAck(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, TCP_QUICKACK, 1); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_TCP.TCP_QUICKACK",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func tcpSynCount(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, TCP_SYNCNT, value); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_TCP.TCP_SYNCNT",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func tcpUserTimeout(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, TCP_USER_TIMEOUT, value); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_TCP.TCP_USER_TIMEOUT",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func tcpWindowClamp(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_WINDOW_CLAMP, value); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_TCP.TCP_WINDOW_CLAMP",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func tcpFastOpen(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, TCP_FASTOPEN, 1); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_TCP.TCP_FASTOPEN",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func tcpFastOpenConnect(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, TCP_FASTOPEN_CONNECT, 1); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_TCP.TCP_FASTOPEN_CONNECT",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func (c *SockUDPOption) Controllers() []Controller {
	var controllers []Controller
	controllers = appendIfNotNil(controllers, udpCORK(c.CORK))
	controllers = appendIfNotNil(controllers, udpSegment(c.Segment))
	controllers = appendIfNotNil(controllers, udpGRO(c.GRO))
	return controllers
}

func udpCORK(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_UDP, UDP_CORK, 1); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_UDP.UDP_CORK",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func udpSegment(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_UDP, UDP_SEGMENT, value); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_UDP.UDP_SEGMENT",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}

func udpGRO(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_UDP, UDP_GRO, 1); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_UDP.UDP_GRO",
			}).Wrap(err)
		}
		runtime.KeepAlive(fd)
		return nil
	}
}
