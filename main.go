package main

import (
	"fmt"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/logging"
	"github.com/DownerCase/ecal-go/ecal/protobuf/publisher"
	"github.com/DownerCase/ecal-go/ecal/registration"
	string_publisher "github.com/DownerCase/ecal-go/ecal/string/publisher"
	"github.com/DownerCase/ecal-go/ecal/string/subscriber"
	"github.com/DownerCase/ecal-go/protos"
)

func main() {
	// eCAL version as string and semantic version components
	fmt.Println(ecal.GetVersionString()) //nolint:forbidigo
	fmt.Println(ecal.GetVersion())       //nolint:forbidigo

	// Initialize eCAL with default config, the unit name "Go eCAL",
	// and the Publisher, Subscriber and Logging components enabled
	initResult := ecal.Initialize(
		ecal.NewConfig(),
		"Go eCAL!",
		ecal.CPublisher|ecal.CSubscriber|ecal.CLogging,
	)

	// Log a message
	logging.Log(logging.LevelInfo, "Initialized: ", initResult)

	defer ecal.Finalize() // Shutdown eCAL at the end of the program

	registration.AddPublisherEventCallback(registrationLogger)

	// Change the unit name
	logging.Debug("Changed name:", ecal.SetUnitName("Go demo"))

	// Check if the eCAL system is Ok.
	// Other eCAL programs can send a message to cause ecal.Ok() to return false
	// Typically used as a condition to terminate daemon-style programs
	logging.Infof("eCAL ok: %t", ecal.Ok())

	// Create new protobuf publisher
	pub, err := publisher.New[protos.Person]()
	if err != nil {
		panic("Failed to make new publisher")
	}
	defer pub.Delete() // Don't forget to delete the publisher when done!

	if pub.Create("person") != nil {
		panic("Failed to Create protobuf publisher")
	}

	stringPublisher, _ := string_publisher.New()
	if stringPublisher.Create("string topic") != nil {
		panic("Failed to Create string publisher")
	}

	sub, _ := subscriber.New()
	if sub.Create("string topic") != nil {
		panic("Failed to Create string subscriber")
	}

	go receiveMessages(sub)

	sendMessages(100, stringPublisher, pub)
}

func sendMessages(
	numToSend int32,
	stringPublisher *string_publisher.Publisher,
	pub *publisher.Publisher[*protos.Person],
) {
	person := &protos.Person{
		Id: 0, Name: "John", Email: "john@doe.net",
		Dog:   &protos.Dog{Name: "Pluto"},
		House: &protos.House{Rooms: 5},
	}

	for idx := range numToSend {
		// Check if program has been requested to stop
		if !ecal.Ok() {
			logging.Warn("eCAL.Ok() is false; shutting down")
			return
		}

		logging.Info("Sending message ", idx)

		// Update message to send
		person.Id = idx

		// Serialize and send protobuf message
		if err := pub.Send(person); err != nil {
			logging.Error(err)
		}

		if err := stringPublisher.Send("Message ", idx); err != nil {
			logging.Error(err)
		}

		// Delay next iteration
		time.Sleep(1 * time.Second)
	}
}

func receiveMessages(s *subscriber.Subscriber) {
	for {
		msg, err := s.Receive(2 * time.Second)
		if err == nil {
			fmt.Println("Received:", msg) //nolint:forbidigo
		} else {
			fmt.Println(err) //nolint:forbidigo
		}
	}
}

func registrationLogger(id registration.TopicID, _ registration.Event) {
	fmt.Println("Received registration sample:", id) //nolint:forbidigo
}
