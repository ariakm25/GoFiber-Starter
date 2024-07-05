package main

import (
	"GoFiber-API/app/auth"
	"GoFiber-API/internal/config"
	internal_logger "GoFiber-API/internal/log"

	"github.com/hibiken/asynq"
)

func main() {
	config.InitConfig(".")

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     config.GetConfig.REDIS_HOST + ":" + config.GetConfig.REDIS_PORT,
			Username: config.GetConfig.REDIS_USERNAME,
			Password: config.GetConfig.REDIS_PASSWORD,
		},
		asynq.Config{
			Concurrency: config.GetConfig.JOB_CONCURRENCY,
		},
	)

	jobHandler := asynq.NewServeMux()

	// Auth Jobs
	jobHandler.HandleFunc(auth.TypeAuthResetPasswordJob, auth.HandleAuthResetPasswordJob)

	if err := srv.Run(jobHandler); err != nil {
		internal_logger.Logger.Sugar().Errorf("Error running jobHandler: %s", err.Error())
	}
}
