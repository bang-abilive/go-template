MIGRATION_TABLE="migrations"

.DEFAULT_GOAL := help

.PHONY: help
help: ## Hiển thị danh sách target và mô tả
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; cyan = "\033[36m"; reset = "\033[0m"} /^[a-zA-Z0-9_\/-]+:.*##/ {printf "  " cyan "%-20s" reset " %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: tools
tools: ## Cài đặt công cụ cần thiết
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/mazrean/kessoku/cmd/kessoku@latest
	go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest

.PHONY: migration
migration: ## Tạo migration mới (name=...)
	@if [ -z "$(name)" ]; then \
		echo "  Migration name is required"; \
		echo "  Usage:"; \
		echo "    make migration name=create_users_table"; \
		exit 1; \
	else \
		goose create $(name) sql; \
	fi
.PHONY: migrate/up
migrate/up: ## Chạy toàn bộ migration up
	goose validate -v
	goose up -v -table=$(MIGRATION_TABLE)

.PHONY: migrate/down
migrate/down: ## Rollback migration gần nhất
	goose validate -v
	goose down -v -table=$(MIGRATION_TABLE)

.PHONY: build
build: ## Build binary API vào ./tmp/main
	go build -race -o ./tmp/main ./cmd/api

.PHONY: run
run: ## Chạy API local bằng go run
	go run -race ./cmd/api

.PHONY: clean
clean: ## Xóa thư mục tạm ./tmp
	rm -rf ./tmp

.PHONY: up
up: ## Chạy docker compose ở background
	docker compose up -d

.PHONY: down
down: ## Dừng docker compose
	docker compose down

.PHONY: test
test: ## Chạy toàn bộ test với race detector
	go test -race -count=1 -v -covermode=atomic ./...

.PHONY: coverage
coverage: ## Sinh báo cáo coverage HTML vào ./tmp
	go test -race -count=1 -v -coverprofile=./tmp/coverage.out -covermode=atomic ./...
	go tool cover -html=./tmp/coverage.out -o ./tmp/coverage.html

.PHONY: tidy
tidy: ## Xóa các dependencies không cần thiết
	go mod tidy
	go mod vendor

.PHONY: di
di: ## Generate kessoku injectors
	go generate ./...

.PHONY: seed
seed: ## Seed the database
	go run -race ./cmd/api -seed

.PHONY: struct/align
struct/align: ## Tự động sắp xếp lại struct để tối ưu hóa bộ nhớ
	fieldalignment -fix ./...

.PHONY: try
try: ## Gửi request thử endpoint register
	curl -X POST http://localhost:8080/api/v1/users/register -H "Content-Type: application/json" -d '{"email": "ndinhbang@example.com", "password": "password"}'