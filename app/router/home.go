package router

import (
	// "fmt"
	// "sort"
	// "strings"

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

func Reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func ReverseCopy[T any](original []T) []T {
	reversed := make([]T, len(original))
	copy(reversed, original)
	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}
	return reversed
}

// func reverse(s []string) {
// 	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
// 		s[i], s[j] = s[j], s[i]
// 	}
// }

func getLinesText(lines []string) string {
	content_line := ""
	Reverse(lines)
	for _, text := range lines {
		content_line += text + "\n"
	}
	return content_line
}

func Log(c *context.Context) {
	lines, err := log.ReverseRead(25)
	// fmt.Println("log:", err)
	if err == nil {
		// content := strings.Join(lines, "\n")
		content := getLinesText(lines)
		c.ReturnJson(0, "ok", content)
		return
	}
	c.ReturnJson(0, "ok", "暂无内容")
}
