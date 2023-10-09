package consumer

import (
	"github.com/IBM/sarama"
)

type KafkaConnectConsumer struct {
	brokers     string
	version     string
	group       string
	assignor    string
	oldest      bool
	verbose     bool
	keepRunning bool
	uuid        string

	consumer sarama.ConsumerGroup
}

type Consumer struct {
	ready    chan bool
	callback func(s *sarama.ConsumerMessage) (err error)
}

