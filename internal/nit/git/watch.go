package git

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type FSWatcher struct {
	w       *fsnotify.Watcher
	events  chan struct{}
	stop    chan struct{}
	stopped chan struct{}
	once    sync.Once
}

func (s Service) NewFSWatcher() (*FSWatcher, error) {
	root, _, err := s.runner.Run("--no-optional-locks", "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, err
	}
	return newFSWatcher(strings.TrimSpace(root))
}

func newFSWatcher(root string) (*FSWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	fw := &FSWatcher{
		w:       w,
		events:  make(chan struct{}, 1),
		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
	if err := fw.addRecursive(root); err != nil {
		_ = w.Close()
		return nil, err
	}
	// Watch the .git directory tree as well for index/HEAD/refs changes.
	_ = fw.addRecursive(filepath.Join(root, ".git"))
	go fw.loop(root)
	return fw, nil
}

func (fw *FSWatcher) Events() <-chan struct{} { return fw.events }

func (fw *FSWatcher) Close() error {
	var err error
	fw.once.Do(func() {
		close(fw.stop)
		<-fw.stopped
		err = fw.w.Close()
	})
	return err
}

func (fw *FSWatcher) loop(root string) {
	defer close(fw.stopped)

	const debounce = 200 * time.Millisecond
	timer := time.NewTimer(time.Hour)
	if !timer.Stop() {
		<-timer.C
	}
	pending := false

	emit := func() {
		select {
		case fw.events <- struct{}{}:
		default:
		}
	}
	schedule := func() {
		pending = true
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(debounce)
	}

	for {
		select {
		case <-fw.stop:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return
		case ev, ok := <-fw.w.Events:
			if !ok {
				return
			}
			if ev.Op&(fsnotify.Create) != 0 {
				if info, err := os.Stat(ev.Name); err == nil && info.IsDir() {
					_ = fw.addRecursive(ev.Name)
				}
			}
			// Skip noisy temporary files outside .git if desired in future.
			if strings.HasPrefix(ev.Name, root) {
				schedule()
			}
		case <-fw.w.Errors:
			// Ignore watcher errors and rely on fallback polling.
		case <-timer.C:
			if pending {
				emit()
				pending = false
			}
		}
	}
}

func (fw *FSWatcher) addRecursive(root string) error {
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return err
	}
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		if base == ".git" && path != root {
			// We already add .git separately from the repo root.
			return filepath.SkipDir
		}
		if base == ".cache" || base == "node_modules" {
			// Avoid huge noisy subtrees; fallback polling still keeps correctness.
			if path != root {
				return filepath.SkipDir
			}
		}
		_ = fw.w.Add(path)
		return nil
	})
}
