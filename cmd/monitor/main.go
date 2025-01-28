package main

import (
	"github.com/DownerCase/ecal-go/ecal"
)

func main() {
	ecal.Initialize(
		ecal.NewConfig(),
		"Go Monitor",
		ecal.CSubscriber|ecal.CMonitoring|ecal.CLogging|ecal.CService,
	)
	defer ecal.Finalize() // Shutdown eCAL at the end of the program
	doCli()
}
