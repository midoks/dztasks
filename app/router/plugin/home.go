package plugin

import (
	"fmt"
	// "net/url"

	"github.com/midoks/dztasks/app/context"
	"github.com/midoks/dztasks/internal/conf"
	// "github.com/midoks/dztasks/app/form"
	
	// "github.com/midoks/dztasks/internal/log"
	// "github.com/midoks/dztasks/internal/tools"
)

const (
	PLUGIN_HOME = "/plugin/index"
)

func PluginHome(c *context.Context) {
	c.Success(PLUGIN_HOME)
}

func PluginList(c *context.Context) {

	pathdir := conf.Plugins.Path

	fmt.Println(pathdir)



	c.Ok("ok")
}
