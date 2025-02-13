package main

import (
	"context"
	"sync"
	"time"
	"unsafe"

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

	for {
		select {
		case <-context.Done():
			return
		default:
			publisher.Messages <- payload
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

	subscriber.Deserialize = func(_ unsafe.Pointer, dataLen int) any {
		bytesReceived += dataLen
		counter++

		return nil
	}

	time.Sleep(2 * time.Second)

	<-time.After(collectionDuration)

	p := message.NewPrinter(language.English)
	p.Printf("Received %d bytes in %v seconds over %d messages\n", bytesReceived, collectionDuration, counter)
	p.Printf("Total: %d MB/s\n", bytesReceived/1024/1024/int(collectionDuration.Seconds()))
	cancel()
}
