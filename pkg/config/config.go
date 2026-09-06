package config

import (
	"encoding/json"
	"os"
	"strconv"
)

// Credentials holds a username/password pair.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Config holds the configuration values for the testing framework.
type Config struct {
	// Core settings
	BaseURL     string `json:"baseUrl"`
	UiURL       string `json:"uiUrl"`
	SeleniumURL string `json:"seleniumUrl"`
	Headless    bool   `json:"headless"`
	Timeout     int    `json:"timeout"`

	// Auth credentials
	AdminCredentials  Credentials `json:"adminCredentials"`
	MemberCredentials Credentials `json:"memberCredentials"`

	// UI test ID placeholders (1-2 per type: buttons, inputs)
	LoginButtonTestID        string `json:"loginButtonTestID"`
	SubmitButtonTestID       string `json:"submitButtonTestID"`
	UsernameInputTestID      string `json:"usernameInputTestID"`
	PasswordInputTestID      string `json:"passwordInputTestID"`

	// Backward-compatible aliases
	AdminLoginButtonTestID         string `json:"adminLoginButtonTestID,omitempty"`
	AdminLoginUsernameInputTestID string `json:"adminLoginUsernameInputTestID,omitempty"`
	AdminLoginPasswordInputTestID string `json:"adminLoginPasswordInputTestID,omitempty"`
	AdminLoginSubmitButtonTestID  string `json:"adminLoginSubmitButtonTestID,omitempty"`
	MemberLoginButtonTestID        string `json:"memberLoginButtonTestID,omitempty"`
	MemberLoginUsernameInputTestID string `json:"memberLoginUsernameInputTestID,omitempty"`
	MemberLoginPasswordInputTestID string `json:"memberLoginPasswordInputTestID,omitempty"`
	MemberLoginSubmitButtonTestID  string `json:"memberLoginSubmitButtonTestID,omitempty"`
	LogoutButtonTestID             string `json:"logoutButtonTestID,omitempty"`
}

// LoadConfig reads the configuration file from path and applies environment overrides.
func LoadConfig(path string) (*Config, error) {
	cfg := &Config{
		BaseURL:     "http://localhost:8080",
		UiURL:       "http://localhost:3000",
		SeleniumURL: "http://localhost:9515",
		Headless:    false,
		Timeout:     10,
	}

	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(cfg); err != nil {
			return nil, err
		}
	}

	// Fallback to generic login/input testIDs if persona-specific fields are unset
	if cfg.AdminLoginButtonTestID == "" {
		cfg.AdminLoginButtonTestID = cfg.LoginButtonTestID
	}
	if cfg.AdminLoginUsernameInputTestID == "" {
		cfg.AdminLoginUsernameInputTestID = cfg.UsernameInputTestID
	}
	if cfg.AdminLoginPasswordInputTestID == "" {
		cfg.AdminLoginPasswordInputTestID = cfg.PasswordInputTestID
	}
	if cfg.AdminLoginSubmitButtonTestID == "" {
		cfg.AdminLoginSubmitButtonTestID = cfg.SubmitButtonTestID
	}

	if cfg.MemberLoginButtonTestID == "" {
		cfg.MemberLoginButtonTestID = cfg.AdminLoginButtonTestID
	}
	if cfg.MemberLoginUsernameInputTestID == "" {
		cfg.MemberLoginUsernameInputTestID = cfg.AdminLoginUsernameInputTestID
	}
	if cfg.MemberLoginPasswordInputTestID == "" {
		cfg.MemberLoginPasswordInputTestID = cfg.AdminLoginPasswordInputTestID
	}
	if cfg.MemberLoginSubmitButtonTestID == "" {
		cfg.MemberLoginSubmitButtonTestID = cfg.AdminLoginSubmitButtonTestID
	}

	// Environment overrides
	if val := os.Getenv("E2E_BASE_URL"); val != "" {
		cfg.BaseURL = val
	}
	if val := os.Getenv("E2E_UI_URL"); val != "" {
		cfg.UiURL = val
	}
	if val := os.Getenv("E2E_SELENIUM_URL"); val != "" {
		cfg.SeleniumURL = val
	}
	if val := os.Getenv("E2E_HEADLESS"); val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			cfg.Headless = boolVal
		}
	}
	if val := os.Getenv("E2E_TIMEOUT"); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			cfg.Timeout = intVal
		}
	}

	return cfg, nil
}
