package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":       "localhost:9092,localhost:9093,localhost:9094", // kafka brokers address
		"group.id":                "at-least-once-group",                          // consumer group ids
		"auto.offset.reset":       "earliest",                                     // Start reading from the earliest offset
		"enable.auto.commit":      false,                                          // automatically commits offset
		"auto.commit.interval.ms": 5000,                                           // Commits offset every 5 seconds
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

	
}
