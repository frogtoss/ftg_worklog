package main

import (
	"fmt"

	"github.com/frogtoss/ftg_worklog/pkg/frontmatter"
	"os"
	"path"

	"github.com/c-bata/go-prompt"
	"strings"
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

func commonOptions() []prompt.Option {
	return []prompt.Option{
		prompt.OptionInputTextColor(prompt.White),
		prompt.OptionInputBGColor(prompt.DarkBlue),
		prompt.OptionPrefixTextColor(prompt.LightGray),
		prompt.OptionSwitchKeyBindMode(prompt.EmacsKeyBind),
	}
}

func limitLengthFilter(buffer *prompt.Buffer) prompt.Buffer {
	if len(buffer.Text()) > 32 {
		buffer.InsertText(buffer.Text()[:32], false, true)
	}
	return *buffer
}

func promptForIncidentDescription() string {

	fmt.Printf("Enter a description that is %d chars or less.\n", MaxDescriptionLength)
	desc := ""
	for len(desc) == 0 || len(desc) > MaxDescriptionLength {

		desc = prompt.Input(
			Prompt,
			func(prompt.Document) []prompt.Suggest {
				return []prompt.Suggest{} // No autocompletion
			},
			commonOptions()...,
		)
	}

	return desc
}

func promptForService(worklogDir string) string {
	fmt.Print("Which service to create incident for?\n")

	var serviceDirs []prompt.Suggest

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
		}
	}

	service := prompt.Input(
		Prompt,
		func(d prompt.Document) []prompt.Suggest {
			return prompt.FilterHasPrefix(serviceDirs, d.Text, true)
		},
		commonOptions()...,
	)

	return service
}

func promptForConfirmation() bool {
	fmt.Println("Continue? (y/N)")

	completer := func(d prompt.Document) []prompt.Suggest {
		suggestions := []prompt.Suggest{
			{Text: "yes", Description: "Continue"},
			{Text: "no", Description: "Abort"},
		}
		return prompt.FilterHasPrefix(suggestions, d.Text, true)
	}

	input := ""
	for input == "" {
		input := prompt.Input(Prompt, completer, commonOptions()...)

		switch input {
		case "yes":
			return true
		case "no":
			return false
		default:
			input = ""
		}
	}

	return false
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
