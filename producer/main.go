package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092,localhost:9093,localhost:9094", // Kafka broker address
		"acks":              1,                                              // waits for leader ack, default behavior
		"retries":           3,                                              // in case of failures may lead to duplicates
	})
	if err != nil {
		log.Fatalf("Failed to created producer: %s", err)
	}
	defer producer.Close()

	topic := "atleast-one-semantic"
	message := "this is a single message"

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)
	if err != nil {
		log.Fatalf("Failed to produce message: %s", err)
	} else {
		log.Printf("Message produced: %s", message)
	}

	// manually wait for producer ack (delivery report)
	// this is akin to manually waiting for producer to commit
	events := <-producer.Events()
	m, ok := events.(*kafka.Message)

	if ok && m.TopicPartition.Error == nil {
		log.Printf("successfully sent message to topic : %s, partition: %d, offset: %d",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	} else {
		log.Printf("Failed to deliver message err: %s", m.TopicPartition.Error.Error())
	}

	// wait for any remaining messages to be delivered
	producer.Flush(15 * 1000) // 15 ms
}
