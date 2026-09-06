# Go E2E Testing Framework

A ready-to-use, modular End-to-End (E2E) testing framework built in Go. It supports automated API testing via a type-safe HTTP client wrapper, UI automation using Selenium WebDriver following the Page Object Model (POM) and Persona/Action pattern, custom HTML/JSON reporting, and granular request/response logging.

---

## Features

- **Extensible API Client**: Auto-marshaling, struct pointer safety checks, and an interface-driven Authentication manager (Bearer, Basic, API Key, mTLS, and SSH signing).
- **Reusable API Action Helpers**: Pre-built assertions and action helpers in `pkg/api/actions/` (`GetAndExpectOK`, `PostAndExpectCreated`, `AssertNotEmpty`, etc.).
- **Selenium UI Integration**: Base Page Object wrappers handling dynamic CSS/XPath element selection, waiting hooks, interaction wrappers, and automated screenshots on test failure.
- **Persona & Action Pattern**: High-level declarative test actions (`pkg/ui/actions/`) representing realistic user journeys (Public, Member, Admin).
- **Batteries-Included Public Tests**: Ready-to-run tests against public internet endpoints (Google, GitHub, example.com for UI; httpbin, GitHub, JSONPlaceholder for API) so you can verify the framework immediately without setting up a backend.
- **Observability Logging**: Automatic date-organized file logging of request and response payloads at `evidence/run-<timestamp>/requests/`.
- **Interactive HTML & Markdown Reporting**: Responsive test results dashboard compiled automatically after each execution.

---

## Directory Structure

```
├── go.mod                     # Go module definitions
├── go.sum                     # Dependency checksums
├── Makefile                   # Execution shortcuts
├── config.json                # Environment configuration (adapt to your project)
├── docker-compose.yml         # Standalone Selenium Chrome container (optional)
├── pkg/
│   ├── config/
│   │   └── config.go          # Config loader and environment variable overrides
│   ├── logger/
│   │   └── logger.go          # Thread-safe logging levels (INFO, DEBUG, WARN, ERROR)
│   ├── client/
│   │   ├── auth.go            # Authentication interface and sub-types
│   │   └── client.go          # Custom HTTP client wrapper and JSON logger
│   ├── api/
│   │   └── actions/
│   │       └── api_actions.go # Reusable API test action helpers
│   ├── ui/
│   │   ├── driver.go          # Selenium WebDriver connection & option manager
│   │   ├── pom.go             # Page Object Model helper wrappers
│   │   ├── pages/
│   │   │   ├── home_page.go          # Example home page object
│   │   │   ├── login_page.go         # Example login page object
│   │   │   └── admin_dashboard_page.go # Example admin dashboard page object
│   │   └── actions/
│   │       ├── public_actions.go     # Persona & public user actions (GoToHome, VerifyPageTitle, etc.)
│   │       ├── member_actions.go     # Member persona skeleton actions
│   │       └── admin_actions.go      # Admin persona skeleton actions
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
│   │   ├── public_api_test.go # Working tests targeting public APIs (httpbin, GitHub, JSONPlaceholder)
│   │   └── example_api_test.go# Example skeleton tests for custom backend
│   └── ui/
│       ├── google/            # Google UI test package (1 test per file)
│       │   ├── main_test.go
│       │   ├── 01_search_journey_test.go
│       │   └── 02_search_direct_test.go
│       ├── github/            # GitHub UI test package (1 test per file)
│       │   ├── main_test.go
│       │   ├── 01_navigation_journey_test.go
│       │   └── 02_public_repo_test.go
│       └── example/           # Example UI test package (1 test per file)
│           ├── main_test.go
│           ├── 01_example_com_test.go
│           ├── 02_public_journey_test.go
│           └── 03_admin_login_journey_test.go
├── fixtures/                  # Generic test data files (CSV, PDF, images)
│   ├── example_bulk_upload.csv
│   ├── example_resource.pdf
│   └── example_image.png
├── docs/
│   ├── MASTER_API_TESTING_PROMPT.md  # AI prompt template for generating API tests
│   └── references/                   # Reference documentation
├── run-tests.sh               # Run both UI and API test suites with Chromedriver management
├── run-api-tests.sh           # Run API test suite
├── run-ui-tests.sh            # Run UI test suite with Chromedriver management
└── scripts/                   # Validation and helper scripts
    ├── validate-template.sh   # Validates template integrity and checks for domain leaks
    ├── install-graphify.sh    # Graphify installation script
    └── token-saver.sh         # Graphify knowledge graph extractor
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
   - Chromedriver is automatically launched and stopped by `./run-ui-tests.sh` and `./run-tests.sh`.

---

## Quick Start

1. **Clone this template**:
   ```bash
   ./clone.sh /path/to/your-new-test-project
   cd /path/to/your-new-test-project
   ```

2. **Run out-of-the-box public tests immediately**:
   ```bash
   # Run public API tests against httpbin, GitHub, JSONPlaceholder
   ./run-api-tests.sh -run TestAPI_Public

   # Run public UI tests against Google, GitHub, example.com
   ./run-ui-tests.sh -run TestUI_Public
   ```

3. **Configure your own application targets** in `config.json`:
   ```json
   {
     "baseUrl": "https://api.yourapp.com",
     "uiUrl": "https://yourapp.com",
     "seleniumUrl": "http://localhost:9515",
     "headless": false,
     "timeout": 10
   }
   ```

4. **Run all tests**:
   ```bash
   ./run-tests.sh
   ```

---

## Writing Tests

### 1. API Test Cases

Use the pre-built helpers in `pkg/api/actions/`:

```go
func TestAPI_GetUser(t *testing.T) {
    tests.RunAPITest(t, "GET /users/1 returns a valid user", func(t *testing.T, tc *tests.TestContext) {
        c := client.NewClient("https://jsonplaceholder.typicode.com", 10*time.Second, tests.ExecutionLogDir)

        var user UserResponse
        actions.GetAndExpectOK(tc, c, "/users/1", &user)
        actions.AssertNotEmpty(tc, "Username", user.Username)
        actions.AssertEquals(tc, "ID", user.ID, 1)
    })
}
```

### 2. UI Test Cases

Use the Persona/Action pattern from `pkg/ui/actions/`:

```go
func TestUI_ExampleComJourney(t *testing.T) {
    tests.RunUITest(t, "example.com Baseline Test", func(t *testing.T, page *ui.Page) {
        persona := actions.NewPublicPersona(page, "https://example.com", 10*time.Second)
        result := actions.NewResult("ExampleCom")

        actions.GoToHome(persona, result)
        actions.VerifyElementVisible(persona, result, "h1", "H1Heading")
        actions.VerifyPageTitle(persona, result, "Example Domain")

        if result.Failed() {
            t.Fatalf("Journey failed: %v", result.Error)
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
| `seleniumUrl` | `http://localhost:9515` | ChromeDriver address |
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

## Reports & Artifacts

After each test run, reports and evidence are automatically generated in `evidence/run-<timestamp>/`:

```
evidence/run-2026-09-06_17-48-05/
├── reports/
│   ├── report.html    # Interactive HTML dashboard (open in browser)
│   └── report.md      # Markdown summary
├── requests/          # Per-test request/response JSON logs
└── screenshots/       # Screenshots and failure captures (UI tests)
```

---

## Template Validation

Run the validation script to verify that the template builds cleanly and contains zero domain-specific leaks:

```bash
./scripts/validate-template.sh
```
