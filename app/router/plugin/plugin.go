package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/midoks/dztasks/app/bgtask"
	"github.com/midoks/dztasks/app/common"
	"github.com/midoks/dztasks/app/context"
	"github.com/midoks/dztasks/app/form"
	"github.com/midoks/dztasks/internal/conf"
	"github.com/midoks/dztasks/internal/log"
	"github.com/midoks/dztasks/internal/tools"
)

const (
	PLUGIN_HOME = "/plugin/index"
	PLUGIN_MENU = "/plugin/menu"
)

func PluginHome(c *context.Context) {
	c.Data["PageIsPlugin"] = true
	c.Success(PLUGIN_HOME)
}

func PluginMenu(c *context.Context, args form.ArgsPluginMenu) {
	c.Data["PageIsPluginMenu"] = true
	c.Data["PageIsPluginMenu_Name"] = args.Name
	c.Data["PageIsPluginMenu_Tag"] = args.Tag

	// fmt.Println(args.Name, args.Tag)
	plugin_dir := conf.Plugins.Path
	list := common.PluginList(plugin_dir)
	c.Data["PluginContent"] = ""

	for _, plugin := range list {
		for _, menu := range plugin.Menu {
			if plugin.Path == args.Name && menu.Tag == args.Tag {
				abs_path := fmt.Sprintf("%s/%s/%s", conf.Plugins.Path, plugin.Path, menu.Path)
				content, _ := ioutil.ReadFile(abs_path)
				c.Data["PluginContent"] = string(content)
			}
		}
	}
	c.Success(PLUGIN_MENU)
}

// 插件列表
func PluginList(c *context.Context) {
	plugin_dir := conf.Plugins.Path
	result := common.PluginList(plugin_dir)
	c.ReturnLayuiJson(0, "ok", len(result), result)
}

// 插件安装
func PluginInstall(c *context.Context, args form.ArgsPluginInstall) {
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
func PluginUninstall(c *context.Context, args form.ArgsPluginUninstall) {
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

// 插件数据请求
func PluginData(c *context.Context, args form.ArgsPluginData) {
	plugin_dir := conf.Plugins.Path
	list := common.PluginList(plugin_dir)

	for _, plugin := range list {
		if plugin.Path == args.Name {
			default_script := fmt.Sprintf("%s/%s/index.py", plugin_dir, plugin.Path)

			script_cmd := make([]string, 0)
			script_cmd = append(script_cmd, default_script)
			script_cmd = append(script_cmd, args.Action)

			script_args := make(map[string]interface{})

			if args.Page > 0 {
				script_args["page"] = args.Page
			}
			if args.Limit > 0 {
				script_args["limit"] = args.Limit
			}

			if !strings.EqualFold(args.Extra, "") {
				script_args["extra"] = args.Extra
			}
			if !strings.EqualFold(args.Args, "") {
				script_args["args"] = args.Args
			}

			post_args, _ := json.Marshal(script_args)

			// fmt.Println("args", args.Args)
			// fmt.Println("extra", args.Extra)

			//展示调用命令
			if conf.Plugins.ShowCmd {
				cmd_args := strings.Join(script_cmd, " ")
				cmd := "python3 " + cmd_args + " '" + string(post_args) + "'"
				// fmt.Println(cmd)
				log.Info(cmd)
			}

			script_cmd = append(script_cmd, string(post_args))
			cmd_data, err := bgtask.ExecInput("python3", script_cmd)

			if err != nil && conf.Plugins.ShowError {
				log.Info(err.Error())
			}

			var plugin_data interface{}
			err = json.Unmarshal(cmd_data, &plugin_data)

			if err != nil && conf.Plugins.ShowError {
				log.Info(err.Error())
			}

			if err == nil {
				c.RenderJson(plugin_data)
				return
			}
			c.Fail(err.Error())
			return
		}
	}
	c.Fail("异常")
}
