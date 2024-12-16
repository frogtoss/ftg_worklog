package main

import (
	"os"

	"bytes"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/pelletier/go-toml"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

type CLI struct {
	// Common flags
	WorklogDir string `help:"The worklog dir" short:"w" type:"existingdir"`

	// Commands
	Incident CLIIncidentCmd `cmd:"" help:"Generate a new incident"`
}

type CLIIncidentCmd struct {
	Service string `help:"The affected service" required:""`
}

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

func fileWouldBeNew(path string) bool {
	_, err := os.Stat(path)
	return err != nil
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

// Encode any into toml, and wrap in +++ frontmatter chars
func EncodeFrontmatter(doc any) string {
	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	encoder.Encode(doc)
	return fmt.Sprintf("+++\n%s+++\n", buf.String())
}

func GenerateWorklogFilename(worklogType, attendingUser, desc string, startTime time.Time) string {
	// worklog filename like:
	//2024.12.08.Incident.HardDriveSpace.MichaelLabbe.md

	fnTrimStr := func(s string) string {
		s = cases.Title(language.English).String(s)
		return strings.ReplaceAll(s, " ", "")
	}

	formattedDate := startTime.Format("2006.01.02")
	worklogType = fnTrimStr(worklogType)
	attendingUser = fnTrimStr(attendingUser)
	desc = fnTrimStr(desc)

	return fmt.Sprintf("%s.%s.%s.%s.md", formattedDate, worklogType, attendingUser, desc)
}

func ServiceDirFromName(serviceName string) string {
	return cases.Lower(language.English).String(serviceName)
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
