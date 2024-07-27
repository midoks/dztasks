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
	Name   string `json:"name"`
	Ps     string `json:"ps"`
	Author string `json:"author"`
}

// json api common data
type LayuiData struct {
	Code  int64       `json:"code"`
	Count int64       `json:"count"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
}

func PluginList(c *context.Context) {

	pathdir := conf.Plugins.Path

	files, _ := ioutil.ReadDir(pathdir)
	result := make([]Plugin, 0)

	for _, file := range files {

		plugin_dir := fmt.Sprintf("%s/%s", pathdir, file.Name())
		plugin_info := fmt.Sprintf("%s/info.json", plugin_dir)
		if !tools.IsExist(plugin_info) {
			continue
		}

		content, err := ioutil.ReadFile(plugin_info)
		fmt.Println(err)
		if err != nil {
			continue
		}

		var p Plugin
		err = json.Unmarshal(content, &p)
		fmt.Println(err)
		fmt.Println(p)
		result = append(result, p)
	}

	fmt.Println(len(result))

	data := LayuiData{Code: 0, Msg: "ok", Count: 1, Data: result}

	c.RenderJson(data)
}
