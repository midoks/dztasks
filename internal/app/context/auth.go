package context

import (
	// "fmt"
	"net/url"

	
	"gopkg.in/macaron.v1"

	"github.com/midoks/dztasks/internal/conf"
)

type ToggleOptions struct {
	SignInRequired  bool
	SignOutRequired bool
	AdminRequired   bool
	DisableCSRF     bool
}

func Toggle(options *ToggleOptions) macaron.Handler {

	return func(c *Context) {
		if !conf.Security.InstallLock {
			// fmt.Println("Toggle:not install")
			c.RedirectSubpath("/install")
			return
		}

		if options.SignInRequired {
			if !c.IsLogged {
				c.SetCookie("redirect_to", url.QueryEscape(conf.Web.Subpath+c.Req.RequestURI), 0, conf.Web.Subpath)
				c.RedirectSubpath("/login")
				return
			}
		}
	}
}
