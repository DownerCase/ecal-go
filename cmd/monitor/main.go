package main

import (
	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/logging"
)

func main() {
	ecal.Initialize(
		ecal.NewConfig(ecal.WithLoggingReceive(true)),
		"Go Monitor",
		ecal.CSubscriber|ecal.CMonitoring|ecal.CLogging|ecal.CService,
	)
	defer ecal.Finalize() // Shutdown eCAL at the end of the program
	logging.SetConsoleFilter(logging.LevelAll)
	doCli()
}
