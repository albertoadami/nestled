# Build stage
FROM golang:1.26-alpine AS builder

ENV GO111MODULE=on
ENV GIN_MODE=release

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main/main.go

# Final stage
FROM alpine:3.19

LABEL maintainer="Alberto Adami <alberto@example.com>"
LABEL org.opencontainers.image.source="https://github.com/albertoadami/nestled"
LABEL org.opencontainers.image.description="Nestled API"

RUN apk add --no-cache ca-certificates curl

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/config.yml .

EXPOSE 8080

CMD ["./main"]