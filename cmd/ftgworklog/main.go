package main

import (
	"os"

	"fmt"
	"github.com/alecthomas/kong"
)

type CLI struct {
	// Common flags
	WorklogDir string `help:"The worklog dir" short:"w" type:"existingdir"`

	// Commands
	Incident IncidentCmd `cmd:"" help:"Generate a new incident"`
}

type IncidentCmd struct{}

func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func findWorklogDir(worklogDir string) (string, error) {

	// default location
	if worklogDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		worklogDir = fmt.Sprintf("%s/worklogs", homeDir)
	}

	exists, err := dirExists(worklogDir)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", fmt.Errorf("Could not find worklog dir '%s'", worklogDir)
	}

	return worklogDir, nil
}

func realMain() int {

	var cli CLI
	ctx := kong.Parse(&cli)
	err := ctx.Run(&cli)
	ctx.FatalIfErrorf(err)

	return 0
}

func main() {
	os.Exit(realMain())

}
