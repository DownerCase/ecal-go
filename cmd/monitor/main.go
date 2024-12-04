package main

import (
	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/logging"
)

func main() {
	ecal.Initialize(ecal.NewConfig(), "Go Monitor", ecal.C_Subscriber|ecal.C_Monitoring|ecal.C_Logging)
	defer ecal.Finalize() // Shutdown eCAL at the end of the program
	logging.SetConsoleFilter(logging.LevelAll)
	doCli()
}
