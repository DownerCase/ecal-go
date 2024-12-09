package ecaltest

import (
	"testing"

	"github.com/DownerCase/ecal-go/ecal"
)

func InitEcal(t *testing.T, opts ...ecal.ConfigOption) {
	initResult := ecal.Initialize(
		ecal.NewConfig(opts...),
		"Go eCAL!",
		ecal.C_Publisher|ecal.C_Subscriber|ecal.C_Logging|ecal.C_Monitoring,
	)
	if initResult != 0 {
		t.Fatal("Failed to initialize", initResult)
	}
}
