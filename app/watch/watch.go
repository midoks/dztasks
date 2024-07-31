package watch

import (
	// "fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"

	"github.com/midoks/dztasks/app/bgtask"
	"github.com/midoks/dztasks/internal/log"
)

func InitWatch(file_pos string) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Infof("init failed to create watcher: %T", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				log.Info("文件发生变动!")
				if !ok {
					return
				}
				bgtask.ResetTask()
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()

	var files []string
	err = filepath.Walk(file_pos, func(file string, info os.FileInfo, err error) error {

		extension := path.Ext(file)
		// fmt.Println(file, extension)
		if strings.EqualFold(extension, ".json") {
			// fmt.Println(file)
			files = append(files, file)
		}

		return nil
	})

	for _, f := range files {
		err = watcher.Add(f)
		if err != nil {
			log.Infof("watching: %T", err)
		}
	}

	<-make(chan struct{})
}
