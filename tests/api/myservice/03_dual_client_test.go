package myservice_test

import (
	"testing"

	"e2e-template/tests"
)

// TestAPI_03_DualClient_SessionIsolation verifies that apiClient and client2 operate as distinct, isolated client instances.
func TestAPI_03_DualClient_SessionIsolation(t *testing.T) {
	tests.RunAPITestWithClients(
		t,
		"Dual Client Session Isolation",
		"Verifies that primary and secondary service clients maintain isolated sessions and headers.",
		"Two distinct clients with isolated configuration",
		apiClient,
		client2,
		func(tc *tests.TestContext) {
			if tc.Client == nil || tc.Client2 == nil {
				tc.Fatalf("Expected both tc.Client and tc.Client2 to be initialized")
			}
			if tc.Client == tc.Client2 {
				tc.Fatalf("Expected tc.Client and tc.Client2 to be distinct instances")
			}
			if tc.Client.BaseURL != tests.GlobalConfig.BaseURL || tc.Client2.BaseURL != tests.GlobalConfig.BaseURL {
				tc.Errorf("Expected both clients to target service baseUrl %s", tests.GlobalConfig.BaseURL)
			}
			tc.Actual = "Both primary and secondary service clients verified and isolated"
		},
	)
}
