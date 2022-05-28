package main

import (
	"encoding/json"
	"fmt"
	"log"

	ebClient "github.com/dmitruk-v/event-bus/client"
)

func main() {
	client, err := ebClient.Connect(":4000")
	if err != nil {
		log.Fatal(err)
	}
	client.Subscribe("bla", func(data []byte) error {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Got:", s)
		return nil
	})
	client.Run()
}
