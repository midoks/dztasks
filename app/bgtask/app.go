package bgtask

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/midoks/dztasks/app/common"
	"github.com/midoks/dztasks/internal/conf"
	"github.com/midoks/dztasks/internal/log"
)

var task *cron.Cron

func execInput(input string) error {
	// Remove the newline character.
	input = strings.TrimSuffix(input, "\n")

	// Prepare the command to execute.
	cmd := exec.Command(input)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	fmt.Println(os.Stdout, os.Stderr)
	// Execute the command and return the error.
	return cmd.Run()
}

func InitTask() {

	plugin_dir := conf.Plugins.Path
	result := common.PluginList(plugin_dir)

	// fmt.Println(result)
	task = cron.New()

	for _, plugin := range result {

		if !plugin.Installed {
			continue
		}

		// fmt.Println(len(plugin.Cron))
		for _, cron := range plugin.Cron {
			// fmt.Println(cron)
			task.AddFunc(cron.Expr, func() {
				msg := fmt.Sprintf("正在执行项目[%s][%s][%s]...", plugin.Name, cron.Name, cron.Expr)
				fmt.Println(msg)
				log.Info(msg)

				run_start := time.Now()

				execInput(cron.Cmd)

				cos := time.Since(run_start)

				msg = fmt.Sprintf("[%s][%s][%s]执行结束,耗时:%s", plugin.Name, cron.Name, cron.Expr, cos)
				fmt.Println(msg)
				log.Info(msg)
			})
		}
	}
	// task.AddFunc("@every 5s", func() { fmt.Println("Every hour thirty") })
	task.Start()
}
