# Module management
go mod init github.com/user/project  # Create new module
go mod tidy                          # Clean up dependencies
go mod download                      # Download dependencies
go mod vendor                        # Copy deps to vendor/
go get -u ./...                      # Update all dependencies

# Building and running
go build ./...                       # Build all packages
go run .                             # Run current package
go install ./...                     # Install binaries to GOPATH/bin

# Testing
go test ./...                        # Run all tests
go test -v ./...                     # Verbose output
go test -race ./...                  # Enable race detector
go test -cover ./...                 # Show coverage percentage
go test -bench=. ./...               # Run benchmarks

# Code quality
go fmt ./...                         # Format all code
go vet ./...                         # Run static analysis
golangci-lint run                    # Run all linters