package main

import (
	"GoFiber-API/app/greeter"
	database "GoFiber-API/external/database/postgres"
	"GoFiber-API/external/database/redis"
	"GoFiber-API/external/mail"
	"GoFiber-API/infra/middleware"
	internal_casbin "GoFiber-API/internal/casbin"
	"GoFiber-API/internal/config"
	internal_logger "GoFiber-API/internal/log"
	"GoFiber-API/internal/migration"
	"flag"
	"log"
	"time"

	"GoFiber-API/app"
	"GoFiber-API/internal/queue"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	is_grpc = flag.Bool("grpc", false, "Use grpc instead of http")
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Init Configuration from .env
	config.InitConfig(".")

	// Init Logger
	internal_logger.InitLogger()

	// Init Redis Store
	redis.NewRedisStore()

	// Init Mail
	mail.NewMail()

	// Init Worker
	queue.InitQueueClient()
	defer queue.QueueClient.Close()

	// Connect to the database
	err := database.ConnectDB(config.GetConfig.DB_HOST, config.GetConfig.DB_PORT, config.GetConfig.DB_USER, config.GetConfig.DB_PASSWORD, config.GetConfig.DB_NAME, config.GetConfig.DB_SSL_MODE)
	if err != nil {
		log.Fatalf("Error connect to Database: %s", err)
		panic(err)
	}

	// Auto Migrate the database
	migration.Migrate()

	// Init Casbin
	internal_casbin.InitAdapter("casbin-rbac.conf", config.GetConfig.DB_HOST, config.GetConfig.DB_PORT, config.GetConfig.DB_USER, config.GetConfig.DB_PASSWORD, config.GetConfig.DB_NAME, config.GetConfig.DB_SSL_MODE)
	middleware.InitRbac("casbin-rbac.conf", internal_casbin.CasbinAdapter)

	if *is_grpc {
		greeter.StartGrpcGreeter()
	} else {
		// Setup Fiber App
		api := fiber.New(fiber.Config{
			Prefork: true,
			AppName: config.GetConfig.APP_NAME,
		})

		// Middleware
		api.Use(recover.New())
		api.Use(logger.New())
		api.Use(cors.New())
		api.Use(helmet.New())
		api.Use(limiter.New(limiter.Config{
			Max:        config.GetConfig.RATE_LIMITER_MAX,
			Expiration: time.Duration(config.GetConfig.RATE_LIMITER_TTL_IN_SECOND) * time.Second,
			// Next: func(c *fiber.Ctx) bool {
			// 	return c.IP() == "127.0.0.1"
			// },
		}))

		// Main Module
		app.MainModule(api)

		api.Listen(":" + config.GetConfig.APP_PORT)
	}

}
