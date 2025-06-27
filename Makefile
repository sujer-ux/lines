build:
	go build -v -o bin/lines .

install-linux:
	sudo cp bin/lines /usr/bin/lines

uninstall-linux:
	sudo rm -f /usr/bin/lines

get:
	go mod tidy