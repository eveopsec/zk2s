package config

import (
	"errors"
	"sync"

	"github.com/nlopes/slack"
	"github.com/vivace-io/gonfig"
)

// ReadConfig accepts a file path and returns the parsed configuration, if
// found. If the file was not found, or could not be parsed, an error is
// returned and the configuration model is nil.
func ReadConfig(filepath string) (cfg *Configuration, err error) {
	cfg = new(Configuration)
	cfg.FilePath = filepath
	cfg.locker = new(sync.RWMutex)
	if err := cfg.Load(); err != nil {
		return nil, err
	}
	return
}

// Configuration holds the configuration for the application.
type Configuration struct {
	locker       *sync.RWMutex
	FilePath     string  `json:"-"`
	TemplateFile string  `json:"templateFile"`
	UserAgent    string  `json:"userAgent"`
	Teams        []*Team `json:"teams"`
}

// File returns the file name/path for gonfig interface
func (cfg *Configuration) File() string {
	return cfg.FilePath
}

// Save the configuration file on the file path set in Configuration.FilePath.
func (cfg *Configuration) Save() error {
	cfg.locker.Lock()
	defer cfg.locker.Unlock()
	if cfg.FilePath == "" {
		return errors.New("Configuration.FilePath was not set")
	}
	return gonfig.Save(cfg)
}

// Load the configuration from the FilePath. This overwrites any unsaved changes.
// Returns an error for an unset file path or any file system permission level
// errors.
func (cfg *Configuration) Load() error {
	cfg.locker.Lock()
	defer cfg.locker.Unlock()
	if cfg.FilePath == "" {
		return errors.New("Configuration.FilePath was not set")
	}
	return gonfig.Load(cfg)
}

// Team is the configuration object for a Slack team.
type Team struct {
	BotToken   string        `json:"botToken"`
	Channels   []Channel     `json:"channels"`
	Bot        *slack.Client `json:"-"` // Used only during run time.
	FailedAuth bool          `json:"-"` // Used only during run time. Denotes a failed token.
}

// Channel defines the configuration for a slack channel in a team, including its filters
type Channel struct {
	Name                string   `json:"channelName"`
	MinimumValue        int      `json:"minimumValue"`
	MaximumValue        int      `json:"maximumValue"`
	IncludeCharacters   []string `json:"includeCharacters"`
	IncludeCorporations []string `json:"includeCorporations"`
	IncludeAlliances    []string `json:"includeAlliances"`
	ExcludedShips       []string `json:"excludedShips"`
}
