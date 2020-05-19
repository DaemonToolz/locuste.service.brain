package main

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher

func initFileWatcher(root string) {
	defer func() {
		if r := recover(); r != nil {
			AddOrUpdateStatus(BrainHttpServer, false)
		}
	}()

	watcher, _ = fsnotify.NewWatcher()

	if err := filepath.Walk(root, watchDir); err != nil {
		failOnError(err, "Couldn't get to the folder")
		AddOrUpdateStatus(BrainWatcher, false)
	} else {
		AddOrUpdateStatus(BrainWatcher, true)
	}
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}