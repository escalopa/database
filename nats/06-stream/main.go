package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("localhost")
	if err != nil {
		log.Fatal(err)
	}

	// Use the JetStream context to produce and consumer messages
	// that have been persisted.
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal(err)
	}

	js.AddStream(&nats.StreamConfig{
		Name:     "example-stream",
		Subjects: []string{"example-subject"},
	})

	js.Publish("example-subject", []byte("Hello JS!"))

	// Publish messages asynchronously.
	for i := 0; i < 500; i++ {
		js.PublishAsync("example-subject", []byte("Hello JS Async!"))
	}
	select {
	case <-js.PublishAsyncComplete():
		fmt.Println("All messages published")
	case <-time.After(5 * time.Second):
		fmt.Println("Did not resolve in time")
	}
}
