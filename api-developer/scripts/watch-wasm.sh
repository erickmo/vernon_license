#!/bin/bash
# Watch and rebuild WASM on changes

set -e

echo "🔍 Watching for WASM changes... (Press Ctrl+C to stop)"
echo "   Monitoring: internal/ui/**/*.go"

last_build=$(date +%s%N)

while true; do
    # Get latest modification time in internal/ui
    latest=$(find internal/ui -name "*.go" -type f -printf '%TY-%Tm-%Td %TH:%TM:%TS\n' | sort -r | head -1)

    if [ -z "$latest" ]; then
        sleep 1
        continue
    fi

    # Compare with last build time
    latest_epoch=$(date -d "$latest" +%s%N 2>/dev/null || echo 0)

    if [ $latest_epoch -gt $last_build ]; then
        echo "📦 Change detected - rebuilding WASM..."
        GOARCH=wasm GOOS=js go build -o web/app.wasm ./cmd/api
        if [ $? -eq 0 ]; then
            echo "✅ WASM rebuilt successfully"
            last_build=$(date +%s%N)
        else
            echo "❌ WASM build failed"
        fi
    fi

    sleep 1
done
