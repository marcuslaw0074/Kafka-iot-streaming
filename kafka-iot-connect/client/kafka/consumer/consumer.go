package consumer

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/gofrs/uuid"
)

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (err error) {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			log.Printf("Message claimed: key = %s, value = %s, timestamp = %v, topic = %s", string(message.Key), string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
			err = consumer.callback(message)
			if err != nil {
				return
			}
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return
		}
	}
}

func InitializeKafkaConnectConsumer(brokers, version, group, assignor string, verbose, oldest bool) *KafkaConnectConsumer {
	if brokers == "" {
		panic("no Kafka bootstrap brokers defined, desc: Kafka bootstrap brokers to connect to, as a comma separated list")
	}
	if group == "" {
		panic("no Kafka consumer group defined, desc: Kafka consumer group definition")
	}

	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	sarama_ver, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	config := sarama.NewConfig()
	config.Version = sarama_ver

	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	}

	if oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	_uid := uuid.Must(uuid.NewV4()).String()

	return &KafkaConnectConsumer{
		brokers:     brokers,
		version:     version,
		group:       group,
		assignor:    assignor,
		oldest:      oldest,
		verbose:     verbose,
		keepRunning: true,
		uuid:        _uid,

		consumer: client,
	}
}

func (k *KafkaConnectConsumer) Consume(ctx context.Context, topics string, callback func(*sarama.ConsumerMessage) error) error {

	consumer := Consumer{
		ready: make(chan bool),
		callback: callback,
	}
	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := k.consumer.Consume(ctx, strings.Split(topics, ","), &consumer); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return err
			}
			log.Printf("Error from consumer: %v\n", err)
			return err
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return nil
		}
		consumer.ready = make(chan bool)
	}
}
