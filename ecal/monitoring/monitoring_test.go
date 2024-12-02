package monitoring

import (
	"testing"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/registration"
	"github.com/DownerCase/ecal-go/internal/ecaltest"
	"github.com/DownerCase/ecal-go/internal/ecaltest/regtest"
	"github.com/DownerCase/ecal-go/internal/ecaltest/testutil_publisher"
	"github.com/DownerCase/ecal-go/internal/ecaltest/testutil_subscriber"
)

func expectTopicPresent(t *testing.T, ts []TopicMon, topic_name string) {
	if len(ts) == 0 {
		t.Error("Monitoring returned no topics")
	}
	for _, topic := range ts {
		if topic.Topic_name == topic_name {
			return
		}
	}
	t.Error("Monitoring does not contain expected topic", topic_name, "\nReceived", ts)
}

func TestPublisherMonitoring(t *testing.T) {
	ecaltest.InitEcal(t)
	defer ecal.Finalize()

	topic := "test_mon_pub"
	channel := make(chan regtest.Callback)
	registration.AddPublisherEventCallback(regtest.EventCallback(topic, channel))

	pub := testutil_publisher.NewGenericPublisher(t, topic)
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

func TestSubscriberMonitoring(t *testing.T) {
	ecaltest.InitEcal(t)
	defer ecal.Finalize()

	topic := "test_mon_sub"
	channel := make(chan regtest.Callback)
	registration.AddSubscriberEventCallback(regtest.EventCallback(topic, channel))

	sub := testutil_subscriber.NewGenericSubscriber(t, topic)
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
