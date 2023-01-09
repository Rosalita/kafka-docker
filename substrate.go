package main

import (
	"context"
	"fmt"
	"github.com/uw-labs/substrate"
	"github.com/uw-labs/substrate/kafka"
	"log"
)

type message struct {
	message string
}

func (m message) Data() []byte {
	return []byte(m.message)
}

func substratePublish(ctx context.Context) {
	// acknowledgments is a channel where messages will go once they have been published.
	acknowledgements := make(chan substrate.Message)

	// toPublish is a channel of messages to be published.
	toPublish := make(chan substrate.Message, 1)

	// add a message to be published.
	toPublish <- message{message: "Hello Substrate"}

	cfg := kafka.AsyncMessageSinkConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "substrate",
	}

	msgSink, err := kafka.NewAsyncMessageSink(cfg)
	if err != nil {
		log.Fatal()
	}

	go msgSink.PublishMessages(ctx, acknowledgements, toPublish)

	fmt.Printf("there are %d acknowledgments\n", len(acknowledgements))

	ack := <-acknowledgements

	fmt.Println("acknowledgement received", ack)
}
