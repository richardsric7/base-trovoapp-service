package sharedconfig

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"

	cloudStorage "cloud.google.com/go/storage"
	firebaseStorage "firebase.google.com/go/storage"
	"github.com/google/uuid"
)

type DocumentUpload struct {
	Content          []byte
	MimeType         string
	OriginalFilename string
	SHA256           string
}

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

type DocumentStorage interface {
	Upload(ctx context.Context, input DocumentUpload) (*StoredDocument, error)
	Open(ctx context.Context, objectID string) (*DocumentDownload, error)
	Delete(ctx context.Context, objectID string) error
}

type GCSPrivateDocumentStorage struct {
	client *firebaseStorage.Client
	bucket string
	prefix string
}

func NewGCSPrivateDocumentStorage(client *firebaseStorage.Client, bucket, prefix string) *GCSPrivateDocumentStorage {
	return &GCSPrivateDocumentStorage{client: client, bucket: strings.TrimSpace(bucket), prefix: strings.Trim(strings.TrimSpace(prefix), "/")}
}

func (s *GCSPrivateDocumentStorage) Upload(ctx context.Context, input DocumentUpload) (*StoredDocument, error) {
	if s == nil || s.client == nil || s.bucket == "" {
		return nil, errors.New("stakeholder document storage is not configured")
	}
	bucket, err := s.client.Bucket(s.bucket)
	if err != nil {
		return nil, fmt.Errorf("open stakeholder document bucket: %w", err)
	}
	objectID := uuid.NewString()
	writer := bucket.Object(s.objectName(objectID)).NewWriter(ctx)
	writer.ContentType = input.MimeType
	writer.CacheControl = "private, no-store"
	writer.Metadata = map[string]string{
		"original-filename": input.OriginalFilename,
		"sha256":            input.SHA256,
	}
	if _, err := io.Copy(writer, bytes.NewReader(input.Content)); err != nil {
		_ = writer.CloseWithError(err)
		return nil, fmt.Errorf("upload stakeholder document: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("finalize stakeholder document upload: %w", err)
	}
	return &StoredDocument{
		ObjectID: objectID, OriginalFilename: input.OriginalFilename, MimeType: input.MimeType,
		SizeBytes: int64(len(input.Content)), SHA256: input.SHA256,
	}, nil
}

func (s *GCSPrivateDocumentStorage) Open(ctx context.Context, objectID string) (*DocumentDownload, error) {
	if s == nil || s.client == nil || s.bucket == "" {
		return nil, errors.New("stakeholder document storage is not configured")
	}
	if _, err := uuid.Parse(objectID); err != nil {
		return nil, errors.New("invalid stakeholder document object id")
	}
	bucket, err := s.client.Bucket(s.bucket)
	if err != nil {
		return nil, fmt.Errorf("open stakeholder document bucket: %w", err)
	}
	object := bucket.Object(s.objectName(objectID))
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return nil, err
	}
	reader, err := object.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	return &DocumentDownload{
		Body: reader, OriginalFilename: attrs.Metadata["original-filename"],
		MimeType: attrs.ContentType, SizeBytes: attrs.Size,
	}, nil
}

func (s *GCSPrivateDocumentStorage) Delete(ctx context.Context, objectID string) error {
	if s == nil || s.client == nil || s.bucket == "" {
		return errors.New("stakeholder document storage is not configured")
	}
	if _, err := uuid.Parse(objectID); err != nil {
		return errors.New("invalid stakeholder document object id")
	}
	bucket, err := s.client.Bucket(s.bucket)
	if err != nil {
		return fmt.Errorf("open stakeholder document bucket: %w", err)
	}
	err = bucket.Object(s.objectName(objectID)).Delete(ctx)
	if errors.Is(err, cloudStorage.ErrObjectNotExist) {
		return nil
	}
	return err
}

func (s *GCSPrivateDocumentStorage) objectName(objectID string) string {
	if s.prefix == "" {
		return objectID
	}
	return path.Join(s.prefix, objectID)
}
