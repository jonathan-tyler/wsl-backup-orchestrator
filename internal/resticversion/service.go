package resticversion

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/example/wsl-backup/internal/config"
	"github.com/example/wsl-backup/internal/prompt"
	"github.com/example/wsl-backup/internal/system"
)

func CheckCompatible(ctx context.Context, cfg config.File, exec system.Executor) error {
	desiredVersion := strings.TrimSpace(cfg.ResticVersion)
	if desiredVersion == "" {
		return nil
	}

	if _, hasWSL := cfg.Profiles["wsl"]; hasWSL {
		if err := checkWSLCompatible(ctx, desiredVersion, exec); err != nil {
			return err
		}
	}

	if _, hasWindows := cfg.Profiles["windows"]; hasWindows {
		if err := checkWindowsCompatible(ctx, desiredVersion, exec); err != nil {
			return err
		}
	}

	return nil
}

func SyncInteractive(ctx context.Context, cfg config.File, exec system.Executor, confirm prompt.ConfirmFunc) error {
	desiredVersion := strings.TrimSpace(cfg.ResticVersion)
	if desiredVersion == "" {
		return nil
	}

	if _, hasWSL := cfg.Profiles["wsl"]; hasWSL {
		if err := syncWSLInteractive(ctx, desiredVersion, exec, confirm); err != nil {
			return err
		}
	}

	if _, hasWindows := cfg.Profiles["windows"]; hasWindows {
		if err := syncWindowsInteractive(ctx, desiredVersion, exec, confirm); err != nil {
			return err
		}
	}

	return nil
}

var versionPattern = regexp.MustCompile(`\b(\d+\.\d+\.\d+)\b`)

func parseResticVersion(output string) (string, error) {
	match := versionPattern.FindStringSubmatch(output)
	if len(match) < 2 {
		return "", fmt.Errorf("could not find version in output %q", output)
	}
	return match[1], nil
}
