package centre

import (
	mq "kafka-iot-connect/client/mqtt"
	"kafka-iot-connect/client/mysql"
	"kafka-iot-connect/client/redis"
	"reflect"
	"sync"
	logging "kafka-iot-connect/log"

	"github.com/gofrs/uuid"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	Id     int              `json:"id" example:"1" format:"int64"`
	Uuid   uuid.UUID        `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
	Name   string           `json:"name" example:"rabb_client_01"`
	Client *mq.MqttConfig `json:"-"`
}

type MysqlEkuiperCentre struct {
	*sync.RWMutex
	Name      string
	Activated bool
	Clients   []*MqttClient
	Mysql     *mysql.MysqlDB
	Redis     *redis.RedisTimeSeriesClient
}

type MqttMessageHandler func(*mq.AllMessageChannels, *logging.MyFileLogger) mqtt.MessageHandler


type logInfo struct{}

var info = logInfo{}

var pkgName = reflect.TypeOf(info).PkgPath()