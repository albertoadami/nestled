package testhelpers

import (
	"context"
	"fmt"
	"testing"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupPostgres(t *testing.T) (*sqlx.DB, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections"),
		),
	}

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	host, err := postgres.Host(ctx)
	assert.NoError(t, err)
	portObj, err := postgres.MappedPort(ctx, "5432")
	assert.NoError(t, err)

	dsn := fmt.Sprintf("host=%s port=%s user=testuser password=testpass dbname=testdb sslmode=disable",
		host, portObj.Port())

	db, err := sqlx.Connect("postgres", dsn)
	assert.NoError(t, err)

	terminate := func() {
		db.Close()
		_ = postgres.Terminate(ctx)
	}

	// run migrations
	runMigrations(t, host, portObj.Port())

	return db, terminate
}

func runMigrations(t *testing.T, host, port string) {
	dbURL := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port)

	m, err := migrate.New(
		"file://../../migrations",
		dbURL,
	)
	assert.NoError(t, err)
	if m == nil {
		t.Fatal("migrate instance is nil, check migrations path and db url")
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		assert.NoError(t, err)
	}
}
