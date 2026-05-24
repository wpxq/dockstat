#!/bin/bash
set -e
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
[ "$ARCH" = "x86_64" ] && ARCH="amd64"
[ "$ARCH" = "aarch64" ] && ARCH="arm64"
URL="https://github.com/wpxq/dockstat/releases/latest/download/dockstat-$OS-$ARCH"
echo "Installing $URL..."
mkdir -p "$HOME/.local/bin"
curl -L "$URL" -o "$HOME/.local/bin/dockstat"
chmod +x "$HOME/.local/bin/dockstat"
echo "dockstat successfully installed!"