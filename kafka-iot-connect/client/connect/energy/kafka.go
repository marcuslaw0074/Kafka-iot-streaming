package energy

import (
	"context"
	"fmt"
	"kafka-iot-connect/client/kafka/consumer"
	"kafka-iot-connect/client/mqtt"
	"kafka-iot-connect/client/mqtt/centre"
	logging "kafka-iot-connect/log"
	"kafka-iot-connect/tool"

	"github.com/IBM/sarama"
	mq "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofrs/uuid"
)


func EnergyᚖMillsᚋKafkaᚋClient(ctx context.Context) {
	log := logging.StartLogger("./log/api.log", 100000)
	log.ClearLog()
	uid := uuid.Must(uuid.NewV4())
	mqConf, err := tool.ParseJsonFile[mqtt.MqttConf]("./mqtt.json")
	if err != nil {
		panic(err)
	}
	conf := mqtt.MqttConfig{
		Host:                  mqConf.Host,
		Port:                  mqConf.Port,
		ClientId:              "kafka-connect-client--" + uid.String(),
		Scheme:                mqtt.GetScheme(mqConf.TlsPath),
		TlsPath:               mqConf.TlsPath,
		Log:                   log,
		UseDefaultHandlers:    false,
		MessageHandler:        centre.EtlMessageHandlerWithClient(ctx, nil),
		OnConnectHandler:      centre.OnConnectHandler,
		ConnectionLostHandler: centre.ConnectionLostHandler,
		AllCh: &mqtt.AllMessageChannels{
			Id: "01",
			Chs: []mqtt.MessageChannel{
				{
					Ch: make(chan []byte, 10),
					Id: "test-stream-02",
					Topics: []mqtt.Topic{
						{
							Name: "etl/energy/bms/marcus/af94c736-8da1-4ab9-f5ee-08db213be062/8",
						},
					},
					IsActive: true,
					CallbackAction: func(c mq.Client, b [][]byte, t []mqtt.Topic, i int) {
						fmt.Printf("received messaged from topic: %s, message: %s\n", t[0].Name, b[0])
					},
				},
			},
		},
	}
	if err := conf.DefaultClient(); err != nil {
		panic(err)
	} else {
		go conf.Subscribe(log)
		go conf.OnKaflaChannelHandler("test-stream-02", nil)
		kc := consumer.InitializeKafkaConnectConsumer("localhost:9092", "3.3.2", "example", "range", true, true)
		kc.Consume(context.Background(), "kafka-streams-kWh-raw-data", func(cm *sarama.ConsumerMessage) error {
			conf.Publish(string(cm.Value), log, "etl/energy/bms/marcus/af94c736-8da1-4ab9-f5ee-08db213be062/8")
			return nil
		})
	}
}