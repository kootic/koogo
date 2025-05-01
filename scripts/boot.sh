#!/bin/sh

echo "🚀 Setting up environment..."

# Install Taskfile
echo "Installing Taskfile..."
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d v3.43.3
echo "✅ Taskfile installed successfully! $(task --version)"

# Install golangci-lint
echo "Installing golangci-lint..."
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.5
echo "✅ golangci-lint installed successfully! $(golangci-lint --version)"

# Temporary workaround for VS Code integration
# See: https://github.com/golang/vscode-go/issues/3732#issuecomment-2758960259
cp ~/go/bin/golangci-lint ~/go/bin/golangci-lint-v2

# Install swaggo
echo "Installing swaggo..."
go install github.com/swaggo/swag/cmd/swag@v1.16.4
echo "✅ swaggo installed successfully! $(swag --version)"

# Install atlas
echo "Installing atlas..."
curl -sSf https://atlasgo.sh | ATLAS_VERSION=v0.32.0 sh -s -- --community
echo "✅ atlas installed successfully! $(atlas version)"

# Starting infrastructure
echo "🚀 Starting infrastructure..."
docker compose up -d
echo "✅ Infrastructure started successfully!"

# Install pre-push hook
./scripts/install_githooks.sh

echo "🏁 Environment setup complete!"
