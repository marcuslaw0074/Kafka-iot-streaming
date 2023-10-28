package main

import (
	"context"
	"encoding/json"
	"fmt"
	"kafka-iot-connect/client/kafka/consumer"
	"kafka-iot-connect/client/kafka/producer"
	"kafka-iot-connect/client/mqtt"
	"kafka-iot-connect/client/mqtt/centre"
	logging "kafka-iot-connect/log"
	"kafka-iot-connect/model"
	"kafka-iot-connect/tool"
	"math/rand"
	"time"

	"github.com/IBM/sarama"
	mq "github.com/eclipse/paho.mqtt.golang"
)

type DataModel struct {
	Name string `json:"Name"`
	Tags struct {
		BuildingName string `json:"buildingName"`
		Id           string `json:"id"`
	} `json:"Tags"`
	Fields struct {
		Reading string  `json:"reading"`
		Value   float64 `json:"value"`
	} `json:"Fields"`
	T string `json:"T"`
}

func DataModelRemoveDuplicate(ds []DataModel) []DataModel {
	allkeys := map[string]bool{}
	ds2 := make([]DataModel, 0)
	for _, i := range ds {
		if _, v := allkeys[i.Tags.Id]; !v {
			allkeys[i.Tags.Id] = true
			ds2 = append(ds2, i)
		}
	}
	return ds2
}

func initializeMqttConnectClient(ctx context.Context) {
	log := logging.StartLogger("./log/api.log", 100000)
	log.ClearLog()
	mqConf, err := tool.ParseJsonFile[mqtt.MqttConf]("./mqtt.json")
	if err != nil {
		panic(err)
	}
	conf := mqtt.MqttConfig{
		Host:                  mqConf.Host,
		Port:                  mqConf.Port,
		ClientId:              "kafka-connect-mqtt-client-01",
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
							Name: "etl/action/bms/etl-214/af94c736-8da1-4ab9-f5ee-08db213be062/8",
						},
					},
					IsActive: true,
					CallbackAction: func(c mq.Client, b [][]byte, t []mqtt.Topic, i int) {
						fmt.Printf("received messaged from topic: %s", t[0].Name)
						kp := producer.InitializeKafkaConnectProducer("localhost:9092", "3.3.2", true, 1)
						for _, bb := range b {
							var data []DataModel
							if err := json.Unmarshal(bb, &data); err != nil {
								log.Log(logging.LogDebug, "main", err)
								continue
							} else {
								fmt.Println(string(bb))
								// fmt.Println("HHHHH")
								// data = tool.RemoveDuplicate(data)
								data = DataModelRemoveDuplicate(data)
								// hh, _ := json.Marshal(data)
								// fmt.Println(string(hh))
								kp.Produce(
									tool.Map(data,
										func(d DataModel, i int) *sarama.ProducerMessage {
											return &sarama.ProducerMessage{
												Topic: "kafka-streams-energy-raw-data",
												Key:   model.BMSDataTypeEncoder(d.Tags.Id, "elec"),
												Value: model.BMSRawDataEncoder(d.Tags.Id, 1, d.Fields.Value),
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
		go func() {
			time.Sleep(3 * time.Hour)
			vals := []float64{
				rand.Float64() * 10,
				rand.Float64() * 10,
				rand.Float64() * 10,
				rand.Float64() * 10,
			}
			ls := []string{
				"esg/NF/Mills/electricity/Device1:Analog_Value_M6_ACB_T4_1_total_kWh",
				"esg/NF/Mills/electricity/Device1:Analog_Value_M6_ACB_T4_2_total_kWh",
				"esg/NF/Mills/electricity/Device1:Analog_Value_M6_ACB_T4_total_kWh",
				"esg/NF/Mills/electricity/Device1:Analog_Value_M6_Shop104A_total_kWh",
			}
			for {
				t := time.Now().UTC().Format("2006-01-02 15:04:05")
				var data []DataModel
				for ind, ele := range ls {
					data = append(data, DataModel{
						Name: "mp_raw_3000014",
						Tags: struct {
							BuildingName string "json:\"buildingName\""
							Id           string "json:\"id\""
						}{
							BuildingName: "mills",
							Id:           ele,
						},
						Fields: struct {
							Reading string  "json:\"reading\""
							Value   float64 "json:\"value\""
						}{
							Reading: "",
							Value:   vals[ind],
						},
						T: t,
					})
					vals[ind] = vals[ind] + rand.Float64()*5
				}
				b, _ := json.Marshal(data)
				conf.Publish(string(b), log, "etl/action/bms/marcus/af94c736-8da1-4ab9-f5ee-08db213be062/8")
				time.Sleep(3 * time.Second)
			}
		}()
	}
}

func initializeKafkaConnectClient(ctx context.Context) {
	log := logging.StartLogger("./log/api.log", 100000)
	log.ClearLog()
	mqConf, err := tool.ParseJsonFile[mqtt.MqttConf]("./mqtt.json")
	if err != nil {
		panic(err)
	}
	conf := mqtt.MqttConfig{
		Host:                  mqConf.Host,
		Port:                  mqConf.Port,
		ClientId:              "kafka-connect-mqtt-client-02",
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
		kc.Consume(context.Background(), "kafka-streams-energy-realtime-data", func(cm *sarama.ConsumerMessage) error {
			conf.Publish(string(cm.Value), log, "etl/energy/bms/marcus/af94c736-8da1-4ab9-f5ee-08db213be062/8")
			return nil
		})
	}
}
