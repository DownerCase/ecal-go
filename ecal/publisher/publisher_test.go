package publisher_test

import (
	"testing"

	"github.com/DownerCase/ecal-go/ecal/publisher"
	"github.com/DownerCase/ecal-go/internal/ecaltest/testutil_publisher"
)

func TestNewPublishers(t *testing.T) {
	for range 100 {
		ptr, err := publisher.New()
		if err != nil {
			t.Error(err)
		}
		defer ptr.Delete()
		t.Log(ptr)
	}
}

func TestPublisher(t *testing.T) {
	pub := testutil_publisher.NewGenericPublisher(t, "testing")
	defer pub.Delete()
	if pub.Messages == nil {
		t.Error("Message channel nil")
	}

}
