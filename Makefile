include .env

APP_NAME=nestled
MAIN=cmd/main/main.go
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: run build clean test lint

run:
	go run $(MAIN)

build:
	go build -o bin/$(APP_NAME) $(MAIN)

clean:
	rm -rf bin/

test:
	go test ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

migrate-up:
	migrate -path migrations/ -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations/ -database "$(DB_URL)" down

migrate-create:
	migrate create -ext sql -dir migrations/ -seq $(name)
