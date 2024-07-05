package main

import (
	database "GoFiber-API/external/database/postgres"
	"GoFiber-API/internal/config"
	internal_logger "GoFiber-API/internal/log"
	"GoFiber-API/internal/migration"
	"flag"
	"log"

	"GoFiber-API/app"
	"GoFiber-API/internal/queue"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Init Configuration from .env
	config.InitConfig(".")

	// Init Logger
	internal_logger.InitLogger()

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

	// Main Module
	app.MainModule(api)

	api.Listen(":" + config.GetConfig.APP_PORT)

}
