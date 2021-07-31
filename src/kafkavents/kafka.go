package kafkavents

import (
    "fmt"
    "github.com/confluentinc/confluent-kafka-go/kafka"
)


// KafkaConf Kafka config definition
type KafkaConf struct {
	BootstrapServers string		`json:"bootstrap.servers"`
	GroupID string    	`json:"group.id"`
    SaslMechanism string`json:"sasl.mechanism"`
    SaslSecurityProtocol string	`json:"security.protocol"`
    SaslUsername string `json:"sasl.username"`
    SaslPassword string `json:"sasl.password"`
    KProducer *kafka.Producer
}

// KafkaVents definition
type KafkaVents struct {
    Config KafkaConf
    Producer *kafka.Producer
}

// KVEvent defines a KafkaVents event
type KVEvent struct {
    Header KVHeader
    Body KVBody
}

// KVHeader defines the header for a KafkaVent
type KVHeader struct {
    SessionID  string `json:"session_id"`
    SessionNum int64    `json:"session_num"`
    Topic string        `json:"topic"`
    Packetnum int       `json:"packetnum"`
    Type string         `json:"type"`
    Source string       `json:"source"`
    Version string      `json:"version"`
    Timestamp string    `json:"timestamp"`
}

// KVBody defines a body for a KafkaVent
type KVBody struct {
    Nodeid string   `json:"nodeid"`
    Domain string   `json:"domain"`
    Nodepath string `json:"nodepath"`
    Status string   `json:"status"`
    Duration float64    `json:"duration"`
    Stdout string       `json:"stdout"`
    Stderr string       `json:"stderr"`
    Message string      `json:"message"`
}

// Send a kafka event
func (kc *KafkaConf) Send(topic string, event []byte) {
    kc.KProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic,
			Partition: kafka.PartitionAny},
		Value: event}, nil)

	// Wait for delivery report
	e := <-kc.KProducer.Events()

	message := e.(*kafka.Message)
	if message.TopicPartition.Error != nil {
		fmt.Printf("failed to deliver message: %v\n",
			message.TopicPartition)
	} else {
		fmt.Printf("delivered to topic %s [%d] at offset %v\n",
			*message.TopicPartition.Topic,
			message.TopicPartition.Partition,
			message.TopicPartition.Offset)
	}
}

// KVProducer gets/makes a Kafka Producer
func (kc *KafkaConf) KVProducer() {
	// Produce a new record to the topic...
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":       kc.BootstrapServers,
		"sasl.mechanisms":         kc.SaslMechanism,
		"security.protocol":       kc.SaslSecurityProtocol,
		"sasl.username":           kc.SaslUsername,
		"sasl.password":           kc.SaslPassword})

	if err != nil {
		panic(fmt.Sprintf("Failed to create producer: %s", err))
	}
    kc.KProducer = producer
}
