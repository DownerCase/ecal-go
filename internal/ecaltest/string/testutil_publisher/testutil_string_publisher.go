package testutilpublisher

import (
	"testing"

	"github.com/DownerCase/ecal-go/ecal/string/publisher"
)

func NewStringPublisher(t *testing.T, topic string) *publisher.Publisher {
	t.Helper()
	pub, err := publisher.New()
	if err != nil {
		t.Error(err)
	}

	if err := pub.Create(topic); err != nil {
		t.Error(err)
	}
	return pub
}
