package router

import (
	"github.com/midoks/dztasks/app/context"
)

const (
	HOME = "/console/index"
)

func Home(c *context.Context) {
	c.Data["PageIsHome"] = true
	c.Success(HOME)
}
