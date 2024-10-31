build:
	go build -o ./bin/SCode

run: build
	./bin/SCode

test:
	go test -v ./...

.PHONY: build run test