// a dump tool that resolves dns A records
// USAGE: lookup github.com
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
	upstreamServers := []string{"127.0.0.1:5300"}

	query := os.Args[1]
	//typ := os.Args[2]

	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{query + ".", dns.TypeA, dns.ClassINET}

	c := new(dns.Client)
	in, rtt, err := c.Exchange(m1, upstreamServers[0])
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", in)
	fmt.Printf("RTT: %s\n", rtt)
}
