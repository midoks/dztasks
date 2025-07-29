package plugin

import (
	"encoding/json"
	"fmt"
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
	PluginHomePath    = "/plugin/index"
	PluginMenuPath    = "/plugin/menu"
	PluginContentPath = "/plugin/content"
)

func PluginHome(c *context.Context) {
	c.Data["PageIsPlugin"] = true
	c.Success(PluginHomePath)
}

func PluginMenu(c *context.Context, args form.ArgsPluginMenu) {
	c.Data["PageIsPluginMenu"] = true
	c.Data["PageIsPluginMenu_Name"] = args.Name
	c.Data["PageIsPluginMenu_Tag"] = args.Tag

	// fmt.Println(args.Name, args.Tag)
	pluginDir := conf.Plugins.Path
	list := common.PluginList(pluginDir)
	c.Data["PluginContent"] = ""

	for _, plugin := range list {
		for _, menu := range plugin.Menu {
			if plugin.Path == args.Name && menu.Tag == args.Tag {
				absPath := fmt.Sprintf("%s/%s/%s", conf.Plugins.Path, plugin.Path, menu.Path)
				content, _ := os.ReadFile(absPath)
				c.Data["PluginContent"] = string(content)
			}
		}
	}
	c.Success(PluginMenuPath)
}

// 插件列表
func PluginList(c *context.Context) {
	pluginDir := conf.Plugins.Path
	result := common.PluginList(pluginDir)
	c.ReturnLayuiJSON(0, "ok", len(result), result)
}

// 插件安装
func PluginInstall(c *context.Context, args form.ArgsPluginInstall) {
	pluginDir := conf.Plugins.Path
	pluginName := fmt.Sprintf("%s/%s", pluginDir, args.Path)
	pluginLock := common.GetPluginInstallLock(pluginName)
	if !tools.IsExist(pluginLock) {
		tools.WriteFile(pluginLock, "ok")
		time.Sleep(2 * time.Second)
		bgtask.ResetTask()
		c.Ok("安装成功")
		return
	}
	c.Ok("已经安装成功")
}

// 插件卸载
func PluginUninstall(c *context.Context, args form.ArgsPluginUninstall) {
	pluginDir := conf.Plugins.Path
	pluginName := fmt.Sprintf("%s/%s", pluginDir, args.Path)
	pluginLock := common.GetPluginInstallLock(pluginName)
	if tools.IsExist(pluginLock) {
		os.Remove(pluginLock)
		time.Sleep(2 * time.Second)
		bgtask.ResetTask()
		c.Ok("卸载成功")
		return
	}
	c.Ok("已经卸载成功")
}

// 插件数据请求
func PluginPage(c *context.Context, args form.ArgsPluginPage) {
	pluginDir := conf.Plugins.Path
	list := common.PluginList(pluginDir)
	c.Data["PluginContent"] = ""
	for _, plugin := range list {
		if plugin.Path == args.Name {
			// fmt.Println(args)
			absPath := fmt.Sprintf("%s/%s/%s", pluginDir, plugin.Path, args.Page)

			content, err := tools.ReadFile(absPath)
			if err == nil {
				c.Data["PluginContent"] = content
				c.Success(PluginContentPath)
				return
			}
			c.Data["PluginContent"] = err.Error()
		}
	}
	c.Success(PluginContentPath)
}

// 插件数据请求
func PluginData(c *context.Context, args form.ArgsPluginData) {
	pluginDir := conf.Plugins.Path
	list := common.PluginList(pluginDir)

	for _, plugin := range list {
		if plugin.Path == args.Name {
			if plugin.Index == "" {
				plugin.Index = "index.py"
			}
			if plugin.Bin == "" {
				plugin.Bin = "python3"
			}

			defaultScript := fmt.Sprintf("%s/%s/%s", pluginDir, plugin.Path, plugin.Index)
			if !strings.EqualFold(plugin.Dir, "") {
				defaultScript = plugin.Index
			}

			fmt.Println(defaultScript)

			scriptCmd := make([]string, 0)
			scriptCmd = append(scriptCmd, defaultScript)
			scriptCmd = append(scriptCmd, args.Action)

			scriptArgs := make(map[string]interface{})

			if args.Page > 0 {
				scriptArgs["page"] = args.Page
			}
			if args.Limit > 0 {
				scriptArgs["limit"] = args.Limit
			}

			if !strings.EqualFold(args.Extra, "") {
				scriptArgs["extra"] = args.Extra
			}
			if !strings.EqualFold(args.Args, "") {
				scriptArgs["args"] = args.Args
			}

			postArgs, _ := json.Marshal(scriptArgs)

			// 展示调用命令
			if conf.Plugins.ShowCmd {
				cmdArgs := strings.Join(scriptCmd, " ")
				cmd := "[CMD]" + plugin.Bin + " " + cmdArgs + " '" + string(postArgs) + "'"
				log.Info(cmd)
			}

			scriptCmd = append(scriptCmd, string(postArgs))
			// cmd_data, err := common.ExecInput(plugin.Bin, scriptCmd)
			cmdData, err := common.ExecPluginCmd(plugin, scriptCmd)

			if err != nil && conf.Plugins.ShowError {
				log.Info(err.Error())
			}

			if !strings.EqualFold(string(cmdData), "") {
				log.Info(string(cmdData) + "\n")
			}

			var pluginData interface{}
			err = json.Unmarshal(cmdData, &pluginData)

			if err != nil && conf.Plugins.ShowError {
				log.Info(err.Error())
			}

			if err == nil {
				c.RenderJSON(pluginData)
				return
			}
			c.Fail(err.Error())
			return
		}
	}
	c.Fail("异常")
}

// 插件数据请求
func PluginFile(c *context.Context, args form.ArgsPluginFile) {
	pluginDir := conf.Plugins.Path
	list := common.PluginList(pluginDir)

	for _, plugin := range list {
		if plugin.Path == args.Name {
			var defaultScript string
			if !strings.EqualFold(plugin.Dir, "") {
				// 如果插件有自定义目录，使用自定义目录
				defaultScript = fmt.Sprintf("%s/%s", plugin.Dir, args.File)
			} else {
				// 否则使用默认的插件目录
				defaultScript = fmt.Sprintf("%s/%s/%s", pluginDir, plugin.Path, args.File)
			}
			content, err := tools.ReadFileByte(defaultScript)
			if err == nil {
				c.RawData(200, content)
				return
			}
			c.Fail(err.Error())
		}
	}
	c.PlainText(400, []byte("异常"))
}
