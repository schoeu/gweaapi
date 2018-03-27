package store

import (
	"../config"
	"fmt"
	"github.com/go-redis/redis"
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

func GetData(key string) string {
	val, err := c.Get(key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func SetData(key string, val string) {
	err := c.Set(key, val, 0).Err()
	fmt.Println(c)
	if err != nil {
		fmt.Println(err)
	}
}
