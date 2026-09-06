package ui_test

// ─────────────────────────────────────────────────────────────────────────────
// TEMPLATE: Example UI Tests (YOUR APPLICATION)
//
// This file shows how to write UI tests against YOUR local/staging web app
// using the action-based persona/result pattern.
//
// Tests are skipped automatically when no server is running at uiUrl.
//
// HOW TO ADAPT:
//  1. Set uiUrl in config.json to your web app's URL.
//  2. Create action helpers in pkg/ui/actions/ for each user persona.
//  3. Write tests as sequential action calls (reads like human steps).
//  4. Add more test files: tests/ui/XX-journey_test.go
//
// Run with:
//   go test -v ./tests/ui/... -run TestUI_01
// ─────────────────────────────────────────────────────────────────────────────

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"e2e-template/pkg/ui"
	"e2e-template/pkg/ui/actions"
	"e2e-template/tests"
)

// isUIServerRunning checks whether the configured uiUrl is reachable.
func isUIServerRunning() bool {
	if tests.GlobalConfig == nil || tests.GlobalConfig.UiURL == "" {
		return false
	}
	uiURL := tests.GlobalConfig.UiURL
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
// Example 1: Public User Journeys (Action-Based Pattern)
//
// This is the recommended pattern for writing UI tests:
//   1. Create a persona (who is performing the actions)
//   2. Create a result collector
//   3. Call sequential action functions
//   4. Assert on the result
// ─────────────────────────────────────────────────────────────────────────────

// TestUI_01_PublicJourneys demonstrates the action-based testing pattern.
// TODO: Replace the placeholder actions with your application's user flows.
func TestUI_01_PublicJourneys(t *testing.T) {
	if !isUIServerRunning() {
		t.Skipf("Skipping: no server running at %s (set uiUrl in config.json)", tests.GlobalConfig.UiURL)
	}
	tests.RunUITest(t, "Verify Public User Journeys", func(t *testing.T, page *ui.Page) {
		cfg := tests.GlobalConfig

		// Setup the persona and the result collector
		pubPersona := actions.NewPublicPersona(page, cfg.UiURL, 5*time.Second)
		result := actions.NewResult("TestUI_01_PublicJourneys")

		// Run simple sequential action calls
		// TODO: Replace these with your application's navigation actions.
		// Example:
		//   actions.GoToHome(pubPersona, result)
		//   actions.GoToPage(pubPersona, result, "testid-products-link", "Products")
		//   actions.GoToHome(pubPersona, result)
		//   actions.GoToPage(pubPersona, result, "testid-about-link", "About")
		actions.GoToHome(pubPersona, result)
		actions.VerifyElementVisible(pubPersona, result, "css:body", "PageBody")

		// Assert on captured results
		if result.Failed() {
			t.Errorf("Test Journey Failed: %v", result.Error)
			t.Errorf("Actions Attempted: %v", result.Actions)
			t.Errorf("Evidence Captured: %v", result.Evidence)
			t.Fatalf("Advice / Remediation: %v", result.Advice)
		}
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Example 2: Admin Login Journey (Skeleton)
// ─────────────────────────────────────────────────────────────────────────────

// TestUI_02_AdminLoginJourney demonstrates the admin login flow.
// TODO: Configure admin credentials in config.json and uncomment the LoginAsAdmin call.
func TestUI_02_AdminLoginJourney(t *testing.T) {
	if !isUIServerRunning() {
		t.Skipf("Skipping: no server running at %s (set uiUrl in config.json)", tests.GlobalConfig.UiURL)
	}
	tests.RunUITest(t, "Admin Login Journey", func(t *testing.T, page *ui.Page) {
		cfg := tests.GlobalConfig

		adminPersona := actions.NewAdminPersona(page, cfg.UiURL, 5*time.Second)
		result := actions.NewResult("TestUI_02_AdminLoginJourney")

		// Navigate to home first
		actions.GoToHome(adminPersona, result)

		// TODO: Uncomment once admin login is configured:
		// actions.LoginAsAdmin(adminPersona, cfg, result)
		// actions.GoToPage(adminPersona, result, "testid-dashboard-link", "Dashboard")

		if result.Failed() {
			t.Errorf("Admin login journey failed: %v\nActions: %v\nAdvice: %v",
				result.Error, result.Actions, result.Advice)
		} else {
			t.Logf("✅ Admin login journey completed (configure admin login to enable)")
		}
	})
}
