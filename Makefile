.PHONY: build
build:
	go build -race -o ./tmp/main ./cmd/api

.PHONY: run
run:
	go run -race ./cmd/api

.PHONY: clean
clean:
	rm -rf ./tmp

.PHONY: test
test:
	go test -race -count=1 -v -covermode=atomic ./...