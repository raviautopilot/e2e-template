package github_test

import (
	"testing"
	"time"

	"e2e-template/pkg/ui"
	"e2e-template/pkg/ui/actions"
	"e2e-template/tests"
)

// TestUI_GitHub_01_NavigationJourney demonstrates verifying GitHub navigation and home page.
func TestUI_GitHub_01_NavigationJourney(t *testing.T) {
	tests.RunUITest(t, "GitHub Navigation Journey", func(t *testing.T, page *ui.Page) {
		persona := actions.NewPublicPersona(page, "https://github.com", 10*time.Second)
		result := actions.NewResult("GitHubJourney")

		// Step 1: Load GitHub homepage and verify title
		actions.GoToHome(persona, result)
		actions.VerifyElementVisible(persona, result, "css:body", "GitHubBody")
		actions.VerifyPageTitle(persona, result, "GitHub")

		if result.Failed() {
			t.Errorf("GitHub journey failed: %v\nActions: %v\nAdvice: %v",
				result.Error, result.Actions, result.Advice)
		} else {
			t.Logf("✅ GitHub homepage journey completed: %v", result.Actions)
		}
	})
}
