#!/usr/bin/env bash

# ==============================================================================
# Script Name: check-entrypoints.sh
# Description: Evaluates the entrypoints tracking file to determine if the 
#              repository should be compiled as an executable binary package
#              or treated strictly as a reusable code library.
# ==============================================================================

set -euo pipefail

ENTRY_FILE=".github/entrypoints.txt"

# Checks if the file exists and contains at least one line of useful text (not just comments or whitespace).
if [ -f "$ENTRY_FILE" ] && grep -v '^#' "$ENTRY_FILE" | grep -q '[^[:space:]]'; then
  echo "[+] Entry points file found with active paths. Enabling GoReleaser engine."
  HAS_MAIN="true"
else
  echo "[-] Entry points file is missing, empty, or contains only comments. Treating repository strictly as a library."
  HAS_MAIN="false"
fi

# Injects the variable into the native GitHub Actions output (if running within the pipeline).
if [ -n "${GITHUB_OUTPUT:-}" ]; then
  echo "has_main=$HAS_MAIN" >> "$GITHUB_OUTPUT"
fi
