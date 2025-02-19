package ecal_test

import (
	"testing"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/internal/ecaltest"
	"github.com/DownerCase/ecal-go/internal/ecaltest/testutil"
	"github.com/DownerCase/ecal-go/protos"
)

func TestSubscriber(t *testing.T) {
	ecaltest.InitEcal(t)

	defer ecal.Finalize() // Shutdown eCAL at the end of the program

	pub := testutil.NewProtobufPublisher[protos.Person](t, "testing_protobuf_subscriber")
	defer pub.Delete()

	sub := testutil.NewProtobufSubscriber[protos.Person](t, "testing_protobuf_subscriber")
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

	sub := testutil.NewProtobufSubscriber[protos.Person](t, "testing_protobuf_subscriber_timeout")
	defer sub.Delete()

	msg, err := sub.Receive(50 * time.Millisecond)
	if err == nil {
		t.Error("Expected timeout, received message:", &msg)
	}
}

func sendMessages(p *ecal.ProtobufPublisher[*protos.Person]) {
	person := &protos.Person{
		Id: 0, Name: "John", Email: "john@doe.net",
		Dog:   &protos.Dog{Name: "Pluto"},
		House: &protos.House{Rooms: 5},
	}

	for !p.IsStopped() {
		p.Send(person)

		person.House.Rooms++

		time.Sleep(10 * time.Millisecond)
	}
}
