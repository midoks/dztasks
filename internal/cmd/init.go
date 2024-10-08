package cmd

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"

// 	"gopkg.in/ini.v1"
// )

// func autoMakeCustomConf(customConf string) error {

// 	if tools.IsExist(customConf) {
// 		return nil
// 	}

// 	// auto make custom conf file
// 	cfg := ini.Empty()
// 	if tools.IsFile(customConf) {
// 		if err := cfg.Append(customConf); err != nil {
// 			return err
// 		}
// 	}

// 	cfg.Section("").Key("app_name").SetValue("vez")
// 	cfg.Section("").Key("run_mode").SetValue("prod")

// 	cfg.Section("web").Key("http_port").SetValue("60000")
// 	cfg.Section("session").Key("provider").SetValue("memory")

// 	cfg.Section("mongodb").Key("addr").SetValue("127.0.0.1:27017")
// 	cfg.Section("mongodb").Key("db").SetValue("vez")

// 	// cfg.Section("image").Key("addr").SetValue("http://0.0.0.0:3333/i/")
// 	// cfg.Section("image").Key("ping").SetValue("http://0.0.0.0:3333/ping")
// 	// cfg.Section("image").Key("ping_response").SetValue("ok")

// 	cfg.Section("security").Key("install_lock").SetValue("true")
// 	secretKey := tools.RandString(15)
// 	cfg.Section("security").Key("secret_key").SetValue(secretKey)

// 	os.MkdirAll(filepath.Dir(customConf), os.ModePerm)
// 	if err := cfg.SaveTo(customConf); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func Init(customConf string) error {
// 	var err error

// 	if customConf == "" {
// 		customConf = filepath.Join(conf.CustomDir(), "conf", "app.conf")
// 	} else {
// 		customConf, err = filepath.Abs(customConf)
// 		if err != nil {
// 			return fmt.Errorf("custom conf file get absolute path: %s", err)
// 		}
// 	}

// 	err = autoMakeCustomConf(customConf)
// 	if err != nil {
// 		return err
// 	}

// 	conf.Init(customConf)
// 	logs.Init()

// 	return nil
// }

// func init() {
// 	Init("")
// }
