package ecaltest

import (
	"testing"

	"github.com/DownerCase/ecal-go/ecal"
)

func InitEcal(t *testing.T, opts ...ecal.ConfigOption) {
	t.Helper()

	if !ecal.Initialize(
		ecal.NewConfig(opts...),
		"Go eCAL!",
		ecal.CPublisher|ecal.CSubscriber|ecal.CLogging|ecal.CMonitoring,
	) {
		t.Fatal("Failed to initialize eCAL")
	}
}
