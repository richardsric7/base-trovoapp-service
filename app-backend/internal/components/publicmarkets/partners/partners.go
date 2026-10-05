// Package partners is how Public Markets talks to its external parties:
// Custodians (who reach CSCS for us), Dealing Members (who trade on NGX and
// FMDQ) and the price feed. Each has an interface, a REST client following
// the integration specification's illustrative contracts, and - in the
// services package - a mock used until a partner's real API is agreed.
package partners

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// CustodianInstruction is a creation (BUY) or redemption (SELL) instruction
// (integration spec §4.3.1-4.3.2).
type CustodianInstruction struct {
	InstructionID             string `json:"instructionId"`
	AssetCode                 string `json:"assetCode"`
	ISIN                      string `json:"isin"`
	Quantity                  string `json:"quantity"`
	Side                      string `json:"side"`
	ReferenceDate             string `json:"referenceDate"`
	TrovotechAccountReference string `json:"trovotechAccountReference"`
}

// DealingOrder is a trade execution instruction (§5.2).
type DealingOrder struct {
	OrderID            string `json:"orderId"`
	AssetCode          string `json:"assetCode"`
	Market             string `json:"market"`
	Side               string `json:"side"`
	Quantity           string `json:"quantity"`
	OrderType          string `json:"orderType"`
	TrovotechReference string `json:"trovotechReference"`
}

// Ack is a partner's synchronous acknowledgement (202 Accepted).
type Ack struct {
	Status     string
	StatusCode int
}

// Position is what a Custodian holds for the pool (§4.3.3).
type Position struct {
	AssetCode string
	UnitsHeld decimal.Decimal
	AsOf      time.Time
}

// Quote is a reference price (§8).
type Quote struct {
	ISIN       string
	Price      decimal.Decimal
	AsOf       time.Time
	MarketOpen bool
}

// Custodian is a Custodian's API (instructions and position queries).
type Custodian interface {
	SendInstruction(ctx context.Context, side string, in CustodianInstruction) (Ack, error)
	Position(ctx context.Context, assetCode string) (Position, error)
}

// DealingMember is a Dealing Member's order API.
type DealingMember interface {
	PlaceOrder(ctx context.Context, o DealingOrder) (Ack, error)
}

// PriceFeed is the NGX/FMDQ reference price vendor.
type PriceFeed interface {
	Quote(ctx context.Context, isin, market string) (Quote, error)
}

// PermanentError is a refusal the partner will repeat (a 4xx): retrying
// cannot help, so the instruction is escalated at once.
type PermanentError struct {
	StatusCode int
	Body       string
}

func (e *PermanentError) Error() string {
	return fmt.Sprintf("partner refused the request (%d): %s", e.StatusCode, e.Body)
}

// IsPermanent reports whether err is a PermanentError.
func IsPermanent(err error) bool {
	var p *PermanentError
	return errors.As(err, &p)
}

// Credentials resolves a credentials reference: "env:NAME" reads NAME
// ("keyId:secret" for HMAC); NAME_CERT and NAME_KEY (PEM files) give a
// client certificate for mTLS. Vault paths are injected into the
// environment at deploy time, like every other managed secret.
type Credentials struct {
	KeyID  string
	Secret string
	Cert   *tls.Certificate
}

// ResolveCredentials loads the credentials a reference names.
func ResolveCredentials(ref string) (Credentials, error) {
	var c Credentials
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return c, nil
	}
	name := strings.TrimPrefix(ref, "env:")
	if strings.HasPrefix(ref, "vault://") {
		// vault://path#FIELD - the deploy injects FIELD into the environment
		if i := strings.LastIndex(ref, "#"); i > 0 {
			name = ref[i+1:]
		} else {
			return c, fmt.Errorf("credentials reference %q names no field (vault://path#FIELD)", ref)
		}
	}
	if v := strings.TrimSpace(os.Getenv(name)); v != "" {
		if i := strings.Index(v, ":"); i > 0 {
			c.KeyID, c.Secret = v[:i], v[i+1:]
		} else {
			c.Secret = v
		}
	}
	certFile, keyFile := os.Getenv(name+"_CERT"), os.Getenv(name+"_KEY")
	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return c, fmt.Errorf("loading the client certificate for %s: %w", name, err)
		}
		c.Cert = &cert
	}
	return c, nil
}

// Sign is the HMAC-SHA256 signature used across Public Markets: hex of
// HMAC(secret, timestamp + "." + body).
func Sign(secret, timestamp string, body []byte) string {
	m := hmac.New(sha256.New, []byte(secret))
	m.Write([]byte(timestamp))
	m.Write([]byte("."))
	m.Write(body)
	return hex.EncodeToString(m.Sum(nil))
}

// VerifySignature checks Sign in constant time.
func VerifySignature(secret, timestamp string, body []byte, signature string) bool {
	if secret == "" || signature == "" {
		return false
	}
	want := Sign(secret, timestamp, body)
	return hmac.Equal([]byte(want), []byte(strings.ToLower(strings.TrimSpace(signature))))
}

// RESTClient calls a partner's REST API with the agreed authentication.
type RESTClient struct {
	BaseURL    string
	AuthScheme string // HMAC | MTLS | NONE
	Creds      Credentials
	HTTP       *http.Client
}

// NewRESTClient builds a client for a partner configuration.
func NewRESTClient(baseURL, authScheme, credentialsRef string) (*RESTClient, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, errors.New("the partner has no base URL")
	}
	creds, err := ResolveCredentials(credentialsRef)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}}
	if strings.EqualFold(authScheme, "MTLS") {
		if creds.Cert == nil {
			return nil, errors.New("mTLS is configured but no client certificate was found")
		}
		transport.TLSClientConfig.Certificates = []tls.Certificate{*creds.Cert}
	}
	return &RESTClient{BaseURL: strings.TrimRight(baseURL, "/"), AuthScheme: strings.ToUpper(authScheme), Creds: creds,
		HTTP: &http.Client{Timeout: 30 * time.Second, Transport: transport}}, nil
}

func (c *RESTClient) do(ctx context.Context, method, path, idempotencyKey string, in, out interface{}) (int, error) {
	var body []byte
	if in != nil {
		var err error
		if body, err = json.Marshal(in); err != nil {
			return 0, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", idempotencyKey)
	}
	if c.AuthScheme == "HMAC" {
		ts := strconv.FormatInt(time.Now().Unix(), 10)
		req.Header.Set("X-Trovotech-Timestamp", ts)
		req.Header.Set("Authorization", fmt.Sprintf("HMAC %s:%s", c.Creds.KeyID, Sign(c.Creds.Secret, ts, body)))
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return 0, err // network: retryable
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests {
		return resp.StatusCode, fmt.Errorf("partner error %d: %s", resp.StatusCode, trim(string(raw), 200))
	}
	if resp.StatusCode >= 400 {
		return resp.StatusCode, &PermanentError{StatusCode: resp.StatusCode, Body: trim(string(raw), 300)}
	}
	if out != nil && len(raw) > 0 {
		if err := json.Unmarshal(raw, out); err != nil {
			return resp.StatusCode, fmt.Errorf("unreadable partner response: %w", err)
		}
	}
	return resp.StatusCode, nil
}

func trim(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

// RESTCustodian is the Custodian contract of §4.
type RESTCustodian struct{ *RESTClient }

// SendInstruction posts a creation or redemption instruction.
func (c RESTCustodian) SendInstruction(ctx context.Context, side string, in CustodianInstruction) (Ack, error) {
	path := "/v1/instructions/creation"
	if strings.EqualFold(side, "SELL") {
		path = "/v1/instructions/redemption"
	}
	var out struct {
		Status string `json:"status"`
	}
	code, err := c.do(ctx, http.MethodPost, path, in.InstructionID, in, &out)
	return Ack{Status: out.Status, StatusCode: code}, err
}

// Position queries the units held for an asset.
func (c RESTCustodian) Position(ctx context.Context, assetCode string) (Position, error) {
	var out struct {
		AssetCode string    `json:"assetCode"`
		UnitsHeld string    `json:"unitsHeld"`
		AsOf      time.Time `json:"asOf"`
	}
	if _, err := c.do(ctx, http.MethodGet, "/v1/positions/"+assetCode, "", nil, &out); err != nil {
		return Position{}, err
	}
	units, err := decimal.NewFromString(out.UnitsHeld)
	if err != nil {
		return Position{}, fmt.Errorf("invalid unitsHeld %q", out.UnitsHeld)
	}
	return Position{AssetCode: out.AssetCode, UnitsHeld: units, AsOf: out.AsOf}, nil
}

// RESTDealingMember is the Dealing Member contract of §5.
type RESTDealingMember struct{ *RESTClient }

// PlaceOrder posts a trade execution instruction.
func (d RESTDealingMember) PlaceOrder(ctx context.Context, o DealingOrder) (Ack, error) {
	var out struct {
		Status string `json:"status"`
	}
	code, err := d.do(ctx, http.MethodPost, "/v1/orders", o.OrderID, o, &out)
	return Ack{Status: out.Status, StatusCode: code}, err
}

// RESTPriceFeed is the price vendor contract of §8.
type RESTPriceFeed struct{ *RESTClient }

// Quote fetches an instrument's reference price.
func (p RESTPriceFeed) Quote(ctx context.Context, isin, market string) (Quote, error) {
	var out struct {
		ISIN       string    `json:"isin"`
		Price      string    `json:"price"`
		AsOf       time.Time `json:"asOf"`
		MarketOpen bool      `json:"marketOpen"`
	}
	if _, err := p.do(ctx, http.MethodGet, "/v1/quotes/"+isin, "", nil, &out); err != nil {
		return Quote{}, err
	}
	price, err := decimal.NewFromString(out.Price)
	if err != nil || !price.IsPositive() {
		return Quote{}, fmt.Errorf("invalid price %q for %s", out.Price, isin)
	}
	return Quote{ISIN: out.ISIN, Price: price, AsOf: out.AsOf, MarketOpen: out.MarketOpen}, nil
}
