package run

import (
	"context"
	"testing"
)

func TestExecuteWSLProfileBackupUsesRunner(t *testing.T) {
	runner := &fakeRunner{}
	args := []string{"backup", "--tag", "cadence=daily", "/data/src"}

	err := executeWSLProfileBackup(context.Background(), args, runner)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(runner.calls) != 1 {
		t.Fatalf("expected one runner call, got %d", len(runner.calls))
	}
}
