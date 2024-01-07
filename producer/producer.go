package main

import (
	"crypto/tls"
	"crypto/x509"
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
type SSLConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CAFile   string `yaml:"caFile"`
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

type KafkaConfig struct {
	BrokerList string    `yaml:"brokerList"`
	Topic      string    `yaml:"topic"`
	SSL        SSLConfig `yaml:"ssl"`
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

func main() {
	// Parse command-line arguments
	topicPtr := flag.String("topic", "", "Kafka topic to produce to")
	messagePtr := flag.String("message", "Hello, Kafka!", "Message to produce to Kafka")
	flag.Parse()

	// Check if the topic flag is provided
	if *topicPtr == "" {
		log.Fatal("Please provide a Kafka topic using the -topic flag")
	}

	// Load Kafka configuration from the YAML file
	config, err := loadConfig("producer_config.yaml")
	if err != nil {
		log.Fatalf("Error loading Kafka configuration: %v", err)
	}

	// Update the Kafka topic from the command-line argument
	config.Topic = *topicPtr

	// Configure the Kafka producer
	configProducer := sarama.NewConfig()
	configProducer.Producer.RequiredAcks = sarama.WaitForAll
	configProducer.Producer.Retry.Max = 5
	configProducer.Producer.Return.Successes = true

	// Configure the Kafka producer with SSL
	if config.SSL.Enabled {
		configProducer.Net.TLS.Enable = true
		configProducer.Net.TLS.Config, err = newTLSConfig(config.SSL.CAFile, config.SSL.CertFile, config.SSL.KeyFile)
		if err != nil {
			log.Fatalf("Error configuring TLS: %v", err)
		}
	}

	// Create a new Kafka producer
	producer, err := sarama.NewSyncProducer(strings.Split(config.BrokerList, ","), configProducer)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalf("Error closing Kafka producer: %v", err)
		}
	}()

	// Produce a message to the Kafka topic
	message := &sarama.ProducerMessage{
		Topic: config.Topic,
		Value: sarama.StringEncoder(*messagePtr),
	}

	// Send the message to the Kafka topic
	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatalf("Error sending message to Kafka: %v", err)
	}

	log.Printf("Message sent successfully to topic %s, partition %d, offset %d\n", config.Topic, partition, offset)

	// Handle shutdown gracefully
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	select {
	case <-sigchan:
		log.Println("Interrupt signal received. Shutting down...")
	}
}
