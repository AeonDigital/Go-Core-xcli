#!/usr/bin/env bash

# ==============================================================================
# Script Name: manage-version-tag.sh
# Description: Automated Git tagging engine. Computes semantic increments 
#              based on commit history headers and publishes the next version.
# ==============================================================================

set -euo pipefail

# Immediate forced synchronization of remote server tags
git fetch --tags --force

# Captures the latest valid tag from the history or initiates the default baseline.
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo "[+] Current repository reference tag: $LATEST_TAG"

# Strips the 'v' literal prefix and extracts the numeric array from the SemVer string.
VERSION_CLEAN=${LATEST_TAG#v}
IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION_CLEAN"

# Captures the latest commit message sent to the current branch.
COMMIT_MSG=$(git log -1 --pretty=%B)
echo "[+] Triggering commit message: $COMMIT_MSG"

# Core Conditional Block for Next-Version Intelligent Calculation
if [[ "$COMMIT_MSG" =~ ^release:\ (v[0-9]+\.[0-9]+\.[0-9]+.*)$ ]]; then
  NEXT_TAG="${BASH_REMATCH[1]}"
  echo "[+] Manual release override detected via commit syntax."
elif [ "$LATEST_TAG" = "v0.0.0" ]; then
  NEXT_TAG="v0.0.1"
  echo "[+] No prior tags found. Initializing repository versioning tracking."
elif [[ "$COMMIT_MSG" == *"BREAKING CHANGE"* ]]; then
  NEXT_TAG="v$((MAJOR + 1)).0.0"
elif [[ "$COMMIT_MSG" =~ ^feat: ]]; then
  NEXT_TAG="v$MAJOR.$((MINOR + 1)).0"
else
  NEXT_TAG="v$MAJOR.$MINOR.$((PATCH + 1))"
fi

echo "[+] Calculated targets for next generation: $NEXT_TAG"

# Registers the variable in the pipeline step's native output for subsequent external use.
if [ -n "${GITHUB_OUTPUT:-}" ]; then
  echo "next_tag=$NEXT_TAG" >> "$GITHUB_OUTPUT"
fi

# Git Runner Bot Identity Settings
git config user.name "github-actions[bot]"
git config user.email "github-actions[bot]@://github.com"

# Local physical implementation and official remote push of the new milestone.
git tag "$NEXT_TAG"
git push origin "$NEXT_TAG"
echo "[+] Success: Version milestone tag $NEXT_TAG published upstream."
