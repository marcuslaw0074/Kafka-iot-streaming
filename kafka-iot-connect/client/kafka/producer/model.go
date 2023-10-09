package producer

import (
	"sync"

	"github.com/IBM/sarama"
	"github.com/rcrowley/go-metrics"
)

type producerProvider struct {
	transactionIdGenerator int32
	producersLock          sync.Mutex
	producers              []sarama.AsyncProducer

	producerProvider func() sarama.AsyncProducer
}

type KafkaConnectProducer struct {
	brokers     string
	version     string
	verbose     bool
	keepRunning bool
	producers   int
	uuid        string
	producer    *producerProvider

	recordsRate metrics.Meter
}

