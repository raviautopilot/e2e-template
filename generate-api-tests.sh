#!/bin/bash

# generate-api-tests.sh — Interactive API test generator from Swagger
# Usage: ./generate-api-tests.sh [flags]
# Interactive mode: ./generate-api-tests.sh  (no flags, prompts for input)
# Non-interactive: ./generate-api-tests.sh --swagger-url URL --service-name NAME --base-url URL

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

echo "========================================="
echo "  API Test Generator from Swagger        "
echo "========================================="
echo ""

# If CLI flags are passed, forward them directly to the Go tool
if [ $# -gt 0 ]; then
    echo "Running in non-interactive mode with provided flags..."
    echo ""
    go run ./cmd/generate-tests "$@"
    exit $?
fi

# ── Interactive Mode ─────────────────────────────────────────────────────────

# 1. Swagger Source
echo "1. Swagger Source"
echo "   Enter swagger URL or local file path:"
read -rp "   > " SWAGGER_SRC

if [ -z "$SWAGGER_SRC" ]; then
    echo "ERROR: Swagger source is required."
    exit 1
fi

# Determine if it's a URL or file
SWAGGER_FLAG=""
if [[ "$SWAGGER_SRC" == http://* ]] || [[ "$SWAGGER_SRC" == https://* ]]; then
    SWAGGER_FLAG="--swagger-url"
    # Auto-detect Swagger UI URL and convert to doc.json if applicable
    if [[ "$SWAGGER_SRC" =~ /index\.html?$ ]]; then
        SWAGGER_JSON_GUESS=$(echo "$SWAGGER_SRC" | sed -E 's|/index\.html?$|/doc.json|')
        echo "   ℹ Detected Swagger UI HTML URL. Auto-converting to: $SWAGGER_JSON_GUESS"
        SWAGGER_SRC="$SWAGGER_JSON_GUESS"
    elif [[ "$SWAGGER_SRC" =~ /swagger/?$ ]]; then
        SWAGGER_JSON_GUESS="${SWAGGER_SRC%/}/doc.json"
        echo "   ℹ Detected Swagger directory. Auto-converting to: $SWAGGER_JSON_GUESS"
        SWAGGER_SRC="$SWAGGER_JSON_GUESS"
    fi
else
    SWAGGER_FLAG="--swagger-file"
    if [ ! -f "$SWAGGER_SRC" ]; then
        echo "ERROR: File not found: $SWAGGER_SRC"
        exit 1
    fi
fi

echo ""

# 2. Is it live?
echo "2. Is this a live endpoint? (y/n):"
read -rp "   > " IS_LIVE
echo ""

# 3. Service name
echo "3. Service name (used for package directory, e.g. configsvc, usersvc):"
read -rp "   > " SERVICE_NAME

if [ -z "$SERVICE_NAME" ]; then
    echo "ERROR: Service name is required."
    exit 1
fi

# Sanitize service name
SERVICE_NAME=$(echo "$SERVICE_NAME" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9_]//g')

echo ""

# 4. Generate scope
echo "4. Generate tests for:"
echo "   a) All endpoints"
echo "   b) Specific tag group(s)"
read -rp "   > " SCOPE

TAGS_FLAG=""
if [ "$SCOPE" = "b" ] || [ "$SCOPE" = "B" ]; then
    echo ""
    echo "   Enter comma-separated tag names (e.g. Dependencies,Modules):"
    read -rp "   > " TAGS
    if [ -n "$TAGS" ]; then
        TAGS_FLAG="--tags $TAGS"
    fi
fi

echo ""

# 5. Base URL
DEFAULT_BASE_URL=""
if [[ "$SWAGGER_SRC" == http://* ]] || [[ "$SWAGGER_SRC" == https://* ]]; then
    # Extract base URL from swagger URL (protocol + host:port)
    DEFAULT_BASE_URL=$(echo "$SWAGGER_SRC" | sed -E 's|(https?://[^/]+).*|\1|')
fi

echo "5. Base URL for tests${DEFAULT_BASE_URL:+ (default: $DEFAULT_BASE_URL)}:"
read -rp "   > " BASE_URL

if [ -z "$BASE_URL" ] && [ -n "$DEFAULT_BASE_URL" ]; then
    BASE_URL="$DEFAULT_BASE_URL"
fi

if [ -z "$BASE_URL" ]; then
    echo "ERROR: Base URL is required."
    exit 1
fi

echo ""

# 6. Output directory
DEFAULT_OUTPUT="tests/api/$SERVICE_NAME"
echo "6. Output directory (default: $DEFAULT_OUTPUT):"
read -rp "   > " OUTPUT_DIR

if [ -z "$OUTPUT_DIR" ]; then
    OUTPUT_DIR="$DEFAULT_OUTPUT"
fi

echo ""
echo "==========================================="
echo "  Summary"
echo "==========================================="
echo "  Swagger:     $SWAGGER_SRC"
echo "  Live:        $IS_LIVE"
echo "  Service:     $SERVICE_NAME"
echo "  Base URL:    $BASE_URL"
echo "  Output:      $OUTPUT_DIR"
echo "  Tags:        ${TAGS:-all}"
echo "==========================================="
echo ""
echo "Generating..."
echo ""

# Build and run the generator
# shellcheck disable=SC2086
go run ./cmd/generate-tests \
    $SWAGGER_FLAG "$SWAGGER_SRC" \
    --service-name "$SERVICE_NAME" \
    --base-url "$BASE_URL" \
    --output-dir "$OUTPUT_DIR" \
    $TAGS_FLAG

echo ""
echo "==========================================="
echo "  Next Steps"
echo "==========================================="
echo "  1. Review generated tests in: $OUTPUT_DIR"
echo "  2. Run tests: ./run-api-tests.sh ./$OUTPUT_DIR/..."
echo "  3. View reports in: evidence/run-*/reports/"
echo "==========================================="
