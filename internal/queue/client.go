package queue

import (
	"GoFiber-API/internal/config"

	"github.com/hibiken/asynq"
)

var QueueClient *asynq.Client

func InitQueueClient() {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     config.GetConfig.REDIS_HOST + ":" + config.GetConfig.REDIS_PORT,
		Username: config.GetConfig.REDIS_USERNAME,
		Password: config.GetConfig.REDIS_PASSWORD,
	})

	QueueClient = client
}
