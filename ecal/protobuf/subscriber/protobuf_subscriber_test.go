package subscriber_test

import (
	"testing"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/protobuf/publisher"
	"github.com/DownerCase/ecal-go/ecal/protobuf/subscriber"
	"github.com/DownerCase/ecal-go/internal/ecaltest"
	testutilpublisher "github.com/DownerCase/ecal-go/internal/ecaltest/protobuf/testutil_publisher"
	"github.com/DownerCase/ecal-go/protos"
)

func newSubscriber[U any, T subscriber.Msg[U]](t *testing.T, topic string) *subscriber.Subscriber[U, T] {
	t.Helper()

	sub, err := subscriber.New[U, T](topic)
	if err != nil {
		t.Error(err)
	}

	return sub
}

func TestSubscriber(t *testing.T) {
	ecaltest.InitEcal(t)

	defer ecal.Finalize() // Shutdown eCAL at the end of the program

	pub := testutilpublisher.NewProtobufPublisher[protos.Person](t, "testing_protobuf_subscriber")
	defer pub.Delete()

	sub := newSubscriber[protos.Person](t, "testing_protobuf_subscriber")
	defer sub.Delete()

	go sendMessages(pub)

	for range 10 {
		msg, err := sub.Receive(2 * time.Second)
		if err != nil {
			t.Error(err)
		}

		if msg.GetId() != 0 {
			t.Error("Wrong ID")
		}

		if msg.GetName() != "John" {
			t.Error("Wrong name")
		}

		if msg.GetEmail() != "john@doe.net" {
			t.Error("Wrong email")
		}

		if msg.GetDog().GetName() != "Pluto" {
			t.Error("Wrong dog")
		}
	}
}

func TestSubscriberTimeout(t *testing.T) {
	ecaltest.InitEcal(t)

	defer ecal.Finalize() // Shutdown eCAL at the end of the program

	sub := newSubscriber[protos.Person](t, "testing_protobuf_subscriber_timeout")
	defer sub.Delete()

	msg, err := sub.Receive(50 * time.Millisecond)
	if err == nil {
		t.Error("Expected timeout, received message:", &msg)
	}
}

func sendMessages(p *publisher.Publisher[*protos.Person]) {
	person := &protos.Person{
		Id: 0, Name: "John", Email: "john@doe.net",
		Dog:   &protos.Dog{Name: "Pluto"},
		House: &protos.House{Rooms: 5},
	}

	for !p.IsStopped() {
		_ = p.Send(person)
		person.House.Rooms++

		time.Sleep(10 * time.Millisecond)
	}
}
