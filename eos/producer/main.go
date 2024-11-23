package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type TransactionalMessage struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func generateTransactionalID() string {
	hostName, _ := os.Hostname()
	time := time.Now().UnixNano()

	return fmt.Sprintf("%s-%d", hostName, time)
}

// what is acks?
// idempotence true
// - Purpose: This setting enables idempotence for the producer,
//  which is essential for achieving exactly-once semantics.
//  With idempotence enabled, Kafka guarantees that even if a message
// is retried due to network issues or failures,
// it will only be written once to Kafka, ensuring no duplicate messages.

func main() {
	transactionalID := generateTransactionalID()

	// kafka producer config with exactly once semantics
	config := &kafka.ConfigMap{
		"bootstrap.servers":                     "localhost:9092,localhost:9093,localhost:9094",
		"acks":                                  "all",
		"enable.idempotence":                    true,            // enable idempotence for exactly once semantics
		"transactional.id":                      transactionalID, // unique transactional ID for producer
		"max.in.flight.requests.per.connection": 1,               // Ensures message order
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatalf("enable to start producer for exactly once semantics err: %+v", err)
	}
	defer producer.Close()

	// initialize kafka transaction
	err = producer.InitTransactions(context.TODO())
	if err != nil {
		log.Fatalf("Failed to initialize transactions: %s", err)
	}

	// begin kafka transaction
	err = producer.BeginTransaction()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %s", err)
	}

	message := TransactionalMessage{ID: "id_001", Content: "hello kafka with exactly once semantics"}

	byteMessage, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("failed to marshal to bytes struct, err: %s", err)
	}

	topic := "exactly-once-topic"
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: byteMessage},
		nil)
	if err != nil {
		log.Fatalf("failed to produce message, err: %s", err.Error())
	}

	err = producer.CommitTransaction(context.TODO())
	if err != nil {
		log.Fatalf("failed to commit transaction, err: %s", err.Error())
	}

	log.Println("Successfully committed transaction kafka")
}
