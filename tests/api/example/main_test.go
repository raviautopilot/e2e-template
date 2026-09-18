package example_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

var (
	apiClient *client.Client
	client2   *client.Client
)

func TestMain(m *testing.M) {
	tests.SetupSuite()
	apiClient, client2 = tests.NewServiceClients(tests.GlobalConfig.BaseURL)
	exitCode := m.Run()
	tests.TeardownSuite()
	os.Exit(exitCode)
}

// isServerRunning does a quick connectivity check against the configured baseUrl.
func isServerRunning() bool {
	if tests.GlobalConfig == nil || tests.GlobalConfig.BaseURL == "" {
		return false
	}
	httpClient := &http.Client{Timeout: 2 * time.Second}
	resp, err := httpClient.Get(tests.GlobalConfig.BaseURL)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return true
}
