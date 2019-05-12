.DEFAULT_GOAL := install

get-pigpio:
	@echo "Fetching and installing pigpio..."
	@git clone https://github.com/joan2937/pigpio.git
	@cd pigpio \
	&& make \
	&& sudo make install
	@sudo rm -rf pigpio
	@echo "Installing pigpio completed!"
get-irslinger:
	@echo "Fetching ir-slinger..."
	@git clone https://github.com/bschwind/ir-slinger.git
	@mv ir-slinger/irslinger.h internal/accontrol/hwinterface/
	@rm -rf ir-slinger
deps: get-pigpio get-irslinger
	@echo ""
	@echo ""
	@echo "===================================="
	@echo "Dependencies fetched and installed!"
build:
	@echo "Building Goira..."
	@go build ./cmd/main.go
	@echo "Build completed!"
install: deps build
