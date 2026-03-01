# wsl-backup

Thin, predictable wrapper around `restic`, run from WSL and usable as a `wsl-sys-cli` extension.

## What it does

- Runs from WSL and targets cross-platform backup flows (WSL + Windows)
- Checks for matching `restic` versions on WSL (via `dnf`) and Windows (via `scoop`), and offers to install/upgrade when mismatched or missing
- Supports *cadences* of `daily`, `weekly`, and `monthly` for backup and reporting commands

## Configuration

- Default config path: `${XDG_CONFIG_HOME:-~/.config}/wsl-backup/config.yaml`
- Optional config override: `BACKUP_CONFIG=/path/to/config.yaml`
- Starter config: [config.example.yaml](config.example.yaml)
- Rule file directory: `~/.config/wsl-backup/`
- Rule naming: `<profile>.<include|exclude>.<daily|weekly|monthly>.txt`
  - Rules are checked for filesystem overlap

## Usage

- This CLI is WSL-only; run it from a WSL shell (not from native Windows or a Dev Container).

```sh
# Run wsl and windows profiles in parallel
backup run <cadence>

# Restore a backup to a target path
backup restore <target>
```

- If installed through `wsl-sys-cli`, run the same commands as `sys backup ...`.

## Caveats

- `restic` stores symlinks as symlinks by default and does not follow them during backup. This behavior helps avoid recursive traversal from link loops. If symlink following is enabled explicitly in a restic invocation, traversal/loop risk must be evaluated separately.
