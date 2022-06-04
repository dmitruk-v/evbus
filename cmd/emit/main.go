package main

import (
	"fmt"
	"time"

	ebus "github.com/dmitruk-v/event-bus/client"
)

func main() {
	ebClient := ebus.NewClient(":4000")

	for i := 0; i < 3; i++ {
		ebClient.Emit("bla", fmt.Sprintf("%d-data", i))
	}

	time.Sleep(time.Second)

	for i := 3; i < 6; i++ {
		ebClient.Emit("bla", fmt.Sprintf("%d-data", i))
	}

	time.Sleep(time.Second)
}
