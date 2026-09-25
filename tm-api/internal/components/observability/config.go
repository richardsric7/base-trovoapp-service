// Package observability is the admin API's gateway to the platform's
// monitoring stack. Trovo Manager screens call these endpoints; the admin API
// queries Prometheus, Loki and GlitchTip server-side over the internal network
// and returns plain JSON.
//
// The browser never talks to those services and never sees their credentials.
// Every query is fixed here - no user-supplied PromQL or LogQL reaches an
// upstream - so a TM user cannot turn this into an arbitrary query console
// against the monitoring stack.
package observability

import (
	"os"
	"strings"
	"time"
)

// Config is read once at startup from the environment. Everything is optional:
// a missing upstream disables its screens rather than breaking the module, so
// a deployment that has not wired up (say) Loki still gets health and issues.
type Config struct {
	PrometheusURL  string
	LokiURL        string
	GlitchTipURL   string
	GlitchTipToken string
	GlitchTipOrg   string

	// GlitchTipPublicURL is where a browser can reach the crash reporter, as
	// distinct from GlitchTipURL which is the internal address this service
	// calls. They differ: the admin API talks to http://glitchtip-web:8000 on
	// the container network, but a "view in the crash reporter" link handed to
	// a user must be a public hostname their browser can resolve.
	GlitchTipPublicURL string

	GrafanaURL string

	// Services maps a display name to the base URL the admin API uses to reach
	// that service's /ready endpoint on the internal network.
	//
	// This is configured rather than derived because container names differ
	// between environments and do not follow one convention - on dev the P2P
	// container is `p2p-dev` while the others are `<name>-api-dev`. Guessing
	// would silently show a healthy service as down.
	Services []ServiceTarget
}

// ServiceTarget is one backend the health screen reports on.
type ServiceTarget struct {
	Name string `json:"name" example:"admin-api"`
	URL  string `json:"-"`
}

// upstreamTimeout bounds every call to Prometheus, Loki or GlitchTip. These
// screens are diagnostic: a slow answer is far better than a hung request, and
// an admin looking at a health page during an incident should not wait.
const upstreamTimeout = 8 * time.Second

// LoadConfig reads the module's configuration from the environment.
//
// OBSERVABILITY_SERVICES is a comma-separated list of name=url pairs, e.g.
//
//	admin-api=http://admin-api-dev:8080,wallet-api=http://wallet-api-dev:8080
func LoadConfig() Config {
	c := Config{
		PrometheusURL:      strings.TrimRight(strings.TrimSpace(os.Getenv("PROMETHEUS_URL")), "/"),
		LokiURL:            strings.TrimRight(strings.TrimSpace(os.Getenv("LOKI_URL")), "/"),
		GlitchTipURL:       strings.TrimRight(strings.TrimSpace(os.Getenv("GLITCHTIP_URL")), "/"),
		GlitchTipToken:     strings.TrimSpace(os.Getenv("GLITCHTIP_API_TOKEN")),
		GlitchTipOrg:       strings.TrimSpace(os.Getenv("GLITCHTIP_ORG")),
		GlitchTipPublicURL: strings.TrimRight(strings.TrimSpace(os.Getenv("GLITCHTIP_PUBLIC_URL")), "/"),
		GrafanaURL:         strings.TrimRight(strings.TrimSpace(os.Getenv("GRAFANA_URL")), "/"),
	}

	for _, pair := range strings.Split(os.Getenv("OBSERVABILITY_SERVICES"), ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		name, url, found := strings.Cut(pair, "=")
		name, url = strings.TrimSpace(name), strings.TrimRight(strings.TrimSpace(url), "/")
		if !found || name == "" || url == "" {
			continue
		}
		c.Services = append(c.Services, ServiceTarget{Name: name, URL: url})
	}
	return c
}

// HasPrometheus reports whether metrics-backed screens can work.
func (c Config) HasPrometheus() bool { return c.PrometheusURL != "" }

// HasLoki reports whether the request-tracing screen can work.
func (c Config) HasLoki() bool { return c.LokiURL != "" }

// HasGlitchTip reports whether the issues screen can work. All three values are
// needed: a token without an org slug cannot address any endpoint.
func (c Config) HasGlitchTip() bool {
	return c.GlitchTipURL != "" && c.GlitchTipToken != "" && c.GlitchTipOrg != ""
}

// PublicGlitchTipURL returns the address to put in a link for a user. Falls
// back to the internal URL only when no public one is configured, so a
// misconfigured deployment shows an obviously wrong link rather than silently
// omitting the link entirely.
func (c Config) PublicGlitchTipURL() string {
	if c.GlitchTipPublicURL != "" {
		return c.GlitchTipPublicURL
	}
	return c.GlitchTipURL
}
