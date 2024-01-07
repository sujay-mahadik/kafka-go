package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"io/ioutil"

	"github.com/IBM/sarama"
	"gopkg.in/yaml.v2"
)

// KafkaConfig represents the Kafka configuration
type SSLConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CAFile   string `yaml:"caFile"`
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

type KafkaConfig struct {
	BrokerList    string    `yaml:"brokerList"`
	SSL           SSLConfig `yaml:"ssl"`
	ConsumerGroup string    `yaml:"consumerGroup"`
}

func newTLSConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}, nil
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

// ConsumerHandler is a simple implementation of sarama.ConsumerGroupHandler
type ConsumerHandler struct{}

// NewConsumerHandler creates a new consumer handler to process Kafka messages
func NewConsumerHandler() sarama.ConsumerGroupHandler {
	return &ConsumerHandler{}
}

// Setup is called once when the consumer group is started.
func (h *ConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	// Not implemented for this example
	return nil
}

// Cleanup is called once when the consumer group is terminated.
func (h *ConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	// Not implemented for this example
	return nil
}

// ConsumeClaim is called for each claim returned by the consumer group
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Process each message from the claim
	for message := range claim.Messages() {
		fmt.Printf("Received message: Topic=%s, Partition=%d, Offset=%d, Key=%s, Value=%s\n",
			message.Topic, message.Partition, message.Offset, message.Key, message.Value)

		// Mark the message as processed
		session.MarkMessage(message, "")
	}

	return nil
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

	// SSL config
	if config.SSL.Enabled {
		configConsumer.Net.TLS.Enable = true
		configConsumer.Net.TLS.Config, err = newTLSConfig(config.SSL.CAFile, config.SSL.CertFile, config.SSL.KeyFile)
		if err != nil {
			log.Fatalf("Error configuring TLS: %v", err)
		}
	}

	// Set up Kafka consumer
	consumer, err := sarama.NewConsumer(strings.Split(config.BrokerList, ","), configConsumer)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Create a new consumer group
	consumerGroup, err := sarama.NewConsumerGroup(strings.Split(config.BrokerList, ","), config.ConsumerGroup, configConsumer)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer group: %v", err)
	}
	defer consumerGroup.Close()

	ctx := context.Background()
	// Consume messages from Kafka topic
	go func() {
		for {
			topic := *topicPtr
			handler := NewConsumerHandler()

			// Consume messages from the specified topic
			err := consumerGroup.Consume(ctx, strings.Split(topic, ","), handler)
			if err != nil {
				log.Fatalf("Error consuming messages: %v", err)
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
