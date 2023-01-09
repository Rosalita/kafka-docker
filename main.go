package main

import (
	"context"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)

	// use kafka-go to create topic, produce and consume synchronously.
	createTopic("kafka-go")  // a topic is a message queue.
	produce(ctx, "kafka-go") // synchronously publish a message.
	consume(ctx, "kafka-go") // synchronously read a message.

	// create a new topic to test substrate.
	createTopic("substrate")

	substratePublish(ctx) // asyncronously publish a message.
}
