package pages

import (
	"time"

	"e2e-template/pkg/ui"
)

// AdminDashboardPage handles actions inside the Admin Panel / Dashboard.
type AdminDashboardPage struct {
	*ui.Page
}

// NewAdminDashboardPage creates a new AdminDashboardPage instance.
func NewAdminDashboardPage(page *ui.Page) *AdminDashboardPage {
	return &AdminDashboardPage{Page: page}
}

// OpenAdminPanel clicks the "Admin Panel" navigation button in the navbar.
func (a *AdminDashboardPage) OpenAdminPanel(adminPanelBtnTestID string, timeout time.Duration) error {
	return a.ClickByTestID(adminPanelBtnTestID, timeout)
}

// ─────────────────────────────────────────────────────────────────────────────
// TEMPLATE: Add your admin dashboard page actions below.
//
// Page Objects encapsulate low-level interactions with a specific page/view.
// They should NOT contain assertions — that's the test's job.
//
// Example methods to add:
//   func (a *AdminDashboardPage) OpenAddEntityModal(testID string, t time.Duration) error { ... }
//   func (a *AdminDashboardPage) FillEntityForm(data FormData, t time.Duration) error { ... }
//   func (a *AdminDashboardPage) SearchEntity(query string, t time.Duration) error { ... }
//   func (a *AdminDashboardPage) DeleteEntity(entityID string, t time.Duration) error { ... }
// ─────────────────────────────────────────────────────────────────────────────
