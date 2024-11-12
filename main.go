package main

import (
	"fmt"
	"github.com/DownerCase/ecal-go/ecal"
	"time"
)

func main() {
	fmt.Println(ecal.GetVersionString())
	version := ecal.GetVersion()
	fmt.Println(version)
	fmt.Println("Init: ", ecal.Initialize(ecal.NewConfig(), "Go eCAL!", ecal.C_Publisher|ecal.C_Subscriber|ecal.C_Logging))
	defer func() { fmt.Println("Finalize: ", ecal.Finalize()) }()
	fmt.Println("Changed name: ", ecal.SetUnitName("Something new"))
	fmt.Println("Am ok?", ecal.Ok())
	publisher, err := ecal.NewPublisher()
	if err != nil {
		panic("Failed to make new publisher")
	}
	defer ecal.DestroyPublisher(&publisher)
	if publisher.Create("example topic") != nil {
		panic("Failed to Create publisher")
	}
	fmt.Println("Got publisher: ", publisher)
	for idx := range 100 {
		fmt.Println("Sending message")
		publisher.Messages <- []byte{'A', 'B', byte(idx)}
		time.Sleep(1 * time.Second)
	}
}
