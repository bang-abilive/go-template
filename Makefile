.PHONY: tools
tools:
	go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: migration
migration:
	@if [ -z "$(name)" ]; then \
		echo "  Migration name is required"; \
		echo "  Usage:"; \
		echo "    make migration name=create_users_table"; \
		exit 1; \
	else \
		goose create $(name) sql; \
	fi
.PHONY: migrate/up
migrate/up:
	goose validate -v
	goose up -v

.PHONY: migrate/down
migrate/down:
	goose validate -v
	goose down -v

.PHONY: build
build:
	go build -race -o ./tmp/main ./cmd/api

.PHONY: run
run:
	go run -race ./cmd/api

.PHONY: clean
clean:
	rm -rf ./tmp

.PHONY: up
up:
	docker compose up -d

.PHONY: down
down:
	docker compose down

.PHONY: test
test:
	go test -race -count=1 -v -covermode=atomic ./...

.PHONY: coverage
coverage:
	go test -race -count=1 -v -coverprofile=./tmp/coverage.out -covermode=atomic ./...
	go tool cover -html=./tmp/coverage.out -o ./tmp/coverage.html

.PHONY: try
try:
	curl -X POST http://localhost:8080/api/v1/users/register -H "Content-Type: application/json" -d '{"email": "ndinhbang@example.com", "password": "password"}'