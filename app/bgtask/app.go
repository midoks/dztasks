package bgtask

import (
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/midoks/dztasks/app/common"
	"github.com/midoks/dztasks/internal/conf"
	"github.com/midoks/dztasks/internal/log"
)

var task *cron.Cron

func clearTask() {
	cronE := task.Entries()
	for _, cron := range cronE {
		task.Remove(cron.ID)
	}
}

func runPluginTask() {
	pluginDir := conf.Plugins.Path
	result := common.PluginList(pluginDir)

	for _, plugin := range result {
		if !plugin.Installed {
			continue
		}

		for _, cron := range plugin.Cron {
			// fmt.Println("cron", plugin.Name, cron.Name, cron.Expr)
			task.AddFunc(cron.Expr, func() {
				msg := ""
				// msg := fmt.Sprintf("正在执行项目[%s][%s][%s]...", plugin.Name, cron.Name, cron.Expr)
				// log.Info(msg)

				runStart := time.Now()

				// fmt.Println(cron.Bin, cron.Args)
				// cronData, err := common.ExecCron(cron.Bin, cron)
				cronData, err := common.ExecPluginCron(plugin, cron)

				if conf.Plugins.ShowError {
					if err != nil {
						log.Info(err.Error())
					}

					if !strings.EqualFold(string(cronData), "") {
						log.Info(string(cronData))
					}
				}

				if err != nil {
					// fmt.Println(err)
					cos := time.Since(runStart)
				msg = fmt.Sprintf("[%s][%s][%s]执行失败,耗时:%s", plugin.Name, cron.Name, cron.Expr, cos)
				log.Info(msg)

					return
				}

				cos := time.Since(runStart)

			msg = fmt.Sprintf("[%s][%s][%s]执行结束,耗时:%s", plugin.Name, cron.Name, cron.Expr, cos)
			log.Info(msg)
			})
		}
	}
}

func InitTask() {
	task = cron.New(cron.WithSeconds())
	runPluginTask()
	task.Start()
}

func ResetTask() {
	log.Info("重置任务")
	clearTask()
	runPluginTask()
	task.Start()
}
