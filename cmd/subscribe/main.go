package main

import (
	"encoding/json"
	"fmt"
	"log"

	ebus "github.com/dmitruk-v/event-bus/client"
)

func main() {
	ebClient := ebus.NewClient(":4000")
	ebClient.Subscribe("bla", func(data []byte) error {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Got:", s)
		return nil
	})
}
