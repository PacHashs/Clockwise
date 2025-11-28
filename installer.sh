#!/usr/bin/env bash
set -euo pipefail

RELEASE_URL="https://github.com/your-org/clockwise/releases/latest/download/cw-linux-amd64.tar.gz"
INSTALL_DIR="/usr/local/bin"
TEMP_DIR=$(mktemp -d)

cleanup() {
  rm -rf "$TEMP_DIR"
}
trap cleanup EXIT

command -v curl >/dev/null 2>&1 || { echo "curl is required"; exit 1; }
command -v tar >/dev/null 2>&1 || { echo "tar is required"; exit 1; }

echo "Downloading Clockwise archive..."
curl -L "$RELEASE_URL" -o "$TEMP_DIR/cw.tar.gz"

echo "Extracting..."
tar -xzf "$TEMP_DIR/cw.tar.gz" -C "$TEMP_DIR"

if [ ! -f "$TEMP_DIR/cw" ]; then
  echo "Expected cw binary not found in archive"
  exit 1
fi

echo "Installing to $INSTALL_DIR (sudo may prompt you)..."
sudo install -m 0755 "$TEMP_DIR/cw" "$INSTALL_DIR/cw"

echo "Clockwise installed to $INSTALL_DIR/cw"

echo "Done. Add $INSTALL_DIR to PATH if it's not already there."
