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

// MemberFormData holds example input field data for creating a new entity.
type MemberFormData struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Status string `json:"status"`
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

	// Login form test IDs (data-testid attributes)
	AdminLoginButtonTestID         string `json:"adminLoginButtonTestID"`
	AdminLoginUsernameInputTestID string `json:"adminLoginUsernameInputTestID"`
	AdminLoginPasswordInputTestID string `json:"adminLoginPasswordInputTestID"`
	AdminLoginSubmitButtonTestID  string `json:"adminLoginSubmitButtonTestID"`
	MemberLoginButtonTestID        string `json:"memberLoginButtonTestID"`
	MemberLoginUsernameInputTestID string `json:"memberLoginUsernameInputTestID"`
	MemberLoginPasswordInputTestID string `json:"memberLoginPasswordInputTestID"`
	MemberLoginSubmitButtonTestID  string `json:"memberLoginSubmitButtonTestID"`
	LogoutButtonTestID             string `json:"logoutButtonTestID"`

	// Example: Create/Add form test IDs
	AdminAddMemberButtonTestID       string `json:"adminAddMemberButtonTestID"`
	AdminAddMemberNameInputTestID    string `json:"adminAddMemberNameInputTestID"`
	AdminAddMemberEmailInputTestID   string `json:"adminAddMemberEmailInputTestID"`
	AdminAddMemberSubmitButtonTestID string `json:"adminAddMemberSubmitButtonTestID"`
	MemberSearchInputTestID          string `json:"memberSearchInputTestID"`
	MemberDeleteButtonTestID         string `json:"memberDeleteButtonTestID"`
	MemberConfirmDeleteButtonTestID  string `json:"memberConfirmDeleteButtonTestID"`
	MemberEditButtonTestID           string `json:"memberEditButtonTestID"`
	MemberSaveEditButtonTestID       string `json:"memberSaveEditButtonTestID"`

	// Example: Bulk upload test IDs
	AdminBulkUploadButtonTestID       string   `json:"adminBulkUploadButtonTestID"`
	AdminBulkUploadFileInputTestID    string   `json:"adminBulkUploadFileInputTestID"`
	AdminBulkUploadSubmitButtonTestID string   `json:"adminBulkUploadSubmitButtonTestID"`
	BulkMemberEmails                  []string `json:"bulkMemberEmails"`
	BulkMemberMobiles                 []string `json:"bulkMemberMobiles"`

	// Example: New entity form data
	NewMemberEmail    string         `json:"newMemberEmail"`
	NewMemberFormData MemberFormData `json:"newMemberFormData"`
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
