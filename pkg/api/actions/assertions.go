package actions

import (
	"fmt"

	"e2e-template/tests"
)

// ─────────────────────────────────────────────────────────────────────────────
// Assertion Helpers
// ─────────────────────────────────────────────────────────────────────────────

// AssertNotEmpty fails the test if the value is empty.
func AssertNotEmpty(tc *tests.TestContext, fieldName string, value string) {
	if value == "" {
		tc.FailureReason = fmt.Sprintf("Expected non-empty %s", fieldName)
		tc.Errorf("Expected non-empty %s", fieldName)
	}
}

// AssertEquals fails the test if got != want.
func AssertEquals(tc *tests.TestContext, fieldName string, got, want string) {
	if got != want {
		tc.FailureReason = fmt.Sprintf("Expected %s=%q, got %q", fieldName, want, got)
		tc.Errorf("Expected %s=%q, got %q", fieldName, want, got)
	}
}

// AssertIntEquals fails the test if got != want.
func AssertIntEquals(tc *tests.TestContext, fieldName string, got, want int) {
	if got != want {
		tc.FailureReason = fmt.Sprintf("Expected %s=%d, got %d", fieldName, want, got)
		tc.Errorf("Expected %s=%d, got %d", fieldName, want, got)
	}
}

// AssertNotZero fails the test if the value is zero.
func AssertNotZero(tc *tests.TestContext, fieldName string, value int) {
	if value == 0 {
		tc.FailureReason = fmt.Sprintf("Expected non-zero %s", fieldName)
		tc.Errorf("Expected non-zero %s", fieldName)
	}
}

// AssertNonEmptyList fails the test if the list is empty.
func AssertNonEmptyList(tc *tests.TestContext, fieldName string, length int) {
	if length == 0 {
		tc.FailureReason = fmt.Sprintf("Expected non-empty %s list", fieldName)
		tc.Errorf("Expected non-empty %s list", fieldName)
	}
}

// AssertListLength fails the test if the list length doesn't match.
func AssertListLength(tc *tests.TestContext, fieldName string, got, want int) {
	if got != want {
		tc.FailureReason = fmt.Sprintf("Expected %s count=%d, got %d", fieldName, want, got)
		tc.Errorf("Expected %s count=%d, got %d", fieldName, want, got)
	}
}

// AssertTrue fails the test if condition is false.
func AssertTrue(tc *tests.TestContext, fieldName string, condition bool) {
	if !condition {
		tc.FailureReason = fmt.Sprintf("Expected %s to be true", fieldName)
		tc.Errorf("Expected %s to be true", fieldName)
	}
}

// AssertFalse fails the test if condition is true.
func AssertFalse(tc *tests.TestContext, fieldName string, condition bool) {
	if condition {
		tc.FailureReason = fmt.Sprintf("Expected %s to be false", fieldName)
		tc.Errorf("Expected %s to be false", fieldName)
	}
}
