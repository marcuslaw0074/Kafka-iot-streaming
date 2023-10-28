package producer

import (
	"fmt"
	"kafka-iot-connect/tool"
	"log"
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/gofrs/uuid"
	"github.com/rcrowley/go-metrics"
)

func (p *producerProvider) borrow() (producer sarama.AsyncProducer) {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	if len(p.producers) == 0 {
		for {
			producer = p.producerProvider()
			if producer != nil {
				return
			}
		}
	}

	index := len(p.producers) - 1
	producer = p.producers[index]
	p.producers = p.producers[:index]
	return
}

func (p *producerProvider) release(producer sarama.AsyncProducer) {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	// If released producer is erroneous close it and don't return it to the producer pool.
	if producer.TxnStatus()&sarama.ProducerTxnFlagInError != 0 {
		// Try to close it
		_ = producer.Close()
		return
	}
	p.producers = append(p.producers, producer)
}

func (p *producerProvider) clear() {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	for _, producer := range p.producers {
		producer.Close()
	}
	p.producers = p.producers[:0]
}

func (p *producerProvider) ProduceRecord(messages ...*sarama.ProducerMessage) {
	producer := p.borrow()
	defer p.release(producer)

	err := producer.BeginTxn()
	if err != nil {
		log.Printf("unable to start txn %s\n", err)
		return
	}

	for _, ele := range messages {
		if ele != nil {
			producer.Input() <- ele
		}
	}

	err = producer.CommitTxn()
	if err != nil {
		log.Printf("Producer: unable to commit txn %s\n", err)
		for {
			if producer.TxnStatus()&sarama.ProducerTxnFlagFatalError != 0 {
				// fatal error. need to recreate producer.
				log.Printf("Producer: producer is in a fatal state, need to recreate it")
				break
			}
			// If producer is in abortable state, try to abort current transaction.
			if producer.TxnStatus()&sarama.ProducerTxnFlagAbortableError != 0 {
				err = producer.AbortTxn()
				if err != nil {
					// If an error occured just retry it.
					log.Printf("Producer: unable to abort transaction: %+v", err)
					continue
				}
				break
			}
			// if not you can retry
			err = producer.CommitTxn()
			if err != nil {
				log.Printf("Producer: unable to commit txn %s\n", err)
				continue
			}
		}
		return
	}
}

func newProducerProvider(brokers []string, producerConfigurationProvider func() *sarama.Config) *producerProvider {
	provider := &producerProvider{}
	provider.producerProvider = func() sarama.AsyncProducer {
		config := producerConfigurationProvider()
		suffix := provider.transactionIdGenerator
		// Append transactionIdGenerator to current config.Producer.Transaction.ID to ensure transaction-id uniqueness.
		if config.Producer.Transaction.ID != "" {
			provider.transactionIdGenerator++
			config.Producer.Transaction.ID = config.Producer.Transaction.ID + "-" + fmt.Sprint(suffix)
		}
		producer, err := sarama.NewAsyncProducer(brokers, config)
		if err != nil {
			return nil
		}
		return producer
	}
	return provider
}

func InitializeKafkaConnectProducer(brokers, version string, verbose bool, producers int) *KafkaConnectProducer {
	if brokers == "" {
		panic("no Kafka bootstrap brokers defined, desc: Kafka bootstrap brokers to connect to, as a comma separated list")
	}
	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	sarama_ver, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	_uid := uuid.Must(uuid.NewV4()).String()

	err = tool.WriteFile("./.kafka_config", []byte(_uid))
	if err != nil {
		log.Panicf("Error writing Kafka uuid: %v", err)
	}

	return &KafkaConnectProducer{
		brokers:     brokers,
		version:     version,
		verbose:     verbose,
		uuid:        _uid,
		keepRunning: true,
		producers:   producers,
		producer: newProducerProvider(strings.Split(brokers, ","), func() *sarama.Config {
			config := sarama.NewConfig()
			config.Version = sarama_ver
			config.Producer.Idempotent = true
			config.Producer.Return.Errors = false
			config.Producer.RequiredAcks = sarama.WaitForAll
			config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
			config.Producer.Transaction.Retry.Backoff = 10
			config.Producer.Transaction.ID = "txn_producer_" + _uid
			config.Net.MaxOpenRequests = 1
			return config
		}),
		recordsRate: metrics.GetOrRegisterMeter("records.rate", nil),
	}
}

func (k *KafkaConnectProducer) Close() {
	k.producer.clear()
}
func (k *KafkaConnectProducer) Produce(messages ...*sarama.ProducerMessage) {
	k.producer.ProduceRecord(messages...)
}