package example_test

import (
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"e2e-template/tests"
)

func TestMain(m *testing.M) {
	tests.SetupSuite()
	exitCode := m.Run()
	tests.TeardownSuite()
	os.Exit(exitCode)
}

// isUIServerRunning checks whether the configured uiUrl is reachable.
func isUIServerRunning() bool {
	if tests.GlobalConfig == nil || tests.GlobalConfig.UiURL == "" {
		return false
	}
	uiURL := tests.GlobalConfig.UiURL
	if strings.Contains(uiURL, "localhost") || strings.Contains(uiURL, "127.0.0.1") {
		c := &http.Client{Timeout: 2 * time.Second}
		resp, err := c.Get(uiURL)
		if err != nil {
			return false
		}
		resp.Body.Close()
	}
	return true
}
