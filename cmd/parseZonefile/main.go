package main

import (
	"bufio"
	"github.com/miekg/dns"
	"log"
	"os"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// create a reader for a given file
	f, err := os.Open("./zone.txt")
	check(err)
	r := bufio.NewReader(f)
	// parseZonefile from reader
	for x := range dns.ParseZone(r, "", "zone.txt") {
		//check(x.Error)
		if x.Error != nil {
			log.Printf("err: '%s'", err.Error)
		}
		log.Printf("%v", x.RR)
	}
}
