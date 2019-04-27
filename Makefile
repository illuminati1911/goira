install-deps:
	git clone https://github.com/joan2937/pigpio.git
	cd pigpio
	make
	sudo make install
build:
	go build ./cmd/main.go