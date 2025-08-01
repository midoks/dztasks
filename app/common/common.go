package common

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/midoks/dztasks/internal/conf"
	"github.com/midoks/dztasks/internal/log"
	"github.com/midoks/dztasks/internal/tools"
)

type PluginCron struct {
	Bin string `json:"bin"`
	Env string `json:"env"`
	Dir string `json:"dir"`

	Name string   `json:"name"`
	Expr string   `json:"expr"`
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

type PluginMenu struct {
	Title string   `json:"title"`
	Path  string   `json:"path"`
	Tag   string   `json:"tag"`
	CSS   []string `json:"css"`
	Js    []string `json:"js"`
}

type Plugin struct {
	Bin string `json:"bin"`
	Env string `json:"env"`
	Dir string `json:"dir"`

	Name      string `json:"name"`
	Ps        string `json:"ps"`
	Author    string `json:"author"`
	Path      string `json:"path"`
	Index     string `json:"index"`
	Icon      string `json:"icon"`
	Installed bool   `json:"installed"`

	Cron []PluginCron `json:"cron"`
	Menu []PluginMenu `json:"menu"`
}

func ExecPluginCron(plugin Plugin, cron PluginCron) ([]byte, error) {
	bin := cron.Bin
	if strings.EqualFold(bin, "") {
		if !strings.EqualFold(plugin.Bin, "") {
			bin = plugin.Bin
		}
	}

	if conf.Plugins.ShowCmd {
		log.Info(bin + " " + strings.Join(cron.Args, " "))
	}

	// Prepare the command to execute.
	cmd := exec.Command(bin, cron.Args...)
	if strings.EqualFold(cron.Dir, "") {
		if !strings.HasSuffix(cron.Dir, "/") {
			cmd.Dir = fmt.Sprintf("%s/%s", conf.Plugins.Path, plugin.Dir)
		} else {
			cmd.Dir = cron.Dir
		}
	} else if !strings.EqualFold(plugin.Dir, "") {
		if !strings.HasSuffix(plugin.Dir, "/") {
			cmd.Dir = fmt.Sprintf("%s/%s", conf.Plugins.Path, plugin.Dir)
		} else {
			cmd.Dir = plugin.Dir
		}
	}

	env := os.Environ()
	if !strings.EqualFold(cron.Env, "") {
		cronEnv := make([]string, 0)
		cronEnv = append(cronEnv, cron.Env)
		env = cronEnv
		// fmt.Println(cmd.Env)
	} else if !strings.EqualFold(plugin.Env, "") {
		env = append(env, plugin.Env)
	}
	cmd.Env = env

	// fmt.Println(os.Stdout, os.Stderr)
	return cmd.CombinedOutput()
}

func ExecCron(bin string, cron PluginCron) ([]byte, error) {
	cmd := exec.Command(bin, cron.Args...)

	if !strings.EqualFold(cron.Dir, "") {
		cmd.Dir = cron.Dir
		if !strings.HasSuffix(cron.Dir, "/") {
			pdir := fmt.Sprintf("%s/%s", conf.Plugins.Path, cron.Dir)
			cmd.Dir = pdir
		} else {
			cmd.Dir = cron.Dir
		}
	}

	if !strings.EqualFold(cron.Env, "") {
		env := make([]string, 0)
		cmd.Env = append(env, cron.Env)
	} else {
		cmd.Env = append(os.Environ())
	}

	// fmt.Println(os.Stdout, os.Stderr)
	return cmd.CombinedOutput()
}

func ExecPluginCmd(plugin Plugin, args []string) ([]byte, error) {
	// Remove the newline character.
	// input = strings.TrimSuffix(input, "\n")

	// Prepare the command to execute.
	cmd := exec.Command(plugin.Bin, args...)
	cmd.Env = append(os.Environ())

	if !strings.EqualFold(plugin.Dir, "") {
		if !strings.HasSuffix(plugin.Dir, "/") {
			pdir := fmt.Sprintf("%s/%s", conf.Plugins.Path, plugin.Dir)
			cmd.Dir = pdir
		} else {
			cmd.Dir = plugin.Dir
		}
	}

	// fmt.Println(os.Stdout, os.Stderr)
	// Execute the command and return the error.
	return cmd.CombinedOutput()
}

func ExecInput(bin string, args []string) ([]byte, error) {
	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ())

	// fmt.Println(os.Stdout, os.Stderr)
	return cmd.CombinedOutput()
}

func GetPluginInstallLock(pluginName string) string {
	lock := fmt.Sprintf("%s/install.lock", pluginName)
	return lock
}

func PluginList(pluginDir string) []Plugin {
	files, _ := os.ReadDir(pluginDir)
	result := make([]Plugin, 0)

	for _, file := range files {

		pluginName := fmt.Sprintf("%s/%s", pluginDir, file.Name())
		pluginInfo := fmt.Sprintf("%s/info.json", pluginName)
		if !tools.IsExist(pluginInfo) {
			continue
		}

		content, err := os.ReadFile(pluginInfo)
		if err != nil {
			continue
		}

		var p Plugin
		err = json.Unmarshal(content, &p)
		if err != nil {
			fmt.Println("err:", err)
			continue
		}
		p.Path = file.Name()

		if p.Icon == "" {
			p.Icon = "layui-icon-tree"
		}

		p.Dir = file.Name()
		pluginLock := GetPluginInstallLock(pluginName)
		if tools.IsExist(pluginLock) {
			p.Installed = true
		}
		result = append(result, p)
	}
	return result
}
