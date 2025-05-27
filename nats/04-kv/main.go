package main

import "github.com/nats-io/nats.go"

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}

	kv, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "service-a",
		Description: "Key-Value store for service A",
		MaxValueSize: 1024, // 1 KB
		History: 1,
		TTL:    60 * 60 * 24, // 1 day
		Replicas: 3,
	})
	if err != nil {
		panic(err)
	}

	kv.
}
