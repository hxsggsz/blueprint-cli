#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

TMP_INSTALL_DIR=$(mktemp -d)
TMP_CONFIG_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_INSTALL_DIR" "$TMP_CONFIG_DIR"' EXIT

BLUEPRINT_INSTALL_DIR="$TMP_INSTALL_DIR" \
BLUEPRINT_CONFIG_DIR="$TMP_CONFIG_DIR" \
BLUEPRINT_SKIP_TEMPLATES=1 \
  ./install.sh

if [ ! -x "$TMP_INSTALL_DIR/blueprint" ]; then
  echo "FAIL: expected executable at $TMP_INSTALL_DIR/blueprint"
  exit 1
fi

"$TMP_INSTALL_DIR/blueprint" version

echo "PASS: install.sh installed a working blueprint binary into a custom dir"
