package google_test

import (
	"testing"
	"time"

	"e2e-template/pkg/ui"
	"e2e-template/pkg/ui/actions"
	"e2e-template/tests"
)

// TestUI_Google_01_SearchJourney demonstrates a complete Google search journey
// using the action-based pattern. Each step reads like a human action.
func TestUI_Google_01_SearchJourney(t *testing.T) {
	tests.RunUITest(t, "Google Search Journey", func(t *testing.T, page *ui.Page) {
		persona := actions.NewPublicPersona(page, "https://www.google.com", 10*time.Second)
		result := actions.NewResult("GoogleJourney")

		// Human-readable journey: Go to Google -> Verify search box -> Type query -> Submit
		actions.GoToHome(persona, result)
		actions.VerifyGoogleSearchBox(persona, result)
		actions.TypeInGoogleSearchBox(persona, result, "e2e testing framework golang")
		actions.SubmitGoogleSearch(persona, result)

		if result.Failed() {
			t.Errorf("Google journey failed: %v\nActions: %v\nAdvice: %v",
				result.Error, result.Actions, result.Advice)
		} else {
			t.Logf("✅ Google search journey completed: %v", result.Actions)
		}
	})
}
