init:
	go mod tidy
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/segmentio/golines@latest
	go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
	go install github.com/jackc/tern@latest

lint:
	# gofumpt is a tool for formatting the code, strickter that go fmt.
	# To find out how to set up autoformatting in your IDE, visit
	# https://github.com/mvdan/gofumpt#visual-studio-code
	golangci-lint run -j8 --enable-only gofumpt ./... --fix
	golangci-lint run -j8 --enable-only gci ./... --fix
	golangci-lint run -j8 ./...

format:
	golines --max-len=120 -w .
	go fmt ./...
