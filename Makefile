init:
	go mod tidy
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/segmentio/golines@latest
	go install github.com/jackc/tern@latest

lint:
	golangci-lint run -j8 --enable-only gofumpt ./... --fix
	golangci-lint run -j8 --enable-only gci ./... --fix
	golangci-lint run -j8 ./...

format:
	golines --max-len=120 -w .
	go fmt ./...
