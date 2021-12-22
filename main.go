package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var domainMode, showCloudflare bool
	var sc *bufio.Scanner

	concurrency := 20

	flag.IntVar(&concurrency, "c", concurrency, "Set the concurrency level")
	flag.BoolVar(&domainMode, "d", false, "Print domains instead of IP addresses")
	flag.BoolVar(&showCloudflare, "cf", false, "Show CloudFlare only")
	flag.Parse()

	jobs := make(chan string)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		go func() {
			for host := range jobs {
				addr, err := net.LookupIP(strings.TrimSpace(host))
				if err != nil {
					continue
				}

				ip := addr[0]
				cf := isCloudflare(ip)

				if !cf && !showCloudflare {
					show(host, ip, domainMode)
				} else if cf && showCloudflare {
					show(host, ip, domainMode)
				}
			}

			defer wg.Done()
		}()
	}

	fn := flag.Arg(0)

	if isStdin() {
		sc = bufio.NewScanner(os.Stdin)
	} else if fn != "" {
		r, err := os.Open(fn)
		if err == nil {
			sc = bufio.NewScanner(r)
		}
	}

	for sc.Scan() {
		i := sc.Text()
		u, err := url.Parse(i)
		if err == nil {
			if u.Host != "" {
				jobs <- u.Host
			} else {
				jobs <- i
			}
		}
	}

	close(jobs)
	wg.Wait()
}

func show(host string, ip net.IP, mode bool) {
	if mode {
		fmt.Println(host)
	} else {
		fmt.Println(ip)
	}
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
	resp, err := http.Get("https://www.cloudflare.com/ips-v4")
	if err != nil {
		log.Fatalln(err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	cidrs := strings.Fields(string(b))

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

func isStdin() bool {
	f, e := os.Stdin.Stat()
	if e != nil {
		return false
	}

	if f.Mode()&os.ModeNamedPipe == 0 {
		return false
	}

	return true
}
