package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	Name      string `json:"name"`
	Ps        string `json:"ps"`
	Author    string `json:"author"`
	Path      string `json:"path"`
	Installed bool   `json:"installed"`
}

func getPluginInstallLock(path string) string {
	pathdir := conf.Plugins.Path
	plugin_dir := fmt.Sprintf("%s/%s", pathdir, path)
	lock := fmt.Sprintf("%s/install.lock", plugin_dir)
	return lock
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

		plugin_lock := getPluginInstallLock(file.Name())
		if tools.IsExist(plugin_lock) {
			p.Installed = true
		}

		p.Path = file.Name()
		result = append(result, p)
	}
	c.ReturnLayuiJson(0, "ok", len(result), result)
}

// 插件安装
func PluginInstall(c *context.Context, args form.PluginInstall) {
	plugin_lock := getPluginInstallLock(args.Path)
	if !tools.IsExist(plugin_lock) {
		tools.WriteFile(plugin_lock, "ok")
		c.Ok("安装成功")
		return
	}
	c.Ok("已经安装成功")
}

// 插件卸载
func PluginUninstall(c *context.Context, args form.PluginUninstall) {
	plugin_lock := getPluginInstallLock(args.Path)
	if tools.IsExist(plugin_lock) {
		os.Remove(plugin_lock)
		c.Ok("卸载成功")
		return
	}
	c.Ok("已经卸载成功")
}
