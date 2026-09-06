package actions

import (
	"fmt"
	"strings"
	"time"

	"e2e-template/pkg/ui"
)

// ─────────────────────────────────────────────────────────────────────────────
// Result — Captures test journey execution details
// ─────────────────────────────────────────────────────────────────────────────

// Result captures test journey execution details, actions, evidence, and advice.
type Result struct {
	TestName string
	Actions  []string
	Evidence []string // File paths of captured screenshots
	Advice   []string // Suggestive messages or remediation advice
	Status   string   // "passed" or "failed"
	Error    error    // The last encountered error, if any
}

// NewResult initializes a test journey result.
func NewResult(testName string) *Result {
	return &Result{
		TestName: testName,
		Status:   "passed",
		Actions:  make([]string, 0),
		Evidence: make([]string, 0),
		Advice:   make([]string, 0),
	}
}

// Failed returns true if any action in the journey failed.
func (r *Result) Failed() bool {
	return r.Status == "failed"
}

// ─────────────────────────────────────────────────────────────────────────────
// Personas — Define who is performing the actions
// ─────────────────────────────────────────────────────────────────────────────

// PublicActionsInterface allows different personas to share public actions.
type PublicActionsInterface interface {
	GetPublicPersona() *PublicPersona
}

// PublicPersona defines the credentials and capabilities of a general visitor.
type PublicPersona struct {
	*ui.Page
	BaseURL        string
	DefaultTimeout time.Duration
}

// NewPublicPersona creates a new PublicPersona.
func NewPublicPersona(page *ui.Page, baseURL string, defaultTimeout time.Duration) *PublicPersona {
	return &PublicPersona{
		Page:           page,
		BaseURL:        baseURL,
		DefaultTimeout: defaultTimeout,
	}
}

// GetPublicPersona implements PublicActionsInterface for PublicPersona.
func (p *PublicPersona) GetPublicPersona() *PublicPersona {
	return p
}

// ensurePage verifies the wrapper and its internal fields are not nil.
func ensurePage(p *PublicPersona) error {
	if p == nil || p.Page == nil {
		return fmt.Errorf("PublicPersona or Page is not initialized")
	}
	if p.Page.Driver == nil {
		return fmt.Errorf("WebDriver is not initialized")
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Persona Wrappers — Inherit public capabilities for specialized roles
// ─────────────────────────────────────────────────────────────────────────────

// MemberPersona embeds PublicPersona to inherit all public capabilities.
type MemberPersona struct {
	*PublicPersona
}

// NewMemberPersona creates a new MemberPersona.
func NewMemberPersona(page *ui.Page, baseURL string, defaultTimeout time.Duration) *MemberPersona {
	return &MemberPersona{
		PublicPersona: NewPublicPersona(page, baseURL, defaultTimeout),
	}
}

// SubscriberPersona embeds MemberPersona.
type SubscriberPersona struct {
	*MemberPersona
}

// NewSubscriberPersona creates a new SubscriberPersona.
func NewSubscriberPersona(page *ui.Page, baseURL string, defaultTimeout time.Duration) *SubscriberPersona {
	return &SubscriberPersona{
		MemberPersona: NewMemberPersona(page, baseURL, defaultTimeout),
	}
}

// AdminPersona embeds PublicPersona.
type AdminPersona struct {
	*PublicPersona
}

// NewAdminPersona creates a new AdminPersona.
func NewAdminPersona(page *ui.Page, baseURL string, defaultTimeout time.Duration) *AdminPersona {
	return &AdminPersona{
		PublicPersona: NewPublicPersona(page, baseURL, defaultTimeout),
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Generic Public Actions — Reusable across any project
// ─────────────────────────────────────────────────────────────────────────────

// GoToHome navigates to the persona's BaseURL and records the outcome.
func GoToHome(pai PublicActionsInterface, r *Result) {
	actionName := "Navigate to Home Page"
	r.Actions = append(r.Actions, actionName)
	if r.Failed() {
		r.Advice = append(r.Advice, fmt.Sprintf("Skipped '%s' because a previous step failed", actionName))
		return
	}

	p := pai.GetPublicPersona()
	if err := ensurePage(p); err != nil {
		r.Status = "failed"
		r.Error = err
		r.Advice = append(r.Advice, "Prerequisite Check: Ensure the Selenium driver is started and passed correctly to the persona")
		return
	}

	if err := p.Page.GoToHome(p.BaseURL); err != nil {
		r.Status = "failed"
		r.Error = err
		r.Advice = append(r.Advice, fmt.Sprintf("Advice: Check if the application is reachable at %s.", p.BaseURL))
		time.Sleep(1 * time.Second)
		if scr, scrErr := p.Page.CaptureScreenshot(r.TestName + "_GoToHome_Failure"); scrErr == nil {
			r.Evidence = append(r.Evidence, scr)
		}
		return
	}

	time.Sleep(2 * time.Second)
	if scr, scrErr := p.Page.CaptureScreenshot(r.TestName + "_GoToHome_Success"); scrErr == nil {
		r.Evidence = append(r.Evidence, scr)
	}
	r.Advice = append(r.Advice, "Home page loaded successfully.")
}

// GoToPage clicks a navigation element by testID and records the outcome.
// This is a generic template action — use it for any navigation link/button.
func GoToPage(pai PublicActionsInterface, r *Result, testID string, pageName string) {
	actionName := fmt.Sprintf("Navigate to %s", pageName)
	r.Actions = append(r.Actions, actionName)
	if r.Failed() {
		r.Advice = append(r.Advice, fmt.Sprintf("Skipped '%s' because a previous step failed", actionName))
		return
	}

	p := pai.GetPublicPersona()
	if err := ensurePage(p); err != nil {
		r.Status = "failed"
		r.Error = err
		return
	}

	if err := p.Page.ClickByTestID(testID, p.DefaultTimeout); err != nil {
		r.Status = "failed"
		r.Error = err
		r.Advice = append(r.Advice, fmt.Sprintf("Verify that '%s' element exists on the page and is clickable.", testID))
		time.Sleep(1 * time.Second)
		if scr, scrErr := p.Page.CaptureScreenshot(r.TestName + "_" + pageName + "_Failure"); scrErr == nil {
			r.Evidence = append(r.Evidence, scr)
		}
		return
	}

	time.Sleep(2 * time.Second)
	if scr, scrErr := p.Page.CaptureScreenshot(r.TestName + "_" + pageName + "_Success"); scrErr == nil {
		r.Evidence = append(r.Evidence, scr)
	}
	r.Advice = append(r.Advice, fmt.Sprintf("%s loaded successfully.", pageName))
}

// VerifyPageTitle checks that the browser's page title contains the expected string.
func VerifyPageTitle(pai PublicActionsInterface, r *Result, expectedTitle string) {
	actionName := fmt.Sprintf("Verify Page Title Contains '%s'", expectedTitle)
	r.Actions = append(r.Actions, actionName)
	if r.Failed() {
		r.Advice = append(r.Advice, fmt.Sprintf("Skipped '%s' because a previous step failed", actionName))
		return
	}

	p := pai.GetPublicPersona()
	if err := ensurePage(p); err != nil {
		r.Status = "failed"
		r.Error = err
		return
	}

	title, err := p.Page.Driver.Title()
	if err != nil {
		r.Status = "failed"
		r.Error = fmt.Errorf("could not read page title: %w", err)
		return
	}

	if !strings.Contains(title, expectedTitle) {
		r.Status = "failed"
		r.Error = fmt.Errorf("expected title to contain %q, got %q", expectedTitle, title)
		r.Advice = append(r.Advice, fmt.Sprintf("Page title was %q but expected it to contain %q.", title, expectedTitle))
		return
	}

	r.Advice = append(r.Advice, fmt.Sprintf("Page title verified: %q", title))
}

// VerifyElementVisible waits for an element to be visible and records the outcome.
func VerifyElementVisible(pai PublicActionsInterface, r *Result, locator string, elementName string) {
	actionName := fmt.Sprintf("Verify '%s' is Visible", elementName)
	r.Actions = append(r.Actions, actionName)
	if r.Failed() {
		r.Advice = append(r.Advice, fmt.Sprintf("Skipped '%s' because a previous step failed", actionName))
		return
	}

	p := pai.GetPublicPersona()
	if err := ensurePage(p); err != nil {
		r.Status = "failed"
		r.Error = err
		return
	}

	if _, err := p.Page.WaitUntilVisible(locator, p.DefaultTimeout); err != nil {
		r.Status = "failed"
		r.Error = fmt.Errorf("%s not found: %w", elementName, err)
		r.Advice = append(r.Advice, fmt.Sprintf("Element '%s' (locator: %s) was not visible within timeout.", elementName, locator))
		if scr, scrErr := p.Page.CaptureScreenshot(r.TestName + "_" + elementName + "_NotVisible"); scrErr == nil {
			r.Evidence = append(r.Evidence, scr)
		}
		return
	}

	if scr, scrErr := p.Page.CaptureScreenshot(r.TestName + "_" + elementName + "_Visible"); scrErr == nil {
		r.Evidence = append(r.Evidence, scr)
	}
	r.Advice = append(r.Advice, fmt.Sprintf("%s is visible.", elementName))
}

// GetElementText reads the text of an element and records the outcome.
func GetElementText(pai PublicActionsInterface, r *Result, locator string, elementName string) string {
	actionName := fmt.Sprintf("Read Text of '%s'", elementName)
	r.Actions = append(r.Actions, actionName)
	if r.Failed() {
		return ""
	}

	p := pai.GetPublicPersona()
	if err := ensurePage(p); err != nil {
		r.Status = "failed"
		r.Error = err
		return ""
	}

	text, err := p.Page.GetText(locator, p.DefaultTimeout)
	if err != nil {
		r.Status = "failed"
		r.Error = fmt.Errorf("could not read text of %s: %w", elementName, err)
		return ""
	}

	r.Advice = append(r.Advice, fmt.Sprintf("%s text: %q", elementName, text))
	return text
}

// TypeIntoElement types text into an element identified by CSS locator.
func TypeIntoElement(pai PublicActionsInterface, r *Result, locator string, text string, elementName string) {
	actionName := fmt.Sprintf("Type into '%s'", elementName)
	r.Actions = append(r.Actions, actionName)
	if r.Failed() {
		r.Advice = append(r.Advice, fmt.Sprintf("Skipped '%s' because a previous step failed", actionName))
		return
	}

	p := pai.GetPublicPersona()
	if err := ensurePage(p); err != nil {
		r.Status = "failed"
		r.Error = err
		return
	}

	if err := p.Page.SendKeys(locator, text, p.DefaultTimeout); err != nil {
		r.Status = "failed"
		r.Error = fmt.Errorf("could not type into %s: %w", elementName, err)
		r.Advice = append(r.Advice, fmt.Sprintf("Verify '%s' is an input/textarea and is interactable.", elementName))
		return
	}

	if scr, scrErr := p.Page.CaptureScreenshot(r.TestName + "_" + elementName + "_Typed"); scrErr == nil {
		r.Evidence = append(r.Evidence, scr)
	}
}

// ClickElement clicks an element by CSS locator and records the outcome.
func ClickElement(pai PublicActionsInterface, r *Result, locator string, elementName string) {
	actionName := fmt.Sprintf("Click '%s'", elementName)
	r.Actions = append(r.Actions, actionName)
	if r.Failed() {
		r.Advice = append(r.Advice, fmt.Sprintf("Skipped '%s' because a previous step failed", actionName))
		return
	}

	p := pai.GetPublicPersona()
	if err := ensurePage(p); err != nil {
		r.Status = "failed"
		r.Error = err
		return
	}

	if err := p.Page.Click(locator, p.DefaultTimeout); err != nil {
		r.Status = "failed"
		r.Error = fmt.Errorf("could not click %s: %w", elementName, err)
		r.Advice = append(r.Advice, fmt.Sprintf("Verify '%s' (locator: %s) exists and is clickable.", elementName, locator))
		if scr, scrErr := p.Page.CaptureScreenshot(r.TestName + "_" + elementName + "_ClickFail"); scrErr == nil {
			r.Evidence = append(r.Evidence, scr)
		}
		return
	}

	time.Sleep(1 * time.Second)
	if scr, scrErr := p.Page.CaptureScreenshot(r.TestName + "_" + elementName + "_Clicked"); scrErr == nil {
		r.Evidence = append(r.Evidence, scr)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Google-Specific Demo Actions — Used by public_ui_test.go
// These show how to write real-world actions for a third-party site.
// ─────────────────────────────────────────────────────────────────────────────

// VerifyGoogleSearchBox waits for Google's search box (textarea or input) to appear.
func VerifyGoogleSearchBox(pai PublicActionsInterface, r *Result) {
	actionName := "Verify Google Search Box Visible"
	r.Actions = append(r.Actions, actionName)
	if r.Failed() {
		return
	}

	p := pai.GetPublicPersona()
	if err := ensurePage(p); err != nil {
		r.Status = "failed"
		r.Error = err
		return
	}

	// Google uses <textarea name="q"> on modern Chrome but <input name="q"> on older layouts
	if _, err := p.Page.WaitUntilVisible("css:textarea[name='q']", p.DefaultTimeout); err == nil {
		if scr, scrErr := p.Page.CaptureScreenshot(r.TestName + "_SearchBox_Visible"); scrErr == nil {
			r.Evidence = append(r.Evidence, scr)
		}
		return
	}
	if _, err := p.Page.WaitUntilVisible("css:input[name='q']", 2*time.Second); err == nil {
		if scr, scrErr := p.Page.CaptureScreenshot(r.TestName + "_SearchBox_Visible"); scrErr == nil {
			r.Evidence = append(r.Evidence, scr)
		}
		return
	}

	r.Status = "failed"
	r.Error = fmt.Errorf("google search box not found (tried textarea[name='q'] and input[name='q'])")
	r.Advice = append(r.Advice, "Google may have changed their search box markup. Check the page source.")
}

// TypeInGoogleSearchBox types text into Google's search box (tries textarea first, then input).
func TypeInGoogleSearchBox(pai PublicActionsInterface, r *Result, query string) {
	actionName := fmt.Sprintf("Type '%s' into Google Search Box", query)
	r.Actions = append(r.Actions, actionName)
	if r.Failed() {
		return
	}

	p := pai.GetPublicPersona()
	if err := ensurePage(p); err != nil {
		r.Status = "failed"
		r.Error = err
		return
	}

	if err := p.Page.SendKeys("css:textarea[name='q']", query, p.DefaultTimeout); err == nil {
		if scr, scrErr := p.Page.CaptureScreenshot(r.TestName + "_Query_Typed"); scrErr == nil {
			r.Evidence = append(r.Evidence, scr)
		}
		return
	}
	if err := p.Page.SendKeys("css:input[name='q']", query, 3*time.Second); err == nil {
		if scr, scrErr := p.Page.CaptureScreenshot(r.TestName + "_Query_Typed"); scrErr == nil {
			r.Evidence = append(r.Evidence, scr)
		}
		return
	}

	r.Status = "failed"
	r.Error = fmt.Errorf("could not type into Google search box")
}

// SubmitGoogleSearch submits the search by pressing Enter.
func SubmitGoogleSearch(pai PublicActionsInterface, r *Result) {
	actionName := "Submit Google Search"
	r.Actions = append(r.Actions, actionName)
	if r.Failed() {
		return
	}

	p := pai.GetPublicPersona()
	if err := ensurePage(p); err != nil {
		r.Status = "failed"
		r.Error = err
		return
	}

	// Submit by pressing Enter into the search box
	if err := p.Page.SendKeys("css:textarea[name='q']", "\n", 3*time.Second); err != nil {
		if err2 := p.Page.SendKeys("css:input[name='q']", "\n", 3*time.Second); err2 != nil {
			r.Status = "failed"
			r.Error = fmt.Errorf("could not submit Google search: %w", err2)
			return
		}
	}

	// Wait for results page to load
	if _, err := p.Page.WaitUntilVisible("css:body", 10*time.Second); err != nil {
		r.Status = "failed"
		r.Error = fmt.Errorf("search results page did not load: %w", err)
		return
	}

	time.Sleep(2 * time.Second)
	if scr, scrErr := p.Page.CaptureScreenshot(r.TestName + "_SearchResults"); scrErr == nil {
		r.Evidence = append(r.Evidence, scr)
	}
	r.Advice = append(r.Advice, "Google search submitted and results page loaded.")
}
