package config

import (
	"log"

	"github.com/spf13/viper"
)

type EnvConfig struct {
	APP_NAME        string
	APP_DESCRIPTION string
	APP_VERSION     string
	APP_PORT        string
	APP_SECRET_KEY  string

	PREFORK_ENABLED bool

	PASETO_LOCAL_SECRET_SYMMETRIC_KEY string
	PASETO_LOCAL_EXPIRATION_HOURS     uint64

	REFRESH_TOKEN_EXPIRATION_DAYS uint16

	DB_HOST                     string
	DB_PORT                     string
	DB_USER                     string
	DB_PASSWORD                 string
	DB_NAME                     string
	DB_SSL_MODE                 string // values are "disable", "require", "verify-ca", "verify-full"
	DB_MAX_IDLE_CONNECTION      uint16
	DB_MAX_OPEN_CONNECTION      uint16
	DB_MAX_LIFETIME_CONNECTION  uint16
	DB_MAX_IDLE_TIME_CONNECTION uint16

	DB_ENABLE_AUTO_MIGRATE bool

	DB_ENABLE_LOG      bool
	REQUEST_ENABLE_LOG bool

	REDIS_PREFIX   string
	REDIS_HOST     string
	REDIS_PORT     string
	REDIS_USERNAME string
	REDIS_PASSWORD string

	JOB_CONCURRENCY         int
	JOB_MONITORING_USERNAME string
	JOB_MONITORING_PASSWORD string

	SMTP_HOST       string
	SMTP_PORT       string
	SMTP_USERNAME   string
	SMTP_PASSWORD   string
	SMTP_FROM_EMAIL string
	SMTP_FROM_NAME  string

	RATE_LIMITER_MAX           int
	RATE_LIMITER_TTL_IN_SECOND uint64

	S3_BUCKET     string
	S3_ENDPOINT   string
	S3_REGION     string
	S3_ACCESS_KEY string
	S3_SECRET_KEY string
}

var GetConfig *EnvConfig

func LoadConfig(path string) (config *EnvConfig) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()

	viper.SetDefault("APP_PORT", "3000")
	viper.SetDefault("APP_NAME", "Go API")
	viper.SetDefault("APP_DESCRIPTION", "Go API Description")
	viper.SetDefault("APP_VERSION", "0.0.0.0")
	viper.SetDefault("APP_SECRET_KEY", "_iWv(UWEp^pf$<?")

	viper.SetDefault("PREFORK_ENABLED", false)

	viper.SetDefault("PASETO_LOCAL_SECRET_SYMMETRIC_KEY", "CX3cZoWd13exnqlxAWMwtj2TvRQXKOKi")
	viper.SetDefault("PASETO_LOCAL_EXPIRATION_HOURS", 9)

	viper.SetDefault("REFRESH_TOKEN_EXPIRATION_DAYS", 14)

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "gofiber_api")
	viper.SetDefault("DB_SSL_MODE", "disable")

	viper.SetDefault("DB_ENABLE_AUTO_MIGRATE", true)

	viper.SetDefault("DB_MAX_IDLE_CONNECTION", 10)
	viper.SetDefault("DB_MAX_OPEN_CONNECTION", 50)
	viper.SetDefault("DB_MAX_LIFETIME_CONNECTION", 60)
	viper.SetDefault("DB_MAX_IDLE_TIME_CONNECTION", 60)

	viper.SetDefault("DB_ENABLE_LOG", true)
	viper.SetDefault("REQUEST_ENABLE_LOG", true)

	viper.SetDefault("REDIS_PREFIX", "gofiber_api_")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_USERNAME", "")
	viper.SetDefault("REDIS_PASSWORD", "")

	viper.SetDefault("JOB_CONCURRENCY", 10)
	viper.SetDefault("JOB_MONITORING_USERNAME", "developer")
	viper.SetDefault("JOB_MONITORING_PASSWORD", "password")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		log.Fatal("Using default config")
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshal config file, %s", err)
		log.Fatal("Using default config")
	}

	return config
}

func InitConfig(path string) {
	GetConfig = LoadConfig(path)
}
