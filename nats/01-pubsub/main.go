package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL,
		nats.Name("nats"),
		nats.DiscoveredServersHandler(discoveredServersHandler),
		nats.ErrorHandler(errorHandler),
	)
	if err != nil {
		log.Panicf("connect to nats: %v", err)
	}

	defer func(nc *nats.Conn) { _ = nc.Drain() }(nc)

	log.Printf("Connected to NATS server: %s\n", nc.ConnectedUrl())
	log.Printf("Discovered servers: %v\n", nc.DiscoveredServers())
	log.Printf("Servers: %v\n", nc.Servers())
	log.Printf("Server name: %s\n", nc.ConnectedServerName())
	log.Printf("Cluster name: %s\n", nc.ConnectedClusterName())
	log.Printf("Server version: %s\n", nc.ConnectedServerVersion())
	log.Printf("Server ID: %s\n", nc.ConnectedServerId())
	log.Printf("Server address: %s\n", nc.ConnectedAddr())
	log.Printf("Auth required: %t\n", nc.AuthRequired())

	clientID, err := nc.GetClientID()
	if err == nil {
		log.Printf("Client ID: %d\n", clientID)
	}
	clientIP, err := nc.GetClientIP()
	if err == nil {
		log.Printf("Client IP: %s\n", clientIP)
	}

	_ = nc.Barrier(func() {
		log.Println("Barrier called")
	})

	buffered, err := nc.Buffered()
	if err == nil {
		log.Printf("Buffered message received: %v\n", buffered)
	}

	now := time.Now()
	err = nc.Flush()
	if err == nil {
		log.Printf("Flush completed in %v\n", time.Since(now))
	}

	sub, err := nc.QueueSubscribe("foo.*", "group", func(msg *nats.Msg) {
		log.Printf("Subject: %s, Reply: %s, Data: %s\n", msg.Subject, msg.Reply, string(msg.Data))
		err := msg.Respond([]byte("Goodbye!"))
		if err != nil {
			log.Printf("Error responding to message: %v\n", err)
		}
	})
	if err != nil {
		log.Panicf("subscribe to foo: %v", err)
	}

	defer func() { _ = sub.Drain() }()

	respInbox := nc.NewRespInbox()
	_, err = nc.Subscribe(respInbox, func(msg *nats.Msg) {
		log.Printf("Response: %s\n", string(msg.Data))
	})
	if err != nil {
		log.Printf("Error subscribing to response inbox: %v\n", err)
	}

	err = nc.PublishRequest("foo.boo", respInbox, []byte("Hello World Again"))
	if err != nil {
		log.Printf("Error publishing request: %v\n", err)
	}

	time.Sleep(2 * time.Second)
}

func errorHandler(nc *nats.Conn, sub *nats.Subscription, err error) {
	log.Printf("Error: %v\n", err)
	if sub != nil {
		log.Printf("Subscription: %v\n", sub.Subject)
	}
	log.Printf("In error handler\n")
	log.Printf("Discovered servers: %v\n", nc.DiscoveredServers())
	log.Printf("Connected servers: %v\n", nc.ConnectedServerName())
	log.Printf("Server URL: %s\n", nc.ConnectedUrl())
}

func discoveredServersHandler(nc *nats.Conn) {
	log.Printf("Discovered servers: %v\n", nc.DiscoveredServers())
}
