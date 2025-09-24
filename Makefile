.PHONY: build test clean install dev-setup

BINARY_NAME=filevault
BUILD_DIR=build

# Development setup
dev-setup:
	go mod tidy
	go mod download

# Build for current platform
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/filevault

# Build for all platforms
build-all:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/filevault
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/filevault
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/filevault

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v ./test/unit/... ./test/integration/... ./test/benchmarks/...

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install locally
install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Security scan
security-scan:
	go list -json -m all | nancy sleuth
	gosec ./...

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Lint code
lint:
	golangci-lint run