package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	logging "kafka-iot-connect/log"
	"kafka-iot-connect/tool"
	"log"
	"os"
	"reflect"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type logInfo struct{}

var info = logInfo{}

var pkgName = reflect.TypeOf(info).PkgPath()

func NewTlsConfig(path string) *tls.Config {
	certpool := x509.NewCertPool()
	ca, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM(ca)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		RootCAs:            certpool,
		InsecureSkipVerify: true,
	}
}

func (m *MqttConfig) Publish(msg string, l *logging.MyFileLogger, topics ...string) {
	if !m.Client.IsConnected() {
		log.Println("Please connect to mqtt server first")
	} else {
		if len(topics) == 0 {
			for _, topic := range m.PublishTopics {
				token := m.Client.Publish(topic.Name, 0, false, msg)
				if l != nil {
					l.Log(logging.LogInfo, pkgName, fmt.Sprintf("Published to topic %s, server: %s:%d", topic.Name, m.Host, m.Port))
				} else {
					log.Printf("Published to topic %s, server: %s:%d\n", topic.Name, m.Host, m.Port)
				}
				token.Wait()
			}
		} else {
			for _, topic := range topics {
				token := m.Client.Publish(topic, 0, false, msg)
				if l != nil {
					l.Log(logging.LogInfo, pkgName, fmt.Sprintf("Published to topic %s, server: %s:%d", topic, m.Host, m.Port))
				} else {
					log.Printf("Published to topic %s, server: %s:%d\n", topic, m.Host, m.Port)
				}
				token.Wait()
			}
		}
	}
}

func (m *MqttConfig) Subscribe(l *logging.MyFileLogger, topics ...string) {
	if !m.Client.IsConnected() {
		log.Println("Please connect to mqtt server first")
	} else {
		if len(topics) == 0 {
			for _, topic := range m.SubcribleTopics {
				token := m.Client.Subscribe(topic.Name, 1, nil)
				token.Wait()
				if l == nil {
					log.Printf("ClientId: %d subscribed to topic: %s, host: %s, port: %d \n", m.Id, topic.Name, m.Host, m.Port)
				} else {
					l.Log(logging.LogInfo, pkgName, fmt.Sprintf("ClientId: %d subscribed to topic: %s, host: %s, port: %d", m.Id, topic.Name, m.Host, m.Port))
				}
			}
		} else {
			for _, topic := range topics {
				if tool.Contain(m.SubcribleTopics, Topic{Name: topic}) > -1 {
					if l == nil {
						log.Println("Already subscribed to topic: ", topic)
					} else {
						l.Log(logging.LogInfo, pkgName, fmt.Sprintf("Already subscribed to topic: %s", topic))
					}
				} else {
					token := m.Client.Subscribe(topic, 1, nil)
					token.Wait()
					m.SubcribleTopics = append(m.SubcribleTopics, Topic{
						Name: topic,
					})
					if l == nil {
						log.Printf("ClientId: %d subscribed to topic: %s, host: %s, port: %d \n", m.Id, topic, m.Host, m.Port)
					} else {
						l.Log(logging.LogInfo, pkgName, fmt.Sprintf("ClientId: %d subscribed to topic: %s, host: %s, port: %d", m.Id, topic, m.Host, m.Port))
					}
				}
			}
		}
	}
}

func (m *MqttConfig) GetSubTopicByName(name string) int {
	for ind, ele := range m.SubcribleTopics {
		if ele.Name == name {
			return ind
		}
	}
	return -1
}

func (m *MqttConfig) RemoveSubTopicByName(name string) error {
	tops := []Topic{}
	for _, ele := range m.SubcribleTopics {
		if ele.Name != name {
			tops = append(tops, ele)
		}
	}
	m.SubcribleTopics = tops
	if len(tops) == len(m.SubcribleTopics) {
		return fmt.Errorf("failed to remove sub topic by name: %s", name)
	} else {
		return nil
	}
}

func (m *MqttConfig) RemoveSubTopicByNames(names ...string) {
	tops := []Topic{}
	for _, ele := range m.SubcribleTopics {
		var remove bool
		for _, name := range names {
			if ele.Name == name {
				remove = true
			}
		}
		if !remove {
			tops = append(tops, ele)
		}
	}
	m.SubcribleTopics = tops
}

func (m *MqttConfig) SubscribeStream(topic string) int {
	if !m.Client.IsConnected() {
		log.Println("Please connect to mqtt server first")
		return -1
	} else {
		if tool.Contain(m.SubcribleTopics, Topic{Name: topic}) > -1 {
			log.Println("Already subscribed to topic: ", topic)
			return 0
		} else {
			token := m.Client.Subscribe(topic, 1, nil)
			token.Wait()
			m.SubcribleTopics = append(m.SubcribleTopics, Topic{
				Name: topic,
			})
			log.Printf("ClientId: %d subscribed to topic: %s, host: %s, port: %d \n", m.Id, topic, m.Host, m.Port)
			return 1
		}
	}
}

func (m *MqttConfig) Unsubscribe(topics ...string) {
	if !m.Client.IsConnected() {
		log.Println("Please connect to mqtt server first")
	} else {
		m.RemoveSubTopicByNames(topics...)
		token := m.Client.Unsubscribe(topics...)
		token.Wait()
	}
}

func (m *MqttConfig) UnsubscribeUnuseTopic(topics ...string) {
	if !m.Client.IsConnected() {
		log.Println("Please connect to mqtt server first")
	} else {
		for _, topic := range topics {
			if m.AllCh.UnsubscribeUnuseTopic(topic) {
				m.Unsubscribe(topic)
			}
		}
	}
}

func (m *MqttConfig) SetClientOptions() {
	opts := mqtt.NewClientOptions()
	if m.Scheme == "" {
		m.Scheme = "tcp"
	}
	opts.AddBroker(fmt.Sprintf("%s://%s:%d", m.Scheme, m.Host, m.Port))
	opts.SetClientID(m.ClientId)
	opts.SetUsername(m.Username)
	opts.SetPassword(m.Password)
	opts.SetCleanSession(false)
	m.Options = opts
}

func GetScheme(tlsPath string) string {
	if tlsPath == "" {
		return "tcp"
	} else {
		return "ssl"
	}
}

func (m *MqttConfig) SetTlsConfig() {
	if m.TlsPath != "" {
		tlsConfig := NewTlsConfig(m.TlsPath)
		if m.Options == nil {
			m.Options = &mqtt.ClientOptions{
				TLSConfig: tlsConfig,
			}
		} else {
			m.Options.SetTLSConfig(tlsConfig)
		}
	}
}

func (m *MqttConfig) SetCallbacks(msgF mqtt.MessageHandler, conF mqtt.OnConnectHandler, conLostF mqtt.ConnectionLostHandler) {
	m.Options.SetDefaultPublishHandler(msgF)
	m.Options.OnConnect = conF
	m.Options.OnConnectionLost = conLostF
}

func (m *MqttConfig) SetClient() error {
	m.Client = mqtt.NewClient(m.Options)
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	} else {
		return nil
	}
}

func (m *MqttConfig) UpdateCallbacks() {
	m.SetCallbacks(
		m.MessageHandler(m.AllCh, m.Log),
		m.OnConnectHandler(m.AllCh, m.Log),
		m.ConnectionLostHandler(m.AllCh, m.Log),
	)
}

func (m *MqttConfig) DefaultClient() error {
	m.SetClientOptions()
	m.SetTlsConfig()
	m.DefaultHandler()
	m.SetCallbacks(
		m.MessageHandler(m.AllCh, m.Log),
		m.OnConnectHandler(m.AllCh, m.Log),
		m.ConnectionLostHandler(m.AllCh, m.Log),
	)
	if err := m.SetClient(); err != nil {
		return err
	}
	m.SubscribeChannels()
	return nil
}

func (m *MqttConfig) SubscribeChannels() {
	for _, ele := range m.AllCh.Chs {
		for _, el := range ele.Topics {
			contains := false
			for _, e := range m.SubcribleTopics {
				if el.Name == e.Name {
					contains = true
				}
			}
			if !contains {
				m.SubcribleTopics = append(m.SubcribleTopics, el)
			}
		}
	}
}

func (m *MqttConfig) DefaultHandler() {
	if m.UseDefaultHandlers {
		m.MessageHandler = func(amc *AllMessageChannels, l *logging.MyFileLogger) mqtt.MessageHandler {
			return amc.DefaultMessageChannelHandler()
		}
		m.OnConnectHandler = func(amc *AllMessageChannels, l *logging.MyFileLogger) mqtt.OnConnectHandler {
			return amc.DefaultChannelConnectHandler()
		}
		m.ConnectionLostHandler = func(amc *AllMessageChannels, l *logging.MyFileLogger) mqtt.ConnectionLostHandler {
			return amc.DefaultChannelConnectLostHandler()
		}
		m.SetCallbacks(
			m.MessageHandler(m.AllCh, m.Log),
			m.OnConnectHandler(m.AllCh, m.Log),
			m.ConnectionLostHandler(m.AllCh, m.Log),
		)
	}
}
