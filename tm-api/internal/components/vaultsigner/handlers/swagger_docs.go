package handlers

// listMySecretsDocs documents GET /me/vault-signer/secrets.
// @Summary      List my signer secrets and assigned positions
// @Tags         Vault Signer - Self Service
// @Security     JwtTokenAuth
// @Security     OrganizationAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      503 {object} map[string]interface{} "Vault signer is not configured on this server"
// @Router       /me/vault-signer/secrets [get]
func listMySecretsDocs() {}

// getMySecretValueDocs documents GET /me/vault-signer/secrets/{secretId}.
// @Summary      Get the current value of my signer slot
// @Tags         Vault Signer - Self Service
// @Security     JwtTokenAuth
// @Security     OrganizationAuth
// @Produce      json
// @Param        secretId path string true "Managed secret ID"
// @Success      200 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Failure      503 {object} map[string]interface{} "Vault signer is not configured on this server"
// @Router       /me/vault-signer/secrets/{secretId} [get]
func getMySecretValueDocs() {}

// putMySecretValueDocs documents PUT /me/vault-signer/secrets/{secretId}.
// @Summary      Submit a new value for my signer slot
// @Description  For an active position (< the managed secret's active signing count), this triggers the on-chain swap flow when the old value is a valid keypair, the new value is a valid+activated keypair, and the two differ. Otherwise it's a plain validated write.
// @Tags         Vault Signer - Self Service
// @Security     JwtTokenAuth
// @Security     OrganizationAuth
// @Accept       json
// @Produce      json
// @Param        secretId path string true "Managed secret ID"
// @Param        payload body putSecretRequest true "New signer value"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Failure      422 {object} map[string]interface{}
// @Failure      503 {object} map[string]interface{} "Vault signer is not configured on this server"
// @Router       /me/vault-signer/secrets/{secretId} [put]
func putMySecretValueDocs() {}

// listPersonalEnvsDocs documents GET /me/vault-signer/personal-envs.
// @Summary      List my personal secret prefixes
// @Description  Never returns a value — only prefix, label, derived key name, and whether it exists.
// @Tags         Vault Signer - Personal Secrets
// @Security     JwtTokenAuth
// @Security     OrganizationAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{} "org member has not linked a Trovo wallet"
// @Failure      503 {object} map[string]interface{} "Vault signer is not configured on this server"
// @Router       /me/vault-signer/personal-envs [get]
func listPersonalEnvsDocs() {}

// createPersonalEnvDocs documents POST /me/vault-signer/personal-envs/{prefix}.
// @Summary      Create my personal secret
// @Description  Create only — not an upsert. 409 if it already exists. Never echoes the value back.
// @Tags         Vault Signer - Personal Secrets
// @Security     JwtTokenAuth
// @Security     OrganizationAuth
// @Accept       json
// @Produce      json
// @Param        prefix path string true "Personal env prefix, e.g. AUTO_APPROVE"
// @Param        payload body createPersonalEnvRequest true "New value"
// @Success      201 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{} "org member has not linked a Trovo wallet"
// @Failure      409 {object} map[string]interface{}
// @Failure      422 {object} map[string]interface{}
// @Failure      503 {object} map[string]interface{} "Vault signer is not configured on this server"
// @Router       /me/vault-signer/personal-envs/{prefix} [post]
func createPersonalEnvDocs() {}

// deletePersonalEnvDocs documents DELETE /me/vault-signer/personal-envs/{prefix}.
// @Summary      Delete my personal secret
// @Description  Fully purges the secret from Vault (not a soft delete). 404 if it doesn't exist.
// @Tags         Vault Signer - Personal Secrets
// @Security     JwtTokenAuth
// @Security     OrganizationAuth
// @Param        prefix path string true "Personal env prefix, e.g. AUTO_APPROVE"
// @Success      204 "No Content"
// @Failure      403 {object} map[string]interface{} "org member has not linked a Trovo wallet"
// @Failure      404 {object} map[string]interface{}
// @Failure      503 {object} map[string]interface{} "Vault signer is not configured on this server"
// @Router       /me/vault-signer/personal-envs/{prefix} [delete]
func deletePersonalEnvDocs() {}

// listManagedSecretsDocs documents GET /admin/vault-signer/managed-secrets.
// @Summary      List managed signer secrets with their current assignments
// @Tags         Vault Signer - Admin
// @Security     JwtTokenAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /admin/vault-signer/managed-secrets [get]
func listManagedSecretsDocs() {}

// registerManagedSecretDocs documents POST /admin/vault-signer/managed-secrets.
// @Summary      Register a new managed signer secret
// @Tags         Vault Signer - Admin
// @Security     JwtTokenAuth
// @Accept       json
// @Produce      json
// @Param        payload body registerManagedSecretRequest true "Managed secret"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /admin/vault-signer/managed-secrets [post]
func registerManagedSecretDocs() {}

// listAssignmentsDocs documents GET /admin/vault-signer/managed-secrets/{id}/assignments.
// @Summary      List current assignments for one managed secret
// @Tags         Vault Signer - Admin
// @Security     JwtTokenAuth
// @Produce      json
// @Param        id path string true "Managed secret ID"
// @Success      200 {object} map[string]interface{}
// @Router       /admin/vault-signer/managed-secrets/{id}/assignments [get]
func listAssignmentsDocs() {}

// createAssignmentDocs documents POST /admin/vault-signer/managed-secrets/{id}/assignments.
// @Summary      Create an assignment (auto-assigned position)
// @Description  No position is supplied — it starts null and is only claimed the first time the owner submits a real signer value. (ManagedSecretID, OwnerRefID) is unique — 409 if this owner already has an assignment on this managed secret.
// @Tags         Vault Signer - Admin
// @Security     JwtTokenAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Managed secret ID"
// @Param        payload body ownerRequest true "Owner identifier (username for trovo_admin, email for org_member) and owner type"
// @Success      201 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Failure      409 {object} map[string]interface{} "owner already assigned to this managed secret"
// @Router       /admin/vault-signer/managed-secrets/{id}/assignments [post]
func createAssignmentDocs() {}

// editAssignmentDocs documents PATCH /admin/vault-signer/managed-secrets/{id}/assignments/{assignmentId}.
// @Summary      Reassign the owner of an assignment
// @Description  Position (null or already claimed) is always retained exactly as-is. Never touches Vault or the chain — the new owner picks up the slot the normal way, via PUT on their own next write, which claims a position first if none has been claimed yet. (ManagedSecretID, OwnerRefID) is unique — 409 if the new owner already has a different assignment on this managed secret.
// @Tags         Vault Signer - Admin
// @Security     JwtTokenAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Managed secret ID"
// @Param        assignmentId path string true "Assignment ID"
// @Param        payload body ownerRequest true "New owner identifier and owner type"
// @Success      200 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Failure      409 {object} map[string]interface{} "owner already assigned to this managed secret"
// @Router       /admin/vault-signer/managed-secrets/{id}/assignments/{assignmentId} [patch]
func editAssignmentDocs() {}

// deleteAssignmentDocs documents DELETE /admin/vault-signer/managed-secrets/{id}/assignments/{assignmentId}.
// @Summary      Delete an assignment
// @Description  If no position was ever claimed, this is a plain row delete. Otherwise it renumbers every later position down by one and collapses the Vault CSV. For an active position with enough redundancy, also removes the signer on-chain before the Vault write — reversed automatically if the Vault write then fails.
// @Tags         Vault Signer - Admin
// @Security     JwtTokenAuth
// @Produce      json
// @Param        id path string true "Managed secret ID"
// @Param        assignmentId path string true "Assignment ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{} "signer index out of range"
// @Failure      404 {object} map[string]interface{} "no assignment found with that ID"
// @Failure      409 {object} map[string]interface{} "would bring the CSV below the minimum entry count"
// @Failure      502 {object} map[string]interface{} "on-chain removal failed, or Vault was unreachable and all changes were reverted"
// @Failure      503 {object} map[string]interface{} "Vault signer is not configured on this server"
// @Router       /admin/vault-signer/managed-secrets/{id}/assignments/{assignmentId} [delete]
func deleteAssignmentDocs() {}

// listAuditLogDocs documents GET /admin/vault-signer/audit-log.
// @Summary      View vault signer change history
// @Tags         Vault Signer - Admin
// @Security     JwtTokenAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /admin/vault-signer/audit-log [get]
func listAuditLogDocs() {}
