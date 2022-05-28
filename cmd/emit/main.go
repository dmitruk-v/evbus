package main

import (
	"fmt"
	"log"
	"time"

	ebClient "github.com/dmitruk-v/event-bus/client"
)

func main() {
	client, err := ebClient.Connect(":4000")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()

	for i := 0; i < 3; i++ {
		if err := client.Emit("bla", fmt.Sprintf("%d-data", i)); err != nil {
			log.Fatal(err)
		}
	}

	time.Sleep(time.Second)

	for i := 3; i < 6; i++ {
		if err := client.Emit("bla", fmt.Sprintf("%d-data", i)); err != nil {
			log.Fatal(err)
		}
	}

	time.Sleep(time.Second * 1)
}
