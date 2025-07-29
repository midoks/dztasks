package router

import (
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

// reverseStrings reverses a slice of strings in place
func reverseStrings(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// getLinesText converts reversed log lines to a single string
func getLinesText(lines []string) string {
	reverseStrings(lines)
	return strings.Join(lines, "\n") + "\n"
}

func Log(c *context.Context) {
	lines, err := log.ReverseRead(25)
	if err != nil {
		c.ReturnJson(1, "读取日志失败", "暂无内容")
		return
	}

	content := getLinesText(lines)
	c.ReturnJson(0, "ok", content)
}
