.PHONY: build run test clean fmt vet race cover

BINARY_NAME=SCode
BINARY_PATH=./bin/$(BINARY_NAME)

build:
   go build -o $(BINARY_PATH)

run: build
   $(BINARY_PATH)

test:
   go test -v ./...

race:
   go test -race -v ./...

cover:
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out

fmt:
   go fmt ./...

vet:
   go vet ./...

clean:
   rm -rf bin/
   rm -f coverage.out

check: fmt vet test race

deps:
   go mod tidy

build-all:
   GOOS=linux GOARCH=amd64 go build -o $(BINARY_PATH)-linux-amd64
   GOOS=darwin GOARCH=amd64 go build -o $(BINARY_PATH)-darwin-amd64
   GOOS=windows GOARCH=amd64 go build -o $(BINARY_PATH)-windows-amd64.exe