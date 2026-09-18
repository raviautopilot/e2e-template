#!/bin/bash

# Ensure script executes in template directory where go.mod resides
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# Configuration
PORT=9515

echo "========================================="
echo " Starting Chromedriver & E2E UI Tests    "
echo "========================================="

# 1. Determine Chromedriver binary location
CHROMEDRIVER_BIN="chromedriver"
if [ -n "$E2E_CHROMEDRIVER_PATH" ]; then
    CHROMEDRIVER_BIN="$E2E_CHROMEDRIVER_PATH"
elif [ -f "config.json" ]; then
    CONFIG_DRIVER=$(grep -o '"chromeDriverPath"[[:space:]]*:[[:space:]]*"[^"]*"' config.json | sed -E 's/"chromeDriverPath"[[:space:]]*:[[:space:]]*"([^"]*)"/\1/')
    if [ -z "$CONFIG_DRIVER" ]; then
        CONFIG_DRIVER=$(grep -o '"chromedriverPath"[[:space:]]*:[[:space:]]*"[^"]*"' config.json | sed -E 's/"chromedriverPath"[[:space:]]*:[[:space:]]*"([^"]*)"/\1/')
    fi
    if [ -n "$CONFIG_DRIVER" ]; then
        CHROMEDRIVER_BIN="$CONFIG_DRIVER"
    fi
fi

# Resolve relative path against script directory
if [[ "$CHROMEDRIVER_BIN" != /* ]] && [ "$CHROMEDRIVER_BIN" != "chromedriver" ]; then
    if [ -f "$SCRIPT_DIR/$CHROMEDRIVER_BIN" ]; then
        CHROMEDRIVER_BIN="$SCRIPT_DIR/$CHROMEDRIVER_BIN"
    fi
fi

# If a local driver file is specified but doesn't exist, warn and fallback to system chromedriver
if [ "$CHROMEDRIVER_BIN" != "chromedriver" ] && [ ! -f "$CHROMEDRIVER_BIN" ]; then
    echo "Warning: Configured ChromeDriver not found at $CHROMEDRIVER_BIN. Falling back to system 'chromedriver' in PATH."
    CHROMEDRIVER_BIN="chromedriver"
fi

if [ -f "$CHROMEDRIVER_BIN" ] && [ ! -x "$CHROMEDRIVER_BIN" ]; then
    chmod +x "$CHROMEDRIVER_BIN"
fi

echo "Using Chromedriver: $CHROMEDRIVER_BIN"
$CHROMEDRIVER_BIN --port=$PORT > /dev/null 2>&1 &
CHROMEDRIVER_PID=$!

# 2. Setup automatic cleanup on exit (trap)
cleanup() {
    echo "Cleaning up: stopping chromedriver (PID: $CHROMEDRIVER_PID)..."
    kill $CHROMEDRIVER_PID 2>/dev/null
    wait $CHROMEDRIVER_PID 2>/dev/null
    echo "Cleanup complete."
}
trap cleanup EXIT

# 3. Wait for Chromedriver to become responsive
echo "Waiting for Chromedriver to start on port $PORT..."
for i in {1..10}; do
    if curl -s http://localhost:$PORT/status | grep -q '"ready":true'; then
        echo "Chromedriver is ready!"
        break
    fi
    if [ $i -eq 10 ]; then
        echo "Error: Chromedriver failed to start on port $PORT."
        exit 1
    fi
    sleep 0.5
done

# 3.5. Parse target/test name from CLI arguments for meaningful evidence directory naming
TARGET_TAG=""
prev=""
for arg in "$@"; do
    if [ "$prev" = "-run" ]; then
        TARGET_TAG="$arg"
        break
    elif [[ "$arg" == -run=* ]]; then
        TARGET_TAG="${arg#-run=}"
        break
    elif [[ "$arg" == ./tests/* ]] || [[ "$arg" == tests/* ]]; then
        pkg="${arg#./tests/ui/}"
        pkg="${pkg#tests/ui/}"
        pkg="${pkg#./tests/api/}"
        pkg="${pkg#tests/api/}"
        pkg="${pkg%/...}"
        pkg="${pkg%/*}"
        TARGET_TAG="$pkg"
    fi
    prev="$arg"
done

if [ -z "$TARGET_TAG" ]; then
    TARGET_TAG="all"
fi

# Sanitize tag for directory safety
TARGET_TAG=$(echo "$TARGET_TAG" | sed 's/[^a-zA-Z0-9_-]/_/g' | sed 's/^_//;s/_$//')
[ -z "$TARGET_TAG" ] && TARGET_TAG="all"

TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
export E2E_SUITE_NAME="UI"
export E2E_TARGET_NAME="$TARGET_TAG"
export E2E_RUN_TIMESTAMP="ui-${TARGET_TAG}-${TIMESTAMP}"

echo "Target Scope     : $TARGET_TAG"
echo "Evidence Dir     : evidence/run-$E2E_RUN_TIMESTAMP"

# 4. Run UI test suite
echo "Running UI test suite..."
HAS_PKG=false
for arg in "$@"; do
    if [[ "$arg" == ./tests/* ]] || [[ "$arg" == tests/* ]]; then
        HAS_PKG=true
        break
    fi
done

if [ "$HAS_PKG" = true ]; then
    go test -v -timeout 30m "$@"
else
    go test -v -timeout 30m ./tests/ui/... "$@"
fi
TEST_EXIT_CODE=$?

echo "========================================="
echo " UI Tests finished with exit code $TEST_EXIT_CODE"
echo " Evidence saved in: evidence/run-$E2E_RUN_TIMESTAMP"
echo "========================================="

exit $TEST_EXIT_CODE

