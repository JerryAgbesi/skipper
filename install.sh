#!/usr/bin/env sh
set -e

REPO="JerryAgbesi/skipper"
BIN_NAME="skipper"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux)  OS="linux"  ;;
  darwin) OS="darwin" ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Fetch the latest release tag from GitHub
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' \
  | sed 's/.*"tag_name": *"\(.*\)".*/\1/')

if [ -z "$LATEST" ]; then
  echo "Could not determine the latest release."
  exit 1
fi

# Sanity-check: version must look like v0.1.2 or 0.1.2
case "$LATEST" in
  v[0-9]*.[0-9]*.[0-9]*|[0-9]*.[0-9]*.[0-9]*) ;;
  *)
    echo "Unexpected version string: '$LATEST'. Aborting."
    exit 1
    ;;
esac

VERSION="${LATEST#v}"
ARCHIVE="${BIN_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${ARCHIVE}"

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading ${BIN_NAME} ${LATEST} (${OS}/${ARCH})..."
curl -fsSL "$URL" -o "${TMP_DIR}/${ARCHIVE}"

tar -xzf "${TMP_DIR}/${ARCHIVE}" -C "$TMP_DIR"

echo "Installing to ${INSTALL_DIR}/${BIN_NAME} ..."
if install -m 755 "${TMP_DIR}/${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}" 2>/dev/null; then
  :
elif command -v sudo >/dev/null 2>&1; then
  echo "(requires sudo)"
  sudo install -m 755 "${TMP_DIR}/${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}"
else
  echo "Permission denied. Re-run as root or install manually:"
  echo "  sudo install -m 755 ${TMP_DIR}/${BIN_NAME} ${INSTALL_DIR}/${BIN_NAME}"
  exit 1
fi

"${INSTALL_DIR}/${BIN_NAME}" --version
