package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ScriptResult represents the execution outcome of a script or shell command.
type ScriptResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
	Error    error
}

// Success returns true if the script exited with status code 0 and without error.
func (r *ScriptResult) Success() bool {
	return r.ExitCode == 0 && r.Error == nil
}

// CombinedOutput returns stdout followed by stderr if stderr is non-empty.
func (r *ScriptResult) CombinedOutput() string {
	if r.Stderr == "" {
		return r.Stdout
	}
	if r.Stdout == "" {
		return r.Stderr
	}
	return r.Stdout + "\n" + r.Stderr
}

// StdoutLines returns stdout split into trimmed lines, filtering empty lines.
func (r *ScriptResult) StdoutLines() []string {
	return splitNonEmptyLines(r.Stdout)
}

// StderrLines returns stderr split into trimmed lines, filtering empty lines.
func (r *ScriptResult) StderrLines() []string {
	return splitNonEmptyLines(r.Stderr)
}

// JSON unmarshals the stdout into the provided target pointer v.
func (r *ScriptResult) JSON(v interface{}) error {
	trimmed := strings.TrimSpace(r.Stdout)
	if trimmed == "" {
		return errors.New("cannot unmarshal empty stdout to JSON")
	}
	return json.Unmarshal([]byte(trimmed), v)
}

// Contains returns true if the script stdout or stderr contains the given substring.
func (r *ScriptResult) Contains(substr string) bool {
	return strings.Contains(r.Stdout, substr) || strings.Contains(r.Stderr, substr)
}

// ScriptOptions defines configuration parameters for executing a script.
type ScriptOptions struct {
	ScriptPath string
	Args       []string
	Env        map[string]string
	Cwd        string
	Timeout    time.Duration
}

// ExecuteScript runs a shell script (.sh) with optional arguments.
// If the script path is relative, it is resolved against the project root.
func ExecuteScript(scriptPath string, args ...string) (*ScriptResult, error) {
	return ExecuteScriptWithOptions(ScriptOptions{
		ScriptPath: scriptPath,
		Args:       args,
	})
}

// ExecuteScriptWithEnv runs a shell script with custom environment variables.
func ExecuteScriptWithEnv(scriptPath string, env map[string]string, args ...string) (*ScriptResult, error) {
	return ExecuteScriptWithOptions(ScriptOptions{
		ScriptPath: scriptPath,
		Env:        env,
		Args:       args,
	})
}

// ExecuteScriptWithTimeout runs a shell script bounded by the given timeout duration.
func ExecuteScriptWithTimeout(scriptPath string, timeout time.Duration, args ...string) (*ScriptResult, error) {
	return ExecuteScriptWithOptions(ScriptOptions{
		ScriptPath: scriptPath,
		Timeout:    timeout,
		Args:       args,
	})
}

// ExecuteScriptWithOptions executes a script according to the detailed ScriptOptions.
func ExecuteScriptWithOptions(opts ScriptOptions) (*ScriptResult, error) {
	resolvedPath := opts.ScriptPath
	if !filepath.IsAbs(resolvedPath) {
		root := GetProjectRoot()
		candidate := filepath.Join(root, resolvedPath)
		if _, err := os.Stat(candidate); err == nil {
			resolvedPath = candidate
		} else if _, err := os.Stat(opts.ScriptPath); err == nil {
			if abs, err := filepath.Abs(opts.ScriptPath); err == nil {
				resolvedPath = abs
			}
		}
	}

	if _, err := os.Stat(resolvedPath); err != nil {
		return nil, fmt.Errorf("script not found at %s: %w", resolvedPath, err)
	}

	// Ensure executable bit if possible
	if info, err := os.Stat(resolvedPath); err == nil && info.Mode()&0111 == 0 {
		_ = os.Chmod(resolvedPath, info.Mode()|0755)
	}

	workingDir := opts.Cwd
	if workingDir == "" {
		workingDir = filepath.Dir(resolvedPath)
	}

	ctx := context.Background()
	var cancel context.CancelFunc
	if opts.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	cmdArgs := append([]string{resolvedPath}, opts.Args...)
	cmd := exec.CommandContext(ctx, "bash", cmdArgs...)
	cmd.Dir = workingDir

	// Setup environment
	cmd.Env = os.Environ()
	for k, v := range opts.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	result := &ScriptResult{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		ExitCode: exitCode,
		Duration: duration,
		Error:    err,
	}

	return result, err
}

// ExecuteCommand runs an arbitrary command and returns the ScriptResult.
func ExecuteCommand(name string, args ...string) (*ScriptResult, error) {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return &ScriptResult{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		ExitCode: exitCode,
		Duration: duration,
		Error:    err,
	}, err
}

func splitNonEmptyLines(s string) []string {
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	return lines
}
