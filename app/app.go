package app

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"

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

// fileSystem implements macaron.TemplateFileSystem for embedded templates
type fileSystem struct {
	files []macaron.TemplateFile
}

// ListFiles returns all template files in the filesystem
func (fs *fileSystem) ListFiles() []macaron.TemplateFile {
	return fs.files
}

// Get retrieves a template file by name
func (fs *fileSystem) Get(name string) (io.Reader, error) {
	for i := range fs.files {
		if fs.files[i].Name()+fs.files[i].Ext() == name {
			return bytes.NewReader(fs.files[i].Data()), nil
		}
	}
	return nil, fmt.Errorf("file %q not found", name)
}

// newTemplateFileSystem returns a macaron.TemplateFileSystem instance for embedded assets.
// The argument "dir" can be used to serve subset of embedded assets. Template file
// found under the "customDir" on disk has higher precedence over embedded assets.
func newTemplateFileSystem(dir, customDir string) macaron.TemplateFileSystem {
	var files []macaron.TemplateFile

	if dir != "" && !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	// Load embedded templates
	allfn := embed.TemplatesAllNames("templates")
	for _, name := range allfn {
		// Filter by directory if specified
		if dir != "" && !strings.HasPrefix(name, "templates/"+dir) {
			continue
		}

		ext := path.Ext(name)
		data, err := embed.Templates.ReadFile(name)
		if err != nil {
			continue // Skip files that can't be read
		}

		name = strings.TrimPrefix(name, "templates")
		name = strings.TrimSuffix(name, ext)

		// Check if custom template exists and use it instead
		if customDir != "" {
			customPath := filepath.Join(customDir, name+ext)
			if customData, err := embed.Templates.ReadFile(customPath); err == nil {
				data = customData
			}
		}

		files = append(files, macaron.NewTplFile(name, data, ext))
	}

	return &fileSystem{files: files}
}

// staticFunc returns an expiration time for static assets (24 hours from now)
func staticFunc() string {
	const rfc1123Format = "Mon, 02 Jan 2006 15:04:05 +0800"
	return time.Now().Add(24 * time.Hour).Format(rfc1123Format)
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
			ETag:        true,
			Expires:     staticFunc,
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
			ETag:        true,
			Expires:     staticFunc,
		},
	))

	// template start
	renderOpt := macaron.RenderOptions{
		Directory:         filepath.Join(conf.WorkDir(), "templates"),
		AppendDirectories: []string{filepath.Join(conf.CustomDir(), "templates")},
		Funcs:             template.FuncMap(),
		Delims:            macaron.Delims{Left: "{[", Right: "]}"},
		IndentJSON:        macaron.Env != macaron.PROD,
	}
	if !conf.Web.LoadAssetsFromDisk {
		renderOpt.TemplateFileSystem = newTemplateFileSystem("", renderOpt.AppendDirectories[0])
	}

	m.Use(macaron.Renderer(renderOpt))
	// template end

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
		m.Get("/plugin/page", reqSignIn, bindIgnErr(form.ArgsPluginPage{}), plugin.PluginPage)
		m.Get("/plugin/file", reqSignIn, bindIgnErr(form.ArgsPluginFile{}), plugin.PluginFile)
		m.Post("/plugin/install", reqSignIn, bindIgnErr(form.ArgsPluginInstall{}), plugin.PluginInstall)
		m.Post("/plugin/uninstall", reqSignIn, bindIgnErr(form.ArgsPluginUninstall{}), plugin.PluginUninstall)
		m.Group("/plugin/data", func() {
			m.Combo("").Get(bindIgnErr(form.ArgsPluginData{}), plugin.PluginData).Post(bindIgnErr(form.ArgsPluginData{}), plugin.PluginData)
		}, reqSignIn)
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
