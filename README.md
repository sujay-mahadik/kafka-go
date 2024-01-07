# kafka-go

This repository provides a simple example of a Kafka producer and consumer implemented in Go, using the Sarama library. The producer sends messages to a specified Kafka topic, and the consumer subscribes to the same topic to receive and log the messages.

![kafka-go-snippet](https://github.com/sujay-mahadik/kafka-go/assets/20441076/beff8e83-e641-42d6-b8b6-6cebff09120a)

## Requirements

- Go installed on your machine.
- Apache Kafka running locally or on a reachable server.
- Sarama library (`github.com/IBM/sarama`) installed. 
You can install it using the following command:

```bash
go get github.com/IBM/sarama
```

## Configuration
Both the producer and consumer rely on configuration files in YAML format. The YAML files include essential Kafka configurations such as broker list. This can be extended to add other configurations like message encoding, ssl and other kafka configurations

### Producer Configuration (producer_config.yaml)
    
```
brokerList: "localhost:9092"

# SSL Configuration
ssl:
  enabled: true
  caFile: "path/to/ca.pem"
  certFile: "path/to/client.pem"
  keyFile: "path/to/client.key"
```
- brokerList: The comma-separated list of Kafka broker addresses.
- ssl.enabled: Set to true to enable SSL.
- ssl.caFile: Path to the Certificate Authority (CA) file.
- ssl.certFile: Path to the client certificate file.
- ssl.keyFile: Path to the client private key file.

### Consumer Configuration (consumer_config.yaml)
    
```
brokerList: "localhost:9092"

# SSL Configuration
ssl:
  enabled: true
  caFile: "path/to/ca.pem"
  certFile: "path/to/client.pem"
  keyFile: "path/to/client.key"

# Consumer Group Configuration
consumerGroup: "your_consumer_group"
```
- brokerList: The comma-separated list of Kafka broker addresses.
- ssl.enabled: Set to true to enable SSL.
- ssl.caFile: Path to the Certificate Authority (CA) file.
- ssl.certFile: Path to the client certificate file.
- ssl.keyFile: Path to the client private key file.
- consumerGroup: Consumer Group name to commit offsets.

## Build Steps

Best part about golang is that the build can be easily done for required OS and Arch

- Build the Kafka Producer
    ```bash
    GOOS={OS_NAME} GOARCH={OS_ARCH} go build -o producer.exe producer.go

    ```

- Build the Kafka Consumer:
    ```bash
    GOOS={OS_NAME} GOARCH={OS_ARCH} go build -o consumer.exe consumer.go
    ```

- Build Parameters

    |Parameter  | Linux     | Windows   |
    |-----------| --------  | -------   |
    |GOOS       | linux     | windows   |
    |GOARCH     | amd64     | amd64     |


## Running the Kafka Producer
To run the Kafka producer, use the following command:

```bash
./producer -topic topic_name -message "Your message here"
```
- -topic: The topic name to sent message to
- -message: The message to be sent to the Kafka topic.

## Running the Kafka Consumer

To run the Kafka consumer, use the following command:
```bash
./consumer -topic topic_name
```
- -topic: The topic name to sent message to

## Help

To get help run


```bash

./producer -help
Usage of producer:
  -message string
        Message to produce to Kafka (default "Hello, Kafka!")
  -topic string
        Kafka topic to produce to

./consumer -help
Usage of consumer:
  -topic string
        Kafka topic to subscribe to
```
