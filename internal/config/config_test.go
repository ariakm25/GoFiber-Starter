package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		InitConfig("../../")
		err := viper.ReadInConfig()

		require.Nil(t, err)
	})

}
