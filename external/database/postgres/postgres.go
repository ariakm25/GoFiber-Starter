package database

import (
	"GoFiber-API/internal/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Connection *gorm.DB

// ConnectDB connects to the database
func ConnectDB(host string, port string, user string, password string, db_name string, ssl_mode string) error {
	// Build the DSN
	build_dsn := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + db_name + " sslmode=" + ssl_mode

	logMode := logger.Default.LogMode(logger.Info)

	if !config.GetConfig.DB_ENABLE_LOG {
		logMode = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(build_dsn), &gorm.Config{
		Logger: logMode,
	})
	if err != nil {
		return err
	}

	Connection = db

	sqlDB, err := Connection.DB()

	if err != nil {
		log.Fatalf("Error get DB: %s", err)
		panic(err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(int(config.GetConfig.DB_MAX_IDLE_CONNECTION))

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(int(config.GetConfig.DB_MAX_OPEN_CONNECTION))

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Duration(config.GetConfig.DB_MAX_LIFETIME_CONNECTION) * time.Second)

	// SetConnMaxIdleTime sets the maximum amount of time a connection may be idle.
	sqlDB.SetConnMaxIdleTime(time.Duration(config.GetConfig.DB_MAX_IDLE_TIME_CONNECTION) * time.Second)

	return nil
}
