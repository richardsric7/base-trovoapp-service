package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	vaultsignermodels "admin-panel-dashboard/internal/components/vaultsigner/models"
	"admin-panel-dashboard/internal/components/vaultsigner/services"
	"admin-panel-dashboard/internal/components/vaultsigner/vaultclient"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// loadOwnedSecret fetches the managed secret and the caller's own assignment
// on it, in one place — every me/vault-signer/secrets/:secretId handler
// needs both. Authorization is a direct comparison against OwnerRefID/
// OwnerType, already available from context with zero extra DB lookups
// beyond this one row (Section 1).
func (h *Handler) loadOwnedSecret(c *gin.Context, secretID string) (*vaultsignermodels.VaultSignerManagedSecret, *vaultsignermodels.VaultSignerAssignment, bool) {
	currentUserID := c.GetString("current_user_id")
	userType := c.GetString("user_type")

	var assignment vaultsignermodels.VaultSignerAssignment
	if err := h.db.Where("managed_secret_id = ? AND owner_ref_id = ? AND owner_type = ?", secretID, currentUserID, userType).
		First(&assignment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found, or not assigned to you"})
		return nil, nil, false
	}
	var secret vaultsignermodels.VaultSignerManagedSecret
	if err := h.db.First(&secret, "id = ?", secretID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "managed secret not found"})
		return nil, nil, false
	}
	return &secret, &assignment, true
}

// ListMySecrets — GET /api/v1/me/vault-signer/secrets (Section 6).
func (h *Handler) ListMySecrets(c *gin.Context) {
	currentUserID := c.GetString("current_user_id")
	userType := c.GetString("user_type")

	var assignments []vaultsignermodels.VaultSignerAssignment
	if err := h.db.Where("owner_ref_id = ? AND owner_type = ?", currentUserID, userType).Find(&assignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	items := make([]gin.H, 0, len(assignments))
	for _, a := range assignments {
		var secret vaultsignermodels.VaultSignerManagedSecret
		if err := h.db.First(&secret, "id = ?", a.ManagedSecretID).Error; err != nil {
			continue
		}
		items = append(items, gin.H{
			"managedSecretId": secret.ID,
			"label":           secret.Label,
			"position":        a.Position,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetMySecretValue — GET /api/v1/me/vault-signer/secrets/:secretId (Section 6).
// If this assignment has never claimed a position (Position is null), there
// is nothing in Vault yet to read — that's reported as "position": null,
// not an error.
func (h *Handler) GetMySecretValue(c *gin.Context) {
	secret, assignment, ok := h.loadOwnedSecret(c, c.Param("secretId"))
	if !ok {
		return
	}
	if assignment.Position == nil {
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"position": nil}})
		return
	}
	value, _, err := vaultclient.ReadValue(c.Request.Context(), h.vc, *secret, *assignment.Position-1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"position": *assignment.Position, "value": value}})
}

type putSecretRequest struct {
	NewValue string `json:"newValue" binding:"required"`
}

// PutMySecretValue — PUT /api/v1/me/vault-signer/secrets/:secretId. If this
// is the first time this assignment's owner has ever submitted a value, it
// claims the lowest unclaimed position first (services.ClaimPosition,
// Section 4/5e) — every CSV entry is already a real, pre-existing on-chain
// signer, so this claim always lands on a genuine key, meaning the write
// that follows always goes through the swap path below, not a first-time
// fallback. From there it branches into the on-chain swap flow for an
// active position, or a plain validated write for a spare one (Section
// 5d/6).
func (h *Handler) PutMySecretValue(c *gin.Context) {
	secret, assignment, ok := h.loadOwnedSecret(c, c.Param("secretId"))
	if !ok {
		return
	}
	var req putSecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "newValue is required"})
		return
	}

	ctx := c.Request.Context()

	if assignment.Position == nil {
		claimed, err := services.ClaimPosition(ctx, h.db, h.vc, assignment.ID)
		if err != nil {
			if errors.Is(err, services.ErrNoPositionAvailable) {
				c.JSON(http.StatusConflict, gin.H{"error": "no_position_available", "message": "Every signer position on this managed secret is already claimed."})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		assignment.Position = &claimed
	}
	position := *assignment.Position
	index := position - 1
	currentUserID := c.GetString("current_user_id")
	userType := c.GetString("user_type")

	if index < secret.ActiveSigningCount {
		result, err := services.SwapSigner(ctx, h.vc, *secret, index, req.NewValue)
		if err != nil {
			writeErrorResponse(c, err)
			return
		}
		kind := "signer_slot"
		var txHash, txStatus *string
		if result.Swapped {
			kind = "signer_swap"
			txHash = stringPtr(result.StellarTxHash)
			txStatus = stringPtr(result.StellarTxStatus)
		}
		writeAuditLog(h.db, vaultsignermodels.VaultSignerAuditLog{
			Kind:               kind,
			ManagedSecretID:    stringPtr(secret.ID),
			Position:           intPtr(position),
			ActorID:            currentUserID,
			ActorType:          userType,
			VaultVersionBefore: intPtr(result.VaultVersionBefore),
			VaultVersionAfter:  intPtr(result.VaultVersionAfter),
			StellarTxHash:      txHash,
			StellarTxStatus:    txStatus,
		})
		resp := gin.H{"vaultVersion": result.VaultVersionAfter}
		if result.Swapped {
			resp["stellarTxHash"] = txHash
			resp["stellarTxStatus"] = result.StellarTxStatus
			if result.SwapError != nil {
				resp["error"] = result.SwapError.Error()
			}
		}
		c.JSON(http.StatusOK, gin.H{"data": resp})
		return
	}

	// Spare slot — plain write, still through the Section 5c validation pipeline.
	if _, err := services.ValidateSignerValue(req.NewValue); err != nil {
		writeErrorResponse(c, err)
		return
	}
	versionBefore, versionAfter, err := vaultclient.WriteAtIndex(ctx, h.vc, *secret, index, strings.TrimSpace(req.NewValue))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	writeAuditLog(h.db, vaultsignermodels.VaultSignerAuditLog{
		Kind:               "signer_slot",
		ManagedSecretID:    stringPtr(secret.ID),
		Position:           intPtr(position),
		ActorID:            currentUserID,
		ActorType:          userType,
		VaultVersionBefore: intPtr(versionBefore),
		VaultVersionAfter:  intPtr(versionAfter),
	})
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"vaultVersion": versionAfter}})
}

// --- Admin endpoints (AuthenticateSuperAdmin-gated, Section 6) ---

type registerManagedSecretRequest struct {
	Label              string `json:"label" binding:"required"`
	VaultMount         string `json:"vaultMount" binding:"required"`
	VaultPath          string `json:"vaultPath" binding:"required"`
	VaultField         string `json:"vaultField" binding:"required"`
	WalletPublicKey    string `json:"walletPublicKey" binding:"required"`
	ActiveSigningCount int    `json:"activeSigningCount" binding:"required"`
}

// RegisterManagedSecret — POST /api/v1/admin/vault-signer/managed-secrets.
func (h *Handler) RegisterManagedSecret(c *gin.Context) {
	var req registerManagedSecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	secret := vaultsignermodels.VaultSignerManagedSecret{
		ID:                 uuid.NewString(),
		Label:              req.Label,
		VaultMount:         req.VaultMount,
		VaultPath:          req.VaultPath,
		VaultField:         req.VaultField,
		WalletPublicKey:    req.WalletPublicKey,
		ActiveSigningCount: req.ActiveSigningCount,
		CreatedAt:          time.Now(),
	}
	if err := h.db.Create(&secret).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": secret})
}

func (h *Handler) decorateAssignments(assignments []vaultsignermodels.VaultSignerAssignment) []gin.H {
	out := make([]gin.H, 0, len(assignments))
	for _, a := range assignments {
		label, err := services.ResolveOwnerLabel(h.db, a.OwnerRefID, a.OwnerType)
		if err != nil {
			label = ""
		}
		out = append(out, gin.H{
			"id":              a.ID,
			"managedSecretId": a.ManagedSecretID,
			"position":        a.Position,
			"ownerRefId":      a.OwnerRefID,
			"ownerType":       a.OwnerType,
			"ownerLabel":      label,
			"assignedAt":      a.AssignedAt,
		})
	}
	return out
}

// ListManagedSecrets — GET /api/v1/admin/vault-signer/managed-secrets, each
// with its current assignment rows (Section 6).
func (h *Handler) ListManagedSecrets(c *gin.Context) {
	var secrets []vaultsignermodels.VaultSignerManagedSecret
	if err := h.db.Find(&secrets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := make([]gin.H, 0, len(secrets))
	for _, s := range secrets {
		var assignments []vaultsignermodels.VaultSignerAssignment
		h.db.Where("managed_secret_id = ?", s.ID).Order("position").Find(&assignments)
		result = append(result, gin.H{"managedSecret": s, "assignments": h.decorateAssignments(assignments)})
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// ListAssignments — GET /api/v1/admin/vault-signer/managed-secrets/:id/assignments.
func (h *Handler) ListAssignments(c *gin.Context) {
	var assignments []vaultsignermodels.VaultSignerAssignment
	if err := h.db.Where("managed_secret_id = ?", c.Param("id")).Order("position").Find(&assignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": h.decorateAssignments(assignments)})
}

type ownerRequest struct {
	OwnerIdentifier string `json:"ownerIdentifier" binding:"required"`
	OwnerType       string `json:"ownerType" binding:"required"`
}

// CreateAssignment — POST /api/v1/admin/vault-signer/managed-secrets/:id/assignments.
// No position accepted — it starts NULL and is only claimed the first time
// the owner submits a real signer value (Section 4/5e).
func (h *Handler) CreateAssignment(c *gin.Context) {
	secretID := c.Param("id")
	var req ownerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ownerRefID, err := services.ResolveOwnerRefID(h.db, req.OwnerIdentifier, req.OwnerType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	assignment, err := services.CreateAssignment(h.db, secretID, ownerRefID, req.OwnerType)
	if err != nil {
		if errors.Is(err, services.ErrOwnerAlreadyAssigned) {
			c.JSON(http.StatusConflict, gin.H{"error": "owner_already_assigned", "message": "This owner already has an assignment on this managed secret."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	writeAuditLog(h.db, vaultsignermodels.VaultSignerAuditLog{
		Kind:            "assignment_create",
		ManagedSecretID: stringPtr(secretID),
		ActorID:         c.GetString("current_user_id"),
		ActorType:       c.GetString("user_type"),
	})
	label, _ := services.ResolveOwnerLabel(h.db, ownerRefID, req.OwnerType)
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{
		"id": assignment.ID, "managedSecretId": assignment.ManagedSecretID, "position": assignment.Position,
		"ownerRefId": assignment.OwnerRefID, "ownerType": assignment.OwnerType, "ownerLabel": label,
		"assignedAt": assignment.AssignedAt,
	}})
}

// EditAssignment — PATCH /api/v1/admin/vault-signer/managed-secrets/:id/assignments/:assignmentId.
// Reassigns the owner only — Position (null or claimed) is always retained
// exactly as-is; this never touches Vault or the chain (Section 6).
func (h *Handler) EditAssignment(c *gin.Context) {
	assignmentID := c.Param("assignmentId")
	var req ownerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ownerRefID, err := services.ResolveOwnerRefID(h.db, req.OwnerIdentifier, req.OwnerType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	assignment, err := services.EditAssignment(h.db, assignmentID, ownerRefID, req.OwnerType)
	if err != nil {
		if errors.Is(err, services.ErrOwnerAlreadyAssigned) {
			c.JSON(http.StatusConflict, gin.H{"error": "owner_already_assigned", "message": "This owner already has an assignment on this managed secret."})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	writeAuditLog(h.db, vaultsignermodels.VaultSignerAuditLog{
		Kind:            "assignment_edit",
		ManagedSecretID: stringPtr(assignment.ManagedSecretID),
		Position:        assignment.Position,
		ActorID:         c.GetString("current_user_id"),
		ActorType:       c.GetString("user_type"),
	})
	label, _ := services.ResolveOwnerLabel(h.db, ownerRefID, req.OwnerType)
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"id": assignment.ID, "managedSecretId": assignment.ManagedSecretID, "position": assignment.Position,
		"ownerRefId": assignment.OwnerRefID, "ownerType": assignment.OwnerType, "ownerLabel": label,
		"assignedAt": assignment.AssignedAt,
	}})
}

// DeleteAssignment — DELETE /api/v1/admin/vault-signer/managed-secrets/:id/assignments/:assignmentId
// (Section 5e/6).
func (h *Handler) DeleteAssignment(c *gin.Context) {
	assignmentID := c.Param("assignmentId")

	result, err := services.DeleteAssignment(c.Request.Context(), h.db, h.vc, assignmentID)
	if err != nil {
		if errors.Is(err, services.ErrCSVFloorViolation) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "assignment_not_found", "message": "No assignment found with that ID."})
			return
		}
		if errors.Is(err, vaultclient.ErrIndexOutOfRange) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		writeAuditLog(h.db, vaultsignermodels.VaultSignerAuditLog{
			Kind:            "assignment_delete",
			ActorID:         c.GetString("current_user_id"),
			ActorType:       c.GetString("user_type"),
			StellarTxStatus: stringPtr("failed"),
		})
		if he, ok := err.(httpError); ok { // e.g. ErrVaultUnreachableReverted
			c.JSON(he.HTTPCode(), he.JSONError())
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   "assignment_delete_failed",
			"message": "Could not remove the signer on-chain; the assignment was not deleted.",
		})
		return
	}

	writeAuditLog(h.db, vaultsignermodels.VaultSignerAuditLog{
		Kind:      "assignment_delete",
		ActorID:   c.GetString("current_user_id"),
		ActorType: c.GetString("user_type"),
		Position: func() *int {
			if result.DeletedPosition == 0 {
				return nil
			}
			return intPtr(result.DeletedPosition)
		}(),
		StellarTxHash:   stringPtr(result.StellarTxHash),
		StellarTxStatus: stringPtr(result.StellarTxStatus),
	})

	if result.DeletedPosition == 0 {
		// Never-claimed assignment — a plain row delete, nothing else happened.
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"claimed": false}})
		return
	}
	var onChainRemoval interface{}
	if result.OnChainAttempted {
		onChainRemoval = gin.H{"attempted": true, "stellarTxHash": result.StellarTxHash, "stellarTxStatus": result.StellarTxStatus}
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"claimed":         true,
		"deletedPosition": result.DeletedPosition,
		"renumberedCount": result.RenumberedCount,
		"vaultCollapsed":  result.VaultCollapsed,
		"onChainRemoval":  onChainRemoval,
	}})
}

// ListAuditLog — GET /api/v1/admin/vault-signer/audit-log.
func (h *Handler) ListAuditLog(c *gin.Context) {
	var rows []vaultsignermodels.VaultSignerAuditLog
	if err := h.db.Order("changed_at desc").Limit(200).Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}
