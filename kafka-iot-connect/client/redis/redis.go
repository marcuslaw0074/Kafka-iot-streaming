package redis

import "github.com/redis/rueidis"

type RedisTimeSeriesClient struct {
	Host      string
	Port      int
	Username  string
	Password  string
	SecretKey string
	Client    rueidis.Client
}