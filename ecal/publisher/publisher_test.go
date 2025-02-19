package publisher_test

import (
	"fmt"
	"testing"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/internal/ecaltest/testutil"
)

func TestNewPublishers(t *testing.T) {
	publishers := make([]*ecal.BinaryPublisher, 100)
	for i := range publishers {
		publishers[i] = testutil.NewBinaryPublisher(t, fmt.Sprintf("testPubTopic-%v", i))
	}

	for _, p := range publishers {
		p.Delete()
	}
}

func TestPublisher(t *testing.T) {
	pub := testutil.NewBinaryPublisher(t, "testing")
	defer pub.Delete()

	if pub.Messages == nil {
		t.Error("Message channel nil")
	}
}
