package redis

import (
	"github.com/go-redis/redis"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/utils"
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
				utils.GetLogger().Info("Redis connection established")
				return nil
			},
		})

		// 尝试 ping redis server，失败时表示连接 Redis 失败
		_, err := redisClient.Ping().Result()
		if err != nil {
			utils.GetLogger().Error("Connect to Redis server failed")
			panic(err)
		}
	}

	return redisClient
}
