package registration_test

import (
	"testing"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/registration"
	"github.com/DownerCase/ecal-go/internal/ecaltest"
	"github.com/DownerCase/ecal-go/internal/ecaltest/regtest"
	testutilpublisher "github.com/DownerCase/ecal-go/internal/ecaltest/testutil_publisher"
	testutilsubscriber "github.com/DownerCase/ecal-go/internal/ecaltest/testutil_subscriber"
)

func TestPublisherCallback(t *testing.T) {
	ecaltest.InitEcal(t)

	defer ecal.Finalize()

	topic := "test_reg_pub"
	channel := make(chan regtest.Callback)

	registration.AddPublisherEventCallback(regtest.EventCallback(topic, channel))

	pub := testutilpublisher.NewGenericPublisher(t, topic)
	defer pub.Delete()

	regtest.ExpectNew(t, topic, channel)

	// Destroy our publisher
	pub.Delete()
	regtest.ExpectDeleted(t, topic, channel)
}

func TestSubscriberCallback(t *testing.T) {
	ecaltest.InitEcal(t)

	defer ecal.Finalize()

	topic := "test_reg_sub"
	channel := make(chan regtest.Callback)

	registration.AddSubscriberEventCallback(regtest.EventCallback(topic, channel))

	sub := testutilsubscriber.NewBinarySubscriber(t, topic)
	defer sub.Delete()

	regtest.ExpectNew(t, topic, channel)

	// Destroy our publisher
	sub.Delete()
	regtest.ExpectDeleted(t, topic, channel)
}
