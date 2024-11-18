package testutil_publisher

import (
	"testing"

	"github.com/DownerCase/ecal-go/ecal/protobuf/publisher"
)

func NewProtobufPublisher[U any, T publisher.Msg[U]](t *testing.T, topic string) *publisher.Publisher[T] {
	pub, err := publisher.New[U, T]()

	if err != nil {
		t.Error(err)
	}

	if err := pub.Create(topic); err != nil {
		t.Error(err)
	}
	return pub
}
