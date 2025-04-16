// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build !linux && !windows

package network

func (c *SockSOOption) Controllers() []Controller {
	return nil
}

func (c *SockIPOption) Controllers() []Controller {
	return nil
}

func (c *SockIPV6Option) Controllers() []Controller {
	return nil
}

func (c *SockTCPOption) Controllers() []Controller {
	return nil
}

func (c *SockUDPOption) Controllers() []Controller {
	return nil
}
