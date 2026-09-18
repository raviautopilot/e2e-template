package myservice_test

import (
	"os"
	"testing"

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
