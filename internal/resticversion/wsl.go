package resticversion

import (
	"context"
	"fmt"

	"github.com/example/wsl-backup/internal/prompt"
	"github.com/example/wsl-backup/internal/system"
)

func checkWSLCompatible(ctx context.Context, desiredVersion string, exec system.Executor) error {
	output, err := exec.RunCapture(ctx, "restic", "version")
	if err != nil {
		return fmt.Errorf("wsl restic is missing; run backup setup")
	}

	installedVersion, parseErr := parseResticVersion(output)
	if parseErr != nil {
		return fmt.Errorf("parse wsl restic version: %w", parseErr)
	}

	if installedVersion != desiredVersion {
		return fmt.Errorf("wsl restic version mismatch: installed=%s required=%s; run backup setup", installedVersion, desiredVersion)
	}

	return nil
}

func syncWSLInteractive(ctx context.Context, desiredVersion string, exec system.Executor, confirm prompt.ConfirmFunc) error {
	output, err := exec.RunCapture(ctx, "restic", "version")
	if err != nil {
		approved, confirmErr := confirm("WSL restic not found. Install via dnf now?")
		if confirmErr != nil {
			return confirmErr
		}
		if !approved {
			return fmt.Errorf("wsl restic is required")
		}
		return exec.Run(ctx, "sudo", "dnf", "install", "-y", "restic")
	}

	installedVersion, parseErr := parseResticVersion(output)
	if parseErr != nil {
		return fmt.Errorf("parse wsl restic version: %w", parseErr)
	}

	if installedVersion == desiredVersion {
		return nil
	}

	approved, confirmErr := confirm(fmt.Sprintf("WSL restic version is %s but config requires %s. Upgrade via dnf now?", installedVersion, desiredVersion))
	if confirmErr != nil {
		return confirmErr
	}
	if !approved {
		return fmt.Errorf("wsl restic version mismatch: installed=%s required=%s", installedVersion, desiredVersion)
	}

	return exec.Run(ctx, "sudo", "dnf", "upgrade", "-y", "restic")
}
