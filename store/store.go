package store

import (
	"../config"
	"fmt"
	"github.com/go-redis/redis"
)

func GetRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err)
	}
	return client
}
