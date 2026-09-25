package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"admin-panel-dashboard/internal/components/stakeholder/models"
	"admin-panel-dashboard/internal/trovosdk"
)

type DistributionPayoutRegistration struct {
	DistributionID   string `json:"distribution_id"`
	TokenizedAssetID string `json:"tokenized_asset_id"`
	Amount           string `json:"amount"`
	Currency         string `json:"currency"`
}

type DistributionPayoutClient interface {
	Register(ctx context.Context, registration DistributionPayoutRegistration) (*models.DistributionPayoutResponse, error)
	Get(ctx context.Context, distributionID string) (*models.DistributionPayoutResponse, error)
}

type WalletDistributionPayoutClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewWalletDistributionPayoutClient(link *trovosdk.ServiceLink) *WalletDistributionPayoutClient {
	if link == nil {
		return &WalletDistributionPayoutClient{client: &http.Client{Timeout: 15 * time.Second}}
	}
	return &WalletDistributionPayoutClient{
		baseURL: strings.TrimRight(link.ApiBaseUrl, "/"), apiKey: link.ApiKey,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *WalletDistributionPayoutClient) Register(ctx context.Context, registration DistributionPayoutRegistration) (*models.DistributionPayoutResponse, error) {
	if err := c.configured(); err != nil {
		return nil, err
	}
	body, err := json.Marshal(registration)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(""), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req, false)
}

func (c *WalletDistributionPayoutClient) Get(ctx context.Context, distributionID string) (*models.DistributionPayoutResponse, error) {
	if err := c.configured(); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint(distributionID), nil)
	if err != nil {
		return nil, err
	}
	return c.do(req, true)
}

func (c *WalletDistributionPayoutClient) configured() error {
	if c == nil || strings.TrimSpace(c.baseURL) == "" || strings.TrimSpace(c.apiKey) == "" {
		return NewHTTPError(http.StatusServiceUnavailable, "Wallet payout integration is not configured")
	}
	return nil
}

func (c *WalletDistributionPayoutClient) endpoint(distributionID string) string {
	endpoint := c.baseURL + "/v1/trovo-api/stakeholder-distributions"
	if strings.TrimSpace(distributionID) != "" {
		endpoint += "/" + path.Base(distributionID)
	}
	return endpoint
}

func (c *WalletDistributionPayoutClient) do(req *http.Request, allowNotFound bool) (*models.DistributionPayoutResponse, error) {
	req.Header.Set("X-TW-SERVICE-LINK-API-KEY", c.apiKey)
	req.Header.Set("User-Agent", "Trovo-Admin-Dashboard")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, NewHTTPError(http.StatusServiceUnavailable, "Wallet payout service is unavailable")
	}
	defer resp.Body.Close()
	if allowNotFound && resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, walletPayoutError(resp)
	}
	var payload struct {
		Data models.DistributionPayoutResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode Wallet payout response: %w", err)
	}
	if strings.TrimSpace(payload.Data.DistributionID) == "" {
		return nil, NewHTTPError(http.StatusBadGateway, "Wallet payout response is missing distribution_id")
	}
	return &payload.Data, nil
}

func walletPayoutError(resp *http.Response) error {
	message := "Wallet payout request failed"
	var payload struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	_ = json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&payload)
	if strings.TrimSpace(payload.Message) != "" {
		message = payload.Message
	} else if strings.TrimSpace(payload.Error) != "" {
		message = payload.Error
	}
	switch resp.StatusCode {
	case http.StatusBadRequest, http.StatusUnprocessableEntity:
		return NewHTTPError(http.StatusUnprocessableEntity, message)
	case http.StatusConflict:
		return NewHTTPError(http.StatusConflict, message)
	case http.StatusUnauthorized, http.StatusForbidden:
		return NewHTTPError(http.StatusBadGateway, "Wallet rejected the configured service link")
	case http.StatusNotFound:
		return NewHTTPError(http.StatusBadGateway, message)
	default:
		return NewHTTPError(http.StatusServiceUnavailable, message)
	}
}
