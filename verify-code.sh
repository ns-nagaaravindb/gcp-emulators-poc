#!/bin/bash

echo "🔍 Verifying Go code compilation..."

# Check gcs.go
echo "Checking gcs.go..."
if go build -o /tmp/gcs gcs.go 2>/dev/null; then
    echo "✅ gcs.go compiles successfully"
    rm -f /tmp/gcs
else
    echo "❌ gcs.go has compilation errors"
    go build gcs.go
fi

# Check pubsubemulator.go  
echo "Checking pubsubemulator.go..."
if go build -o /tmp/pubsub pubsubemulator.go 2>/dev/null; then
    echo "✅ pubsubemulator.go compiles successfully"
    rm -f /tmp/pubsub
else
    echo "❌ pubsubemulator.go has compilation errors"
    go build pubsubemulator.go
fi

echo ""
echo "✨ Code verification complete!"
echo ""
echo "Next steps:"
echo "1. Start Docker Desktop"  
echo "2. Run: ./run-demo.sh"
echo "   OR"
echo "3. Run: make demo"