init:
	go mod tidy
	go mod download
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go install github.com/segmentio/golines@latest
	go install github.com/jackc/tern@latest

lint:
	golangci-lint run --fix

format:
	golangci-lint fmt
