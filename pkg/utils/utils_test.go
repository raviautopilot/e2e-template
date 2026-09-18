package utils_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"e2e-template/pkg/utils"
)

func TestGenerators(t *testing.T) {
	// Test GenerateRandomString
	s10 := utils.GenerateRandomString(10)
	if len(s10) != 10 {
		t.Errorf("Expected length 10, got %d", len(s10))
	}
	s0 := utils.GenerateRandomString(0)
	if s0 != "" {
		t.Errorf("Expected empty string for len 0, got %s", s0)
	}
	sNeg := utils.GenerateRandomString(-5)
	if sNeg != "" {
		t.Errorf("Expected empty string for len -5, got %s", sNeg)
	}

	// Uniqueness check
	s10b := utils.GenerateRandomString(10)
	if s10 == s10b {
		t.Errorf("Expected different random strings, got identical %s", s10)
	}

	// Test GenerateRandomStringwithTimeStamp
	tsStr := utils.GenerateRandomStringwithTimeStamp(8)
	parts := strings.Split(tsStr, "_")
	if len(parts) < 3 {
		t.Fatalf("Expected at least 3 parts (YYYYMMDD_HHMMSS_random), got %s", tsStr)
	}
	randomPart := parts[len(parts)-1]
	if len(randomPart) != 8 {
		t.Errorf("Expected random part length 8, got %d in %s", len(randomPart), tsStr)
	}

	// Test PascalCase alias
	tsAlias := utils.GenerateRandomStringWithTimeStamp(6)
	if !strings.Contains(tsAlias, "_") {
		t.Errorf("Expected timestamped string with _, got %s", tsAlias)
	}

	// Test Email generator
	email := utils.GenerateRandomEmail("myorg.test")
	if !strings.HasSuffix(email, "@myorg.test") {
		t.Errorf("Expected suffix @myorg.test, got %s", email)
	}

	// Test Int generator
	val := utils.GenerateRandomInt(5, 10)
	if val < 5 || val > 10 {
		t.Errorf("Expected integer between 5 and 10, got %d", val)
	}

	// Test Phone generator
	phone := utils.GenerateRandomPhone()
	if !strings.HasPrefix(phone, "+1") || len(phone) != 12 {
		t.Errorf("Expected +1 followed by 10 digits, got %s", phone)
	}
}

func TestPrettyPrint(t *testing.T) {
	type SampleData struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	data := SampleData{Name: "testing", Count: 42}

	// Test PrettyJSON
	out, err := utils.PrettyJSON(data)
	if err != nil {
		t.Fatalf("Unexpected error from PrettyJSON: %v", err)
	}
	if !strings.Contains(out, "  \"name\": \"testing\"") {
		t.Errorf("Expected indented JSON output, got:\n%s", out)
	}

	// Test MustPrettyJSON
	mustOut := utils.MustPrettyJSON(data)
	if mustOut != out {
		t.Errorf("MustPrettyJSON mismatch:\n%s\nvs\n%s", mustOut, out)
	}

	// Test PrettyJSONString
	raw := `{"status":"ok","items":[1,2,3]}`
	prettyRaw, err := utils.PrettyJSONString(raw)
	if err != nil {
		t.Fatalf("Unexpected error from PrettyJSONString: %v", err)
	}
	if !strings.Contains(prettyRaw, "  \"status\": \"ok\"") {
		t.Errorf("Expected indented raw JSON, got:\n%s", prettyRaw)
	}

	// Test FprintJSON
	var buf bytes.Buffer
	if err := utils.FprintJSON(&buf, data); err != nil {
		t.Fatalf("FprintJSON failed: %v", err)
	}
	if !strings.Contains(buf.String(), "testing") {
		t.Errorf("Expected buffer to contain formatted JSON, got: %s", buf.String())
	}
}

func TestProject(t *testing.T) {
	// Test FindProjectRoot
	root, err := utils.FindProjectRoot()
	if err != nil {
		t.Fatalf("FindProjectRoot failed: %v", err)
	}
	if !utils.FileExists(filepath.Join(root, "go.mod")) {
		t.Errorf("Expected go.mod to exist at project root %s", root)
	}

	// Test GetProjectRoot
	if r := utils.GetProjectRoot(); r != root {
		t.Errorf("GetProjectRoot mismatch: %s vs %s", r, root)
	}

	// Test GetModuleName
	moduleName, err := utils.GetModuleName()
	if err != nil {
		t.Fatalf("GetModuleName failed: %v", err)
	}
	if moduleName != "e2e-template" {
		t.Errorf("Expected module name 'e2e-template', got '%s'", moduleName)
	}

	// Test GetConfigPath
	cfgPath, err := utils.GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}
	if !strings.HasSuffix(cfgPath, "config.json") {
		t.Errorf("Expected path ending with config.json, got %s", cfgPath)
	}

	// Test GetEvidenceDir
	evDir, err := utils.GetEvidenceDir("2026-09-18")
	if err != nil {
		t.Fatalf("GetEvidenceDir failed: %v", err)
	}
	if !strings.HasSuffix(evDir, "run-2026-09-18") {
		t.Errorf("Expected evidence dir to end with run-2026-09-18, got %s", evDir)
	}

	// Test Git Branch
	branch := utils.GetGitBranch()
	if branch == "" {
		t.Logf("Git branch returned empty (may not be git repo or detached HEAD)")
	} else {
		t.Logf("Detected git branch: %s", branch)
	}

	// Test FileExists & DirExists
	if !utils.FileExists(filepath.Join(root, "go.mod")) {
		t.Errorf("Expected go.mod to be recognized as existing file")
	}
	if utils.FileExists(root) {
		t.Errorf("Project root directory should not be recognized as a regular file")
	}
	if !utils.DirExists(root) {
		t.Errorf("Expected project root to be recognized as existing directory")
	}
}

func TestScriptExecutor(t *testing.T) {
	// Create a temporary shell script for testing
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test_runner.sh")
	scriptContent := `#!/bin/bash
echo "Hello from script: $1"
echo "ENV_FOO is: $CUSTOM_FOO"
>&2 echo "Notice: logging to stderr"
echo '{"status":"completed","code":200}'
`
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("Failed to write test script: %v", err)
	}

	// Test ExecuteScriptWithEnv
	env := map[string]string{
		"CUSTOM_FOO": "bar_value_123",
	}
	res, err := utils.ExecuteScriptWithEnv(scriptPath, env, "arg_alpha")
	if err != nil {
		t.Fatalf("ExecuteScriptWithEnv failed: %v", err)
	}

	if !res.Success() {
		t.Errorf("Expected script to succeed, exit code: %d", res.ExitCode)
	}
	if !res.Contains("Hello from script: arg_alpha") {
		t.Errorf("Expected stdout to contain passed arg, got: %s", res.Stdout)
	}
	if !res.Contains("ENV_FOO is: bar_value_123") {
		t.Errorf("Expected stdout to contain env variable, got: %s", res.Stdout)
	}
	if !strings.Contains(res.Stderr, "Notice: logging to stderr") {
		t.Errorf("Expected stderr to contain notice, got: %s", res.Stderr)
	}
	if len(res.StdoutLines()) < 3 {
		t.Errorf("Expected at least 3 non-empty stdout lines, got %d", len(res.StdoutLines()))
	}

	// Test JSON parsing on script stdout
	lines := res.StdoutLines()
	jsonLine := lines[len(lines)-1]
	var payload struct {
		Status string `json:"status"`
		Code   int    `json:"code"`
	}
	resJSON := &utils.ScriptResult{Stdout: jsonLine}
	if err := resJSON.JSON(&payload); err != nil {
		t.Fatalf("Failed to parse script JSON output: %v", err)
	}
	if payload.Status != "completed" || payload.Code != 200 {
		t.Errorf("Unexpected JSON payload values: %+v", payload)
	}

	// Test ExecuteScriptWithTimeout
	slowScriptPath := filepath.Join(tmpDir, "slow_script.sh")
	slowContent := `#!/bin/bash
sleep 2
echo "Finished"
`
	if err := os.WriteFile(slowScriptPath, []byte(slowContent), 0755); err != nil {
		t.Fatalf("Failed to write slow script: %v", err)
	}

	timeoutRes, _ := utils.ExecuteScriptWithTimeout(slowScriptPath, 50*time.Millisecond)
	if timeoutRes == nil || timeoutRes.Success() {
		t.Errorf("Expected slow script to fail/timeout, got success")
	}
}
