package store

import (
	"../config"
	"fmt"
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
	if err != nil {
		fmt.Println(err)
	}
	c = client
	fmt.Println(c)
}

func GetData(key string) interface{} {
	val, err := c.Get(key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func SetData(key string, val interface{}, dur time.Duration) {
	err := c.Set(key, val, dur).Err()
	fmt.Println(c)
	if err != nil {
		fmt.Println(err)
	}
}
