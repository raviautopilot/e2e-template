#!/bin/bash
# ==============================================================================
# Script Name: validate-template.sh
# Description: Validates that the e2e-template repository builds cleanly, passes
#              vet checks, and contains zero domain-specific leaks (nammataga,
#              taga, office bearer, grievance, tbf, etc.).
# ==============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

ERRORS=0

log_pass() {
    echo -e "${GREEN}✅ PASS:${NC} $1"
}

log_fail() {
    echo -e "${RED}❌ FAIL:${NC} $1"
    ERRORS=$((ERRORS + 1))
}

log_info() {
    echo -e "${YELLOW}ℹ️  INFO:${NC} $1"
}

echo "=================================================="
echo "🔍 E2E Template Integrity & Non-Duplication Audit"
echo "=================================================="

# 1. Compilation check
log_info "Verifying Go compilation (go build ./...)..."
if go build ./... > /dev/null 2>&1; then
    log_pass "go build ./... compiled successfully"
else
    log_fail "go build ./... failed"
fi

# 2. Go test compilation check
log_info "Verifying Go test compilation (go test -run=^$ ./...)..."
if go test -run=^$ ./... > /dev/null 2>&1; then
    log_pass "all test files compile cleanly"
else
    log_fail "test files failed to compile"
fi

# 3. Go vet check
log_info "Running go vet ./... inspection..."
if go vet ./... > /dev/null 2>&1; then
    log_pass "go vet ./... passed with zero warnings"
else
    log_fail "go vet ./... found issues"
fi

# 4. Domain Leak Audit (Grep check across tracked repository files)
log_info "Scanning for domain-specific leakage..."

FORBIDDEN_PATTERNS=(
    "nammataga"
    "office[ -]?bearer"
    "grievance"
    "\btbf\b"
    "\btaga\b"
    "taga-"
    "taga_"
    "tagaId"
)

# Files/directories to exclude from text scan
EXCLUDES=(
    ":!evidence"
    ":!graphify-out"
    ":!.agents"
    ":!.git"
    ":!scripts/validate-template.sh"
)

for pattern in "${FORBIDDEN_PATTERNS[@]}"; do
    MATCHES=$(git grep -inE "$pattern" -- . "${EXCLUDES[@]}" 2>/dev/null || true)
    if [ -n "$MATCHES" ]; then
        log_fail "Found forbidden domain pattern '$pattern':"
        echo "$MATCHES" | head -n 10
        if [ $(echo "$MATCHES" | wc -l) -gt 10 ]; then
            echo "... and more"
        fi
    else
        log_pass "Zero matches for '$pattern'"
    fi
done

# 5. Check for deleted mock directories
log_info "Verifying obsolete fixture directories are absent..."
if [ -d "fixtures/about" ] || [ -d "fixtures/health" ]; then
    log_fail "Obsolete fixture directories (fixtures/about or fixtures/health) still exist"
else
    log_pass "Obsolete fixture directories are absent"
fi

# 6. Check for obsolete deployment scripts or archives
if [ -d "archives" ]; then
    log_fail "Obsolete archives directory still exists"
else
    log_pass "Obsolete archives directory is absent"
fi

echo "=================================================="
if [ "$ERRORS" -eq 0 ]; then
    echo -e "${GREEN}🎉 TEMPLATE VALIDATION PASSED — Zero domain leakage detected!${NC}"
    echo "=================================================="
    exit 0
else
    echo -e "${RED}💥 TEMPLATE VALIDATION FAILED — $ERRORS check(s) failed!${NC}"
    echo "=================================================="
    exit 1
fi
