package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Repositories struct {
	Daily   string `yaml:"daily"`
	Weekly  string `yaml:"weekly"`
	Monthly string `yaml:"monthly"`
}

func (r Repositories) ForCadence(cadence string) (string, error) {
	switch cadence {
	case "daily":
		return strings.TrimSpace(r.Daily), nil
	case "weekly":
		return strings.TrimSpace(r.Weekly), nil
	case "monthly":
		return strings.TrimSpace(r.Monthly), nil
	default:
		return "", fmt.Errorf("unsupported cadence %q", cadence)
	}
}

type Profile struct {
	Repositories  Repositories `yaml:"repositories"`
	UseFSSnapshot bool         `yaml:"use_fs_snapshot"`
}

func (p Profile) RepositoryFor(cadence string) (string, error) {
	return p.Repositories.ForCadence(cadence)
}

type File struct {
	ResticVersion string             `yaml:"restic_version"`
	Profiles      map[string]Profile `yaml:"profiles"`

	path string
}

func (f File) Path() string {
	return f.path
}

func (f File) Dir() string {
	if f.path == "" {
		return ""
	}
	return filepath.Dir(f.path)
}

type Loader struct {
	ReadFile func(string) ([]byte, error)
	Getenv   func(string) string
}

func NewLoader() Loader {
	return Loader{ReadFile: os.ReadFile, Getenv: os.Getenv}
}

func (l Loader) Load() (File, error) {
	path, err := l.ResolvePath()
	if err != nil {
		return File{}, err
	}

	content, err := l.ReadFile(path)
	if err != nil {
		return File{}, fmt.Errorf("read config %q: %w", path, err)
	}

	var loaded File
	if err := yaml.Unmarshal(content, &loaded); err != nil {
		return File{}, fmt.Errorf("parse config %q: %w", path, err)
	}

	loaded.path = path
	if err := validate(loaded); err != nil {
		return File{}, err
	}

	return loaded, nil
}

func (l Loader) ResolvePath() (string, error) {
	if override := l.Getenv("BACKUP_CONFIG"); override != "" {
		return override, nil
	}

	xdgConfig := l.Getenv("XDG_CONFIG_HOME")
	if xdgConfig != "" {
		return filepath.Join(xdgConfig, "wsl-backup", "config.yaml"), nil
	}

	home := l.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("cannot resolve config path: HOME is unset")
	}

	return filepath.Join(home, ".config", "wsl-backup", "config.yaml"), nil
}

func validate(cfg File) error {
	if len(cfg.Profiles) == 0 {
		return fmt.Errorf("config has no profiles")
	}

	for profileName, profile := range cfg.Profiles {
		for _, cadence := range []string{"daily", "weekly", "monthly"} {
			repository, err := profile.RepositoryFor(cadence)
			if err != nil {
				return fmt.Errorf("profile %q repository lookup failed: %w", profileName, err)
			}
			if repository == "" {
				return fmt.Errorf("profile %q has empty %s repository", profileName, cadence)
			}
		}
		if profile.UseFSSnapshot && profileName != "windows" {
			return fmt.Errorf("profile %q enables use_fs_snapshot, but use_fs_snapshot is supported only for the windows profile", profileName)
		}
	}

	return nil
}

func FileWithPathForTest(file File, path string) File {
	file.path = path
	return file
}
