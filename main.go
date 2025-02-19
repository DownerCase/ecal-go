package main

import (
	"fmt"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/ecallog"
	"github.com/DownerCase/ecal-go/ecal/ecaltypes"
	"github.com/DownerCase/ecal-go/ecal/registration"
	"github.com/DownerCase/ecal-go/protos"
)

func main() {
	// eCAL version as string and semantic version components
	fmt.Println(ecal.GetVersionString()) //nolint:forbidigo
	fmt.Println(ecal.GetVersion())       //nolint:forbidigo

	// Initialize eCAL with default config, the unit name "Go eCAL",
	// and the Publisher, Subscriber and Logging components enabled
	initResult := ecal.Initialize(
		ecal.NewConfig(ecal.WithConsoleLogging(true), ecal.WithConsoleLogAll(), ecal.WithUDPLogAll()),
		"Go eCAL!",
		ecal.CPublisher|ecal.CSubscriber|ecal.CLogging,
	)
	defer ecal.Finalize() // Shutdown eCAL at the end of the program

	ecal.SetState(ecal.ProcSevHealthy, ecal.ProcSevLevel1, "Running the ecal-go demo")

	// Log a message
	ecallog.Infof("Initialized: %t", initResult)
	ecallog.Warn("Consider this a warning")
	ecallog.Error("This is an error")
	ecallog.Log(ecallog.LogLevelFatal, "When things are really bad")
	ecallog.Log(ecallog.LogLevelDebug1, "Level 1 debug message")
	ecallog.Log(ecallog.LogLevelDebug2, "Level 2 debug message")
	ecallog.Log(ecallog.LogLevelDebug3, "Level 3 debug message")
	ecallog.Log(ecallog.LogLevelDebug4, "Level 4 debug message")

	registration.AddPublisherEventCallback(registrationLogger)

	// Check if the eCAL system is Ok.
	// Other eCAL programs can send a message to cause ecal.Ok() to return false
	// Typically used as a condition to terminate daemon-style programs
	ecallog.Infof("eCAL ok: %t", ecal.Ok())

	// Create new protobuf publisher
	pub, err := ecal.NewProtobufPublisher[protos.Person]("person")
	if err != nil {
		panic("Failed to make new publisher")
	}
	defer pub.Delete() // Don't forget to delete the publisher when done!

	protoSub, err := ecal.NewProtobufSubscriber[protos.Person]("person")
	if err != nil {
		panic("Failed to make protobuf subscriber")
	}
	defer protoSub.Delete()

	stringPublisher, err := ecal.NewStringPublisher("string topic")
	if err != nil {
		panic("Failed to make string publisher")
	}

	sub, err := ecal.NewStringSubscriber("string topic")
	if err != nil {
		panic("Failed to make string subscriber")
	}

	go receiveMessages(sub, protoSub)

	sendMessages(100, stringPublisher, pub)
}

func sendMessages(
	numToSend int32,
	stringPublisher *ecal.StringPublisher,
	pub *ecal.ProtobufPublisher[*protos.Person],
) {
	person := &protos.Person{
		Id: 0, Name: "John", Email: "john@doe.net",
		Dog:   &protos.Dog{Name: "Pluto"},
		House: &protos.House{Rooms: 5},
	}

	for idx := range numToSend {
		// Check if program has been requested to stop
		if !ecal.Ok() {
			ecallog.Warn("eCAL.Ok() is false; shutting down")
			return
		}

		ecallog.Info("Sending message ", idx)

		// Update message to send
		person.Id = idx

		// Serialize and send protobuf message
		pub.Send(person)

		stringPublisher.Send(fmt.Sprint("Message ", idx))

		// Delay next iteration
		time.Sleep(1 * time.Second)
	}
}

func receiveMessages(s *ecal.StringSubscriber, s2 *ecal.ProtobufSubscriber[*protos.Person]) {
	for {
		msg, err := s.Receive(2 * time.Second)
		if err == nil {
			fmt.Println("Received:", msg) //nolint:forbidigo
		} else {
			fmt.Println(err) //nolint:forbidigo
		}

		msg2, err := s2.Receive(2 * time.Second)
		if err == nil {
			fmt.Println("Received:", msg2) //nolint:forbidigo
		} else {
			fmt.Println(err) //nolint:forbidigo
		}
	}
}

func registrationLogger(id ecaltypes.TopicID, _ registration.Event) {
	fmt.Println("Received registration sample:", id) //nolint:forbidigo
}
