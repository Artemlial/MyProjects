package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"text/tabwriter"
)

type result struct {
	IP       []string
	Hostname string
}

func lookupA(fqdn string) ([]string, error) {
	var ips []string
	in, err := net.LookupIP(fqdn)
	if err != nil {
		return ips, err
	}
	if len(in) < 1 {
		return ips, errors.New("no answer")
	}
	for _, ip := range in {
		ips = append(ips, ip.String())
	}
	return ips, nil
}

func lookupCNAME(fqdn string) (string, error) {
	in, err := net.LookupCNAME(fqdn)
	if err != nil {
		return "", err
	}
	if in == "" {
		return in, errors.New("no CNAME")
	}
	return in, nil
}

func lookup(fqdn string) []result {
	var results []result
	var cfqdn = fqdn
	for {
		cname, err := lookupCNAME(cfqdn)
		if err == nil {
			cfqdn = cname
			continue
		}
		ips, err := lookupA(cfqdn)
		if err != nil {
			break
		}
		results = append(results, result{
			IP:       ips,
			Hostname: cfqdn,
		})
		break
	}
	return results
}

type empty struct{}

func worker(tracker chan empty, fqdns chan string, gather chan []result) {
	for fqdn := range fqdns {
		results := lookup(fqdn)
		if len(results) > 0 {
			gather <- results
		}
	}
	var e empty
	tracker <- e
}

func main() {
	var (
		flDomain      = flag.String("domain", "", "The domain to perform guessing against.")
		flWordlist    = flag.String("wordlist", "", "The wordlist to use for guessing.")
		flWorkerCount = flag.Int("c", 100, "The amount of workers to use.")
		// flServerAddr = flag.String("server", "8.8.8.8:53", "The DNS server to use.")
	)
	flag.Parse()
	if *flDomain == "" || *flWordlist == "" {
		fmt.Println("-domain and -wordlist are required")
		os.Exit(1)
	}
	fmt.Println(*flWorkerCount)
	var results []result
	fqdns := make(chan string, *flWorkerCount)
	gather := make(chan []result)
	tracker := make(chan empty)
	fh, err := os.Open(*flWordlist)
	if err != nil {
		panic(err)
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)

	// r := &net.Resolver{
	//     PreferGo: true,
	//     Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
	//         d := net.Dialer{
	//             Timeout: time.Millisecond * time.Duration(10000),
	//         }
	//         return d.DialContext(ctx, network, *flServerAddr)
	//     },
	// }
	for i := 0; i < *flWorkerCount; i++ {
		go worker(tracker, fqdns, gather)
	}
	for scanner.Scan() {
		fqdns <- fmt.Sprintf("%s.%s", scanner.Text(), *flDomain)
	}

	go func() {
		for i := range gather {
			results = append(results, i...)
		}
		var e empty
		tracker <- e
	}()

	close(fqdns)
	for i := 0; i < *flWorkerCount; i++ {
		<-tracker
	}
	close(gather)
	<-tracker

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', 0)

	for _, res := range results {
		fmt.Fprintf(w, "%s{\n\t%v\n}\n", res.Hostname, res.IP)
	}
	w.Flush()
}
