package bgtask

import (
	"fmt"
	"os"
	"os/exec"
	// "strings"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/midoks/dztasks/app/common"
	"github.com/midoks/dztasks/internal/conf"
	"github.com/midoks/dztasks/internal/log"
)

var task *cron.Cron

func ExecInput(bin string, args []string) ([]byte, error) {
	// Remove the newline character.
	// input = strings.TrimSuffix(input, "\n")

	// Prepare the command to execute.
	cmd := exec.Command(bin, args...)

	cmd.Env = append(os.Environ())

	// fmt.Println(os.Stdout, os.Stderr)
	// Execute the command and return the error.
	return cmd.CombinedOutput()
}

func clearTask() {
	cronE := task.Entries()
	for _, cron := range cronE {
		task.Remove(cron.ID)
	}
}

func runPluginTask() {
	plugin_dir := conf.Plugins.Path
	result := common.PluginList(plugin_dir)

	for _, plugin := range result {

		if !plugin.Installed {
			continue
		}

		for _, cron := range plugin.Cron {
			// fmt.Println("cron", plugin.Name, cron.Name, cron.Expr)
			task.AddFunc(cron.Expr, func() {
				msg := ""
				// msg := fmt.Sprintf("正在执行项目[%s][%s][%s]...", plugin.Name, cron.Name, cron.Expr)
				// fmt.Println(msg)
				// log.Info(msg)

				run_start := time.Now()

				cronData, err := ExecInput(cron.Bin, cron.Args)

				if conf.Plugins.ShowCmd {
					log.Info(string(cronData))
				}

				if err != nil {
					// fmt.Println(err)
					cos := time.Since(run_start)
					msg = fmt.Sprintf("[%s][%s][%s]执行失败,耗时:%s", plugin.Name, cron.Name, cron.Expr, cos)
					// fmt.Println(msg)
					log.Info(msg)
					log.Info(err.Error())
					return
				}

				cos := time.Since(run_start)

				msg = fmt.Sprintf("[%s][%s][%s]执行结束,耗时:%s", plugin.Name, cron.Name, cron.Expr, cos)
				// fmt.Println(msg)
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
