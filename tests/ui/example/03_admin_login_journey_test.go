package example_test

import (
	"testing"
	"time"

	"e2e-template/pkg/ui"
	"e2e-template/pkg/ui/actions"
	"e2e-template/tests"
)

// TestUI_Example_03_AdminLoginJourney demonstrates the admin login flow skeleton.
// TODO: Configure admin credentials in config.json and uncomment the LoginAsAdmin call.
func TestUI_Example_03_AdminLoginJourney(t *testing.T) {
	if !isUIServerRunning() {
		t.Skipf("Skipping: no server running at %s (set uiUrl in config.json)", tests.GlobalConfig.UiURL)
	}
	tests.RunUITest(t, "Admin Login Journey", func(t *testing.T, page *ui.Page) {
		cfg := tests.GlobalConfig

		adminPersona := actions.NewAdminPersona(page, cfg.UiURL, 5*time.Second)
		result := actions.NewResult("TestUI_02_AdminLoginJourney")

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
