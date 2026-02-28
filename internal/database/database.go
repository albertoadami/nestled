package database

import (
	"fmt"

	"github.com/albertoadami/nestled/internal/config"
	"github.com/jmoiron/sqlx"
)

func Connect(config *config.DatabaseConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name,
	)
	return sqlx.Connect("postgres", dsn)
}
