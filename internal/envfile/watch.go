package envfile

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// WatchEvent describes a change detected in a watched env file.
type WatchEvent struct {
	Path    string
	OldHash string
	NewHash string
	Env     EnvMap
}

// Watcher polls one or more env files for changes and emits events.
type Watcher struct {
	paths    []string
	hashes   map[string]string
	interval time.Duration
	Events   chan WatchEvent
	Errors   chan error
	stop     chan struct{}
	mu       sync.Mutex
}

// NewWatcher creates a Watcher that checks the given paths every interval.
func NewWatcher(interval time.Duration, paths ...string) *Watcher {
	return &Watcher{
		paths:    paths,
		hashes:   make(map[string]string),
		interval: interval,
		Events:   make(chan WatchEvent, 16),
		Errors:   make(chan error, 8),
		stop:     make(chan struct{}),
	}
}

// Start begins polling in a background goroutine.
func (w *Watcher) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.poll()
			case <-w.stop:
				return
			}
		}
	}()
}

// Stop shuts down the polling goroutine.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) poll() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, p := range w.paths {
		newHash, err := hashFile(p)
		if err != nil {
			w.Errors <- fmt.Errorf("watch: %w", err)
			continue
		}
		oldHash := w.hashes[p]
		if oldHash == newHash {
			continue
		}
		env, err := ParseFile(p)
		if err != nil {
			w.Errors <- fmt.Errorf("watch parse %s: %w", p, err)
			continue
		}
		w.hashes[p] = newHash
		w.Events <- WatchEvent{Path: p, OldHash: oldHash, NewHash: newHash, Env: env}
	}
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
