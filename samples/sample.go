package main

import (
	"fmt"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
)

func main() {
	// Initialize eCAL
	ecal.Initialize(
		ecal.NewConfig(ecal.WithConsoleLogging(true), ecal.WithConsoleLogAll()),
		"Go eCAL!",
		ecal.CDefault,
	)
	defer ecal.Finalize() // Shutdown eCAL at the end of the program

	// Send messages
	go func() {
		publisher, err := ecal.NewStringPublisher("string topic")
		if err != nil {
			panic("Failed to make string publisher")
		}
		defer publisher.Delete()

		for idx := 0; true; idx++ {
			publisher.Send(fmt.Sprint("This is message ", idx))

			time.Sleep(time.Second)
		}
	}()

	// Receive messages
	subscriber, err := ecal.NewStringSubscriber("string topic")
	if err != nil {
		panic("Failed to Create string subscriber")
	}
	defer subscriber.Delete()

	for ecal.Ok() {
		msg, err := subscriber.Receive(time.Second * 2)
		if err == nil {
			fmt.Println("Received:", msg) //nolint:forbidigo
		} else {
			fmt.Println(err) //nolint:forbidigo
		}
	}
}
