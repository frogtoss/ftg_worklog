package main

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// LaunchEditorForFile launches the configured text editor
// for a given file.  On Windows, this uses file associations.
// On Non-Windows, this uses the EDITOR env var, falling back to
// whichever vi is in path.
func LaunchEditorForFile(path string) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "start", path)
	} else {

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}

		// handle EDITOR with command line args
		editorParts := strings.Fields(editor)
		exe := editorParts[0]
		args := append(editorParts[1:], path)

		cmd = exec.Command(exe, args...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
