#!/usr/bin/env bash
set -euo pipefail

REPO="hxsggsz/blueprint-cli"
INSTALL_DIR="${BLUEPRINT_INSTALL_DIR:-$HOME/.local/bin}"
CONFIG_DIR="${BLUEPRINT_CONFIG_DIR:-$HOME/.config/blueprint}"
CONFIG_FILE="${CONFIG_DIR}/config.json"

detect_platform() {
  local os arch
  os=$(uname -s)
  arch=$(uname -m)

  case "$os" in
    Linux) os="linux" ;;
    Darwin) os="darwin" ;;
    *)
      echo "Error: unsupported OS '$os'. Supported: Linux, Darwin." >&2
      exit 1
      ;;
  esac

  case "$arch" in
    x86_64|amd64) arch="amd64" ;;
    arm64|aarch64) arch="arm64" ;;
    *)
      echo "Error: unsupported architecture '$arch'. Supported: amd64, arm64." >&2
      exit 1
      ;;
  esac

  echo "${os} ${arch}"
}

resolve_version() {
  if [ -n "${BLUEPRINT_VERSION:-}" ]; then
    echo "$BLUEPRINT_VERSION"
    return
  fi

  local latest
  latest=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')

  if [ -z "$latest" ]; then
    echo "Error: could not resolve the latest release from GitHub API." >&2
    exit 1
  fi

  echo "$latest"
}

write_config() {
  local templates_path="$1"
  mkdir -p "$CONFIG_DIR"

  if [ -f "$CONFIG_FILE" ]; then
    echo "Config already exists at ${CONFIG_FILE}."
    echo "Make sure 'template_path' points to ${templates_path}."
    return
  fi

  cat > "$CONFIG_FILE" <<EOF
{
  "template_path": "${templates_path}",
  "ignore_file_paths": []
}
EOF

  echo "Wrote config to ${CONFIG_FILE}"
}

setup_templates() {
  if [ "${BLUEPRINT_SKIP_TEMPLATES:-}" = "1" ]; then
    echo "Skipping template setup."
    return
  fi

  local templates_path="${BLUEPRINT_TEMPLATES_PATH:-}"

  if [ -z "$templates_path" ]; then
    echo
    echo "Where are your templates stored?"
    printf "Templates directory path (press Enter to skip): "
    # Read from /dev/tty so the prompt works when the script is
    # piped (e.g. `curl ... | bash`), where stdin is the script itself.
    if [ -t 0 ]; then
      read -r templates_path
    elif [ -c /dev/tty ]; then
      read -r templates_path < /dev/tty
    else
      templates_path=""
    fi
  fi

  if [ -z "$templates_path" ]; then
    echo "No templates path provided. Set 'template_path' in ${CONFIG_FILE} later."
    return
  fi

  if [ ! -d "$templates_path" ]; then
    echo "Warning: '${templates_path}' is not an existing directory." >&2
  fi

  write_config "$templates_path"
}

main() {
  local platform
  platform=$(detect_platform) || exit 1
  read -r OS ARCH <<< "$platform"
  VERSION=$(resolve_version)

  if [[ "$VERSION" != v* ]]; then
    VERSION="v${VERSION}"
  fi

  local asset="blueprint-${VERSION}-${OS}-${ARCH}.tar.gz"
  local base_url="https://github.com/${REPO}/releases/download/${VERSION}"
  tmp_dir=$(mktemp -d)
  staged_bin=""
  trap 'rm -f "$staged_bin"; rm -rf "$tmp_dir"' EXIT

  echo "Installing blueprint ${VERSION} (${OS}/${ARCH})..."

  if ! curl -fsSL -o "${tmp_dir}/${asset}" "${base_url}/${asset}"; then
    echo "Error: release asset not found: ${base_url}/${asset}" >&2
    echo "Check that version '${VERSION}' exists: https://github.com/${REPO}/releases" >&2
    exit 1
  fi

  if ! curl -fsSL -o "${tmp_dir}/checksums.txt" "${base_url}/checksums.txt"; then
    echo "Error: checksums file not found: ${base_url}/checksums.txt" >&2
    echo "Check that version '${VERSION}' exists: https://github.com/${REPO}/releases" >&2
    exit 1
  fi

  echo "Verifying checksum..."
  (cd "$tmp_dir" && grep -F -- "$asset" checksums.txt | sha256sum -c -) || {
    echo "Error: checksum verification failed for ${asset}. Aborting install." >&2
    exit 1
  }

  tar -xzf "${tmp_dir}/${asset}" -C "$tmp_dir"

  mkdir -p "$INSTALL_DIR"

  # Stage the new binary inside INSTALL_DIR itself (not the /tmp-based
  # mktemp dir) so the final `mv` below is a same-filesystem rename, which
  # is atomic. A cross-filesystem mv degrades to copy+unlink and can leave
  # a truncated/corrupt binary in place if interrupted mid-copy.
  staged_bin=$(mktemp "${INSTALL_DIR}/.blueprint.XXXXXX")
  cp "${tmp_dir}/blueprint-${VERSION}-${OS}-${ARCH}/blueprint" "$staged_bin"
  chmod +x "$staged_bin"
  mv "$staged_bin" "${INSTALL_DIR}/blueprint"

  echo "Installed blueprint ${VERSION} to ${INSTALL_DIR}/blueprint"

  case ":$PATH:" in
    *":${INSTALL_DIR}:"*) ;;
    *)
      echo "Warning: ${INSTALL_DIR} is not in your PATH. Add it with:"
      echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
      ;;
  esac

  setup_templates
}

main
