package plugin

import (
	"fmt"
	// "os"
	"encoding/json"
	"io/ioutil"
	// "path/filepath"
	// "net/url"

	"github.com/midoks/dztasks/app/context"
	"github.com/midoks/dztasks/app/form"
	"github.com/midoks/dztasks/internal/conf"

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
	Path   string `json:"path"`
}

// 插件列表
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
		if err != nil {
			continue
		}

		var p Plugin
		err = json.Unmarshal(content, &p)
		if err != nil {
			continue
		}

		p.Path = file.Name()
		result = append(result, p)
	}
	c.ReturnLayuiJson(0, "ok", len(result), result)
}

// 插件安装
func PluginInstall(c *context.Context, args form.PluginInstall) {
	pathdir := conf.Plugins.Path
	plugin_dir := fmt.Sprintf("%s/%s", pathdir, args.Path)
	plugin_install := fmt.Sprintf("%s/install.lock", plugin_dir)

	if !tools.IsExist(plugin_install) {
		tools.WriteFile(plugin_install, "ok")
	}
	c.Ok("安装成功")
}

// 插件卸载
func PluginUninstall(c *context.Context) {

}
