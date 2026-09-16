#!/usr/bin/env bash
set -euo pipefail

if [ $# -ne 3 ]; then
  echo "Usage: $0 <goos> <goarch> <version>" >&2
  echo "Example: $0 linux amd64 v1.2.0" >&2
  exit 1
fi

GOOS="$1"
GOARCH="$2"
VERSION="$3"

cd "$(dirname "$0")/.."

OUT_DIR="dist/blueprint-${VERSION}-${GOOS}-${GOARCH}"
mkdir -p "$OUT_DIR"

echo "Building blueprint ${VERSION} for ${GOOS}/${GOARCH}..."
GOOS="$GOOS" GOARCH="$GOARCH" go build \
  -ldflags "-X main.Version=${VERSION}" \
  -o "${OUT_DIR}/blueprint" \
  .

tar -czf "dist/blueprint-${VERSION}-${GOOS}-${GOARCH}.tar.gz" -C "dist" "blueprint-${VERSION}-${GOOS}-${GOARCH}"

echo "Built ${OUT_DIR}/blueprint and dist/blueprint-${VERSION}-${GOOS}-${GOARCH}.tar.gz"
