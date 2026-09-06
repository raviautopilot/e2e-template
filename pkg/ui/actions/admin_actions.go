package actions

import (
	"fmt"
	"time"

	"e2e-template/pkg/config"
	"e2e-template/pkg/ui"
	"e2e-template/pkg/ui/pages"
)

// AdminActionsInterface allows AdminPersona to execute admin-specific actions.
type AdminActionsInterface interface {
	PublicActionsInterface
	GetAdminPersona() *AdminPersona
}

// GetAdminPersona implements AdminActionsInterface for AdminPersona.
func (p *AdminPersona) GetAdminPersona() *AdminPersona {
	return p
}

// captureScreenshot captures a PNG screenshot with a small render delay and appends it to evidence.
func captureScreenshot(ap *AdminPersona, r *Result, name string) {
	if ap == nil || ap.Page == nil {
		return
	}
	time.Sleep(500 * time.Millisecond)
	if scr, scrErr := ap.Page.CaptureScreenshot(name); scrErr == nil {
		r.Evidence = append(r.Evidence, scr)
	}
}

// captureNetworkEvidence retrieves intercepted network errors and appends formatted logs to Advice.
func captureNetworkEvidence(ap *AdminPersona, r *Result) {
	if ap == nil || ap.Page == nil || ap.Page.Driver == nil {
		return
	}
	if raw := ap.Page.RetrieveNetworkErrors(); raw != "" {
		if formatted := ui.FormatNetworkErrors(raw); formatted != "" {
			r.Advice = append(r.Advice, formatted)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// TEMPLATE: Add your admin action functions below.
//
// Each action function should:
//   1. Accept an AdminActionsInterface + *config.Config + *Result
//   2. Append the action name to result.Actions
//   3. Check result.Failed() to skip if a previous step failed
//   4. Perform the action (click, type, verify)
//   5. Capture a screenshot for evidence
//   6. On failure: set result.Status = "failed", result.Error = err
//
// Example:
//   func LoginAsAdmin(aai AdminActionsInterface, cfg *config.Config, r *Result) { ... }
//   func AddNewEntity(aai AdminActionsInterface, cfg *config.Config, r *Result) { ... }
//   func DeleteEntity(aai AdminActionsInterface, cfg *config.Config, id string, r *Result) { ... }
// ─────────────────────────────────────────────────────────────────────────────

// LoginAsAdmin opens Admin login, populates credentials, submits, and verifies login.
// TODO: Replace the page objects and testIDs with your application's admin login flow.
func LoginAsAdmin(aai AdminActionsInterface, cfg *config.Config, r *Result) {
	actionName := "Login as Admin"
	r.Actions = append(r.Actions, actionName)
	if r.Failed() {
		r.Advice = append(r.Advice, fmt.Sprintf("Skipped '%s' because a previous step failed", actionName))
		return
	}

	ap := aai.GetAdminPersona()
	if err := ensurePage(ap.PublicPersona); err != nil {
		r.Status = "failed"
		r.Error = err
		return
	}

	homePage := pages.NewHomePage(ap.Page)
	loginPage := pages.NewLoginPage(ap.Page)

	// Step 1: Open admin login
	if err := homePage.OpenAdminLogin(cfg.AdminLoginButtonTestID, ap.DefaultTimeout); err != nil {
		r.Status = "failed"
		r.Error = err
		r.Advice = append(r.Advice, "Verify the admin login button testID exists on the page.")
		captureScreenshot(ap, r, "AdminLogin_OpenFailure")
		return
	}

	// Step 2: Fill credentials and submit
	if err := loginPage.FillAndSubmitLogin(
		cfg.AdminLoginUsernameInputTestID,
		cfg.AdminLoginPasswordInputTestID,
		cfg.AdminLoginSubmitButtonTestID,
		cfg.AdminCredentials.Username,
		cfg.AdminCredentials.Password,
		ap.DefaultTimeout,
	); err != nil {
		r.Status = "failed"
		r.Error = err
		r.Advice = append(r.Advice, "Verify admin credentials in config.json and login form testIDs.")
		captureScreenshot(ap, r, "AdminLogin_SubmitFailure")
		return
	}

	// Step 3: Wait for dashboard to load
	// TODO: Replace with your dashboard's indicator element.
	time.Sleep(2 * time.Second)
	captureScreenshot(ap, r, "AdminLogin_Success")
	r.Advice = append(r.Advice, "Admin login completed successfully.")
}
