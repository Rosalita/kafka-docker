# kafka-docker

For learning purposes, I wanted to try run Kafka locally and write some Go code to interact with it.

# Kafka and Zookeeper
Kafka has a number of brokers, one of which is the leader. 
Zookeeper handles leadership election, it also manages service discovery. 
Zookeeper can exist without Kafka, however Kafka can not exist without Zookeeper.
This repo contains a Docker Compose file which will create Kafka and Zookeeper in separate containers.
To start the containers, install Docker Desktop then install Task with:
`go install github.com/go-task/task/v3/cmd/task@latest`
Run Docker Desktop, then `task start`

# Kafka UI
It would be nice to get a web ui set up for kafka as unlike RabbitMQ it doesn't seem to come with one by default. 
I decided to set up [kafka-ui](https://github.com/provectus/kafka-ui/) and added it to a new container in the docker compose file. 
Configuring kafka-ui to run against a container required a bit more reading up on [kafka-listeners](https://www.confluent.io/en-gb/blog/kafka-listeners-explained/).
Once kafka-ui was running it was possible to create a new topic via the ui :tada:

To view Kafka ui start the containers with `task start` then navigate to `localhost:8080`.

# Go
To interact with Kafka there is [confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go) which appears to be 
a cgo wrapper around [librdkafka](https://github.com/confluentinc/librdkafka) but I decided to use [kafka-go](https://github.com/segmentio/kafka-go)
as it is written entirely in Go. 

When using kafka-go, only the leader broker can create topics. I wrote some code in `kafka.go` to create a topic. 
I can verify the topic is created with `docker logs broker` and I can also see it appear on the web ui too :heart_eyes:

# Producing messages
Produce is a publish that adds a message to a topic. Unlike other pub/sub systems, once a message is published it is saved and persists.
RabbitMQ can persist messages but it doesn't by default. 
Kafka always persists messages for a configurable duration of time, regardless of whether they have been consumed or not.
A producer has to send a message to a broker. A cluster can have one or more brokers. The broker is responsible for storing the message.

When a producer sends a message it includes a key, the key is hashed to create a partition assignment. 
Keys are used to put related events into the same partition.
For example a producer could publish messages about the weather to the topic `London` with the keys of `temperature` and `windspeed`. 
The `temperature` data would all exist in one partition and the `windspeed` data would all existing in another partition.
When Kafka has multiple brokers, it would keep replica partitions across multiple brokers. 
If the broker that held the `temperature` partition failed, Kafka would start serving consumers `temperature` information from a 
partition replica in a different broker.
Partitions are how Kafka allows paralellism as data about `temperature` can be worked on at the same time as data about `windspeed`

Once a producer has sent a message, it can't be changed.
Kafka guarantees each message is saved exactly once.
Brokers also guarantee to store messages in the exact order they are received. 
Kafka does this through idempotence and transactions.

# Consuming messages
Consume is a read of a message. Consuming a message does not delete the message.
A consumer connects to a broker and requests available messages for a topic.
Consumers have to pull messages off a topic. 
Consumers can't edit or update messages, only receive them.
A consumer can read messages from multiple topics. 
Consumers also have a `groupID` which allows them to share work between a group of consumers.

A consumer keeps an offset cursor which is used to track which messages it has read.
After reading a message, a consumer must advance it's cursor to the next offset position and continue.
It is the consumers responsibility to advance it's cursor and remember it's last read offset within a partition.
A consumer only tracks one offset per partition. 
When a consumer commits an offset for a partition, it only commits the value of the last message.
Doing this automatically causes all earlier messages to be committed to it's offset.

Kafka brokers track the offsets with an internal topic (that consumers can't write to) called `__consumer_offsets`
This means that if a consumer group goes down, when it comes back, it can request its offset from kafka (read only) and pick up where it left off.

Multiple consumers outside the same group can consume a single partition tracking their own individual offset.
Multiple consumers inside the same group will work together and share an offset.
All this guarantees that each message in a topic partition is read exactly once.
