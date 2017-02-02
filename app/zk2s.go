package app

import (
	"github.com/eveopsec/zk2s/app/config"
	"github.com/eveopsec/zk2s/app/slack"
	"github.com/eveopsec/zk2s/app/tmpl"
	"github.com/vivace-io/evelib/redisq"
)

type App struct {
	slackClient *slack.Client
}

func NewApp(cfg *config.Configuration) (app *App, err error) {
	app = new(App)
	if err = tmpl.Init(cfg.TemplateFile); err != nil {
		return nil, err
	}
	if app.slackClient, err = slack.NewClient(cfg); err != nil {
		return nil, err
	}
	return
}

// Start the application, accepting new killmails from RedisQ and posting them
// to the configured teams and channels. This function does not block.
func (app *App) Start() (err error) {
	client, err := redisq.NewClient(redisq.DefaultOptions())
	if err != nil {
		return err
	}
	client.AddFunc(app.slackClient.Recieve)
	client.Listen()
	return nil
}
