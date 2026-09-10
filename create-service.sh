#!/bin/bash
# ==============================================================================
# Service Test Boilerplate Generator for Go E2E Testing Framework
# ==============================================================================
# Quickly scaffolds clean, production-ready API and UI test packages.
#
# Usage:
#   ./create-service.sh <service-name> [--api|--ui|--all] [--force]
#
# Examples:
#   ./create-service.sh payments
#   ./create-service.sh orders --all
#   ./create-service.sh auth --force
# ==============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors for terminal output
BOLD="\033[1m"
GREEN="\033[0;32m"
BLUE="\033[0;34m"
YELLOW="\033[0;33m"
RED="\033[0;31m"
NC="\033[0m" # No Color

print_help() {
    echo -e "${BOLD}Go E2E Framework — Service Test Boilerplate Generator${NC}"
    echo ""
    echo -e "${BOLD}Usage:${NC}"
    echo "  ./create-service.sh <service-name> [options]"
    echo ""
    echo -e "${BOLD}Options:${NC}"
    echo "  --api        Generate API test package (default)"
    echo "  --ui         Generate UI test package & Page Object"
    echo "  --all        Generate both API and UI test packages"
    echo "  --force      Overwrite existing files if directory already exists"
    echo "  -h, --help   Show this help message"
    echo ""
    echo -e "${BOLD}Examples:${NC}"
    echo "  ./create-service.sh orders"
    echo "  ./create-service.sh payments --api"
    echo "  ./create-service.sh auth --all"
    echo "  make new-service name=billing"
    echo ""
}

# Parse Arguments
SERVICE_NAME=""
GEN_MODE=""
FORCE=false
INTERACTIVE_MODE=false

while [[ $# -gt 0 ]]; do
    case "$1" in
        -h|--help)
            print_help
            exit 0
            ;;
        --api)
            GEN_MODE="api"
            shift
            ;;
        --ui)
            GEN_MODE="ui"
            shift
            ;;
        --all)
            GEN_MODE="all"
            shift
            ;;
        --force)
            FORCE=true
            shift
            ;;
        -*)
            echo -e "${RED}Error: Unknown option $1${NC}"
            print_help
            exit 1
            ;;
        *)
            if [ -z "$SERVICE_NAME" ]; then
                SERVICE_NAME="$1"
            fi
            shift
            ;;
    esac
done

# Interactive Prompt Mode if service name was not provided via arguments
if [ -z "$SERVICE_NAME" ]; then
    INTERACTIVE_MODE=true
    echo -e "${BLUE}=====================================================${NC}"
    echo -e "${BOLD} ✨ Interactive Service Test Boilerplate Generator${NC}"
    echo -e "${BLUE}=====================================================${NC}"
    echo ""

    # 1. Ask for service name
    while [ -z "$SERVICE_NAME" ]; do
        echo -ne "${BOLD}👉 Enter the service name (e.g. orders, payments, auth): ${NC}"
        read -r INPUT_NAME
        SERVICE_NAME=$(echo "$INPUT_NAME" | xargs)
        if [ -z "$SERVICE_NAME" ]; then
            echo -e "${RED}Service name cannot be empty. Please try again.${NC}"
        fi
    done

    # 2. Ask for test suite type
    echo ""
    echo -e "${BOLD}👉 Select test suite type to generate:${NC}"
    echo -e "   ${BOLD}[1]${NC} API tests only ${GREEN}(Default)${NC}"
    echo -e "   ${BOLD}[2]${NC} UI tests only (Page Object + Journey test)"
    echo -e "   ${BOLD}[3]${NC} Both (API + UI tests)"
    echo -ne "${BOLD}Enter choice [1-3] (Default: 1): ${NC}"
    read -r TYPE_CHOICE

    case "$TYPE_CHOICE" in
        2|"ui"|"UI")
            GEN_MODE="ui"
            ;;
        3|"both"|"all"|"ALL")
            GEN_MODE="all"
            ;;
        1|"api"|"API"|"")
            GEN_MODE="api"
            ;;
        *)
            echo -e "${YELLOW}Notice: Unknown selection '$TYPE_CHOICE'. Defaulting to [1] API tests only.${NC}"
            GEN_MODE="api"
            ;;
    esac
    echo ""
else
    # Default to api if not specified via flags in CLI mode
    if [ -z "$GEN_MODE" ]; then
        GEN_MODE="api"
    fi
fi

# Sanitize service name (lowercase, letters/digits/hyphens/underscores)
RAW_NAME="$SERVICE_NAME"
SERVICE_SLUG=$(echo "$RAW_NAME" | tr '[:upper:]' '[:lower:]' | tr ' ' '_' | sed 's/[^a-z0-9_-]//g')
PKG_NAME=$(echo "$SERVICE_SLUG" | tr '-' '_')

# TitleCase for Go types and test names
# e.g., "user-auth" -> "UserAuth", "orders" -> "Orders"
PASCAL_NAME=$(echo "$SERVICE_SLUG" | awk -F'[-_]' '{for(i=1;i<=NF;i++) printf toupper(substr($i,1,1)) substr($i,2)}')

echo -e "${BLUE}=====================================================${NC}"
echo -e "${BOLD} Generating E2E Test Suite for: ${GREEN}${SERVICE_SLUG}${NC}"
echo -e " Package Name: ${BOLD}${PKG_NAME}_test${NC} | Mode: ${BOLD}${GEN_MODE}${NC}"
echo -e "${BLUE}=====================================================${NC}"

# ─────────────────────────────────────────────────────────────────────────────
# 1. Generate API Test Suite
# ─────────────────────────────────────────────────────────────────────────────
if [ "$GEN_MODE" = "api" ] || [ "$GEN_MODE" = "all" ]; then
    API_DIR="tests/api/${SERVICE_SLUG}"

    if [ -d "$API_DIR" ] && [ "$FORCE" = false ]; then
        if [ "$INTERACTIVE_MODE" = true ]; then
            echo -ne "${YELLOW}Directory $API_DIR already exists. Overwrite? [y/N]: ${NC}"
            read -r CONFIRM_OVERWRITE
            if [[ "$CONFIRM_OVERWRITE" =~ ^[Yy]$ ]]; then
                FORCE=true
            else
                echo -e "${YELLOW}Skipping existing API directory $API_DIR.${NC}"
            fi
        else
            echo -e "${YELLOW}Warning: Directory $API_DIR already exists. Use --force to overwrite.${NC}"
        fi
    fi

    if [ ! -d "$API_DIR" ] || [ "$FORCE" = true ]; then
        mkdir -p "$API_DIR"
        echo -e "${GREEN}[+]${NC} Created API directory: ${BOLD}${API_DIR}${NC}"

        # 1.1 main_test.go
        cat <<EOF > "$API_DIR/main_test.go"
package ${PKG_NAME}_test

import (
	"os"
	"testing"

	"e2e-template/tests"
)

// TestMain bootstraps the test suite for the ${SERVICE_SLUG} service.
// It initializes global configs, sets up reporting, and handles test suite lifecycle.
func TestMain(m *testing.M) {
	tests.SetupSuite()
	exitCode := m.Run()
	tests.TeardownSuite()
	os.Exit(exitCode)
}
EOF
        echo -e "    └─ ${BOLD}main_test.go${NC} (suite lifecycle hook)"

        # 1.2 types_test.go
        cat <<EOF > "$API_DIR/types_test.go"
package ${PKG_NAME}_test

// ─────────────────────────────────────────────────────────────────────────────
// Data Models for ${PASCAL_NAME} Service
// ─────────────────────────────────────────────────────────────────────────────

// HealthResponse represents the health/status check response.
type HealthResponse struct {
	Status    string \`json:"status,omitempty"\`
	Message   string \`json:"message,omitempty"\`
	Version   string \`json:"version,omitempty"\`
	Timestamp string \`json:"timestamp,omitempty"\`
}

// ${PASCAL_NAME}Request represents the payload for creating or querying ${SERVICE_SLUG}.
type ${PASCAL_NAME}Request struct {
	Name        string                 \`json:"name,omitempty"\`
	Email       string                 \`json:"email,omitempty"\`
	Description string                 \`json:"description,omitempty"\`
	Metadata    map[string]interface{} \`json:"metadata,omitempty"\`
}

// ${PASCAL_NAME}Response represents the API response for ${SERVICE_SLUG}.
type ${PASCAL_NAME}Response struct {
	ID        string                 \`json:"id,omitempty"\`
	Status    string                 \`json:"status,omitempty"\`
	Message   string                 \`json:"message,omitempty"\`
	Data      map[string]interface{} \`json:"data,omitempty"\`
	CreatedAt string                 \`json:"created_at,omitempty"\`
}

// ErrorResponse represents a standard error response payload.
type ErrorResponse struct {
	Error   string \`json:"error,omitempty"\`
	Message string \`json:"message,omitempty"\`
	Code    int    \`json:"code,omitempty"\`
}
EOF
        echo -e "    └─ ${BOLD}types_test.go${NC} (JSON request/response models)"

        # 1.3 01_health_check_test.go
        cat <<EOF > "$API_DIR/01_health_check_test.go"
package ${PKG_NAME}_test

import (
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_${PASCAL_NAME}_01_HealthCheck verifies service availability and health.
func TestAPI_${PASCAL_NAME}_01_HealthCheck(t *testing.T) {
	baseURL := tests.GlobalConfig.BaseURL
	apiClient := client.NewClient(baseURL, 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(
		t,
		"${PASCAL_NAME} - Service Health Check",
		"Verifies that the ${SERVICE_SLUG} health check endpoint responds with 200 OK.",
		"HTTP 200 OK with healthy status payload",
		func(tc *tests.TestContext) {
			apiClient.SetTestName("${PASCAL_NAME} - Service Health Check")
			tc.Client = apiClient

			var resp HealthResponse
			// Adjust endpoint path to match your service (/health, /api/health, /ping, etc.)
			actions.GetAndExpectOK(tc, apiClient, "/health", &resp)
			actions.AssertNotEmpty(tc, "status", resp.Status)
		},
	)
}
EOF
        echo -e "    └─ ${BOLD}01_health_check_test.go${NC} (baseline health check)"

        # 1.4 02_parameterized_test.go
        cat <<EOF > "$API_DIR/02_parameterized_test.go"
package ${PKG_NAME}_test

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// ${PASCAL_NAME}TestCase defines a table-driven scenario for ${SERVICE_SLUG} API testing.
type ${PASCAL_NAME}TestCase struct {
	Name           string
	Description    string
	Payload        ${PASCAL_NAME}Request
	WantStatusCode int
	ExpectSuccess  bool
}

// TestAPI_${PASCAL_NAME}_02_Parameterized executes table-driven test scenarios.
func TestAPI_${PASCAL_NAME}_02_Parameterized(t *testing.T) {
	baseURL := tests.GlobalConfig.BaseURL
	apiClient := client.NewClient(baseURL, 15*time.Second, tests.ExecutionLogDir)

	scenarios := []${PASCAL_NAME}TestCase{
		{
			Name:        "Valid Resource Creation",
			Description: "Submit valid request payload and expect successful creation.",
			Payload: ${PASCAL_NAME}Request{
				Name:        "Sample ${PASCAL_NAME}",
				Email:       tests.GlobalConfig.AdminCredentials.Username,
				Description: "Automated E2E Test Entry",
			},
			WantStatusCode: 200,
			ExpectSuccess:  true,
		},
		{
			Name:        "Missing Required Field",
			Description: "Submit empty name field and expect validation rejection.",
			Payload: ${PASCAL_NAME}Request{
				Name:        "",
				Email:       tests.GlobalConfig.AdminCredentials.Username,
				Description: "Missing name payload",
			},
			WantStatusCode: 400,
			ExpectSuccess:  false,
		},
		{
			Name:        "Malformed Email Syntax",
			Description: "Submit invalid email format to verify input sanitizer.",
			Payload: ${PASCAL_NAME}Request{
				Name:        "Malformed Email Tester",
				Email:       "invalid-email-syntax",
				Description: "Syntax test",
			},
			WantStatusCode: 400,
			ExpectSuccess:  false,
		},
		{
			Name:        "Special Characters / SQL Injection Pattern",
			Description: "Submit SQL injection payload to test backend sanitization.",
			Payload: ${PASCAL_NAME}Request{
				Name:        "'; DROP TABLE users; --",
				Email:       tests.GlobalConfig.AdminCredentials.Username,
				Description: "Security injection probe",
			},
			WantStatusCode: 400,
			ExpectSuccess:  false,
		},
	}

	for _, sc := range scenarios {
		sc := sc
		t.Run(sc.Name, func(t *testing.T) {
			expectedText := fmt.Sprintf("HTTP %d Status", sc.WantStatusCode)
			if sc.ExpectSuccess {
				expectedText = "HTTP 200 OK with valid response payload"
			}

			testName := fmt.Sprintf("${PASCAL_NAME} - %s", sc.Name)

			tests.RunAPITestWithDetails(
				t,
				testName,
				sc.Description,
				expectedText,
				func(tc *tests.TestContext) {
					apiClient.SetTestName(testName)
					tc.Client = apiClient

					var resp ${PASCAL_NAME}Response
					// Adjust endpoint path to your target API route (e.g. /api/v1/${SERVICE_SLUG})
					targetEndpoint := "/api/v1/${SERVICE_SLUG}"

					if sc.ExpectSuccess {
						actions.PostAndExpectOK(tc, apiClient, targetEndpoint, &sc.Payload, &resp)
					} else {
						actions.PostAndExpectStatus(tc, apiClient, targetEndpoint, &sc.Payload, sc.WantStatusCode)
					}
				},
			)
		})
	}
}
EOF
        echo -e "    └─ ${BOLD}02_parameterized_test.go${NC} (table-driven test matrix)"
    fi
fi

# ─────────────────────────────────────────────────────────────────────────────
# 2. Generate UI Test Suite & Page Object
# ─────────────────────────────────────────────────────────────────────────────
if [ "$GEN_MODE" = "ui" ] || [ "$GEN_MODE" = "all" ]; then
    UI_DIR="tests/ui/${SERVICE_SLUG}"
    PAGE_FILE="pkg/ui/pages/${SERVICE_SLUG}_page.go"

    # 2.1 Page Object
    if [ -f "$PAGE_FILE" ] && [ "$FORCE" = false ]; then
        echo -e "${YELLOW}Warning: Page object $PAGE_FILE already exists. Skipped.${NC}"
    else
        mkdir -p "pkg/ui/pages"
        cat <<EOF > "$PAGE_FILE"
package pages

import (
	"time"

	"e2e-template/pkg/ui"
)

// ${PASCAL_NAME}Page represents the Page Object Model (POM) for the ${SERVICE_SLUG} UI view.
type ${PASCAL_NAME}Page struct {
	*ui.Page

	// Element Locators (CSS, XPath, or testid)
	HeadingText   string
	SubmitButton  string
	InputField    string
	ResultCard    string
}

// New${PASCAL_NAME}Page initializes the page object with selectors.
func New${PASCAL_NAME}Page(page *ui.Page) *${PASCAL_NAME}Page {
	return &${PASCAL_NAME}Page{
		Page:         page,
		HeadingText:  "css:h1",
		SubmitButton: "css:button[type='submit']",
		InputField:   "css:input[name='query']",
		ResultCard:   "css:.result-item",
	}
}

// SearchAndSubmit enters a search term and clicks submit.
func (p *${PASCAL_NAME}Page) SearchAndSubmit(query string, timeout time.Duration) error {
	if err := p.SendKeys(p.InputField, query, timeout); err != nil {
		return err
	}
	return p.Click(p.SubmitButton, timeout)
}

// WaitForHeading waits until the heading element is visible.
func (p *${PASCAL_NAME}Page) WaitForHeading(timeout time.Duration) error {
	_, err := p.WaitUntilVisible(p.HeadingText, timeout)
	return err
}
EOF
        echo -e "${GREEN}[+]${NC} Created Page Object: ${BOLD}${PAGE_FILE}${NC}"
    fi

    # 2.2 UI Test Directory
    if [ -d "$UI_DIR" ] && [ "$FORCE" = false ]; then
        echo -e "${YELLOW}Warning: Directory $UI_DIR already exists. Skipped.${NC}"
    else
        mkdir -p "$UI_DIR"
        echo -e "${GREEN}[+]${NC} Created UI test directory: ${BOLD}${UI_DIR}${NC}"

        # main_test.go
        cat <<EOF > "$UI_DIR/main_test.go"
package ${PKG_NAME}_test

import (
	"os"
	"testing"

	"e2e-template/tests"
)

func TestMain(m *testing.M) {
	tests.SetupSuite()
	exitCode := m.Run()
	tests.TeardownSuite()
	os.Exit(exitCode)
}
EOF
        echo -e "    └─ ${BOLD}main_test.go${NC}"

        # 01_journey_test.go
        cat <<EOF > "$UI_DIR/01_${SERVICE_SLUG}_journey_test.go"
package ${PKG_NAME}_test

import (
	"testing"
	"time"

	"e2e-template/pkg/ui"
	"e2e-template/pkg/ui/actions"
	"e2e-template/tests"
)

// TestUI_${PASCAL_NAME}_01_Journey verifies navigation and core interactions.
func TestUI_${PASCAL_NAME}_01_Journey(t *testing.T) {
	tests.RunUITest(t, "${PASCAL_NAME} User Journey", func(t *testing.T, page *ui.Page) {
		targetURL := tests.GlobalConfig.UiURL
		persona := actions.NewPublicPersona(page, targetURL, 10*time.Second)
		result := actions.NewResult("${PASCAL_NAME}Journey")

		actions.GoToHome(persona, result)
		actions.VerifyElementVisible(persona, result, "css:h1", "MainHeading")

		if result.Failed() {
			t.Errorf("${PASCAL_NAME} UI test failed: %v", result.Error)
		} else {
			t.Logf("✅ ${PASCAL_NAME} journey completed successfully.")
		}
	})
}
EOF
        echo -e "    └─ ${BOLD}01_${SERVICE_SLUG}_journey_test.go${NC}"
    fi
fi

# ─────────────────────────────────────────────────────────────────────────────
# 3. Verify Code Compilation
# ─────────────────────────────────────────────────────────────────────────────
echo ""
echo -e "${BLUE}Checking generated Go code syntax & compilation...${NC}"

COMPILE_SUCCESS=true
if [ "$GEN_MODE" = "api" ] || [ "$GEN_MODE" = "all" ]; then
    if ! go test -c "./tests/api/${SERVICE_SLUG}" -o /dev/null 2>/dev/null; then
        echo -e "${RED}❌ API compilation check failed. Please inspect generated files.${NC}"
        COMPILE_SUCCESS=false
    fi
fi

if [ "$GEN_MODE" = "ui" ] || [ "$GEN_MODE" = "all" ]; then
    if ! go test -c "./tests/ui/${SERVICE_SLUG}" -o /dev/null 2>/dev/null; then
        echo -e "${RED}❌ UI compilation check failed. Please inspect generated files.${NC}"
        COMPILE_SUCCESS=false
    fi
fi

if [ "$COMPILE_SUCCESS" = true ]; then
    echo -e "${GREEN}✅ All generated Go files compiled cleanly (0 syntax errors).${NC}"
fi

# ─────────────────────────────────────────────────────────────────────────────
# 4. Completion & Next Steps
# ─────────────────────────────────────────────────────────────────────────────
echo ""
echo -e "${GREEN}🎉 Boilerplate generated successfully for '${SERVICE_SLUG}'!${NC}"
echo ""
echo -e "${BOLD}Next steps:${NC}"
if [ "$GEN_MODE" = "api" ] || [ "$GEN_MODE" = "all" ]; then
    echo -e "  1. Customize endpoints & structs in: ${BOLD}tests/api/${SERVICE_SLUG}/types_test.go${NC}"
    echo -e "  2. Run your new API test suite:"
    echo -e "     ${BLUE}./run-api-tests.sh ${SERVICE_SLUG}${NC}"
    echo -e "  3. Or pass custom environment credentials:"
    echo -e "     ${BLUE}E2E_BASE_URL=https://api.yourdomain.com ./run-api-tests.sh ${SERVICE_SLUG}${NC}"
fi
if [ "$GEN_MODE" = "ui" ] || [ "$GEN_MODE" = "all" ]; then
    echo -e "  4. Adjust UI selectors in: ${BOLD}pkg/ui/pages/${SERVICE_SLUG}_page.go${NC}"
    echo -e "  5. Run your UI test suite:"
    echo -e "     ${BLUE}./run-ui-tests.sh ${SERVICE_SLUG}${NC}"
fi
echo ""
