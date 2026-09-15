package health

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// Serve starts a small HTTP server exposing the probes.
//
// This service has no HTTP server of its own - it is a worker - so one is added
// purely for monitoring. It listens on HEALTH_PORT (default 8080) and serves
// only these endpoints, which keeps the surface small enough to expose on the
// internal network without further thought.
//
// Runs in its own goroutine and never blocks the worker: a monitoring endpoint
// that could stall payment processing would be worse than no endpoint at all.
func Serve(s *Service) {
	port := os.Getenv("HEALTH_PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, s.Liveness())
	})
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		report, status := s.Readiness(r.Context())
		writeJSON(w, status, report)
	})

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("[health] probes listening on :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// A failure to bind must not take the engine down: payment
			// processing matters more than the ability to report on it.
			log.Printf("[health] probe server stopped: %v", err)
		}
	}()
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("[health] failed to encode response: %v", err)
	}
}
