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
