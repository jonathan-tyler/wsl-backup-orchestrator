package run

import (
	"context"

	"github.com/example/wsl-backup/internal/restic"
)

func executeWSLProfileBackup(ctx context.Context, resticArgs []string, runner restic.Executor) error {
	return runner.Run(ctx, resticArgs...)
}
