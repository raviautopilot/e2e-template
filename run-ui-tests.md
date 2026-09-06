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

# Run only public site tests
./run-ui-tests.sh -run TestUI_Public

# Run a specific public test
./run-ui-tests.sh -run TestUI_Public_05_ExampleCom
```

### Option B: Using Standard `go test`
If you already have `chromedriver` running on port 9515:
```bash
go test -v ./tests/ui/... -run TestUI_Public -count=1
```

---

## 4. Test Catalog

### 🌐 Public Site UI Tests (`tests/ui/public_ui_test.go`)
Real, working browser automation tests targeting public internet websites:

| Test Name | Target Site | What It Tests |
|---|---|---|
| `TestUI_Public_01_GoogleJourney` | Google (`https://www.google.com`) | Page load, search input visibility, typing search query, submitting search |
| `TestUI_Public_02_GoogleSearchDirect` | Google (`https://www.google.com/search?q=...`) | Direct search URL navigation and results page title verification |
| `TestUI_Public_03_GitHubJourney` | GitHub (`https://github.com`) | Homepage load, main heading/header verification, title check |
| `TestUI_Public_04_GitHubPublicRepo` | GitHub (`https://github.com/octocat/Hello-World`) | Public repository page navigation, repo title & readme container verification |
| `TestUI_Public_05_ExampleCom` | example.com (`https://example.com`) | Lightweight baseline connectivity, H1 heading visibility and text check |

### 🛠️ Example Skeleton UI Tests (`tests/ui/example_ui_test.go`)
Template tests for your custom web application (configured via `uiUrl` in `config.json`):

| Test Name | Purpose |
|---|---|
| `TestUI_01_PublicJourneys` | Skeleton demonstrating public persona navigation across app pages |
| `TestUI_02_AdminLoginJourney` | Skeleton demonstrating admin persona login flow and dashboard access |

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

Every UI test execution records artifacts under `evidence/run-<timestamp>/`:
- **Screenshots**: Automatically captured on key transitions and errors at `evidence/run-<timestamp>/screenshots/<TestName>/`
- **Interactive HTML Report**: `evidence/run-<timestamp>/reports/report.html` (includes embedded screenshots and execution breakdown)
- **Markdown Report**: `evidence/run-<timestamp>/reports/report.md`
