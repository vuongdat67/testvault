#!/bin/bash

# FileVault Cross-Platform Build Script
# This script builds FileVault for multiple operating systems and architectures

set -e

# Configuration
BINARY_NAME="filevault"
BUILD_DIR="build"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE} -s -w"

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# Platform configurations
# Format: "GOOS GOARCH SUFFIX"
PLATFORMS=(
    "windows amd64 .exe"
    "windows arm64 .exe"
    "linux amd64 "
    "linux arm64 "
    "darwin amd64 "
    "darwin arm64 "
    "freebsd amd64 "
    "openbsd amd64 "
)

echo "🚀 Building FileVault v${VERSION} (${COMMIT}) for multiple platforms..."
echo "📅 Build date: ${DATE}"
echo ""

# Build for each platform
for platform in "${PLATFORMS[@]}"; do
    read -r GOOS GOARCH SUFFIX <<< "$platform"
    
    OUTPUT_NAME="${BUILD_DIR}/${BINARY_NAME}-${GOOS}-${GOARCH}${SUFFIX}"
    
    echo "🔨 Building for ${GOOS}/${GOARCH}..."
    
    env GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build \
        -ldflags "${LDFLAGS}" \
        -o "${OUTPUT_NAME}" \
        ./cmd/filevault
    
    if [ $? -eq 0 ]; then
        FILE_SIZE=$(du -h "${OUTPUT_NAME}" | cut -f1)
        echo "   ✅ Built: ${OUTPUT_NAME} (${FILE_SIZE})"
    else
        echo "   ❌ Failed to build for ${GOOS}/${GOARCH}"
        exit 1
    fi
done

echo ""
echo "📦 Creating additional packages..."

# Create checksums
echo "🔐 Generating checksums..."
cd ${BUILD_DIR}
shasum -a 256 * > SHA256SUMS
cd ..
echo "   ✅ SHA256 checksums created"

# Create archives for easier distribution
echo "📋 Creating release archives..."
cd ${BUILD_DIR}

# Windows ZIP archives
for file in *windows*.exe; do
    if [ -f "$file" ]; then
        zip "${file%.exe}.zip" "$file" SHA256SUMS
        echo "   ✅ Created ${file%.exe}.zip"
    fi
done

# Unix tar.gz archives  
for file in *linux* *darwin* *freebsd* *openbsd*; do
    if [ -f "$file" ] && [[ "$file" != *.zip ]]; then
        tar -czf "${file}.tar.gz" "$file" SHA256SUMS
        echo "   ✅ Created ${file}.tar.gz"
    fi
done

cd ..

echo ""
echo "📊 Build Summary:"
echo "=================="
ls -lh ${BUILD_DIR}/ | grep -E '(filevault|SHA256SUMS)'

echo ""
echo "🎉 Cross-platform build completed successfully!"
echo "   📁 Binaries available in: ${BUILD_DIR}/"
echo "   🔒 Checksums: ${BUILD_DIR}/SHA256SUMS"
echo ""

# Optional: Upload to release (if running in CI)
if [ "${CI}" = "true" ] && [ -n "${GITHUB_TOKEN}" ]; then
    echo "🚢 CI environment detected. Binaries ready for release upload."
fi
