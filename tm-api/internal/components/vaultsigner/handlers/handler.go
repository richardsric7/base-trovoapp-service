package handlers

import (
	"net/http"

	"admin-panel-dashboard/internal/components/vaultsigner/services"

	"github.com/gin-gonic/gin"
	vaultapi "github.com/hashicorp/vault/api"
	"gorm.io/gorm"
)

// Handler holds the dependencies every vaultsigner handler needs. It is
// constructed once in controllers.Init and shared across requests — the
// Vault client and the parsed PERSONAL_ENVS prefixes are both safe to reuse
// (Section 2/3).
type Handler struct {
	db                  *gorm.DB
	vc                  *vaultapi.Client
	personalEnvPrefixes []services.PersonalEnvPrefix
}

func NewHandler(db *gorm.DB, vc *vaultapi.Client, personalEnvPrefixes []services.PersonalEnvPrefix) *Handler {
	return &Handler{db: db, vc: vc, personalEnvPrefixes: personalEnvPrefixes}
}

// httpError is implemented by every internal/errors type (CustomError,
// ErrorInvalidPublicKey, ErrorBlockchainAccountNotActivated, ...) — reusing
// their HTTPCode()/JSONError() rather than inventing a second error
// envelope for this feature (Section 6).
type httpError interface {
	HTTPCode() int
	JSONError() gin.H
}

func writeErrorResponse(c *gin.Context, err error) {
	if he, ok := err.(httpError); ok {
		c.JSON(he.HTTPCode(), he.JSONError())
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func (h *Handler) isKnownPersonalEnvPrefix(prefix string) bool {
	for _, p := range h.personalEnvPrefixes {
		if p.Prefix == prefix {
			return true
		}
	}
	return false
}
