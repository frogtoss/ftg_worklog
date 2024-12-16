package main

import (
	"fmt"

	"github.com/frogtoss/ftg_worklog/pkg/frontmatter"
	"os"
	"path"
)

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

func (i *CLIIncidentCmd) Run(cli *CLI) error {

	//
	// compute full path to generated file
	worklogDir, err := findWorklogDir(cli.WorklogDir)
	if err != nil {
		return err
	}

	incident := frontmatter.NewIncidentWithService(cli.Incident.Service)

	incident.Description = "system test"

	worklogFilename := GenerateWorklogFilename("incident",
		"mlabbe",
		incident.Description,
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

	fmt.Printf("Creating '%s'\n", fullPath)

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
		// not a total failure -- file was already successfully written
		fmt.Fprintf(os.Stderr, "Failed to launch editor for file '%s':\n%+v\nLaunch it manually\n",
			fullPath, err)
	}

	return nil
}
