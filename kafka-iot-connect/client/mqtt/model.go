package mqtt

import (
	"sync"
	"time"
	logging "kafka-iot-connect/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Topic struct {
	DataSoruceId int
	ProjectId    string
	Name         string
	Id           int
	Size         int
	Uuid         string
}

type MessageChannel struct {
	Ch             chan []byte                               `json:"-"`
	Id             string                                    `json:"id,omitempty"`
	Topics         []Topic                                   `json:"topics,omitempty"`
	PublishTopics  []Topic                                   `json:"publishTopics,omitempty"`
	Data           [][]byte                                  `json:"-"`
	WaitingTime    time.Duration                             `json:"-"`
	DefaultTime    time.Duration                             `json:"-"`
	CallbackAction func(mqtt.Client, [][]byte, []Topic, int) `json:"-"`
	IsActive       bool                                      `json:"isActive"`
}

type AllMessageChannels struct {
	mu           sync.RWMutex
	Id           string           `json:"id,omitempty"`
	Chs          []MessageChannel `json:"chs,omitempty"`
	MaxGoroutine int              `json:"maxGoroutine,omitempty"`
}

type MqttConfig struct {
	Id                    int                                                                         `json:"id,omitempty"`
	Host                  string                                                                      `json:"host,omitempty"`
	Port                  int                                                                         `json:"port,omitempty"`
	ClientId              string                                                                      `json:"clientId,omitempty"`
	Username              string                                                                      `json:"username,omitempty"`
	Password              string                                                                      `json:"password,omitempty"`
	TlsPath               string                                                                      `json:"tlsPath,omitempty"`
	Scheme                string                                                                      `json:"scheme,omitempty"`
	Log                   *logging.MyFileLogger                                                       `json:"-"`
	PublishTopics         []Topic                                                                     `json:"publishTopics,omitempty"`
	SubcribleTopics       []Topic                                                                     `json:"subcribleTopics,omitempty"`
	Options               *mqtt.ClientOptions                                                         `json:"-"`
	Client                mqtt.Client                                                                 `json:"-"`
	UseDefaultHandlers    bool                                                                        `json:"useDefaultHandlers,omitempty"`
	MessageHandler        func(*AllMessageChannels, *logging.MyFileLogger) mqtt.MessageHandler        `json:"-"`
	OnConnectHandler      func(*AllMessageChannels, *logging.MyFileLogger) mqtt.OnConnectHandler      `json:"-"`
	ConnectionLostHandler func(*AllMessageChannels, *logging.MyFileLogger) mqtt.ConnectionLostHandler `json:"-"`
	AllCh                 *AllMessageChannels                                                         `json:"allCh,omitempty"`
}
