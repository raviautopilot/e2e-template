package client

import (
	"crypto/tls"
	"encoding/base64"
	"net/http"
)

// Authenticator dictates how to apply authentication to an http.Request or http.Client.
type Authenticator interface {
	Apply(req *http.Request) error
}

// BearerTokenAuth implements Bearer Token authentication.
type BearerTokenAuth struct {
	Token string
}

// Apply adds the Bearer Token to the Authorization header.
func (a *BearerTokenAuth) Apply(req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+a.Token)
	return nil
}

// BasicAuth implements HTTP Basic authentication.
type BasicAuth struct {
	Username string
	Password string
}

// Apply adds the Basic Auth credentials to the Authorization header.
func (a *BasicAuth) Apply(req *http.Request) error {
	auth := a.Username + ":" + a.Password
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	return nil
}

// APIKeyAuth implements custom API Key authentication.
type APIKeyAuth struct {
	Key   string
	Value string
	In    string // "header" or "query"
}

// Apply injects the API key in the specified header or query param.
func (a *APIKeyAuth) Apply(req *http.Request) error {
	if a.In == "query" {
		q := req.URL.Query()
		q.Add(a.Key, a.Value)
		req.URL.RawQuery = q.Encode()
	} else {
		req.Header.Set(a.Key, a.Value)
	}
	return nil
}

// SSHKeyAuth represents signing request headers via an SSH private key.
type SSHKeyAuth struct {
	KeyName   string
	Signature string
}

// Apply adds custom headers for the signed SSH request.
func (a *SSHKeyAuth) Apply(req *http.Request) error {
	req.Header.Set("X-SSH-Key-Name", a.KeyName)
	req.Header.Set("X-SSH-Signature", a.Signature)
	return nil
}

// ClientCertAuth implements transport layer client certificate (mTLS) configuration.
type ClientCertAuth struct {
	Certificate tls.Certificate
}

// Apply is a no-op on the request; the transport loader handles its TLS configuration.
func (a *ClientCertAuth) Apply(req *http.Request) error {
	return nil
}

// VaultTokenAuth implements HashiCorp Vault token authentication via X-Vault-Token header.
type VaultTokenAuth struct {
	Token     string
	Namespace string // optional HashiCorp Vault namespace (X-Vault-Namespace)
}

// Apply injects the X-Vault-Token and optional X-Vault-Namespace headers.
func (a *VaultTokenAuth) Apply(req *http.Request) error {
	req.Header.Set("X-Vault-Token", a.Token)
	if a.Namespace != "" {
		req.Header.Set("X-Vault-Namespace", a.Namespace)
	}
	return nil
}

// CustomHeaderAuth applies arbitrary headers for authentication.
type CustomHeaderAuth struct {
	Headers map[string]string
}

// Apply injects custom header key-value pairs into the request.
func (a *CustomHeaderAuth) Apply(req *http.Request) error {
	for k, v := range a.Headers {
		req.Header.Set(k, v)
	}
	return nil
}

// MultiAuth combines multiple authenticators into a single execution chain.
type MultiAuth struct {
	Authenticators []Authenticator
}

// Apply executes each authenticator in sequence.
func (m *MultiAuth) Apply(req *http.Request) error {
	for _, auth := range m.Authenticators {
		if auth != nil {
			if err := auth.Apply(req); err != nil {
				return err
			}
		}
	}
	return nil
}
