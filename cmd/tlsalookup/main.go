// tlsalookup - takes a hostname and looks up the TLSA record for port 443 tcp
// USAGE: tlsalookup 0x41414141.de
//
package main

import (
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"os"
)

func init() {
}

func main() {
	flag.Parse()
	upstreamServers := []string{"8.8.8.8:53", "8.8.4.4:53"}

	query := os.Args[1]
	//typ := os.Args[2]

	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{"_443._tcp." + query + ".", dns.TypeTLSA, dns.ClassINET}

	c := new(dns.Client)
	in, rtt, err := c.Exchange(m1, upstreamServers[0])
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", in)
	fmt.Printf("RTT: %s\n", rtt)
}
