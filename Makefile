

install:
	go mod tidy
	go mod download
	go mod verify

run:
	go version
	go run bloom.go