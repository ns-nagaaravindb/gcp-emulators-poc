#!/bin/bash

echo "🚀 Starting GCS and Pub/Sub Emulator Demo"
echo "========================================"

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Stop any existing containers to avoid port conflicts
echo "🧹 Cleaning up any existing containers..."
docker-compose down >/dev/null 2>&1

# Check for port conflicts
echo "🔍 Checking for port conflicts..."
if lsof -i :4443 >/dev/null 2>&1; then
    echo "⚠️  Warning: Port 4443 is in use. This may cause conflicts."
    echo "   You may need to stop the process using this port."
fi

if lsof -i :8085 >/dev/null 2>&1; then
    echo "⚠️  Warning: Port 8085 is in use. This may cause conflicts."
    echo "   You may need to stop the process using this port."
fi

# Start the emulators
echo "📦 Starting emulators with Docker Compose..."
docker-compose up -d

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 15

# Check if services are running
echo "🔍 Checking service status..."
if curl -s http://localhost:4443 >/dev/null; then
    echo "✅ Fake GCS Server is running on http://localhost:4443"
else
    echo "❌ Fake GCS Server is not responding"
    echo "   Try: docker-compose logs fake-gcs-server"
fi

if curl -s http://localhost:8085 >/dev/null; then
    echo "✅ Pub/Sub Emulator is running on http://localhost:8085"
else
    echo "❌ Pub/Sub Emulator is not responding"
    echo "   Try: docker-compose logs pubsub-emulator"
fi

# Check Go version consistency
echo "🔧 Checking Go environment..."
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
GO_ROOT_VERSION=$(basename $(go env GOROOT))
if [ "$GO_VERSION" != "$GO_ROOT_VERSION" ]; then
    echo "⚠️  Warning: Go version mismatch detected!"
    echo "   go version: $GO_VERSION"
    echo "   GOROOT version: $GO_ROOT_VERSION"
    echo "   This may cause compilation issues."
    echo "   Consider running: go clean -cache && go clean -modcache"
fi

# Download Go dependencies
echo "📥 Downloading Go dependencies..."
go mod tidy

# Set environment variables for emulators
export STORAGE_EMULATOR_HOST=localhost:4443
export PUBSUB_EMULATOR_HOST=localhost:8085

echo ""
echo "🎯 Running GCS Demo..."
echo "====================="
go run gcs.go

echo ""
echo "🎯 Running Pub/Sub Demo..."
echo "=========================="
timeout 15s go run pubsubemulator.go

echo ""
echo "✨ Demo completed!"
echo ""
echo "To stop the emulators, run: docker-compose down"
echo "To view logs, run: docker-compose logs -f"