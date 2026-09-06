package example_test

import (
	"testing"
	"time"

	"e2e-template/pkg/ui"
	"e2e-template/pkg/ui/actions"
	"e2e-template/tests"
)

// TestUI_Example_01_ExampleCom verifies example.com loads with correct content.
func TestUI_Example_01_ExampleCom(t *testing.T) {
	tests.RunUITest(t, "example.com Baseline Test", func(t *testing.T, page *ui.Page) {
		persona := actions.NewPublicPersona(page, "https://example.com", 5*time.Second)
		result := actions.NewResult("ExampleCom")

		actions.GoToHome(persona, result)
		actions.VerifyElementVisible(persona, result, "css:h1", "H1Heading")
		actions.VerifyPageTitle(persona, result, "Example Domain")

		h1Text := actions.GetElementText(persona, result, "css:h1", "H1Text")
		if h1Text != "" {
			t.Logf("✅ example.com h1: %q", h1Text)
		}

		if result.Failed() {
			t.Errorf("example.com test failed: %v\nAdvice: %v", result.Error, result.Advice)
		} else {
			t.Logf("✅ example.com baseline passed: %v", result.Actions)
		}
	})
}
