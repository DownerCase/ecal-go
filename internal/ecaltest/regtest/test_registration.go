package regtest

import (
	"testing"
	"time"

	"github.com/DownerCase/ecal-go/ecal/registration"
)

type Callback struct {
	Event registration.Event
	ID    registration.TopicID
}

func EventCallback(topic string, channel chan Callback) func(registration.TopicID, registration.Event) {
	return func(id registration.TopicID, event registration.Event) {
		if id.TopicName == topic {
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
	if response.ID.TopicName != topic {
		// Should be pre-filtered by callback
		t.Error("Unexpected event for topic", response.ID.TopicName)
	} else if response.Event != event {
		t.Error("Expected event", event, "actual", response.Event)
	} else {
		time.Sleep(50 * time.Millisecond) // Small delay to allow eCAL to finish
	}
}

func ExpectNew(t *testing.T, topic string, channel chan Callback) {
	expectEvent(registration.EntityNew, t, topic, channel)
}

func ExpectDeleted(t *testing.T, topic string, channel chan Callback) {
	expectEvent(registration.EntityDeleted, t, topic, channel)
}
