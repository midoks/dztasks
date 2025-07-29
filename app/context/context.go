// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"gopkg.in/macaron.v1"

	"github.com/go-macaron/cache"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"

	"github.com/midoks/dztasks/app/common"
	"github.com/midoks/dztasks/app/form"
	"github.com/midoks/dztasks/app/template"
	"github.com/midoks/dztasks/internal/conf"
	"github.com/midoks/dztasks/internal/log"
)

type ToggleOptions struct {
	SignInRequired  bool
	SignOutRequired bool
	AdminRequired   bool
	DisableCSRF     bool
}

func Toggle(options *ToggleOptions) macaron.Handler {
	return func(c *Context) {
		if options.SignInRequired {
			if !c.IsLogged {
				c.SetCookie("redirect_to", url.QueryEscape(conf.Web.Subpath+c.Req.RequestURI), 0, conf.Web.Subpath)
				c.RedirectSubpath("/login")
				return
			}
		}
	}
}

// Context represents context of a request.
type Context struct {
	*macaron.Context
	Cache   cache.Cache
	csrf    csrf.CSRF
	Flash   *session.Flash
	Session session.Store

	Link        string // Current request URL
	IsLogged    bool
	IsBasicAuth bool
	IsTokenAuth bool
}

// JSONData represents the standard JSON API response structure
type JSONData struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// LayuiData represents the Layui table JSON response structure
type LayuiData struct {
	Code  int64       `json:"code"`
	Count int         `json:"count"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
}

// RawTitle sets the "Title" field in template data.
func (c *Context) RawTitle(title string) {
	c.Data["Title"] = title
}

// Title localizes the "Title" field in template data.
func (c *Context) Title(locale string) {
	c.RawTitle(c.Tr(locale))
}

// RenderWithErr used for page has form validation but need to prompt error to users.
func (c *Context) RenderWithErr(msg, tpl string, f interface{}) {
	if f != nil {
		form.Assign(f, c.Data)
	}
	c.Flash.ErrorMsg = msg
	c.Data["Flash"] = c.Flash
	c.HTML(http.StatusOK, tpl)
}

// HasError returns true if error occurs in form validation.
func (c *Context) HasError() bool {
	hasErr, ok := c.Data["HasError"]
	if !ok {
		return false
	}
	c.Flash.ErrorMsg = c.Data["ErrorMsg"].(string)
	c.Data["Flash"] = c.Flash
	return hasErr.(bool)
}

// PageIs sets "PageIsxxx" field in template data.
func (c *Context) PageIs(name string) {
	c.Data["PageIs"+name] = true
}

// Require sets "Requirexxx" field in template data.
func (c *Context) Require(name string) {
	c.Data["Require"+name] = true
}

// FormErr sets "Err_xxx" field in template data.
func (c *Context) FormErr(names ...string) {
	for i := range names {
		c.Data["Err_"+names[i]] = true
	}
}

func (c *Context) GetErrMsg() string {
	return c.Data["ErrorMsg"].(string)
}

// HasValue returns true if value of given name exists.
func (c *Context) HasValue(name string) bool {
	_, ok := c.Data[name]
	return ok
}

// HTML responses template with given status.
func (c *Context) HTML(status int, name string) {
	// log.Infof("Template:%s", name)
	c.Context.HTML(status, name)
}

func (c *Context) HTMLString(name string, content string) {
	c.Context.HTMLString(name, content)
}

// TEXT responses template with given status.
func (c *Context) PlainText(status int, name []byte) {
	c.Context.PlainText(status, name)
}

func (c *Context) RawData(status int, v []byte) {
	c.Context.RawData(status, v)
}

// Success responses template with status http.StatusOK.
func (c *Context) Success(name string) {
	c.HTML(http.StatusOK, name)
}

// RenderJSON renders JSON response with status http.StatusOK
func (c *Context) RenderJSON(data interface{}) {
	c.JSON(http.StatusOK, data)
}

// ReturnJSON returns a standard JSON API response
func (c *Context) ReturnJSON(code int64, msg string, data interface{}) {
	c.RenderJSON(JSONData{Code: code, Msg: msg, Data: data})
}

// ReturnLayuiJSON returns a Layui table JSON response
func (c *Context) ReturnLayuiJSON(code int64, msg string, count int, data interface{}) {
	c.RenderJSON(LayuiData{Code: code, Msg: msg, Count: count, Data: data})
}

// Ok returns a success JSON response
func (c *Context) Ok(msg string) {
	c.ReturnJSON(1, msg, "")
}

// Fail returns a failure JSON response
func (c *Context) Fail(msg string) {
	c.ReturnJSON(-1, msg, "")
}

// NotFound renders the 404 page.
func (c *Context) NotFound() {
	c.HTML(http.StatusNotFound, fmt.Sprintf("status/%d", http.StatusNotFound))
}

// Error renders the 500 page.
func (c *Context) Error(err error, msg string) {
	// c.Title("status.internal_server_error")

	// Only in non-production mode or admin can see the actual error message.
	if !conf.IsProdMode() || (c.IsLogged) {
		c.Data["ErrorMsg"] = err
	}
	c.HTML(http.StatusInternalServerError, fmt.Sprintf("status/%d", http.StatusInternalServerError))
}

// Errorf renders the 500 response with formatted message.
func (c *Context) Errorf(err error, format string, args ...interface{}) {
	c.Error(err, fmt.Sprintf(format, args...))
}

// NotFoundOrError responses with 404 page for not found error and 500 page otherwise.
func (c *Context) NotFoundOrError(err error, msg string) {
	// if errutil.IsNotFound(err) {
	// 	c.NotFound()
	// 	return
	// }
	c.Error(err, msg)
}

// RedirectSubpath responses redirection with given location and status.
// It prepends setting.Server.Subpath to the location string.
func (c *Context) RedirectSubpath(location string, status ...int) {
	c.Redirect(conf.Web.Subpath+location, status...)
}

// Contexter initializes a classic context for a request.
func Contexter() macaron.Handler {
	return func(ctx *macaron.Context, cache cache.Cache, sess session.Store, f *session.Flash, x csrf.CSRF) {
		c := &Context{
			Context: ctx,
			Cache:   cache,
			csrf:    x,
			Flash:   f,
			Session: sess,
		}

		ctx.Map(c)
		c.Data["PageStartTime"] = time.Now()

		if len(conf.Web.AccessControlAllowOrigin) > 0 {
			c.Header().Set("Access-Control-Allow-Origin", conf.Web.AccessControlAllowOrigin)
			c.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Header().Set("Access-Control-Max-Age", "3600")
			c.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		}

		isLogged := c.Session.Get("login")
		name := c.Session.Get("name")
		if isLogged == true {
			c.IsLogged = true
			c.Data["LoggedUserName"] = name
			c.Data["IsAdmin"] = true
		}

		c.Data["CSRFToken"] = x.GetToken()
		c.Data["CSRFTokenHTML"] = template.Safe(`<input type="hidden" name="_csrf" value="` + x.GetToken() + `">`)

		if !conf.IsProdMode() {
			log.Debugf("Session ID: %s", sess.ID())
			log.Debugf("CSRF Token: %s", c.Data["CSRFToken"])
		}
		pluginDir := conf.Plugins.Path
		c.Data["PluginMenu"] = common.PluginList(pluginDir)

		// avoid iframe use by other.
		// c.Header().Set("X-Content-Type-Options", "nosniff")
		// c.Header().Set("X-Frame-Options", "DENY")
	}
}
