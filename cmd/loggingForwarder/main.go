// a logging dns fowarder
//
// starts a dns server and forwards all queries to an upstream dns server.
// answers from the upstream server are replyed back to the client.
//
package main

import (
	"flag"
	"github.com/miekg/dns"
	"github.com/op/go-logging"
)

var upstreamDNS string
var listenAddr string
var logLevel string

var log = logging.MustGetLogger("loggingForwarder")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x} %{message} %{color:reset}`,
)

func init() {
	logging.SetFormatter(format)
	flag.StringVar(&upstreamDNS, "upstream", "8.8.8.8:53", "upstream DNS to forward requets to")
	flag.StringVar(&listenAddr, "addr", "127.0.0.1:5301", "port the forwarder should listen for dns requests")
	flag.StringVar(&logLevel, "logLevel", "INFO", "log level to be used")
}

// handleDNSRequest - sends queries to upstream servers and response back to client.
func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	log.Debugf("recieved msg %v from %s", r, w.RemoteAddr())
	// generate a reply msg
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	// if it is a query,
	// forward to upstream dns and write answer back to client
	case dns.OpcodeQuery:
		log.Debugf("contacting upstream '%s' to forward query", upstreamDNS)
		c := new(dns.Client)
		in, _, err := c.Exchange(r, upstreamDNS)
		if err != nil {
			log.Errorf("%s", err.Error())
		}
		log.Debugf("got reply from upstream")
		// append answers to the reply
		for _, rr := range in.Answer {
			m.Answer = append(m.Answer, rr)
		}
	}
	// write reply to client
	w.WriteMsg(m)
	// log request and reply
	log.Infof("%s %s %s", w.RemoteAddr(), m.Question[0].Name, m.Answer)
}

func main() {
	flag.Parse()
	level, err := logging.LogLevel(logLevel)
	if err != nil {
		log.Errorf("%s", err.Error())
	}
	logging.SetLevel(level, "loggingForwarder")
	// attach request handler func
	dns.HandleFunc(".", handleDnsRequest)

	// start server
	//port := 5301
	server := &dns.Server{Addr: listenAddr, Net: "udp"}
	log.Debugf("Starting at %s", listenAddr)
	err = server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s ", err.Error())

	}
}
