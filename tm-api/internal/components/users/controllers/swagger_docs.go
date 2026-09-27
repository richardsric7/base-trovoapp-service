package users

// This file holds swaggo doc comments for the routes registered directly in
// Init (main.go) as inline gin.HandlerFunc closures rather than named
// package functions. swag matches an @Router entry to its documented path
// and method regardless of which function the comment sits above, so these
// no-op stub functions exist purely to carry the doc comment block - the
// same pattern used in internal/components/vaultsigner/handlers/swagger_docs.go
// and internal/components/stakeholder/handlers/swagger_docs.go.
//
// Note: every route in this file is registered on the bare router (not the
// /api/v1 group used everywhere else), so its real path does not start with
// /api/v1 despite the @BasePath declared in main.go - the same pre-existing
// mismatch as /organizations/members/reset-password in the organizations
// component. The @Router path below is the literal path gin serves.

// @Summary Get a user's own wallet detail
// @Description Returns the caller's own Trovo Wallet user record (username/email/phone, balances metadata, KYC level, etc). targetUser is looked up and then compared against the caller's own token subject - a request for anyone else's targetUser is rejected with 401. Cached for 2 minutes per identifier. Reads from the shared TrovoWalletDB/P2P schema.
// @Tags Users
// @Produce json
// @Param targetUser path string true "Username, email, or user ID to look up (must resolve to the caller's own account)"
// @Success 200 {object} models.UserInfo
// @Failure 400 {object} map[string]string "error": "user not found or lookup failed"
// @Failure 401 {object} map[string]string "error": "unauthorized access"
// @Security JwtTokenAuth
// @Router /v1/users/detail/{targetUser} [get]
func getUserDetailDocs() {}

// @Summary Update the caller's own contact info
// @Description Updates the caller's own Trovo Wallet user record (e.g. contact phone) in the shared TrovoWalletDB/P2P schema, then broadcasts the change to any connected clients and invalidates cached lookups for the user.
// @Tags Users
// @Accept json
// @Produce json
// @Param body body models.UserUpdateInfo true "Fields to update on the caller's own account"
// @Success 200 {object} models.UserInfo
// @Failure 400 {object} map[string]string "error": "invalid payload or update failed"
// @Security JwtTokenAuth
// @Router /v1/users/update [put]
func updateUserDocs() {}

// @Summary Update a user's KYC level (admin)
// @Description Admin-only. Updates the KYC level and related fields of the user identified by targetUser in the shared TrovoWalletDB/P2P schema. The caller must be a wallet-side admin (AdminLevel > 0 on their own user record, not a tm-api AdminUser role) - a non-admin caller gets 401. Logged to the audit trail.
// @Tags Users
// @Accept json
// @Produce json
// @Param targetUser path string true "Username or email of the user whose KYC level is being changed"
// @Param body body models.KYCUpdateInfo true "New KYC fields"
// @Success 200 {object} models.UserInfo
// @Failure 400 {object} map[string]string "error": "invalid payload or update failed"
// @Failure 401 {object} map[string]string "error": "Unauthorized access. Only admins can access this endpoint"
// @Security JwtTokenAuth
// @Router /v1/admin/users/update/{targetUser} [put]
func adminUpdateUserKycDocs() {}

// @Summary Toggle the caller's own merchant online/offline status
// @Description Merchants only (KYCLevel must be non-zero, i.e. merchant-tier KYC). Flips the caller's own online/offline flag and broadcasts a userOnline/userOffline event to their connection streams. Logged to the audit trail.
// @Tags Users
// @Produce json
// @Success 200 {object} models.UserInfo
// @Failure 400 {object} map[string]string "error": "Offline/Online toggle is for merchants only, or toggle failed"
// @Security JwtTokenAuth
// @Router /v1/users/toggle [put]
func toggleUserOnlineDocs() {}
