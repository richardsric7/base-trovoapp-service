package servicelinks

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type fakeDocumentStorage struct {
	upload       sharedconfig.DocumentUpload
	download     *sharedconfig.DocumentDownload
	deletedID    string
	uploadCalled int
}

func (f *fakeDocumentStorage) Upload(_ context.Context, input sharedconfig.DocumentUpload) (*sharedconfig.StoredDocument, error) {
	f.uploadCalled++
	f.upload = input
	return &sharedconfig.StoredDocument{
		ObjectID: uuid.NewString(), OriginalFilename: input.OriginalFilename,
		MimeType: input.MimeType, SizeBytes: int64(len(input.Content)), SHA256: input.SHA256,
	}, nil
}

func (f *fakeDocumentStorage) Open(_ context.Context, _ string) (*sharedconfig.DocumentDownload, error) {
	return f.download, nil
}

func (f *fakeDocumentStorage) Delete(_ context.Context, objectID string) error {
	f.deletedID = objectID
	return nil
}

func TestStakeholderDocumentUploadDownloadDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	storage := &fakeDocumentStorage{download: &sharedconfig.DocumentDownload{
		Body: io.NopCloser(strings.NewReader("private file")), OriginalFilename: "report.pdf",
		MimeType: "application/pdf", SizeBytes: 12,
	}}
	gc := &sharedconfig.GlobalConfig{StakeholderDocumentStorage: storage}
	router := gin.New()
	router.POST("/documents", postStakeholderDocumentHandler(gc))
	router.GET("/documents/:objectID", getStakeholderDocumentHandler(gc))
	router.DELETE("/documents/:objectID", deleteStakeholderDocumentHandler(gc))

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, err := w.CreateFormFile("document_file", "valuation.pdf")
	if err != nil {
		t.Fatal(err)
	}
	content := []byte("%PDF-1.4\nprivate report")
	_, _ = part.Write(content)
	_ = w.Close()
	req := httptest.NewRequest(http.MethodPost, "/documents", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	if response.Code != http.StatusCreated {
		t.Fatalf("upload status=%d body=%s", response.Code, response.Body.String())
	}
	if storage.uploadCalled != 1 || storage.upload.OriginalFilename != "valuation.pdf" || storage.upload.MimeType != "application/pdf" || !bytes.Equal(storage.upload.Content, content) || len(storage.upload.SHA256) != 64 {
		t.Fatalf("unexpected upload: %+v", storage.upload)
	}
	var payload struct {
		Data sharedconfig.StoredDocument `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil || payload.Data.ObjectID == "" {
		t.Fatalf("invalid response: %v %s", err, response.Body.String())
	}

	objectID := uuid.NewString()
	response = httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/documents/"+objectID, nil))
	if response.Code != http.StatusOK || response.Body.String() != "private file" || response.Header().Get("Content-Type") != "application/pdf" {
		t.Fatalf("download status=%d headers=%v body=%q", response.Code, response.Header(), response.Body.String())
	}

	response = httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodDelete, "/documents/"+objectID, nil))
	if response.Code != http.StatusNoContent || storage.deletedID != objectID {
		t.Fatalf("delete status=%d deleted=%q", response.Code, storage.deletedID)
	}
}

func TestStakeholderDocumentUploadRejectsUnsupportedContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	storage := &fakeDocumentStorage{}
	router := gin.New()
	router.POST("/documents", postStakeholderDocumentHandler(&sharedconfig.GlobalConfig{StakeholderDocumentStorage: storage}))

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, _ := w.CreateFormFile("document_file", "payload.txt")
	_, _ = part.Write([]byte("plain text"))
	_ = w.Close()
	req := httptest.NewRequest(http.MethodPost, "/documents", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	if response.Code != http.StatusBadRequest || storage.uploadCalled != 0 {
		t.Fatalf("status=%d body=%s uploadCalled=%d", response.Code, response.Body.String(), storage.uploadCalled)
	}
}
