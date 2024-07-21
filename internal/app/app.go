package app


import (
	"net/http"
	"path/filepath"

	"gopkg.in/macaron.v1"

	// "github.com/go-macaron/binding"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	// "github.com/go-macaron/csrf"
	"github.com/go-macaron/gzip"
	// "github.com/go-macaron/session"

	// "github.com/midoks/dztasks/internal/app/context"
	// "github.com/midoks/dztasks/internal/app/router"
	"github.com/midoks/dztasks/internal/app/template"
	"github.com/midoks/dztasks/internal/assets/templates"
	"github.com/midoks/dztasks/internal/conf"
)

func newMacaron() *macaron.Macaron {
	m := macaron.New()

	if !conf.Web.DisableRouterLog {
		m.Use(macaron.Logger())
	}

	m.Use(macaron.Recovery())

	if conf.Web.EnableGzip {
		m.Use(gzip.Gziper())
	}

	m.Use(macaron.Static(
		filepath.Join(conf.CustomDir(), "public"),
		macaron.StaticOptions{
			SkipLogging: conf.Web.DisableRouterLog,
		},
	))

	var publicFs http.FileSystem
	if !conf.Web.LoadAssetsFromDisk {
		// publicFs = public.NewFileSystem()
		publicFs = http.FS(conf.App.PublicFs)
	}

	m.Use(macaron.Static(
		filepath.Join(conf.WorkDir(), "public"),
		macaron.StaticOptions{
			FileSystem:  publicFs,
			SkipLogging: conf.Web.DisableRouterLog,
		},
	))

	//template start
	renderOpt := macaron.RenderOptions{
		Directory:         filepath.Join(conf.WorkDir(), "templates"),
		AppendDirectories: []string{filepath.Join(conf.CustomDir(), "templates")},
		Funcs:             template.FuncMap(),
		IndentJSON:        macaron.Env != macaron.PROD,
	}
	if !conf.Web.LoadAssetsFromDisk {
		// renderOpt.TemplateFileSystem = http.FS(conf.App.TemplateFs)
		renderOpt.TemplateFileSystem = templates.NewTemplateFileSystem("", renderOpt.AppendDirectories[0])
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


func Start(port int) {
	m := newMacaron()
	// m = setRouter(m)
	m.Run(port)
}