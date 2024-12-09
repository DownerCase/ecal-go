package regtest

import (
	"testing"
	"time"

	"github.com/DownerCase/ecal-go/ecal/registration"
)

type Callback struct {
	Event registration.Event
	Id    registration.TopicId
}

func EventCallback(topic string, channel chan Callback) func(registration.TopicId, registration.Event) {
	return func(id registration.TopicId, event registration.Event) {
		if id.Topic_name == topic {
			channel <- Callback{event, id}
		}
	}
}

func expectEvent(event registration.Event, t *testing.T, topic string, channel chan Callback) {
	var response Callback
	select {
	case response = <-channel:
	case <-time.After(3 * time.Second):
		t.Error("Registration timeout")
		return
	}
	if response.Id.Topic_name != topic {
		// Should be pre-filtered by callback
		t.Error("Unexpected event for topic", response.Id.Topic_name)
	} else if response.Event != event {
		t.Error("Expected event", event, "actual", response.Event)
	} else {
		time.Sleep(50 * time.Millisecond) // Small delay to allow eCAL to finish
	}
}

func ExpectNew(t *testing.T, topic string, channel chan Callback) {
	expectEvent(registration.ENTITY_NEW, t, topic, channel)
}

func ExpectDeleted(t *testing.T, topic string, channel chan Callback) {
	expectEvent(registration.ENTITY_DELETED, t, topic, channel)
}
