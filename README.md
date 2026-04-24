# Module management
```bash
go mod init github.com/user/project  # Create new module
go mod tidy                          # Clean up dependencies
go mod download                      # Download dependencies
go mod vendor                        # Copy deps to vendor/
go get -u ./...                      # Update all dependencies
```
# Building and running
```bash
go build ./...                       # Build all packages
go run .                             # Run current package
go install ./...                     # Install binaries to GOPATH/bin
```
# Testing
```bash
go test ./...                        # Run all tests
go test -v ./...                     # Verbose output
go test -race ./...                  # Enable race detector
go test -cover ./...                 # Show coverage percentage
go test -bench=. ./...               # Run benchmarks
```
# Code quality
```bash
go fmt ./...                         # Format all code
go vet ./...                         # Run static analysis
golangci-lint run                    # Run all linters
```

# Project structure (Hexagonal/Clean Architecture)

```bash
.
├── cmd/
│   └── api/
│       └── main.go             # Nơi duy nhất thực hiện Dependency Injection
├── internal/
│   ├── domain/                 # TẦNG 1: BUSINESS LOGIC (Core)
│   │   ├── entity/             # Các đối tượng nghiệp vụ (User, Article)
│   │   │   └── user.go      
│   │   ├── values/             # Value Objects (Email, Password, etc.)
│   │   │   └── email.go     
│   │   ├── errors/             # Errors
│   ├── usecase/                # TẦNG 2: APPLICATION LOGIC - Điều phối dữ liệu (Service Layer)
│   │   ├── user/               # Chia theo module nghiệp vụ (user, article, etc.)
│   │   │   ├── input.go        # Interface/DTO đầu vào của UseCase (RegisterUserRequest, etc.)
│   │   │   ├── output.go       # Interface/DTO đầu ra của UseCase (RegisterUserResponse, etc.)
│   │   │   ├── repository.go   # Interface định nghĩa các phương thức lưu trữ (UserRepository, etc.)
│   │   │   └── service.go      # Triển khai nghiệp vụ (ví dụ: Register, Login)
│   ├── delivery/               # TẦNG 3: INTERFACE ADAPTERS (Input)
│   │   ├── grpc/               # gRPC
│   │   └── http/               # Echo Web Framework (Handlers, Middlewares)
│   │       ├── middleware/     # Middleware
│   │       └── v1/             # Versioning (API)
│   │           ├── handler/    # Handler (Request)
│   │           └── presenter/  # Presenter (Response)
│   └── repository/             # TẦNG 3: INTERFACE ADAPTERS (Output)
│       └── postgres/           # Implement Repository Interface bằng pgx (Postgres)
│           └── user_pg.go
├── pkg/                        # THƯ VIỆN DÙNG CHUNG - Các tiện ích không chứa logic nghiệp vụ
│    └── pgsql/                  
├── migrations/                 # Database migrations
├── compose.yml                 # Docker compose file
├── .env                        # Environment variables
├── .env.example                # Environment variables example
├── .gitignore                  # Git ignore file
├── Makefile                    # Makefile
├── README.md                   # README file
├── go.mod                      # Go module file
├── go.sum                      # Go sum file
├── .air.toml                   # Air configuration file
```