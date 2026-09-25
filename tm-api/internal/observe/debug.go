package observe

import (
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"time"
)

// ServeDebug exposes Go's pprof profiles - goroutine, heap, CPU, block - on a
// loopback-only listener, so a leak can be diagnosed with one request instead
// of a day of code reading.
//
// Loopback only is the whole point: the profiles reveal internal structure and
// can be expensive to produce, so they must never be reachable from the edge.
// Binding to 127.0.0.1 inside the container makes them reachable solely via
// `docker exec <container> wget -qO- localhost:6060/debug/pprof/goroutine?debug=1`
// and nothing else - no Traefik rule to get wrong, nothing on the overlay network.
//
// DEBUG_ADDR overrides the address; leave it unset in every deployed environment.
func ServeDebug() {
	addr := os.Getenv("DEBUG_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6060"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		log.Printf("[observe] pprof listening on %s (loopback only)", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Diagnostics must never take the service down with them.
			log.Printf("[observe] pprof server stopped: %v", err)
		}
	}()
}
