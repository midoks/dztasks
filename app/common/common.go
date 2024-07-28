package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/midoks/dztasks/internal/tools"
)

type PluginCron struct {
	Name string `json:"name"`
	Expr string `json:"expr"`
	Cmd  string `json:"cmd"`
}

type Plugin struct {
	Name      string       `json:"name"`
	Ps        string       `json:"ps"`
	Author    string       `json:"author"`
	Path      string       `json:"path"`
	Installed bool         `json:"installed"`
	Cron      []PluginCron `json:"cron"`
}

func GetPluginInstallLock(plugin_name string) string {
	lock := fmt.Sprintf("%s/install.lock", plugin_name)
	return lock
}

func PluginList(plugin_dir string) []Plugin {
	files, _ := ioutil.ReadDir(plugin_dir)
	result := make([]Plugin, 0)

	for _, file := range files {

		plugin_name := fmt.Sprintf("%s/%s", plugin_dir, file.Name())
		plugin_info := fmt.Sprintf("%s/info.json", plugin_name)
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

		plugin_lock := GetPluginInstallLock(plugin_name)
		if tools.IsExist(plugin_lock) {
			p.Installed = true
		}
		result = append(result, p)
	}
	return result
}
