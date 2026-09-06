#!/bin/bash

# Ensure script executes in template directory where go.mod resides
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

echo "========================================="
echo " Starting E2E API Test Suite             "
echo "========================================="

# 1. Parse target/test name from CLI arguments for meaningful evidence directory naming
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
        pkg="${arg#./tests/api/}"
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
export E2E_SUITE_NAME="API"
export E2E_TARGET_NAME="$TARGET_TAG"
export E2E_RUN_TIMESTAMP="api-${TARGET_TAG}-${TIMESTAMP}"

echo "Target Scope     : $TARGET_TAG"
echo "Evidence Dir     : evidence/run-$E2E_RUN_TIMESTAMP"

# 2. Run API test suite
echo "Running API test suite..."
HAS_PKG=false
for arg in "$@"; do
    if [[ "$arg" == ./tests/* ]] || [[ "$arg" == tests/* ]]; then
        HAS_PKG=true
        break
    fi
done

if [ "$HAS_PKG" = true ]; then
    go test -v "$@"
else
    go test -v ./tests/api/... "$@"
fi
TEST_EXIT_CODE=$?

echo "========================================="
echo " API Tests finished with exit code $TEST_EXIT_CODE"
echo " Evidence saved in: evidence/run-$E2E_RUN_TIMESTAMP"
echo "========================================="

exit $TEST_EXIT_CODE

