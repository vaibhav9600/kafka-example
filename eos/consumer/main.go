package main

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// TransactionalMessage represents the structure of a message in a transaction
type TransactionalMessage struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

// offset of each partition increases by 2

func processMessage(t TransactionalMessage) {
	log.Printf("Processing message, Transactional ID: %s, Content: %s", t.ID, t.Content)
}
func main() {
	config := &kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9092,localhost:9093,localhost:9094",
		"group.id":           "exactly-once-consumer-group",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,            // disable auto commits of offsets
		"isolation.level":    "read_committed", // Only read messages from committed transactions
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatalf("error while initializing consumer group, err: %+v", err)
	}
	defer consumer.Close()

	topic := "exactly-once-topic"
	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("failed to subscribe to topic: %s, err: %s", topic, err.Error())
	}

	log.Println("Consumer started, waiting for messages...")

	for {
		msg, err := consumer.ReadMessage(-1) // blocking call for read messages
		if err != nil {
			log.Printf("consumer error:%v", err)
			continue
		}

		var transactionalMessage TransactionalMessage
		err = json.Unmarshal(msg.Value, &transactionalMessage)
		if err != nil {
			log.Printf("Failed to deserialize message: %v", err)
			continue
		}

		processMessage(transactionalMessage)

		// Commit the message's offset after processing
		kafkaCommit, err := consumer.CommitMessage(msg)
		if err != nil {
			log.Printf("Failed to commit message: %v", err)
		}
		log.Printf("commit successful, offset: %d, partition: %d", kafkaCommit[0].Offset, kafkaCommit[0].Partition)
	}
}
