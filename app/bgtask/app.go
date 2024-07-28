package bgtask

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

var task *cron.Cron

func InitTask() {
	task = cron.New()

	task.AddFunc("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
	task.AddFunc("@hourly", func() { fmt.Println("Every hour") })
	task.AddFunc("@every 5s", func() { fmt.Println("Every hour thirty") })
	task.Start()
}
