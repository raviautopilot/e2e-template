package client

import (
	"errors"
	"fmt"
	"testing"
)

func TestAsHttpError(t *testing.T) {
	// 1. Nil error
	if httpErr, ok := AsHttpError(nil); ok || httpErr != nil {
		t.Errorf("expected nil/false for nil error, got %v, %v", httpErr, ok)
	}

	// 2. Direct HttpError
	underlying := errors.New("underlying failure")
	origErr := &httpErrorImpl{
		statusCode: 404,
		respBody:   `{"error": "not found"}`,
		err:        underlying,
	}

	httpErr, ok := AsHttpError(origErr)
	if !ok || httpErr == nil {
		t.Fatalf("expected true for direct HttpError")
	}
	if httpErr.StatusCode() != 404 {
		t.Errorf("StatusCode = %d, want 404", httpErr.StatusCode())
	}
	if httpErr.ResponseBody() != `{"error": "not found"}` {
		t.Errorf("ResponseBody = %q, want %q", httpErr.ResponseBody(), `{"error": "not found"}`)
	}
	if httpErr.UnderlyingError() != underlying {
		t.Errorf("UnderlyingError = %v, want %v", httpErr.UnderlyingError(), underlying)
	}

	// 3. Wrapped HttpError (fmt.Errorf with %w)
	wrappedErr := fmt.Errorf("action failed: %w", origErr)
	extracted, ok := AsHttpError(wrappedErr)
	if !ok || extracted == nil {
		t.Fatalf("expected true for wrapped HttpError")
	}
	if extracted.StatusCode() != 404 {
		t.Errorf("StatusCode = %d, want 404", extracted.StatusCode())
	}

	// 4. Non-HttpError
	plainErr := errors.New("regular standard error")
	if _, ok := AsHttpError(plainErr); ok {
		t.Errorf("expected false for standard non-HttpError")
	}

	// 5. Unwrap test
	if !errors.Is(wrappedErr, underlying) {
		t.Errorf("expected errors.Is to find underlying error through Unwrap chain")
	}
}
