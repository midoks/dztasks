package conf

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"net/url"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/ini.v1"

	"github.com/midoks/dztasks/embed"
	"github.com/midoks/dztasks/internal/tools"
)

// File is the configuration object.
var File *ini.File

func autoMakeCustomConf(customConf string) error {

	if tools.IsExist(customConf) {
		return nil
	}

	// auto make custom conf file
	cfg := ini.Empty()
	if tools.IsFile(customConf) {
		if err := cfg.Append(customConf); err != nil {
			return err
		}
	}

	cfg.Section("").Key("app_name").SetValue("dztasks")
	cfg.Section("").Key("run_mode").SetValue("prod")

	cfg.Section("web").Key("http_port").SetValue("11011")
	cfg.Section("session").Key("provider").SetValue("memory")

	cfg.Section("admin").Key("user").SetValue(tools.RandString(8))
	cfg.Section("admin").Key("pass").SetValue(tools.RandString(10))

	cfg.Section("plugins").Key("path").SetValue("plugins")

	os.MkdirAll(filepath.Dir(customConf), os.ModePerm)
	if err := cfg.SaveTo(customConf); err != nil {
		return err
	}

	return nil
}

func Init(customConf string) error {

	data, _ := embed.Conf.ReadFile("conf/app.conf")

	// fmt.Println(data)
	File, err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, data)
	if err != nil {
		return errors.Wrap(err, "parse 'conf/app.conf'")
	}

	if customConf == "" {
		customConf = filepath.Join(CustomDir(), "conf", "app.conf")
		autoMakeCustomConf(customConf)
	} else {
		customConf, err = filepath.Abs(customConf)
		if err != nil {
			return errors.Wrap(err, "get absolute path")
		}
	}
	CustomConf = customConf

	if tools.IsFile(customConf) {
		if err = File.Append(customConf); err != nil {
			return errors.Wrapf(err, "append %q", customConf)
		}
	} else {
		log.Println("Custom config ", customConf, " not found. Ignore this warning if you're running for the first time")
	}

	File.NameMapper = ini.TitleUnderscore

	if err = File.Section(ini.DefaultSection).MapTo(&App); err != nil {
		return errors.Wrap(err, "mapping default section")
	}

	// ***************************
	// ----- Log settings -----
	// ***************************
	if err = File.Section("log").MapTo(&Log); err != nil {
		return errors.Wrap(err, "mapping [log] section")
	}

	// ****************************
	// ----- Web settings -----
	// ****************************

	if err = File.Section("web").MapTo(&Web); err != nil {
		return errors.Wrap(err, "mapping [web] section")
	}

	Web.AppDataPath = ensureAbs(Web.AppDataPath)

	if !strings.HasSuffix(Web.ExternalURL, "/") {
		Web.ExternalURL += "/"
	}
	Web.URL, err = url.Parse(Web.ExternalURL)
	if err != nil {
		return errors.Wrapf(err, "parse '[server] EXTERNAL_URL' %q", err)
	}

	// Subpath should start with '/' and end without '/', i.e. '/{subpath}'.
	Web.Subpath = strings.TrimRight(Web.URL.Path, "/")
	Web.SubpathDepth = strings.Count(Web.Subpath, "/")

	unixSocketMode, err := strconv.ParseUint(Web.UnixSocketPermission, 8, 32)
	if err != nil {
		return errors.Wrapf(err, "parse '[server] unix_socket_permission' %q", Web.UnixSocketPermission)
	}
	if unixSocketMode > 0777 {
		unixSocketMode = 0666
	}
	Web.UnixSocketMode = os.FileMode(unixSocketMode)

	// ****************************
	// ----- Session settings -----
	// ****************************

	if err = File.Section("session").MapTo(&Session); err != nil {
		return errors.Wrap(err, "mapping [session] section")
	}

	// ****************************
	// ----- Session settings -----
	// ****************************

	if err = File.Section("admin").MapTo(&Admin); err != nil {
		return errors.Wrap(err, "mapping [admin] section")
	}

	// ****************************
	// ----- Session settings -----
	// ****************************

	if err = File.Section("plugins").MapTo(&Plugins); err != nil {
		return errors.Wrap(err, "mapping [plugins] section")
	}

	// *****************************
	// ----- Security settings -----
	// *****************************

	if err = File.Section("security").MapTo(&Security); err != nil {
		return errors.Wrap(err, "mapping [security] section")
	}

	// Check run user when the install is locked.
	if Security.InstallLock {
		currentUser, match := CheckRunUser(App.RunUser)
		if !match {
			return fmt.Errorf("user configured to run imail is %q, but the current user is %q", App.RunUser, currentUser)
		}
	}

	return nil
}
