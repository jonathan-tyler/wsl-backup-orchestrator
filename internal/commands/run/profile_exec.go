package run

import (
	"context"
	"strings"

	"github.com/example/wsl-backup/internal/restic"
	"github.com/example/wsl-backup/internal/system"
)

func executeProfileBackup(ctx context.Context, profileName string, resticArgs []string, runner restic.Executor, exec system.Executor) error {
	if strings.EqualFold(profileName, "windows") {
		return executeWindowsProfileBackup(ctx, resticArgs, exec)
	}

	return executeWSLProfileBackup(ctx, resticArgs, runner)
}
