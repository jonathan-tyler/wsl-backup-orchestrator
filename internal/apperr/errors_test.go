package apperr

import "testing"

func TestUsageErrorImplementsErrorMessage(t *testing.T) {
	err := UsageError{Message: "bad usage"}
	if err.Error() != "bad usage" {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}
