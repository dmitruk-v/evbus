package main

import (
	"log"
	"os"

	ebServer "github.com/dmitruk-v/event-bus/server"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile | log.LstdFlags)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	server := ebServer.New(":4000", logger)
	if err := server.Listen(); err != nil {
		return err
	}
	return nil
}
