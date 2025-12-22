package blocker

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func (b *Blocker) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err := watcher.Add(b.file); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if ev.Op&fsnotify.Write != 0 {
					log.Println("[blocker] rules updated, reloading")
					_ = b.LoadFile(b.file)
				}
			case err := <-watcher.Errors:
				log.Println("[blocker] watcher error:", err)
			}
		}
	}()

	return nil
}
