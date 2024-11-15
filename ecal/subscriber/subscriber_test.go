package subscriber

import (
	"reflect"
	"testing"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/publisher"
)

var TEST_MESSAGE = []byte{4, 15, 80}

func TestSubscriber(t *testing.T) {
	initResult := ecal.Initialize(
		ecal.NewConfig(),
		"Go eCAL!",
		ecal.C_Publisher|ecal.C_Subscriber|ecal.C_Logging,
	)
	if initResult != 0 {
		t.Fatal("Failed to initialize", initResult)
	}
	defer ecal.Finalize() // Shutdown eCAL at the end of the program

	pub, err := publisher.New()
	if err != nil {
		t.Error(err)
	}
	defer pub.Delete()

	if err := pub.Create("testing_subscriber", DataType{}); err != nil {
		t.Error(err)
	}

	sub, err := New()
	if err != nil {
		t.Error(err)
	}
	defer sub.Delete()
	if err := sub.Create("testing_subscriber", DataType{}); err != nil {
		t.Error(err)
	}

	go sendMessages(pub)
	for range 10 {
		msg := sub.Receive()
		if msg == nil {
			t.Error("Nil message received:")
		}
		if len(msg) != len(TEST_MESSAGE) {
			t.Error("Expected message of length", len(TEST_MESSAGE), "Received:", len(msg))
		}
		if !reflect.DeepEqual(msg, TEST_MESSAGE) {
			t.Error(msg, "!=", TEST_MESSAGE)
		}
	}
}

func sendMessages(p *publisher.Publisher) {
	for !p.IsStopped() {
		p.Messages <- TEST_MESSAGE
		time.Sleep(10 * time.Millisecond)
	}
}
