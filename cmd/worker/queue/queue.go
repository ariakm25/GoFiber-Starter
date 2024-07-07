package main

import (
	"GoFiber-API/app/auth"
	database "GoFiber-API/external/database/postgres"
	"GoFiber-API/external/mail"
	"GoFiber-API/internal/config"
	internal_logger "GoFiber-API/internal/log"
	"log"

	"github.com/hibiken/asynq"
)

func main() {
	config.InitConfig(".")
	internal_logger.InitLogger()
	mail.NewMail()
	// Connect to the database
	err := database.ConnectDB(config.GetConfig.DB_HOST, config.GetConfig.DB_PORT, config.GetConfig.DB_USER, config.GetConfig.DB_PASSWORD, config.GetConfig.DB_NAME, config.GetConfig.DB_SSL_MODE)
	if err != nil {
		log.Fatalf("Error connect to Database: %s", err)
		panic(err)
	}

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
