package main

import (
	"fmt"
	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/publisher"
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
	pub, err := publisher.New()
	if err != nil {
		panic("Failed to make new publisher")
	}
	defer publisher.Destroy(&pub)
	if pub.Create("example topic") != nil {
		panic("Failed to Create publisher")
	}
	fmt.Println("Got publisher: ", pub)
	for idx := range 100 {
		fmt.Println("Sending message")
		pub.Messages <- []byte{'A', 'B', byte(idx)}
		time.Sleep(1 * time.Second)
	}
}
