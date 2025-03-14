package main

import (
	"github.com/DownerCase/ecal-go/ecal"
)

func main() {
	ecal.Initialize(
		ecal.NewConfig(
			ecal.WithUDPLogReceiving(true),
		),
		"Go Monitor",
		ecal.CSubscriber|ecal.CMonitoring|ecal.CLogging|ecal.CService,
	)
	defer ecal.Finalize() // Shutdown eCAL at the end of the program

	ecal.SetState(ecal.ProcSevHealthy, ecal.ProcSevLevel1, "Monitoring eCAL")
	doCli()
}
