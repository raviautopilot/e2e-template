# E2E Framework Configuration & Customization Guide

This guide walks you through cloning, configuring, and extending the E2E testing framework for any web application and backend API target.

---

## 1. Cloning and Initial Setup

Clone this repository and verify dependencies:

```bash
# Clone using the built-in clone script (excludes gitignored files)
./clone.sh /path/to/your-new-test-project
cd /path/to/your-new-test-project

# Install Go module dependencies
make deps
```

---

## 2. Scaffolding a New Service Test Suite (Interactive Generator)

Instead of manually creating test folders and copying files, use the built-in interactive generator script:

```bash
# 1. Interactive Mode (prompts for service name and type: API, UI, or Both)
./create-service.sh

# 2. Non-interactive CLI Mode (for automated scripts or CI/CD)
./create-service.sh orders --api
./create-service.sh payments --ui
./create-service.sh checkout --all

# 3. Makefile Shortcut
make new-service name=billing
```

### What It Generates Automatically:
* **For API**:
  * `tests/api/<service>/main_test.go` — Test suite lifecycle bootstrap (`TestMain`).
  * `tests/api/<service>/types_test.go` — Request/Response models with JSON tags.
  * `tests/api/<service>/01_health_check_test.go` — Health endpoint test.
  * `tests/api/<service>/02_parameterized_test.go` — Table-driven parameterized test matrix.
* **For UI**:
  * `pkg/ui/pages/<service>_page.go` — Page Object Model (POM) with selectors and action methods.
  * `tests/ui/<service>/main_test.go` — UI suite lifecycle runner.
  * `tests/ui/<service>/01_<service>_journey_test.go` — Browser journey test.

---

## 3. Configuring Targets & Credentials

Update `config.json` in the root of the repository with your application's URLs and placeholder credentials:

```json
{
  "baseUrl": "https://api.yourdomain.com",
  "uiUrl": "https://yourdomain.com",
  "seleniumUrl": "http://localhost:9515",
  "headless": false,
  "timeout": 10,
  "adminCredentials": {
    "username": "admin@yourdomain.com",
    "password": "REPLACE_WITH_YOUR_PASSWORD"
  }
}
```

### 🔒 Best Practice: Supplying Credentials Securely
To prevent committing passwords to GitHub, keep `config.json` with placeholder values and pass credentials dynamically via environment variables:

```bash
E2E_BASE_URL=https://api.yourdomain.com \
E2E_ADMIN_USERNAME=admin@yourdomain.com \
E2E_ADMIN_PASSWORD=your-secret-password \
./run-api-tests.sh <service-name>
```

---

## 4. Writing API Tests (`tests/api/<service>/`)

Every service has its own dedicated package directory under `tests/api/<service>/`.

### 4.1 Request & Response Models (`types_test.go`)
Define clean Go structs mapping your JSON request and response payloads:

```go
package myservice_test

type LoginRequest struct {
    Email    string `json:"email,omitempty"`
    Password string `json:"password,omitempty"`
}

type LoginResponse struct {
    Token       string                 `json:"token,omitempty"`
    AccessToken string                 `json:"access_token,omitempty"`
    User        map[string]interface{} `json:"user,omitempty"`
}
```

### 4.2 Parameterized / Table-Driven Tests (`01_login_parameterized_test.go`)
Use the `tests.RunAPITestWithDetails` runner for rich failure tracking and automatic request/response logging:

```go
package myservice_test

import (
    "fmt"
    "testing"
    "time"

    "e2e-template/pkg/api/actions"
    "e2e-template/pkg/client"
    "e2e-template/tests"
)

func TestAPI_01_AdminLogin(t *testing.T) {
    baseURL := tests.GlobalConfig.BaseURL
    apiClient := client.NewClient(baseURL, 15*time.Second, tests.ExecutionLogDir)

    tests.RunAPITestWithDetails(
        t,
        "Admin Login - Valid Credentials",
        "Submits admin credentials loaded dynamically from configuration.",
        "HTTP 200 OK with Auth Token",
        func(tc *tests.TestContext) {
            testName := "Admin Login - Valid Credentials"
            apiClient.SetTestName(testName)
            tc.Client = apiClient

            req := LoginRequest{
                Email:    tests.GlobalConfig.AdminCredentials.Username,
                Password: tests.GlobalConfig.AdminCredentials.Password,
            }
            var resp LoginResponse

            // Pass pointer to struct &req
            actions.PostAndExpectOK(tc, apiClient, "/api/v1/auth/login", &req, &resp)
        },
    )
}
```

---

## 5. Writing UI Tests (`tests/ui/<service>/`)

UI tests use Selenium WebDriver with the **Page Object Model (POM)** pattern.

### 5.1 Page Object (`pkg/ui/pages/<service>_page.go`)
Define element locators (CSS, XPath, or `data-testid`) and helper interaction methods:

```go
package pages

import (
    "time"
    "e2e-template/pkg/ui"
)

type LoginPage struct {
    *ui.Page
    EmailInput    string
    PasswordInput string
    SubmitBtn     string
}

func NewLoginPage(page *ui.Page) *LoginPage {
    return &LoginPage{
        Page:          page,
        EmailInput:    "css:input[name='email']",
        PasswordInput: "css:input[name='password']",
        SubmitBtn:     "css:button[type='submit']",
    }
}

func (p *LoginPage) Login(email, password string, timeout time.Duration) error {
    if err := p.SendKeys(p.EmailInput, email, timeout); err != nil {
        return err
    }
    if err := p.SendKeys(p.PasswordInput, password, timeout); err != nil {
        return err
    }
    return p.Click(p.SubmitBtn, timeout)
}
```

### 5.2 UI Journey Test (`tests/ui/<service>/01_login_journey_test.go`)
```go
package myservice_test

import (
    "testing"
    "time"

    "e2e-template/pkg/ui"
    "e2e-template/pkg/ui/pages"
    "e2e-template/tests"
)

func TestUI_LoginJourney(t *testing.T) {
    tests.RunUITest(t, "Admin Login Journey", func(t *testing.T, page *ui.Page) {
        loginPage := pages.NewLoginPage(page)
        
        // Navigate to UI URL
        if err := page.GoToHome(tests.GlobalConfig.UiURL + "/login"); err != nil {
            t.Fatalf("Failed to navigate: %v", err)
        }

        err := loginPage.Login(
            tests.GlobalConfig.AdminCredentials.Username,
            tests.GlobalConfig.AdminCredentials.Password,
            5*time.Second,
        )
        if err != nil {
            t.Fatalf("Login action failed: %v", err)
        }
    })
}
```

---

## 6. Running Tests & Viewing Reports

### 6.1 Running API Tests
```bash
# Run a specific service test package:
./run-api-tests.sh myservice

# Run all API tests in the repository:
./run-api-tests.sh

# Pass runtime credentials:
E2E_ADMIN_PASSWORD=yourpassword ./run-api-tests.sh myservice
```

### 6.2 Running UI Tests
1. Start Selenium / Chromedriver (or Docker standalone container):
   ```bash
   # Option A: Local chromedriver
   chromedriver --port=9515

   # Option B: Docker container
   docker compose up -d
   ```
2. Execute UI tests:
   ```bash
   # Run a specific service UI test:
   ./run-ui-tests.sh myservice

   # Run headless mode (no browser popup):
   E2E_HEADLESS=true ./run-ui-tests.sh myservice
   ```

### 6.3 Viewing Reports & Evidence
After every test run, reports and evidence are automatically compiled into the `evidence/` directory:
- **HTML Dashboard**: `evidence/run-<timestamp>/reports/report.html` (interactive dashboard with clickable request/response logs and failure screenshots)
- **Markdown Summary**: `evidence/run-<timestamp>/reports/report.md`
- **JSON Raw Log**: `evidence/run-<timestamp>/reports/report.json`
- **Raw Request/Response Payloads**: `evidence/run-<timestamp>/requests/`
- **Failure Screenshots**: `evidence/run-<timestamp>/screenshots/`

---

## 7. Package and Directory Layout Conventions

Always keep the repository clean and modular:

```text
tests/
├── api/
│   ├── example/          # Reference template API tests
│   └── <service-name>/   # Your custom service API tests
└── ui/
    ├── example/          # Reference template UI tests
    └── <service-name>/   # Your custom service UI tests
pkg/
├── api/actions/          # HTTP verb action helpers (GetAndExpectOK, etc.)
├── client/               # Custom HTTP client & request/response logger
├── config/               # Configuration & env var override loader
├── ui/pages/             # Page Object Model definitions
└── report/               # HTML & Markdown report compiler
```
