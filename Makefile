.PHONY: build
build:
	go build -race -o ./tmp/main ./cmd/api

.PHONY: run
run:
	go run -race ./cmd/api

.PHONY: clean
clean:
	rm -rf ./tmp