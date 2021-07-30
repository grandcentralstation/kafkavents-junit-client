package kafkavents

// KafkaConf Kafka config definition
type KafkaConf struct {
	BootstrapServers string		`json:"bootstrap.servers"`
	GroupID string    	`json:"group.id"`
    SaslMechanism string	`json:"sasl.mechanism"`
    SaslSecurityProtocol string	`json:"security.protocol"`
    SaslUsername string `json:"sasl.username"`
    SaslPassword string `json:"sasl.password"`
}
