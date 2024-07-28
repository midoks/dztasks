package plugin

import (
	"fmt"
	"os"
	"time"

	"github.com/midoks/dztasks/app/bgtask"
	"github.com/midoks/dztasks/app/common"
	"github.com/midoks/dztasks/app/context"
	"github.com/midoks/dztasks/app/form"
	"github.com/midoks/dztasks/internal/conf"
	"github.com/midoks/dztasks/internal/tools"
)

const (
	PLUGIN_HOME = "/plugin/index"
)

func PluginHome(c *context.Context) {
	c.Data["PageIsPlugin"] = true
	c.Success(PLUGIN_HOME)
}

// 插件列表
func PluginList(c *context.Context) {
	plugin_dir := conf.Plugins.Path
	result := common.PluginList(plugin_dir)
	c.ReturnLayuiJson(0, "ok", len(result), result)
}

// 插件安装
func PluginInstall(c *context.Context, args form.PluginInstall) {
	plugin_dir := conf.Plugins.Path
	plugin_name := fmt.Sprintf("%s/%s", plugin_dir, args.Path)
	plugin_lock := common.GetPluginInstallLock(plugin_name)
	if !tools.IsExist(plugin_lock) {
		tools.WriteFile(plugin_lock, "ok")
		time.Sleep(2)
		bgtask.ResetTask()
		c.Ok("安装成功")
		return
	}
	c.Ok("已经安装成功")
}

// 插件卸载
func PluginUninstall(c *context.Context, args form.PluginUninstall) {
	plugin_dir := conf.Plugins.Path
	plugin_name := fmt.Sprintf("%s/%s", plugin_dir, args.Path)
	plugin_lock := common.GetPluginInstallLock(plugin_name)
	if tools.IsExist(plugin_lock) {
		os.Remove(plugin_lock)
		time.Sleep(2)
		bgtask.ResetTask()
		c.Ok("卸载成功")
		return
	}
	c.Ok("已经卸载成功")
}
