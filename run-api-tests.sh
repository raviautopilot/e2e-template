#!/bin/bash

# Ensure script executes in template directory where go.mod resides
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

echo "========================================="
echo " Starting E2E API Test Suite             "
echo "========================================="

# 1. Parse arguments, resolve shorthand package names, and set up target name for evidence directory
GO_TEST_ARGS=()
PKG_TAG=""
RUN_TAG=""
HAS_PKG=false
prev=""

for arg in "$@"; do
    if [ "$prev" = "-run" ]; then
        RUN_TAG="$arg"
        GO_TEST_ARGS+=("$arg")
        prev="$arg"
        continue
    elif [[ "$arg" == -run=* ]]; then
        RUN_TAG="${arg#-run=}"
        GO_TEST_ARGS+=("$arg")
        prev="$arg"
        continue
    fi

    # Check if arg is a shorthand package name (e.g. "twincore" -> tests/api/twincore)
    if [ -d "tests/api/$arg" ]; then
        HAS_PKG=true
        [ -z "$PKG_TAG" ] && PKG_TAG="$arg"
        GO_TEST_ARGS+=("./tests/api/$arg/...")
    # Check if arg is an explicit package path
    elif [[ "$arg" == tests/api/* ]] || [[ "$arg" == ./tests/api/* ]]; then
        HAS_PKG=true
        pkg="${arg#./tests/api/}"
        pkg="${pkg#tests/api/}"
        pkg="${pkg%/...}"
        pkg="${pkg%/*}"
        [ -z "$PKG_TAG" ] && PKG_TAG="$pkg"
        if [[ "$arg" != ./* ]]; then
            GO_TEST_ARGS+=("./$arg")
        else
            GO_TEST_ARGS+=("$arg")
        fi
    elif [[ "$arg" == tests/* ]] || [[ "$arg" == ./tests/* ]]; then
        HAS_PKG=true
        pkg="${arg#./tests/}"
        pkg="${pkg#tests/}"
        pkg="${pkg%/...}"
        pkg="${pkg%/*}"
        [ -z "$PKG_TAG" ] && PKG_TAG="$pkg"
        if [[ "$arg" != ./* ]]; then
            GO_TEST_ARGS+=("./$arg")
        else
            GO_TEST_ARGS+=("$arg")
        fi
    else
        GO_TEST_ARGS+=("$arg")
    fi
    prev="$arg"
done

# Determine TARGET_TAG
TARGET_TAG=""
if [ -n "$PKG_TAG" ] && [ -n "$RUN_TAG" ]; then
    TARGET_TAG="${PKG_TAG}-${RUN_TAG}"
elif [ -n "$PKG_TAG" ]; then
    TARGET_TAG="$PKG_TAG"
elif [ -n "$RUN_TAG" ]; then
    TARGET_TAG="$RUN_TAG"
else
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
if [ "$HAS_PKG" = true ]; then
    go test -v "${GO_TEST_ARGS[@]}"
else
    go test -v ./tests/api/... "${GO_TEST_ARGS[@]}"
fi
TEST_EXIT_CODE=$?

echo "========================================="
echo " API Tests finished with exit code $TEST_EXIT_CODE"
echo " Evidence saved in: evidence/run-$E2E_RUN_TIMESTAMP"
echo "========================================="

exit $TEST_EXIT_CODE

