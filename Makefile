APP_NAME := homevision
BIN_DIR  := bin

.PHONY: run build test fmt vet tidy clean

run:
	go run ./cmd/api

build:
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/api

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf $(BIN_DIR)
