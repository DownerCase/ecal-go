package main

import (
	"context"
	"sync"
	"time"
	"unsafe"

	"github.com/DownerCase/ecal-go/ecal"
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
	publisher, err := ecal.NewBinaryPublisher(topic)
	if err != nil {
		panic(err)
	}

	defer wg.Done()
	defer publisher.Delete()

	payload := make([]byte, 8*1024*1024)

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
	mutex := sync.Mutex{}
	bytesReceived := 0
	counter := 0

	subscriber, err := subscriber.NewGenericSubscriber(
		topic,
		ecal.DataType{},
		func(_ unsafe.Pointer, dataLen int) any {
			mutex.Lock()
			bytesReceived += dataLen
			counter++
			mutex.Unlock()

			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	defer wg.Done()
	defer subscriber.Delete()

	time.Sleep(2 * time.Second)

	mutex.Lock()
	bytesReceived = 0
	counter = 0
	before := time.Now()
	mutex.Unlock()

	<-time.After(40 * time.Second)

	mutex.Lock()
	after := time.Now()
	bytesSnapshot := bytesReceived
	counterSnapshot := counter
	mutex.Unlock()

	p := message.NewPrinter(language.English)
	captureDuration := after.Sub(before).Seconds()
	p.Printf("Received %d bytes in %.2f seconds over %d messages\n", bytesSnapshot, captureDuration, counterSnapshot)
	p.Printf("Total: %.0f MB/s\n", float64(bytesSnapshot/1024/1024)/captureDuration)
	cancel()
}
