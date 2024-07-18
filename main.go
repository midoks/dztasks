package main

import (
	// "fmt"
	"log"
	"os"

	"dztasks/internal/cmd"

	"github.com/urfave/cli"
)


const Version = "0.0.3"
const AppName = "dztasks"

func init() {
	// conf.App.Version = Version
	// conf.App.Name = AppName
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
		log.Println("Failed to start application: ", err)
	}

}