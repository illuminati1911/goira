.DEFAULT_GOAL := all
all: build test
build:
	@echo "Building Goira..."
	@go build -o goira ./cmd/main.go
	@echo "Build completed!"
test:
	@go test -race ./...
	@go vet ./...
	@bash -c "diff -u <(echo -n) <(gofmt -d -s .)"