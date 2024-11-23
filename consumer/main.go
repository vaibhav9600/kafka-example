package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// each consumer group is assigned to different partitions, this happens dynamically, lets say we had one consumer group initially
// it will consume all the partitions but if new consumer comes into picture both will be divided different partitions to listen to
// and will not consume same message,
// for example: consumer 1 is assigned with partition 2 and the other consumer is assigned with partition 0 and partition 1

// how many consumers we can have for each consumer group? max number can be = number of partitions for a topic
// Each partition in a Kafka topic can be consumed by only one consumer in a consumer group at any given time.
// 

func main() {
	// Commits offsets every 5 seconds, meaning that if the consumer crashes, it might reprocess the last few messages.
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092,localhost:9093,localhost:9094", // kafka brokers address
		"group.id":          "at-least-once-group",                          // consumer group ids
		// If there’s no previous offset, it starts consuming from the earliest offset.
		"auto.offset.reset":       "earliest", // Start reading from the earliest offset
		"enable.auto.commit":      false,      // automatically commits offset
		"auto.commit.interval.ms": 5000,       // Commits offset every 5 seconds
	})
	if err != nil {
		log.Fatalf("unable to start consumer, %s", err)
	}
	defer consumer.Close()

	// subscribe to topic
	topic := "atleast-one-semantic"
	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Printf("failed subscription to topic, err: %s", err.Error())
	}

	log.Println("Consumer started, waiting for messages...")

	// poll for new messages
	for {
		msg, err := consumer.ReadMessage(-1) // blocking call to read message
		if err == nil {
			log.Printf("consumer message: %s, Partition: %d, Topic: %s, Offset: %d",
				string(msg.Value), msg.TopicPartition.Partition, *msg.TopicPartition.Topic, msg.TopicPartition.Offset)
		} else {
			log.Printf("failed to read message, err: %s", err.Error())
		}
	}
}
