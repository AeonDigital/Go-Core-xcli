#!/usr/bin/env bash

# ==============================================================================
# Script Name: trigger-release.sh
# Description: Triggers the GitHub Actions CI/CD pipeline by pushing an empty 
#              commit to the main branch. Supports both automatic semantic 
#              incrementing and explicit manual version targeting.
# ==============================================================================

# Exit immediately if any command fails or if an uninitialized variable is used
set -euo pipefail

# ------------------------------------------------------------------------------
# Configuration & Safety Checks
# ------------------------------------------------------------------------------
TARGET_BRANCH="main"

# Ensure the execution happens from within a Git repository repository root
if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
  echo "[-] Error: This script must be executed inside a Git repository." >&2
  exit 1
fi

# Ensure the local repository is currently on the targeted main branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "$TARGET_BRANCH" ]; then
  echo "[-] Error: You are on branch '$CURRENT_BRANCH'. Switch to '$TARGET_BRANCH' first." >&2
  exit 1
fi

# Check if there are any uncommitted local architectural changes
if ! git diff-index --quiet HEAD --; then
  echo "[-] Warning: You have uncommitted changes. Stash or commit them before triggering a release." >&2
  exit 1
fi

# ------------------------------------------------------------------------------
# Version Logic Processing
# ------------------------------------------------------------------------------
# If an argument is provided (e.g., ./trigger-release.sh v1.2.0), use it as manual override.
# Otherwise, create a technical standard commit to let the pipeline handle automatic incrementing (+1 patch).
if [ $# -gt 0 ]; then
  VERSION_INPUT="$1"
  
  # Enforce standard formatting for safety (ensures input starts with v followed by a number)
  if [[ ! "$VERSION_INPUT" =~ ^v[0-9]+ ]]; then
    echo "[-] Error: Invalid version format. Example usage: $0 v1.0.0" >&2
    exit 1
  fi
  
  COMMIT_MSG="release: $VERSION_INPUT"
  echo "[+] Preparing manual targeted release: $VERSION_INPUT"
else
  COMMIT_MSG="chore: trigger automatic release pipeline"
  echo "[+] Preparing automated semantic patch increment (+1)..."
fi

# ------------------------------------------------------------------------------
# Execution Phase
# ------------------------------------------------------------------------------
echo "[+] Creating empty infrastructure commit..."
git commit --allow-empty -m "$COMMIT_MSG"

echo "[+] Uploading deployment trigger upstream to origin/$TARGET_BRANCH..."
git push origin "$TARGET_BRANCH"

echo "[+] Success! The CI/CD workflow pipeline has been successfully triggered."
