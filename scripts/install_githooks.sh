#!/bin/sh

echo "Installing pre-push hook..."
cp scripts/githooks/pre_push.sh .git/hooks/pre-push
chmod +x .git/hooks/pre-push
echo "✅ Pre-push hook installed successfully!"
