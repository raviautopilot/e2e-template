package config_test

import (
	"path/filepath"
	"testing"

	"e2e-template/pkg/config"
)

func TestLoadConfig_PlaceholdersAndFallbacks(t *testing.T) {
	configPath := filepath.Join("..", "..", "config.json")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.LoginButtonTestID == "" {
		t.Errorf("Expected loginButtonTestID to be populated, got empty")
	}
	if cfg.SubmitButtonTestID == "" {
		t.Errorf("Expected submitButtonTestID to be populated, got empty")
	}
	if cfg.UsernameInputTestID == "" {
		t.Errorf("Expected usernameInputTestID to be populated, got empty")
	}
	if cfg.PasswordInputTestID == "" {
		t.Errorf("Expected passwordInputTestID to be populated, got empty")
	}

	// Verify backward-compatible fallbacks
	if cfg.AdminLoginButtonTestID != cfg.LoginButtonTestID {
		t.Errorf("Expected AdminLoginButtonTestID to fallback to %s, got %s", cfg.LoginButtonTestID, cfg.AdminLoginButtonTestID)
	}
	if cfg.AdminLoginUsernameInputTestID != cfg.UsernameInputTestID {
		t.Errorf("Expected AdminLoginUsernameInputTestID to fallback to %s, got %s", cfg.UsernameInputTestID, cfg.AdminLoginUsernameInputTestID)
	}
	if cfg.AdminLoginPasswordInputTestID != cfg.PasswordInputTestID {
		t.Errorf("Expected AdminLoginPasswordInputTestID to fallback to %s, got %s", cfg.PasswordInputTestID, cfg.AdminLoginPasswordInputTestID)
	}
	if cfg.AdminLoginSubmitButtonTestID != cfg.SubmitButtonTestID {
		t.Errorf("Expected AdminLoginSubmitButtonTestID to fallback to %s, got %s", cfg.SubmitButtonTestID, cfg.AdminLoginSubmitButtonTestID)
	}

	if cfg.ChromeDriverPath != "lib/chromedriver" {
		t.Errorf("Expected ChromeDriverPath to be 'lib/chromedriver', got '%s'", cfg.ChromeDriverPath)
	}
}

func TestLoadConfig_ChromeDriverPath_EnvOverride(t *testing.T) {
	configPath := filepath.Join("..", "..", "config.json")
	customPath := "/custom/bin/chromedriver"
	t.Setenv("E2E_CHROMEDRIVER_PATH", customPath)

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.ChromeDriverPath != customPath {
		t.Errorf("Expected ChromeDriverPath to be '%s', got '%s'", customPath, cfg.ChromeDriverPath)
	}
}
