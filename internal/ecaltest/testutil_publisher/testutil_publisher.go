package testutilpublisher

import (
	"testing"

	"github.com/DownerCase/ecal-go/ecal/publisher"
)

func NewGenericPublisher(t *testing.T, topic string) *publisher.Publisher {
	t.Helper()

	pub, err := publisher.New(topic, publisher.DataType{})
	if err != nil {
		t.Error(err)
	}

	return pub
}
