package main

import "net"

func isCloudflare(ip net.IP) bool {
	for _, c := range cidrs {
		if c == "" {
			continue
		}

		hosts, err := hosts(c)
		if err != nil {
			continue
		}

		for _, host := range hosts {
			if host == ip.String() {
				return true
			}
		}
	}
	return false
}
