# Nestled

**A family expense tracking backend service** for managing household finances, tracking expenses, and setting budget alerts on expense categories.

> ⚠️ **Status**: Implementation in progress

## Overview

Nestled is the backend service that enables families to:

- Track income and expenses
- Organize transactions by categories
- Define and monitor budgets for each category
- Receive alerts when spending approaches or exceeds budget limits
- View financial reports and spending insights

## Tech Stack

- **Language**: Go 1.26.0
- **Framework**: [Gin Web Framework](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL 17
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Testing**: Go's built-in testing + [testify](https://github.com/stretchr/testify) + [testcontainers](https://github.com/testcontainers/testcontainers-go)
- **ORM**: [sqlx](https://github.com/jmoiron/sqlx)

## Prerequisites

Before you begin, ensure you have installed:

- **Go**: 1.26.0 or later ([download](https://golang.org/dl/))
- **PostgreSQL**: 17 or later (for local development)
  - OR **Docker & Docker Compose** (for containerized setup)
- **migrate CLI**: For running database migrations locally ([installation](https://github.com/golang-migrate/migrate#cli-usage))

## Getting Started

### 1. Clone the repository

```bash
git clone <repo-url>
cd nestled
```

### 2. Set up environment variables

Copy the default configuration:

```bash
cp config.yml .env
```

Or create a `.env` file with your database credentials:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=nestled_db
GIN_MODE=debug
```

### 3. Install dependencies

```bash
go mod download
go mod tidy
```

## Building the Application

### Build with Make

```bash
make build
```

This compiles the application and outputs a binary to `bin/nestled`.

### Build with Go directly

```bash
go build -o bin/nestled cmd/main/main.go
```

## Running the Application

### Option 1: Local Development (with local PostgreSQL)

#### Start PostgreSQL (if not running)

```bash
# Using Homebrew on macOS
brew services start postgresql

# Or using Docker
docker run -d \
  --name nestled-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=nestled_db \
  -p 5432:5432 \
  postgres:17
```

#### Run migrations

```bash
make migrate-up
```

#### Start the application

```bash
make run
```

The API will be available at `http://localhost:8080`

### Option 2: Docker Compose (Recommended)

Start all services (PostgreSQL, migrations, and the app) with one command:

```bash
docker-compose up --build
```

This will:
1. Start PostgreSQL database
2. Run all pending migrations
3. Start the application on port 8080

Stop all services:

```bash
docker-compose down
```

View logs:

```bash
docker-compose logs -f app
```

## Running Tests

### Run all tests

```bash
make test
```

### Run specific package tests

```bash
go test ./internal/handlers -v
```

### Run tests with coverage

```bash
go test ./... -cover
```

Integration tests use [testcontainers](https://www.testcontainers.org/) to automatically spin up a PostgreSQL container, run migrations, and tear down after completion.

## Database Migrations

Migrations are SQL scripts stored in the `migrations/` directory and are automatically versioned.

### Run pending migrations

```bash
make migrate-up
```

Or with environment variables manually:

```bash
migrate -path migrations/ -database "postgres://postgres:password@localhost:5432/nestled_db?sslmode=disable" up
```

### Rollback last migration

```bash
make migrate-down
```

### Create a new migration

```bash
make migrate-create name=add_user_preferences
```

This creates two files:
- `migrations/000003_add_user_preferences.up.sql` (applies the change)
- `migrations/000003_add_user_preferences.down.sql` (reverts the change)

## Project Structure

```
nestled/
├── cmd/
│   └── main/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/                     # Configuration management
│   ├── crypto/                     # Cryptography utilities
│   ├── database/                   # Database setup
│   ├── dto/                        # Data transfer objects
│   ├── errors/                     # Error definitions
│   ├── handlers/                   # HTTP route handlers
│   ├── model/                      # Domain models
│   ├── repositories/               # Data access layer
│   ├── routes/                     # Route definitions
│   ├── services/                   # Business logic
│   └── testhelpers/                # Testing utilities
├── migrations/                     # SQL migrations
├── bin/                            # Build outputs
├── config.yml                      # Configuration template
├── docker-compose.yml              # Docker Compose setup
├── Dockerfile                      # Application Docker image
├── Makefile                        # Build and run targets
├── go.mod                          # Go module definition
└── README.md                       # This file
```

## Common Tasks

### Lint the codebase

```bash
make lint
```

Requires [golangci-lint](https://golangci-lint.run/).

### Clean build artifacts

```bash
make clean
```

### Update dependencies

```bash
make tidy
```

## API Endpoints

The API is built with Gin and provides RESTful endpoints for:

- **Health Check**: `GET /health`
- **User Management**: Registration, authentication (in development)
- **Expense Tracking**: CRUD operations for expenses
- **Budget Management**: Set and monitor budgets per category
- **Reports**: Financial summaries and insights (planned)

See the handlers in `internal/handlers/` for detailed endpoint implementations.

## Development

### Git Workflow

```bash
# Create a feature branch
git checkout -b feature/your-feature-name

# Make changes, commit, and push
git add .
git commit -m "feat: add your feature"
git push origin feature/your-feature-name

# Create a pull request
```

### Code Style

This project follows Go conventions:
- Format code with `go fmt`
- Follow [Effective Go](https://golang.org/doc/effective_go)
- Write clear, idiomatic Go

## Troubleshooting

### Migration failures in tests

If you see `"unknown driver file"` errors during tests, ensure the test helpers use the correct import path for golang-migrate (v4-compatible).

### PostgreSQL connection refused

- Verify PostgreSQL is running: `psql -U postgres`
- Check environment variables match your database setup
- Ensure the database exists: `createdb nestled_db`

## Contributing

1. Create a feature branch
2. Make your changes and write/update tests
3. Run `make test` and `make lint` to verify
4. Submit a pull request

## License

See [LICENSE](LICENSE) for details.

## Support

For issues or questions, please open an issue on the repository.