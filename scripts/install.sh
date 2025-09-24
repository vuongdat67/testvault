#!/bin/bash
# FileVault Installation Script

set -e

# Configuration
BINARY_NAME="filevault"
INSTALL_DIR="/usr/local/bin"
VERSION="latest"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64) ARCH="arm64" ;;
    *) echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}"; exit 1 ;;
esac

echo -e "${BLUE}🔐 FileVault Installer${NC}"
echo -e "${BLUE}OS: $OS${NC}"
echo -e "${BLUE}Architecture: $ARCH${NC}"
echo ""

# Check if running as root for system-wide installation
if [[ $EUID -ne 0 ]] && [[ "$INSTALL_DIR" == "/usr/local/bin" ]]; then
    echo -e "${YELLOW}⚠️  Installing to user directory instead of system-wide${NC}"
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

# Download or copy binary
if [[ -f "build/filevault-$OS-$ARCH" ]]; then
    echo -e "${BLUE}📦 Installing from local build...${NC}"
    cp "build/filevault-$OS-$ARCH" "$INSTALL_DIR/$BINARY_NAME"
elif [[ -f "$BINARY_NAME" ]]; then
    echo -e "${BLUE}📦 Installing local binary...${NC}"
    cp "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
else
    echo -e "${RED}❌ Binary not found. Please build first with: ./scripts/build.sh${NC}"
    exit 1
fi

# Make executable
chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Verify installation
if command -v "$BINARY_NAME" >/dev/null 2>&1; then
    VERSION_OUTPUT=$($BINARY_NAME version 2>/dev/null | head -n1 || echo "unknown")
    echo -e "${GREEN}✅ FileVault installed successfully!${NC}"
    echo -e "${GREEN}📍 Location: $INSTALL_DIR/$BINARY_NAME${NC}"
    echo -e "${GREEN}🔢 Version: $VERSION_OUTPUT${NC}"
    echo ""
    echo -e "${BLUE}💡 Usage examples:${NC}"
    echo -e "   filevault encrypt document.pdf"
    echo -e "   filevault decrypt document.pdf.enc"
    echo -e "   filevault --help"
else
    echo -e "${YELLOW}⚠️  Installation completed but binary not in PATH${NC}"
    echo -e "${YELLOW}📍 Binary location: $INSTALL_DIR/$BINARY_NAME${NC}"
    echo -e "${YELLOW}💡 Add to PATH: export PATH=\"$INSTALL_DIR:\$PATH\"${NC}"
fi