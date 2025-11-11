# Этап 1: Сборка
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# Установка goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go

# Этап 2: Запуск
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
# Копируем goose
COPY --from=builder /go/bin/goose .
# Копируем миграции и .env
COPY migrations ./migrations
COPY .env .
EXPOSE 8080
# Запускаем миграции, а затем основное приложение
CMD ["sh", "-c", "/app/goose up && /app/main"]