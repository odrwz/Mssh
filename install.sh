#!/bin/bash
set -e

REPO="odrwz/CLImssh"
BINARY_NAME="climssh"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Build binary name
BINARY="${BINARY_NAME}-${OS}-${ARCH}"

# Download URL (requires GitHub Release)
URL="https://github.com/${REPO}/releases/latest/download/${BINARY}"

echo "Downloading ${BINARY} from ${URL}..."
curl -L -o "/tmp/${BINARY_NAME}" "${URL}" || {
    echo "Download failed. Trying to build from source..."
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo "Go is not installed. Please install Go first:"
        echo "  brew install go"
        exit 1
    fi
    
    # Build from source
    echo "Building from source..."
    go install "github.com/${REPO}@latest"
    echo "Installed to $(go env GOPATH)/bin/${BINARY_NAME}"
    exit 0
}

# Install
chmod +x "/tmp/${BINARY_NAME}"
echo "Installing to ${INSTALL_DIR} (may require sudo)..."
sudo mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"

echo "✅ ${BINARY_NAME} installed successfully!"
echo "Run '${BINARY_NAME}' to start."
