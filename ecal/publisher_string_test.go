package ecal_test

import (
	"testing"

	"github.com/DownerCase/ecal-go/internal/ecaltest/testutil"
)

func TestStringPublisher(t *testing.T) {
	pub := testutil.NewStringPublisher(t, "test_string_publisher")
	defer pub.Delete()

	if pub.Messages == nil {
		t.Error("Message channel nil")
	}

	// TODO: Check datatype information
	pub.Send("my message")
}
