package config

import (
	"sync"

	"github.com/nlopes/slack"
	"github.com/urfave/cli"
	"github.com/vivace-io/evelib/zkill"
	"github.com/vivace-io/gonfig"
)

var (
	CONFIG *Configuration
)

func Init(c *cli.Context) error {
	CONFIG = &Configuration{
		app: new(Application),
	}
	return gonfig.Load(CONFIG.app)
}

// Configuration contains returns the Application configuration in its current state,
// and protects it to be safe for use in goroutines.
type Configuration struct {
	sync.RWMutex
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
	UserAgent string  `json:"userAgent"`
	Teams     []*Team `json:"teams"`
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
	BotToken   string        `json:"botToken"`
	Channels   []*Channel    `json:"channels"`
	Bot        *slack.Client `json:"-"`
	FailedAuth bool          `json:"-"`
}

// Channel defines the configuration for a slack channel in a team, including its filters
type Channel struct {
	Name                string       `json:"channelName"`
	InBulk              bool         `json:"inBulk"`
	BulkPostInterval    int          `json:"bulkPostInterval"`
	MinimumValue        int          `json:"minimumValue"`
	MaximumValue        int          `json:"maximumValue"`
	IncludeCharacters   []string     `json:"includeCharacters"`
	IncludeCorporations []string     `json:"includeCorporations"`
	IncludeAlliances    []string     `json:"includeAlliance"`
	ExcludedShips       []string     `json:"excludedShips"`
	PendingBulk         []zkill.Kill `json:"-"`
}
