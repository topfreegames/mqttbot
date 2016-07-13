// mqttbot
// https://github.com/topfreegames/arkadiko
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package redisclient

import (
	"fmt"
	"sync"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
)

var once sync.Once
var client *RedisClient

type RedisClient struct {
	Pool *redis.Pool
}

func GetRedisClient(redisHost string, redisPort int, redisPass string) *RedisClient {
	once.Do(func() {
		client = &RedisClient{}
		redisAddress := fmt.Sprintf("%s:%d", redisHost, redisPort)
		client.Pool = redis.NewPool(func() (redis.Conn, error) {
			if viper.GetString("redis.password") != "" {
				c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", viper.GetString("redis.host"),
					viper.GetInt("redis.port")), redis.DialPassword(viper.GetString("redis.password")))
				if err != nil {
					logger.Logger.Error(err.Error())
				}
				return c, err
			}
			c, err := redis.Dial("tcp", redisAddress)
			if err != nil {
				if err != nil {
					logger.Logger.Error(err.Error())
				}
			}
			return c, err
		}, viper.GetInt("redis.maxPoolSize"))
	})
	return client
}
