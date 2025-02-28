#!/bin/bash

# SpringWell CLI installer
set -e

# Determine system type
SYSTEM=$(uname -s)
ARCH=$(uname -m)

# Translate architecture names
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

# Set platform-specific variables
case "$SYSTEM" in
    "Linux")
        PLATFORM="linux"
        INSTALL_DIR="/usr/local/bin"
        ;;
    "Darwin")
        PLATFORM="darwin"
        INSTALL_DIR="/usr/local/bin"
        ;;
    *)
        echo "Unsupported system: $SYSTEM"
        exit 1
        ;;
esac

# Binary name
BINARY_NAME="springwell"
BINARY_PATH="$INSTALL_DIR/$BINARY_NAME"

# Print banner
echo "╔═══════════════════════════════════════════╗"
echo "║        SpringWell CLI Installer           ║"
echo "╚═══════════════════════════════════════════╝"
echo

# Check if we have the binary already built
if [ -f "$BINARY_NAME" ]; then
    echo "✓ Found local binary, installing..."
    sudo cp "$BINARY_NAME" "$BINARY_PATH"
    sudo chmod +x "$BINARY_PATH"
else
    echo "Building SpringWell CLI from source..."
    # Check if Go is installed
    if ! command -v go &>/dev/null; then
        echo "❌ Go is not installed. Please install Go first."
        exit 1
    fi
    
    # Build the binary
    echo "⚙️ Building for $PLATFORM-$ARCH..."
    GOOS=$PLATFORM GOARCH=$ARCH go build -o "$BINARY_NAME" ./cmd/springwell
    
    # Install the binary
    echo "📦 Installing to $INSTALL_DIR..."
    sudo cp "$BINARY_NAME" "$BINARY_PATH"
    sudo chmod +x "$BINARY_PATH"
    
    # Clean up
    rm "$BINARY_NAME"
fi

echo
echo "🎉 Installation complete! SpringWell CLI is now available."
echo "Run 'springwell --help' to get started." 