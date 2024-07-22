package cmd

import (
	// "fmt"
	// "net/http"
	// _ "net/http/pprof"
	// "path/filepath"
	// "strings"
	// "time"

	"github.com/urfave/cli"

	"github.com/midoks/dztasks/internal/conf"
	"github.com/midoks/dztasks/app"
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
	conf.Init("")

	app.Start(9011)
	return nil
}
