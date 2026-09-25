package handlers

import (
	"errors"
	"net/http"

	vaultsignermodels "admin-panel-dashboard/internal/components/vaultsigner/models"
	"admin-panel-dashboard/internal/components/vaultsigner/services"
	"admin-panel-dashboard/internal/components/vaultsigner/vaultclient"

	"github.com/gin-gonic/gin"
)

// resolveTrovoUsernameOrRespond runs Section 5b's gate. On failure it writes
// the friendly response itself and returns ok=false — every personal-secrets
// handler (list, create, delete) starts with this, uniformly.
func (h *Handler) resolveTrovoUsernameOrRespond(c *gin.Context) (username string, ok bool) {
	username, ok = services.ResolveTrovoUsername(c, h.db)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "trovo_wallet_not_linked",
			"message": services.FriendlyMessageNotLinked,
		})
		return "", false
	}
	return username, true
}

// ListPersonalEnvs — GET /api/v1/me/vault-signer/personal-envs (Section 6).
// Never a value — only prefix, label, derived key name, and exists.
func (h *Handler) ListPersonalEnvs(c *gin.Context) {
	username, ok := h.resolveTrovoUsernameOrRespond(c)
	if !ok {
		return
	}
	items, err := services.ListPersonalEnvs(c.Request.Context(), h.vc, h.personalEnvPrefixes, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

type createPersonalEnvRequest struct {
	NewValue string `json:"newValue" binding:"required"`
}

// CreatePersonalEnv — POST /api/v1/me/vault-signer/personal-envs/:prefix.
// Create-only — 409 if it already exists (Section 6).
func (h *Handler) CreatePersonalEnv(c *gin.Context) {
	username, ok := h.resolveTrovoUsernameOrRespond(c)
	if !ok {
		return
	}
	prefix := c.Param("prefix")
	if !h.isKnownPersonalEnvPrefix(prefix) {
		c.JSON(http.StatusNotFound, gin.H{"error": "unknown personal secret prefix"})
		return
	}
	var req createPersonalEnvRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "newValue is required"})
		return
	}

	vaultKey, version, err := services.CreatePersonalEnv(c.Request.Context(), h.vc, prefix, username, req.NewValue)
	if err != nil {
		if errors.Is(err, vaultclient.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "personal secret already exists"})
			return
		}
		writeErrorResponse(c, err)
		return
	}

	writeAuditLog(h.db, vaultsignermodels.VaultSignerAuditLog{
		Kind:              "personal_env_create",
		VaultKey:          stringPtr(vaultKey),
		ActorID:           username,
		ActorType:         c.GetString("user_type"),
		VaultVersionAfter: intPtr(version),
	})
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"prefix": prefix, "vaultKey": vaultKey, "vaultVersion": version}})
}

// DeletePersonalEnv — DELETE /api/v1/me/vault-signer/personal-envs/:prefix.
// Fully purges the secret — 404 if it doesn't exist (Section 6).
func (h *Handler) DeletePersonalEnv(c *gin.Context) {
	username, ok := h.resolveTrovoUsernameOrRespond(c)
	if !ok {
		return
	}
	prefix := c.Param("prefix")
	if !h.isKnownPersonalEnvPrefix(prefix) {
		c.JSON(http.StatusNotFound, gin.H{"error": "unknown personal secret prefix"})
		return
	}

	vaultKey, err := services.DeletePersonalEnv(c.Request.Context(), h.vc, prefix, username)
	if err != nil {
		if errors.Is(err, vaultclient.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "personal secret not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	writeAuditLog(h.db, vaultsignermodels.VaultSignerAuditLog{
		Kind:      "personal_env_delete",
		VaultKey:  stringPtr(vaultKey),
		ActorID:   username,
		ActorType: c.GetString("user_type"),
	})
	c.Status(http.StatusNoContent)
}
