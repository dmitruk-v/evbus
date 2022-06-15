package main

import (
	"fmt"
	"time"

	ebclient "github.com/dmitruk-v/evbus/client"
)

func main() {
	ebClient := ebclient.New(":4000")
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
