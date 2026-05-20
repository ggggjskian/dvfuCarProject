
FROM golang:1.25 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Переключили путь сборки на ./src
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./src

# Финальный образ тоже делаем на полноценном golang, отказываясь от alpine
FROM golang:1.25
WORKDIR /root/

COPY --from=builder /app/main .
EXPOSE 8000

CMD ["./main"]