APP_NAME=nestled
MAIN=cmd/main/main.go

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