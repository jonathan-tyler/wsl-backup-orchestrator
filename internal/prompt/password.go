package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

type PasswordFunc func(message string) (string, error)

var isTerminalInput = term.IsTerminal
var readTerminalPassword = term.ReadPassword

func NewPasswordPrompt(input io.Reader, output io.Writer) PasswordFunc {
	reader := bufio.NewReader(input)
	inputFile, _ := input.(*os.File)

	return func(message string) (string, error) {
		promptText := fmt.Sprintf("%s: ", message)
		useHiddenInput := inputFile != nil && isTerminalInput(int(inputFile.Fd()))
		if useHiddenInput {
			promptText = fmt.Sprintf("%s (input hidden; not saved in shell history): ", message)
		}

		if _, err := io.WriteString(output, promptText); err != nil {
			return "", err
		}

		var line string
		if useHiddenInput {
			passwordBytes, err := readTerminalPassword(int(inputFile.Fd()))
			if _, writeErr := io.WriteString(output, "\n"); writeErr != nil {
				return "", writeErr
			}
			if err != nil {
				return "", err
			}
			line = string(passwordBytes)
		} else {
			readLine, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				return "", err
			}
			line = readLine
		}

		password := strings.TrimSpace(line)
		if password == "" {
			return "", fmt.Errorf("restic password is empty")
		}

		return password, nil
	}
}
