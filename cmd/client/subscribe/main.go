package main

import (
	"encoding/json"
	"fmt"
	"log"

	ebclient "github.com/dmitruk-v/evbus/client"
)

func main() {
	ebClient := ebclient.New(":4000")
	ebClient.Subscribe("bla", func(data []byte) error {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Got:", s)
		return nil
	})
}
