package router

import (
	// "fmt"
	"strings"

	"github.com/midoks/dztasks/app/context"
	"github.com/midoks/dztasks/internal/log"
)

const (
	HOME = "/console/index"
)

func Home(c *context.Context) {
	c.Data["PageIsHome"] = true
	c.Success(HOME)
}

func Log(c *context.Context) {
	lines, err := log.ReverseRead(25)
	// fmt.Println("log:", err)
	if err == nil {
		content := strings.Join(lines, "\n")
		c.ReturnJson(0, "ok", content)
		return
	}
	c.ReturnJson(0, "ok", "暂无内容")
}
