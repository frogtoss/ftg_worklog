package frontmatter

// frontmatter is the structured data that is prefixed before each
// markdown file

import (
	"os/user"
	"time"
)

const FrontmatterVersion = uint(1)

type IncidentService struct {
	Name string   `toml:"name"`
	Envs []string `toml:"envs"`
}

type IncidentTime struct {
	LogStart time.Time `toml:"log_start"`
}

type IncidentPersonnel struct {
	Authors   []string `toml:"authors"`
	Attending []string `toml:"attending"`
}

type IncidentRunbook struct {
	Path  string `toml:"path"`
	Title string `toml:"title"`
}

type Incident struct {
	Description string `toml:"description"`

	Type    string `toml:"type"`
	Version uint   `toml:"version"`

	Service IncidentService `toml:"service"`

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

	Time IncidentTime `toml:"time"`

	Personnel IncidentPersonnel `toml:"personnel"`

	Runbook []IncidentRunbook `toml:"runbook"`
}

func NewIncidentWithService(serviceName string) Incident {

	userData, err := user.Current()

	var username string
	if err != nil {
		username = "unknown"
	} else {
		username = userData.Username
	}

	return Incident{
		Type:        "incident",
		Description: "",
		Version:     FrontmatterVersion,

		Service: IncidentService{
			Name: serviceName,
			Envs: []string{""},
		},

		Time: IncidentTime{
			LogStart: time.Now(),
		},

		Personnel: IncidentPersonnel{
			Authors:   []string{username},
			Attending: []string{username},
		},

		Runbook: []IncidentRunbook{
			{
				Path:  "",
				Title: "",
			},
		},
	}
}
