#!/usr/bin/env bash

# ==============================================================================
# Script Name: pre-build-verify.sh
# Description: Validates compilation integrity across all internal submodules
#              to prevent broken code structures from triggering releases.
# ==============================================================================

set -euo pipefail

echo "========================================="
echo "Running pre-flight verification build..."
echo "========================================="

# Recursively searches for all go.mod files.
find . -name "go.mod" | while read -r modfile; do
  dir=$(dirname "$modfile")
  echo "[+] Validating Go compilation in: $dir"
  
  # Enter the subdirectory and perform a clean build of the local dependencies.
  cd "$dir"
  go build ./...
  
  # Returns to the starting point of the loop execution.
  cd - > /dev/null
done
