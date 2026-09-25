package handlers

import (
	server "admin-panel-dashboard/internal/components/general/services"
	"admin-panel-dashboard/internal/components/organizations/services"
	stakeholderServices "admin-panel-dashboard/internal/components/stakeholder/services"
	usermetricsDB "admin-panel-dashboard/internal/components/usermetrics/db"
	usermetrics "admin-panel-dashboard/internal/components/usermetrics/services"
	userServices "admin-panel-dashboard/internal/components/users/services"
	p2pErrors "admin-panel-dashboard/internal/errors"
	"admin-panel-dashboard/internal/mail"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverModels "admin-panel-dashboard/internal/server/models"
	serverResponse "admin-panel-dashboard/internal/server/response"
	"admin-panel-dashboard/internal/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupUserPassword sets up the user's password after OTP verification
// @Summary Setup root user password
// @Description Sets up the root user's password after OTP verification. Requires organization ID, user ID, and OTP in path parameters, and password in request body.
// @Tags organizations
// @Accept json
// @Produce json
// @Param invite_id path string true "invite_id"
// @Param request body models.SetupUserPasswordRequest true "Password setup request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /organizations/setup-password/{invite_id} [post]
func SetupUserPassword(server *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get path parameters
		inviteID := c.Param("invite_id")
		if inviteID == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Missing invite_id path parameter"})
			return
		}
		//VALIDATE OTP too
		var req models.SetupUserPasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		// Validate request body
		if strings.TrimSpace(req.Password) == "" || len(req.OTP) <= 0 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Password and OTP is required"})
			return
		}

		token, err := setupPasswordService(c, server, inviteID, req)
		if err != nil {
			log.Printf("[ORGANIZATION] Error setting up password: %v", err)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "Password setup successful", gin.H{"token": token}, nil)
	}
}

func setupPasswordService(c *gin.Context, server *serverModels.Server, inviteID string, req models.SetupUserPasswordRequest) (string, error) {
	invite, err := services.ValidateInviteIDExpiryAndOtp(server, inviteID, req)
	if err != nil {
		slog.Error("[ORGANIZATION] Error checking invite", "error", err)
		return "", fmt.Errorf("invalid or expired invite")
	}

	active, err := services.IsStatusActive(server, invite.Email)
	if err != nil {
		slog.Error("[ORGANIZATION] Error checking member status", "error", err)
		return "", fmt.Errorf("internal server error")
	}
	if active {
		slog.Error("[ORGANIZATION] Member already active", "invite_id", inviteID)
		return "", fmt.Errorf("member already active, please login")
	}
	// IsEmailStatusValid
	//isEmailStatusAccepted, err := services.IsEmailStatusValid(server, inviteID)
	//if err != nil {
	//	slog.Error("[ORGANIZATION] Error checking email status", "error", err)
	//	return "", fmt.Errorf("internal server error")
	//}
	//if !isEmailStatusAccepted {
	//	slog.Error("[ORGANIZATION] Email status not accepted", "invite_id", inviteID)
	//	return "", fmt.Errorf("email not verified, please validate your email first")
	//}

	member, err := validateInviteForPasswordReset(server, inviteID)
	if err != nil {
		slog.Error("[ORGANIZATION] Error validating invite", "error", err)
		return "", err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		slog.Error("[ORGANIZATION] Error hashing password", "error", err)
		return "", fmt.Errorf("internal server error")
	}

	now := time.Now()
	if err := completePasswordSetup(server.AdminDB, invite, &member, hashedPassword, now); err != nil {
		log.Printf("[ORGANIZATION] Error completing password setup: %v", err)
		slog.Error("[ORGANIZATION] Error completing password setup", "error", err)
		return "", err
	}

	// Generate JWT token
	token, err := middleware.GenerateOrganizationTokenWithDB(server.AdminDB, &member)
	if err != nil {
		log.Printf("[ORGANIZATION] Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to generate authentication token"})
		return "", fmt.Errorf("internal server error")
	}
	return token, nil
}

func completePasswordSetup(db *gorm.DB, invite *models.OrganizationInvite, member *models.OrganizationMember, hashedPassword string, now time.Time) error {
	return db.Transaction(func(tx *gorm.DB) error {
		member.Password = hashedPassword
		member.Status = models.OrganizationStatusActive
		member.UpdatedAt = now
		member.VerificationOTP = ""
		member.OTPExpiresAt = time.Time{}
		member.KYCStatus = models.PASSWORD_RESET_KYC_STATUS

		if err := tx.Save(&member).Error; err != nil {
			return fmt.Errorf("failed to activate organization member: %w", err)
		}

		if err := tx.Model(&models.OrganizationInvite{}).
			Where("id = ?", invite.ID).
			Updates(map[string]interface{}{
				"status":                  models.InviteStatusAccepted,
				"email_validation_status": models.EndInviteStatusValidEmail,
				"email_validated_at":      now,
				"updated_at":              now,
			}).Error; err != nil {
			return fmt.Errorf("failed to accept organization invite: %w", err)
		}

		// Completing the root-admin invite is the organization activation event.
		// A regular member accepting an invite must not change organization state.
		if invite.Role == string(models.OrganizationSuperAdmin) {
			result := tx.Model(&models.Organization{}).
				Where("id = ? AND status = ?", invite.OrganizationID, models.OrganizationStatusPending).
				Updates(map[string]interface{}{
					"status":     models.OrganizationStatusActive,
					"updated_at": now,
				})
			if result.Error != nil {
				return fmt.Errorf("failed to activate organization: %w", result.Error)
			}
			if result.RowsAffected == 0 {
				var organization models.Organization
				if err := tx.Select("status").First(&organization, "id = ?", invite.OrganizationID).Error; err != nil {
					return fmt.Errorf("failed to find organization during activation: %w", err)
				}
				if organization.Status != models.OrganizationStatusActive {
					return fmt.Errorf("organization cannot be activated from status %s", organization.Status)
				}
			}
		}

		return nil
	})
}

func validateInviteForPasswordReset(server *serverModels.Server, inviteID string) (models.OrganizationMember, error) {
	// Get organization member
	var member models.OrganizationMember
	if err := server.AdminDB.Where("invite_id = ? AND status = ?", inviteID, models.InviteStatusPending).First(&member).Error; err != nil {
		slog.Error("Failed to find organization member", "error", err, "invite_id", inviteID)
		return models.OrganizationMember{}, fmt.Errorf("invalid user or invite not found")
	}
	//if member.Status != models.EndInviteStatusValidEmail {
	//	return models.OrganizationMember{}, fmt.Errorf("please validate your email first")
	//}
	return member, nil
}

// ResendPasswordSetupToken resends the password setup token
// @Summary Resend password setup token
// @Description Resends the password setup token for a verified organization
// @Tags organizations

//// @Accept json
//// @Produce json
//// @Param request body models.ResendPasswordSetupRequest true "Resend token request"
//// @Success 200 {object} map[string]string
//// @Failure 400 {object} models.ErrorResponse
//// @Failure 500 {object} models.ErrorResponse
//// @Router /organizations/resend-password-setup [post]
//func ResendPasswordSetupToken(server *serverModels.Server) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		var req models.ResendPasswordSetupRequest
//		if err := c.ShouldBindJSON(&req); err != nil {
//			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
//			return
//		}
//
//		token, err := services.SendPasswordSetupToken(server, req.Email)
//		if err != nil {
//			log.Printf("[ORGANIZATION] Error resending password setup token: %v", err)
//			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
//			return
//		}
//
//		serverResponse.JSON(c, http.StatusOK, "Password setup token resent successfully", gin.H{"token": token}, nil)
//	}
//}

// InviteOrganizationMember invites a member to an organization
// @Summary Invite organization member
// @Description Invites a member to join the organization. Requires JWT token in Authorization header: 'Authorization: Bearer <token>'
// @Tags organizations
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Accept json
// @Produce json
// @Param request body models.OrganizationInviteRequest true "Member invitation request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /organizations/members [post]
func InviteOrganizationMembe(server *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get organization details from JWT token
		orgID := c.GetString("organization_id")
		orgEmail := c.GetString("organization_email")
		orgType := c.GetString("organization_type")
		memberID := c.GetString("member_id")
		memberEmail := c.GetString("member_email")
		memberRole := c.GetString("member_role")
		organizationName := c.GetString("organization_name")

		if orgID == "" || orgEmail == "" || orgType == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
			return
		}

		// Check if the inviter is a root user (super admin)
		if memberRole != string(models.OrganizationSuperAdmin) {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Only root users can invite new members"})
			return
		}

		// Parse request body
		var req models.OrganizationInviteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		inviteOrAdminPayload := &models.InviteOrgMemberPayload{
			OrganizationName: organizationName,
			OrganizationType: orgType,
			AdminID:          memberID,
			//AdminFirstName:   req.AdminFirstName,
			//AdminLastName:    req.AdminLastName,
			AdminEmail: memberEmail,
		}
		// Call service to handle invitation
		invite, err := services.InviteOrganizationMember(server, orgID, &req, inviteOrAdminPayload)
		if err != nil {
			log.Printf("[ORGANIZATION] Error inviting member: %v", err)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("You have successfully sent an invitation to %s to join Trovo Manager as %s",
				memberEmail,
				models.OrganizationMemberRole),
			"invite_id":       invite.InviteID,
			"organization_id": invite.OrganizationID,
		})
		//serverResponse.JSON(c, http.StatusOK, "Invitation sent successfully", nil, nil)
	}
}

// InitiateWalletLinkHandler (Step 1) links Trovo Manager with Trovo Wallet in the DB only. Allowed once per user.
func InitiateWalletLinkHandler(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LinkWalletRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		walletUsername := strings.TrimSpace(strings.ToLower(req.Username))
		if walletUsername != "" && !strings.Contains(walletUsername, "@trovo") {
			walletUsername = walletUsername + "@trovo"
		}

		err := services.LinkWalletOnly(s, c, walletUsername)
		if err != nil {
			log.Printf("Failed to link wallet: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Wallet linked successfully"})
	}
}

// RequestWalletLinkAuthHandler (Step 1) starts authorization for an unlinked organization member.
// @Summary Request wallet authorization (Step 1)
// @Description Starts wallet linking for the authenticated organization member. Returns authId, deeplink, and QR code. Complete approval in Trovo Wallet, then call the verify endpoint.
// @Tags organizations
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param request body models.WalletLinkAuthorizationRequest true "Trovo Wallet username"
// @Success 200 {object} trovosdk.TrovoWalletAuthorizationData
// @Failure 400 {object} models.ErrorResponse
// @Router /organizations/members/link-wallet/request-auth [post]
func RequestWalletLinkAuthHandler(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.WalletLinkAuthorizationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "wallet_username is required"})
			return
		}
		authData, err := services.SendWalletLinkAuthorizationRequest(s, c, req.WalletUsername)
		if err != nil {
			log.Printf("Failed to send wallet auth request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, authData)
	}
}

// RequestWalletLinkLoginHandler sends a login request for the linked wallet; returns loginId, deeplink, QR (for use with verify-login to get token).
func RequestWalletLinkLoginHandler(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		loginData, err := services.SendWalletLinkLoginRequest(s, c)
		if err != nil {
			log.Printf("Failed to send wallet login request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, loginData)
	}
}

// VerifyWalletLinkLoginHandler (Step 3) verifies the login request and returns tokens (VerifyLoginRequest).
func VerifyWalletLinkLoginHandler(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			LoginID string `json:"login_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "login_id is required"})
			return
		}

		resp, err := services.VerifyWalletLinkLoginRequest(s, c, req.LoginID)
		if err != nil {
			log.Printf("Failed to verify wallet login: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// VerifyWalletLinkHandler (Step 2) verifies authorization and persists the member wallet link.
// @Summary Verify wallet authorization (Step 2)
// @Description Verifies approval using auth_id and links the Trovo Wallet to the authenticated organization member.
// @Tags organizations
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param request body object{auth_id=string,wallet_username=string} true "Verification details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Router /organizations/members/link-wallet/verify [post]
func VerifyWalletLinkHandler(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberID := c.GetString("member_id")
		if memberID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "organization member context required"})
			return
		}

		var req struct {
			AuthID         string `json:"auth_id" binding:"required"`
			WalletUsername string `json:"wallet_username" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		err := services.VerifyWalletLink(s, memberID, req.WalletUsername, req.AuthID)
		if err != nil {
			log.Printf("Failed to verify wallet link: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Wallet linked successfully"})
	}
}

// InviteOrganizationMember invites a new member to an organization
// @Summary Invite organization member
// @Description Invites a member to join the organization. Supports both Trovo admin and organization admin authentication.
// @Tags organizations
// @Security JwtTokenAuth
// @Security OrganizationAuth
// @Accept json
// @Produce json
// @Param request body models.OrganizationInviteRequest true "Member invitation request. Only trovo_admin must provide organization_id in the request body, while organization_admin will use their token to supply the org_id."
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /organizations/members [post]
func InviteOrganizationMember(server *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Bind the incoming JSON payload
		var req models.OrganizationInviteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		// 2. Check organization role permission (org members only, admins bypass)
		if err := middleware.CheckOrganizationPermission(c, string(models.OrganizationSuperAdmin)); err != nil {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Error: err.Error()})
			return
		}

		// 3. Get organization ID (from token for org members, from request for admins)
		orgID, err := middleware.GetOrganizationIDFromContext(c)
		if err != nil {
			// For Trovo admins, try to get org_id from request body
			if req.OrganizationID != "" {
				orgID = req.OrganizationID
			} else {
				c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
				return
			}
		}

		// 4. Get normalized user info
		inviterID := c.GetString("current_user_id")
		inviterEmail := c.GetString("current_user_email")

		// 5. Get organization details
		var org models.Organization
		if err := server.AdminDB.First(&org, "id = ?", orgID).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Organization not found"})
			return
		}
		organizationName := org.Name
		orgType := string(org.Type)
		// 4. Execute the common business logic
		inviteOrAdminPayload := &models.InviteOrgMemberPayload{
			OrganizationName: organizationName,
			OrganizationType: orgType,
			AdminID:          inviterID,
			AdminEmail:       inviterEmail,
		}

		invite, err := services.InviteOrganizationMember(server, orgID, &req, inviteOrAdminPayload)
		if err != nil {
			log.Printf("[ORGANIZATION] Error inviting member to org %s: %v", orgID, err)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		// 5. Return a successful response
		c.JSON(http.StatusOK, gin.H{
			"message":         fmt.Sprintf("An invitation has been successfully sent to %s to join %s.", req.Email, organizationName),
			"invite_id":       invite.InviteID,
			"organization_id": invite.OrganizationID,
		})
	}
}

type InviteOrgAdminPayload struct {
	OrganizationName string `json:"organization_name" binding:"required" example:"Acme Asset Management Ltd"`
	AdminFirstName   string `json:"admin_first_name" binding:"required" example:"John"`
	AdminLastName    string `json:"admin_last_name" binding:"required" example:"Doe"`
	AdminEmail       string `json:"admin_email" binding:"required,email" example:"john.doe@acme.com"`
	// StakeholderType is provided via query param and set server-side; hidden from Swagger body
	StakeholderType string `json:"-" swaggerignore:"true"`
	Address         string `json:"address" example:"123 Main St"`
	Country         string `json:"country" example:"USA"`
	FeePercent      string `json:"fee_percent" example:"0.5"`
	FeeFixed        string `json:"fee_fixed" example:"10.00"`
	// StakeholderID is set internally after partner creation, not provided by user
	StakeholderID *uint64 `json:"-"`
	// Level is the compliance requirement tier this organisation onboards at
	// (see PRD-kyc-compliance-templates.md). Optional; defaults to 1 (Basic).
	Level *int `json:"level" example:"1"`
}

// validateStakeholderLinkage validates stakeholder linkage using existing partner DB functions
func validateStakeholderLinkage(stakeholderType *string, stakeholderID *uint64, trovoDB *gorm.DB) error {
	// If both are nil, no linkage is being set - this is valid
	if stakeholderType == nil && stakeholderID == nil {
		return nil
	}

	// If one is provided but not the other, this is an error
	if stakeholderType == nil && stakeholderID != nil {
		return fmt.Errorf("stakeholder_type is required when stakeholder_id is provided")
	}
	if stakeholderType != nil && stakeholderID == nil {
		return fmt.Errorf("stakeholder_id is required when stakeholder_type is provided")
	}

	// Validate stakeholder_type value
	validTypes := []string{
		models.StakeholderTypeAssetManager,
		models.StakeholderTypeIssuingHouse,
		models.StakeholderTypeCustodian,
		models.StakeholderTypeLegal,
		models.StakeholderTypeRatingAgency,
		models.StakeholderTypeTrustee,
		models.StakeholderTypeLegalAdviser,
		models.StakeholderTypeFinancialAdviser,
	}
	isValid := false
	for _, validType := range validTypes {
		if *stakeholderType == validType {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("invalid stakeholder_type '%s'. Valid values are: asset_manager, asset_issuing_house, approved_asset_custodian, legal_and_professionals, rating_agency, trustees, legal_adviser, financial_adviser", *stakeholderType)
	}

	// Both are provided and type is valid - validate stakeholder exists using EXISTING db function
	_, err := usermetricsDB.GetPartnerEntityByID(*stakeholderType, int(*stakeholderID), trovoDB)
	if err != nil {
		return fmt.Errorf("stakeholder validation failed: %s with ID %d does not exist. Please create the stakeholder first using /partners/save endpoint", *stakeholderType, *stakeholderID)
	}

	return nil
}

// STEP 1
// AddAdmin handles the invitation of a new admin user.
// @Summary      Add Root Admin
// @Description  Usable by Trovo Admin Only. Invite an organization by providing their admin's email and role is by default set as: "ROOT_SUPER_ADMIN".
// @Description  Valid organization types: ASSET_MANAGER, ASSET_CUSTODIAN, RATING_AGENCY, REGULATOR, LEGAL_AGENCY, PROFESSIONAL_AGENCY
// @Description
// @Description  **Optional Stakeholder Linkage:**
// @Description  You can optionally link the organization to an existing stakeholder in the Trovo database by providing both stakeholder_id and stakeholder_type.
// @Description  Valid stakeholder_type values: asset_manager, asset_issuing_house, approved_asset_custodian, legal_and_professionals, rating_agency, trustees, legal_adviser, financial_adviser
// @Description
// @Description  Note: Both stakeholder_id and stakeholder_type must be provided together, or both omitted. The stakeholder must exist in the /partners endpoint.
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        stakeholder_type query string true "Stakeholder type" Enums(asset_manager,asset_issuing_house,approved_asset_custodian,legal_and_professionals,rating_agency,trustees,legal_adviser,financial_adviser)
// @Param        payload  body      InviteOrgAdminPayload  true  "Admin invitation payload (exclude stakeholder_type from body)"
// @Success      200      {object}  map[string]string   "Admin successfully created"
// @Failure      400      {object}  map[string]string   "Invalid payload or user does not exist"
// @Failure      500      {object}  map[string]string   "Internal server error"
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Router       /organizations/invite/root-user [post]
func TrovoAdminInvitesOrganizationRootUser(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID
		adminInfo, done, _ := server.GetUserFromContext(userID, s.TrovoWalletDB, c)
		if done {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		isAdmin, err := server.IsAdmin(adminInfo.Email, s.AdminDB)
		if err != nil || !isAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			return
		}

		var payload InviteOrgAdminPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}

		// Read stakeholder_type from query param and set it on the payload
		stakeholderType := strings.TrimSpace(strings.ToLower(c.Query("stakeholder_type")))
		if stakeholderType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "stakeholder_type is required as a query parameter"})
			return
		}
		// Validate against allowed list
		switch stakeholderType {
		case models.StakeholderTypeAssetManager, models.StakeholderTypeIssuingHouse, models.StakeholderTypeCustodian, models.StakeholderTypeLegal, models.StakeholderTypeRatingAgency, models.StakeholderTypeTrustee, models.StakeholderTypeLegalAdviser, models.StakeholderTypeFinancialAdviser:
			payload.StakeholderType = stakeholderType
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stakeholder_type. Valid: asset_manager, asset_issuing_house, approved_asset_custodian, legal_and_professionals, rating_agency, trustees, legal_adviser, financial_adviser"})
			return
		}

		// Validate Email
		if strings.TrimSpace(payload.AdminEmail) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or username must be provided"})
			return
		}

		partnerRequest := &usermetrics.SavePartnerRequest{
			PartnerType:     payload.StakeholderType,
			Action:          "create",
			UserInfo:        adminInfo,
			WalletDB:        s.TrovoWalletDB,
			UseDirectFields: true,
			PartnerName:     payload.OrganizationName,
			Address:         payload.Address,
			Country:         payload.Country,
			FeePercent:      payload.FeePercent,
			FeeFixed:        payload.FeeFixed,
		}

		response, err := inviteUserFlow(s, payload, adminInfo, partnerRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Println("Organization invite email sent successfully:", response)

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("You have successfully sent an invitation to %s to join Trovo Manager as %s",
				payload.AdminEmail,
				payload.StakeholderType),
			"invite_id":       response.InviteID,
			"organization_id": response.OrganizationID,
		})
	}
}

func inviteUserFlow(s *serverModels.Server, payload InviteOrgAdminPayload, adminInfo models.UserInfo, partnerRequest *usermetrics.SavePartnerRequest) (*models.InviteUserFlowStruct, error) {
	onboardingLevel := 1
	if payload.Level != nil {
		if *payload.Level <= 0 {
			return nil, fmt.Errorf("level must be a positive integer")
		}
		onboardingLevel = *payload.Level
	}
	//usermetrics.DBMutex.Lock()
	_, savedID, err := usermetrics.SaveAndUpdatePartnerData(partnerRequest)
	if err != nil {
		log.Printf("[PARTNER] Error saving/updating partner data: %v", err)
		return nil, fmt.Errorf("failed to save or update partner data: %v", err)
	}

	// Set the stakeholder ID from the saved partner data
	payload.StakeholderID = &savedID
	log.Printf("[PARTNER] Successfully saved partner data with ID: %d", savedID)

	// Validate that the newly created stakeholder exists and is valid
	if err := validateStakeholderLinkage(&payload.StakeholderType, payload.StakeholderID, s.TrovoWalletDB); err != nil {
		return nil, err
	}

	// Check for existing organization root admin with this email
	var invitedRootUser models.OrganizationMember
	if err := s.AdminDB.Where("email = ?", payload.AdminEmail).First(&invitedRootUser).Error; err == nil {
		msg := fmt.Sprintf("This ROOT admin already exists with role %s", invitedRootUser.Role)
		return nil, fmt.Errorf(msg)
	}

	// Generate organization ID
	orgID := uuid.New().String()
	inviteExpiresAt := time.Now().Add(models.OrganizationInviteExpiryHour * time.Hour)

	// Create organization
	org := &models.Organization{
		ID:              orgID,
		Name:            payload.OrganizationName,
		Email:           payload.AdminEmail,
		Type:            payload.StakeholderType,
		Status:          models.OrganizationStatusPending,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		CreatedBy:       adminInfo.Email,
		StakeholderID:   payload.StakeholderID,
		StakeholderType: &payload.StakeholderType,
		Level:           &onboardingLevel,
	}

	// Create invite
	invite := &models.OrganizationInvite{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Email:          payload.AdminEmail,
		Role:           string(models.OrganizationSuperAdmin),
		Status:         models.InviteStatusPending,
		ExpiresAt:      inviteExpiresAt,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		InvitedBy:      adminInfo.Email,
	}
	// Create organization member record for root user
	orgMember := &models.OrganizationMember{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Email:          payload.AdminEmail,
		Role:           string(models.OrganizationSuperAdmin),
		Status:         models.OrganizationStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		CreatedBy:      adminInfo.Email,
		FirstName:      payload.AdminFirstName,
		LastName:       payload.AdminLastName,
		InviteID:       invite.ID,
	}

	if err := s.AdminDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(org).Error; err != nil {
			return fmt.Errorf("create organization: %w", err)
		}
		if err := stakeholderServices.AutoAssignComplianceRequirements(context.Background(), tx, org.ID, payload.StakeholderType, onboardingLevel); err != nil {
			return fmt.Errorf("assign compliance requirements: %w", err)
		}
		if err := tx.Create(invite).Error; err != nil {
			return fmt.Errorf("create organization invite: %w", err)
		}
		if err := tx.Create(orgMember).Error; err != nil {
			return fmt.Errorf("create organization member: %w", err)
		}
		return nil
	}); err != nil {
		log.Printf("[ORGANIZATION] Error creating organization onboarding records: %v", err)
		return nil, fmt.Errorf("failed to create organization onboarding records: %w", err)
	}

	// Send email notification using the proper template
	// Get base URL with fallback
	frontendBaseURL := utils.GetTrovomanagerFrontendBaseUrl()
	inviteLink := fmt.Sprintf("%s/organizations/invite/validate-email/%s", frontendBaseURL, invite.ID)

	// Create admin user object for email
	adminUser := models.AdminUser{
		Email:     payload.AdminEmail,
		FirstName: payload.AdminFirstName,
		LastName:  payload.AdminLastName,
	}

	_, _, err = mail.SendOrganizationInviteEmail(adminUser, inviteLink, payload.OrganizationName, payload.StakeholderType, inviteExpiresAt)
	if err != nil {
		log.Println("Error sending organization invite email:", err)
		msg := fmt.Sprintf("Failed to send invitation email to %s: %v", payload.AdminEmail, err)
		return nil, fmt.Errorf(msg)
	}

	inviteUserFlowStruct := &models.InviteUserFlowStruct{
		OrganizationName: payload.OrganizationName,
		OrganizationType: payload.StakeholderType,
		InviteID:         invite.ID,
		OrganizationID:   orgID,
	}

	return inviteUserFlowStruct, nil
}

// inviteUserFlowStruct

//// STEP 2
//// OrganizationUserValidatesEmail Root UserValidatesEmailInvite
//// @Summary Validate root user email
//// @Description Validates a root user's email using the organization ID and OTP code
//// @Tags organizations
//// @Accept json
//// @Produce json
//// @Param invite_id path string true "Invite ID"
//// @Param request body models.RootUserValidatesEmail true "OTP verification request"
//// @Success 200 {object} map[string]interface{}
//// @Failure 400 {object} models.ErrorResponse
//// @Failure 500 {object} models.ErrorResponse
//// @Router /organizations/invite/validate-email/{invite_id} [post]
//func OrganizationUserValidatesEmail(server *serverModels.Server) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// Get organization ID from path parameter
//		invite_id := c.Param("invite_id")
//		if invite_id == "" {
//			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Organization ID is required"})
//			return
//		}
//
//		var req models.RootUserValidatesEmail
//		if err := c.ShouldBindJSON(&req); err != nil {
//			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
//			return
//		}
//
//		// Check for invited admin with valid OTP and unexpired invite
//		invite, err := services.CheckValidAdminInvite(server, invite_id, req.OTP)
//		if err != nil {
//			log.Printf("[ORGANIZATION] Error checking invite: %v", err)
//			c.JSON(http.StatusMethodNotAllowed, models.ErrorResponse{Error: "Invalid or expired or accepted invite"})
//			return
//		}
//		if invite == nil {
//			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid or expired invite"})
//			return
//		}
//		// Check if the invite is already accepted
//		if invite.Status == models.InviteStatusAccepted {
//			c.JSON(http.StatusMethodNotAllowed, models.ErrorResponse{Error: "email already validated"})
//			return
//		}
//
//		// Update status of invite and members table to VALID_EMAIL and ACTIVE
//		if err := services.UpdateInviteStatuses(server.AdminDB, invite.ID); err != nil {
//			log.Printf("[ORGANIZATION] Failed to update invite status: %v", err)
//			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update invite status"})
//			return
//		}
//
//		c.JSON(http.StatusOK, gin.H{
//			"message": "Email  validated successfully",
//			"status":  "success",
//		})
//	}
//}

// RootUserValidateInvite validates a root user's invite
// @Summary Validate root user invite
// @Description Validates a root user's invite using the invite ID and organization ID
// @Tags organizations
// @Accept json
// @Produce json
// @Param invite_id path string true "Invite ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /organizations/invite/validate/{invite_id} [post]
func RootUserValidateInvite(server *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get organization ID from path parameter
		inviteID := c.Param("invite_id")
		if inviteID == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid or expired invite"})
			return
		}

		// Check for invited admin with unexpired invite
		validateInviteResponse, err := userValidatesInvite(server, inviteID)
		if err != nil {
			log.Printf("[ORGANIZATION] Error checking invite: %v", err)
			c.JSON(http.StatusMethodNotAllowed, models.ErrorResponse{Error: "Invalid or expired or accepted invite"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":           validateInviteResponse.Message,
			"organization_name": validateInviteResponse.OrganizationName,
			"admin_email":       validateInviteResponse.AdminEmail,
			"invite_id":         validateInviteResponse.InviteID,
			"organization_id":   validateInviteResponse.OrganizationID,
			"role":              validateInviteResponse.Role,
		})
	}
}

func userValidatesInvite(server *serverModels.Server, inviteID string) (*ValidateInviteResponse, error) {
	invite, err := services.ValidateInviteIDAndStatusExpiry(server, inviteID)
	if err != nil {
		log.Printf("[ORGANIZATION] Error checking invite: %v", err)
		return nil, fmt.Errorf("invalid or already accepted or expired invite")
	}

	// Check if OTP was already sent within the last 24 hours
	otpAlreadySent := false
	if !invite.OTPExpiresAt.IsZero() && time.Now().Before(invite.OTPExpiresAt) {
		otpAlreadySent = true
	}

	// Only send new OTP if one hasn't been sent recently
	if !otpAlreadySent {
		// Update status to PENDING
		invite.Status = models.InviteStatusPending
		if err = services.UpdateInviteStatus(server.AdminDB, invite.ID, models.InviteStatusPending); err != nil {
			log.Printf("[ORGANIZATION] Failed to update invite status: %v", err)
			return nil, fmt.Errorf("failed to update invite status")
		}

		// Generate and send new OTP
		invite.VerificationOTP, _ = utils.GenerateOTP(6)
		invite.OTPExpiresAt = time.Now().Add(models.OrganizationInviteOTPExpiryHour * time.Hour)

		if err := utils.SendOTPEmail(server, invite); err != nil {
			log.Printf("[ORGANIZATION] Failed to send OTP email: %v", err)
			return nil, fmt.Errorf("failed to send OTP email")
		}

		// Store the OTP and its expiry in the invite
		if err := server.AdminDB.Save(&invite).Error; err != nil {
			log.Printf("[ORGANIZATION] Failed to save invite with OTP: %v", err)
			return nil, fmt.Errorf("failed to save invite with OTP")
		}
	}

	// Fetch organization details
	org, err := services.GetOrganization(server, invite.OrganizationID)
	if err != nil {
		log.Printf("[ORGANIZATION] Error fetching organization details: %v", err)
		//c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch organization details"})
		return nil, fmt.Errorf("failed to fetch organization details")
	}

	// Return different messages based on whether OTP was sent
	var message string
	if otpAlreadySent {
		message = "OTP already generated, please check your email"
	} else {
		message = "OTP has been generated and sent to your email"
	}

	return &ValidateInviteResponse{
		Message:          message,
		OrganizationName: org.Name,
		AdminEmail:       invite.Email,
		InviteID:         invite.ID,
		OrganizationID:   invite.OrganizationID,
		Role:             invite.Role,
	}, nil
}

//validate invite response
/*
c.JSON(http.StatusOK, gin.H{
			"message":           message,
			"organization_name": org.Name,
			"admin_email":       invite.Email,
			"invite_id":         invite.ID,
			"organization_id":   invite.OrganizationID,
			"role":              invite.Role,
		})
*/
type ValidateInviteResponse struct {
	Message          string `json:"message"`
	OrganizationName string `json:"organization_name"`
	AdminEmail       string `json:"admin_email"`
	InviteID         string `json:"invite_id"`
	OrganizationID   string `json:"organization_id"`
	Role             string `json:"role"`
}

// ListOrganizationMembers lists all members of an organization
// @Summary List organization members
// @Description Lists organization members. For organization members: lists their own organization. For Trovo admins: requires organization_id query parameter to specify which organization.
// @Tags organizations
// @Security JwtTokenAuth
// @Security OrganizationAuth
// @Produce json
// @Param organization_id query string false "Organization ID (required for Trovo admins, ignored for organization members)"
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Number of items per page (default: 10)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /organizations/members [get]
func ListOrganizationMembers(server *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Check organization role permission (org members only, admins bypass)
		if err := middleware.CheckOrganizationPermission(c, string(models.OrganizationSuperAdmin)); err != nil {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Error: err.Error()})
			return
		}

		// 2. Get organization ID (from token for org members, from query for admins)
		orgID, err := middleware.GetOrganizationIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		// 3. Verify the organization exists (optional additional check)
		var org models.Organization
		if err := server.AdminDB.First(&org, "id = ?", orgID).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Organization not found"})
			return
		}

		// Pagination
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

		members, total, err := services.ListOrganizationMembersPaginated(server, orgID, page, pageSize)
		if err != nil {
			log.Printf("[ORGANIZATION] Error listing members: %v", err)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		totalPages := (int(total) + pageSize - 1) / pageSize
		serverResponse.JSON(c, http.StatusOK, "Members retrieved successfully", gin.H{
			"members": members,
			"pagination": gin.H{
				"page":       page,
				"pageSize":   pageSize,
				"total":      total,
				"totalPages": totalPages,
			},
		}, nil)
	}
}

// GetOrganizationDetails gets the details of an organization
// @Summary Get organization details
// @Description Gets organization details. For organization members: gets their own organization details. For Trovo admins: requires organization_id query parameter to specify which organization.
// @Tags organizations
// @Security JwtTokenAuth
// @Security OrganizationAuth
// @Produce json
// @Param organization_id query string false "Organization ID (required for Trovo admins, ignored for organization members)"
// @Success 200 {object} models.Organization
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /organizations/details [get]
func GetOrganizationDetails(server *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get organization ID (from token for org members, from query for admins)
		orgID, err := middleware.GetOrganizationIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		orgDetails, err := services.GetOrganizationDetails(server, orgID)
		if err != nil {
			log.Printf("[ORGANIZATION] Error getting organization details: %v", err)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "Organization details retrieved successfully", orgDetails, nil)
	}
}

// OrganizationMemberLogin handles member login
// @Summary Member login
// @Description Authenticates an organization member and returns a JWT token
// @Tags organizations
// @Accept json
// @Produce json
// @Param request body models.OrganizationMemberLoginRequest true "Login request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /organizations/member/login [post]
func OrganizationMemberLogin(server *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.OrganizationMemberLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		token, err := services.OrganizationMemberLogin(server, req.Email, req.Password)
		if err != nil {
			log.Printf("[ORGANIZATION] Error in member login: %v", err)
			// unauthorized error only
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid email or password"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "Login successful", gin.H{"token": token}, nil)
	}
}

// @Summary Request member password reset
// @Description Initiates the password reset process for a member
// @Tags organizations
// @Accept json
// @Produce json
// @Param request body models.OrganizationMemberPasswordResetRequest true "Password reset request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /organizations/members/reset-password [post]
func RequestMemberPasswordReset(c *gin.Context) {
	var req models.OrganizationMemberPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	if err := services.RequestMemberPasswordReset(c.MustGet("server").(*serverModels.Server), &req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	serverResponse.JSON(c, http.StatusOK, "Password reset instructions sent to your email", nil, nil)
}

//// VerifyTeamMemberOTP godoc
//// @Summary Verify OTP for team member
//// @Description Verifies the OTP sent to a team member's email
//// @Tags organizations
//// @Accept json
//// @Produce json
//// @Param request body models.TeamMemberVerifyOTPRequest true "OTP verification request"
//// @Success 200 {object} models.OrganizationMember
//// @Failure 400 {object} models.ErrorResponse
//// @Failure 401 {object} models.ErrorResponse
//// @Failure 500 {object} models.ErrorResponse
//// @Router /organizations/team-member/verify-otp [post]
//func VerifyTeamMemberOTP(server *serverModels.Server) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		var req models.TeamMemberVerifyOTPRequest
//		if err := c.ShouldBindJSON(&req); err != nil {
//			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request"})
//			return
//		}
//
//		member, err := services.VerifyTeamMemberOTP(server, &req)
//		if err != nil {
//			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
//			return
//		}
//
//		c.JSON(http.StatusOK, member)
//	}
//}

//// ResendTeamMemberOTP godoc
//// @Summary Resend OTP for team member
//// @Description Resends the verification OTP to a team member's email
//// @Tags organizations
//// @Accept json
//// @Produce json
//// @Param request body models.TeamMemberResendOTPRequest true "Resend OTP request"
//// @Success 200 {object} models.SuccessResponse
//// @Failure 400 {object} models.ErrorResponse
//// @Failure 401 {object} models.ErrorResponse
//// @Failure 500 {object} models.ErrorResponse
//// @Router /organizations/team-member/resend-otp [post]
//func ResendTeamMemberOTP(server *serverModels.Server) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		var req models.TeamMemberResendOTPRequest
//		if err := c.ShouldBindJSON(&req); err != nil {
//			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request"})
//			return
//		}
//
//		if err := services.ResendTeamMemberOTP(server, &req); err != nil {
//			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
//			return
//		}
//
//		c.JSON(http.StatusOK, models.SuccessResponse{Message: "OTP resent successfully"})
//	}
//}

// ListOrganizationsPaginatedHandler lists all organizations with pagination and filters
// @Summary List organizations
// @Description Lists all organizations with pagination, filtering, and search. Usable by Trovo Admin only.
// @Tags organizations
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param name query string false "Organization name (partial match)"
// @Param email query string false "Organization email (partial match)"
// @Param type query string false "Organization type (exact match)"
// @Param status query string false "Organization status (exact match)"
// @Param createdBy query string false "Created by (partial match)"
// @Param search query string false "General search (matches name, email, type, status, createdBy)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /organizations [get]
func ListOrganizationsPaginatedHandler(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID
		adminInfo, done, _ := server.GetUserFromContext(userID, s.TrovoWalletDB, c)
		if done {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		isAdmin, err := server.IsAdmin(adminInfo.Email, s.AdminDB)
		if err != nil || !isAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			return
		}

		// Parse pagination and filter params
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
		name := c.Query("name")
		email := c.Query("email")
		orgType := c.Query("type")
		status := c.Query("status")
		createdBy := c.Query("createdBy")
		search := c.Query("search")

		orgs, total, err := services.ListOrganizationsPaginated(s, page, pageSize, name, email, orgType, status, createdBy, search)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organizations"})
			return
		}

		totalPages := (int(total) + pageSize - 1) / pageSize
		c.JSON(http.StatusOK, gin.H{
			"organizations": orgs,
			"total":         total,
			"pagination": gin.H{
				"page":       page,
				"pageSize":   pageSize,
				"total":      total,
				"totalPages": totalPages,
			},
		})
	}
}

// TrovoAdminDeactivatesOrganization deactivates an organization
// @Summary Deactivate organization
// @Description Deactivates an organization by setting its status to INACTIVE. Usable by Trovo Admin only.
// @Tags organizations
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param organization_id path string true "Organization ID"
// @Success 200 {object} map[string]interface{} "Success message with organization details"
// @Failure 400 {object} map[string]string "Bad Request (invalid data, missing fields, etc.)"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not Found (organization not found)"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /organizations/deactivate/{organization_id} [put]
func TrovoAdminDeactivatesOrganization(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		// --- Auth and User Info --- Start
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
			return
		}
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, s.TrovoWalletDB)
		if err != nil {
			log.Println("[DEACTIVATE_ORGANIZATION] error retrieving user info for UserID:", userID, "error:", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": "failed to retrieve user info"}
			}
			c.JSON(statusCode, response)
			return
		}

		// Check if user is admin
		isAdmin, err := server.IsAdmin(userInfo.Email, s.AdminDB)
		if err != nil || !isAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
			return
		}
		// --- Auth and User Info --- End

		// Get organization ID from path
		organizationID := c.Param("organization_id")
		if organizationID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "organization_id path parameter is required"})
			return
		}

		// Call service to deactivate organization
		deactivatedOrg, err := services.DeactivateOrganization(s, organizationID)
		if err != nil {
			log.Printf("[DEACTIVATE_ORGANIZATION] Error deactivating organization: %v", err)

			// Handle specific error cases
			if err.Error() == "organization not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
				return
			}
			if err.Error() == "organization is already inactive" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "organization is already inactive"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deactivate organization", "details": err.Error()})
			return
		}

		log.Printf("[DEACTIVATE_ORGANIZATION] Successfully deactivated organization %s (User: %s)", organizationID, userInfo.Username)

		c.JSON(http.StatusOK, gin.H{
			"message":           "Organization deactivated successfully",
			"organization_id":   organizationID,
			"organization_name": deactivatedOrg.Name,
			"status":            models.OrganizationStatusInactive,
		})
	}
}

// UpdateOrganizationPayload represents the request payload for updating an organization
type UpdateOrganizationPayload struct {
	OrganizationName string `json:"organization_name" binding:"required" example:"Acme Asset Management Ltd"`
	OrganizationType string `json:"-" swaggerignore:"true"`
	AdminFirstName   string `json:"admin_first_name" binding:"required" example:"John"`
	AdminLastName    string `json:"admin_last_name" binding:"required" example:"Doe"`
	// AdminEmail is intentionally excluded - email should not be updated via this endpoint
	StakeholderType string  `json:"-" swaggerignore:"true"` // Set from query param
	Address         string  `json:"address" example:"123 Main St"`
	Country         string  `json:"country" example:"USA"`
	FeePercent      string  `json:"fee_percent" example:"0.5"`
	FeeFixed        string  `json:"fee_fixed" example:"10.00"`
	StakeholderID   *uint64 `json:"-"` // Set internally after partner update, not provided by user
	// Level is the compliance requirement tier for this organisation (see
	// PRD-kyc-compliance-templates.md §6.2). Omit to leave unchanged.
	Level *int `json:"level" example:"2"`
}

// TrovoAdminUpdatesOrganizationRootUser updates an existing organization and its linked stakeholder
// @Summary Update organization and stakeholder data
// @Description Updates an existing organization and its linked stakeholder data. Requires stakeholder_type as query parameter and organization_id in path. Note: admin_email cannot be updated via this endpoint as it's a sensitive identifier.
// @Tags organizations
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param organization_id path string true "Organization ID"
// @Param stakeholder_type query string true "Stakeholder type" Enums(asset_manager,asset_issuing_house,approved_asset_custodian,legal_and_professionals,rating_agency,trustees,legal_adviser,financial_adviser)
// @Param payload body UpdateOrganizationPayload true "Organization update payload (exclude stakeholder_type and admin_email from body)"
// @Success 200 {object} map[string]interface{} "Success message with organization and stakeholder IDs"
// @Failure 400 {object} map[string]string "Bad Request (invalid data, missing fields, etc.)"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not Found (organization not found)"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /organizations/update/{organization_id} [put]
func TrovoAdminUpdatesOrganizationRootUser(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		// --- Auth and User Info --- Start
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
			return
		}
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, s.TrovoWalletDB)
		if err != nil {
			log.Println("[UPDATE_ORGANIZATION] error retrieving user info for UserID:", userID, "error:", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": "failed to retrieve user info"}
			}
			c.JSON(statusCode, response)
			return
		}
		// --- Auth and User Info --- End

		// Get organization ID from path
		organizationID := c.Param("organization_id")
		if organizationID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "organization_id path parameter is required"})
			return
		}

		// Check if organization exists
		var existingOrg models.Organization
		if err := s.AdminDB.Where("id = ?", organizationID).First(&existingOrg).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
				return
			}
			log.Printf("[UPDATE_ORGANIZATION] Error finding organization: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find organization"})
			return
		}

		// Check if organization has existing stakeholder linkage
		if existingOrg.StakeholderID == nil || existingOrg.StakeholderType == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "organization does not have a linked stakeholder to update"})
			return
		}

		// Use stakeholder_type from organization record
		stakeholderType := strings.TrimSpace(strings.ToLower(*existingOrg.StakeholderType))

		// Optionally validate query parameter matches stored type (if provided)
		queryStakeholderType := strings.TrimSpace(strings.ToLower(c.Query("stakeholder_type")))
		if queryStakeholderType != "" && queryStakeholderType != stakeholderType {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("stakeholder_type mismatch. Organization is linked to '%s', but query parameter specifies '%s'", stakeholderType, queryStakeholderType)})
			return
		}

		// Validate stakeholder_type
		switch stakeholderType {
		case models.StakeholderTypeAssetManager, models.StakeholderTypeIssuingHouse, models.StakeholderTypeCustodian, models.StakeholderTypeLegal, models.StakeholderTypeRatingAgency, models.StakeholderTypeTrustee, models.StakeholderTypeLegalAdviser, models.StakeholderTypeFinancialAdviser:
			// Valid stakeholder type
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid stakeholder_type '%s' stored in organization. Valid values are: asset_manager, asset_issuing_house, approved_asset_custodian, legal_and_professionals, rating_agency, trustees, legal_adviser, financial_adviser", stakeholderType)})
			return
		}

		// Verify stakeholder exists before attempting update
		if err := validateStakeholderLinkage(&stakeholderType, existingOrg.StakeholderID, s.TrovoWalletDB); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "stakeholder validation failed", "details": err.Error()})
			return
		}

		// Bind request payload
		var payload UpdateOrganizationPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format", "details": err.Error()})
			return
		}
		if payload.Level != nil && *payload.Level <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "level must be a positive integer"})
			return
		}

		// Set stakeholder type from organization record
		payload.StakeholderType = stakeholderType

		// Create partner request for updating stakeholder data
		partnerRequest := &usermetrics.SavePartnerRequest{
			PartnerType:     stakeholderType,
			Action:          "update",
			UserInfo:        userInfo,
			WalletDB:        s.TrovoWalletDB,
			UseDirectFields: true,
			PartnerName:     payload.OrganizationName,
			Address:         payload.Address,
			Country:         payload.Country,
			FeePercent:      payload.FeePercent,
			FeeFixed:        payload.FeeFixed,
			PartnerID:       uint64(*existingOrg.StakeholderID),
		}

		// Update stakeholder data
		_, updatedStakeholderID, err := usermetrics.SaveAndUpdatePartnerData(partnerRequest)
		if err != nil {
			log.Printf("[UPDATE_ORGANIZATION] Error updating partner data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update stakeholder data", "details": err.Error()})
			return
		}

		// Set the updated stakeholder ID
		payload.StakeholderID = &updatedStakeholderID

		// Validate that the updated stakeholder exists and is valid
		if err := validateStakeholderLinkage(&payload.StakeholderType, payload.StakeholderID, s.TrovoWalletDB); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "stakeholder validation failed", "details": err.Error()})
			return
		}

		// Update organization record (email is not updated - it's a sensitive identifier)
		updateData := map[string]interface{}{
			"name":             payload.OrganizationName,
			"type":             payload.StakeholderType,
			"stakeholder_id":   payload.StakeholderID,
			"stakeholder_type": &payload.StakeholderType,
			"updated_at":       time.Now(),
		}
		levelChanged := payload.Level != nil && (existingOrg.Level == nil || *existingOrg.Level != *payload.Level)
		if levelChanged {
			updateData["level"] = *payload.Level
		}

		if err := s.AdminDB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&existingOrg).Updates(updateData).Error; err != nil {
				return fmt.Errorf("update organization: %w", err)
			}
			if levelChanged {
				if err := stakeholderServices.AutoAssignComplianceRequirements(c.Request.Context(), tx, organizationID, stakeholderType, *payload.Level); err != nil {
					return fmt.Errorf("assign compliance requirements: %w", err)
				}
			}

			var orgMember models.OrganizationMember
			err := tx.Where("organization_id = ? AND email = ?", organizationID, existingOrg.Email).First(&orgMember).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			if err != nil {
				return fmt.Errorf("load organization member: %w", err)
			}
			return tx.Model(&orgMember).Updates(map[string]interface{}{
				"first_name": payload.AdminFirstName,
				"last_name":  payload.AdminLastName,
				"updated_at": time.Now(),
			}).Error
		}); err != nil {
			log.Printf("[UPDATE_ORGANIZATION] Error updating organization records: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update organization"})
			return
		}

		log.Printf("[UPDATE_ORGANIZATION] Successfully updated organization %s and stakeholder %d (User: %s)", organizationID, updatedStakeholderID, userInfo.Username)

		c.JSON(http.StatusOK, gin.H{
			"message":          "Organization and stakeholder data updated successfully",
			"organization_id":  organizationID,
			"stakeholder_id":   updatedStakeholderID,
			"stakeholder_type": payload.StakeholderType,
		})
	}
}

// OrganizationMemberLogout logs out an organization member.
//
// @Summary Organization member logout
// @Description Revokes every currently issued organization JWT for the authenticated member. The client must also discard its token; a new login issues a token for the new session version.
// @Tags organizations
// @Produce json
// @Security OrganizationAuth
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} object "Successfully logged out"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Router /organizations/member/logout [post]
func OrganizationMemberLogout(server *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberID := c.GetString(middleware.StakeholderContextMemberID)
		organizationID := c.GetString(middleware.StakeholderContextOrganizationID)
		if err := revokeOrganizationMemberSessions(c.Request.Context(), server.AdminDB, memberID, organizationID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Organization member not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to revoke organization session"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "successfully logged out", nil, nil)
	}
}

func revokeOrganizationMemberSessions(ctx context.Context, database *gorm.DB, memberID, organizationID string) error {
	if database == nil || strings.TrimSpace(memberID) == "" || strings.TrimSpace(organizationID) == "" {
		return gorm.ErrRecordNotFound
	}
	result := database.WithContext(ctx).Model(&models.OrganizationMember{}).
		Where("id = ? AND organization_id = ?", memberID, organizationID).
		UpdateColumn("session_version", gorm.Expr("session_version + 1"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
