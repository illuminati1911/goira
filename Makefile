.DEFAULT_GOAL := all

get-pigpio:
	@echo "Fetching and installing pigpio..."
	@git clone https://github.com/joan2937/pigpio.git
	@cd pigpio \
	&& make \
	&& sudo make install
	@sudo rm -rf pigpio
	@echo "Installing pigpio completed!"
update-irslinger:
	@echo "Fetching ir-slinger..."
	@git clone https://github.com/bschwind/ir-slinger.git
	@rm -f internal/accontrol/hwinterface/irslinger.h
	@mv ir-slinger/irslinger.h internal/accontrol/hwinterface/
	@rm -rf ir-slinger
deps: get-pigpio
	@echo ""
	@echo ""
	@echo "===================================="
	@echo "Dependencies fetched and installed!"
build:
	@echo "Building Goira..."
	@go build ./cmd/main.go
	@echo "Build completed!"
all: deps build test
test:
	@go test -race ./...
	@go vet ./...
	@bash -c "diff -u <(echo -n) <(gofmt -d -s .)"