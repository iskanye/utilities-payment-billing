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
COPY --from=builder /app/docker-entrypoint.sh ./
ENTRYPOINT ["/bin/sh", "./docker-entrypoint.sh"]