package main

import (
	"flag"
	"log"

	"github.com/unix-world/smartgoplus/cloud/msgauth/dmarc"
)

func main() {
	flag.Parse()

	domain := flag.Arg(0)
	if domain == "" {
		log.Fatal("usage: dmarc-lookup <domain>")
	}

	rec, err := dmarc.Lookup(domain)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", rec)
}
