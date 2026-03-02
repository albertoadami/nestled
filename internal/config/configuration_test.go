package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigSuccessfully(t *testing.T) {
	viper.AddConfigPath("../..") // add project root as config path

	configuration, error := LoadConfig()

	assert.NoError(t, error)

	dbConfig := configuration.Database

	assert.Equal(t, dbConfig.Host, "localhost")
	assert.Equal(t, dbConfig.Port, 5432)
	assert.Equal(t, dbConfig.User, "postgres")
	assert.Equal(t, dbConfig.Password, "password")
	assert.Equal(t, dbConfig.Name, "nestled_db")

	jwtConfig := configuration.JWT

	assert.Equal(t, "your_secret_key", jwtConfig.Secret)
	assert.Equal(t, 6, jwtConfig.Expiration)

}
