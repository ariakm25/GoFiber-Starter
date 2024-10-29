package redis

import (
	"GoFiber-API/internal/config"
	"context"
	"log"
	"runtime"
	"strconv"

	"github.com/gofiber/storage/redis/v3"
)

var RedisStore *redis.Storage

func NewRedisStore() error {

	port, err := strconv.Atoi(config.GetConfig.REDIS_PORT)

	if err != nil {
		port = 6379
	}

	RedisStore = redis.New(redis.Config{
		Host:      config.GetConfig.REDIS_HOST,
		Port:      port,
		Username:  "",
		Password:  config.GetConfig.REDIS_PASSWORD,
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})

	val, err := RedisStore.Conn().Ping(context.Background()).Result()

	log.Printf("Redis Ping: %s", val)

	if err != nil {
		return err
	}

	if val != "PONG" {
		return err
	}

	return nil
}
