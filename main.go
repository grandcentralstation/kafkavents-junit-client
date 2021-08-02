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
	"github.com/google/uuid"
	"strings"
	"time"
)

// GetProducer gets/makes a Kafka Producer
func GetProducer(config kafkavents.KafkaConf) kafka.Producer {
	// Produce a new record to the topic...
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":       config.BootstrapServers,
		"sasl.mechanisms":         config.SaslMechanism,
		"security.protocol":       config.SaslSecurityProtocol,
		"sasl.username":           config.SaslUsername,
		"sasl.password":           config.SaslPassword})

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

	SessionID := uuid.NewString()
	fmt.Printf("SessionID: %s\n", SessionID)
	var SessionNum int64 = 0
	var PacketNum int64 = 0

	kafkaConfFile, _ := os.Open(*kafkaConfigPath)
	kafkaByteValue, _ := ioutil.ReadAll(kafkaConfFile)

	kafkaconf := kafkavents.KafkaConf{}
	kafkaConfFile.Close()

	json.Unmarshal(kafkaByteValue, &kafkaconf)
	//fmt.Println(kafkaconf.BootstrapServers)

	kafkaconf.KVProducer()
	producer := kafkaconf.KProducer

	junitXMLFile, _ := os.Open(*junitXMLFilePath)
	byteValue, _ := ioutil.ReadAll(junitXMLFile)
	var testsuites kafkavents.Testsuites
	xml.Unmarshal(byteValue, &testsuites)
	//fmt.Println(len(testsuites.Testsuites))

	/*
	jsontext, err := json.MarshalIndent(testsuites, "", "  ")
	if err != nil {
		fmt.Println("ERROR on JSON Marshal")
	}
	junitXMLFile.Close()

	//fmt.Println(string(jsontext))

	kafkaconf.Send(*kafkaTopic, jsontext)
	*/

	// refactor into function w/ start struct
	kventStart := kafkavents.KVEvent{}
	kventStart.Header.Topic = *kafkaTopic
	kventStart.Header.SessionID = SessionID
	kventStart.Header.Type = "sessionstart"
	kventStart.Header.Source = "junit-kafkavents"
	kventStart.Header.Version = "0.01"
	kventStart.Header.Timestamp = time.Now()
	kventStart.Header.SessionNum = SessionNum
	kventStart.Header.Packetnum = PacketNum
	kventStart.Body.Name = "My Test"
	//kventStartMessage, err := json.MarshalIndent(kventStart, "", "  ")
	kventStartMessage, err := json.Marshal(kventStart)
	if err != nil {
		fmt.Println("ERROR on JSON Marshal")
	}
	kafkaconf.Send(*kafkaTopic, []byte(kventStartMessage))

	kvent := kafkavents.KVEvent{}
	for _, testsuite := range testsuites.Testsuites {
		//fmt.Printf("%s", testsuite)
		for _, testcase := range testsuite.Testcases {
			kvent.Header.Topic = *kafkaTopic
			kvent.Header.SessionID = SessionID
			kvent.Header.Type = "testresult"
			kvent.Header.Source = "junit-kafkavents"
			kvent.Header.Version = "0.01"
			kvent.Header.Timestamp = time.Now()
			kvent.Header.SessionNum = SessionNum
			kvent.Header.Packetnum = PacketNum
			// build nodepath from Classname and Name
			tempName := strings.Replace(testcase.Name, "[", ".", 1)
			tempName = strings.Replace(tempName, "]", "", 1)
			kvent.Body.Nodepath = testcase.Classname + "." + tempName
			kvent.Body.Domain = kvent.Body.Nodepath
			status := "passed"
			if testcase.Skipped != nil {
				status = "skipped"
				kvent.Body.Message = testcase.Skipped.Message
				kvent.Body.Stderr = testcase.Skipped.Text
			}
			if testcase.Failure != nil {
				status = "failed"
				kvent.Body.Message = testcase.Failure.Message
				kvent.Body.Stderr = testcase.Failure.Text
			}
			kvent.Body.Status = status
			kvent.Body.Duration = testcase.Time
			kventMessage, err := json.MarshalIndent(kvent, "", "  ")
			//kventMessage, err := json.Marshal(kvent)
			if err != nil {
				fmt.Println("ERROR on JSON Marshal")
			}
			kafkaconf.Send(*kafkaTopic, []byte(kventMessage))
			// move this to send to keep track of actual sends
			PacketNum++
			SessionNum++
		}
	}

	// refactor into function w/ end struct
	kventEnd := kafkavents.KVEvent{}
	kventEnd.Header.Topic = *kafkaTopic
	kventEnd.Header.SessionID = SessionID
	kventEnd.Header.Type = "sessionend"
	kventEnd.Header.Source = "junit-kafkavents"
	kventEnd.Header.Version = "0.01"
	kventEnd.Header.Timestamp = time.Now()
	kventEnd.Header.SessionNum = SessionNum
	kventEnd.Header.Packetnum = PacketNum
	kventEnd.Body.Name = "My Test"
	//kventEndMessage, err := json.MarshalIndent(kventEnd, "", "  ")
	kventEndMessage, err := json.Marshal(kventEnd)
	if err != nil {
		fmt.Println("ERROR on JSON Marshal")
	}
	kafkaconf.Send(*kafkaTopic, []byte(kventEndMessage))

	fmt.Println("Closing producer")
	producer.Flush(5000)
	producer.Close()
	fmt.Println("producer closed")
}
