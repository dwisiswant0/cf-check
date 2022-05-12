package main

import (
	"bufio"
	"flag"
	"net"
	"net/url"
	"os"
	"strings"
)

func main() {
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
		if err != nil {
			panic(err)
		}

		sc = bufio.NewScanner(r)
	} else {
		return
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
