package resticversion

import (
	"context"
	"strings"
	"testing"

	"github.com/example/wsl-backup/internal/config"
)

func TestCheckCompatibleFailsWithSetupHintOnWindowsMismatch(t *testing.T) {
	exec := &fakeSystemExecutor{
		captureOutput: map[string]string{
			"powershell.exe -NoProfile -Command restic version": "restic 0.17.3 compiled with go1.24",
		},
	}

	err := CheckCompatible(context.Background(), config.File{
		ResticVersion: "0.18.1",
		Profiles:      map[string]config.Profile{"windows": {Repository: `C:\repo`}},
	}, exec)

	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "run backup setup") {
		t.Fatalf("expected setup hint, got %v", err)
	}
}

func TestSyncInteractiveWindowsUpdateDeclinedReturnsError(t *testing.T) {
	exec := &fakeSystemExecutor{
		captureOutput: map[string]string{
			"powershell.exe -NoProfile -Command restic version": "restic 0.17.3 compiled with go1.24",
		},
		captureErr: map[string]error{},
		runErr:     map[string]error{},
	}

	err := SyncInteractive(context.Background(), config.File{
		ResticVersion: "0.18.1",
		Profiles:      map[string]config.Profile{"windows": {Repository: `C:\repo`}},
	}, exec, func(string) (bool, error) {
		return false, nil
	})

	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "version mismatch") {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exec.runCalls) != 0 {
		t.Fatalf("did not expect update command when declined")
	}
}
