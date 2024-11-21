package registration

import (
	"fmt"
	"testing"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/internal/ecaltest"
	"github.com/DownerCase/ecal-go/internal/ecaltest/testutil_publisher"
	"github.com/DownerCase/ecal-go/internal/ecaltest/testutil_subscriber"
)

type callback struct {
	event Event
	id    TopicId
}

func expectEvent(event Event, t *testing.T, topic string, channel chan callback) {
	var response callback
	select {
	case response = <-channel:
	case <-time.After(3 * time.Second):
		t.Error("Registration timeout")
		return
	}
	if response.id.Topic_name != topic {
		// Should be pre-filtered by callback
		t.Error("Unexpected event for topic", response.id.Topic_name)
	}
	if response.event != event {
		t.Error("Expected event", event, "actual", response.event)
	}
}

func TestPublisherCallback(t *testing.T) {
	ecaltest.InitEcal(t)
	defer ecal.Finalize()

	topic := "test_reg_pub"
	channel := make(chan callback)

	AddPublisherEventCallback(func(id TopicId, event Event) {
		fmt.Println("Event:", event, "Received registration sample:", id)
		if id.Topic_name == topic {
			channel <- callback{event, id}
		}
	})

	pub := testutil_publisher.NewGenericPublisher(t, topic)
	defer pub.Delete()

	expectEvent(ENTITY_NEW, t, topic, channel)

	// Destroy our publisher
	pub.Delete()
	expectEvent(ENTITY_DELETED, t, topic, channel)
}

func TestSubscriberCallback(t *testing.T) {
	ecaltest.InitEcal(t)
	defer ecal.Finalize()

	topic := "test_reg_sub"
	channel := make(chan callback)

	AddSubscriberEventCallback(func(id TopicId, event Event) {
		fmt.Println("Event:", event, "Received registration sample:", id)
		if id.Topic_name == topic {
			channel <- callback{event, id}
		}
	})

	sub := testutil_subscriber.NewGenericSubscriber(t, topic)
	defer sub.Delete()

	expectEvent(ENTITY_NEW, t, topic, channel)

	// Destroy our publisher
	sub.Delete()
	expectEvent(ENTITY_DELETED, t, topic, channel)
}
