#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."
rm -rf dist

./scripts/build.sh linux amd64 v9.9.9

BIN_DIR="dist/blueprint-v9.9.9-linux-amd64"
TARBALL="dist/blueprint-v9.9.9-linux-amd64.tar.gz"

if [ ! -f "$BIN_DIR/blueprint" ]; then
  echo "FAIL: expected binary at $BIN_DIR/blueprint"
  exit 1
fi

if [ ! -f "$TARBALL" ]; then
  echo "FAIL: expected tarball at $TARBALL"
  exit 1
fi

ACTUAL_VERSION=$("$BIN_DIR/blueprint" version)
if [ "$ACTUAL_VERSION" != "v9.9.9" ]; then
  echo "FAIL: expected version v9.9.9, got $ACTUAL_VERSION"
  exit 1
fi

echo "PASS: build.sh produced a working v9.9.9 linux/amd64 binary"
rm -rf dist
