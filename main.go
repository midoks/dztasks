package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/midoks/dztasks/internal/cmd"
	"github.com/midoks/dztasks/internal/conf"
)

const (
	Version = "1.2"
	AppName = "dztasks"
)

func init() {
	conf.App.Version = Version
	conf.App.Name = AppName
}

func main() {
	app := cli.NewApp()
	app.Name = AppName
	app.Version = Version
	app.Usage = "a dztasks service"
	app.Commands = []cli.Command{
		cmd.Web,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
