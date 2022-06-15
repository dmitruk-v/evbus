package main

import (
	"flag"
	"log"

	ebserv "github.com/dmitruk-v/evbus/server"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile | log.LstdFlags)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	addr := flag.String("addr", ":3366", "Address of the server. Examples: 172.16.1.14:4001, :4002")
	flag.Parse()

	server := ebserv.New(*addr)
	if err := server.Listen(); err != nil {
		return err
	}
	return nil
}
