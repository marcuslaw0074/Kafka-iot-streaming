package mqtt

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (m *AllMessageChannels) DefaultMessageChannelHandler() mqtt.MessageHandler {
	var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		for _, ele := range m.Chs {
			for _, el := range ele.Topics {
				if msg.Topic() == el.Name {
					ele.Ch <- msg.Payload()
				}
			}
		}
	}
	return messagePubHandler
}

func (m *AllMessageChannels) DefaultChannelConnectHandler() mqtt.OnConnectHandler {
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		for _, ele := range m.Chs {
			log.Printf("Ch id: %s connected\n", ele.Id)
		}
	}
	return connectHandler
}

func (m *AllMessageChannels) DefaultChannelConnectLostHandler() mqtt.ConnectionLostHandler {
	var connectHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Printf("Lost connection: %s\n", err.Error())
	}
	return connectHandler
}