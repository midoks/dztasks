package user

import (
	"fmt"
	"net/url"

	"github.com/go-macaron/captcha"

	"github.com/midoks/dztasks/app/context"
	"github.com/midoks/dztasks/app/form"
	"github.com/midoks/dztasks/internal/conf"
	"github.com/midoks/dztasks/internal/log"
	"github.com/midoks/dztasks/internal/tools"
)

const (
	LOGIN                    = "/user/login"
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

	// uid := c.Session.Get("uid")
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


	// // Check auto-login
	// isSucceed, err := AutoLogin(c)
	// if err != nil {
	// 	c.Error(err, "auto login")
	// 	return
	// }

	// redirectTo := c.Query("redirect_to")
	// if len(redirectTo) > 0 {
	// 	c.SetCookie("redirect_to", redirectTo, 0, conf.Web.Subpath)
	// } else {
	// 	redirectTo, _ = url.QueryUnescape(c.GetCookie("redirect_to"))
	// }

	// if isSucceed {
	// 	if tools.IsSameSiteURLPath(redirectTo) {
	// 		c.Redirect(redirectTo)
	// 	} else {
	// 		c.RedirectSubpath("/")
	// 	}
	// 	c.SetCookie("redirect_to", "", -1, conf.Web.Subpath)
	// 	return
	// }

	c.Success(LOGIN)
}

func LoginPost(c *context.Context, f form.SignIn) {

	fmt.Println(f.UserName, f.Password)

	// loginBool, uid := db.LoginByUserPassword(f.UserName, f.Password)

	// if !loginBool {
	// 	c.FormErr("UserName", "Password")
	// 	c.RenderWithErr(c.Tr("form.username_password_incorrect"), LOGIN, &f)
	// 	return
	// }

	if c.HasError() {
		c.Success(LOGIN)
		return
	}

	// u, err := db.UserGetByName(f.UserName)
	// if err != nil {
	// 	c.FormErr("UserName", "Password")
	// 	c.RenderWithErr(c.Tr("form.username_password_incorrect"), LOGIN, &f)
	// 	return
	// }

	// if f.Remember {
	// 	days := 86400 * conf.Security.LoginRememberDays
	// 	c.SetCookie(conf.Security.CookieUsername, u.Name, days, conf.Web.Subpath, "", conf.Security.CookieSecure, true)
	// 	c.SetSuperSecureCookie(u.Salt+u.Password, conf.Security.CookieRememberName, u.Name, days, conf.Web.Subpath, "", conf.Security.CookieSecure, true)
	// }

	// session.Start(c)
	// c.Session.Set("uid", uid)
	// c.Session.Set("uname", f.UserName)

	// Clear whatever CSRF has right now, force to generate a new one
	c.SetCookie(conf.Session.CSRFCookieName, "", -1, conf.Web.Subpath)
	if conf.Security.EnableLoginStatusCookie {
		c.SetCookie(conf.Security.LoginStatusCookieName, "true", 0, conf.Web.Subpath)
	}

	redirectTo, _ := url.QueryUnescape(c.GetCookie("redirect_to"))
	c.SetCookie("redirect_to", "", -1, conf.Web.Subpath)

	if tools.IsSameSiteURLPath(redirectTo) {
		c.Redirect(redirectTo)
		return
	}

	c.RedirectSubpath("/")
}

func SignOut(c *context.Context) {
	_ = c.Session.Flush()
	_ = c.Session.Destory(c.Context)
	c.SetCookie(conf.Security.CookieUsername, "", -1, conf.Web.Subpath)
	c.SetCookie(conf.Security.CookieRememberName, "", -1, conf.Web.Subpath)
	c.SetCookie(conf.Session.CSRFCookieName, "", -1, conf.Web.Subpath)
	c.RedirectSubpath("/")
}


func SignUp(c *context.Context) {
	c.Data["EnableCaptcha"] = conf.Auth.EnableRegistrationCaptcha

	if conf.Auth.DisableRegistration {
		c.Data["DisableRegistration"] = true
		c.Success(LOGIN)
		return
	}

	c.Success(LOGIN)
}

func SignUpPost(c *context.Context, cpt *captcha.Captcha, f form.Register) {

	c.Data["EnableCaptcha"] = conf.Auth.EnableRegistrationCaptcha

	if conf.Auth.DisableRegistration {
		c.Status(403)
		return
	}

	if c.HasError() {
		c.Success(LOGIN)
		return
	}

	if conf.Auth.EnableRegistrationCaptcha && !cpt.VerifyReq(c.Req) {
		c.FormErr("Captcha")
		c.RenderWithErr("验证码错误!", LOGIN, &f)
		return
	}

	// if f.Password != f.Retype {
	// 	c.FormErr("Password")
	// 	c.RenderWithErr(c.Tr("form.password_not_match"), SIGNUP, &f)
	// 	return
	// }

	// u := &db.User{
	// 	Name:     f.UserName,
	// 	Password: f.Password,
	// 	IsActive: false,
	// }


	// if err := db.CreateUser(u); err != nil {
	// 	c.Error(err, "create user")
	// 	return
	// }

	// log.Debugf("Account created: %s", u.Name)

	// Auto-set admin for the only user.
	// if db.UsersCount() == 1 {
	// 	u.IsAdmin = true
	// 	u.IsActive = true
	// 	if err := db.UpdateUser(u); err != nil {
	// 		c.Error(err, "update user")
	// 		return
	// 	}
	// }

	// Send confirmation email.
	// if conf.Auth.RequireEmailConfirmation {
	// 	email.SendActivateAccountMail(c.Context, db.NewMailerUser(u))
	// 	c.Data["IsSendRegisterMail"] = true
	// 	c.Data["Email"] = u.Email
	// 	c.Data["Hours"] = conf.Auth.ActivateCodeLives / 60
	// 	c.Success(ACTIVATE)

	// 	if err := c.Cache.Put(u.MailResendCacheKey(), 1, 180); err != nil {
	// 		log.Error("Failed to put cache key 'mail resend': %v", err)
	// 	}
	// 	return
	// }

	c.RedirectSubpath("/login")
}
