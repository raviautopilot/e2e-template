# Comprehensive API Authentication Guide

This guide provides an end-to-end reference for implementing, configuring, and testing all major API authentication strategies within the **E2E Testing Framework**.

---

## Table of Contents

1. [Architecture & The `Authenticator` Interface](#1-architecture--the-authenticator-interface)
2. [Authentication Strategies Matrix](#2-authentication-strategies-matrix)
3. [HTTP Basic Authentication (`BasicAuth`)](#3-http-basic-authentication-basicauth)
4. [Bearer Token & JWT Authentication (`BearerTokenAuth`)](#4-bearer-token--jwt-authentication-bearertokenauth)
5. [HashiCorp Vault Token Authentication (`VaultTokenAuth`)](#5-hashicorp-vault-token-authentication-vaulttokenauth)
6. [API Key Authentication (`APIKeyAuth`)](#6-api-key-authentication-apikeyauth)
7. [Custom Headers & Session / Cookie Auth (`CustomHeaderAuth`)](#7-custom-headers--session--cookie-auth-customheaderauth)
8. [Mutual TLS / Client Certificates (`ClientCertAuth`)](#8-mutual-tls--client-certificates-clientcertauth)
9. [Chained Multi-Authentication (`MultiAuth`)](#9-chained-multi-authentication-multiauth)
10. [Negative & RBAC/ABAC Security Boundary Testing](#10-negative--rbacabac-security-boundary-testing)
11. [Quick Reference Cheat Sheet](#11-quick-reference-cheat-sheet)

---

## 1. Architecture & The `Authenticator` Interface

In `e2e-template`, all authentication mechanisms adhere to the [`client.Authenticator`](file:///home/ubuntu/code/github/raviautopilot/templates/e2e-template/pkg/client/auth.go) interface:

```go
package client

import "net/http"

// Authenticator dictates how to apply authentication to an http.Request or http.Client.
type Authenticator interface {
    Apply(req *http.Request) error
}
```

When an HTTP request is executed through [`client.Client.SendHttpRequest`](file:///home/ubuntu/code/github/raviautopilot/templates/e2e-template/pkg/client/client.go) or any of the higher-level action helpers ([`actions.AuthenticatedGet`](file:///home/ubuntu/code/github/raviautopilot/templates/e2e-template/pkg/api/actions/get.go), [`actions.AuthenticatedPost`](file:///home/ubuntu/code/github/raviautopilot/templates/e2e-template/pkg/api/actions/post.go), [`actions.AuthenticatedPut`](file:///home/ubuntu/code/github/raviautopilot/templates/e2e-template/pkg/api/actions/put.go), [`actions.AuthenticatedPatch`](file:///home/ubuntu/code/github/raviautopilot/templates/e2e-template/pkg/api/actions/patch.go), [`actions.AuthenticatedDelete`](file:///home/ubuntu/code/github/raviautopilot/templates/e2e-template/pkg/api/actions/delete.go)), the framework:

1. Validates request/response pointers.
2. Marshals JSON payloads (if present).
3. Configures TLS/mTLS if [`ClientCertAuth`](file:///home/ubuntu/code/github/raviautopilot/templates/e2e-template/pkg/client/auth.go) is passed.
4. Invokes `auth.Apply(httpReq)` to inject required headers, query parameters, or signatures.
5. Logs request headers, sanitized parameters, response payloads, and latency to `evidence/run-<timestamp>/logs/`.

---

## 2. Authentication Strategies Matrix

| Strategy | Header / Format | Go Authenticator Type | Typical Use Case |
|---|---|---|---|
| **Basic Auth** | `Authorization: Basic <base64(user:pass)>` | `&client.BasicAuth{Username, Password}` | Legacy APIs, Admin panels, Internal services |
| **Bearer Token** | `Authorization: Bearer <token>` | `&client.BearerTokenAuth{Token}` | JWT, OAuth2, GitHub PAT, Auth0, Keycloak |
| **HashiCorp Vault** | `X-Vault-Token: <token>`<br>`X-Vault-Namespace: <ns>` | `&client.VaultTokenAuth{Token, Namespace}` | HashiCorp Vault HTTP API v1, secrets engines |
| **API Key (Header)** | `X-API-Key: <key>` / custom key | `&client.APIKeyAuth{Key, Value, In: "header"}` | Third-party APIs (Stripe, SendGrid, Datadog) |
| **API Key (Query)** | `?api_key=<key>` / custom query | `&client.APIKeyAuth{Key, Value, In: "query"}` | CDN endpoints, Google Maps API, Webhooks |
| **Custom Headers** | Arbitrary key-value pairs | `&client.CustomHeaderAuth{Headers}` | Multi-tenant (`X-Tenant-ID`), CSRF, Session IDs |
| **Mutual TLS (mTLS)** | TLS Client Certificate handshake | `&client.ClientCertAuth{Certificate}` | Banking, Open Banking APIs, Zero-Trust |
| **Multi-Auth** | Multiple authenticators chained | `&client.MultiAuth{Authenticators}` | Enterprise services needing token + tenant + vault |

---

## 3. HTTP Basic Authentication (`BasicAuth`)

### 3.1 Overview
HTTP Basic Auth encodes credentials as `base64(username + ":" + password)` and transmits them via the standard `Authorization` header:
```http
GET /api/v1/protected HTTP/1.1
Host: api.example.com
Authorization: Basic YWRtaW46c2VjcmV0UGFzc3dvcmQ=
```

### 3.2 Usage with `client.BasicAuth`
```go
package myservice_test

import (
    "testing"

    "e2e-template/pkg/api/actions"
    "e2e-template/pkg/client"
    "e2e-template/tests"
)

func TestAPI_BasicAuth_Example(t *testing.T) {
    tests.RunAPITestWithClients(
        t,
        "Protected Endpoint - HTTP Basic Auth",
        "Authenticates using username and password loaded from config.",
        "HTTP 200 OK with user profile payload",
        apiClient,
        client2,
        func(tc *tests.TestContext) {
            // 1. Initialize BasicAuth
            auth := &client.BasicAuth{
                Username: tests.GlobalConfig.AdminCredentials.Username,
                Password: tests.GlobalConfig.AdminCredentials.Password,
            }

            // 2. Send authenticated GET
            var profile UserProfileResponse
            actions.AuthenticatedGet(tc, apiClient, "/api/v1/profile", auth, &profile)

            // 3. Assertions
            if profile.Email != tests.GlobalConfig.AdminCredentials.Username {
                tc.Errorf("Expected email %s, got %s", tests.GlobalConfig.AdminCredentials.Username, profile.Email)
            }
        },
    )
}
```

### 3.3 Dynamic Config & Environment Variables
Credentials can be configured in `config.json` and overridden dynamically at runtime:
```json
{
  "adminCredentials": {
    "username": "admin@example.com",
    "password": "REPLACE_WITH_ACTUAL_PASSWORD"
  }
}
```
CLI override:
```bash
E2E_ADMIN_USERNAME=ci_user E2E_ADMIN_PASSWORD=ci_secret_pass ./run-api-tests.sh myservice
```

---

## 4. Bearer Token & JWT Authentication (`BearerTokenAuth`)

### 4.1 Overview
Bearer tokens (JWTs or OAuth2 access tokens) are sent via:
```http
GET /api/v1/orders HTTP/1.1
Host: api.example.com
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 4.2 Standard Flow: Login -> Capture Token -> Authenticate Calls

```go
package myservice_test

import (
    "testing"

    "e2e-template/pkg/api/actions"
    "e2e-template/pkg/client"
    "e2e-template/tests"
)

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    AccessToken string `json:"access_token"`
    TokenType   string `json:"token_type"`
}

type OrderResponse struct {
    ID     string  `json:"id"`
    Amount float64 `json:"amount"`
}

func TestAPI_BearerAuth_LoginAndFetch(t *testing.T) {
    tests.RunAPITestWithClients(
        t,
        "Bearer Auth - Login and Fetch Protected Resource",
        "Logs in with admin credentials, extracts JWT, and queries orders API.",
        "HTTP 200 OK with orders list",
        apiClient,
        client2,
        func(tc *tests.TestContext) {
            // Step 1: Login to get the JWT
            loginBody := LoginRequest{
                Username: tests.GlobalConfig.AdminCredentials.Username,
                Password: tests.GlobalConfig.AdminCredentials.Password,
            }
            var loginResp LoginResponse
            actions.PostAndExpectOK(tc, apiClient, "/api/v1/auth/login", &loginBody, &loginResp)

            tc.AssertNotEmpty(loginResp.AccessToken, "Access token must not be empty")

            // Step 2: Construct BearerTokenAuth
            bearerAuth := &client.BearerTokenAuth{
                Token: loginResp.AccessToken,
            }

            // Step 3: Access protected endpoint
            var orders []OrderResponse
            actions.AuthenticatedGet(tc, apiClient, "/api/v1/orders", bearerAuth, &orders)

            tc.Actual = "Successfully fetched orders using Bearer JWT"
        },
    )
}
```

### 4.3 Static Personal Access Token (PAT) / API Bearer Token
If your API uses a long-lived secret token (e.g. GitHub PAT or service account token):
```go
token := os.Getenv("GITHUB_TOKEN")
auth := &client.BearerTokenAuth{Token: token}
actions.AuthenticatedGet(tc, apiClient, "/user/repos", auth, &repos)
```

---

## 5. HashiCorp Vault Token Authentication (`VaultTokenAuth`)

### 5.1 Overview
HashiCorp Vault uses the `X-Vault-Token` header for authentication on its HTTP API (`/v1/...`). For Vault Enterprise, multi-tenancy is controlled via the `X-Vault-Namespace` header:
```http
GET /v1/secret/data/database/creds HTTP/1.1
Host: vault.example.com:8200
X-Vault-Token: s.2N9Y5xP6Z0c8t3V7l9Q4R2T1
X-Vault-Namespace: engineering/payments
```

### 5.2 Using `client.VaultTokenAuth`

The framework includes dedicated first-class support for HashiCorp Vault via [`client.VaultTokenAuth`](file:///home/ubuntu/code/github/raviautopilot/templates/e2e-template/pkg/client/auth.go):

```go
type VaultTokenAuth struct {
    Token     string // Required: Token value for X-Vault-Token
    Namespace string // Optional: Namespace for X-Vault-Namespace (Enterprise)
}
```

### 5.3 Complete Vault KV v2 Test Example

```go
package vault_test

import (
    "os"
    "testing"

    "e2e-template/pkg/api/actions"
    "e2e-template/pkg/client"
    "e2e-template/tests"
)

// Vault KV v2 Data Structures
type VaultKVWriteRequest struct {
    Data map[string]interface{} `json:"data"`
}

type VaultKVReadResponse struct {
    Data struct {
        Data     map[string]interface{} `json:"data"`
        Metadata map[string]interface{} `json:"metadata"`
    } `json:"data"`
}

func TestAPI_Vault_ReadWriteSecrets(t *testing.T) {
    tests.RunAPITestWithClients(
        t,
        "HashiCorp Vault - Read and Write KV v2 Secret",
        "Authenticates using X-Vault-Token and writes/reads secrets.",
        "HTTP 200 OK with verified secret payload",
        apiClient,
        client2,
        func(tc *tests.TestContext) {
            vaultToken := os.Getenv("VAULT_TOKEN")
            if vaultToken == "" {
                vaultToken = "root" // Default local dev token
            }

            // 1. Initialize Vault Authenticator
            vaultAuth := &client.VaultTokenAuth{
                Token:     vaultToken,
                Namespace: os.Getenv("VAULT_NAMESPACE"), // e.g. "admin/dev" or ""
            }

            // 2. Write a secret to KV v2: POST /v1/secret/data/my-app/config
            secretPath := "/v1/secret/data/my-app/config"
            writePayload := VaultKVWriteRequest{
                Data: map[string]interface{}{
                    "db_user": "postgres_admin",
                    "api_key": "super-secret-key-42",
                },
            }
            var writeResp map[string]interface{}
            actions.AuthenticatedPost(tc, apiClient, secretPath, vaultAuth, &writePayload, &writeResp)

            // 3. Read back the secret: GET /v1/secret/data/my-app/config
            var readResp VaultKVReadResponse
            actions.AuthenticatedGet(tc, apiClient, secretPath, vaultAuth, &readResp)

            // 4. Assertions
            tc.AssertEqual(readResp.Data.Data["db_user"], "postgres_admin", "db_user match")
            tc.AssertEqual(readResp.Data.Data["api_key"], "super-secret-key-42", "api_key match")
        },
    )
}
```

### 5.4 Vault AppRole Login Flow
When testing automated service logins via Vault AppRole:
```go
type AppRoleLoginRequest struct {
    RoleId   string `json:"role_id"`
    SecretId string `json:"secret_id"`
}

type AppRoleLoginResponse struct {
    Auth struct {
        ClientToken string   `json:"client_token"`
        Policies    []string `json:"policies"`
        LeaseDuration int    `json:"lease_duration"`
    } `json:"auth"`
}

// 1. Exchange RoleID + SecretID for a token
var loginResp AppRoleLoginResponse
actions.PostAndExpectOK(tc, apiClient, "/v1/auth/approle/login", &appRoleReq, &loginResp)

// 2. Use the client_token with VaultTokenAuth
vaultAuth := &client.VaultTokenAuth{Token: loginResp.Auth.ClientToken}
```

---

## 6. API Key Authentication (`APIKeyAuth`)

### 6.1 Header-Based API Key
Many services expect keys like `X-API-Key`, `X-Api-Token`, or `apikey`:
```http
GET /v1/data HTTP/1.1
X-API-Key: live_sk_99a8b7c6d5e4
```

```go
auth := &client.APIKeyAuth{
    Key:   "X-API-Key",
    Value: "live_sk_99a8b7c6d5e4",
    In:    "header", // Injects into http.Header
}
actions.AuthenticatedGet(tc, apiClient, "/v1/data", auth, &resp)
```

### 6.2 Query-Parameter-Based API Key
Public or webhook endpoints frequently require query parameters:
```http
GET /maps/api/geocode/json?address=Paris&key=AIzaSy... HTTP/1.1
```

```go
auth := &client.APIKeyAuth{
    Key:   "key",
    Value: "AIzaSyD-ExampleKey",
    In:    "query", // Appends to req.URL.RawQuery
}
actions.AuthenticatedGet(tc, apiClient, "/maps/api/geocode/json?address=Paris", auth, &resp)
```

---

## 7. Custom Headers & Session / Cookie Auth (`CustomHeaderAuth`)

### 7.1 Using `client.CustomHeaderAuth`
When your API expects multiple custom headers or proprietary tokens:

```go
auth := &client.CustomHeaderAuth{
    Headers: map[string]string{
        "X-Tenant-ID":      "org_889900",
        "X-Session-ID":     "sess_abc123xyz",
        "X-Correlation-ID": "corr-uuid-4",
    },
}
actions.AuthenticatedGet(tc, apiClient, "/api/v1/tenant/settings", auth, &settings)
```

### 7.2 Cookie-Based Authentication
For session-based or cookie-authenticated APIs:
```go
auth := &client.CustomHeaderAuth{
    Headers: map[string]string{
        "Cookie":       "sessionId=s%3A9Z8Y7X; Path=/; HttpOnly",
        "X-CSRF-Token": "csrf_token_secret_val",
    },
}
actions.AuthenticatedPost(tc, apiClient, "/api/v1/account/update", auth, &updateReq, &resp)
```

---

## 8. Mutual TLS / Client Certificates (`ClientCertAuth`)

### 8.1 Overview
In Zero-Trust and high-security enterprise environments (Open Banking, Payment Gateways), the client authenticates via an X.509 certificate during the TLS handshake.

### 8.2 Usage with `client.ClientCertAuth`

```go
import (
    "crypto/tls"
    "e2e-template/pkg/client"
)

func loadClientCertificate(certFile, keyFile string) (*client.ClientCertAuth, error) {
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        return nil, err
    }
    return &client.ClientCertAuth{Certificate: cert}, nil
}

// In test execution:
certAuth, err := loadClientCertificate("fixtures/certs/client.crt", "fixtures/certs/client.key")
if err != nil {
    tc.Fatalf("Failed to load mTLS client certificate: %v", err)
}

actions.AuthenticatedGet(tc, apiClient, "/secure/transfers", certAuth, &transfers)
```
*Note: `client.Client.SendHttpRequest` automatically injects the certificate into the underlying HTTP Transport's `tls.Config.Certificates`.*

---

## 9. Chained Multi-Authentication (`MultiAuth`)

### 9.1 Overview
Complex enterprise gateways often require combining multiple authentication mechanisms simultaneously. For example:
- A user **Bearer JWT** (for user identity)
- A **HashiCorp Vault Token** (to dynamically access backend secrets)
- An **X-Tenant-ID** header (for multi-tenant routing)

### 9.2 Usage with `client.MultiAuth`

```go
multiAuth := &client.MultiAuth{
    Authenticators: []client.Authenticator{
        &client.BearerTokenAuth{Token: userJwtToken},
        &client.VaultTokenAuth{Token: vaultToken, Namespace: "engineering"},
        &client.CustomHeaderAuth{
            Headers: map[string]string{
                "X-Tenant-ID": "tenant-42",
            },
        },
    },
}

actions.AuthenticatedGet(tc, apiClient, "/api/v2/secure-gateway/data", multiAuth, &response)
```

---

## 10. Negative & RBAC/ABAC Security Boundary Testing

A robust test suite must verify that the API properly rejects invalid, missing, or insufficient credentials.

### 10.1 Testing 401 Unauthorized (Invalid / Missing Auth)

```go
func TestAPI_Auth_Negative_401(t *testing.T) {
    tests.RunAPITestWithClients(
        t,
        "Auth Boundary - Invalid Token Rejection",
        "Ensures protected endpoints reject invalid or corrupted tokens.",
        "HTTP 401 Unauthorized",
        apiClient,
        client2,
        func(tc *tests.TestContext) {
            invalidAuth := &client.BearerTokenAuth{Token: "invalid-bogus-token"}

            // Call SendHttpRequest directly to verify status
            var resp map[string]interface{}
            err := apiClient.SendHttpRequest("GET", "/api/v1/orders", nil, nil, &resp, invalidAuth)

            if err == nil {
                tc.Fatalf("Expected HTTP 401 Unauthorized, but request succeeded!")
            }
            if err.StatusCode() != 401 {
                tc.Errorf("Expected status 401, got %d. Body: %s", err.StatusCode(), err.ResponseBody())
            } else {
                tc.Actual = "HTTP 401 Unauthorized verified"
            }
        },
    )
}
```

### 10.2 Testing 403 Forbidden with Dual-Client RBAC

Use the two pre-initialized clients (`apiClient` and `client2`) to verify that standard members cannot access administrative resources:

```go
func TestAPI_RBAC_Boundary_AdminVsMember(t *testing.T) {
    tests.RunAPITestWithClients(
        t,
        "RBAC Boundary - Member Cannot Access Admin Endpoint",
        "Uses apiClient as Admin and client2 as Member to verify 403 Forbidden enforcement.",
        "Admin gets 200 OK, Member gets 403 Forbidden",
        apiClient,
        client2,
        func(tc *tests.TestContext) {
            adminAuth := &client.BasicAuth{
                Username: tests.GlobalConfig.AdminCredentials.Username,
                Password: tests.GlobalConfig.AdminCredentials.Password,
            }
            memberAuth := &client.BasicAuth{
                Username: tests.GlobalConfig.MemberCredentials.Username,
                Password: tests.GlobalConfig.MemberCredentials.Password,
            }

            adminPath := "/api/v1/admin/audit-logs"

            // 1. Admin should succeed (200 OK)
            var logs []map[string]interface{}
            actions.AuthenticatedGet(tc, apiClient, adminPath, adminAuth, &logs)

            // 2. Member should fail with 403 Forbidden
            var dummy map[string]interface{}
            err := client2.SendHttpRequest("GET", adminPath, nil, nil, &dummy, memberAuth)

            tc.AssertEqual(err != nil && err.StatusCode() == 403, true, "Member must receive HTTP 403 Forbidden")
            tc.Actual = "Admin accessed logs (200 OK); Member was blocked (403 Forbidden)"
        },
    )
}
```

---

## 11. Quick Reference Cheat Sheet

| Authentication Method | Code Snippet |
|---|---|
| **HTTP Basic** | `auth := &client.BasicAuth{Username: "user", Password: "pwd"}` |
| **Bearer Token (JWT)** | `auth := &client.BearerTokenAuth{Token: "eyJhbGciOi..."}` |
| **HashiCorp Vault** | `auth := &client.VaultTokenAuth{Token: "s.abc123", Namespace: "dev"}` |
| **API Key (Header)** | `auth := &client.APIKeyAuth{Key: "X-API-Key", Value: "key123", In: "header"}` |
| **API Key (Query)** | `auth := &client.APIKeyAuth{Key: "token", Value: "val456", In: "query"}` |
| **Custom Headers** | `auth := &client.CustomHeaderAuth{Headers: map[string]string{"X-Token": "v"}}` |
| **Cookie / Session** | `auth := &client.CustomHeaderAuth{Headers: map[string]string{"Cookie": "sess=1"}}` |
| **Mutual TLS (mTLS)** | `auth := &client.ClientCertAuth{Certificate: cert}` |
| **Multi-Chained** | `auth := &client.MultiAuth{Authenticators: []client.Authenticator{auth1, auth2}}` |

### Action Helpers Quick Invocation

```go
// GET
actions.AuthenticatedGet(tc, apiClient, "/path", auth, &respStruct)

// POST
actions.AuthenticatedPost(tc, apiClient, "/path", auth, &reqStruct, &respStruct)

// PUT
actions.AuthenticatedPut(tc, apiClient, "/path", auth, &reqStruct, &respStruct)

// PATCH
actions.AuthenticatedPatch(tc, apiClient, "/path", auth, &reqStruct, &respStruct)

// DELETE
actions.AuthenticatedDelete(tc, apiClient, "/path", auth, &respStruct)
```
