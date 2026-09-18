package github_test

import (
	"os"
	"testing"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

const baseURL = "https://api.github.com"

var (
	apiClient *client.Client
	client2   *client.Client
)

func TestMain(m *testing.M) {
	tests.SetupSuite()
	apiClient, client2 = tests.NewServiceClients(baseURL)
	exitCode := m.Run()
	tests.TeardownSuite()
	os.Exit(exitCode)
}
