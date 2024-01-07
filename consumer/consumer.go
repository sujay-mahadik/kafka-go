package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"

	"io/ioutil"

	"github.com/IBM/sarama"
	"gopkg.in/yaml.v2"
)

// KafkaConfig represents the Kafka configuration
type KafkaConfig struct {
	BrokerList string `yaml:"brokerList"`
}

func loadConfig(filename string) (*KafkaConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config KafkaConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	// Parse command-line arguments
	topicPtr := flag.String("topic", "", "Kafka topic to subscribe to")
	flag.Parse()

	// Check if the topic flag is provided
	if *topicPtr == "" {
		log.Fatal("Please provide a Kafka topic using the -topic flag")
	}

	// Load Kafka configuration from the YAML file
	config, err := loadConfig("consumer_config.yaml")
	if err != nil {
		log.Fatalf("Error loading Kafka configuration: %v", err)
	}

	// Configure the Kafka consumer
	configConsumer := sarama.NewConfig()
	configConsumer.Consumer.Return.Errors = true

	// Create a new Kafka consumer
	consumer, err := sarama.NewConsumer(strings.Split(config.BrokerList, ","), configConsumer)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("Error closing Kafka consumer: %v", err)
		}
	}()

	// Subscribe to the specified Kafka topic
	partitionConsumer, err := consumer.ConsumePartition(*topicPtr, 0, sarama.OffsetNewest)
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
