package protobuf_test

import (
	"testing"

	"github.com/DownerCase/ecal-go/internal/protobuf"
	"github.com/DownerCase/ecal-go/protos"
)

func TestFullName(t *testing.T) {
	t.Parallel()

	expectedName := "pb.People.Person"
	if fn := protobuf.GetFullName(&protos.Person{}); fn != expectedName {
		t.Error("Expected: ", expectedName, " Actual: ", fn)
	}
}
