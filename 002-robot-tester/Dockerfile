# Используем официальный образ Go
FROM golang:1.23-alpine as builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o 002-robot-tester

# Используем более легкий образ для запуска
FROM alpine:latest

# Копируем собранное приложение из образа builder
COPY --from=builder /app/002-robot-tester .

# Открываем порт 5434
EXPOSE 5434

# Запускаем приложение
CMD ["./002-robot-tester"]
