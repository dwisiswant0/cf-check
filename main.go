package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	concurrency := 20
	jobs := make(chan string)
	var wg sync.WaitGroup
	var domainMode bool

	flag.IntVar(&concurrency, "c", 20, "Set the concurrency level")
	flag.BoolVar(&domainMode, "d", false, "Prints domain instead of IP address")
	flag.Parse()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			for host := range jobs {
				addr, err := net.LookupIP(strings.TrimSpace(host))
				if err != nil {
					continue
				}

				if !isCloudflare(addr[0]) {
					if domainMode {
						fmt.Println(host)
					} else {
						fmt.Println(addr[0])
					}
				}
			}
			wg.Done()
		}()
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		jobs <- sc.Text()
	}

	close(jobs)

	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %s\n", err)
	}

	wg.Wait()
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	lenIPs := len(ips)
	switch {
	case lenIPs < 2:
		return ips, nil

	default:
		return ips[1 : len(ips)-1], nil
	}
}

func isCloudflare(ip net.IP) bool {
	cidrs := []string{"173.245.48.0/20", "103.21.244.0/22", "103.22.200.0/22", "103.31.4.0/22", "141.101.64.0/18", "108.162.192.0/18", "190.93.240.0/20", "188.114.96.0/20", "197.234.240.0/22", "198.41.128.0/17", "162.158.0.0/15", "104.16.0.0/12", "172.64.0.0/13", "131.0.72.0/22"}
	for i := range cidrs {
		hosts, err := hosts(cidrs[i])
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
