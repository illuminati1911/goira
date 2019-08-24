.DEFAULT_GOAL := all
all: build test
.PHONY: build
build:
	@echo "Building Goira..."
	@go build -o build/goira ./cmd/main.go
	@echo "Build completed!"
frontend:
	@echo "Fetching goira frontend..."
	@git clone https://github.com/illuminati1911/goira-frontend.git
	@echo "Fetching finished"
	@echo "Building frontend"
	@npm install --prefix goira-frontend/
	@npm run build --prefix goira-frontend/
	@echo "Combining with backend"
	@mkdir -p build/frontend
	@cp -R goira-frontend/build/* build/frontend
	@echo "Cleanup..."
	@rm -rf goira-frontend
	@echo "Build completed!"
full: build frontend
test:
	@go test -race ./...
	@go vet ./...
	@bash -c "diff -u <(echo -n) <(gofmt -d -s .)"
clean:
	@rm -rf build/
