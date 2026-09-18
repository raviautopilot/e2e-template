package tests

import (
	"path/filepath"
	"testing"

	"e2e-template/pkg/config"
)

func TestResolveChromeDriverPath(t *testing.T) {
	moduleRoot, err := findModuleRoot()
	if err != nil {
		t.Fatalf("Failed to find module root: %v", err)
	}

	origConfig := GlobalConfig
	defer func() {
		GlobalConfig = origConfig
	}()

	// Case 1: Empty / not set -> returns "chromedriver"
	GlobalConfig = &config.Config{
		ChromeDriverPath: "",
	}
	if p := resolveChromeDriverPath(); p != "chromedriver" {
		t.Errorf("Expected 'chromedriver', got '%s'", p)
	}

	// Case 2: Explicit "chromedriver" -> returns "chromedriver"
	GlobalConfig.ChromeDriverPath = "chromedriver"
	if p := resolveChromeDriverPath(); p != "chromedriver" {
		t.Errorf("Expected 'chromedriver', got '%s'", p)
	}

	// Case 3: Relative path "lib/chromedriver" -> returns moduleRoot/lib/chromedriver
	GlobalConfig.ChromeDriverPath = "lib/chromedriver"
	expectedPath := filepath.Join(moduleRoot, "lib", "chromedriver")
	if p := resolveChromeDriverPath(); p != expectedPath {
		t.Errorf("Expected '%s', got '%s'", expectedPath, p)
	}

	// Case 4: Absolute path
	absPath := "/opt/custom/chromedriver"
	GlobalConfig.ChromeDriverPath = absPath
	if p := resolveChromeDriverPath(); p != absPath {
		t.Errorf("Expected '%s', got '%s'", absPath, p)
	}
}

func TestNewServiceClient_And_NewServiceClients(t *testing.T) {
	if GlobalConfig == nil {
		SetupSuite()
	}

	testURL := "https://api.testservice.local"
	c := NewServiceClient(testURL)
	if c == nil {
		t.Fatalf("Expected non-nil client from NewServiceClient")
	}
	if c.BaseURL != testURL {
		t.Errorf("Expected BaseURL %s, got %s", testURL, c.BaseURL)
	}

	c1, c2 := NewServiceClients(testURL)
	if c1 == nil || c2 == nil {
		t.Fatalf("Expected non-nil clients from NewServiceClients")
	}
	if c1 == c2 {
		t.Errorf("Expected c1 and c2 to be distinct instances, got identical pointer")
	}
	if c1.BaseURL != testURL || c2.BaseURL != testURL {
		t.Errorf("Expected BaseURL %s for both clients", testURL)
	}
}

func TestRunAPITestWithClients_DualClientContext(t *testing.T) {
	c1, c2 := NewServiceClients("https://api.testservice.local")

	RunAPITestWithClients(t, "Dual Client Test Demo", "Validates c1 and c2 are injected into TestContext", "Both clients present in context", c1, c2, func(tc *TestContext) {
		if tc.Client == nil {
			tc.Errorf("Expected tc.Client to be populated")
		}
		if tc.Client2 == nil {
			tc.Errorf("Expected tc.Client2 to be populated")
		}
		if tc.Client == tc.Client2 {
			tc.Errorf("Expected tc.Client and tc.Client2 to be distinct instances")
		}
	})
}

