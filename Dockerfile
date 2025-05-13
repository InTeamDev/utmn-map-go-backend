# Dockerfile
############################
# 1. Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Аргумент SERVICE: adminapi или publicapi
ARG SERVICE

# Собираем статический бинарь
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/${SERVICE} ./cmd/${SERVICE}

############################
# 2. Runtime stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app

# Аргумент SERVICE
ARG SERVICE

# Копируем бинарь
COPY --from=builder /app/bin/${SERVICE} ./

# Копируем соответствующий конфиг
COPY config/${SERVICE}.yaml ./config.yaml

# По умолчанию не экспонируем порт —
# docker-compose сам проксирует нужный.
# EXPOSE 8000
# EXPOSE 8001

# Запуск: binary --config config.yaml
ENTRYPOINT ["./adminapi", "--config", "config.yaml"]
# В docker-compose укажем override на publicapi
