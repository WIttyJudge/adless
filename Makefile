NAME=barrier

SRC_PATH=./cmd/$(NAME)
BIN_PATH=./build/$(NAME)

VERSION=$(shell git describe --abbrev=0 2>/dev/null || echo -n "unknown")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo -n "unknown")
BUILD_DATE=$(shell date +%FT%T%z)

LDFLAGS=-w -s \
		-X main.version=$(VERSION) \
		-X main.gitCommit=$(GIT_COMMIT) \
		-X main.buildDate=$(BUILD_DATE)

.PHONY: run build test coverage clean

run:
	go run $(SRC_PATH)

build: build-linux

build-linux:
	GOOS=linux go build -ldflags "$(LDFLAGS)" -o $(BIN_PATH) $(SRC_PATH)

test:
	go test -v ./...

coverage:
	go test ./... -coverprofile=coverage.out =covermode=atomic
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf "${BIN_PATH}"
