package subscriber

import (
	"testing"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/protobuf/publisher"
	"github.com/DownerCase/ecal-go/protos"
)

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

	pub, err := publisher.New[protos.Person]()
	if err != nil {
		t.Error(err)
	}
	defer pub.Delete()

	if err := pub.Create("testing_string_subscriber"); err != nil {
		t.Error(err)
	}

	sub, err := New[protos.Person]()
	if err != nil {
		t.Error(err)
	}
	defer sub.Delete()
	if err := sub.Create("testing_string_subscriber"); err != nil {
		t.Error(err)
	}

	go sendMessages(pub)
	for range 10 {
		msg, err := sub.Receive()

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

func sendMessages(p *publisher.Publisher[*protos.Person]) {
	person := &protos.Person{Id: 0, Name: "John", Email: "john@doe.net",
		Dog:   &protos.Dog{Name: "Pluto"},
		House: &protos.House{Rooms: 5},
	}

	for !p.IsStopped() {
		p.Send(person)
		person.House.Rooms += 1
		time.Sleep(10 * time.Millisecond)
	}
}
