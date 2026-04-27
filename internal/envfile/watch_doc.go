// Package envfile — Watch
//
// The Watcher type provides file-based change detection for env files.
// It polls the filesystem at a configurable interval and emits WatchEvent
// values whenever a file's MD5 hash changes.
//
// Usage:
//
//	w := envfile.NewWatcher(500*time.Millisecond, ".env", ".env.local")
//	w.Start()
//	defer w.Stop()
//
//	for {
//		select {
//		case ev := <-w.Events:
//			fmt.Printf("changed: %s\n", ev.Path)
//		case err := <-w.Errors:
//			log.Println("watch error:", err)
//		}
//	}
//
// WatchEvent carries the file path, old and new MD5 hashes, and the
// freshly parsed EnvMap so callers can react immediately without a
// second parse call.
package envfile
