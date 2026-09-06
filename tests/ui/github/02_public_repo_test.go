package github_test

import (
	"testing"
	"time"

	"e2e-template/pkg/ui"
	"e2e-template/pkg/ui/actions"
	"e2e-template/tests"
)

// TestUI_GitHub_02_PublicRepo verifies a known public repository page.
func TestUI_GitHub_02_PublicRepo(t *testing.T) {
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
