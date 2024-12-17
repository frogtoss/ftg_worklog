package main

import (
	"fmt"

	"github.com/frogtoss/ftg_worklog/pkg/frontmatter"
	"os"
	"path"
	"strings"

	// docs https://pkg.go.dev/github.com/elk-language/go-prompt
	"github.com/elk-language/go-prompt"
	istrings "github.com/elk-language/go-prompt/strings"
)

const MaxDescriptionLength = 32
const Prompt = "> "

const IncidentTMPL = `
# Incident Response Worklog #

## Means of Discovery ##

_How was the incident discovered?_

## Progress Notes ##

_Raw notes during progress to be preserved_

## Changes ##

_Set of changes that were made, and which envs they were deployed to_

## Reflections ##

_How to catch this sooner / prevent_

`

func commonOptions(aborted *bool) []prompt.Option {

	// ctrl c sets the bool
	ctrlCBind := prompt.KeyBind{
		Key: prompt.ControlC,
		Fn: func(p *prompt.Prompt) bool {
			*aborted = true
			return false
		},
	}

	// and then exitchecker checks if it was set on each input
	exitChecker := func(in string, breakline bool) bool {
		// exit prompt when ctrl-c was pressed
		return *aborted
	}

	return []prompt.Option{
		prompt.WithPrefix(Prompt),
		prompt.WithInputTextColor(prompt.White),
		prompt.WithInputBGColor(prompt.DarkBlue),
		prompt.WithPrefixTextColor(prompt.LightGray),
		prompt.WithKeyBind(ctrlCBind),
		prompt.WithExitChecker(exitChecker),
	}
}

// InteractivePrompt issues prompt with an optional completer,
// returning whether it was aborted, and then the string
func InteractivePrompt(completer prompt.Completer) (bool, string) {
	aborted := false

	options := commonOptions(&aborted)
	if completer != nil {
		options = append(options, prompt.WithCompleter(completer))
	}

	userStr := prompt.Input(options...)

	return aborted, userStr
}

func promptForIncidentDescription() string {

	fmt.Printf("Enter a description that is %d chars or less.\n", MaxDescriptionLength)
	fmt.Println("eg: \"registry down\", or \"slow response time\"")

	var desc string
	var aborted bool
	for len(desc) == 0 || len(desc) > MaxDescriptionLength {

		aborted, desc = InteractivePrompt(nil)
		handlePromptAbort(aborted)
	}

	return desc
}

func handlePromptAbort(abort bool) {
	if abort {
		fmt.Println("aborted")
		os.Exit(1)
	}
}

func promptForService(worklogDir string) string {
	fmt.Print("Which service to create incident for?\n")

	var serviceDirs []prompt.Suggest
	dirNames := make(map[string]int)

	// suggestions from directory names
	entries, err := os.ReadDir(worklogDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			if strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			serviceDirs = append(serviceDirs,
				prompt.Suggest{Text: entry.Name(), Description: ""},
			)

			dirNames[entry.Name()] = 1
		}
	}

	completer := func(d prompt.Document) (suggestions []prompt.Suggest, startChar, endChar istrings.RuneNumber) {
		endIndex := d.CurrentRuneIndex()
		w := d.GetWordBeforeCursor()
		startIndex := endIndex - istrings.RuneCount([]byte(w))

		return prompt.FilterHasPrefix(serviceDirs, w, true), startIndex, endIndex
	}

	var userService string
	for {

		aborted, userService := InteractivePrompt(completer)
		handlePromptAbort(aborted)

		_, match := dirNames[userService]
		if !match {
			fmt.Printf("service '%s' not found.\n", userService)
		} else {
			break
		}
	}

	return userService
}

func promptForConfirmation() bool {
	fmt.Println("Continue? (yes/no)")

	yesNo := []prompt.Suggest{
		{Text: "yes", Description: ""},
		{Text: "no", Description: "abort"},
	}

	completer := func(d prompt.Document) (suggestions []prompt.Suggest, startChar, endChar istrings.RuneNumber) {
		endIndex := d.CurrentRuneIndex()
		w := d.GetWordBeforeCursor()
		startIndex := endIndex - istrings.RuneCount([]byte(w))

		return prompt.FilterHasPrefix(yesNo, w, true), startIndex, endIndex
	}

	for {
		aborted, confirm := InteractivePrompt(completer)
		handlePromptAbort(aborted)

		switch confirm {
		case "yes":
			return true
		case "no":
			return false
		}

	}
}

func (i *CLIIncidentCmd) Run(cli *CLI) error {

	//
	// compute full path to generated file
	worklogDir, err := findWorklogDir(cli.WorklogDir)
	if err != nil {
		return err
	}

	serviceName := cli.Incident.Service
	if len(cli.Incident.Service) == 0 {
		serviceName = promptForService(worklogDir)
	}

	incident := frontmatter.NewIncidentWithService(serviceName)

	description := cli.Incident.Description
	if len(cli.Incident.Description) == 0 {
		description = promptForIncidentDescription()
	}

	worklogFilename := GenerateWorklogFilename("incident",
		"mlabbe",
		description,
		incident.Time.LogStart)

	serviceDir := ServiceDirFromName(incident.Service.Name)

	fullDir := path.Join(worklogDir, serviceDir)
	{
		exists, err := dirExists(fullDir)
		if !exists || err != nil {
			return fmt.Errorf("Directory '%s' does not exist.  Wrong service?\n%+v",
				fullDir, err)
		}
	}
	fullPath := path.Join(fullDir, worklogFilename)

	if !fileWouldBeNew(fullPath) {
		return fmt.Errorf("File '%s' already exists.", fullPath)
	}

	fmt.Printf("Create '%s'\n", fullPath)
	if !promptForConfirmation() {
		return fmt.Errorf("aborted")
	}

	//
	// generate full file content
	body := EncodeFrontmatter(incident)
	body += IncidentTMPL

	err = os.WriteFile(fullPath, []byte(body), 0644)
	if err != nil {
		return err
	}

	err = LaunchEditorForFile(fullPath)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Failed to launch editor for file '%s':\n%+v\nLaunch it manually\n",
			fullPath, err)

		// not a total loss -- file was already successfully written
		// so return nil here

	}

	return nil
}
