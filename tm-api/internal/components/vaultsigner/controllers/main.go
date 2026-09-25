package controllers

import (
	"log"
	"net/http"
	"os"

	"admin-panel-dashboard/internal/components/accesslog"
	"admin-panel-dashboard/internal/components/vaultsigner/handlers"
	"admin-panel-dashboard/internal/components/vaultsigner/services"
	"admin-panel-dashboard/internal/components/vaultsigner/vaultclient"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverModels "admin-panel-dashboard/internal/server/models"

	vaultapi "github.com/hashicorp/vault/api"

	"github.com/gin-gonic/gin"
)

// requireVaultClient rejects every vault-signer request with a clear 503 when Vault isn't
// configured, instead of letting the request reach a handler and panic on a nil-pointer
// dereference deep inside the Vault SDK.
func requireVaultClient(vc *vaultapi.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if vc == nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error":   "vault_unavailable",
				"message": "Vault signer is not configured on this server.",
			})
			return
		}
		c.Next()
	}
}

// Init registers every vault-signer route (Section 6). Building the Vault
// client and parsing PERSONAL_ENVS happen once here, not per-request.
func Init(router *gin.Engine, s *serverModels.Server) {
	vc, err := vaultclient.NewClient()
	if err != nil {
		log.Printf("[vaultsigner] Vault client not configured, Vault-backed routes will return 503 until it is: %v\n", err)
	}

	personalEnvPrefixes, err := services.ParsePersonalEnvs(os.Getenv("PERSONAL_ENVS"))
	if err != nil {
		log.Fatalf("[vaultsigner] invalid PERSONAL_ENVS: %v\n", err)
	}

	h := handlers.NewHandler(s.AdminDB, vc, personalEnvPrefixes)

	apiV1 := router.Group("/api/v1")

	// Self-service — Trovo Admin or Organization member (Section 1/6).
	me := apiV1.Group("/me/vault-signer", middleware.AllowOrgOrTrovoAdminNormalized(s.AdminDB), requireVaultClient(vc))
	me.GET("/secrets", h.ListMySecrets)
	me.GET("/secrets/:secretId", h.GetMySecretValue)
	me.PUT("/secrets/:secretId", accesslog.Audit(s.AdminDB, models.EventSecretWrite, models.AccessCategorySecret, accesslog.Param("secretId")), h.PutMySecretValue)
	me.GET("/personal-envs", h.ListPersonalEnvs)
	me.POST("/personal-envs/:prefix", accesslog.Audit(s.AdminDB, models.EventPersonalEnvCreate, models.AccessCategorySecret, accesslog.Param("prefix")), h.CreatePersonalEnv)
	me.DELETE("/personal-envs/:prefix", accesslog.Audit(s.AdminDB, models.EventPersonalEnvDelete, models.AccessCategorySecret, accesslog.Param("prefix")), h.DeletePersonalEnv)

	// Admin-only — Trovo SuperAdmin (Section 1/6).
	admin := apiV1.Group("/admin/vault-signer", middleware.AuthenticateSuperAdmin(s.AdminDB))
	admin.GET("/managed-secrets", h.ListManagedSecrets)
	admin.POST("/managed-secrets", accesslog.Audit(s.AdminDB, models.EventManagedSecretRegister, models.AccessCategorySecret, accesslog.BodyField("label")), h.RegisterManagedSecret)
	admin.GET("/managed-secrets/:id/assignments", h.ListAssignments)
	admin.POST("/managed-secrets/:id/assignments", accesslog.Audit(s.AdminDB, models.EventSecretAssignmentGrant, models.AccessCategorySecret, accesslog.BodyField("ownerIdentifier")), h.CreateAssignment)
	admin.PATCH("/managed-secrets/:id/assignments/:assignmentId", accesslog.Audit(s.AdminDB, models.EventSecretAssignmentEdit, models.AccessCategorySecret, accesslog.Param("assignmentId")), h.EditAssignment)
	// Only DeleteAssignment actually touches Vault (it removes the signer on-chain via the
	// vault client) — guarded individually so the rest of the admin CRUD surface stays usable
	// during a Vault outage.
	admin.DELETE("/managed-secrets/:id/assignments/:assignmentId", requireVaultClient(vc), accesslog.Audit(s.AdminDB, models.EventSecretAssignmentRevoke, models.AccessCategorySecret, accesslog.Param("assignmentId")), h.DeleteAssignment)
	admin.GET("/audit-log", h.ListAuditLog)
}
