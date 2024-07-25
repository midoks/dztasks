package plugin

import (
	// "fmt"
	// "net/url"

	"github.com/midoks/dztasks/app/context"

	// "github.com/midoks/dztasks/app/form"
	// "github.com/midoks/dztasks/internal/conf"
	// "github.com/midoks/dztasks/internal/log"
	// "github.com/midoks/dztasks/internal/tools"
)

const (
	PLUGIN_PAGE = "/plugin/index"
)

func Home(c *context.Context) {
	c.Success(PLUGIN_PAGE)
}
