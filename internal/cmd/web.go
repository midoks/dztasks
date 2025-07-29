package cmd

import (
	"github.com/urfave/cli"

	"github.com/midoks/dztasks/app"
	"github.com/midoks/dztasks/app/bgtask"
	"github.com/midoks/dztasks/app/watch"
	"github.com/midoks/dztasks/internal/conf"
	"github.com/midoks/dztasks/internal/log"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "this command start web services",
	Description: `start web services`,
	Action:      runWeb,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "custom configuration file path"),
	},
}

func runWeb(c *cli.Context) error {
	_ = conf.Init("")
	log.Init()
	go watch.InitWatch(conf.Plugins.Path)
	bgtask.InitTask()

	app.Start(conf.Web.HTTPPort)
	return nil
}
