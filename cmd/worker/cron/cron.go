package main

import (
	"GoFiber-API/internal/config"
	"GoFiber-API/internal/log"

	"github.com/hibiken/asynq"
)

func main() {
	config.InitConfig(".")

	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     config.GetConfig.REDIS_HOST + ":" + config.GetConfig.REDIS_PORT,
			Username: config.GetConfig.REDIS_USERNAME,
			Password: config.GetConfig.REDIS_PASSWORD,
		},
		nil,
	)

	// Tasks

	// emailData, _ := json.Marshal(auth.ResetPasswordJobPayload{Email: "ariakm25@gmail.com"})

	// taskID, err := scheduler.Register("@every 24h", asynq.NewTask(
	// 	auth.TypeAuthResetPasswordJob,
	// 	emailData,
	// ),
	// 	asynq.Retention(1*time.Hour),
	// )

	// if err != nil {
	// 	internal_log.Logger.Sugar().Errorf("Error registering an task: %s", err.Error())
	// }
	// log.Printf("registered an task: %q\n", taskID)

	if err := scheduler.Run(); err != nil {
		internal_log.Logger.Sugar().Errorf("Error running scheduler: %s", err.Error())
	}
}
