#!/bin/bash
# FileVault Release Script

set -e

VERSION=${1:-"v1.0.0"}
BUILD_DIR="build"

echo "🚀 Creating FileVault release $VERSION..."

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# Validate version format
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}❌ Invalid version format. Use: v1.0.0${NC}"
    exit 1
fi

# Run tests
echo -e "${BLUE}🧪 Running tests...${NC}"
./scripts/test.sh

# Build all platforms
echo -e "${BLUE}🔨 Building all platforms...${NC}"
./scripts/build.sh

# Create release archive
echo -e "${BLUE}📦 Creating release archive...${NC}"
cd $BUILD_DIR
zip -r filevault-$VERSION.zip ./*
cd ..

# Generate checksums
echo -e "${BLUE}🔐 Generating checksums...${NC}"
cd $BUILD_DIR
sha256sum * > SHA256SUMS
cd ..

echo -e "${GREEN}✅ Release $VERSION created successfully!${NC}"
echo -e "${GREEN}📁 Files available in: $BUILD_DIR/${NC}"
echo -e "${GREEN}📦 Archive: $BUILD_DIR/filevault-$VERSION.zip${NC}"