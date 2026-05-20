# Используем образ Go
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Сначала копируем только файлы зависимостей
COPY go.mod go.sum ./
# Скачиваем зависимости (это закэшируется, чтобы не качать каждый раз)
RUN go mod download

# Теперь копируем весь остальной код бэкенда (папки cmd и internal)
COPY . .

# СОБИРАЕМ ПРОЕКТ ИЗ ПАПКИ cmd/
RUN go build -o main ./cmd/main.go

# Финальный легкий образ
FROM alpine:latest
WORKDIR /root/

# Забираем собранный бинарник из первого этапа
COPY --from=builder /app/main .

# Открываем порт
EXPOSE 8000

# Запускаем
CMD ["./main"]