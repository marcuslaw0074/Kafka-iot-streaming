package centre

import (
	"context"
	"fmt"
	mq "kafka-iot-connect/client/mqtt"
	"kafka-iot-connect/client/redis"
	logging "kafka-iot-connect/log"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func EtlMessageHandlerWithClient(ctx context.Context, r *redis.RedisTimeSeriesClient) MqttMessageHandler {
	return func(amc *mq.AllMessageChannels, l *logging.MyFileLogger) mqtt.MessageHandler {
		return func(client mqtt.Client, msg mqtt.Message) {
			l.LogLim(logging.LogInfo, pkgName, fmt.Sprintf("Received message: %s from topic: %s", msg.Payload(), msg.Topic()), 100)
			for _, ele := range amc.Chs {
				for _, el := range ele.Topics {
					if el.MatchTopics(msg.Topic()) {
						if r != nil {
							go r.XaddTrim(ctx, fmt.Sprintf("%s-streamId-%d", r.SecretKey, el.Id), fmt.Sprintf("%d", el.Size), "*", "value", string(msg.Payload()))
						}
						ele.Ch <- msg.Payload()
					}
				}
			}
		}
	}
}

func OnConnectHandler(amc *mq.AllMessageChannels, l *logging.MyFileLogger) mqtt.OnConnectHandler {
    var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
        for _, ele := range amc.Chs {
            if l == nil {
                log.Printf("Channel id: %s connected\n", ele.Id)
            } else {
                l.Log(logging.LogInfo, pkgName, fmt.Sprintf("Channel id: %s connected", ele.Id))
            }
        }
    }
    return connectHandler
}

func ConnectionLostHandler(amc *mq.AllMessageChannels, l *logging.MyFileLogger) mqtt.ConnectionLostHandler {
    var connectHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
        if l == nil {
            log.Printf("Lost connection: %s\n", err.Error())
        } else {
            l.Log(logging.LogError, pkgName, fmt.Sprintf("Lost connection: %s", err.Error()))
            client.Connect().Wait()
        }
    }
    return connectHandler
}