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
	fmt.Println("Init: ", ecal.Initialize(ecal.NewConfig(), "Go eCAL!", ecal.Publisher|ecal.Subscriber|ecal.Logging))
	fmt.Println("Changed name: ", ecal.SetUnitName("Something new"))
	fmt.Println("Am ok?", ecal.Ok())
	time.Sleep(10 * time.Second)
	fmt.Println("Finalize: ", ecal.Finalize())
}
