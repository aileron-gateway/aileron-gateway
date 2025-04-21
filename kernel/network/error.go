// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package network

const (
	ErrPkg = "network"

	ErrTypeSocket    = "socket"
	ErrTypePackConn  = "packet conn"
	ErrTypeDialer    = "dialer"
	ErrTypeListener  = "listener"
	ErrTypeTLSCert   = "tls certification"
	ErrTypeTLS       = "tls"
	ErrTypeTransport = "transport"
	ErrTypeContainer = "container"

	// ErrDscNil is a error description.
	// This description indicates failure of
	// applying socket option.
	ErrDscSockOpts = "failed to set socket option."

	// ErrDscPackConn is a error description.
	// This description indicates the failure of
	// packet connection.
	ErrDscPackConn = "packet conn error."

	// ErrDscDialer is a error description.
	// This description indicates the failure of
	// dialer. Crete or dial, etc.
	ErrDscDialer = "dialer error."

	// ErrDscListener is a error description.
	// This description indicates the failure of
	// listener. Crete or accept, etc.
	ErrDscListener = "listener error."

	// ErrDscTLSCert is a error description.
	// This description indicates the failure of
	// TLS certifications.
	ErrDscTLSCert = "TLS certification error."

	// ErrDscTLS is a error description.
	// This description indicates the failure of
	// TLS configuration.
	ErrDscTLS = "TLS configuration failed."

	// ErrDscNewTransport is a error description.
	// This description indicates failure of creation
	// of transport layer component.
	ErrDscNewTransport = "failed to create a new transport."

	// ErrDscContainer is a error description.
	// This description indicates failure of network address
	// container creation.
	ErrDscContainer = "network address container creation failed."
)
