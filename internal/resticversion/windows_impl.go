package resticversion

import (
	"context"
	"fmt"

	"github.com/example/wsl-backup/internal/prompt"
	"github.com/example/wsl-backup/internal/system"
)

func checkWindowsCompatible(ctx context.Context, desiredVersion string, exec system.Executor) error {
	output, err := exec.RunCapture(ctx, "powershell.exe", "-NoProfile", "-Command", "restic version")
	if err != nil {
		return fmt.Errorf("windows restic is missing; run backup setup")
	}

	installedVersion, parseErr := parseResticVersion(output)
	if parseErr != nil {
		return fmt.Errorf("parse windows restic version: %w", parseErr)
	}

	if installedVersion != desiredVersion {
		return fmt.Errorf("windows restic version mismatch: installed=%s required=%s; run backup setup", installedVersion, desiredVersion)
	}

	return nil
}

func syncWindowsInteractive(ctx context.Context, desiredVersion string, exec system.Executor, confirm prompt.ConfirmFunc) error {
	output, err := exec.RunCapture(ctx, "powershell.exe", "-NoProfile", "-Command", "restic version")
	if err != nil {
		approved, confirmErr := confirm("Windows restic not found. Install via scoop now?")
		if confirmErr != nil {
			return confirmErr
		}
		if !approved {
			return fmt.Errorf("windows restic is required")
		}
		return exec.Run(ctx, "powershell.exe", "-NoProfile", "-Command", "scoop install restic")
	}

	installedVersion, parseErr := parseResticVersion(output)
	if parseErr != nil {
		return fmt.Errorf("parse windows restic version: %w", parseErr)
	}

	if installedVersion == desiredVersion {
		return nil
	}

	approved, confirmErr := confirm(fmt.Sprintf("Windows restic version is %s but config requires %s. Update via scoop now?", installedVersion, desiredVersion))
	if confirmErr != nil {
		return confirmErr
	}
	if !approved {
		return fmt.Errorf("windows restic version mismatch: installed=%s required=%s", installedVersion, desiredVersion)
	}

	return exec.Run(ctx, "powershell.exe", "-NoProfile", "-Command", "scoop update restic")
}
