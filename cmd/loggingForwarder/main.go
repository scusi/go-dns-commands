// a logging dns fowarder
//
// starts a dns server and forwards all queries to an upstream dns server.
// answers from the upstream server are replyed back to the client.
//
package main

import (
	"flag"
	"github.com/miekg/dns"
	"log"
	//"strconv"
)

var upstreamDNS string
var listenAddr string

func init() {
	flag.StringVar(&upstreamDNS, "upstream", "8.8.8.8:53", "upstream DNS to forward requets to")
	flag.StringVar(&listenAddr, "addr", "127.0.0.1:5301", "port the forwarder should listen for dns requests")
}

// handleDNSRequest - sends queries to upstream servers and response back to client.
func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	// generate a reply msg
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	// if it is a query,
	// forward to upstream dns and write answer back to client
	case dns.OpcodeQuery:
		c := new(dns.Client)
		in, _, err := c.Exchange(r, upstreamDNS)
		if err != nil {
			log.Printf("%s", err.Error())
		}
		// append answers to the reply
		for _, rr := range in.Answer {
			m.Answer = append(m.Answer, rr)
		}
	}
	// write reply to client
	w.WriteMsg(m)
	// log request and reply
	log.Printf("%s %s %s\n", w.RemoteAddr(), m.Question[0].Name, m.Answer)
}

func main() {
	flag.Parse()
	// attach request handler func
	dns.HandleFunc(".", handleDnsRequest)

	// start server
	//port := 5301
	server := &dns.Server{Addr: listenAddr, Net: "udp"}
	log.Printf("Starting at %d\n", listenAddr)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())

	}

}
