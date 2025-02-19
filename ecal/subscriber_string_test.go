package ecal_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/internal/ecaltest"
	"github.com/DownerCase/ecal-go/internal/ecaltest/testutil"
)

const TestMessage = "Test string"

func TestStringSubscriber(t *testing.T) {
	ecaltest.InitEcal(t)

	defer ecal.Finalize() // Shutdown eCAL at the end of the program

	pub := testutil.NewStringPublisher(t, "testing_string_subscriber")
	defer pub.Delete()

	sub := testutil.NewStringSubscriber(t, "testing_string_subscriber")
	defer sub.Delete()

	go sendStringMessages(pub)

	for range 10 {
		msg, err := sub.Receive(2 * time.Second)
		if err != nil {
			t.Error(err)
		}

		if len(msg) != len(TestMessage) {
			t.Error("Expected message of length", len(TestMessage), "Received:", len(msg))
		}

		if !reflect.DeepEqual(msg, TestMessage) {
			t.Error(msg, "!=", TestMessage)
		}
	}
}

func TestStringSubscriberTimeout(t *testing.T) {
	ecaltest.InitEcal(t)

	defer ecal.Finalize() // Shutdown eCAL at the end of the program

	sub := testutil.NewStringSubscriber(t, "testing_string_subscriber_timeout")
	defer sub.Delete()

	msg, err := sub.Receive(50 * time.Millisecond)
	if err == nil {
		t.Error("Expected timeout, received message:", msg)
	}
}

func sendStringMessages(p *ecal.StringPublisher) {
	for !p.IsStopped() {
		p.Messages <- TestMessage

		time.Sleep(10 * time.Millisecond)
	}
}
