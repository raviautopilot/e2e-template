package google_test

import (
	"testing"
	"time"

	"e2e-template/pkg/ui"
	"e2e-template/pkg/ui/actions"
	"e2e-template/tests"
)

// TestUI_Google_02_SearchDirect demonstrates navigating directly
// to a search results URL (bypasses homepage consent dialogs).
func TestUI_Google_02_SearchDirect(t *testing.T) {
	tests.RunUITest(t, "Google Direct Search Results", func(t *testing.T, page *ui.Page) {
		persona := actions.NewPublicPersona(page, "https://www.google.com/search?q=e2e+testing+framework+golang&hl=en", 10*time.Second)
		result := actions.NewResult("GoogleDirectSearch")

		actions.GoToHome(persona, result)
		actions.VerifyElementVisible(persona, result, "css:body", "SearchResultsBody")

		if result.Failed() {
			t.Errorf("Direct search failed: %v\nAdvice: %v", result.Error, result.Advice)
		} else {
			t.Logf("✅ Google direct search page loaded: %v", result.Actions)
		}
	})
}
