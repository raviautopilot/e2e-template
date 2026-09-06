package example_test

import (
	"testing"
	"time"

	"e2e-template/pkg/ui"
	"e2e-template/pkg/ui/actions"
	"e2e-template/tests"
)

// TestUI_Example_02_PublicJourneys demonstrates the action-based testing pattern on a custom app.
// TODO: Replace the placeholder actions with your application's user flows.
func TestUI_Example_02_PublicJourneys(t *testing.T) {
	if !isUIServerRunning() {
		t.Skipf("Skipping: no server running at %s (set uiUrl in config.json)", tests.GlobalConfig.UiURL)
	}
	tests.RunUITest(t, "Verify Public User Journeys", func(t *testing.T, page *ui.Page) {
		cfg := tests.GlobalConfig

		pubPersona := actions.NewPublicPersona(page, cfg.UiURL, 5*time.Second)
		result := actions.NewResult("TestUI_01_PublicJourneys")

		actions.GoToHome(pubPersona, result)
		actions.VerifyElementVisible(pubPersona, result, "css:body", "PageBody")

		if result.Failed() {
			t.Errorf("Test Journey Failed: %v", result.Error)
			t.Errorf("Actions Attempted: %v", result.Actions)
			t.Errorf("Evidence Captured: %v", result.Evidence)
			t.Fatalf("Advice / Remediation: %v", result.Advice)
		}
	})
}
