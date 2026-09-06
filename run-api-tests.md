# E2E API Test Suite Execution Guide

This document outlines how to compile, execute, configure, and expand the API automation tests in this framework.

---

## 1. Prerequisites & Environment Setup

Ensure that Go (1.20+) is installed and that the target API endpoint is reachable.

### Check Go Installation:
```bash
go version
```

### Precompile the API Test Binary (Optional, Fast CI/CD Execution):
Building a precompiled test binary avoids recompiling on every invocation:
```bash
go test -c ./tests/api -o api.test
```

---

## 2. Test Configuration

The API test suite reads its configuration from `config.json`. You can modify this file directly or override parameters via environment variables:

| `config.json` Field | Environment Variable | Default Value | Description |
|---|---|---|---|
| `baseUrl` | `E2E_BASE_URL` | `http://localhost:8080` | Target URL of the backend REST API |
| `timeout` | `E2E_TIMEOUT` | `10` | Request timeout in seconds |

---

## 3. How to Run the Tests

### Option A: Using the `run-api-tests.sh` Script (Recommended)
```bash
# Run all API tests
./run-api-tests.sh

# Run all API tests in a specific package
./run-api-tests.sh ./tests/api/httpbin/...
./run-api-tests.sh ./tests/api/github/...
./run-api-tests.sh ./tests/api/jsonplaceholder/...
./run-api-tests.sh ./tests/api/example/...

# Run a specific test
./run-api-tests.sh -run TestAPI_HttpBin_01_GetEcho
./run-api-tests.sh -run TestAPI_GitHub_01_User
./run-api-tests.sh -run TestAPI_JSONPlaceholder_01_Posts
```

### Option B: Using Standard `go test`
```bash
# Run all API tests with verbose output
go test -v ./tests/api/...

# Run a specific package
go test -v ./tests/api/httpbin/...
go test -v ./tests/api/github/...
go test -v ./tests/api/jsonplaceholder/...
go test -v ./tests/api/example/...
```

---

## 4. Test Catalog

Each test suite resides in its own package under `tests/api/`, with **one test per file**:

### 🌐 httpbin Package (`tests/api/httpbin/`)
| File | Test Function | What It Tests |
|---|---|---|
| `01_get_echo_test.go` | `TestAPI_HttpBin_01_GetEcho` | GET request echo and query parameter handling |
| `02_status_codes_test.go` | `TestAPI_HttpBin_02_StatusCodes` | Standard status codes (200 OK, 404 Not Found, 500 Error) |
| `03_post_json_test.go` | `TestAPI_HttpBin_03_PostJSON` | POST request with JSON payload echoed in response |
| `04_delay_test.go` | `TestAPI_HttpBin_04_Delay` | Delayed response handling without premature timeout |
| `05_auth_test.go` | `TestAPI_HttpBin_05_Auth` | Bearer Token and HTTP Basic authentication |

### 🌐 GitHub Package (`tests/api/github/`)
| File | Test Function | What It Tests |
|---|---|---|
| `01_user_test.go` | `TestAPI_GitHub_01_User` | User profile retrieval and 404 on nonexistent user |
| `02_repo_test.go` | `TestAPI_GitHub_02_Repo` | Public repo metadata retrieval and 404 on nonexistent repo |
| `03_rate_limit_test.go` | `TestAPI_GitHub_03_RateLimit` | GitHub API rate limit endpoint validation |

### 🌐 JSONPlaceholder Package (`tests/api/jsonplaceholder/`)
| File | Test Function | What It Tests |
|---|---|---|
| `01_posts_test.go` | `TestAPI_JSONPlaceholder_01_Posts` | Post list (100 items), single post, and 404 for missing post |
| `02_create_post_test.go` | `TestAPI_JSONPlaceholder_02_CreatePost` | Creating a post with JSON payload (201 Created) |
| `03_users_test.go` | `TestAPI_JSONPlaceholder_03_Users` | Users list retrieval and user field validation |

### 🛠️ Example Package (`tests/api/example/`)
Template tests for your custom application backend (configured via `baseUrl` in `config.json`):

| File | Test Function | Purpose |
|---|---|---|
| `01_health_check_test.go` | `TestAPI_Example_01_HealthCheck` | Verifies health / ping endpoint of target backend |
| `02_public_endpoints_test.go` | `TestAPI_Example_02_PublicEndpoints` | Verifies public / unauthenticated endpoints |

---

## 5. Writing New API Tests

Tests use the reusable action helpers from `pkg/api/actions/`:

```go
package api_test

import (
    "testing"
    "time"

    "e2e-template/pkg/api/actions"
    "e2e-template/pkg/client"
    "e2e-template/tests"
)

func TestMyAPI(t *testing.T) {
    tests.RunAPITest(t, "My Endpoint Test", func(t *testing.T, tc *tests.TestContext) {
        c := client.NewClient("https://api.example.com", 10*time.Second, tests.ExecutionLogDir)

        // GET and expect 200 OK
        var resp MyResponseType
        actions.GetAndExpectOK(tc, c, "/items", &resp)
        actions.AssertNotEmpty(tc, "Items", resp.Items)

        // POST and expect 201 Created
        body := map[string]string{"name": "test"}
        var created MyResponseType
        actions.PostAndExpectCreated(tc, c, "/items", body, &created)
        actions.AssertEquals(tc, "Name", created.Name, "test")
    })
}
```

---

## 6. Test Evidence & Reports

Every test execution automatically collects evidence in `evidence/run-<timestamp>/`:
- **API Request/Response logs**: `evidence/run-<timestamp>/requests/`
- **Interactive HTML Report**: `evidence/run-<timestamp>/reports/report.html`
- **Markdown Report**: `evidence/run-<timestamp>/reports/report.md`
