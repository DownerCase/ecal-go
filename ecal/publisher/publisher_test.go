package publisher_test

import (
	"fmt"
	"testing"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/publisher"
	testutilpublisher "github.com/DownerCase/ecal-go/internal/ecaltest/testutil_publisher"
)

func TestNewPublishers(t *testing.T) {
	for i := range 100 {
		ptr, err := publisher.New(fmt.Sprintf("testPubTopic-%v", i), ecal.DataType{})
		if err != nil {
			t.Error(err)
		}
		defer ptr.Delete()
	}
}

func TestPublisher(t *testing.T) {
	pub := testutilpublisher.NewGenericPublisher(t, "testing")
	defer pub.Delete()

	if pub.Messages == nil {
		t.Error("Message channel nil")
	}
}
