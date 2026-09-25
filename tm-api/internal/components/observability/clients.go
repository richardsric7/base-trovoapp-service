package observability

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// httpClient is shared by every upstream call. A single client reuses
// connections rather than opening one per request.
var httpClient = &http.Client{Timeout: upstreamTimeout}

// getJSON performs a GET and decodes the body into out. headers may be nil.
func getJSON(ctx context.Context, endpoint string, headers map[string]string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("upstream returned %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// ---- Prometheus -------------------------------------------------------------

// promInstantResult is the shape Prometheus returns for an instant query.
type promInstantResult struct {
	Data struct {
		Result []struct {
			Metric map[string]string `json:"metric"`
			Value  [2]any            `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// promRangeResult is the shape Prometheus returns for a range query.
type promRangeResult struct {
	Data struct {
		Result []struct {
			Metric map[string]string `json:"metric"`
			Values [][2]any          `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// queryInstant runs a single-point PromQL query.
//
// The query is always built by this package - never taken from the request -
// so a TM user cannot run arbitrary PromQL against the monitoring stack.
func (c Config) queryInstant(ctx context.Context, promQL string) (promInstantResult, error) {
	var out promInstantResult
	endpoint := c.PrometheusURL + "/api/v1/query?query=" + url.QueryEscape(promQL)
	err := getJSON(ctx, endpoint, nil, &out)
	return out, err
}

// queryRange runs a PromQL query over a window, returning a series per metric.
func (c Config) queryRange(ctx context.Context, promQL string, window time.Duration, step time.Duration) (promRangeResult, error) {
	var out promRangeResult
	end := time.Now()
	start := end.Add(-window)
	endpoint := fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%d&end=%d&step=%d",
		c.PrometheusURL, url.QueryEscape(promQL), start.Unix(), end.Unix(), int(step.Seconds()))
	err := getJSON(ctx, endpoint, nil, &out)
	return out, err
}

// sampleValue pulls the numeric value out of Prometheus's [timestamp, "value"]
// pair. Prometheus encodes the value as a string, so it needs parsing.
func sampleValue(pair [2]any) (float64, bool) {
	if len(pair) < 2 {
		return 0, false
	}
	s, ok := pair[1].(string)
	if !ok {
		return 0, false
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || isNotNumber(f) {
		return 0, false
	}
	return f, true
}

// isNotNumber reports NaN, which Prometheus returns for an undefined
// expression (e.g. a rate over a series with no samples). Rendering NaN as 0
// would be a lie; callers treat it as "no data".
func isNotNumber(f float64) bool { return f != f }

// sampleTime pulls the timestamp out of a Prometheus sample pair.
func sampleTime(pair [2]any) (time.Time, bool) {
	if len(pair) < 1 {
		return time.Time{}, false
	}
	ts, ok := pair[0].(float64)
	if !ok {
		return time.Time{}, false
	}
	return time.Unix(int64(ts), 0).UTC(), true
}

// ---- Loki -------------------------------------------------------------------

type lokiResult struct {
	Data struct {
		Result []struct {
			Stream map[string]string `json:"stream"`
			Values [][2]string       `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// queryLoki runs a LogQL query over a window.
//
// As with Prometheus, the selector is always constructed here. The only
// user-supplied value that reaches Loki is a request id, and that is validated
// against a strict pattern first (see validRequestID) - without which a
// crafted value could inject LogQL operators.
func (c Config) queryLoki(ctx context.Context, logQL string, window time.Duration, limit int) (lokiResult, error) {
	var out lokiResult
	end := time.Now()
	start := end.Add(-window)
	endpoint := fmt.Sprintf("%s/loki/api/v1/query_range?query=%s&start=%d&end=%d&limit=%d&direction=forward",
		c.LokiURL, url.QueryEscape(logQL), start.UnixNano(), end.UnixNano(), limit)
	err := getJSON(ctx, endpoint, nil, &out)
	return out, err
}

// ---- GlitchTip --------------------------------------------------------------

// glitchTipIssue is the subset of GlitchTip's issue payload the UI needs.
type glitchTipIssue struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Culprit   string `json:"culprit"`
	Level     string `json:"level"`
	Status    string `json:"status"`
	Count     any    `json:"count"`
	FirstSeen string `json:"firstSeen"`
	LastSeen  string `json:"lastSeen"`
	Permalink string `json:"permalink"`
	Project   struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"project"`
}

// fetchIssues lists open issues across the organisation, newest activity first.
func (c Config) fetchIssues(ctx context.Context, limit int, project, query string) ([]glitchTipIssue, error) {
	endpoint := fmt.Sprintf("%s/api/0/organizations/%s/issues/?limit=%d&sort=last_seen",
		c.GlitchTipURL, url.PathEscape(c.GlitchTipOrg), limit)
	if project != "" {
		endpoint += "&project=" + url.QueryEscape(project)
	}
	if query != "" {
		endpoint += "&query=" + url.QueryEscape(query)
	}

	var out []glitchTipIssue
	err := getJSON(ctx, endpoint, map[string]string{
		"Authorization": "Bearer " + c.GlitchTipToken,
	}, &out)
	return out, err
}

// issueCount normalises GlitchTip's event count, which arrives as either a
// JSON number or a string depending on version.
func issueCount(v any) int64 {
	switch t := v.(type) {
	case float64:
		return int64(t)
	case string:
		n, err := strconv.ParseInt(strings.TrimSpace(t), 10, 64)
		if err != nil {
			return 0
		}
		return n
	default:
		return 0
	}
}
