package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	//"io"
	"io/ioutil"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"kafkavents"
	"os"	
	//"strings"
)

// GetProducer gets/makes a Kafka Producer
func GetProducer(kafkaconf kafkavents.KafkaConf) kafka.Producer {
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

	return *producer
}

func main() {
	var kafkaConfigPath = flag.String("k", "examples/kafka_conf.json", "kafka config file")
	var junitXMLFilePath = flag.String("f", "examples/junit.xml", "junit xml file")
	var kafkaTopic = flag.String("t", "kafkavents", "the Kafka topic")
	flag.Parse()

	kafkaConfFile, _ := os.Open(*kafkaConfigPath)
	kafkaByteValue, _ := ioutil.ReadAll(kafkaConfFile)

	kafkaconf := kafkavents.KafkaConf{}
	kafkaConfFile.Close()

	json.Unmarshal(kafkaByteValue, &kafkaconf)
	fmt.Println(kafkaconf.BootstrapServers)

	kafkaconf.KVProducer()
	producer := kafkaconf.KProducer

	junitXMLFile, _ := os.Open(*junitXMLFilePath)
	byteValue, _ := ioutil.ReadAll(junitXMLFile)
	var testsuites kafkavents.Testsuites
	xml.Unmarshal(byteValue, &testsuites)
	fmt.Println(len(testsuites.Testsuites))

	jsontext, err := json.MarshalIndent(testsuites, "", "  ")
	if err != nil {
		fmt.Println("ERROR on JSON Marshal")
	}
	junitXMLFile.Close()

	fmt.Println(string(jsontext))

	//kafkaconf.Send(*kafkaTopic, json)

	for _, testsuite := range testsuites.Testsuites {
		fmt.Printf("%s", testsuite)
		for _, testcase := range testsuite.Testcases {
			message, err := json.MarshalIndent(testcase, "", "  ")
			if err != nil {
				fmt.Println("ERROR on JSON Marshal")
			}
			fmt.Printf("%s", message)
			kafkaconf.Send(*kafkaTopic, []byte(message))
		}
	}

	fmt.Printf("Closing producer")
	producer.Flush(5000)
	producer.Close()
	fmt.Printf("producer closed")
}
