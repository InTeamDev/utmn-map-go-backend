FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/${SERVICE} ./cmd/${SERVICE}

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app

ARG SERVICE

COPY --from=builder /app/bin/${SERVICE} ./

COPY config/${SERVICE}.docker.yaml ./config.yaml

ENTRYPOINT ["./adminapi", "--config", "config.yaml"]
