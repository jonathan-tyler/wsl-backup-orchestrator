package main

import (
	"io"
	"os"
	"reflect"
	"testing"
)

func TestMainForwardsArgsAndExitsWithCode(t *testing.T) {
	originalArgs := os.Args
	originalMain := cliMain
	originalExit := exitFunc
	t.Cleanup(func() {
		os.Args = originalArgs
		cliMain = originalMain
		exitFunc = originalExit
	})

	os.Args = []string{"backup", "run", "daily"}

	var gotArgs []string
	var gotCode int
	cliMain = func(args []string, _ io.Writer, _ io.Writer) int {
		gotArgs = append([]string{}, args...)
		return 7
	}
	exitFunc = func(code int) {
		gotCode = code
	}

	main()

	if !reflect.DeepEqual(gotArgs, []string{"run", "daily"}) {
		t.Fatalf("unexpected args forwarded to cli.Main: %#v", gotArgs)
	}
	if gotCode != 7 {
		t.Fatalf("expected exit code 7, got %d", gotCode)
	}
}
