package app

import (
	"fmt"
	"io"
	// "io/fs"
	// "io/ioutil"
	"bytes"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/macaron.v1"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/gzip"
	"github.com/go-macaron/session"

	"github.com/midoks/dztasks/app/context"
	"github.com/midoks/dztasks/app/form"
	"github.com/midoks/dztasks/app/router"
	"github.com/midoks/dztasks/app/router/plugin"
	"github.com/midoks/dztasks/app/router/user"
	"github.com/midoks/dztasks/app/template"
	"github.com/midoks/dztasks/embed"
	"github.com/midoks/dztasks/internal/conf"
)

type fileSystem struct {
	files []macaron.TemplateFile
}

func (fs *fileSystem) ListFiles() []macaron.TemplateFile {
	// for i := range fs.files {
	// 	fmt.Println(fs.files[i].Name())
	// }
	return fs.files
}

func (fs *fileSystem) Get(name string) (io.Reader, error) {

	for i := range fs.files {
		if fs.files[i].Name()+fs.files[i].Ext() == name {
			return bytes.NewReader(fs.files[i].Data()), nil
		}
	}
	return nil, fmt.Errorf("file %q not found", name)
}

// NewTemplateFileSystem returns a macaron.TemplateFileSystem instance for embedded assets.
// The argument "dir" can be used to serve subset of embedded assets. Template file
// found under the "customDir" on disk has higher precedence over embedded assets.
func newTemplateFileSystem(dir, customDir string) macaron.TemplateFileSystem {
	var files []macaron.TemplateFile

	if dir != "" && !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	allfn := embed.TemplatesAllNames("templates")

	for _, name := range allfn {
		ext := path.Ext(name)
		data, _ := embed.Templates.ReadFile(name)

		name = strings.TrimPrefix(name, "templates")
		name = strings.TrimSuffix(name, ext)
		files = append(files, macaron.NewTplFile(name, data, ext))
	}

	return &fileSystem{files: files}
}

func bootstrapMacaron() *macaron.Macaron {
	m := macaron.New()

	if !conf.Web.DisableRouterLog {
		m.Use(macaron.Logger())
	}

	m.Use(macaron.Recovery())

	if conf.Web.EnableGzip {
		m.Use(gzip.Gziper())
	}

	m.Use(macaron.Static(
		filepath.Join(conf.CustomDir(), "static"),
		macaron.StaticOptions{
			SkipLogging: conf.Web.DisableRouterLog,
		},
	))

	var staticFs http.FileSystem
	if !conf.Web.LoadAssetsFromDisk {
		staticFs = http.FS(embed.Static)
	}

	m.Use(macaron.Static(
		filepath.Join(conf.WorkDir(), "static"),
		macaron.StaticOptions{
			FileSystem:  staticFs,
			SkipLogging: conf.Web.DisableRouterLog,
		},
	))

	//template start
	renderOpt := macaron.RenderOptions{
		Directory:         filepath.Join(conf.WorkDir(), "templates"),
		AppendDirectories: []string{filepath.Join(conf.CustomDir(), "templates")},
		Funcs:             template.FuncMap(),
		Delims:            macaron.Delims{"{[", "]}"},
		IndentJSON:        macaron.Env != macaron.PROD,
	}
	if !conf.Web.LoadAssetsFromDisk {
		renderOpt.TemplateFileSystem = newTemplateFileSystem("", renderOpt.AppendDirectories[0])
	}

	m.Use(macaron.Renderer(renderOpt))
	//template end

	m.Use(cache.Cacher(cache.Options{
		Adapter:       conf.Cache.Adapter,
		AdapterConfig: conf.Cache.Host,
		Interval:      conf.Cache.Interval,
	}))

	m.Use(captcha.Captchaer(captcha.Options{
		SubURL: conf.Web.Subpath,
	}))

	return m
}

func bootstrapRouter(m *macaron.Macaron) *macaron.Macaron {

	reqSignIn := context.Toggle(&context.ToggleOptions{SignInRequired: true})
	reqSignOut := context.Toggle(&context.ToggleOptions{SignOutRequired: true})

	bindIgnErr := binding.BindIgnErr
	m.SetAutoHead(true)

	m.Group("", func() {

		m.Group("", func() {
			m.Group("/login", func() {
				m.Combo("").Get(user.Login).Post(bindIgnErr(form.SignIn{}), user.LoginPost)
			})
			m.Get("/logout", user.SignOut)
		}, reqSignOut)
		m.Get("/", reqSignIn, router.Home)
		m.Get("/log", reqSignIn, router.Log)
		m.Get("/plugin", reqSignIn, plugin.PluginHome)
		m.Get("/plugin/list", reqSignIn, plugin.PluginList)
		m.Get("/plugin/menu", reqSignIn, bindIgnErr(form.ArgsPluginMenu{}), plugin.PluginMenu)
		m.Post("/plugin/install", reqSignIn, bindIgnErr(form.ArgsPluginInstall{}), plugin.PluginInstall)
		m.Post("/plugin/uninstall", reqSignIn, bindIgnErr(form.ArgsPluginUninstall{}), plugin.PluginUninstall)

	}, session.Sessioner(session.Options{
		Provider:       conf.Session.Provider,
		ProviderConfig: conf.Session.ProviderConfig,
		CookieName:     conf.Session.CookieName,
		CookiePath:     conf.Web.Subpath,
		Gclifetime:     conf.Session.GCInterval,
		Maxlifetime:    conf.Session.MaxLifeTime,
		Secure:         conf.Session.CookieSecure,
	}), csrf.Csrfer(csrf.Options{
		Secret:         conf.Security.SecretKey,
		Header:         "X-CSRF-Token",
		Cookie:         conf.Session.CSRFCookieName,
		CookieDomain:   conf.Web.URL.Hostname(),
		CookiePath:     conf.Web.Subpath,
		CookieHttpOnly: true,
		SetCookie:      true,
		Secure:         conf.Web.URL.Scheme == "https",
	}), context.Contexter())
	return m
}

func Start(port int) {
	boot := bootstrapMacaron()
	boot = bootstrapRouter(boot)
	boot.Run(port)
}
