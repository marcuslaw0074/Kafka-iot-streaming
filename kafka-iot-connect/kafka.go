package main

import (
	"context"
	"kafka-iot-connect/client/kafka/consumer"
	"kafka-iot-connect/client/kafka/producer"
	"kafka-iot-connect/model"

	"github.com/IBM/sarama"
)

func initializeKafkaConnect() {
	kp := producer.InitializeKafkaConnectProducer("localhost:9092", "3.3.2", true, 1)
	kp.Produce(
		&sarama.ProducerMessage{Topic: "kafka-streams-energy-raw-data", Key: model.BMSDataTypeEncoder("user1", "elec"), Value: model.BMSRawDataEncoder("user1", 1, 839.543)},
		&sarama.ProducerMessage{Topic: "kafka-streams-energy-raw-data", Key: model.BMSDataTypeEncoder("user1", "elec"), Value: model.BMSRawDataEncoder("user1", 1, 869.543)},
	)

	kc := consumer.InitializeKafkaConnectConsumer("localhost:9092", "3.3.2", "example", "range", true, true)
	kc.Consume(context.Background(), "kafka-streams-energy-realtime-data", func(cm *sarama.ConsumerMessage) error {
		return nil
	})
}
