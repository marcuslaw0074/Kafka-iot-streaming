package main

import (
	"kafka-iot-connect/client/kafka/producer"
	"kafka-iot-connect/model"
	"kafka-iot-connect/tool"
	"time"

	"github.com/IBM/sarama"
)
/*
docker exec -it kafka-broker bash   

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic test-kafka-streams-energy-raw-data

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic test-kafka-streams-energy-realtime-data

kafka-console-producer \
  --topic test-kafka-streams-energy-raw-data \
  --bootstrap-server localhost:9092 \
  --property parse.key=true \
  --property key.separator=":::"

kafka-console-consumer --bootstrap-server localhost:9092 \
--topic test-kafka-streams-energy-raw-data \
--from-beginning \
--formatter kafka.tools.DefaultMessageFormatter \
--property print.timestamp=true \
--property print.key=true \
--property print.value=true

kafka-console-consumer --bootstrap-server localhost:9092 \
--topic test-kafka-streams-energy-realtime-data \
--from-beginning \
--formatter kafka.tools.DefaultMessageFormatter \
--property print.timestamp=true \
--property print.key=true \
--property print.value=true

*/
func main() {
	msgs, err := tool.ParseJsonFile[DataParser]("./data/data.json")
	if err != nil {
		panic(err)
	} else {
		kp := producer.InitializeKafkaConnectProducer("localhost:9092", "3.3.2", true, 1)
		for _, msg := range msgs.Data {
			kp.Produce(&sarama.ProducerMessage{
				Topic: "test-kafka-streams-energy-raw-data",
				Key:   model.BMSDataTypeEncoder("", "elec"),
				Value: msg,
			})
			time.Sleep(time.Duration(msgs.Interval) * time.Second)
		}
	}
}
