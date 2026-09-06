package ui_test

// ─────────────────────────────────────────────────────────────────────────────
// PUBLIC SITE UI TESTS
//
// These tests target real, publicly-available websites so the template works
// out-of-the-box without any local server running.
//
// Sites tested:
//   - https://google.com  — Search interaction
//   - https://github.com  — Navigation and page verification
//   - https://example.com — Minimal baseline / sanity test
//
// Tests use the action-based persona/result pattern:
//   1. Create a persona (who is performing the actions)
//   2. Create a result collector (tracks actions, evidence, advice)
//   3. Call sequential action functions (reads like human steps)
//   4. Assert on the result at the end
//
// Run with:
//   go test -v ./tests/ui/... -run TestUI_Public
// ─────────────────────────────────────────────────────────────────────────────

import (
	"testing"
	"time"

	"e2e-template/pkg/ui"
	"e2e-template/pkg/ui/actions"
	"e2e-template/tests"
)

// ─────────────────────────────────────────────────────────────────────────────
// 01 · Google — Search Journey
// ─────────────────────────────────────────────────────────────────────────────

// TestUI_Public_01_GoogleJourney demonstrates a complete Google search journey
// using the action-based pattern. Each step reads like a human action.
func TestUI_Public_01_GoogleJourney(t *testing.T) {
	tests.RunUITest(t, "Google Search Journey", func(t *testing.T, page *ui.Page) {
		persona := actions.NewPublicPersona(page, "https://www.google.com", 10*time.Second)
		result := actions.NewResult("GoogleJourney")

		// Human-readable journey: Go to Google → Verify search box → Type query → Submit
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

// TestUI_Public_02_GoogleSearchDirect demonstrates navigating directly
// to a search results URL (bypasses homepage consent dialogs).
func TestUI_Public_02_GoogleSearchDirect(t *testing.T) {
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

// ─────────────────────────────────────────────────────────────────────────────
// 02 · GitHub — Navigation Journey
// ─────────────────────────────────────────────────────────────────────────────

// TestUI_Public_03_GitHubJourney demonstrates verifying multiple GitHub pages.
func TestUI_Public_03_GitHubJourney(t *testing.T) {
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

// TestUI_Public_04_GitHubPublicRepo verifies a known public repository page.
func TestUI_Public_04_GitHubPublicRepo(t *testing.T) {
	tests.RunUITest(t, "GitHub Public Repo Page (octocat/Hello-World)", func(t *testing.T, page *ui.Page) {
		persona := actions.NewPublicPersona(page, "https://github.com/octocat/Hello-World", 10*time.Second)
		result := actions.NewResult("GitHubPublicRepo")

		actions.GoToHome(persona, result)
		actions.VerifyElementVisible(persona, result, "css:body", "RepoPageBody")
		actions.VerifyPageTitle(persona, result, "Hello-World")

		if result.Failed() {
			t.Errorf("GitHub repo page failed: %v\nAdvice: %v", result.Error, result.Advice)
		} else {
			t.Logf("✅ GitHub public repo page loaded: %v", result.Actions)
		}
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// 03 · Example.com — Simplest possible baseline test
// ─────────────────────────────────────────────────────────────────────────────

// TestUI_Public_05_ExampleCom verifies example.com loads with correct content.
func TestUI_Public_05_ExampleCom(t *testing.T) {
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
