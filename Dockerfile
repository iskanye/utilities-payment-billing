# Сборка
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/billing ./cmd/billing/main.go
RUN go build -o /bin/migrator ./cmd/migrator/main.go

# Запуск
FROM alpine
USER root
WORKDIR /home/app
COPY --from=builder /bin/billing ./
COPY --from=builder /bin/migrator ./
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations
RUN mkdir storage
RUN ./migrator --uri=postgres:postgres@localhost:5430/postgres --migrations-path=./migrations
ENTRYPOINT ["./billing"]
CMD ["-config", "./config/dev.yaml"]