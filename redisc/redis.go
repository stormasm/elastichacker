package redisc

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type RedisConfig struct {
	Hostname string
	Port     string
}

func (c *RedisConfig) Connect_string() string {
	connect := fmt.Sprint(c.Hostname, ":", c.Port)
	return connect
}

func NewRedisConfig() *RedisConfig {
	cfg := &RedisConfig{
		Hostname: "localhost",
		Port:     "6379",
	}
	return cfg
}

func getRedisConn() (c redis.Conn) {

	cfg := NewRedisConfig()
	connect_string := cfg.Connect_string()
	c, err := redis.Dial("tcp", connect_string)
	if err != nil {
		panic(err)
	}
	return c
	// defer c.Close()
}
