package main

import (
	"fmt"
	"log"
	"time"

	ebclient "github.com/dmitruk-v/evbus/client"
	"github.com/dmitruk-v/retry"
)

func main() {
	var ebClient *ebclient.EventBusClient
	err := retry.OnPanic(3, time.Second/2, func() {
		ebClient = ebclient.New(":3366")
	})
	if err != nil {
		log.Fatal(err)
	}
	go ebClient.Run()

	for i := 0; i < 3; i++ {
		ebClient.Emit("bla", fmt.Sprintf("%d-data", i))
	}

	time.Sleep(time.Second)

	for i := 3; i < 6; i++ {
		ebClient.Emit("bla", fmt.Sprintf("%d-data", i))
	}

	time.Sleep(time.Second)
}
