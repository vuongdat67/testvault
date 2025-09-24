#!/bin/bash
# FileVault Test Runner Script

set -e

echo "🧪 Running FileVault Test Suite..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test categories
echo -e "${BLUE}📋 Running unit tests...${NC}"
go test -v ./test/unit/... || exit 1

echo -e "${BLUE}📋 Running integration tests...${NC}" 
go test -v ./test/integration/... || exit 1

echo -e "${BLUE}📋 Running benchmarks...${NC}"
go test -bench=. ./test/benchmarks/...

echo -e "${BLUE}📋 Running race condition tests...${NC}"
go test -race ./...

echo -e "${BLUE}📋 Generating coverage report...${NC}"
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

echo -e "${GREEN}✅ All tests passed successfully!${NC}"
echo -e "${GREEN}📊 Coverage report generated: coverage.html${NC}"