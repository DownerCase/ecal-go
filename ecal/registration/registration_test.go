package registration

import (
	"fmt"
	"testing"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/internal/ecaltest"
	"github.com/DownerCase/ecal-go/internal/ecaltest/testutil_publisher"
)

type callback struct {
	event Event
	id TopicId
}

// WARNING: This test relies on no other publishers being active whilst running
func TestPublisherCallback(t *testing.T) {
	ecaltest.InitEcal(t)
	defer ecal.Finalize()

	channel := make(chan callback)

	AddPublisherEventCallback(func(id TopicId, event Event) {
		registrationLogger(channel, id, event)
	})

	pub := testutil_publisher.NewGenericPublisher(t, "test_reg_pub")
	defer pub.Delete()

	added := <- channel
	if added.event != ENTITY_NEW {
		t.Error("Expected new entity, got", added.event)
	}
	if added.id.Topic_name != "test_reg_pub" {
		t.Error("Unexpected registration for", added.id.Topic_name)
	}

	// Destroy our publisher
	pub.Delete()
	removed := <- channel
	if removed.event != ENTITY_DELETED {
		t.Error("Expected deleted entity, got", removed.event)
	}
	if removed.id.Topic_name != "test_reg_pub" {
		t.Error("Unexpected registration for", removed.id.Topic_name)
	}
}

func registrationLogger(channel chan callback, id TopicId, event Event) {
	fmt.Println("Event:", event, "Received registration sample:", id)
	channel <- callback {event, id }
}
