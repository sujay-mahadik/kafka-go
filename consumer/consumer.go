package main

import (
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/IBM/sarama"
)

const (
	brokerList = "localhost:9092"    // Replace with your Kafka broker address
	topic      = "quickstart-events" // Replace with your Kafka topic
)

func main() {
	// Configure the Kafka consumer
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create a new Kafka consumer
	consumer, err := sarama.NewConsumer(strings.Split(brokerList, ","), config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("Error closing Kafka consumer: %v", err)
		}
	}()

	// Subscribe to the Kafka topic
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error subscribing to Kafka topic: %v", err)
	}

	// Handle messages in a separate goroutine
	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				log.Printf("Received message: Topic=%s, Partition=%d, Offset=%d, Key=%s, Value=%s\n",
					msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			case err := <-partitionConsumer.Errors():
				log.Printf("Error: %v\n", err)
			}
		}
	}()

	// Handle shutdown gracefully
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	select {
	case <-sigchan:
		log.Println("Interrupt signal received. Shutting down...")
	}
}
