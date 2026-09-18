package client

import (
	"encoding/base64"
	"net/http/httptest"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/api", nil)
	auth := &BasicAuth{
		Username: "admin",
		Password: "secretPassword123",
	}

	if err := auth.Apply(req); err != nil {
		t.Fatalf("unexpected error applying BasicAuth: %v", err)
	}

	expectedVal := "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:secretPassword123"))
	if got := req.Header.Get("Authorization"); got != expectedVal {
		t.Errorf("BasicAuth Authorization header = %q, want %q", got, expectedVal)
	}
}

func TestBearerTokenAuth(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/api", nil)
	auth := &BearerTokenAuth{
		Token: "my-jwt-token-xyz",
	}

	if err := auth.Apply(req); err != nil {
		t.Fatalf("unexpected error applying BearerTokenAuth: %v", err)
	}

	expectedVal := "Bearer my-jwt-token-xyz"
	if got := req.Header.Get("Authorization"); got != expectedVal {
		t.Errorf("BearerTokenAuth Authorization header = %q, want %q", got, expectedVal)
	}
}

func TestVaultTokenAuth(t *testing.T) {
	t.Run("without namespace", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://vault.local:8200/v1/secret/data/creds", nil)
		auth := &VaultTokenAuth{
			Token: "s.abcdef1234567890",
		}

		if err := auth.Apply(req); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := req.Header.Get("X-Vault-Token"); got != "s.abcdef1234567890" {
			t.Errorf("X-Vault-Token header = %q, want s.abcdef1234567890", got)
		}
		if got := req.Header.Get("X-Vault-Namespace"); got != "" {
			t.Errorf("unexpected X-Vault-Namespace header: %q", got)
		}
	})

	t.Run("with namespace", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://vault.local:8200/v1/secret/data/creds", nil)
		auth := &VaultTokenAuth{
			Token:     "hvs.enterpriseToken999",
			Namespace: "admin/engineering",
		}

		if err := auth.Apply(req); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := req.Header.Get("X-Vault-Token"); got != "hvs.enterpriseToken999" {
			t.Errorf("X-Vault-Token header = %q, want hvs.enterpriseToken999", got)
		}
		if got := req.Header.Get("X-Vault-Namespace"); got != "admin/engineering" {
			t.Errorf("X-Vault-Namespace header = %q, want admin/engineering", got)
		}
	})
}

func TestAPIKeyAuth(t *testing.T) {
	t.Run("header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/api", nil)
		auth := &APIKeyAuth{
			Key:   "X-API-Key",
			Value: "my-secret-api-key",
			In:    "header",
		}

		if err := auth.Apply(req); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := req.Header.Get("X-API-Key"); got != "my-secret-api-key" {
			t.Errorf("X-API-Key header = %q, want my-secret-api-key", got)
		}
	})

	t.Run("query", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/api?filter=active", nil)
		auth := &APIKeyAuth{
			Key:   "apiKey",
			Value: "query123",
			In:    "query",
		}

		if err := auth.Apply(req); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := req.URL.Query().Get("apiKey"); got != "query123" {
			t.Errorf("query parameter apiKey = %q, want query123", got)
		}
		if got := req.URL.Query().Get("filter"); got != "active" {
			t.Errorf("existing query parameter filter = %q, want active", got)
		}
	})
}

func TestCustomHeaderAuth(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/api", nil)
	auth := &CustomHeaderAuth{
		Headers: map[string]string{
			"X-Session-ID": "session_889900",
			"X-Tenant-ID":  "tenant_us_east",
		},
	}

	if err := auth.Apply(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := req.Header.Get("X-Session-ID"); got != "session_889900" {
		t.Errorf("X-Session-ID = %q, want session_889900", got)
	}
	if got := req.Header.Get("X-Tenant-ID"); got != "tenant_us_east" {
		t.Errorf("X-Tenant-ID = %q, want tenant_us_east", got)
	}
}

func TestMultiAuth(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/api", nil)
	auth := &MultiAuth{
		Authenticators: []Authenticator{
			&BearerTokenAuth{Token: "jwt-token-123"},
			&VaultTokenAuth{Token: "s.vault-456", Namespace: "dev"},
		},
	}

	if err := auth.Apply(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := req.Header.Get("Authorization"); got != "Bearer jwt-token-123" {
		t.Errorf("Authorization = %q, want Bearer jwt-token-123", got)
	}
	if got := req.Header.Get("X-Vault-Token"); got != "s.vault-456" {
		t.Errorf("X-Vault-Token = %q, want s.vault-456", got)
	}
	if got := req.Header.Get("X-Vault-Namespace"); got != "dev" {
		t.Errorf("X-Vault-Namespace = %q, want dev", got)
	}
}
