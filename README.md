# Go E2E Testing Framework

A ready-to-use, modular End-to-End (E2E) testing framework built in Go. It supports automated API testing via a type-safe HTTP client wrapper, UI automation using Selenium WebDriver following the Page Object Model (POM) pattern, custom HTML/JSON reporting, and granular request/response logging.

---

## Features

- **Extensible API Client**: Auto-marshaling, struct pointer safety checks, and an interface-driven Authentication manager (Bearer, Basic, API Key, mTLS, and SSH signing).
- **Selenium UI Integration**: Base Page Object wrappers handling dynamic CSS/XPath element selection, waiting hooks, interaction wrappers, and automated screenshots on test failure.
- **Observability Logging**: Automatic date-organized file logging of request and response payloads at `evidence/run-<timestamp>/requests/`.
- **Beautiful HTML/JSON Reporting**: Interactive responsive test results dashboard compiled automatically after each execution.

---

## Directory Structure

```
├── go.mod                     # Go module definitions
├── go.sum                     # Dependency checksums
├── Makefile                   # Execution shortcuts
├── config.json                # Environment configuration (adapt to your project)
├── pkg/
│   ├── config/
│   │   └── config.go          # Config loader and environment variable overrides
│   ├── logger/
│   │   └── logger.go          # Thread-safe logging levels (INFO, DEBUG, WARN, ERROR)
│   ├── client/
│   │   ├── auth.go            # Authentication interface and sub-types
│   │   └── client.go          # Custom HTTP client wrapper and JSON logger
│   ├── ui/
│   │   ├── driver.go          # Selenium WebDriver connection & option manager
│   │   ├── pom.go             # Page Object Model helper wrappers
│   │   ├── pages/
│   │   │   ├── home_page.go          # Example home page object
│   │   │   └── login_page.go         # Example login page object
│   │   └── actions/
│   │       ├── admin_actions.go      # Example admin user actions
│   │       ├── member_actions.go     # Example member user actions
│   │       └── public_actions.go     # Example public user actions
│   └── report/
│       ├── report.go          # Result collector and HTML compiler
│       ├── template.html      # Visual dashboard layout template
│       └── template.md        # Markdown report template
├── tests/
│   ├── main_test.go           # Suite bootstrap (TestMain)
│   ├── helpers.go             # RunAPITest, RunUITest, seedTestData, etc.
│   ├── api/
│   │   ├── main_test.go       # API package bootstrap
│   │   ├── types_test.go      # Response type definitions
│   │   └── example_api_test.go  # Example API test cases
│   └── ui/
│       ├── main_test.go       # UI package bootstrap
│       └── example_ui_test.go   # Example UI test cases
├── fixtures/                  # Test data files (CSV, JSON, PDF, images)
│   ├── example_bulk_upload.csv
│   ├── example_resource.pdf
│   └── example_image.png
├── docs/
│   ├── MASTER_API_TESTING_PROMPT.md  # AI prompt template for generating API tests
│   └── references/                   # Reference documentation
├── archives/                  # Historical reference documents
└── scripts/                   # Helper scripts (publishing, npm, graphify, etc.)
```

---

## Prerequisites

1. **Golang**: Ensure Go 1.21+ is installed.
2. **Google Chrome & Chromedriver**: Install Chrome and Chromedriver on your local machine or server.
   - For Ubuntu/Linux:
     ```bash
     sudo apt-get update
     sudo apt-get install -y chromium-browser chromium-chromedriver
     ```
   - Ensure `chromedriver` is available in your PATH or running on port `9515`.

---

## Quick Start

1. **Clone this template**:
   ```bash
   ./clone.sh /path/to/your-new-test-project
   cd /path/to/your-new-test-project
   ```

2. **Configure your targets** in `config.json`:
   ```json
   {
     "baseUrl": "https://api.yourapp.com",
     "uiUrl": "https://yourapp.com",
     "seleniumUrl": "http://localhost:9515",
     "headless": false,
     "timeout": 10
   }
   ```

3. **Install dependencies**:
   ```bash
   make deps
   ```

4. **Start Chromedriver**:
   ```bash
   chromedriver --port=9515
   ```

5. **Run example tests**:
   ```bash
   make test-all
   ```

6. **Run in Headless Mode (CI/VMs)**:
   ```bash
   E2E_HEADLESS=true make test-all
   ```

---

## Writing Tests

### 1. API Test Cases

Create tests using `RunAPITestWithDetails` for rich reporting, or `RunAPITest` for simpler cases:

```go
func TestAPI_GetUser(t *testing.T) {
    tests.RunAPITestWithDetails(
        t,
        "GET /users/1 returns a valid user",
        "Verifies that fetching user ID 1 returns a non-empty name and email.",
        "HTTP 200 OK with non-empty name and email",
        func(tc *tests.TestContext) {
            var user UserResponse
            err := tc.Client.SendHttpRequest("GET", "/users/1", nil, nil, &user, nil)
            if err != nil {
                tc.FailureReason = fmt.Sprintf("Request failed: %v", err)
                tc.Fatalf("Request failed: %v", err)
            }
            tc.Actual = fmt.Sprintf("HTTP 200 OK, name=%q, email=%q", user.Name, user.Email)
            if user.Name == "" {
                tc.Errorf("Expected non-empty name")
            }
        },
    )
}
```

### 2. UI Test Cases

Create page objects in `pkg/ui/pages/` and run them with `RunUITest`:

```go
func TestUI_LoginFlow(t *testing.T) {
    tests.RunUITest(t, "Admin Login Flow", func(t *testing.T, page *ui.Page) {
        cfg := tests.GlobalConfig

        if err := page.Navigate(cfg.UiURL + "/login"); err != nil {
            t.Fatalf("Failed to navigate: %v", err)
        }

        page.Click("testid:" + cfg.AdminLoginButtonTestID)
        page.TypeText("testid:" + cfg.AdminLoginUsernameInputTestID, cfg.AdminCredentials.Username)
        page.TypeText("testid:" + cfg.AdminLoginPasswordInputTestID, cfg.AdminCredentials.Password)
        page.Click("testid:" + cfg.AdminLoginSubmitButtonTestID)

        // Verify dashboard loaded
        if err := page.WaitForElement("testid:dashboard-header"); err != nil {
            t.Errorf("Dashboard did not load after login: %v", err)
        }
    })
}
```

*Note: If a UI test fails, the framework automatically writes a screenshot to `evidence/run-<timestamp>/screenshots/` and includes an inline link/preview in the HTML report.*

---

## Configuration

### config.json Fields

| Field | Default | Description |
|---|---|---|
| `baseUrl` | `http://localhost:8080` | Your API's base URL |
| `uiUrl` | `http://localhost:3000` | Your web app's URL |
| `seleniumUrl` | `http://localhost:9515` | ChromeDriver WebSocket URL |
| `headless` | `false` | Set `true` for headless Chrome (CI) |
| `timeout` | `10` | Default timeout in seconds |
| `adminCredentials` | — | Admin username/password for tests |
| `memberCredentials` | — | Member/user credentials for tests |
| `*TestID` fields | — | `data-testid` attribute values for UI elements |

### Environment Variable Overrides

| Config JSON | Env Variable |
|---|---|
| `baseUrl` | `E2E_BASE_URL` |
| `uiUrl` | `E2E_UI_URL` |
| `seleniumUrl` | `E2E_SELENIUM_URL` |
| `headless` | `E2E_HEADLESS` |
| `timeout` | `E2E_TIMEOUT` |

---

## Reports

After each test run, reports are generated in `evidence/run-<timestamp>/`:

```
evidence/run-2026-01-15_14-30-00/
├── reports/
│   ├── report.html    # Interactive HTML dashboard (open in browser)
│   └── report.md      # Markdown summary
├── requests/          # Per-test request/response JSON logs
└── screenshots/       # Failure screenshots (UI tests)
```
