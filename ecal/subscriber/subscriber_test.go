package subscriber_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/internal/ecaltest"
	"github.com/DownerCase/ecal-go/internal/ecaltest/testutil"
)

func GetTestMessage() []byte {
	return []byte{4, 15, 80}
}

func TestSubscriber(t *testing.T) {
	ecaltest.InitEcal(t)

	defer ecal.Finalize() // Shutdown eCAL at the end of the program

	pub := testutil.NewBinaryPublisher(t, "testing_subscriber")
	defer pub.Delete()

	sub := testutil.NewBinarySubscriber(t, "testing_subscriber")
	defer sub.Delete()

	go sendMessages(pub)

	TestMessage := GetTestMessage()

	for range 10 {
		// TODO: Reduce the propagation delay for when the subscriber gets
		// connected to the publisher
		msg, err := sub.Receive(2 * time.Second)
		if err != nil {
			t.Error("Received err:", err)
		}

		if msg == nil {
			t.Error("Nil message received:")
		}

		if len(msg) != len(TestMessage) {
			t.Error("Expected message of length", len(TestMessage), "Received:", len(msg))
		}

		if !reflect.DeepEqual(msg, TestMessage) {
			t.Error(msg, "!=", TestMessage)
		}
	}
}

func TestSubscriberTimeout(t *testing.T) {
	ecaltest.InitEcal(t)

	defer ecal.Finalize() // Shutdown eCAL at the end of the program

	sub := testutil.NewBinarySubscriber(t, "testing_subscriber_timeout")
	defer sub.Delete()

	msg, err := sub.Receive(50 * time.Millisecond)
	if err == nil {
		t.Error("Expected timeout, received message:", msg)
	}
}

func sendMessages(p *ecal.BinaryPublisher) {
	TestMessage := GetTestMessage()
	for !p.IsStopped() {
		p.Messages <- TestMessage

		time.Sleep(10 * time.Millisecond)
	}
}
