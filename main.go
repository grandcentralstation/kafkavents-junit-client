package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	//"io"
	"io/ioutil"
	"kafkavents"
	"os"	
	//"strings"
)


func main() {
	var kafkaConfigPath = flag.String("k", "examples/kafka_conf.json", "kafka config file")
	var junitXMLFilePath = flag.String("f", "examples/junit.xml", "junit xml file")
	flag.Parse()
	junitXMLFile, _ := os.Open(*junitXMLFilePath)
	kafkaConfFile, _ := os.Open(*kafkaConfigPath)
	kafkaByteValue, _ := ioutil.ReadAll(kafkaConfFile)
	kafkaConfFile.Close()

	var kafkaconf kafkavents.KafkaConf

	json.Unmarshal(kafkaByteValue, &kafkaconf)
	fmt.Println(kafkaconf.BootstrapServers)

	topic := "kafkavents"
	// Produce a new record to the topic...
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":       kafkaconf.BootstrapServers,
		"sasl.mechanisms":         kafkaconf.SaslMechanism,
		"security.protocol":       kafkaconf.SaslSecurityProtocol,
		"sasl.username":           kafkaconf.SaslUsername,
		"sasl.password":           kafkaconf.SaslPassword})

	if err != nil {
		panic(fmt.Sprintf("Failed to create producer: %s", err))
	}
	
	byteValue, _ := ioutil.ReadAll(junitXMLFile)

	// we initialize our Users array
	var testsuites kafkavents.Testsuites
	xml.Unmarshal(byteValue, &testsuites)
	fmt.Println(len(testsuites.Testsuites))
	fmt.Println("Testsuite: " + testsuites.Testsuites[0].Name)

	json, err := json.MarshalIndent(testsuites, "", "  ")
	if err != nil {
		fmt.Println("ERROR on JSON Marshall")
	}
	junitXMLFile.Close()

	fmt.Println(string(json))

	//value := json
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic,
			Partition: kafka.PartitionAny},
		Value: []byte(json)}, nil)

	// Wait for delivery report
	e := <-producer.Events()

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


	producer.Close()
}
