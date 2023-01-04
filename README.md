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

# Go
To interact with Kafka there is [confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go) which appears to be 
a cgo wrapper around [librdkafka](https://github.com/confluentinc/librdkafka) but I decided to use [kafka-go](https://github.com/segmentio/kafka-go)
as it is written entirely in Go. 

When using kafka-go, only the leader broker create topics. I wrote some code in `kafka.go` to create a topic. 
I can verify the topic is created with `docker logs broker` and I can also see it appear on the web ui too :heart_eyes:
