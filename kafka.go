package main

import (
	"log"
	"net"
	"strconv"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	createTopic("foo")
}

func createTopic(topic string){
	// kafka.Dial randomly picks one of the brokers in the cluster.
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

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
