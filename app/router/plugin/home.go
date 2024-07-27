package plugin

import (
	"fmt"
	// "os"
	"io/ioutil"
	// "path/filepath"
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
	c.Data["PageIsPlugin"] = true
	c.Success(PLUGIN_HOME)
}

type PluginTask struct {
	Name string
}

type Plugin struct {
	Name string
	Ps  string
}

func PluginList(c *context.Context) {

	pathdir := conf.Plugins.Path

	files, err := ioutil.ReadDir(pathdir)
	fmt.Println(files,err)

	for _, file := range files {
        fmt.Println(file.Name())

        plugin_name := fmt.Sprintf("%s/%s", pathdir, file.Name())
        fmt.Println(plugin_name)
    }

	c.Ok("ok")
}
