package user

import (
	// "fmt"

	"github.com/midoks/dztasks/app/context"
	"github.com/midoks/dztasks/app/form"
	"github.com/midoks/dztasks/internal/conf"
)

const (
	LOGIN = "/user/login"
)

func Login(c *context.Context) {
	isLogged := c.Session.Get("login")
	if isLogged == true {
		c.RedirectSubpath("/")
	}
	c.Success(LOGIN)
}

func DoLogin(c *context.Context, f form.SignIn) {
	// fmt.Println(f)
	if conf.Admin.User == f.Username && conf.Admin.Pass == f.Password {
		name := conf.Admin.User
		pass := conf.Admin.Pass
		if f.Remember {
			days := 86400 * conf.Security.LoginRememberDays
			c.SetCookie(conf.Security.CookieUsername, name, days, conf.Web.Subpath, "", conf.Security.CookieSecure, true)
			c.SetSuperSecureCookie(pass, conf.Security.CookieRememberName, name, days, conf.Web.Subpath, "", conf.Security.CookieSecure, true)
		}

		_ = c.Session.Set("name", name)
		_ = c.Session.Set("login", true)

		c.SetCookie(conf.Session.CSRFCookieName, "", -1, conf.Web.Subpath)
		if conf.Security.EnableLoginStatusCookie {
			c.SetCookie(conf.Security.LoginStatusCookieName, "true", 0, conf.Web.Subpath)
		}

		c.Redirect("/")
	} else {
		c.Fail("认证失败")
	}
}

func LoginPost(c *context.Context, f form.SignIn) {
	// fmt.Println(f)
	if conf.Admin.User == f.Username && conf.Admin.Pass == f.Password {
		name := conf.Admin.User
		pass := conf.Admin.Pass
		if f.Remember {
			days := 86400 * conf.Security.LoginRememberDays
			c.SetCookie(conf.Security.CookieUsername, name, days, conf.Web.Subpath, "", conf.Security.CookieSecure, true)
			c.SetSuperSecureCookie(pass, conf.Security.CookieRememberName, name, days, conf.Web.Subpath, "", conf.Security.CookieSecure, true)
		}

		_ = c.Session.Set("name", name)
		_ = c.Session.Set("login", true)

		c.SetCookie(conf.Session.CSRFCookieName, "", -1, conf.Web.Subpath)
		if conf.Security.EnableLoginStatusCookie {
			c.SetCookie(conf.Security.LoginStatusCookieName, "true", 0, conf.Web.Subpath)
		}

		c.Ok("登录成功")
	} else {
		c.Fail("认证失败")
	}
}

func SignOut(c *context.Context) {
	// 清除会话中的登录状态
	_ = c.Session.Delete("login")
	_ = c.Session.Delete("name")
	_ = c.Session.Flush()
	_ = c.Session.Release()

	// 清除所有相关的 cookies，包括会话 cookie
	c.SetCookie(conf.Session.CookieName, "", -1, conf.Web.Subpath)
	c.SetCookie(conf.Security.CookieUsername, "", -1, conf.Web.Subpath)
	c.SetCookie(conf.Security.CookieRememberName, "", -1, conf.Web.Subpath)
	c.SetCookie(conf.Session.CSRFCookieName, "", -1, conf.Web.Subpath)
	if conf.Security.EnableLoginStatusCookie {
		c.SetCookie(conf.Security.LoginStatusCookieName, "", -1, conf.Web.Subpath)
	}

	c.RedirectSubpath("/")
}
