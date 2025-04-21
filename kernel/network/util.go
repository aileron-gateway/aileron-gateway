// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package network

import (
	"math"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
)

// VerboseLogs is the flag to enable
// networking debug logs.
var VerboseLogs bool

func init() {
	e := os.Getenv("GODEBUG")
	VerboseLogs = strings.Contains(e, "network=1")
}

// Container checks if the IP and Port
// is contained in the implementers
// network address.
type Container interface {
	// Contains returns if the given ip and port
	// is contained the networks that implementers defined.
	// This method must not panic even nil or invalid
	// ip or port was given.
	Contains(ip net.IP, port int) bool
}

// netContainer holds IP and Port range.
// This implements Container interface.
type netContainer struct {
	IPNet    *net.IPNet
	PortFrom int
	PortTo   int
}

// Contains returns is the given ip and port
// are in the range of this container.
func (c *netContainer) Contains(ip net.IP, port int) bool {
	if port < c.PortFrom || port > c.PortTo {
		return false
	}
	return c.IPNet.Contains(ip)
}

// netContainers returns containers from given addresses.
// Example of supported formats are as follows.
//   - CIDR: "192.168.0.0/16", "fc00::/7"
//   - CIDR with port: "192.168.0.0/16:8080", "[fc00::/7]:8080"
//   - CIDR with port range: "192.168.0.0/16:8080-9090", "[fc00::/7]:8080-9090"
//
// Following formats are NOT supported.
//   - Port only: ":8080"
//   - Port range only: ":8080-9090"
//
// Use "0.0.0.0/0" or "::/0" to allow all IPs.
// Note that IPv6 expression does not work for IPv4 mapped addresses.
// All ports are allowed if port or port range is not specified.
func netContainers(addresses []string) ([]Container, error) {
	containers := make([]Container, 0, len(addresses))
	for _, addr := range addresses {
		container := &netContainer{
			PortFrom: 0,
			PortTo:   math.MaxInt,
		}

		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			host = addr // Port may not be in the addr.
		}

		_, nw, err := net.ParseCIDR(host)
		if err != nil {
			return nil, (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeContainer,
				Description: ErrDscContainer,
			}).Wrap(err)
		}
		container.IPNet = nw

		if port != "" {
			if !strings.Contains(port, "-") {
				port = port + "-" + port // Change port to port range. For example, "8080" to "8080-8080".
			}
			arr := strings.Split(port, "-")
			if container.PortFrom, err = strconv.Atoi(arr[0]); err != nil {
				return nil, (&er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeContainer,
					Description: ErrDscContainer,
				}).Wrap(err)
			}
			if container.PortTo, err = strconv.Atoi(arr[1]); err != nil {
				return nil, (&er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeContainer,
					Description: ErrDscContainer,
				}).Wrap(err)
			}
		}

		containers = append(containers, container)
	}
	return slices.Clip(containers), nil
}

// splitHostPort split the host and
// This function returns nil as IP when an invalid address was given.
// This function uses net.SplitHostPort to parse address with ip and port,
// net.ParseIP to parse ip only address.
// Supported formats are the same with net.SplitHostPort and net.ParseIP,
// This function does not return any error but return nil IP and port 0 instead.
func splitHostPort(address string) (net.IP, int) {
	if address == "" {
		return nil, 0
	}
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		host = address
	}
	if port != "" {
		portNum, _ := strconv.Atoi(port)
		return net.ParseIP(host), portNum
	}
	return net.ParseIP(host), 0
}
