// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build windows

package network

import (
	"runtime"
	"syscall"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
)

const (
	SO_DEBUG = 0x1
)

func (c *SockSOOption) Controllers() []Controller {
	var controllers []Controller
	controllers = appendIfNotNil(controllers, soDebug(c.Debug))
	controllers = appendIfNotNil(controllers, soKeepAlive(c.KeepAlive))
	controllers = appendIfNotNil(controllers, soLinger(c.Linger))
	controllers = appendIfNotNil(controllers, soRcvbuf(c.ReceiveBuffer))
	controllers = appendIfNotNil(controllers, soReuseaddr(c.ReuseAddr))
	controllers = appendIfNotNil(controllers, soSndbuf(c.SendBuffer))
	return controllers
}

func soDebug(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, SO_DEBUG, 1); err != nil {
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

func soKeepAlive(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, 1); err != nil {
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
		Linger: int32(value),
	}
	return func(fd uintptr) error {
		if err := setsockoptLinger(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_LINGER, l); err != nil {
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

func soRcvbuf(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF, value); err != nil {
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

func soSndbuf(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_SNDBUF, value); err != nil {
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

func soReuseaddr(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
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

func (c *SockIPOption) Controllers() []Controller {
	var controllers []Controller
	controllers = appendIfNotNil(controllers, ipTTL(c.TTL))
	return controllers
}

func ipTTL(ttl int) Controller {
	if ttl <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(syscall.Handle(fd), syscall.IPPROTO_IP, syscall.IP_TTL, ttl); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeSocket,
				Description: ErrDscSockOpts,
				Detail:      "IPPROTO_TCP.IP_TTL",
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
	controllers = appendIfNotNil(controllers, tcpNoDelay(c.NoDelay))
	return controllers
}

func tcpNoDelay(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(syscall.Handle(fd), syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1); err != nil {
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

func (c *SockUDPOption) Controllers() []Controller {
	return nil
}
