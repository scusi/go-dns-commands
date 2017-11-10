package main

import (
	"github.com/miekg/dns"
	"log"
	"strconv"
)

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
		in, _, err := c.Exchange(r, "127.0.0.1:53")
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
	// attach request handler func
	dns.HandleFunc(".", handleDnsRequest)

	// start server
	port := 5301
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())

	}

}
