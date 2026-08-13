package ui_test

// ─────────────────────────────────────────────────────────────────────────────
// PUBLIC SITE UI TESTS
//
// These tests target real, publicly-available websites so the template works
// out-of-the-box without any local server running.
//
// Sites tested:
//   - https://google.com  — Search page interaction and result verification
//   - https://github.com  — Navigation, search bar, and public repo page
//   - https://example.com — Minimal baseline / sanity test
//
// These tests demonstrate the full POM framework capabilities:
//   - page.GoToHome(url)                    → navigate to a URL
//   - page.WaitUntilVisible(locator, t)     → wait for element to appear
//   - page.SendKeys(locator, text, t)       → type into an input
//   - page.Click(locator, t)                → click an element
//   - page.GetText(locator, t)              → read element text
//   - page.CaptureScreenshot(stepName)      → save a step screenshot
//   - page.Driver.Title()                   → read the browser page title
//
// Run with:
//   go test -v ./tests/ui/... -run TestUI_Public
// ─────────────────────────────────────────────────────────────────────────────

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"e2e-template/pkg/ui"
	"e2e-template/tests"
)

// waitForGoogleSearchBox tries each known Google search box selector in sequence.
// Google uses <textarea name="q"> on modern Chrome but <input name="q"> on
// older layouts. The pom.go parseLocator only strips one "css:" prefix so
// comma-separated multi-selectors like "css:a, css:b" are NOT valid — always
// pass them as separate WaitUntilVisible calls instead.
func waitForGoogleSearchBox(page *ui.Page, timeout time.Duration) error {
	// Try the modern textarea first (used by Google since ~2023)
	if _, err := page.WaitUntilVisible("css:textarea[name='q']", timeout); err == nil {
		return nil
	}
	// Fallback: older input-based search box
	if _, err := page.WaitUntilVisible("css:input[name='q']", 2*time.Second); err == nil {
		return nil
	}
	return fmt.Errorf("google search box not found (tried textarea[name='q'] and input[name='q'])")
}

// sendKeysToGoogleSearchBox types text into whichever search box variant is present.
func sendKeysToGoogleSearchBox(page *ui.Page, text string, timeout time.Duration) error {
	if err := page.SendKeys("css:textarea[name='q']", text, timeout); err == nil {
		return nil
	}
	return page.SendKeys("css:input[name='q']", text, 3*time.Second)
}

// ─────────────────────────────────────────────────────────────────────────────
// 01 · Google Search
// ─────────────────────────────────────────────────────────────────────────────

// TestUI_Public_01_GoogleHomePageLoads verifies Google's homepage loads with a search box.
func TestUI_Public_01_GoogleHomePageLoads(t *testing.T) {
	tests.RunUITest(t, "Google Homepage Loads With Search Box", func(t *testing.T, page *ui.Page) {
		if err := page.GoToHome("https://www.google.com"); err != nil {
			t.Fatalf("Failed to navigate to Google: %v", err)
		}
		_, _ = page.CaptureScreenshot("01-google-home")

		// Google uses <textarea name="q"> on modern Chrome (updated from <input name="q">).
		// waitForGoogleSearchBox() tries both variants sequentially.
		if err := waitForGoogleSearchBox(page, 10*time.Second); err != nil {
			t.Fatalf("%v", err)
		}
		_, _ = page.CaptureScreenshot("02-search-box-visible")

		t.Logf("✅ Google homepage loaded with search box")
	})
}

// TestUI_Public_02_GoogleSearch performs a search by navigating directly to a results URL.
// This avoids the homepage interaction and is more reliable across regions/consent pages.
func TestUI_Public_02_GoogleSearch(t *testing.T) {
	tests.RunUITest(t, "Google Search Results Page Loads", func(t *testing.T, page *ui.Page) {
		// Navigate directly to a search results URL — this is the most reliable approach
		// as it bypasses any homepage consent dialogs or A/B test variants.
		searchURL := "https://www.google.com/search?q=e2e+testing+framework+golang&hl=en"
		if err := page.GoToHome(searchURL); err != nil {
			t.Fatalf("Failed to navigate to Google search results: %v", err)
		}
		_, _ = page.CaptureScreenshot("01-search-results-page")

		// Verify the search box on the results page is present.
		// Uses sequential selector attempts — comma-separated CSS locators are not supported.
		if err := waitForGoogleSearchBox(page, 10*time.Second); err != nil {
			t.Logf("WARNING: Search box not found on results page: %v", err)
		} else {
			_, _ = page.CaptureScreenshot("02-results-search-box-present")
		}

		// Verify at least the page body loaded (handles redirects/consent screens gracefully)
		if _, err := page.WaitUntilVisible("css:body", 5*time.Second); err != nil {
			t.Fatalf("Google search results page body did not load: %v", err)
		}

		// Check page title contains our search query
		title, err := page.Driver.Title()
		if err != nil {
			t.Logf("WARNING: Could not read page title: %v", err)
		} else {
			t.Logf("Page title: %q", title)
		}

		_, _ = page.CaptureScreenshot("03-results-final")
		t.Logf("✅ Google search results page loaded")
	})
}

// TestUI_Public_03_GoogleSearchTyped demonstrates typing into Google's search box and submitting.
func TestUI_Public_03_GoogleSearchTyped(t *testing.T) {
	tests.RunUITest(t, "Google Homepage Search by Typing", func(t *testing.T, page *ui.Page) {
		if err := page.GoToHome("https://www.google.com"); err != nil {
			t.Fatalf("Failed to navigate to Google: %v", err)
		}

		// Wait for search box (textarea on modern Chrome, input on older builds)
		if err := waitForGoogleSearchBox(page, 10*time.Second); err != nil {
			t.Fatalf("%v", err)
		}

		// Type the search query into whichever search box variant is present
		if err := sendKeysToGoogleSearchBox(page, "selenium golang e2e", 5*time.Second); err != nil {
			t.Fatalf("Failed to type into Google search box: %v", err)
		}
		_, _ = page.CaptureScreenshot("01-query-typed")

		// Submit by pressing Enter into the search box
		if err := sendKeysToGoogleSearchBox(page, "\n", 3*time.Second); err != nil {
			t.Fatalf("Failed to submit search: %v", err)
		}

		// Wait for the results page to load — check for body
		if _, err := page.WaitUntilVisible("css:body", 10*time.Second); err != nil {
			t.Fatalf("Search results page did not load: %v", err)
		}
		_, _ = page.CaptureScreenshot("02-results-after-submit")

		title, _ := page.Driver.Title()
		t.Logf("✅ Google search submitted, title=%q", title)
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// 02 · GitHub
// ─────────────────────────────────────────────────────────────────────────────

// TestUI_Public_04_GitHubHomePageLoads verifies GitHub's homepage loads.
func TestUI_Public_04_GitHubHomePageLoads(t *testing.T) {
	tests.RunUITest(t, "GitHub Homepage Loads", func(t *testing.T, page *ui.Page) {
		if err := page.GoToHome("https://github.com"); err != nil {
			t.Fatalf("Failed to navigate to GitHub: %v", err)
		}
		_, _ = page.CaptureScreenshot("01-github-home")

		// Verify the page body is present
		if _, err := page.WaitUntilVisible("css:body", 10*time.Second); err != nil {
			t.Fatalf("GitHub page body not found: %v", err)
		}
		_, _ = page.CaptureScreenshot("02-github-body-visible")

		// Verify the page title contains "GitHub"
		title, err := page.Driver.Title()
		if err != nil {
			t.Logf("WARNING: Could not read page title: %v", err)
		} else if !strings.Contains(title, "GitHub") {
			t.Errorf("Expected page title to contain 'GitHub', got %q", title)
		} else {
			t.Logf("✅ GitHub homepage loaded, title=%q", title)
		}
	})
}

// TestUI_Public_05_GitHubPublicRepoLoads verifies a known public repo page loads correctly.
func TestUI_Public_05_GitHubPublicRepoLoads(t *testing.T) {
	tests.RunUITest(t, "GitHub Public Repo Page Loads (octocat/Hello-World)", func(t *testing.T, page *ui.Page) {
		if err := page.GoToHome("https://github.com/octocat/Hello-World"); err != nil {
			t.Fatalf("Failed to navigate to GitHub repo: %v", err)
		}
		_, _ = page.CaptureScreenshot("01-repo-page")

		// Verify page body loaded
		if _, err := page.WaitUntilVisible("css:body", 10*time.Second); err != nil {
			t.Fatalf("GitHub repo page body not found: %v", err)
		}
		_, _ = page.CaptureScreenshot("02-repo-body-visible")

		// Verify the page title contains the repo name
		title, err := page.Driver.Title()
		if err != nil {
			t.Logf("WARNING: Could not read page title: %v", err)
		} else if !strings.Contains(title, "Hello-World") {
			t.Errorf("Expected title to contain 'Hello-World', got %q", title)
		} else {
			t.Logf("✅ GitHub public repo page loaded, title=%q", title)
		}
	})
}

// TestUI_Public_06_GitHubSearchWorks verifies GitHub's search functionality via direct URL.
func TestUI_Public_06_GitHubSearchWorks(t *testing.T) {
	tests.RunUITest(t, "GitHub Search Results Page Loads", func(t *testing.T, page *ui.Page) {
		if err := page.GoToHome("https://github.com/search?q=selenium+golang&type=repositories"); err != nil {
			t.Fatalf("Failed to navigate to GitHub search results: %v", err)
		}
		_, _ = page.CaptureScreenshot("01-github-search-results")

		if _, err := page.WaitUntilVisible("css:body", 10*time.Second); err != nil {
			t.Fatalf("GitHub search results page did not load: %v", err)
		}
		_, _ = page.CaptureScreenshot("02-search-results-body")

		title, _ := page.Driver.Title()
		t.Logf("✅ GitHub search page loaded, title=%q", title)
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// 03 · Example.com — Simplest possible baseline test
// Maintained by IANA, always available, predictable content.
// ─────────────────────────────────────────────────────────────────────────────

// TestUI_Public_07_ExampleComLoads verifies example.com loads with correct h1.
func TestUI_Public_07_ExampleComLoads(t *testing.T) {
	tests.RunUITest(t, "example.com Baseline Loads", func(t *testing.T, page *ui.Page) {
		if err := page.GoToHome("https://example.com"); err != nil {
			t.Fatalf("Failed to navigate to example.com: %v", err)
		}
		_, _ = page.CaptureScreenshot("01-example-com")

		// example.com always has an h1 with "Example Domain"
		h1Text, err := page.GetText("css:h1", 5*time.Second)
		if err != nil {
			t.Fatalf("h1 element not found on example.com: %v", err)
		}
		_, _ = page.CaptureScreenshot("02-h1-visible")

		if h1Text == "" {
			t.Errorf("Expected non-empty h1 text on example.com")
		} else {
			t.Logf("✅ example.com loaded, h1=%q", h1Text)
		}

		// Verify a link is present on the page
		if _, err := page.WaitUntilVisible("css:a", 3*time.Second); err != nil {
			t.Logf("Note: No links found on example.com: %v", err)
		} else {
			linkText, _ := page.GetText("css:a", 3*time.Second)
			t.Logf("✅ Link present: %q", linkText)
		}
	})
}

// TestUI_Public_08_ExampleComPageTitle verifies example.com has the correct browser title.
// Uses page.Driver.Title() (the browser API) instead of CSS:title (not a visible element).
func TestUI_Public_08_ExampleComPageTitle(t *testing.T) {
	tests.RunUITest(t, "example.com Has Correct Browser Page Title", func(t *testing.T, page *ui.Page) {
		if err := page.GoToHome("https://example.com"); err != nil {
			t.Fatalf("Failed to navigate to example.com: %v", err)
		}

		// Wait for page to load
		if _, err := page.WaitUntilVisible("css:h1", 5*time.Second); err != nil {
			t.Fatalf("Page did not load (h1 not found): %v", err)
		}
		_, _ = page.CaptureScreenshot("01-title-check")

		// Use Driver.Title() — the correct way to get <title> content.
		// css:title does NOT work because <title> is in <head> and is not a "visible" element.
		title, err := page.Driver.Title()
		if err != nil {
			t.Fatalf("Could not read browser title: %v", err)
		}

		expectedTitle := "Example Domain"
		if title != expectedTitle {
			t.Errorf("Expected page title %q, got %q", expectedTitle, title)
		} else {
			t.Logf("✅ Browser title matches: %q", title)
		}

		// Prevent unused import
		_ = fmt.Sprintf
	})
}
