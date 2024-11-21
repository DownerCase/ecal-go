package main

import (
	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/logging"
	"github.com/DownerCase/ecal-go/ecal/registration"
)

type topicEntry struct {
	topic string
	info  registration.QualityTopicInfo
}

// The string is the TopicId.EntityId
var topicMap map[string]topicEntry
var topicMapHasUpdate bool = false

func main() {

	topicMap = make(map[string]topicEntry)

	ecal.Initialize(ecal.NewConfig(), "Go Monitor", ecal.C_Subscriber|ecal.C_Monitoring)
	defer ecal.Finalize() // Shutdown eCAL at the end of the program
	logging.SetConsoleFilter(logging.LevelAll)
	doCli()
}

