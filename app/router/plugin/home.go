package plugin

import (
	"fmt"
	// "os"
	"encoding/json"
	"io/ioutil"
	// "path/filepath"
	// "net/url"

	"github.com/midoks/dztasks/app/context"
	"github.com/midoks/dztasks/internal/conf"
	// "github.com/midoks/dztasks/app/form"

	// "github.com/midoks/dztasks/internal/log"
	"github.com/midoks/dztasks/internal/tools"
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
	Name   string
	Ps     string
	Author string
}

func PluginList(c *context.Context) {

	pathdir := conf.Plugins.Path

	files, _ := ioutil.ReadDir(pathdir)

	for _, file := range files {
		fmt.Println(file.Name())

		plugin_name := fmt.Sprintf("%s/%s", pathdir, file.Name())
		plugin_info := fmt.Sprintf("%s/info.json", plugin_name)
		fmt.Println(plugin_name, plugin_info)

		if !tools.IsExist(plugin_info) {
			continue
		}

		content, err := ioutil.ReadFile(plugin_info)
		fmt.Println(err)
		if err != nil {
			continue
		}

		var payload Plugin
		err = json.Unmarshal(content, &payload)
		fmt.Println(err)

	}

	c.Ok("ok")
}
