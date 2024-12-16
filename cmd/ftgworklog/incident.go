package main

import (
	"fmt"
	"os/user"
	"time"
)

type PrefixIncidentService struct {
	Name string   `toml:"name"`
	Envs []string `toml:"envs"`
}

type PrefixIncidentTime struct {
	LogStart time.Time `toml:"log_start"`
}

type PrefixIncidentPersonnel struct {
	Authors   []string `toml:"authors"`
	Attending []string `toml:"attending"`
}

type PrefixIncidentRunbook struct {
	Path  string `toml:"path"`
	Title string `toml:"title"`
}

type PrefixIncident struct {
	Type        string `toml:"type"`
	Description string `toml:"description"`

	Service PrefixIncidentService `toml:"service"`

	Severity struct {
		Rating int `toml:"rating"`

		DeveloperCritical bool `toml:"developer_critical"`
		DeveloperPartial  bool `toml:"developer_partial"`
		CustomerCritical  bool `toml:"customer_critical"`
		CustomerPartial   bool `toml:"customer_partial"`
		HardwareFailure   bool `toml:"hardware_failure"`
		AffectsRevenue    bool `toml:"affects_revenue"`
		DataCorruption    bool `toml:"data_corruption"`
	} `toml:"severity"`

	Time PrefixIncidentTime `toml:"time"`

	Personnel PrefixIncidentPersonnel `toml:"personnel"`

	Runbook []PrefixIncidentRunbook `toml:"runbook"`
}

func NewPrefixIncidentWithService(serviceName string) PrefixIncident {

	userData, err := user.Current()

	var username string
	if err != nil {
		username = "unknown"
	} else {
		username = userData.Username
	}

	return PrefixIncident{
		Type:        "incident",
		Description: "",

		Service: PrefixIncidentService{
			Name: serviceName,
			Envs: []string{""},
		},

		Time: PrefixIncidentTime{
			LogStart: time.Now(),
		},

		Personnel: PrefixIncidentPersonnel{
			Authors:   []string{username},
			Attending: []string{username},
		},

		Runbook: []PrefixIncidentRunbook{
			{
				Path:  "<under ftg_sites>",
				Title: "",
			},
		},
	}
}

func (i *IncidentCmd) Run(cli *CLI) error {
	fmt.Println("Worklog Path:", cli.WorklogDir)

	worklogDir, err := findWorklogDir(cli.WorklogDir)
	if err != nil {
		return err
	}

	fmt.Printf("found worklog dir: '%s'\n", worklogDir)

	// todo: get service from command line, and it has to match a dir under worklog dir

	fmt.Printf("%+v\n", NewPrefixIncidentWithService("fuckchew"))

	return nil
}
