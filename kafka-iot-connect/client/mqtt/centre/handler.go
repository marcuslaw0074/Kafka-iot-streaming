package centre

import (
	"context"
	"fmt"
	mq "kafka-iot-connect/client/mqtt"
	"kafka-iot-connect/client/redis"
	logging "kafka-iot-connect/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func EtlMessageHandlerWithClient(ctx context.Context, r *redis.RedisTimeSeriesClient) MqttMessageHandler {
	return func(amc *mq.AllMessageChannels, l *logging.MyFileLogger) mqtt.MessageHandler {
		return func(client mqtt.Client, msg mqtt.Message) {
			l.LogLim(logging.LogInfo, pkgName, fmt.Sprintf("Received message: %s from topic: %s", msg.Payload(), msg.Topic()), 100)
			for _, ele := range amc.Chs {
				for _, el := range ele.Topics {
					if el.MatchTopics(msg.Topic()) {
						// go r.XaddTrim(ctx, fmt.Sprintf("%s-streamId-%d", r.SecretKey, el.Id), fmt.Sprintf("%d", el.Size), "*", "value", string(msg.Payload()))
						ele.Ch <- msg.Payload()
					}
				}
			}
		}
	}
}
