package cmd

import (
	"fmt"
	// "net/http"
	// _ "net/http/pprof"
	// "path/filepath"
	// "strings"
	// "time"

	"github.com/urfave/cli"
	"github.com/gin-gonic/gin"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "This command start web services",
	Description: `Start Web Services`,
	Action:      runWeb,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "Custom configuration file path"),
	},
}


func runWeb(c *cli.Context) error {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
	return nil
}
