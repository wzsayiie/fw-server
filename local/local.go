package local

import (
	"fmt"
	"net"
	"strings"
)

// get ipv4 address and ipv6 link-local address of current host.
func HostIPs(ip4 *string, ip6 *string) error {
	if ip4 == nil && ip6 == nil {
		return fmt.Errorf("parameters are both nil")
	}

	var addrs = favAddrs()
	if len(addrs) == 0 {
		return fmt.Errorf("no available addresses")
	}

	var v4 string
	var v6 string
	for _, v := range addrs {
		var host = pickHost(v)

		if isIPv4(host) {
			v4 = host
		} else if isIPv6(host) {
			v6 = host
		}
	}

	if ip4 != nil {
		*ip4 = v4
	}
	if ip6 != nil {
		*ip6 = v6
	}
	return nil
}

func favAddrs() []net.Addr {

	var intfs, intfErr = net.Interfaces()
	if intfErr != nil {
		return nil
	}

	for _, v := range intfs {
		if v.Flags&net.FlagUp == 0 {
			// isn't working.
			continue
		}
		if v.Flags&net.FlagLoopback != 0 {
			// is local loopback address.
			continue
		}

		var addrs, addrErr = v.Addrs()
		if addrErr != nil {
			continue
		}

		// there are usually two valus,
		// ipv4 address and ipv6 link-local addresses.
		if len(addrs) != 0 {
			return addrs
		}
	}
	return nil
}

func pickHost(addr net.Addr) string {
	var str = addr.String()

	for i, v := range str {
		if v == '/' {
			return str[:i]
		}
	}
	return ""
}

func isIPv4(ip string) bool {
	return strings.Count(ip, ".") == 3
}

func isIPv6(ip string) bool {
	return strings.Count(ip, ":") >= 2
}
