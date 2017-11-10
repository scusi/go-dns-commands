// logging dns server that can overwrite arbitrary A records
package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/miekg/dns"
)

// A records listed here will be served with priority over real values (if any)
// You can use the table to create your own tld or you can use it
// to overwrite real values.
var records = map[string]string{
	"test.scusi.":         "192.168.188.2",
	"analyzr.scusi.":      "192.168.188.46",
	"rechenknecht.scusi.": "192.168.188.45",
	"ida32.scusi.":        "192.168.188.37",
	"ida64.scusi.":        "192.168.188.55",
	"www.spiegel.de.":     "127.0.0.1",
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			ip := records[q.Name]
			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)

				}

			} else {
				m1 := new(dns.Msg)
				m1.Id = dns.Id()
				m1.RecursionDesired = true
				m1.Question = make([]dns.Question, 1)
				m1.Question[0] = dns.Question{q.Name, dns.TypeA, dns.ClassINET}
				c := new(dns.Client)
				in, _, err := c.Exchange(m1, "127.0.0.1:53")
				if err != nil {
					log.Printf("%s", err.Error())
				}
				for _, rr := range in.Answer {
					m.Answer = append(m.Answer, rr)
				}
			}
		}

	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)

	}

	w.WriteMsg(m)
	log.Printf("%s %s %s\n", w.RemoteAddr(), m.Question[0].Name, m.Answer)
}

func main() {
	// attach request handler func
	dns.HandleFunc(".", handleDnsRequest)

	// start server
	port := 5300
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())

	}

}
