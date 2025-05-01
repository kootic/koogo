#!/bin/sh

echo "Installing pre-push hook..."
cp scripts/pre_push_hook.sh .git/hooks/pre-push
chmod +x .git/hooks/pre-push
echo "✅ Pre-push hook installed successfully!"
