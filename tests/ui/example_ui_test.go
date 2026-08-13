package ui_test

// ─────────────────────────────────────────────────────────────────────────────
// TEMPLATE: Example UI Tests (Page Object Model)
//
// This file shows how to write UI/browser tests using the e2e-template framework.
// Tests use the Page Object Model (POM) pattern via pkg/ui/pages and pkg/ui/actions.
//
// HOW TO ADAPT:
//  1. Set uiUrl in config.json to your web app's URL.
//  2. Create page objects in pkg/ui/pages/ for each page of your application.
//  3. Create action helpers in pkg/ui/actions/ for each user persona.
//  4. Add more test files following this pattern: XX-journey_test.go
//
// KEY PAGE METHODS (from pkg/ui/pom.go):
//  - page.GoToHome(url)                                   → navigate to a URL
//  - page.WaitUntilVisible(locator, timeout)              → wait for element
//  - page.Click(locator, timeout)                         → click element
//  - page.SendKeys(locator, text, timeout)                → type text into input
//  - page.GetText(locator, timeout)                       → read element text
//  - page.CaptureScreenshot(stepName)                     → save a step screenshot
//  - page.ClickByTestID(testID, timeout)                  → click by data-testid
//  - page.SendKeysByTestID(testID, text, timeout)         → type by data-testid
//  - page.VerifyFormElements(testIDs, timeout)            → check elements exist
//
// LOCATOR SYNTAX:
//  - "css:.my-class"            → CSS selector
//  - "css:#my-id"               → CSS id
//  - "xpath://button[text()='OK']" → XPath
//  - "testid:my-testid"         → data-testid attribute (preferred)
//
// Run with:
//   go test -v ./tests/ui/...
// ─────────────────────────────────────────────────────────────────────────────

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"e2e-template/pkg/ui"
	"e2e-template/tests"
)

// isUIServerRunning checks whether the configured uiUrl is reachable.
// Returns false when uiUrl points to localhost/127.0.0.1 and no server answers.
// This lets the template example tests skip gracefully in a fresh checkout.
func isUIServerRunning() bool {
	if tests.GlobalConfig == nil || tests.GlobalConfig.UiURL == "" {
		return false
	}
	uiURL := tests.GlobalConfig.UiURL
	// If it's a localhost URL and nothing answers, skip
	if strings.Contains(uiURL, "localhost") || strings.Contains(uiURL, "127.0.0.1") {
		c := &http.Client{Timeout: 2 * time.Second}
		resp, err := c.Get(uiURL)
		if err != nil {
			return false
		}
		resp.Body.Close()
	}
	return true
}

// ─────────────────────────────────────────────────────────────────────────────
// Example 1: Verify the home page loads and contains expected content.
// ─────────────────────────────────────────────────────────────────────────────

// TestUI_01_HomePageLoads navigates to the configured uiUrl and verifies the page loads.
// TODO: Replace the CSS/XPath selectors with ones from your own application.
func TestUI_01_HomePageLoads(t *testing.T) {
	if !isUIServerRunning() {
		t.Skipf("Skipping: no server running at %s (set uiUrl in config.json)", tests.GlobalConfig.UiURL)
	}
	tests.RunUITest(t, "Home Page Loads Successfully", func(t *testing.T, page *ui.Page) {
		cfg := tests.GlobalConfig

		// Navigate to the home page
		if err := page.GoToHome(cfg.UiURL); err != nil {
			t.Fatalf("Failed to navigate to home page %s: %v", cfg.UiURL, err)
		}

		// Capture initial screenshot
		if _, err := page.CaptureScreenshot("01-home-page-loaded"); err != nil {
			t.Logf("WARNING: Could not capture screenshot: %v", err)
		}

		// TODO: Replace "css:body" with a selector specific to your application's landing element.
		// Example: page.WaitUntilVisible("css:body", 5*time.Second)
		//      or: page.WaitUntilVisible("testid:your-hero-section", 5*time.Second)
		if _, err := page.WaitUntilVisible("css:body", 5*time.Second); err != nil {
			t.Errorf("Home page body element not found: %v", err)
			return
		}

		t.Logf("✅ Home page loaded at %s", cfg.UiURL)
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Example 2: Verify that login UI elements are present.
// Uses the login page configuration from config.json.
// ─────────────────────────────────────────────────────────────────────────────

// TestUI_02_LoginPageElements verifies login form elements exist on the login page.
// TODO: Configure the selector test IDs in config.json (adminLoginUsernameInputTestID, etc.).
func TestUI_02_LoginPageElements(t *testing.T) {
	if !isUIServerRunning() {
		t.Skipf("Skipping: no server running at %s (set uiUrl in config.json)", tests.GlobalConfig.UiURL)
	}
	tests.RunUITest(t, "Login Page Elements Are Present", func(t *testing.T, page *ui.Page) {
		cfg := tests.GlobalConfig

		// Navigate to the home page first
		if err := page.GoToHome(cfg.UiURL); err != nil {
			t.Fatalf("Failed to navigate to %s: %v", cfg.UiURL, err)
		}

		// Take a screenshot of the initial state
		_, _ = page.CaptureScreenshot("01-landing")

		// Uncomment and adapt the code below once you've configured your login test IDs.
		// The example uses ClickByTestID which maps to data-testid attributes.
		//
		// Example: Click the admin login button to open the login form
		// if cfg.AdminLoginButtonTestID != "" {
		// 	if err := page.ClickByTestID(cfg.AdminLoginButtonTestID, 5*time.Second); err != nil {
		// 		t.Errorf("Admin login button not found: %v", err)
		// 		return
		// 	}
		// 	_, _ = page.CaptureScreenshot("02-admin-login-modal-open")
		//
		// 	// Verify login form fields are visible
		// 	loginTestIDs := cfg.AdminLoginTestIDs
		// 	if err := page.VerifyFormElements(loginTestIDs, 5*time.Second); err != nil {
		// 		t.Errorf("Login form elements not found: %v", err)
		// 	}
		// }

		t.Logf("✅ Login page elements check passed (configure test IDs in config.json to enable)")
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Example 3: Multi-step user journey (skeleton)
// Shows how to chain multiple page interactions in a single test.
// ─────────────────────────────────────────────────────────────────────────────

// TestUI_03_ExampleJourney demonstrates a multi-step navigation journey.
// Replace each step with real interactions for your application.
func TestUI_03_ExampleJourney(t *testing.T) {
	if !isUIServerRunning() {
		t.Skipf("Skipping: no server running at %s (set uiUrl in config.json)", tests.GlobalConfig.UiURL)
	}
	tests.RunUITest(t, "Example Multi-Step User Journey", func(t *testing.T, page *ui.Page) {
		cfg := tests.GlobalConfig

		// Step 1: Navigate to the app
		if err := page.GoToHome(cfg.UiURL); err != nil {
			t.Fatalf("Failed to navigate: %v", err)
		}
		_, _ = page.CaptureScreenshot("step-01-home")
		t.Logf("Step 1: Navigated to %s", cfg.UiURL)

		// Step 2: TODO - Add your application-specific interactions here.
		// Examples using the pom.go API:
		//
		//   // Click a nav link by testid
		//   page.ClickByTestID(cfg.AdminLoginButtonTestID, 5*time.Second)
		//
		//   // Type into an input field
		//   page.SendKeysByTestID(cfg.AdminLoginUsernameInputTestID, cfg.AdminCredentials.Username, 5*time.Second)
		//   page.SendKeysByTestID(cfg.AdminLoginPasswordInputTestID, cfg.AdminCredentials.Password, 5*time.Second)
		//
		//   // Submit the form
		//   page.ClickByTestID(cfg.AdminLoginSubmitButtonTestID, 5*time.Second)
		//
		//   // Wait for the dashboard to appear
		//   page.WaitUntilVisible("testid:dashboard-header", 10*time.Second)

		_, _ = page.CaptureScreenshot("step-02-completed")
		t.Logf("Step 2: Journey completed (configure your steps in this file)")

		// Prevent unused import errors while the steps above are commented out.
		_ = time.Second
	})
}
