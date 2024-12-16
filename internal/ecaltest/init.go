package ecaltest

import (
	"testing"

	"github.com/DownerCase/ecal-go/ecal"
)

func InitEcal(t *testing.T, opts ...ecal.ConfigOption) {
	initResult := ecal.Initialize(
		ecal.NewConfig(opts...),
		"Go eCAL!",
		ecal.CPublisher|ecal.CSubscriber|ecal.CLogging|ecal.CMonitoring,
	)
	if initResult != 0 {
		t.Fatal("Failed to initialize", initResult)
	}
}
