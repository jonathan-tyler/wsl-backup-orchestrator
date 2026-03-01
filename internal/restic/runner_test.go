package restic

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatCommandQuotesWhitespace(t *testing.T) {
	formatted := formatCommand([]string{"backup", "--target", "/tmp/my folder"})
	if formatted != "backup --target \"/tmp/my folder\"" {
		t.Fatalf("unexpected format: %s", formatted)
	}
}

func TestOSRunnerPrintsCommandAndStreamsOutput(t *testing.T) {
	original := commandContext
	commandContext = fakeExecCommand
	t.Cleanup(func() {
		commandContext = original
	})

	configPath := writeConfigFile(t, `restic_version: "0.18.1"
keepassxc_database: /tmp/vault.kdbx
keepassxc_entry: restic/main
profiles:
  wsl:
    repository: /repo/wsl
    use_fs_snapshot: false
`)
	t.Setenv("BACKUP_CONFIG", configPath)
	t.Setenv(KeepassDatabaseEnv, "")
	t.Setenv(KeepassEntryEnv, "")

	var stdout strings.Builder
	var stderr strings.Builder
	runner := NewOSRunner(&stdout, &stderr)

	err := runner.Run(context.Background(), "snapshots")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "$ keepassxc-cli show -q -a Password /tmp/vault.kdbx restic/main") {
		t.Fatalf("expected keepass lookup command in stdout, got %q", out)
	}
	if !strings.Contains(out, "$ restic snapshots") {
		t.Fatalf("expected echoed command in stdout, got %q", out)
	}
	if !strings.Contains(out, "helper stdout") {
		t.Fatalf("expected command stdout in stdout writer, got %q", out)
	}
	if !strings.Contains(stderr.String(), "helper stderr") {
		t.Fatalf("expected command stderr in stderr writer, got %q", stderr.String())
	}
}

func TestOSRunnerRejectsEmptyArgs(t *testing.T) {
	var stdout strings.Builder
	var stderr strings.Builder
	runner := NewOSRunner(&stdout, &stderr)
	err := runner.Run(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestOSRunnerReadsKeepassSettingsFromEnvOverrides(t *testing.T) {
	original := commandContext
	commandContext = fakeExecCommand
	t.Cleanup(func() {
		commandContext = original
	})

	configPath := writeConfigFile(t, `restic_version: "0.18.1"
keepassxc_database: /tmp/config-vault.kdbx
keepassxc_entry: config/restic
profiles:
  wsl:
    repository: /repo/wsl
    use_fs_snapshot: false
`)
	t.Setenv("BACKUP_CONFIG", configPath)
	t.Setenv(KeepassDatabaseEnv, "/tmp/env-vault.kdbx")
	t.Setenv(KeepassEntryEnv, "env/restic")

	var stdout strings.Builder
	var stderr strings.Builder
	runner := NewOSRunner(&stdout, &stderr)

	err := runner.Run(context.Background(), "snapshots")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	out := stdout.String()
	if !strings.Contains(out, "$ keepassxc-cli show -q -a Password /tmp/env-vault.kdbx env/restic") {
		t.Fatalf("expected env override keepass command, got %q", out)
	}
}

func TestOSRunnerFailsWhenKeepassSettingsMissing(t *testing.T) {
	original := commandContext
	commandContext = fakeExecCommand
	t.Cleanup(func() {
		commandContext = original
	})

	configPath := writeConfigFile(t, `restic_version: "0.18.1"
profiles:
  wsl:
    repository: /repo/wsl
    use_fs_snapshot: false
`)
	t.Setenv("BACKUP_CONFIG", configPath)
	t.Setenv(KeepassDatabaseEnv, "")
	t.Setenv(KeepassEntryEnv, "")

	var stdout strings.Builder
	var stderr strings.Builder
	runner := NewOSRunner(&stdout, &stderr)

	err := runner.Run(context.Background(), "snapshots")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "missing KeepassXC lookup settings") {
		t.Fatalf("expected missing settings error, got %v", err)
	}
}

func TestOSRunnerFailsWhenKeepassLookupFails(t *testing.T) {
	original := commandContext
	commandContext = fakeExecCommand
	t.Cleanup(func() {
		commandContext = original
	})

	t.Setenv(KeepassDatabaseEnv, "/tmp/vault.kdbx")
	t.Setenv(KeepassEntryEnv, "restic/main")
	t.Setenv("FAKE_KEEPASS_FAIL", "1")

	var stdout strings.Builder
	var stderr strings.Builder
	runner := NewOSRunner(&stdout, &stderr)

	err := runner.Run(context.Background(), "snapshots")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "keepassxc-cli password lookup failed") {
		t.Fatalf("expected keepass lookup error, got %v", err)
	}
}

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}
	return path
}

func fakeExecCommand(_ context.Context, name string, args ...string) *exec.Cmd {
	allArgs := []string{"-test.run=TestHelperProcess", "--"}
	allArgs = append(allArgs, name)
	allArgs = append(allArgs, args...)
	cmd := exec.Command(os.Args[0], allArgs...)
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	if len(os.Args) < 4 {
		os.Exit(2)
	}

	commandName := os.Args[3]
	if commandName == "keepassxc-cli" {
		if os.Getenv("FAKE_KEEPASS_FAIL") == "1" {
			fmt.Fprintln(os.Stderr, "database is locked")
			os.Exit(1)
		}

		fmt.Fprintln(os.Stdout, "test-password")
		os.Exit(0)
	}

	if commandName == "restic" {
		if os.Getenv("RESTIC_PASSWORD") == "" {
			fmt.Fprintln(os.Stderr, "missing RESTIC_PASSWORD")
			os.Exit(2)
		}

		fmt.Fprintln(os.Stdout, "helper stdout")
		fmt.Fprintln(os.Stderr, "helper stderr")
		os.Exit(0)
	}

	fmt.Fprintln(os.Stderr, "unknown helper command")
	os.Exit(2)
}
