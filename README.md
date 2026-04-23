## Module management

```bash
go mod init github.com/user/project  # Create new module
go mod tidy                          # Clean up dependencies
go mod download                      # Download dependencies
go mod vendor                        # Copy deps to vendor/
go get -u ./...                      # Update all dependencies
```

## Building and running

```bash
go build ./...                       # Build all packages
go run .                             # Run current package
go install ./...                     # Install binaries to GOPATH/bin
```

## Testing

```bash
go test ./...                        # Run all tests
go test -v ./...                     # Verbose output
go test -race ./...                  # Enable race detector
go test -cover ./...                 # Show coverage percentage
go test -bench=. ./...               # Run benchmarks
```

## Code quality

```bash
go fmt ./...                         # Format all code
go vet ./...                         # Run static analysis
golangci-lint run                    # Run all linters
```

## Project structure (Hexagonal/Clean Architecture)

```text
.
├── cmd/
│   └── api/
│       └── main.go           # Entry point: Nơi khởi tạo Container và chạy Server
├── internal/
│   ├── domain/               # TẦNG 1: BUSINESS LOGIC (Core) - Không phụ thuộc bên ngoài
│   │   └── values/           # Value Objects (Email, Password, etc.)
│   │   ├── entity/           # Các đối tượng nghiệp vụ (User, Article) và Value Objects
│   │   └── repository/       # Interface định nghĩa các phương thức lưu trữ
│   ├── usecase/              # TẦNG 2: APPLICATION LOGIC - Điều phối dữ liệu (Service Layer)
│   │   ├── dto/              # Data Transfer Objects (Request/Response)
│   │   └── user_uc.go        # Triển khai nghiệp vụ (ví dụ: Register, Login)
│   ├── infrastructure/       # TẦNG 3: EXTERNAL TOOLS - Triển khai kỹ thuật chi tiết
│   │   └── persistence/      # Implement Repository Interface bằng pgx (Postgres)
│   │       └── pgx_user.go
│   └── interface/            # TẦNG 4: DELIVERY MECHANISM - Giao tiếp với thế giới
│       ├── http/             # Echo Web Framework (Handlers, Middlewares)
│       │   ├── middleware/
│       │   ├── handler/
│       │   ├── presenter/
│       │   └── router/
│       └── grpc/             # gRPC
├── pkg/                      # THƯ VIỆN DÙNG CHUNG - Các tiện ích không chứa logic nghiệp vụ
│   └── pgsql/                # Khởi tạo pgxpool
├── .env                      # Biến môi trường
├── go.mod                    # Module dependencies
└── go.sum                    # Module checksums
```