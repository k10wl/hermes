package cli

import (
	"io"
	"os"
	"os/exec"
)

func OpenInEditor(
	input string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) (string, error) {
	tmpFile, err := os.CreateTemp("", "hermes.*.tmp")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(input)
	if err != nil {
		return "", err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	// Run the editor command and wait for it to complete
	if err := cmd.Run(); err != nil {
		return "", err
	}

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", err
	}

	return string(content), nil
}
