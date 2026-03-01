package prompt

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestNewYesNoConfirmAcceptsYes(t *testing.T) {
	input := strings.NewReader("yes\n")
	var output strings.Builder
	confirm := NewYesNoConfirm(input, &output)

	ok, err := confirm("Install now?")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !ok {
		t.Fatalf("expected confirmation to be true")
	}
	if !strings.Contains(output.String(), "Install now?") {
		t.Fatalf("expected prompt text in output")
	}
}

func TestNewYesNoConfirmDefaultsToNo(t *testing.T) {
	input := strings.NewReader("\n")
	var output strings.Builder
	confirm := NewYesNoConfirm(input, &output)

	ok, err := confirm("Install now?")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if ok {
		t.Fatalf("expected confirmation to be false")
	}
}

func TestNewYesNoConfirmAcceptsShortY(t *testing.T) {
	input := strings.NewReader("Y\n")
	var output strings.Builder
	confirm := NewYesNoConfirm(input, &output)

	ok, err := confirm("Install now?")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !ok {
		t.Fatalf("expected confirmation to be true")
	}
}

func TestNewYesNoConfirmAcceptsEOFWithoutNewline(t *testing.T) {
	input := strings.NewReader("yes")
	var output strings.Builder
	confirm := NewYesNoConfirm(input, &output)

	ok, err := confirm("Install now?")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !ok {
		t.Fatalf("expected confirmation to be true")
	}
}

func TestNewYesNoConfirmReturnsOutputWriteError(t *testing.T) {
	confirm := NewYesNoConfirm(strings.NewReader("yes\n"), failingWriter{err: errors.New("write fail")})

	_, err := confirm("Install now?")
	if err == nil || !strings.Contains(err.Error(), "write fail") {
		t.Fatalf("expected write error, got %v", err)
	}
}

func TestNewYesNoConfirmReturnsInputReadError(t *testing.T) {
	confirm := NewYesNoConfirm(failingReader{err: errors.New("read fail")}, io.Discard)

	_, err := confirm("Install now?")
	if err == nil || !strings.Contains(err.Error(), "read fail") {
		t.Fatalf("expected read error, got %v", err)
	}
}

type failingWriter struct {
	err error
}

func (w failingWriter) Write(_ []byte) (int, error) {
	return 0, w.err
}

type failingReader struct {
	err error
}

func (r failingReader) Read(_ []byte) (int, error) {
	return 0, r.err
}
