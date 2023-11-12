package energy

import (
	"context"
	"encoding/json"
	"fmt"
	"kafka-iot-connect/client/kafka/producer"
	"kafka-iot-connect/client/mqtt"
	"kafka-iot-connect/client/mqtt/centre"
	logging "kafka-iot-connect/log"
	"kafka-iot-connect/model"
	"kafka-iot-connect/tool"

	"github.com/IBM/sarama"
	mq "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofrs/uuid"
)

func EnergyᚖMillsᚋMqttᚋClient(ctx context.Context) {
	log := logging.StartLogger("./log/api.log", 100000)
	uid := uuid.Must(uuid.NewV4())
	log.ClearLog()
	mqConf, err := tool.ParseJsonFile[mqtt.MqttConf]("./mqtt.json")
	if err != nil {
		panic(err)
	}
	conf := mqtt.MqttConfig{
		Host:                  mqConf.Host,
		Port:                  mqConf.Port,
		ClientId:              "mqtt-connect-client--" + uid.String(),
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
					Id: "test-stream-01",
					Topics: []mqtt.Topic{
						{
							Name: "N/ELEC/NFMILLS/D/0001/test",
						},
					},
					IsActive: true,
					CallbackAction: func(c mq.Client, b [][]byte, t []mqtt.Topic, i int) {
						fmt.Printf("received messaged from topic: %s", t[0].Name)
						kp := producer.InitializeKafkaConnectProducer("localhost:9092", "3.3.2", true, 1)
						for _, bb := range b {
							var data []TimeSeriesPoint
							if err := json.Unmarshal(bb, &data); err != nil {
								log.Log(logging.LogDebug, "main", err)
								continue
							} else {
								fmt.Println(string(bb))
								kp.Produce(
									tool.Map(data,
										func(d TimeSeriesPoint, i int) *sarama.ProducerMessage {
											return &sarama.ProducerMessage{
												Topic: "kafka-etl-energy-raw",
												Key:   model.BMSDataTypeEncoder(d.Id, "elec"),
												Value: model.BMSDeltaDataEncoder(d.Id, 1, d.Value),
											}
										},
									)...,
								)
								kp.Close()
							}
						}
					},
				},
			},
		},
	}
	if err := conf.DefaultClient(); err != nil {
		panic(err)
	} else {
		go conf.Subscribe(log)
		go conf.OnKaflaChannelHandler("test-stream-01", nil)
	}
}
