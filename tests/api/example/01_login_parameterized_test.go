package example_test

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// LoginTestCase defines a table-driven scenario for API authentication testing.
type LoginTestCase struct {
	Name           string
	Description    string
	Email          string
	Password       string
	WantStatusCode int
	ExpectSuccess  bool
}

// TestAPI_01_AdminLogin_Parameterized runs high-level parameterized login scenarios against the backend API.
// Default Target Host: dev.varam.srilakshmiretail.com
// Default Admin Credentials: sudhan.admin@gmail.com / testlaxmi
func TestAPI_01_AdminLogin_Parameterized(t *testing.T) {
	baseURL := "https://dev.varam.srilakshmiretail.com"
	if tests.GlobalConfig != nil && tests.GlobalConfig.BaseURL != "" && tests.GlobalConfig.BaseURL != "https://api.yourapp.com" {
		baseURL = tests.GlobalConfig.BaseURL
	}

	apiClient := client.NewClient(baseURL, 15*time.Second, tests.ExecutionLogDir)

	scenarios := []LoginTestCase{
		{
			Name:           "Valid Admin Credentials",
			Description:    "Submit valid admin email and password (sudhan.admin@gmail.com / testlaxmi)",
			Email:          "sudhan.admin@gmail.com",
			Password:       "testlaxmi",
			WantStatusCode: 200,
			ExpectSuccess:  true,
		},
		{
			Name:           "Invalid Password",
			Description:    "Submit valid admin email with incorrect password",
			Email:          "sudhan.admin@gmail.com",
			Password:       "wrongpassword123",
			WantStatusCode: 401,
			ExpectSuccess:  false,
		},
		{
			Name:           "Non-Existent User Email",
			Description:    "Submit non-existent user email with valid password structure",
			Email:          "nonexistent.admin.xyz@gmail.com",
			Password:       "testlaxmi",
			WantStatusCode: 401,
			ExpectSuccess:  false,
		},
		{
			Name:           "Empty Password Field",
			Description:    "Submit valid admin email with blank password",
			Email:          "sudhan.admin@gmail.com",
			Password:       "",
			WantStatusCode: 401,
			ExpectSuccess:  false,
		},
		{
			Name:           "Empty Email Field",
			Description:    "Submit blank email with valid admin password",
			Email:          "",
			Password:       "testlaxmi",
			WantStatusCode: 401,
			ExpectSuccess:  false,
		},
		{
			Name:           "Both Fields Empty",
			Description:    "Submit empty request body for authentication endpoint",
			Email:          "",
			Password:       "",
			WantStatusCode: 401,
			ExpectSuccess:  false,
		},
		{
			Name:           "Malformed Email Format",
			Description:    "Submit invalid email syntax (missing @ domain)",
			Email:          "sudhan.admin.invalidemail",
			Password:       "testlaxmi",
			WantStatusCode: 401,
			ExpectSuccess:  false,
		},
		{
			Name:           "SQL Injection in Email Field",
			Description:    "Submit malicious SQL injection payload in email field",
			Email:          "' OR '1'='1",
			Password:       "testlaxmi",
			WantStatusCode: 401,
			ExpectSuccess:  false,
		},
		{
			Name:           "XSS Script Tag Payload in Email",
			Description:    "Submit XSS payload script tag in email field",
			Email:          "<script>alert('xss')</script>",
			Password:       "testlaxmi",
			WantStatusCode: 401,
			ExpectSuccess:  false,
		},
	}

	for _, sc := range scenarios {
		sc := sc
		t.Run(sc.Name, func(t *testing.T) {
			expectedText := fmt.Sprintf("HTTP Status %d", sc.WantStatusCode)
			if sc.ExpectSuccess {
				expectedText = "HTTP 200 OK with valid Auth Token / User Session"
			}

			tests.RunAPITestWithDetails(
				t,
				fmt.Sprintf("Admin Login - %s", sc.Name),
				sc.Description,
				expectedText,
				func(tc *tests.TestContext) {
					testName := fmt.Sprintf("Admin Login - %s", sc.Name)
					apiClient.SetTestName(testName)
					tc.Client = apiClient
					req := LoginRequest{
						Email:    sc.Email,
						Password: sc.Password,
					}
					var resp LoginResponse

					if sc.ExpectSuccess {
						actions.PostAndExpectOK(tc, apiClient, "/api/v1/auth/login", &req, &resp)
					} else {
						actions.PostAndExpectStatus(tc, apiClient, "/api/v1/auth/login", &req, sc.WantStatusCode)
					}
				},
			)
		})
	}
}
