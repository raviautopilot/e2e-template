package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var (
	cachedProjectRoot string
	cachedRootMu      sync.RWMutex
)

// FindProjectRoot searches upwards from the current working directory to locate the project root
// identified by the presence of a `go.mod` file or a `.git` directory.
func FindProjectRoot() (string, error) {
	cachedRootMu.RLock()
	if cachedProjectRoot != "" {
		defer cachedRootMu.RUnlock()
		return cachedProjectRoot, nil
	}
	cachedRootMu.RUnlock()

	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			cachedRootMu.Lock()
			cachedProjectRoot = dir
			cachedRootMu.Unlock()
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			cachedRootMu.Lock()
			cachedProjectRoot = dir
			cachedRootMu.Unlock()
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("project root not found (searched for go.mod and .git)")
}

// GetProjectRoot returns the project root path, falling back to current working directory if not found.
func GetProjectRoot() string {
	root, err := FindProjectRoot()
	if err != nil || root == "" {
		if cwd, err := os.Getwd(); err == nil {
			return cwd
		}
		return "."
	}
	return root
}

// GetModuleName parses the `go.mod` file at project root and extracts the module path name.
func GetModuleName() (string, error) {
	root, err := FindProjectRoot()
	if err != nil {
		return "", err
	}

	goModPath := filepath.Join(root, "go.mod")
	file, err := os.Open(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to open go.mod at %s: %w", goModPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed reading go.mod: %w", err)
	}

	return "", fmt.Errorf("module declaration not found in %s", goModPath)
}

// GetConfigPath locates the active config.json path, checking the project root and current working directory.
func GetConfigPath() (string, error) {
	root, err := FindProjectRoot()
	if err == nil {
		candidate := filepath.Join(root, "config.json")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	// Check relative paths
	relCandidates := []string{"config.json", "../config.json", "../../config.json", "../../../config.json"}
	for _, rel := range relCandidates {
		if _, err := os.Stat(rel); err == nil {
			abs, err := filepath.Abs(rel)
			if err == nil {
				return abs, nil
			}
			return rel, nil
		}
	}

	return "", fmt.Errorf("config.json not found")
}

// GetEvidenceDir returns the path to the evidence directory for a given run timestamp.
// If timestamp is not provided, returns the base evidence directory.
func GetEvidenceDir(runTimestamp ...string) (string, error) {
	root, err := FindProjectRoot()
	if err != nil {
		return "", err
	}
	evidenceBase := filepath.Join(root, "evidence")
	if len(runTimestamp) > 0 && runTimestamp[0] != "" {
		return filepath.Join(evidenceBase, "run-"+runTimestamp[0]), nil
	}
	return evidenceBase, nil
}

// GetGitBranch returns the current git branch name, or empty string if not in a git repo.
func GetGitBranch() string {
	root, err := FindProjectRoot()
	if err != nil {
		return ""
	}
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// FileExists checks if a file exists and is not a directory.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
