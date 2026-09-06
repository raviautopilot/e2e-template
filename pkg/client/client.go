package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"e2e-template/pkg/logger"
)

// HttpError defines the interface for HTTP errors returned by the client.
type HttpError interface {
	error
	StatusCode() int
	ResponseBody() string
	UnderlyingError() error
}

type httpErrorImpl struct {
	statusCode int
	respBody   string
	err        error
}

func (e *httpErrorImpl) Error() string {
	if e.err != nil {
		return fmt.Sprintf("HTTP Error: status=%d, err=%s", e.statusCode, e.err.Error())
	}
	return fmt.Sprintf("HTTP Error: status=%d", e.statusCode)
}

func (e *httpErrorImpl) StatusCode() int {
	return e.statusCode
}

func (e *httpErrorImpl) ResponseBody() string {
	return e.respBody
}

func (e *httpErrorImpl) UnderlyingError() error {
	return e.err
}

// Client represents the custom HTTP client.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	LogDir     string
	TestName   string
	mu         sync.Mutex
	LastError  error
}

// NewClient initializes a new client with the given base URL and request logging directory.
func NewClient(baseURL string, timeout time.Duration, logDir string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
		LogDir: logDir,
	}
}

// SetTestName sets the test name for meaningful request/response file naming.
func (c *Client) SetTestName(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.TestName = name
}

// SendHttpRequest executes the HTTP request, validates input/output pointers, performs auth, and logs details.
func (c *Client) SendHttpRequest(method string, path string, headers map[string]string, reqBodyPtr interface{}, respBodyPtr interface{}, auth Authenticator) (httpErr HttpError) {
	defer func() {
		if httpErr != nil {
			c.LastError = httpErr
		}
	}()
	// 1. Validate pointers
	if reqBodyPtr != nil {
		val := reflect.ValueOf(reqBodyPtr)
		if val.Kind() != reflect.Ptr {
			return &httpErrorImpl{statusCode: 0, err: errors.New("request body must be a pointer to a struct/value")}
		}
		if val.IsNil() {
			return &httpErrorImpl{statusCode: 0, err: errors.New("request body pointer must not be nil")}
		}
	}

	if respBodyPtr != nil {
		val := reflect.ValueOf(respBodyPtr)
		if val.Kind() != reflect.Ptr {
			return &httpErrorImpl{statusCode: 0, err: errors.New("response body must be a pointer to a struct/value")}
		}
		if val.IsNil() {
			return &httpErrorImpl{statusCode: 0, err: errors.New("response body pointer must not be nil")}
		}
	}

	// 2. Configure transport for mTLS if ClientCertAuth is provided
	if auth != nil {
		if certAuth, ok := auth.(*ClientCertAuth); ok {
			c.mu.Lock()
			if c.HTTPClient.Transport == nil {
				c.HTTPClient.Transport = &http.Transport{
					TLSClientConfig: &tls.Config{
						Certificates: []tls.Certificate{certAuth.Certificate},
					},
				}
			} else {
				if transport, ok := c.HTTPClient.Transport.(*http.Transport); ok {
					if transport.TLSClientConfig == nil {
						transport.TLSClientConfig = &tls.Config{
							Certificates: []tls.Certificate{certAuth.Certificate},
						}
					} else {
						transport.TLSClientConfig.Certificates = []tls.Certificate{certAuth.Certificate}
					}
				}
			}
			c.mu.Unlock()
		}
	}

	// 3. Marshal Request Body
	var bodyReader io.Reader
	var requestRawBody []byte
	if reqBodyPtr != nil {
		var err error
		requestRawBody, err = json.Marshal(reqBodyPtr)
		if err != nil {
			return &httpErrorImpl{statusCode: 0, err: fmt.Errorf("failed to marshal request body: %w", err)}
		}
		bodyReader = bytes.NewReader(requestRawBody)
	}

	// 4. Create request
	fullURL := c.BaseURL + path
	httpReq, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return &httpErrorImpl{statusCode: 0, err: fmt.Errorf("failed to create http request: %w", err)}
	}

	// 5. Apply Headers
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}

	// 6. Apply Authentication
	if auth != nil {
		if err := auth.Apply(httpReq); err != nil {
			return &httpErrorImpl{statusCode: 0, err: fmt.Errorf("failed to apply authentication: %w", err)}
		}
	}

	// 7. Execute Request
	logger.Debug("Sending API request: %s %s", method, fullURL)
	startTime := time.Now()
	httpResp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		c.logExchange(startTime, httpReq, requestRawBody, nil, nil, err)
		return &httpErrorImpl{statusCode: 0, err: fmt.Errorf("network execution failed: %w", err)}
	}
	defer httpResp.Body.Close()

	// 8. Read Response Body
	responseRawBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		c.logExchange(startTime, httpReq, requestRawBody, httpResp, nil, err)
		return &httpErrorImpl{statusCode: httpResp.StatusCode, err: fmt.Errorf("failed to read response body: %w", err)}
	}

	// Log the API exchange
	c.logExchange(startTime, httpReq, requestRawBody, httpResp, responseRawBody, nil)

	// 9. Unmarshal Response Body if successful and pointer provided
	if httpResp.StatusCode >= 200 && httpResp.StatusCode < 300 {
		if respBodyPtr != nil && len(responseRawBody) > 0 {
			if err := json.Unmarshal(responseRawBody, respBodyPtr); err != nil {
				return &httpErrorImpl{
					statusCode: httpResp.StatusCode,
					respBody:   string(responseRawBody),
					err:        fmt.Errorf("failed to unmarshal response body: %w", err),
				}
			}
		}
		return nil
	}

	// 10. Return HTTP error for non-2xx status
	return &httpErrorImpl{
		statusCode: httpResp.StatusCode,
		respBody:   string(responseRawBody),
		err:        fmt.Errorf("request failed with status %d", httpResp.StatusCode),
	}
}

// detectCallerTestName inspects the call stack to find the active Go test function name.
func detectCallerTestName() string {
	for skip := 1; skip < 30; skip++ {
		pc, _, _, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fnName := runtime.FuncForPC(pc).Name()
		// e.g. "e2e-template/tests/api/httpbin_test.TestAPI_HttpBin_01_GetEcho.func1"
		parts := strings.Split(fnName, "/")
		lastPart := parts[len(parts)-1]
		dotParts := strings.Split(lastPart, ".")
		for _, dotPart := range dotParts {
			if strings.HasPrefix(dotPart, "Test") && !strings.HasPrefix(dotPart, "TestContext") && !strings.HasPrefix(dotPart, "Testing") {
				return dotPart
			}
		}
	}
	return ""
}

// logExchange logs the request/response details to disk.
func (c *Client) logExchange(startTime time.Time, req *http.Request, reqBody []byte, resp *http.Response, respBody []byte, err error) {
	now := time.Now()
	dateDir := now.Format("2006-01-02")
	timePrefix := now.Format("15-04-05")

	// Append suffix to avoid file collision during parallel runs
	rand.Seed(time.Now().UnixNano())
	suffix := fmt.Sprintf("%06d", rand.Intn(1000000))
	baseDir := filepath.Join(c.LogDir, dateDir)

	if err := os.MkdirAll(baseDir, 0777); err != nil {
		logger.Error("Failed to create log directory: %s", err)
		return
	}

	testTag := c.TestName
	if testTag == "" {
		testTag = detectCallerTestName()
	}
	if testTag != "" {
		testTag = strings.ReplaceAll(testTag, "/", "_")
		testTag = strings.ReplaceAll(testTag, " ", "_")
	}

	var reqFileName, respFileName string
	if testTag != "" {
		reqFileName = fmt.Sprintf("%s-%s-%s-request.json", testTag, timePrefix, suffix)
		respFileName = fmt.Sprintf("%s-%s-%s-response.json", testTag, timePrefix, suffix)
	} else {
		reqFileName = fmt.Sprintf("%s-%s-request.json", timePrefix, suffix)
		respFileName = fmt.Sprintf("%s-%s-response.json", timePrefix, suffix)
	}

	// Prepare Request JSON Log
	reqHeadersMap := make(map[string][]string)
	for k, v := range req.Header {
		reqHeadersMap[k] = v
	}

	var reqBodyJSON interface{}
	if len(reqBody) > 0 {
		_ = json.Unmarshal(reqBody, &reqBodyJSON)
	}

	reqLog := map[string]interface{}{
		"timestamp": startTime.Format(time.RFC3339Nano),
		"test_name": testTag,
		"method":    req.Method,
		"url":       req.URL.String(),
		"headers":   reqHeadersMap,
		"body":      reqBodyJSON,
	}

	reqFilePath := filepath.Join(baseDir, reqFileName)
	reqFile, fileErr := os.Create(reqFilePath)
	if fileErr == nil {
		encoder := json.NewEncoder(reqFile)
		encoder.SetIndent("", "  ")
		_ = encoder.Encode(reqLog)
		reqFile.Close()
	}

	// Prepare Response JSON Log
	var respLog map[string]interface{}
	if resp != nil {
		respHeadersMap := make(map[string][]string)
		for k, v := range resp.Header {
			respHeadersMap[k] = v
		}

		var respBodyJSON interface{}
		if len(respBody) > 0 {
			_ = json.Unmarshal(respBody, &respBodyJSON)
		}

		respLog = map[string]interface{}{
			"timestamp":   now.Format(time.RFC3339Nano),
			"test_name":   testTag,
			"status_code": resp.StatusCode,
			"headers":     respHeadersMap,
			"body":        respBodyJSON,
			"latency_ms":  now.Sub(startTime).Milliseconds(),
		}
	} else {
		errMsg := "No response (network error)"
		if err != nil {
			errMsg = err.Error()
		}
		respLog = map[string]interface{}{
			"timestamp": now.Format(time.RFC3339Nano),
			"test_name": testTag,
			"error":     errMsg,
		}
	}

	respFilePath := filepath.Join(baseDir, respFileName)
	respFile, fileErr := os.Create(respFilePath)
	if fileErr == nil {
		encoder := json.NewEncoder(respFile)
		encoder.SetIndent("", "  ")
		_ = encoder.Encode(respLog)
		respFile.Close()
	}
}
