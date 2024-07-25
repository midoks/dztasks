package user

import (
	"fmt"
	// "net/url"

	"github.com/midoks/dztasks/app/context"
	"github.com/midoks/dztasks/app/form"
	"github.com/midoks/dztasks/internal/conf"
	"github.com/midoks/dztasks/internal/log"
	// "github.com/midoks/dztasks/internal/tools"
)

const (
	LOGIN = "/user/login"
)

// AutoLogin reads cookie and try to auto-login.
func AutoLogin(c *context.Context) (bool, error) {

	uname := c.GetCookie(conf.Security.CookieUsername)
	if len(uname) == 0 {
		return false, nil
	}

	isSucceed := false
	defer func() {
		if !isSucceed {
			log.Infof("auto-login cookie cleared: %s", uname)
			c.SetCookie(conf.Security.CookieUsername, "", -1, conf.Web.Subpath)
			c.SetCookie(conf.Security.CookieRememberName, "", -1, conf.Web.Subpath)
			c.SetCookie(conf.Security.LoginStatusCookieName, "", -1, conf.Web.Subpath)
		}
	}()

	name := c.Session.Get("name")
	fmt.Println("login:",name)
	// if uid != nil {
	// 	u, err := db.UserGetById(uid.(int64))
	// 	if err != nil {
	// 		return false, nil
	// 	}

	// 	if val, ok := c.GetSuperSecureCookie(u.Salt+u.Password, conf.Security.CookieRememberName); !ok || val != u.Name {
	// 		return false, nil
	// 	}

	// 	isSucceed = true
	// 	c.Session.Set("uid", u.Id)
	// 	c.Session.Set("uname", u.Name)
	// 	c.SetCookie(conf.Session.CSRFCookieName, "", -1, conf.Web.Subpath)
	// 	if conf.Security.EnableLoginStatusCookie {
	// 		c.SetCookie(conf.Security.LoginStatusCookieName, "true", 0, conf.Web.Subpath)
	// 	}
	// }

	return true, nil
}

func Login(c *context.Context) {
	c.Success(LOGIN)
}

func LoginPost(c *context.Context, f form.SignIn) {

	fmt.Println("api",f.Username, f.Password)
	fmt.Println("conf",conf.Admin.User, conf.Admin.Pass)


	if (conf.Admin.User == f.Username && conf.Admin.Pass == f.Password){
		name := conf.Admin.User
		pass := conf.Admin.Pass
		if f.Remember {
			days := 86400 * conf.Security.LoginRememberDays
			c.SetCookie(conf.Security.CookieUsername, name, days, conf.Web.Subpath, "", conf.Security.CookieSecure, true)
			c.SetSuperSecureCookie(pass, conf.Security.CookieRememberName, name, days, conf.Web.Subpath, "", conf.Security.CookieSecure, true)
		}

		c.Session.Set("name", name)
		c.Session.Set("login", true)

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
	_ = c.Session.Flush()
	_ = c.Session.Destory(c.Context)
	c.SetCookie(conf.Security.CookieUsername, "", -1, conf.Web.Subpath)
	c.SetCookie(conf.Security.CookieRememberName, "", -1, conf.Web.Subpath)
	c.SetCookie(conf.Session.CSRFCookieName, "", -1, conf.Web.Subpath)
	c.RedirectSubpath("/")
}
