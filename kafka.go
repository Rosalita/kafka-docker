package main

import (
	"context"
	"log"
	"net"
	"strconv"

	kafka "github.com/segmentio/kafka-go"
)

func produce(ctx context.Context, topic string) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   topic,
	})

	msg := kafka.Message{
		Key:   []byte("A Key"),
		Value: []byte("A Value"),
	}

	if err := writer.WriteMessages(ctx, msg); err != nil {
		log.Fatal(err)
	}
}

func consume(ctx context.Context, topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   topic,
		GroupID: "my-group",
	})

	msg, err := reader.ReadMessage(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("message received: %+v\n", msg)
}

func createTopic(topic string) {
	// kafka.Dial randomly picks one of the brokers in the cluster.
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	brokers, err := conn.Brokers()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("brokers: %+v", brokers)

	// Use the connection to randomly chosen broker to get the leader broker.
	leader, err := conn.Controller()
	if err != nil {
		log.Fatal(err)
	}
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(leader.Host, strconv.Itoa(leader.Port)))
	if err != nil {
		log.Fatal(err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{{Topic: topic, NumPartitions: 1, ReplicationFactor: 1}}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Fatal(err)
	}
}
