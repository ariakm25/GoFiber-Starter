package monitoring

import (
	"GoFiber-API/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
)

func Module(router fiber.Router) {
	route := router.Group("/monitoring")

	route.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			config.GetConfig.JOB_MONITORING_USERNAME: config.GetConfig.JOB_MONITORING_PASSWORD,
		},
	}))

	// Jobs Monitoring
	route.All("/jobs/*", adaptor.HTTPHandler(asynqmon.New(asynqmon.Options{
		RootPath: "/monitoring/jobs",
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     config.GetConfig.REDIS_HOST + ":" + config.GetConfig.REDIS_PORT,
			Username: config.GetConfig.REDIS_USERNAME,
			Password: config.GetConfig.REDIS_PASSWORD,
		},
	})))
}
