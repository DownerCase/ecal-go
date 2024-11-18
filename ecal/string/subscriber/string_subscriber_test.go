package subscriber

import (
	"reflect"
	"testing"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/string/publisher"
	"github.com/DownerCase/ecal-go/internal/ecaltest"
	"github.com/DownerCase/ecal-go/internal/ecaltest/string/testutil_publisher"
)

var TEST_MESSAGE = "Test string"

func newSubscriber(t *testing.T, topic string) *Subscriber {
	sub, err := New()
	if err != nil {
		t.Error(err)
	}
	if err := sub.Create(topic); err != nil {
		t.Error(err)
	}
	return sub
}

// Export for testing
var NewSubscriber = newSubscriber

func TestSubscriber(t *testing.T) {
	ecaltest.InitEcal(t)
	defer ecal.Finalize() // Shutdown eCAL at the end of the program

	pub := testutilpublisher.NewStringPublisher(t, "testing_string_subscriber")
	defer pub.Delete()

	sub := newSubscriber(t, "testing_string_subscriber")
	defer sub.Delete()

	go sendMessages(pub)
	for range 10 {
		msg, err := sub.Receive(2 * time.Second)
		if err != nil {
			t.Error(err)
		}
		if len(msg) != len(TEST_MESSAGE) {
			t.Error("Expected message of length", len(TEST_MESSAGE), "Received:", len(msg))
		}
		if !reflect.DeepEqual(msg, TEST_MESSAGE) {
			t.Error(msg, "!=", TEST_MESSAGE)
		}
	}
}

func TestSubscriberTimeout(t *testing.T) {
	ecaltest.InitEcal(t)
	defer ecal.Finalize() // Shutdown eCAL at the end of the program
	sub := newSubscriber(t, "testing_string_subscriber_timeout")
	defer sub.Delete()
	msg, err := sub.Receive(50 * time.Millisecond)
	if err == nil {
		t.Error("Expected timeout, received message:", msg)
	}
}

func sendMessages(p *publisher.Publisher) {
	for !p.IsStopped() {
		p.Messages <- []byte(TEST_MESSAGE)
		time.Sleep(10 * time.Millisecond)
	}
}
