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

# Run only public API tests
./run-api-tests.sh -run TestAPI_Public

# Run a specific test
./run-api-tests.sh -run TestAPI_Public_01_HttpBin
```

### Option B: Using Standard `go test`
```bash
# Run all API tests with verbose output
go test -v ./tests/api/...

# Run only public tests
go test -v ./tests/api/... -run TestAPI_Public -count=1
```

### Option C: Using Precompiled Binary (`api.test`)
```bash
./api.test -test.v -test.run TestAPI_Public
```

---

## 4. Test Catalog

### 🌐 Public API Tests (`tests/api/public_api_test.go`)
Real, working tests that execute against public internet APIs without requiring a local backend:

| Test Name | Target API | Endpoints Tested |
|---|---|---|
| `TestAPI_Public_01_HttpBin` | [httpbin.org](https://httpbin.org) | GET echo, query params, status codes (200/404/500), POST JSON, delay/timeout, Bearer auth, Basic auth |
| `TestAPI_Public_02_GitHub` | [api.github.com](https://api.github.com) | GET user, 404 nonexistent user, GET repo, 404 nonexistent repo, rate limit status |
| `TestAPI_Public_03_JSONPlaceholder` | [jsonplaceholder.typicode.com](https://jsonplaceholder.typicode.com) | GET post list (100 posts), GET single post, 404 nonexistent post, POST create (201 Created), GET users |

### 🛠️ Example Skeleton Tests (`tests/api/example_api_test.go`)
Template tests for your custom application backend (configured via `baseUrl` in `config.json`):

| Test Name | Purpose |
|---|---|
| `TestAPI_01_HealthCheck` | Verifies health / ping endpoint of target backend |
| `TestAPI_02_PublicEndpoints` | Verifies public / unauthenticated endpoints |

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
