# E2E Framework Configuration & Customization Guide

This guide walks you through cloning, configuring, and extending the E2E testing framework for a new web app and API target. We use `https://yourapp.com` (Frontend UI) and `https://api.yourapp.com` (Backend API) as placeholder examples — replace these with your actual URLs.

---

## 1. Cloning and Initial Setup

Clone this repository and verify dependencies:

```bash
# Clone using the built-in clone script (excludes gitignored files)
./clone.sh /path/to/your-new-test-project
cd /path/to/your-new-test-project

# Install dependencies (Selenium bindings and other packages)
make deps
```

---

## 2. Configuring Targets

Update `config.json` in the root of the repository with your application's URLs:

```json
{
  "baseUrl": "https://api.yourapp.com",
  "uiUrl": "https://yourapp.com",
  "seleniumUrl": "http://localhost:9515",
  "headless": false,
  "timeout": 10,
  "adminCredentials": {
    "username": "admin@yourapp.com",
    "password": "your-admin-password"
  }
}
```

*Note: You can override these variables on the fly in CI environments using environment variables:*
```bash
E2E_BASE_URL=https://api.yourapp.com \
E2E_UI_URL=https://yourapp.com \
E2E_HEADLESS=true \
make test-all
```

---

## 3. UI Element Capturing Techniques

We use the **Page Object Model (POM)** pattern. To automate pages, you must identify stable element locators.

### Finding Elements via Chrome DevTools
1. Open Chrome, navigate to your app, right-click on any element (e.g., a Sign In button), and select **Inspect**.
2. In the **Elements** panel of DevTools:
   - Press `Ctrl + F` (or `Cmd + F` on macOS) to open the search bar.
   - Test your selector candidates (CSS or XPath) to ensure they match exactly **1 element**.

### Selector Strategy Guidelines
* **Prefer data-testid attributes**: Add `data-testid="my-button"` to your UI elements and use `testid:my-button` in tests. This is the most stable strategy.
* **Use IDs**: Selectors like `css:#login-email` or `css:#submit-btn` are fast and stable.
* **Use Clean CSS Classes**: Look for semantic class names like `css:.navbar-brand` or `css:.btn-primary`.
* **Avoid Auto-Generated Classes**: Avoid dynamic hashes (e.g. `css:.StyledButton-sc-1234a-0`). These change with every frontend build.
* **Use Attributes**: Elements without IDs often have descriptive attributes:
  - CSS: `css:input[name='email']` or `css:button[type='submit']`
* **XPath for Text/Hierarchies**:
  - Text Match: `xpath://button[contains(text(), 'Sign In')]`
  - Sibling navigation: `xpath://div[@class='card-body']/following-sibling::div/button`

---

## 4. Coding Page Objects (UI Testing)

Create a new file in `pkg/ui/pages/` representing a web page.

### Create Page Object: `pkg/ui/pages/profile_page.go`
```go
package pages

import (
    "time"

    "github.com/tebeka/selenium"
    "e2e-template/pkg/ui"
)

type ProfilePage struct {
    *ui.Page
    AvatarIcon   string
    ProfileName  string
    LogoutButton string
}

func NewProfilePage(driver selenium.WebDriver, screenshotDir string) *ProfilePage {
    return &ProfilePage{
        Page:         ui.NewPage(driver, screenshotDir),
        AvatarIcon:   "css:.user-avatar",
        ProfileName:  "css:#profile-header-name",
        LogoutButton: "xpath://button[text()='Logout']",
    }
}

func (p *ProfilePage) GetUsername(timeout time.Duration) (string, error) {
    return p.GetText(p.ProfileName, timeout)
}
```

### Write UI E2E Test: `tests/ui/profile_ui_test.go`
```go
package ui_test

import (
    "testing"
    "time"

    "e2e-template/pkg/ui"
    "e2e-template/pkg/ui/pages"
    "e2e-template/tests"
)

func TestUI_LoginFlow(t *testing.T) {
    tests.RunUITest(t, "User Login Flow", func(t *testing.T, page *ui.Page) {
        cfg := tests.GlobalConfig

        // Navigate to login page
        if err := page.Navigate(cfg.UiURL + "/login"); err != nil {
            t.Fatalf("Failed to load login page: %v", err)
        }

        loginPage := pages.NewLoginPage(page.Driver, page.ScreenshotDir)
        profilePage := pages.NewProfilePage(page.Driver, page.ScreenshotDir)

        // Perform login
        err := loginPage.Login(cfg.AdminCredentials.Username, cfg.AdminCredentials.Password, 5*time.Second)
        if err != nil {
            t.Fatalf("Login attempt failed: %v", err)
        }

        // Verify the user's name is displayed after login
        name, err := profilePage.GetUsername(5 * time.Second)
        if err != nil {
            t.Fatalf("Failed to fetch profile username: %v", err)
        }

        if name == "" {
            t.Errorf("Expected a non-empty profile username after login")
        }
    })
}
```

---

## 5. Structuring API Models (API Testing)

API models map backend JSON keys to Go structs.

### Where to put models?
* **Shared Models**: Place in a new folder `pkg/models/` (e.g. `pkg/models/user.go`) if reused across multiple test suites.
* **Test-Specific Models**: Place directly inside the test file (e.g. `tests/api/auth_test.go`) if only used in a single validation suite.

### Example: API Authentication Structs

#### Struct Design: `tests/api/types_test.go`
```go
package api_test

// LoginRequest defines the payload sent to POST /v1/auth/login
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// UserProfile models user metadata returned in the login response.
type UserProfile struct {
    ID        string `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
}

// LoginResponse defines the payload returned from POST /v1/auth/login
type LoginResponse struct {
    AccessToken string      `json:"access_token"`
    TokenType   string      `json:"token_type"`
    ExpiresIn   int         `json:"expires_in"`
    User        UserProfile `json:"user"`
}
```

#### Writing API E2E Test: `tests/api/auth_api_test.go`
```go
package api_test

import (
    "fmt"
    "testing"

    "e2e-template/tests"
)

func TestAPI_AuthenticateUser(t *testing.T) {
    tests.RunAPITestWithDetails(
        t,
        "POST /auth/login returns an access token",
        "Verifies that valid admin credentials produce a non-empty access token.",
        "HTTP 200 OK with non-empty access_token",
        func(tc *tests.TestContext) {
            cfg := tests.GlobalConfig
            reqPayload := &LoginRequest{
                Email:    cfg.AdminCredentials.Username,
                Password: cfg.AdminCredentials.Password,
            }
            var resp LoginResponse

            err := tc.Client.SendHttpRequest("POST", "/v1/auth/login", nil, reqPayload, &resp, nil)
            if err != nil {
                tc.FailureReason = fmt.Sprintf("Login request failed: %v", err)
                tc.Fatalf("Login request failed: %v", err)
            }
            tc.Actual = fmt.Sprintf("HTTP 200 OK, token=%q", resp.AccessToken)
            if resp.AccessToken == "" {
                tc.FailureReason = "Expected non-empty access_token"
                tc.Errorf("Expected non-empty access_token")
            }
        },
    )
}
```

---

## 6. Execution Flow and Verification

1. Start your local Chromedriver:
   ```bash
   chromedriver --port=9515
   ```

2. Run your tests:
   ```bash
   # Run all test suites (API + UI)
   make test-all

   # Run only API tests
   make test-api

   # Run only UI tests
   make test-ui

   # Run a specific test by name
   go test -v ./tests/... -run=TestAPI_AuthenticateUser
   ```

3. Open the generated HTML dashboard to review results:
   ```bash
   # The path is printed to terminal after the run
   open evidence/run-<timestamp>/reports/report.html
   ```

---

## 7. Seeding Test Data

If your test suite needs pre-populated data (e.g., a test user in the database), override the `seedTestData()` function in `tests/helpers.go`:

```go
// In tests/helpers.go — replace the empty stub:
func seedTestData() {
    if err := createTestUser(GlobalConfig); err != nil {
        fmt.Printf("WARNING: Failed to seed test user: %v\n", err)
    }
}
```

---

## 8. Adding New Test Files

Follow this naming convention for test files:
- API tests: `tests/api/XX-feature-name_test.go` (e.g. `tests/api/02-auth_test.go`)
- UI tests: `tests/ui/XX-journey-name_test.go` (e.g. `tests/ui/02-login-flow_test.go`)
- Page objects: `pkg/ui/pages/feature_page.go`
- Action helpers: `pkg/ui/actions/persona_actions.go`
