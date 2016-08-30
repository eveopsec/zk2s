package config

import (
	"sync"

	"github.com/urfave/cli"
	"github.com/vivace-io/gonfig"
)

var (
	CONFIG *Configuration
)

func Init(c *cli.Context) error {
	return nil
}

// Configuration contains returns the Application configuration in its current state,
// and protects it to be safe for use in goroutines.
type Configuration struct {
	*sync.RWMutex
	app *Application
}

func (this *Configuration) Get() (app Application, err error) {
	if this.app == nil {
		this.Lock()
		defer this.Unlock()
		return *this.app, gonfig.Load(this.app)
	}
	this.RLock()
	defer this.RUnlock()
	return *this.app, nil
}

type Application struct {
	UserAgent string `json:"userAgent"`
	Teams     []Team `json:"teams"`
}

// File returns the file name/path for gonfig interface
func (this *Application) File() string {
	return "cfg.zk2s.json"
}

// Save the configuration file
func (this *Application) Save() error {
	return gonfig.Save(this)
}

type Team struct {
	BotToken string    `json:"botToken"`
	Channels []Channel `json:"channels"`
}

// Channel defines the configuration for a slack channel in a team, including its filters
type Channel struct {
	Name                string   `json:"channelName"`
	MinimumValue        int      `json:"minimumValue"`
	MaximumValue        int      `json:"maximumValue"`
	IncludeCharacters   []string `json:"includeCharacters"`
	IncludeCorporations []string `json:"includeCorporations"`
	IncludeAlliances    []string `json:"includeAlliance"`
	ExcludedShips       []string `json:"excludedShips"`
}
