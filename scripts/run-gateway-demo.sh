#!/bin/bash
# Run the complete gateway demo with backends

set -e

echo "🚀 Starting Rate-Limited API Gateway Demo"
echo "=========================================="
echo ""

# Check if running from project root
if [ ! -f "go.mod" ]; then
    echo "❌ Error: Please run this script from the project root directory"
    echo "   cd /path/to/goroutine-3000 && ./scripts/run-gateway-demo.sh"
    exit 1
fi

echo "📦 Step 1: Starting backend services..."
echo ""

# Start backend 1
PORT=8081 go run examples/backend/main.go &
BACKEND1_PID=$!
echo "✓ Backend 1 started on port 8081 (PID: $BACKEND1_PID)"

# Start backend 2
PORT=8082 go run examples/backend/main.go &
BACKEND2_PID=$!
echo "✓ Backend 2 started on port 8082 (PID: $BACKEND2_PID)"

# Wait for backends to start
echo ""
echo "⏳ Waiting for backends to be ready..."
sleep 3

echo ""
echo "🚪 Step 2: Starting gateway on port 8080..."
echo ""

# Start gateway
go run examples/gateway/main.go &
GATEWAY_PID=$!

# Wait for gateway to start
sleep 2

echo ""
echo "✅ All services running!"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎯 Try these commands in another terminal:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "  # Make a simple request"
echo "  curl http://localhost:8080/api/hello"
echo ""
echo "  # Get data from backend"
echo "  curl http://localhost:8080/api/data"
echo ""
echo "  # Check gateway statistics"
echo "  curl http://localhost:8080/stats"
echo ""
echo "  # Test rate limiting (429 after 10 requests)"
echo "  for i in {1..15}; do curl http://localhost:8080/api/hello; echo; done"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "⚠️  Press Ctrl+C to stop all services"
echo ""

# Cleanup function
cleanup() {
    echo ""
    echo "🛑 Stopping all services..."
    kill $GATEWAY_PID 2>/dev/null || true
    kill $BACKEND1_PID 2>/dev/null || true
    kill $BACKEND2_PID 2>/dev/null || true
    echo "✓ All services stopped"
    exit 0
}

# Trap signals
trap cleanup INT TERM

# Keep script running
wait
