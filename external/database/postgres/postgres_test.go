package database_test

import (
	postgres "GoFiber-API/external/database/postgres"
	"GoFiber-API/internal/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	config.InitConfig("../../../")
}

func TestConnectDB(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := postgres.ConnectDB(config.GetConfig.DB_HOST, config.GetConfig.DB_PORT, config.GetConfig.DB_USER, config.GetConfig.DB_PASSWORD, config.GetConfig.DB_NAME, config.GetConfig.DB_SSL_MODE)
		require.Nil(t, err)
	})

	t.Run("failed", func(t *testing.T) {
		err := postgres.ConnectDB("localhost", "5432", "postgres", "postgres", "invalid", "disable")
		require.NotNil(t, err)
	})
}
