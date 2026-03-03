package setup

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/jonathan-tyler/wsl-backup-orchestrator/internal/apperr"
	"github.com/jonathan-tyler/wsl-backup-orchestrator/internal/config"
	"github.com/jonathan-tyler/wsl-backup-orchestrator/internal/prompt"
	"github.com/jonathan-tyler/wsl-backup-orchestrator/internal/restic"
	"github.com/jonathan-tyler/wsl-backup-orchestrator/internal/resticversion"
	"github.com/jonathan-tyler/wsl-backup-orchestrator/internal/system"
)

type ConfigLoader interface {
	Load() (config.File, error)
}

type Dependencies struct {
	Loader  ConfigLoader
	System  system.Executor
	Confirm prompt.ConfirmFunc
}

var (
	setupStdout io.Writer = os.Stdout
	setupStderr io.Writer = os.Stderr
	setupStdin  io.Reader = os.Stdin

	newConfigLoader = func() ConfigLoader {
		return config.NewLoader()
	}
	newSystemExecutor = func(stdout io.Writer, stderr io.Writer) system.Executor {
		return system.NewOSExecutor(stdout, stderr)
	}
	newConfirmFunc = func(input io.Reader, output io.Writer) prompt.ConfirmFunc {
		return prompt.NewYesNoConfirm(input, output)
	}
	handleWithReportFunc = handleWithReport
)

func Handle(ctx context.Context, args []string, _ restic.Executor) error {
	if len(args) != 0 {
		return apperr.UsageError{Message: "setup does not take positional arguments"}
	}

	fmt.Fprintln(setupStdout, "Running wsl-backup setup checks and installers...")

	deps := Dependencies{
		Loader:  newConfigLoader(),
		System:  newSystemExecutor(setupStdout, setupStderr),
		Confirm: newConfirmFunc(setupStdin, setupStdout),
	}

	report, err := handleWithReportFunc(ctx, deps)
	printSetupReport(report)
	if err == nil {
		fmt.Fprintln(setupStdout, "wsl-backup setup completed successfully.")
	} else {
		fmt.Fprintf(setupStdout, "wsl-backup setup failed: %v\n", err)
	}
	return err
}

func HandleWith(ctx context.Context, deps Dependencies) error {
	_, err := handleWithReport(ctx, deps)
	return err
}

func handleWithReport(ctx context.Context, deps Dependencies) (resticversion.SetupReport, error) {
	if deps.Loader == nil {
		deps.Loader = newConfigLoader()
	}
	if deps.System == nil {
		deps.System = newSystemExecutor(setupStdout, setupStderr)
	}
	if deps.Confirm == nil {
		deps.Confirm = func(string) (bool, error) { return false, nil }
	}

	cfg, err := deps.Loader.Load()
	if err != nil {
		return resticversion.SetupReport{}, err
	}

	return resticversion.SyncInteractiveWithReport(ctx, cfg, deps.System, deps.Confirm)
}

func printSetupReport(report resticversion.SetupReport) {
	if len(report.Items) == 0 {
		fmt.Fprintln(setupStdout, "setup report: no profile checks were executed")
		return
	}

	fmt.Fprintln(setupStdout, "setup report:")
	for _, item := range report.Items {
		fmt.Fprintf(setupStdout, "- %s: %s (%s)\n", item.Platform, item.Status, item.Message)
	}
}
