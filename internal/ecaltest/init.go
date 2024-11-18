package ecaltest

import (
	"testing"

	"github.com/DownerCase/ecal-go/ecal"
)

func InitEcal(t *testing.T) {
	initResult := ecal.Initialize(
		ecal.NewConfig(),
		"Go eCAL!",
		ecal.C_Publisher|ecal.C_Subscriber|ecal.C_Logging,
	)
	if initResult != 0 {
		t.Fatal("Failed to initialize", initResult)
	}
}
