# kafka-go

This repository provides a simple example of a Kafka producer and consumer implemented in Go, using the Sarama library. The producer sends messages to a specified Kafka topic, and the consumer subscribes to the same topic to receive and log the messages.

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
```
- brokerList: The comma-separated list of Kafka broker addresses.

### Consumer Configuration (consumer_config.yaml)
    
```
brokerList: "localhost:9092"
```
- brokerList: The comma-separated list of Kafka broker addresses.


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
