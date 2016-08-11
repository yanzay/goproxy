package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	remoteHost      = flag.String("host", "", "Remote host")
	remotePort      = flag.Int("port", 0, "Remote port")
	listen          = flag.String("listen", ":4242", "Local address to listen")
	dump            = flag.String("dump", "", "Write dump to file")
	skipHealthcheck = flag.Bool("skip-healthcheck", false, "Skip healthcheck")
)

func main() {
	flag.Parse()
	remoteAddr := fmt.Sprintf("%s:%d", *remoteHost, *remotePort)
	proxy := &proxyServer{localAddr: *listen, remoteAddr: remoteAddr, dumpTo: dumpTo(*dump)}
	var err error
	if !*skipHealthcheck {
		err = proxy.healthcheck()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Healthcheck to %s OK", remoteAddr)
	}
	err = proxy.start()
	if err != nil {
		log.Fatal(err)
	}
}

func dumpTo(filename string) *os.File {
	dumpTo := os.Stdout
	if len(filename) > 0 {
		file, err := os.Create(filename)
		if err != nil {
			log.Printf("Fail to open file %s, fallback to stdout", filename)
		} else {
			dumpTo = file
		}
	}
	return dumpTo
}
