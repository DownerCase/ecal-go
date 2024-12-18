package monitoring

import (
	"os"
	"testing"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/registration"
	"github.com/DownerCase/ecal-go/internal/ecaltest"
	"github.com/DownerCase/ecal-go/internal/ecaltest/regtest"
	testutilpublisher "github.com/DownerCase/ecal-go/internal/ecaltest/testutil_publisher"
	testutilsubscriber "github.com/DownerCase/ecal-go/internal/ecaltest/testutil_subscriber"
)

func expectTopicPresent(t *testing.T, ts []TopicMon, topicName string) {
	t.Helper()

	if len(ts) == 0 {
		t.Error("Monitoring returned no topics")
	}

	for _, topic := range ts {
		if topic.TopicName == topicName {
			return
		}
	}

	t.Error("Monitoring does not contain expected topic", topicName, "\nReceived", ts)
}

func TestPublisherMonitoring(t *testing.T) {
	ecaltest.InitEcal(t)

	defer ecal.Finalize()

	topic := "test_mon_pub"
	channel := make(chan regtest.Callback)
	registration.AddPublisherEventCallback(regtest.EventCallback(topic, channel))

	pub := testutilpublisher.NewGenericPublisher(t, topic)
	defer pub.Delete()

	mon := GetMonitoring(MonitorHost)
	if len(mon.Publishers) > 0 {
		t.Error("Monitoring returned publishers when not requested")
	}

	// Wait for publisher to be registered
	regtest.ExpectNew(t, topic, channel)

	// Get its info
	mon = GetMonitoring(MonitorPublisher)
	expectTopicPresent(t, mon.Publishers, topic)
}

func expectPid(t *testing.T, pid int, procs []ProcessMon) {
	t.Helper()

	hostname, err := os.Hostname()
	if err != nil {
		t.Error("Could not get hostname")
	}

	for _, proc := range procs {
		if pid == int(proc.Pid) {
			if proc.HostName != hostname {
				t.Error("Expected hostname", hostname, "got", proc.HostName)
			}

			return
		}
	}

	t.Error("Could not find self in process list")
}

func TestSubscriberMonitoring(t *testing.T) {
	ecaltest.InitEcal(t)

	defer ecal.Finalize()

	topic := "test_mon_sub"
	channel := make(chan regtest.Callback)
	registration.AddSubscriberEventCallback(regtest.EventCallback(topic, channel))

	sub := testutilsubscriber.NewGenericSubscriber(t, topic)
	defer sub.Delete()

	mon := GetMonitoring(MonitorHost)
	if len(mon.Publishers) > 0 {
		t.Error("Monitoring returned publishers when not requested")
	}

	// Wait for publisher to be registered
	regtest.ExpectNew(t, topic, channel)

	// Get its info
	mon = GetMonitoring(MonitorSubscriber)
	expectTopicPresent(t, mon.Subscribers, topic)
}

func TestProcessMonitoring(t *testing.T) {
	// Given: eCAL Initialized
	ecaltest.InitEcal(t)

	defer ecal.Finalize()

	time.Sleep(1500 * time.Millisecond) // Propagation delay...

	// When: Requesting the processes
	mon := GetMonitoring(MonitorProcess)

	// Expect: This program
	expectPid(t, os.Getpid(), mon.Processes)
}
