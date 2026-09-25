package observe

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// requestsTotal counts requests by route, method and status. Labelled by
	// gin's route template (/users/:id), never the raw path, so that an
	// endpoint hit with many different IDs stays one time series instead of
	// exploding cardinality.
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests by route, method and status code.",
		},
		[]string{"service", "method", "route", "status"},
	)

	// requestDuration records latency distribution per route, which is what
	// makes "the API got slower this week" answerable.
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency by route and method.",
			Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"service", "method", "route"},
	)

	// requestsInFlight shows concurrent load, which distinguishes "slow
	// endpoint" from "server saturated".
	requestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "HTTP requests currently being served.",
		},
	)

	// panicsTotal counts recovered panics. Any non-zero rate is worth an alert.
	panicsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_panics_total",
			Help: "Panics recovered while serving HTTP requests.",
		},
		[]string{"service", "route"},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal, requestDuration, requestsInFlight, panicsTotal)
}

// Metrics records Prometheus metrics for every request.
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip the scrape endpoint itself so Prometheus does not measure its
		// own polling.
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()
		requestsInFlight.Inc()
		defer requestsInFlight.Dec()

		c.Next()

		// FullPath is the route template and is empty for unmatched requests;
		// grouping those under a constant keeps 404 scans from creating a new
		// time series per probed URL.
		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}

		requestsTotal.WithLabelValues(serviceName, c.Request.Method, route, strconv.Itoa(c.Writer.Status())).Inc()
		requestDuration.WithLabelValues(serviceName, c.Request.Method, route).Observe(time.Since(start).Seconds())
	}
}

// RecordPanic increments the panic counter for a route.
func RecordPanic(route string) {
	if route == "" {
		route = "unmatched"
	}
	panicsTotal.WithLabelValues(serviceName, route).Inc()
}

// MetricsHandler serves the Prometheus scrape endpoint.
func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
