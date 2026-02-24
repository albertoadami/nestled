FROM alpine:latest

RUN apk add --no-cache ca-certificates curl

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate

COPY migrations/ /migrations

WORKDIR /migrations

ENTRYPOINT ["migrate"]
CMD ["--help"]