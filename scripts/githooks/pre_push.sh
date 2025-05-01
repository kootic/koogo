#!/bin/sh

set -eo pipefail

# Get the path to the pre-push hook
HOOK_PATH="$(git rev-parse --git-dir)/hooks/pre-push"
SCRIPT_PATH="$(realpath "$0")"

# Check if the hook has been updated
if [ -f "$HOOK_PATH" ] && ! cmp -s "$SCRIPT_PATH" "$HOOK_PATH"; then
  echo "⚠️ The pre-push hook has been updated!"
  echo "Please run ./scripts/install_githooks.sh to update your local hook:"
  echo "This will ensure you have the latest version of the pre-push hook."
  exit 1
fi

task lint

task generate

# Check if there are any changes in the docs directory
if git diff --name-only | grep -q "^docs/"; then
  echo "📝 Swagger docs have been updated. Staging changes..."
  git add docs/
  echo "❗ Please review the staged changes and commit them again."
  exit 1
else
  echo "✅ No changes detected in Swagger docs."
fi

task test

# Check if there are any changes that are not committed and warn the user
# Save git diff result to show the user what changes are uncommitted
UNCOMMITTED_CHANGES=$(git diff --name-only)
if [ -n "$UNCOMMITTED_CHANGES" ]; then
  echo "❗ There are uncommitted changes:"
  echo "$UNCOMMITTED_CHANGES"
  echo ""
  echo "Are you sure you want to push?"
  read -p "Continue? (y/n): " confirm
  if [ "$confirm" != "y" ]; then
    echo "❌ Push aborted."
    exit 1
  fi
else
  echo "✅ No uncommitted changes detected."
fi
