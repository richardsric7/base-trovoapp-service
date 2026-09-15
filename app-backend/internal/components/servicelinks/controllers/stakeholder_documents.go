package servicelinks

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	servicelinkServices "trovo-wallet-api/internal/components/servicelinks/services"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const maxStakeholderDocumentSize int64 = 10 << 20

var stakeholderDocumentMimeTypes = map[string]struct{}{
	"application/pdf": {},
	"image/jpeg":      {},
	"image/png":       {},
}

// postStakeholderDocumentHandler stores a private stakeholder document.
// @Summary Upload private stakeholder document
// @Tags servicelinks
// @Accept multipart/form-data
// @Produce json
// @Param document_file formData file true "PDF, JPEG, or PNG document (max 10MB)"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /v1/trovo-api/stakeholder-documents [post]
func postStakeholderDocumentHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if gc.StakeholderDocumentStorage == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "stakeholder document storage is not configured"})
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxStakeholderDocumentSize+(1<<20))
		if err := c.Request.ParseMultipartForm(maxStakeholderDocumentSize); err != nil {
			var maxBytesErr *http.MaxBytesError
			if errors.As(err, &maxBytesErr) {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "document exceeds 10MB"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart request or document exceeds 10MB"})
			return
		}
		file, header, err := c.Request.FormFile("document_file")
		if errors.Is(err, http.ErrMissingFile) {
			file, header, err = c.Request.FormFile("documentFile")
		}
		if err != nil {
			var maxBytesErr *http.MaxBytesError
			if errors.As(err, &maxBytesErr) {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "document exceeds 10MB"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "document_file is required"})
			return
		}
		defer file.Close()

		content, err := io.ReadAll(io.LimitReader(file, maxStakeholderDocumentSize+1))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read document"})
			return
		}
		if int64(len(content)) > maxStakeholderDocumentSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "document exceeds 10MB"})
			return
		}
		mimeType := http.DetectContentType(content)
		if _, ok := stakeholderDocumentMimeTypes[mimeType]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported document type; use PDF, JPEG, or PNG"})
			return
		}
		filename := filepath.Base(strings.ReplaceAll(strings.TrimSpace(header.Filename), `\`, "/"))
		filename = strings.NewReplacer(`"`, "", "\r", "", "\n", "").Replace(filename)
		if filename == "." || filename == "" {
			filename = "document"
		}
		hash := sha256.Sum256(content)
		stored, err := gc.StakeholderDocumentStorage.Upload(c.Request.Context(), sharedconfig.DocumentUpload{
			Content: content, MimeType: mimeType, OriginalFilename: filename, SHA256: hex.EncodeToString(hash[:]),
		})
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "unable to store stakeholder document"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"data": stored})
	}
}

// getStakeholderDocumentHandler streams a private stakeholder document.
// @Summary Download private stakeholder document
// @Tags servicelinks
// @Produce application/octet-stream
// @Param objectID path string true "Storage object ID"
// @Success 200 {file} binary
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /v1/trovo-api/stakeholder-documents/{objectID} [get]
func getStakeholderDocumentHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if gc.StakeholderDocumentStorage == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "stakeholder document storage is not configured"})
			return
		}
		objectID := c.Param("objectID")
		if _, err := uuid.Parse(objectID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document object id"})
			return
		}
		document, err := gc.StakeholderDocumentStorage.Open(c.Request.Context(), objectID)
		if err != nil {
			if errors.Is(err, storage.ErrObjectNotExist) {
				c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
				return
			}
			c.JSON(http.StatusBadGateway, gin.H{"error": "unable to load stakeholder document"})
			return
		}
		defer document.Body.Close()
		filename := strings.NewReplacer("\"", "", "\r", "", "\n", "").Replace(document.OriginalFilename)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		c.Header("Content-Type", document.MimeType)
		c.Header("Content-Length", strconv.FormatInt(document.SizeBytes, 10))
		c.Status(http.StatusOK)
		_, _ = io.Copy(c.Writer, document.Body)
	}
}

// deleteStakeholderDocumentHandler deletes an unreferenced private stakeholder document.
// @Summary Delete private stakeholder document
// @Tags servicelinks
// @Param objectID path string true "Storage object ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /v1/trovo-api/stakeholder-documents/{objectID} [delete]
func deleteStakeholderDocumentHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if gc.StakeholderDocumentStorage == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "stakeholder document storage is not configured"})
			return
		}
		objectID := c.Param("objectID")
		if _, err := uuid.Parse(objectID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document object id"})
			return
		}
		if err := gc.StakeholderDocumentStorage.Delete(c.Request.Context(), objectID); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "unable to delete stakeholder document"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func requireActiveServiceLink(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceLink, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid service link"})
			return
		}
		if serviceLink.Inactive != 0 || serviceLink.Suspended != 0 || serviceLink.Verified == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "service link is not active"})
			return
		}
		c.Next()
	}
}
