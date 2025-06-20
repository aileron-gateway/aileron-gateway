// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package network_test

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"github.com/aileron-gateway/aileron-gateway/kernel/network"
)

// ServerCert is a self-signed server certification.
// This cert was generated by the command below and checked
// by the command `openssl x509 -text -in server.crt -noout`.
// This cert include "localhost", "*.sock"(unix domain socket), "127.0.0.0" as allowed hosts.
//   - openssl req -newkey rsa:4096 -nodes -keyout server.key -x509 -days 36500 -out server.crt -addext 'subjectAltName = DNS:localhost,DNS:*.sock,IP:127.0.0.1' -subj '/CN=127.0.0.1'
var serverCert = `
-----BEGIN CERTIFICATE-----
MIIFLzCCAxegAwIBAgIUS824ePmTCuw0pAGZBvd4u971LR4wDQYJKoZIhvcNAQEL
BQAwFDESMBAGA1UEAwwJMTI3LjAuMC4xMCAXDTI0MDIyNjEwMzQyN1oYDzIxMjQw
MjAyMTAzNDI3WjAUMRIwEAYDVQQDDAkxMjcuMC4wLjEwggIiMA0GCSqGSIb3DQEB
AQUAA4ICDwAwggIKAoICAQCp3YPwQW1u0vrHezhSzMwd+W543E2p91vlgtRQw3DY
Gp2NPjsYekdOAjzMe7QiQskI4zXTWwLQLlETvMXdUQZCkLZd+j0ORc2lZO7byfBO
t+Ct9KUBtrCe3P0ZuL3zatJLDrZDa9K5ySj8FbUpFBurWM0Nd180ETPfniUgzli1
Lra40inqAf3NNMu/Q+BaWJ4Dq2pgeHNLcj9k7XdUC9rHO5qkwCtVZaI4R36d1k/W
+Ow9QtfEKhBrR0b8GjdK+U7vgs8T6eFzT3+ULdPnw9UeUPIOvErDS2/TCmUtfMEi
4wA5G6I9VWr+1XPvg3uLxj+InzPmqWyXxS4Z4GlyNmnWDT0STeng6RmmBjHfF0ao
Z3DsO7t4E2XVvfZnOxUfdrbW+GYYQkCNEvcLQvLjVfMtIynY1PV8InKyAfTP3VdC
6UqV7SqgeCW3m0turbUBF62iTcpY8fkpISONsSHO5yN9JRXbTLc2mRvwRdO6t0oj
L8IVGLVhbiQGbYjjab22yMJXNnd5tMuB9aHrccuA/Z1Rd0hrDUJ2lfjfohDGla4p
jkhQFXZlKjqZULJv3Q02fiY0u5NReFynWdTy+jktTCkqIm/WOgphFXP2bEFo0jo+
tjEIqUZwuSjBaM7rIQiYqZVPc7nCB7u4Ci7eJT+jwmC+3OdURLUFPMSkSJPlGvKj
kwIDAQABo3cwdTAdBgNVHQ4EFgQUjEXgiVrZAnZCKQg1hICw29D8n54wHwYDVR0j
BBgwFoAUjEXgiVrZAnZCKQg1hICw29D8n54wDwYDVR0TAQH/BAUwAwEB/zAiBgNV
HREEGzAZgglsb2NhbGhvc3SCBiouc29ja4cEfwAAATANBgkqhkiG9w0BAQsFAAOC
AgEAA0C3U/yJgXNASugOzxlTowjtHmIrRgBVUKT0SLgbYXlVATjaqm7fEI2ZASbf
y2+YWHi9hZDZagjSic11rVJ8ARIc33AGBr1yr46LmdQpEdMndzTWo9vct6i+LfWo
DC+2t5lTLngXbvW9zwZeKgbJNwb/18qv+hP1cKLAWdv6yZrpJ4cK5dAFgRoDvtjo
B7ZOkDhHbVQedHqHCgtkxRJe5s6Vt9kRZ3gGBoZmk0ZZ3Q0jlCvzdo8gSjb67Gye
aa5epYI3wdyRLlodpPdb1j3bBR38kzZ46i+BRF5EQhAdQwcFpqonWwU+9e5sGy0S
RyD61ab3Au8GVGwe9ov/CnEtxfdAPQkRaQCvMdW5lqEc5oI4cJbjdoCW4O7ZI0DV
i/jaOjT3DoMOEP4KMMdMid7zM1ZRtT6ySXiMGy2NDwyqVB5qkkk4ZbsB6Rf5RkuT
KlRNI3jBBNACBB1XDKyrmFyJp8LxPi/DuEO7h8Yc1jKy9XxJUCXPsJm9pRKyiZjJ
07E5PtYEKBLw0ueULuWnq3qF3AqADwAqnlVjCaDfSAWsD7JwKtwYn2BdCAKRGkRK
Y6mUJt/NfBChb1UOXmj6kdeTlVkyALbH7RvgxFfEtCY/BoUzNfBVVFeEI1jwczvo
/dZadUHIBqea8MjGQnXCKMgtO9O7GqkOqJPXVyAP6LTmGb0=
-----END CERTIFICATE-----
`

// serverKey is the signing key for the serverCert.
// This key was generated by the command below.
//   - openssl req -newkey rsa:4096 -nodes -keyout server.key -x509 -days 36500 -out server.crt -addext 'subjectAltName = DNS:localhost,DNS:*.sock,IP:127.0.0.1' -subj '/CN=127.0.0.1'
var serverKey = `
-----BEGIN PRIVATE KEY-----
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
o6Sye62iEmy7VFhLWsDQe6PpvMJx
-----END PRIVATE KEY-----
`

// x509KeyPair is the X509 cert and key pair.
// This will be used in the example tests.
var x509KeyPair, _ = tls.X509KeyPair([]byte(serverCert), []byte(serverKey))

// rootCAs returns root CAs corresponding to the
// x509KeyPair created above.
var rootCAs = func() *x509.CertPool {
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(serverCert))
	return pool
}

func ExampleNewListener_listenTCP() {
	available, _ := net.Listen("tcp", "127.0.0.1:0")
	available.Close()
	testAddr := available.Addr().String()

	// Listen TCP at the local address.
	lc := &network.ListenConfig{
		Address: "tcp://" + testAddr,
	}
	ln, _ := network.NewListener(lc)
	defer ln.Close()

	go func() {
		dc := &network.DialConfig{}
		dl, _ := network.NewDialer(dc)
		conn, _ := dl.Dial("tcp", testAddr)
		conn.Write([]byte("hello!!"))
		conn.Close()
	}()

	// Accept the connection created in the goroutine above.
	conn, _ := ln.Accept()
	defer conn.Close()

	b := make([]byte, 7) // Expect "hello!!".
	conn.Read(b)         // Receive "hello!!".

	fmt.Println(string(b))
	// Output:
	// hello!!
}

func ExampleNewListener_listenTLS() {
	available, _ := net.Listen("tcp", "127.0.0.1:0")
	available.Close()
	testAddr := available.Addr().String()

	lc := &network.ListenConfig{
		Address: "tcp://" + testAddr,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{x509KeyPair},
		},
	}
	ln, _ := network.NewListener(lc)
	defer ln.Close()

	go func() {
		dc := &network.DialConfig{
			TLSConfig: &tls.Config{
				RootCAs: rootCAs(),
			},
		}
		dl, _ := network.NewDialer(dc)
		conn, _ := dl.Dial("tcp", testAddr)
		conn.Write([]byte("hello!!"))
		conn.Close()
	}()
	// Accept the connection created in the goroutine above.
	conn, _ := ln.Accept()
	defer conn.Close()
	b := make([]byte, 7) // Expect "hello!!".
	conn.Read(b)         // Receive "hello!!".

	fmt.Println(string(b))
	// Output:
	// hello!!
}

func ZExampleNewListener_listenTCP_PathNameSocket() {
	// PathNameSocket is like "/var/run/test.sock"
	// AbstractSocket is like "@test"
	lc := &network.ListenConfig{
		Address: "unix://test.sock",
	}
	ln, _ := network.NewListener(lc)
	defer func() {
		ln.Close()
		os.Remove("test.sock") // Make sure to remove socket file.
	}()

	go func() {
		dc := &network.DialConfig{}
		dl, _ := network.NewDialer(dc)
		conn, _ := dl.Dial("unix", "test.sock")
		conn.Write([]byte("hello!!"))
		conn.Close()
	}()

	// Accept the connection created in the goroutine above.
	conn, _ := ln.Accept()
	defer conn.Close()

	b := make([]byte, 7) // Expect "hello!!".
	conn.Read(b)         // Receive "hello!!".

	fmt.Println(string(b))
	// Output:
	// hello!!
}

func ExampleNewPacketConn_connectUDP() {
	// Get a local address that is available for the test.
	available, _ := net.ListenPacket("udp", "127.0.0.1:0")
	available.Close()
	testAddr := available.LocalAddr().String()

	pc := &network.PacketConnConfig{
		Network: "udp",
		Address: testAddr,
	}
	conn, _ := network.NewPacketConn(pc)
	defer conn.Close()

	go func() {
		// Create a new UDP dialer..
		dc := &network.DialConfig{}
		dl, _ := network.NewDialer(dc)
		conn, _ := dl.Dial("udp", testAddr)
		conn.Write([]byte("hello!!"))
		conn.Close()
	}()

	b := make([]byte, 7) // Expect "hello!!".
	conn.ReadFrom(b)     // Read message from the connection.

	fmt.Println(string(b))
	// Output:
	// hello!!
}
