# Start from golang base image
FROM golang:1.26-alpine as builder
# Enable go modules
ENV GO111MODULE=on
# Enable release mode
ENV GIN_MODE=release

LABEL maintainer="Alberto Adami <alberto@example.com>"
LABEL org.opencontainers.image.source="https://github.com/albertoadami/nestled"
LABEL org.opencontainers.image.description="Nestled API"

# Install bash. (alpine image does not have bash in it)
RUN apk update && apk add git && apk add bash

# Set current working directory
WORKDIR /app

# Note here: To avoid downloading dependencies every time we
# build image. Here, we are caching all the dependencies by
# first copying go.mod and go.sum files and downloading them,
# to be used every time we build the image if the dependencies
# are not changed.

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download all dependencies.
RUN go mod download

# Now, copy the source code
COPY . .

# Note here: CGO_ENABLED is disabled for cross system compilation
# It is also a common best practise.

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main/main.go

EXPOSE 8080
# Run executable
CMD ["./main"]