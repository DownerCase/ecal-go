package registration

import (
	"fmt"
	"testing"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/internal/ecaltest"
	"github.com/DownerCase/ecal-go/internal/ecaltest/testutil_publisher"
	"github.com/DownerCase/ecal-go/internal/ecaltest/testutil_subscriber"
)

type callback struct {
	event Event
	id    TopicId
}

func expectNew(t *testing.T, topic string, channel chan callback) {
	added := <-channel
	if added.event != ENTITY_NEW {
		t.Error("Expected new entity, got", added.event)
	}
	if added.id.Topic_name != topic {
		t.Error("Unexpected registration for", added.id.Topic_name)
	}
}

func expectDeleted(t *testing.T, topic string, channel chan callback) {
	removed := <-channel
	if removed.event != ENTITY_DELETED {
		t.Error("Expected deleted entity, got", removed.event)
	}
	if removed.id.Topic_name != topic {
		t.Error("Unexpected registration for", removed.id.Topic_name)
	}
}

// WARNING: This test relies on no other topics being active whilst running
func TestPublisherCallback(t *testing.T) {
	ecaltest.InitEcal(t)
	defer ecal.Finalize()

	channel := make(chan callback)

	AddPublisherEventCallback(func(id TopicId, event Event) {
		registrationLogger(channel, id, event)
	})

	topic := "test_reg_pub"
	pub := testutil_publisher.NewGenericPublisher(t, topic)
	defer pub.Delete()

	expectNew(t, topic, channel)

	// Destroy our publisher
	pub.Delete()
	expectDeleted(t, topic, channel)
}

func TestSubscriberCallback(t *testing.T) {
	ecaltest.InitEcal(t)
	defer ecal.Finalize()

	channel := make(chan callback)

	AddSubscriberEventCallback(func(id TopicId, event Event) {
		registrationLogger(channel, id, event)
	})

	topic := "test_reg_sub"
	sub := testutil_subscriber.NewGenericSubscriber(t, topic)
	defer sub.Delete()

	expectNew(t, topic, channel)

	// Destroy our publisher
	sub.Delete()
	expectDeleted(t, topic, channel)
}

func registrationLogger(channel chan callback, id TopicId, event Event) {
	fmt.Println("Event:", event, "Received registration sample:", id)
	channel <- callback{event, id}
}
