#!/usr/bin/env bash

# ==============================================================================
# Script Name: load-metadata.sh
# Description: Automatically extracts Go version/module name from go.mod and
#              loads custom pipeline configurations into GITHUB_ENV.
# ==============================================================================

set -euo pipefail

# 1. Automatic extraction of go.mod from the project root
if [ -f "go.mod" ]; then
  MODULE_PREFIX=$(awk '/^module / {print $2}' go.mod)
  GO_VERSION=$(awk '/^go / {print $2}' go.mod | cut -d. -f1,2)
  
  echo "MODULE_PREFIX=$MODULE_PREFIX" >> "$GITHUB_ENV"
  echo "GO_VERSION=$GO_VERSION" >> "$GITHUB_ENV"
  echo "[+] Auto-detected from go.mod: Module=$MODULE_PREFIX, Go=$GO_VERSION"
else
  echo "[-] Error: go.mod not found at root level." >&2
  exit 1
fi

# 2. Reading the custom configuration file (.github/config.txt)
CONFIG_FILE=".github/config.txt"
if [ -f "$CONFIG_FILE" ]; then
  echo "[+] Custom configuration file found. Loading additional variables..."
  while IFS= read -r line || [ -n "$line" ]; do
    # Ignore comments and blank lines
    [[ "$line" =~ ^[[:space:]]*# ]] && continue
    [[ -z "${line//[:space:]/}" ]] && continue
    
    echo "$line" >> "$GITHUB_ENV"
    echo "[+] Loaded custom config: $line"
  done < "$CONFIG_FILE"
else
  echo "[-] Warning: .github/config.txt missing. Applying default pipeline thresholds."
  echo "COVERAGE_THRESHOLD=80" >> "$GITHUB_ENV"
fi
