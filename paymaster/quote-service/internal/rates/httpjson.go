package rates

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// http-json:
//
//	{
//	  "type": "http-json",
//	  "url": "https://api.example.com/ticker?pair=USDTNGN",
//	  "path": "data.ticker.last",                  // dot path; numeric segments index arrays
//	  "headers": {"X-Api-Key": "vault:EXAMPLE_API_KEY"}
//	}
//
// Reads one number (JSON number or numeric string) from a JSON HTTP API:
// an exchange ticker, an FX feed, or the stablecoin issuer's rate endpoint.
// Header values of the form "vault:KEY" are read from the configuration
// secret, so API keys never appear in the pair definitions.
type httpJSONSource struct {
	url     string
	path    []string
	headers map[string]string
	client  *http.Client
}

func init() {
	Register("http-json", func(def json.RawMessage, env Env) (Source, error) {
		var c struct {
			URL     string            `json:"url"`
			Path    string            `json:"path"`
			Headers map[string]string `json:"headers"`
		}
		if err := json.Unmarshal(def, &c); err != nil {
			return nil, err
		}
		if !strings.HasPrefix(c.URL, "https://") && !strings.HasPrefix(c.URL, "http://") {
			return nil, fmt.Errorf(`"url" must be an http(s) URL`)
		}
		if c.Path == "" {
			return nil, fmt.Errorf(`"path" is required`)
		}
		headers := map[string]string{}
		for k, v := range c.Headers {
			resolved, err := resolveSecret(v, env)
			if err != nil {
				return nil, err
			}
			headers[k] = resolved
		}
		client := env.HTTP
		if client == nil {
			client = http.DefaultClient
		}
		return &httpJSONSource{url: c.URL, path: strings.Split(c.Path, "."), headers: headers, client: client}, nil
	})
}

func (s *httpJSONSource) Fetch(ctx context.Context) (decimal.Decimal, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url, nil)
	if err != nil {
		return decimal.Zero, err
	}
	req.Header.Set("Accept", "application/json")
	for k, v := range s.headers {
		req.Header.Set(k, v)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return decimal.Zero, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return decimal.Zero, err
	}
	if res.StatusCode != http.StatusOK {
		return decimal.Zero, fmt.Errorf("HTTP %d", res.StatusCode)
	}
	var doc interface{}
	dec := json.NewDecoder(strings.NewReader(string(body)))
	dec.UseNumber()
	if err := dec.Decode(&doc); err != nil {
		return decimal.Zero, fmt.Errorf("response is not JSON: %w", err)
	}
	return extractNumber(doc, s.path)
}

func extractNumber(doc interface{}, path []string) (decimal.Decimal, error) {
	cur := doc
	for i, seg := range path {
		switch node := cur.(type) {
		case map[string]interface{}:
			next, ok := node[seg]
			if !ok {
				return decimal.Zero, fmt.Errorf("path %s: key %q not found", strings.Join(path[:i+1], "."), seg)
			}
			cur = next
		case []interface{}:
			idx, err := strconv.Atoi(seg)
			if err != nil || idx < 0 || idx >= len(node) {
				return decimal.Zero, fmt.Errorf("path %s: index %q out of range", strings.Join(path[:i+1], "."), seg)
			}
			cur = node[idx]
		default:
			return decimal.Zero, fmt.Errorf("path %s: cannot descend into a scalar", strings.Join(path[:i+1], "."))
		}
	}
	switch v := cur.(type) {
	case json.Number:
		return decimal.NewFromString(v.String())
	case string:
		d, err := decimal.NewFromString(strings.TrimSpace(v))
		if err != nil {
			return decimal.Zero, fmt.Errorf("value %q is not a number", v)
		}
		return d, nil
	default:
		return decimal.Zero, fmt.Errorf("value at path is not a number")
	}
}
