package store

import (
	"../utils"
	"../config"
	"github.com/go-redis/redis"
	"time"
)

var c *redis.Client

func GetRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := client.Ping().Result()
	utils.ErrHandle(err)
	c = client
}

func GetData(key string) interface{} {
	val, err := c.Get(key).Result()
	utils.ErrHandle(err)
	return val
}

func SetData(key string, val interface{}, dur time.Duration) {
	err := c.Set(key, val, dur).Err()
	utils.ErrHandle(err)
}
