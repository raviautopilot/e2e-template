package actions

import (
	"fmt"
	"time"

	"e2e-template/pkg/config"
	"e2e-template/pkg/ui/pages"
)

// MemberActionsInterface allows MemberPersona to execute member-specific actions.
type MemberActionsInterface interface {
	PublicActionsInterface
	GetMemberPersona() *MemberPersona
}

// GetMemberPersona implements MemberActionsInterface for MemberPersona.
func (p *MemberPersona) GetMemberPersona() *MemberPersona {
	return p
}

// ─────────────────────────────────────────────────────────────────────────────
// TEMPLATE: Add your member action functions below.
//
// Each action function should:
//   1. Accept a MemberActionsInterface + *config.Config + *Result
//   2. Append the action name to result.Actions
//   3. Check result.Failed() to skip if a previous step failed
//   4. Perform the action (click, type, verify)
//   5. Capture a screenshot for evidence
//   6. On failure: set result.Status = "failed", result.Error = err
//
// Example:
//   func MemberLogin(mai MemberActionsInterface, cfg *config.Config, r *Result) { ... }
//   func MemberViewProfile(mai MemberActionsInterface, cfg *config.Config, r *Result) { ... }
//   func MemberSubmitForm(mai MemberActionsInterface, cfg *config.Config, data FormData, r *Result) { ... }
// ─────────────────────────────────────────────────────────────────────────────

// MemberLogin opens the member login modal, enters credentials, and submits.
// TODO: Replace the testIDs with your application's member login flow.
func MemberLogin(mai MemberActionsInterface, cfg *config.Config, r *Result) {
	actionName := "Member Login"
	r.Actions = append(r.Actions, actionName)
	if r.Failed() {
		r.Advice = append(r.Advice, fmt.Sprintf("Skipped '%s' because a previous step failed", actionName))
		return
	}

	mp := mai.GetMemberPersona()
	homePage := pages.NewHomePage(mp.Page)
	loginPage := pages.NewLoginPage(mp.Page)

	// Step 1: Open member login
	if err := homePage.OpenMemberLogin(cfg.MemberLoginButtonTestID, mp.DefaultTimeout); err != nil {
		r.Status = "failed"
		r.Error = err
		r.Advice = append(r.Advice, "Verify the member login button testID exists on the page.")
		return
	}

	// Step 2: Fill credentials and submit
	if err := loginPage.FillAndSubmitLogin(
		cfg.MemberLoginUsernameInputTestID,
		cfg.MemberLoginPasswordInputTestID,
		cfg.MemberLoginSubmitButtonTestID,
		cfg.MemberCredentials.Username,
		cfg.MemberCredentials.Password,
		mp.DefaultTimeout,
	); err != nil {
		r.Status = "failed"
		r.Error = err
		r.Advice = append(r.Advice, "Verify member credentials in config.json and login form testIDs.")
		return
	}

	// Step 3: Wait for logged-in state
	// TODO: Replace with your app's logged-in indicator element.
	time.Sleep(2 * time.Second)
	if scr, scrErr := mp.Page.CaptureScreenshot("MemberLogin_Success"); scrErr == nil {
		r.Evidence = append(r.Evidence, scr)
	}
	r.Advice = append(r.Advice, "Member login completed successfully.")
}
