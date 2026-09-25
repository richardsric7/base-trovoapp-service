package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path"
	"strings"
	"time"

	"admin-panel-dashboard/internal/trovosdk"
)

type StoredDocument struct {
	ObjectID         string `json:"object_id"`
	OriginalFilename string `json:"original_filename"`
	MimeType         string `json:"mime_type"`
	SizeBytes        int64  `json:"size_bytes"`
	SHA256           string `json:"sha256"`
}

type DocumentDownload struct {
	Body             io.ReadCloser
	OriginalFilename string
	MimeType         string
	SizeBytes        int64
}

type DocumentStorageClient interface {
	Upload(ctx context.Context, filename, mimeType string, content []byte) (*StoredDocument, error)
	Download(ctx context.Context, objectID string) (*DocumentDownload, error)
	Delete(ctx context.Context, objectID string) error
}

type WalletDocumentStorageClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewWalletDocumentStorageClient(link *trovosdk.ServiceLink) *WalletDocumentStorageClient {
	if link == nil {
		return &WalletDocumentStorageClient{client: &http.Client{Timeout: time.Minute}}
	}
	return &WalletDocumentStorageClient{
		baseURL: strings.TrimRight(link.ApiBaseUrl, "/"),
		apiKey:  link.ApiKey,
		client:  &http.Client{Timeout: time.Minute},
	}
}

func (c *WalletDocumentStorageClient) Upload(ctx context.Context, filename, mimeType string, content []byte) (*StoredDocument, error) {
	if err := c.configured(); err != nil {
		return nil, err
	}
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="document_file"; filename="%s"`, safeDocumentFilename(filename)))
	h.Set("Content-Type", mimeType)
	part, err := w.CreatePart(h)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(content); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(""), &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, walletStorageError(resp)
	}
	var payload struct {
		Data StoredDocument `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode Wallet document response: %w", err)
	}
	if strings.TrimSpace(payload.Data.ObjectID) == "" {
		return nil, fmt.Errorf("Wallet document response is missing object_id")
	}
	return &payload.Data, nil
}

func (c *WalletDocumentStorageClient) Download(ctx context.Context, objectID string) (*DocumentDownload, error) {
	if err := c.configured(); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint(objectID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, walletStorageError(resp)
	}
	return &DocumentDownload{
		Body: resp.Body, OriginalFilename: responseFilename(resp),
		MimeType: resp.Header.Get("Content-Type"), SizeBytes: resp.ContentLength,
	}, nil
}

func (c *WalletDocumentStorageClient) Delete(ctx context.Context, objectID string) error {
	if err := c.configured(); err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.endpoint(objectID), nil)
	if err != nil {
		return err
	}
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return walletStorageError(resp)
	}
	return nil
}

func (c *WalletDocumentStorageClient) configured() error {
	if c == nil || strings.TrimSpace(c.baseURL) == "" || strings.TrimSpace(c.apiKey) == "" {
		return NewHTTPError(http.StatusServiceUnavailable, "document storage is not configured")
	}
	return nil
}

func (c *WalletDocumentStorageClient) endpoint(objectID string) string {
	endpoint := c.baseURL + "/v1/trovo-api/stakeholder-documents"
	if objectID != "" {
		endpoint += "/" + path.Base(objectID)
	}
	return endpoint
}

func (c *WalletDocumentStorageClient) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-TW-SERVICE-LINK-API-KEY", c.apiKey)
	req.Header.Set("User-Agent", "Trovo-Admin-Dashboard")
	return c.client.Do(req)
}

func walletStorageError(resp *http.Response) error {
	message := "Wallet document storage request failed"
	var payload struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	_ = json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&payload)
	if payload.Message != "" {
		message = payload.Message
	} else if payload.Error != "" {
		message = payload.Error
	}
	if resp.StatusCode == http.StatusNotFound {
		return NewHTTPError(http.StatusNotFound, "document file not found")
	}
	if resp.StatusCode == http.StatusBadRequest {
		return NewHTTPError(http.StatusBadRequest, message)
	}
	if resp.StatusCode == http.StatusRequestEntityTooLarge {
		return NewHTTPError(http.StatusRequestEntityTooLarge, message)
	}
	if resp.StatusCode == http.StatusServiceUnavailable {
		return NewHTTPError(http.StatusServiceUnavailable, message)
	}
	return NewHTTPError(http.StatusBadGateway, message)
}

func responseFilename(resp *http.Response) string {
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err == nil {
		return safeDocumentFilename(params["filename"])
	}
	return "document"
}

func safeDocumentFilename(filename string) string {
	filename = path.Base(strings.ReplaceAll(strings.TrimSpace(filename), `\`, "/"))
	filename = strings.NewReplacer(`"`, "", "\r", "", "\n", "").Replace(filename)
	if filename == "." || filename == "" {
		return "document"
	}
	return filename
}
