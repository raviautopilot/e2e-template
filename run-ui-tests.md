# E2E UI Test Suite Execution Guide

This document outlines how to compile, execute, configure, and expand the UI automation tests using Selenium WebDriver and the Persona/Action Pattern.

---

## 1. Prerequisites & Environment Setup

Ensure `chromedriver` is installed and available on your system.

### Install Chromedriver (Ubuntu/Debian):
```bash
sudo apt-get update
sudo apt-get install -y chromium-browser chromium-chromedriver
```

### Automatic Chromedriver Management:
The `./run-ui-tests.sh` script automatically starts `chromedriver` on port `9515` before tests and cleanly shuts it down when tests finish.

---

## 2. Test Configuration

The test suite reads configuration from `config.json`. You can modify this file directly or override values via environment variables:

| `config.json` Field | Environment Variable | Default Value | Description |
|---|---|---|---|
| `uiUrl` | `E2E_UI_URL` | `http://localhost:3000` | Target URL of your web frontend |
| `seleniumUrl` | `E2E_SELENIUM_URL` | `http://localhost:9515` | Selenium WebDriver address |
| `headless` | `E2E_HEADLESS` | `false` | Run browser in headless mode (`true` for CI, `false` to watch browser) |
| `timeout` | `E2E_TIMEOUT` | `10` | Timeout in seconds for page element waits |

---

## 3. How to Run the Tests

### Option A: Using the `run-ui-tests.sh` Script (Recommended)
Automatically manages Chromedriver lifecycle and handles traps:

```bash
# Run all UI tests (headless mode by default in CI)
./run-ui-tests.sh

# Run with visible Chrome browser window (headed mode)
E2E_HEADLESS=false ./run-ui-tests.sh

# Run all UI tests in a specific package
./run-ui-tests.sh ./tests/ui/google/...
./run-ui-tests.sh ./tests/ui/github/...
./run-ui-tests.sh ./tests/ui/example/...

# Run a specific test
./run-ui-tests.sh -run TestUI_Google_01_SearchJourney
./run-ui-tests.sh -run TestUI_GitHub_01_NavigationJourney
./run-ui-tests.sh -run TestUI_Example_01_ExampleCom
```

### Option B: Using Standard `go test`
If you already have `chromedriver` running on port 9515:
```bash
go test -v ./tests/ui/...
go test -v ./tests/ui/google/...
go test -v ./tests/ui/github/...
go test -v ./tests/ui/example/...
```

---

## 4. Test Catalog

Each test suite resides in its own package under `tests/ui/`, with **one test per file** for readability:

### 🌐 Google Package (`tests/ui/google/`)
| File | Test Function | What It Tests |
|---|---|---|
| `01_search_journey_test.go` | `TestUI_Google_01_SearchJourney` | Page load, search input visibility, typing search query, submitting search |
| `02_search_direct_test.go` | `TestUI_Google_02_SearchDirect` | Direct search URL navigation and results page title verification |

### 🌐 GitHub Package (`tests/ui/github/`)
| File | Test Function | What It Tests |
|---|---|---|
| `01_navigation_journey_test.go` | `TestUI_GitHub_01_NavigationJourney` | Homepage load, main heading/header verification, title check |
| `02_public_repo_test.go` | `TestUI_GitHub_02_PublicRepo` | Public repository page navigation, repo title & readme container verification |

### 🛠️ Example Package (`tests/ui/example/`)
| File | Test Function | What It Tests |
|---|---|---|
| `01_example_com_test.go` | `TestUI_Example_01_ExampleCom` | Lightweight baseline connectivity and H1 heading check against example.com |
| `02_public_journey_test.go` | `TestUI_Example_02_PublicJourneys` | Skeleton demonstrating public persona navigation across your app |
| `03_admin_login_journey_test.go` | `TestUI_Example_03_AdminLoginJourney` | Skeleton demonstrating admin login journey across your app |

---

## 5. Writing New UI Tests

Tests follow the **Persona/Action Pattern**, which separates high-level user intentions from raw DOM/Selenium commands:

```go
package ui_test

import (
    "testing"
    "time"

    "e2e-template/pkg/ui"
    "e2e-template/pkg/ui/actions"
    "e2e-template/tests"
)

func TestMyPageJourney(t *testing.T) {
    tests.RunUITest(t, "My Page Journey", func(t *testing.T, page *ui.Page) {
        persona := actions.NewPublicPersona(page, "https://example.com", 10*time.Second)
        result := actions.NewResult("MyJourney")

        // Declarative human-readable actions
        actions.GoToHome(persona, result)
        actions.VerifyPageTitle(persona, result, "Example Domain")
        actions.VerifyElementVisible(persona, result, "h1", "H1Heading")

        // Fail test if any action in the journey encountered an error
        if result.Failed() {
            t.Fatalf("Journey failed: %v\nActions: %v", result.Error, result.Actions)
        }
    })
}
```

---

## 6. Test Evidence & Screenshots

Every UI test execution records artifacts in meaningful, test-scoped directories under `evidence/`:
- **Evidence Run Directory**: `evidence/run-ui-<target>-<timestamp>/`
  - When running a specific test via `-run <TestName>`, `<target>` reflects the test name (e.g., `evidence/run-ui-TestUI_Google_01_SearchJourney-2026-09-06_18-45-00/`).
  - When targeting a package (e.g. `./tests/ui/google/...`), `<target>` reflects the package (e.g., `evidence/run-ui-google-2026-09-06_18-45-00/`).
  - When running all tests, `<target>` defaults to `all` (e.g., `evidence/run-ui-all-2026-09-06_18-45-00/`).

### Evidence Artifacts:
- **Screenshots**: Captured on transitions and failures at `evidence/run-ui-<target>-<timestamp>/screenshots/<TestName>/`
- **Network Requests**: Intercepted browser request traces at `evidence/run-ui-<target>-<timestamp>/requests/<TestName>/`
- **Interactive HTML Report**: `evidence/run-ui-<target>-<timestamp>/reports/report.html` (includes lightbox viewer and step screenshots linked to each test)
- **Markdown Report**: `evidence/run-ui-<target>-<timestamp>/test-report.md` (summary table with evidence links)

