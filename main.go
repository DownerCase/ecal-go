package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/protobuf/publisher"
	string_publisher "github.com/DownerCase/ecal-go/ecal/string/publisher"
	"github.com/DownerCase/ecal-go/protos"
)

func main() {
	// eCAL version as string and semantic version components
	fmt.Println(ecal.GetVersionString())
	fmt.Println(ecal.GetVersion())

	// Initialize eCAL with default config, the unit name "Go eCAL",
	// and the Publisher, Subscriber and Logging components enabled
	initResult := ecal.Initialize(
		ecal.NewConfig(),
		"Go eCAL!",
		ecal.C_Publisher|ecal.C_Subscriber|ecal.C_Logging,
	)
	fmt.Println("Init:", initResult)

	defer ecal.Finalize() // Shutdown eCAL at the end of the program

	// Change the unit name
	fmt.Println("Changed name:", ecal.SetUnitName("Go demo"))

	// Check if the eCAL system is Ok.
	// Other eCAL programs can send a message to cause ecal.Ok() to return false
	// Typically used as a condition to terminate daemon-style programs
	fmt.Println("eCAL ok?", ecal.Ok())

	// Create new protobuf publisher
	pub, err := publisher.New(&protos.Person{})
	if err != nil {
		panic("Failed to make new publisher")
	}
	defer pub.Delete() // Don't forget to delete the publisher when done!

	person := &protos.Person{Id: 0, Name: "John", Email: "john@doe.net",
		Dog:   &protos.Dog{Name: "Pluto"},
		House: &protos.House{Rooms: 5},
	}

	if pub.Create("person") != nil {
		panic("Failed to Create protobuf publisher")
	}

	string_pub, _ := string_publisher.New()
	if string_pub.Create("string topic") != nil {
		panic("Failed to Create string publisher")
	}

	for idx := range 100 {
		// Check if program has been requested to stop
		if !ecal.Ok() {
			fmt.Println("eCAL.Ok() is false; shutting down")
			return
		}

		fmt.Println("Sending message ", idx)

		// Update message to send
		person.Id = int32(idx)

		// Serialize and send protobuf message
		if err := pub.Send(person); err != nil {
			fmt.Println("Error: ", err)
		}

		string_msg := "Sent " + strconv.Itoa(idx) + " messages"
		if err = string_pub.Send(string_msg); err != nil {
			fmt.Println("Error: ", err)
		}

		// Delay next iteration
		time.Sleep(1 * time.Second)
	}
}
