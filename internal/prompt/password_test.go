package prompt

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func TestNewPasswordPromptAcceptsInput(t *testing.T) {
	input := strings.NewReader("secret-value\n")
	var output strings.Builder
	prompt := NewPasswordPrompt(input, &output)

	password, err := prompt("Restic password")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if password != "secret-value" {
		t.Fatalf("unexpected password %q", password)
	}
	if !strings.Contains(output.String(), "Restic password") {
		t.Fatalf("expected prompt text in output")
	}
}

func TestNewPasswordPromptAcceptsEOFWithoutNewline(t *testing.T) {
	input := strings.NewReader("secret-value")
	prompt := NewPasswordPrompt(input, io.Discard)

	password, err := prompt("Restic password")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if password != "secret-value" {
		t.Fatalf("unexpected password %q", password)
	}
}

func TestNewPasswordPromptRejectsEmptyInput(t *testing.T) {
	input := strings.NewReader("   \n")
	prompt := NewPasswordPrompt(input, io.Discard)

	_, err := prompt("Restic password")
	if err == nil || !strings.Contains(err.Error(), "restic password is empty") {
		t.Fatalf("expected empty password error, got %v", err)
	}
}

func TestNewPasswordPromptUsesHiddenTerminalInput(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "password-input-*")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	t.Cleanup(func() {
		_ = file.Close()
	})

	originalIsTerminalInput := isTerminalInput
	originalReadTerminalPassword := readTerminalPassword
	t.Cleanup(func() {
		isTerminalInput = originalIsTerminalInput
		readTerminalPassword = originalReadTerminalPassword
	})

	isTerminalInput = func(fd int) bool {
		return fd == int(file.Fd())
	}
	readTerminalPassword = func(fd int) ([]byte, error) {
		if fd != int(file.Fd()) {
			t.Fatalf("unexpected file descriptor %d", fd)
		}
		return []byte("secret-value"), nil
	}

	var output strings.Builder
	prompt := NewPasswordPrompt(file, &output)

	password, err := prompt("Restic password")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if password != "secret-value" {
		t.Fatalf("unexpected password %q", password)
	}
	if !strings.Contains(output.String(), "input hidden") {
		t.Fatalf("expected hidden input note, got %q", output.String())
	}
	if !strings.Contains(output.String(), "shell history") {
		t.Fatalf("expected shell history note, got %q", output.String())
	}
	if !strings.HasSuffix(output.String(), "\n") {
		t.Fatalf("expected trailing newline after hidden input, got %q", output.String())
	}
}

func TestNewPasswordPromptReturnsHiddenInputError(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "password-input-*")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	t.Cleanup(func() {
		_ = file.Close()
	})

	originalIsTerminalInput := isTerminalInput
	originalReadTerminalPassword := readTerminalPassword
	t.Cleanup(func() {
		isTerminalInput = originalIsTerminalInput
		readTerminalPassword = originalReadTerminalPassword
	})

	isTerminalInput = func(fd int) bool {
		return fd == int(file.Fd())
	}
	readTerminalPassword = func(int) ([]byte, error) {
		return nil, errors.New("hidden read fail")
	}

	prompt := NewPasswordPrompt(file, io.Discard)

	_, err = prompt("Restic password")
	if err == nil || !strings.Contains(err.Error(), "hidden read fail") {
		t.Fatalf("expected hidden input error, got %v", err)
	}
}

func TestNewPasswordPromptReturnsOutputWriteError(t *testing.T) {
	prompt := NewPasswordPrompt(strings.NewReader("secret\n"), failingWriter{err: errors.New("write fail")})

	_, err := prompt("Restic password")
	if err == nil || !strings.Contains(err.Error(), "write fail") {
		t.Fatalf("expected write error, got %v", err)
	}
}

func TestNewPasswordPromptReturnsInputReadError(t *testing.T) {
	prompt := NewPasswordPrompt(failingReader{err: errors.New("read fail")}, io.Discard)

	_, err := prompt("Restic password")
	if err == nil || !strings.Contains(err.Error(), "read fail") {
		t.Fatalf("expected read error, got %v", err)
	}
}
