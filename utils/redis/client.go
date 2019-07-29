package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ipweb-group/go-sdk/conf"
)

var redisClient *redis.Client

func GetClient() *redis.Client {
	if redisClient == nil {
		config := conf.GetConfig().RedisConfig

		redisClient = redis.NewClient(&redis.Options{
			Addr:     config.Addr,
			Password: config.Password,
			DB:       config.DB,
			OnConnect: func(conn *redis.Conn) error {
				fmt.Println("Redis connection established")
				return nil
			},
		})

		// 尝试 ping redis server，失败时表示连接 Redis 失败
		_, err := redisClient.Ping().Result()
		if err != nil {
			fmt.Println("Connect to Redis server failed")
			panic(err)
		}
	}

	return redisClient
}
