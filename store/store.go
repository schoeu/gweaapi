package store

import (
	"../config"
	"fmt"
	"github.com/go-redis/redis"
)

var client *redis.Client

func GetRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.DB,
	})

	pong, err := client.Ping().Result()
}

func SetData(k string, v string) {
	err := client.Set(k, v, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}
