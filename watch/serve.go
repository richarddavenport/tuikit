package watch

import (
	"context"
	"fmt"
	"html"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

// server holds the page and tells browsers when it changes.
type server struct {
	mu      sync.RWMutex
	page    []byte
	version int
	err     error
}

func (s *server) set(page []byte, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.page, s.err, s.version = page, err, s.version+1
}

func (s *server) read() ([]byte, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.page, s.version, s.err
}

// Run watches, rebuilds and serves until ctx is done. It returns the address it
// listened on through addr, which the caller prints.
func (c Config) Run(ctx context.Context, addr string, ready func(string)) error {
	srv := &server{}
	poll := NewPoller(c.Dir)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page, version, buildErr := srv.read()
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if buildErr != nil {
			// The error goes ON the page. Leaving the last good frames up when
			// the code no longer compiles is a page that lies about the tool.
			_, _ = fmt.Fprint(w, errorPage(buildErr, version))
			return
		}
		_, _ = w.Write(append(page, []byte(reloadScript(version))...))
	})

	// The browser asks for the version rather than holding a socket open. A
	// long poll would be fewer requests; a number every half second is fewer
	// moving parts, and a dev server's job is to be boring.
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		_, version, _ := srv.read()
		w.Header().Set("Cache-Control", "no-store")
		_, _ = fmt.Fprint(w, version)
	})

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer func() { _ = listener.Close() }()

	if ready != nil {
		ready("http://" + listener.Addr().String())
	}

	srvHTTP := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() { _ = srvHTTP.Serve(listener) }()
	defer func() { _ = srvHTTP.Close() }()

	// Debounced: an editor writing a file in two syscalls is one save, and
	// rebuilding twice makes the page flicker through a state nobody asked for.
	var pending bool
	var settleAt time.Time
	ticker := time.NewTicker(c.interval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}

		changed, err := poll.Changed()
		if err != nil {
			srv.set(nil, err)
			continue
		}
		if changed {
			pending, settleAt = true, time.Now().Add(c.interval())
			continue
		}
		if !pending || time.Now().Before(settleAt) {
			continue
		}
		pending = false

		start := time.Now()
		if err := c.Build(); err != nil {
			c.logf("failed in %s", time.Since(start).Round(time.Millisecond))
			srv.set(nil, err)
			continue
		}
		page, err := os.ReadFile(c.Out)
		srv.set(page, err)
		c.logf("rebuilt in %s", time.Since(start).Round(time.Millisecond))
	}
}

// reloadScript polls the version and reloads when it moves.
func reloadScript(version int) string {
	return fmt.Sprintf(`<script>
(function () {
  var seen = %d;
  setInterval(function () {
    fetch("/version", { cache: "no-store" })
      .then(function (r) { return r.text(); })
      .then(function (v) { if (+v !== seen) location.reload(); })
      .catch(function () {});
  }, 500);
})();
</script>`, version)
}

// errorPage shows what went wrong, in the frames' own monospace so a Go error
// with a file and a line is readable as one.
func errorPage(err error, version int) string {
	return fmt.Sprintf(`<!doctype html>
<html lang="en"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>capture failed</title>
<style>
:root { --ground:#fbfafb; --ink:#22202a; --bad:#b3261e; --line:#e6e3ec; }
@media (prefers-color-scheme: dark) { :root:not([data-theme="light"]) {
  --ground:#131218; --ink:#e8e6f0; --bad:#ff8a80; --line:#2b2936; } }
:root[data-theme="dark"] { --ground:#131218; --ink:#e8e6f0; --bad:#ff8a80; --line:#2b2936; }
body { margin:0; padding:40px; background:var(--ground); color:var(--ink);
  font:16px/1.6 ui-sans-serif,-apple-system,"Segoe UI",system-ui,sans-serif; }
h1 { font-size:22px; margin:0 0 6px; color:var(--bad); }
p { margin:0 0 20px; color:var(--ink); opacity:.75; max-width:64ch; }
pre { font-family:ui-monospace,"SF Mono",Menlo,Consolas,monospace; font-size:13px;
  line-height:1.45; white-space:pre-wrap; background:#0c0c0e; color:#f0d0d0;
  border:1px solid #26262c; border-radius:7px; padding:16px 18px; overflow-x:auto; }
</style></head><body>
<h1>The capture failed</h1>
<p>The frames from before this change are not shown. They would describe a tool
that no longer exists, which is worse than showing nothing.</p>
<pre>%s</pre>
%s
</body></html>`, html.EscapeString(err.Error()), reloadScript(version))
}
