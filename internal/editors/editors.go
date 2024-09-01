package editors

import (
	"os"
	"os/exec"
	"runtime"
)

func DefaultEditor() string {
	switch runtime.GOOS {
	case "linux":
		return "vim"
	case "windows":
		return "notepad"
	case "darwin":
		return "nano"
	default:
		return "nano"
	}
}

func OpenEditor(editor, file string) *exec.Cmd {
	cmd := exec.Command(editor, file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd
}
