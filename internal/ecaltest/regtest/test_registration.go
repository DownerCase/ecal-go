package regtest

import (
	"testing"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/registration"
)

const (
	RegistrationTimeout  = 3 * time.Second
	SynchronizationDelay = 50 * time.Millisecond
)

type Callback struct {
	Event registration.Event
	ID    ecal.TopicID
}

func EventCallback(topic string, channel chan Callback) func(ecal.TopicID, registration.Event) {
	return func(id ecal.TopicID, event registration.Event) {
		if id.TopicName == topic {
			channel <- Callback{event, id}
		}
	}
}

func expectEvent(t *testing.T, event registration.Event, topic string, channel chan Callback) {
	t.Helper()

	var response Callback
	select {
	case response = <-channel:
	case <-time.After(RegistrationTimeout):
		t.Error("Registration timeout")
		return
	}

	switch {
	case response.ID.TopicName != topic:
		// Should be pre-filtered by callback
		t.Error("Unexpected event for topic", response.ID.TopicName)
	case response.Event != event:
		t.Error("Expected event", event, "actual", response.Event)
	default:
		time.Sleep(SynchronizationDelay) // Small delay to allow eCAL to finish
	}
}

func ExpectNew(t *testing.T, topic string, channel chan Callback) {
	t.Helper()
	expectEvent(t, registration.EntityNew, topic, channel)
}

func ExpectDeleted(t *testing.T, topic string, channel chan Callback) {
	t.Helper()
	expectEvent(t, registration.EntityDeleted, topic, channel)
}
