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
	// Configure the Kafka producer
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	// Create a new Kafka producer
	producer, err := sarama.NewSyncProducer(strings.Split(brokerList, ","), config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalf("Error closing Kafka producer: %v", err)
		}
	}()

	// Create a Kafka message
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder("Hello, Kafka!"),
	}

	// Send the message to the Kafka topic
	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatalf("Error sending message to Kafka: %v", err)
	}

	log.Printf("Message sent successfully to topic %s, partition %d, offset %d\n", topic, partition, offset)

	// Handle shutdown gracefully
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	select {
	case <-sigchan:
		log.Println("Interrupt signal received. Shutting down...")
	}
}
