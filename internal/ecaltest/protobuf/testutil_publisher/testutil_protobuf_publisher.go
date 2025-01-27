package testutilpublisher

import (
	"testing"

	"github.com/DownerCase/ecal-go/ecal/protobuf/publisher"
)

func NewProtobufPublisher[U any, T publisher.Msg[U]](t *testing.T, topic string) *publisher.Publisher[T] {
	t.Helper()

	pub, err := publisher.New[U, T](topic)
	if err != nil {
		t.Error(err)
	}

	return pub
}
