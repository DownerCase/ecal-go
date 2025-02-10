package main

import (
	"context"
	"sync"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/publisher"
	"github.com/DownerCase/ecal-go/ecal/subscriber"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const topic = "go-benchmark"

func main() {
	ecal.Initialize(
		ecal.NewConfig(),
		"Go Throughput",
		ecal.CSubscriber|ecal.CPublisher,
	)
	defer ecal.Finalize() // Shutdown eCAL at the end of the program

	ecal.SetState(ecal.ProcSevHealthy, ecal.ProcSevLevel1, "Measuring eCAL throughput")

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)

	go measurePublish(ctx, &wg)
	go measureReceive(&wg, cancel)
	wg.Wait()
}

func measurePublish(context context.Context, wg *sync.WaitGroup) {
	publisher, err := publisher.New(topic, ecal.DataType{})
	if err != nil {
		panic(err)
	}

	defer wg.Done()
	defer publisher.Delete()

	payload := make([]byte, 4*1024*1024)
	counter := 0

	p := message.NewPrinter(language.English)
	defer func() { p.Printf("Sent %d messages\n", counter) }()

	for {
		select {
		case <-context.Done():
			return
		default:
			if len(publisher.Messages) < 100 {
				publisher.Messages <- payload

				counter++
			}
		}
	}
}

func measureReceive(wg *sync.WaitGroup, cancel context.CancelFunc) {
	subscriber, err := subscriber.New(topic, ecal.DataType{})
	if err != nil {
		panic(err)
	}

	defer wg.Done()
	defer subscriber.Delete()

	collectionDuration := 4 * time.Second
	bytesReceived := 0
	counter := 0

	time.Sleep(2 * time.Second)

	timeout := time.After(collectionDuration)

	for {
		select {
		case <-timeout:
			p := message.NewPrinter(language.English)
			p.Printf("Received %d bytes in %v seconds over %d messages\n", bytesReceived, collectionDuration, counter)
			p.Printf("Total: %d MB/s\n", bytesReceived/1000/1000/int(collectionDuration.Seconds()))
			cancel()

			return
		case msg := <-subscriber.Messages:
			bytesReceived += len(msg.([]byte))
			counter++
		}
	}
}
