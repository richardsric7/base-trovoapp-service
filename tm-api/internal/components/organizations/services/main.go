package services

import (
	"admin-panel-dashboard/internal/mail"
	"admin-panel-dashboard/internal/models"
	"admin-panel-dashboard/internal/utils"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	"gorm.io/gorm"

	serverModels "admin-panel-dashboard/internal/server/models"

	usermetricsDB "admin-panel-dashboard/internal/components/usermetrics/db"
	"admin-panel-dashboard/internal/middleware"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CreateOrganization creates a new organization and its root user
func CreateOrganization(server *serverModels.Server, req *models.OrganizationCreateRequest, adminUser interface{}) (*models.Organization, error) {
	// Start a transaction
	tx := server.AdminDB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Generate OTP
	otp, _ := utils.GenerateOTP(6)
	// otpExpiresAt := time.Now().Add(15 * time.Minute)

	// Create organization
	org := &models.Organization{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Type:      req.Type,
		Status:    models.ACTIVE,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: string(models.OrganizationSuperAdmin),
	}

	if err := tx.Create(org).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.RootUserPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// Create root user (organization admin)
	rootUser := &models.OrganizationMember{
		ID:              uuid.New().String(),
		OrganizationID:  org.ID,
		Email:           req.RootUserEmail,
		Role:            string(models.OrganizationSuperAdmin),
		Status:          models.ACTIVE,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		CreatedBy:       "SYSTEM",
		VerificationOTP: otp,
		//OTPExpiresAt:    otpExpiresAt,
		FirstName: req.RootUserFirstName,
		LastName:  req.RootUserLastName,
		Password:  string(hashedPassword),
	}

	if err := tx.Create(rootUser).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Send verification email
	emailBody := fmt.Sprintf(
		"Hi,\nWelcome to %s!\n\nYour verification code is: %s\nYour Organization ID is: %s\n\nPlease use these details to verify your organization email.\nThe Trovo Team",
		org.Name, otp, org.ID,
	)

	// Create admin user for email sending
	admin := models.AdminUser{
		Email:     org.Email,
		FirstName: org.Name,
	}

	// Send email using Mailgun template
	_, _, err = mail.SendEmailWithMailgunTemplate(admin, emailBody, "Organization Email Verification")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return org, nil
}

// VerifyOrganizationEmail verifies the organization's email using OTP
func VerifyOrganizationEmail(server *serverModels.Server, req *models.OrganizationVerifyRequest) (string, error) {
	// Find the root user
	var rootUser models.OrganizationMember
	if err := server.AdminDB.Where("organization_id = ? AND role = ?", req.OrganizationID, "ADMIN").First(&rootUser).Error; err != nil {
		return "", errors.New("organization not found")
	}

	if rootUser.Status != "PENDING" {
		return "", errors.New("organization already verified")
	}

	if time.Now().After(rootUser.OTPExpiresAt) {
		return "", errors.New("OTP expired")
	}

	if rootUser.VerificationOTP != req.OTP {
		return "", errors.New("invalid OTP")
	}

	// Generate temporary token for password setup
	token := uuid.New().String()
	expiresAt := time.Now().Add(10 * time.Minute)

	// Update root user status and save token
	rootUser.Status = "VERIFIED"
	rootUser.UpdatedAt = time.Now()
	rootUser.PasswordSetupToken = token
	rootUser.PasswordSetupExpiresAt = expiresAt
	if err := server.AdminDB.Save(&rootUser).Error; err != nil {
		return "", err
	}

	return token, nil
}

// SetOrganizationPassword sets the organization's password
func SetOrganizationPassword(server *serverModels.Server, req *models.OrganizationSetPasswordRequest) (string, error) {
	// Find the root user
	var rootUser models.OrganizationMember
	if err := server.AdminDB.Where("email = ? AND role = ?", req.Email, "ADMIN").First(&rootUser).Error; err != nil {
		return "", errors.New("organization not found")
	}

	if rootUser.Status != "VERIFIED" {
		return "", errors.New("organization not verified")
	}

	if time.Now().After(rootUser.PasswordSetupExpiresAt) {
		return "", errors.New("token expired")
	}

	if rootUser.PasswordSetupToken != req.Token {
		return "", errors.New("invalid token")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Update root user password and status
	rootUser.Password = string(hashedPassword)
	rootUser.Status = "ACTIVE"
	rootUser.PasswordSetupToken = ""
	rootUser.PasswordSetupExpiresAt = time.Time{}
	rootUser.UpdatedAt = time.Now()
	if err := server.AdminDB.Save(&rootUser).Error; err != nil {
		return "", err
	}

	// Generate JWT token
	token, err := middleware.GenerateOrganizationToken(&rootUser)
	if err != nil {
		return "", err
	}

	return token, nil
}

// InviteOrganizationMember invites a member to an organization
func InviteOrganizationMember(server *serverModels.Server, orgID string, req *models.OrganizationInviteRequest, adminInfo *models.InviteOrgMemberPayload) (*models.InviteUserFlowStruct, error) {
	// Check if member already exists
	var existingMember models.OrganizationMember
	if err := server.AdminDB.Where("email = ? AND organization_id = ?", req.Email, orgID).First(&existingMember).Error; err == nil {
		return nil, fmt.Errorf("member already exists in this organization")
	}

	inviteExpiresAt := time.Now().Add(models.OrganizationInviteExpiryHour * time.Hour)

	// Create invite
	invite := &models.OrganizationInvite{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Email:          req.Email,
		Role:           string(models.OrganizationMemberRole),
		Status:         models.InviteStatusPending,
		ExpiresAt:      inviteExpiresAt,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		InvitedBy:      adminInfo.AdminEmail,
	}
	//tx.Create(&invite)
	if err := server.AdminDB.Create(&invite).Error; err != nil {
		log.Printf("[ORGANIZATION] Error creating invite: %v", err)
		return nil, fmt.Errorf("failed to create invite: %v", err)
	}

	// Create organization member record for root user
	orgMember := &models.OrganizationMember{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Email:          req.Email,
		Role:           string(models.OrganizationMemberRole),
		Status:         models.OrganizationStatusPending,
		CreatedAt:      time.Now(), // correct timezone handling
		UpdatedAt:      time.Now(),
		CreatedBy:      adminInfo.AdminEmail,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		InviteID:       invite.ID,
	}

	if err := server.AdminDB.Create(&orgMember).Error; err != nil {
		log.Printf("[ORGANIZATION] Error creating organization member: %v", err)
		return nil, fmt.Errorf("failed to create organization member: %v", err)
	}

	// Send email notification using the proper template
	// Get base URL with fallback
	baseURL := utils.GetTrovomanagerFrontendBaseUrl()
	inviteLink := fmt.Sprintf("%s/organizations/invite/validate-email/%s", baseURL, invite.ID)

	// Create admin user object for email
	adminUser := models.AdminUser{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	_, _, err := mail.SendOrganizationInviteEmail(adminUser, inviteLink, adminInfo.OrganizationName, adminInfo.OrganizationType, inviteExpiresAt)
	if err != nil {
		log.Println("Error sending organization invite email:", err)
		return nil, fmt.Errorf("failed to send invite email: %v", err)
	}

	inviteUserFlowStruct := &models.InviteUserFlowStruct{
		OrganizationName: adminInfo.OrganizationName,
		OrganizationType: adminInfo.OrganizationType,
		InviteID:         invite.ID,
		OrganizationID:   orgID,
	}

	//inviteUserFlowStruct := &InviteUserFlowStruct{
	//	OrganizationName: payload.OrganizationName,
	//	OrganizationType: payload.OrganizationType,
	//	InviteID:         invite.ID,
	//	OrganizationID:   orgID,
	//}

	//return inviteUserFlowStruct, nil

	//// Generate OTP and expiration time
	//otp, _ := utils.GenerateOTP(6)
	//otpExpiresAt := time.Now().Add(24 * time.Hour)
	//
	//// Create new member record
	//member := &models.OrganizationMember{
	//	ID:              uuid.New().String(),
	//	OrganizationID:  orgID,
	//	Email:           req.Email,
	//	Role:            string(models.OrganizationMemberRole),
	//	Status:          models.OrganizationStatusPending,
	//	CreatedAt:       time.Now(),
	//	UpdatedAt:       time.Now(),
	//	CreatedBy:       "SYSTEM", // This will be updated with the actual inviter's email
	//	VerificationOTP: otp,
	//	OTPExpiresAt:    otpExpiresAt,
	//	FirstName:       req.FirstName,
	//	LastName:        req.LastName,
	//}
	//
	//if err := server.AdminDB.Create(&member).Error; err != nil {
	//	return fmt.Errorf("failed to create member: %v", err)
	//}
	//
	//// Send invitation email with OTP
	//go func() {
	//	emailBody := fmt.Sprintf(
	//		"You have been invited to join an organization. Please use the following OTP to verify your email:\n\nOTP: %s\n\nThis OTP will expire in 24 hours.",
	//		otp,
	//	)
	//	_, response, err := mail.SendEmailWithMailgunTemplate(models.AdminUser{Email: req.Email}, emailBody, "Organization Invitation")
	//	if err != nil {
	//		log.Println("Error sending email:", err)
	//		return
	//	}
	//	log.Println("Email sent successfully:", response)
	//}()

	return inviteUserFlowStruct, nil
}

// ListOrganizationMembers lists all members of an organization
func ListOrganizationMembers(server *serverModels.Server, orgID string) ([]models.OrganizationMemberResponse, error) {
	var members []models.OrganizationMember

	if err := server.AdminDB.Where("organization_id = ?", orgID).
		Order("created_at DESC").
		Find(&members).Error; err != nil {
		return nil, err
	}

	// Convert to response format
	response := make([]models.OrganizationMemberResponse, len(members))
	for i, member := range members {
		response[i] = models.OrganizationMemberResponse{
			ID:        member.ID,
			Email:     member.Email,
			FirstName: member.FirstName,
			LastName:  member.LastName,
			Role:      member.Role,
			Status:    member.Status,
			CreatedAt: member.CreatedAt,
		}
	}

	return response, nil
}

// ListOrganizationMembersPaginated lists all members of an organization with pagination
func ListOrganizationMembersPaginated(server *serverModels.Server, orgID string, page, pageSize int) ([]models.OrganizationMemberResponse, int64, error) {
	var members []models.OrganizationMember
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// Count total members
	if err := server.AdminDB.Model(&models.OrganizationMember{}).Where("organization_id = ?", orgID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated members
	if err := server.AdminDB.Where("organization_id = ?", orgID).
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&members).Error; err != nil {
		return nil, 0, err
	}

	// Convert to response format
	response := make([]models.OrganizationMemberResponse, len(members))
	for i, member := range members {
		response[i] = models.OrganizationMemberResponse{
			ID:        member.ID,
			Email:     member.Email,
			FirstName: member.FirstName,
			LastName:  member.LastName,
			Role:      member.Role,
			Status:    member.Status,
			CreatedAt: member.CreatedAt,
		}
	}

	return response, total, nil
}

// getOrganizationByID is a helper method to fetch organization by ID
func getOrganizationByID(server *serverModels.Server, orgID string) (*models.Organization, error) {
	var org models.Organization
	if err := server.AdminDB.Where("id = ?", orgID).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organization not found")
		}
		return nil, fmt.Errorf("failed to find organization: %w", err)
	}
	return &org, nil
}

// GetOrganizationAdmin fetches the organization admin (root user) and sanitizes sensitive data
func GetOrganizationAdmin(server *serverModels.Server, orgID string) (*models.OrganizationMember, error) {
	var orgAdmin models.OrganizationMember
	if err := server.AdminDB.Where("organization_id = ? AND role = ?", orgID, models.OrganizationSuperAdmin).First(&orgAdmin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("organization admin not found")
		}
		return nil, err
	}

	// Clear sensitive data
	orgAdmin.Password = ""
	orgAdmin.VerificationOTP = ""
	orgAdmin.PasswordSetupToken = ""
	orgAdmin.PasswordResetOTP = ""

	return &orgAdmin, nil
}

// OrganizationAdminResponse represents the limited details of an organization admin
type OrganizationAdminResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    string `json:"status"`
}

// OrganizationDetailsResponse represents the response for organization details including stakeholder info
type OrganizationDetailsResponse struct {
	Organization      *models.Organization       `json:"organization"`
	Stakeholder       interface{}                `json:"stakeholder,omitempty"`
	OrganizationAdmin *OrganizationAdminResponse `json:"organization_admin,omitempty"`
}

// GetOrganizationDetails gets the details of an organization and its linked stakeholder if available
func GetOrganizationDetails(server *serverModels.Server, orgID string) (*OrganizationDetailsResponse, error) {
	// Get organization using helper method
	org, err := getOrganizationByID(server, orgID)
	if err != nil {
		return nil, err
	}

	response := &OrganizationDetailsResponse{
		Organization: org,
	}

	// Fetch organization admin details
	if orgAdmin, err := GetOrganizationAdmin(server, orgID); err == nil {
		response.OrganizationAdmin = &OrganizationAdminResponse{
			FirstName: orgAdmin.FirstName,
			LastName:  orgAdmin.LastName,
			Email:     orgAdmin.Email,
			Role:      orgAdmin.Role,
			Status:    orgAdmin.Status,
		}
	} else {
		log.Printf("[GET_ORGANIZATION_DETAILS] Error fetching organization admin (org_id: %s): %v", orgID, err)
	}

	// If organization has stakeholder linkage, fetch stakeholder details
	if org.StakeholderID != nil && org.StakeholderType != nil {
		stakeholderType := *org.StakeholderType
		stakeholderID := int(*org.StakeholderID)

		stakeholder, err := usermetricsDB.GetPartnerEntityByID(stakeholderType, stakeholderID, server.TrovoWalletDB)
		if err != nil {
			// Log the error but don't fail the entire request - just return org without stakeholder
			log.Printf("[GET_ORGANIZATION_DETAILS] Error fetching stakeholder (type: %s, id: %d): %v", stakeholderType, stakeholderID, err)
		} else {
			response.Stakeholder = stakeholder
		}
	}

	return response, nil
}

// SendPasswordSetupToken resends the password setup token
func SendPasswordSetupToken(server *serverModels.Server, email string) (string, error) {
	// Find the root user
	var rootUser models.OrganizationMember
	if err := server.AdminDB.Where("email = ? AND role = ?", email, "ADMIN").First(&rootUser).Error; err != nil {
		return "", errors.New("organization not found")
	}

	if rootUser.Status != "VERIFIED" {
		return "", errors.New("organization not verified")
	}

	// Generate new token
	token := uuid.New().String()
	expiresAt := time.Now().Add(10 * time.Minute)

	// Update root user with new token
	rootUser.PasswordSetupToken = token
	rootUser.PasswordSetupExpiresAt = expiresAt
	rootUser.UpdatedAt = time.Now()
	if err := server.AdminDB.Save(&rootUser).Error; err != nil {
		return "", err
	}

	return token, nil
}

// SendPasswordToken resends the password setup token
func SendPasswordToken(server *serverModels.Server, email string) (string, error) {
	// Find the root user
	var rootUser models.OrganizationInvite
	if err := server.AdminDB.
		Where("email = ? AND role IN ?", email, []string{string(models.OrganizationSuperAdmin), string(models.OrganizationMemberRole)}).
		First(&rootUser).Error; err != nil {
		return "", errors.New("user not found")
	}
	if rootUser.Status != "VERIFIED" {
		return "", errors.New("user not invited")
	}

	// Generate new token
	token, _ := utils.GenerateOTP(6)
	expiresAt := time.Now().Add(20 * time.Minute)

	// Update root user with new token
	rootUser.VerificationOTP = token
	rootUser.ExpiresAt = expiresAt
	rootUser.UpdatedAt = time.Now()
	if err := server.AdminDB.Save(&rootUser).Error; err != nil {
		return "", err
	}

	return token, nil
}

// OrganizationMemberLogin handles member login
func OrganizationMemberLogin(server *serverModels.Server, email, password string) (string, error) {
	var member models.OrganizationMember
	if err := server.AdminDB.Where("email = ?", email).First(&member).Error; err != nil {
		return "", errors.New("member not found")
	}

	if member.Status != models.OrganizationStatusActive {
		return "", errors.New("member account is not ACTIVE")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	// Generate JWT token
	token, err := middleware.GenerateOrganizationToken(&member)
	if err != nil {
		return "", err
	}

	return token, nil
}

// RequestMemberPasswordReset initiates the password reset process for a member
func RequestMemberPasswordReset(server *serverModels.Server, req *models.OrganizationMemberPasswordResetRequest) error {
	// Find the member
	var member models.OrganizationMember
	if err := server.AdminDB.Where("organization_id = ? AND email = ?", req.OrganizationID, req.Email).First(&member).Error; err != nil {
		return errors.New("member not found")
	}

	// Generate OTP
	otp, _ := utils.GenerateOTP(6)
	otpExpiresAt := time.Now().Add(15 * time.Minute)

	// Update member with OTP
	member.PasswordResetOTP = otp
	member.PasswordResetExpiresAt = otpExpiresAt
	member.UpdatedAt = time.Now()

	if err := server.AdminDB.Save(&member).Error; err != nil {
		return err
	}

	// Send password reset email
	emailBody := fmt.Sprintf(
		"Hi,\nYou have requested to reset your password.\n\nYour verification code is: %s\nYour Organization ID is: %s\n\nPlease use these details to reset your password.\nThe Trovo Team",
		otp, req.OrganizationID,
	)

	// Create admin user for email sending
	admin := models.AdminUser{
		Email:     member.Email,
		FirstName: member.Email, // Using email as first name since we don't have member's name
	}

	// Send email using Mailgun template
	_, _, err := mail.SendEmailWithMailgunTemplate(admin, emailBody, "Password Reset Request")
	if err != nil {
		return err
	}

	return nil
}

// SetMemberNewPassword sets a new password for a member after OTP verification
func SetMemberNewPassword(server *serverModels.Server, req *models.OrganizationMemberSetPasswordRequest) error {
	// Find the member
	var member models.OrganizationMember
	if err := server.AdminDB.Where("organization_id = ? AND email = ?", req.OrganizationID, req.Email).First(&member).Error; err != nil {
		return errors.New("member not found")
	}

	// Verify OTP
	if member.PasswordResetOTP != req.OTP {
		return errors.New("invalid OTP")
	}

	if time.Now().After(member.PasswordResetExpiresAt) {
		return errors.New("OTP expired")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update member password and clear OTP
	member.Password = string(hashedPassword)
	member.PasswordResetOTP = ""
	member.PasswordResetExpiresAt = time.Time{}
	member.UpdatedAt = time.Now()

	if err := server.AdminDB.Save(&member).Error; err != nil {
		return err
	}

	return nil
}

// CheckValidAdminInvite checks if an admin invite is valid
func CheckValidAdminInvite(server *serverModels.Server, invite_id string, otp string) (*models.OrganizationInvite, error) {
	slog.Info("Checking OTP", "invite_id", invite_id, "otp", otp)
	var invite models.OrganizationInvite
	if err := server.AdminDB.Where("id = ? AND verification_otp = ? AND status = ?", invite_id, otp, models.InviteStatusPending).First(&invite).Error; err != nil {
		slog.Error("Error checking invite", "error", err)
		return nil, fmt.Errorf("invalid or expired invite")
	}

	// Check if OTP has expired
	if time.Now().After(invite.OTPExpiresAt) {
		return nil, fmt.Errorf("OTP has expired")
	}

	return &invite, nil
}

// ValidateInviteIDExpiryAndOtp validates invite ID and checks if it's expired // status must be EMAIL_VALID ALSO
func ValidateInviteIDExpiryAndOtp(server *serverModels.Server, inviteID string, req models.SetupUserPasswordRequest) (*models.OrganizationInvite, error) {
	var invite models.OrganizationInvite
	if err := server.AdminDB.Where("id = ? AND status = ? AND expires_at > ?",
		inviteID,
		models.InviteStatusPending,
		time.Now(),
	).First(&invite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invite not found or expired or already accepted")
		}
		return nil, fmt.Errorf("error fetching invite: %v", err)
	}

	// check if OTP is valid
	if invite.VerificationOTP != req.OTP {
		return nil, fmt.Errorf("invalid or expired invite")
	}

	if time.Now().After(invite.OTPExpiresAt) {
		return nil, fmt.Errorf("OTP has expired, please request a new one")
	}

	return &invite, nil
}

// IsEmailStatusValid checks if the email status is accepted for an invite
func IsEmailStatusValid(server *serverModels.Server, inviteID string) (bool, error) {
	var invite models.OrganizationMember
	if err := server.AdminDB.Where("invite_id = ?", inviteID).First(&invite).Error; err != nil {
		return false, err
	}
	return invite.Status == models.EndInviteStatusValidEmail, nil
}

// is status active in  models.OrganizationMember?
func IsStatusActive(server *serverModels.Server, email string) (bool, error) {
	var member models.OrganizationMember
	if err := server.AdminDB.Where("email = ?", email).First(&member).Error; err != nil {
		return false, err
	}
	return member.Status == models.OrganizationStatusActive, nil
}

// UpdateInviteStatuses updates both invite and member statuses
func UpdateInviteStatuses(db *gorm.DB, inviteID string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Update invite status
		if err := tx.Model(&models.OrganizationInvite{}).
			Where("id = ?", inviteID).
			Updates(map[string]interface{}{
				"email_validation_status": models.EndInviteStatusValidEmail,
				"email_validated_at":      time.Now(),
				"status":                  models.InviteStatusPending,
			}).Error; err != nil {
			return err
		}

		// Update member status
		if err := tx.Model(&models.OrganizationMember{}).
			Where("invite_id = ?", inviteID).
			Update("status", models.EndInviteStatusValidEmail).Error; err != nil {
			return err
		}

		return nil
	})
}

// ValidateInviteIDAndStatusExpiry validates invite ID, status, and expiry
func ValidateInviteIDAndStatusExpiry(server *serverModels.Server, inviteID string) (*models.OrganizationInvite, error) {
	var invite models.OrganizationInvite
	if err := server.AdminDB.Where("id = ? AND status = ? AND expires_at > ?",
		inviteID,
		models.InviteStatusPending,
		time.Now(),
	).First(&invite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invite not found or expired or already accepted")
		}
		return nil, fmt.Errorf("error fetching invite: %v", err)
	}
	return &invite, nil
}

// method to validate invite id and expiry of the invite id
func ValidateInviteIDAndExpiry(server *serverModels.Server, inviteID string) (*models.OrganizationInvite, error) {
	var invite models.OrganizationInvite
	if err := server.AdminDB.Where("id = ? AND status = ? AND expires_at > ?",
		inviteID,
		models.InviteStatusPending,
		time.Now(),
	).First(&invite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Invite not found or expired", err)
			return nil, fmt.Errorf("invite not found or expired or already accepted")
		}
		log.Println("Invite ERROR", err)
		return nil, fmt.Errorf("error fetching invite: %v", err)
	}
	return &invite, nil
}

// CheckValidAdminInvite checks if there's a non-expired ADMIN invite matching the email and token
func CheckValidAdminInvite1(server *serverModels.Server, payload *models.RootUserValidatesEmail) (*models.OrganizationInvite, error) {
	var invite models.OrganizationInvite

	err := server.AdminDB.
		Where("token = ? AND role IN ? AND status = ? AND expires_at > ?",
			// payload.OTP,
			[]string{string(models.OrganizationSuperAdmin), string(models.OrganizationMemberRole)},
			"PENDING",
			time.Now(),
		).
		First(&invite).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found — return nil instead of error
		}
		return nil, err // DB or other error
	}

	return &invite, nil
}

func UpdateInviteStatus(db *gorm.DB, inviteID string, newStatus string) error {
	return db.Model(&models.OrganizationInvite{}).
		Where("id = ?", inviteID).
		Update("status", newStatus).Error
}

func UpdateEmailInviteStatus_(db *gorm.DB, inviteID string, newStatus string) error {
	return db.Model(&models.OrganizationInvite{}).
		Where("id = ?", inviteID).
		Updates(map[string]interface{}{
			"email_validation_status": newStatus,
			"email_validated_at":      time.Now(),
		}).Error
}

func UpdateEmailInviteStatus(db *gorm.DB, inviteID string, newStatus string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Step 1: Update the invite status and timestamp
		if err := tx.Model(&models.OrganizationInvite{}).
			Where("id = ?", inviteID).
			Updates(map[string]interface{}{
				"email_validation_status": newStatus,
				"email_validated_at":      time.Now(),
				"status":                  models.InviteStatusAccepted,
			}).Error; err != nil {
			return fmt.Errorf("failed to update invite status: %w", err)
		}

		// Step 2: Update the corresponding member's status to ACTIVE
		result := tx.Model(&models.OrganizationMember{}).
			Where("invite_id = ?", inviteID).
			Update("status", newStatus)

		if result.Error != nil {
			return fmt.Errorf("failed to update member status: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("no member found for invite ID: %s", inviteID)
		}

		return nil // commit transaction
	})
}

// GetOrganization by ID
func GetOrganization(server *serverModels.Server, orgID string) (*models.Organization, error) {
	var org models.Organization
	if err := server.AdminDB.Where("id = ?", orgID).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

// GET USER by ID
func GetUserByID(server *serverModels.Server, userID string) (*models.OrganizationMember, error) {
	var user models.OrganizationMember
	if err := server.AdminDB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// get user by email
func GetUserByEmail(server *serverModels.Server, email string) (*models.OrganizationMember, error) {
	var user models.OrganizationMember
	if err := server.AdminDB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// get user by invite id
func GetUserByInviteID(server *serverModels.Server, inviteID string) (*models.OrganizationMember, error) {
	var invite models.OrganizationMember
	if err := server.AdminDB.Where("invite_id = ?", inviteID).First(&invite).Error; err != nil {
		return nil, err
	}
	return &invite, nil
}

func IsInviteTokenExpired(db *gorm.DB, token string) (bool, error) {
	var invite models.OrganizationInvite

	err := db.Select("expires_at").
		Where("token = ?", token).
		First(&invite).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return time.Now().After(invite.ExpiresAt), nil
}

// SetupRootUserProfile sets up the root user's profile after accepting the invitation
func SetupRootUserProfile(server *serverModels.Server, req *models.RootUserSetupProfileRequest) (string, error) {
	// Find the invite
	var invite models.OrganizationInvite
	if err := server.AdminDB.Where("email = ? AND token = ? AND role = ? AND status = ? AND expires_at > ?",
		req.Email, req.Token, models.OrganizationSuperAdmin, "ACCEPTED", time.Now()).
		First(&invite).Error; err != nil {
		return "", errors.New("invalid or expired invite")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Start a transaction
	tx := server.AdminDB.Begin()
	if tx.Error != nil {
		return "", tx.Error
	}

	// Create organization first
	org := &models.Organization{
		ID:        uuid.New().String(),
		Name:      req.OrganizationName,
		Email:     req.OrganizationEmail, // Using organization's contact email
		Type:      req.OrganizationType,
		Status:    "ACTIVE",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: "SYSTEM",
	}

	if err := tx.Create(org).Error; err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to create organization: %v", err)
	}

	// Create root user with the same organization ID
	rootUser := &models.OrganizationMember{
		ID:             uuid.New().String(),
		OrganizationID: org.ID,    // Using the same ID as the organization
		Email:          req.Email, // Using member's personal email
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Password:       string(hashedPassword),
		Role:           string(models.OrganizationSuperAdmin),
		Status:         "ACTIVE",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		CreatedBy:      "SYSTEM",
	}

	if err := tx.Create(rootUser).Error; err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to create root user: %v", err)
	}

	// Update invite status to COMPLETED
	if err := tx.Model(&models.OrganizationInvite{}).
		Where("id = ?", invite.ID).
		Update("status", "COMPLETED").Error; err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to update invite status: %v", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return "", fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Generate JWT token
	token, err := middleware.GenerateOrganizationToken(rootUser)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}

// ResendRootUserOTP resends the OTP for root user verification
func ResendRootUserOTP(server *serverModels.Server, req *models.ResendOTPRequest) error {
	// Find the invite
	var invite models.OrganizationInvite
	if err := server.AdminDB.Where("email = ? AND token = ? AND role = ? AND status = ? AND expires_at > ?",
		req.Email, req.Token, models.OrganizationSuperAdmin, "ACCEPTED", time.Now()).
		First(&invite).Error; err != nil {
		return errors.New("invalid or expired invite")
	}

	// Generate new OTP
	otp, _ := utils.GenerateOTP(6)
	otpExpiresAt := time.Now().Add(15 * time.Minute)

	// Update invite with new OTP
	invite.VerificationOTP = otp
	invite.OTPExpiresAt = otpExpiresAt
	invite.UpdatedAt = time.Now()

	if err := server.AdminDB.Save(&invite).Error; err != nil {
		return err
	}

	// Send OTP email
	emailBody := fmt.Sprintf(
		"Hi,\nYou have requested a new verification code.\n\nYour new verification code is: %s\n\nPlease use this code to verify your email.\nThe Trovo Team",
		otp,
	)

	// Create admin user for email sending
	admin := models.AdminUser{
		Email:     req.Email,
		FirstName: req.Email, // Using email as first name since we don't have user's name yet
	}

	// Send email using Mailgun template
	_, _, err := mail.SendEmailWithMailgunTemplate(admin, emailBody, "New Verification Code")
	if err != nil {
		return err
	}

	return nil
}

// VerifyTeamMemberOTP verifies the OTP for a team member
func VerifyTeamMemberOTP(server *serverModels.Server, req *models.TeamMemberVerifyOTPRequest) (*models.OrganizationMember, error) {
	// Find the invite
	var invite models.OrganizationInvite
	if err := server.AdminDB.Where("email = ? AND token = ? AND role = ? AND status = ? AND expires_at > ?",
		req.Email, req.Token, "MEMBER", "PENDING", time.Now()).First(&invite).Error; err != nil {
		return nil, fmt.Errorf("invalid or expired invite")
	}

	// Verify OTP
	if invite.VerificationOTP != req.OTP {
		return nil, fmt.Errorf("invalid OTP")
	}

	if time.Now().After(invite.OTPExpiresAt) {
		return nil, fmt.Errorf("OTP has expired")
	}

	// Create or update member
	var member models.OrganizationMember
	result := server.AdminDB.Where("email = ? AND organization_id = ?", req.Email, invite.OrganizationID).First(&member)
	if result.Error != nil {
		// Create new member
		member = models.OrganizationMember{
			ID:             uuid.New().String(),
			OrganizationID: invite.OrganizationID,
			Email:          req.Email,
			Role:           invite.Role,
			Status:         "PENDING", // Will be updated to ACTIVE after profile setup
			CreatedBy:      invite.InvitedBy,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := server.AdminDB.Create(&member).Error; err != nil {
			return nil, fmt.Errorf("failed to create member: %v", err)
		}
	}

	// Update invite status
	invite.Status = "OTP_VERIFIED"
	if err := server.AdminDB.Save(&invite).Error; err != nil {
		return nil, fmt.Errorf("failed to update invite status: %v", err)
	}

	return &member, nil
}

// SetupTeamMemberProfile sets up a team member's profile after OTP verification
func SetupTeamMemberProfile(server *serverModels.Server, req *models.TeamMemberSetupProfileRequest) (string, error) {
	// Find the invite
	var invite models.OrganizationInvite
	if err := server.AdminDB.Where("email = ? AND token = ? AND role = ? AND status = ? AND expires_at > ?",
		req.Email, req.Token, "MEMBER", "OTP_VERIFIED", time.Now()).First(&invite).Error; err != nil {
		return "", fmt.Errorf("invalid or expired invite")
	}

	// Find the member
	var member models.OrganizationMember
	if err := server.AdminDB.Where("email = ? AND organization_id = ?", req.Email, invite.OrganizationID).First(&member).Error; err != nil {
		return "", fmt.Errorf("member not found")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	// Update member
	member.Password = string(hashedPassword)
	member.FirstName = req.FirstName
	member.LastName = req.LastName
	member.Status = "ACTIVE"
	member.UpdatedAt = time.Now()

	if err := server.AdminDB.Save(&member).Error; err != nil {
		return "", fmt.Errorf("failed to update member: %v", err)
	}

	// Update invite status
	invite.Status = "COMPLETED"
	if err := server.AdminDB.Save(&invite).Error; err != nil {
		return "", fmt.Errorf("failed to update invite status: %v", err)
	}

	// Generate JWT token
	token, err := middleware.GenerateOrganizationToken(&member)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}

// ResendTeamMemberOTP resends the OTP for team member verification
func ResendTeamMemberOTP(server *serverModels.Server, req *models.TeamMemberResendOTPRequest) error {
	// Find the invite
	var invite models.OrganizationInvite
	if err := server.AdminDB.Where("email = ? AND token = ? AND role = ? AND status = ? AND expires_at > ?",
		req.Email, req.Token, "MEMBER", "PENDING", time.Now()).First(&invite).Error; err != nil {
		return fmt.Errorf("invalid or expired invite")
	}

	// Generate new OTP
	otp, _ := utils.GenerateOTP(6)
	otpExpiresAt := time.Now().Add(15 * time.Minute)

	// Update invite with new OTP
	invite.VerificationOTP = otp
	invite.OTPExpiresAt = otpExpiresAt
	invite.Email = req.Email
	if err := server.AdminDB.Save(&invite).Error; err != nil {
		return fmt.Errorf("failed to update invite: %v", err)
	}

	// Send OTP via email
	if err := utils.SendOTPEmail(server, &invite); err != nil {
		return fmt.Errorf("failed to send OTP email: %v", err)
	}

	return nil
}

// SendEmailVerificationOTP sends an OTP to verify the user's email
func SendEmailVerificationOTP(server *serverModels.Server, req *models.SendEmailVerificationOTPRequest) error {
	// Find the invite
	var invite models.OrganizationInvite
	if err := server.AdminDB.Where("email = ? AND token = ? AND status = ? AND expires_at > ?",
		req.Email, req.Token, "PENDING", time.Now()).First(&invite).Error; err != nil {
		return fmt.Errorf("invalid or expired invite")
	}

	// Generate new OTP
	otp, _ := utils.GenerateOTP(6)
	otpExpiresAt := time.Now().Add(15 * time.Minute)

	// Update invite with new OTP
	invite.VerificationOTP = otp
	invite.OTPExpiresAt = otpExpiresAt
	invite.Email = req.Email
	if err := server.AdminDB.Save(&invite).Error; err != nil {
		return fmt.Errorf("failed to update invite: %v", err)
	}

	// Send OTP via email
	if err := utils.SendOTPEmail(server, &invite); err != nil {
		return fmt.Errorf("failed to send OTP email: %v", err)
	}

	return nil
}

// VerifyEmailOTP verifies the OTP sent to the user's email
func VerifyEmailOTP(server *serverModels.Server, req *models.VerifyEmailOTPRequest) error {
	// Find the invite
	var invite models.OrganizationInvite
	if err := server.AdminDB.Where("email = ? AND status = ? AND expires_at > ?",
		req.Email, "VERIFIED", time.Now()).First(&invite).Error; err != nil {
		return fmt.Errorf("invalid or expired invite")
	}

	// Verify OTP
	if invite.VerificationOTP != req.OTP {
		return fmt.Errorf("invalid OTP")
	}

	if time.Now().After(invite.OTPExpiresAt) {
		return fmt.Errorf("OTP has expired")
	}

	// Update invite status
	invite.Status = "EMAIL_VERIFIED"
	if err := server.AdminDB.Save(&invite).Error; err != nil {
		return fmt.Errorf("failed to update invite status: %v", err)
	}

	return nil
}

// ValidateMemberInvite validates a member's invite token
func ValidateMemberInvite(server *serverModels.Server, organizationID string, userID string, otp string) (*models.OrganizationMember, error) {
	// Find the member with the verification OTP and matching organization ID
	var member models.OrganizationMember
	if err := server.AdminDB.Where("verification_otp = ? AND status = ? AND organization_id = ? AND id = ?",
		otp,
		models.OrganizationStatusPending,
		organizationID,
		userID,
	).First(&member).Error; err != nil {
		return nil, fmt.Errorf("invalid or expired invite")
	}

	// Check if OTP has expired
	if time.Now().After(member.OTPExpiresAt) {
		return nil, fmt.Errorf("invite has expired. Please request a new invitation")
	}

	// Return member details needed for password setup
	return &member, nil
}

// ListOrganizationsPaginated returns a paginated list of organizations with filters
func ListOrganizationsPaginated(server *serverModels.Server, page, pageSize int, name, email, orgType, status, createdBy, search string) ([]models.OrganizationResponse, int64, error) {
	var organizations []models.Organization
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	query := server.AdminDB.Model(&models.Organization{})

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if email != "" {
		query = query.Where("email ILIKE ?", "%"+email+"%")
	}
	if orgType != "" {
		query = query.Where("type = ?", orgType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if createdBy != "" {
		query = query.Where("created_by ILIKE ?", "%"+createdBy+"%")
	}
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR type ILIKE ? OR status ILIKE ? OR created_by ILIKE ?", like, like, like, like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&organizations).Error; err != nil {
		return nil, 0, err
	}
	mapOrganizations := make([]models.OrganizationResponse, len(organizations))
	for i, org := range organizations {
		orgAdmin, err := GetOrganizationRootAdmin(server, org.ID)
		if err != nil {
			log.Printf("[LIST_ORGANIZATIONS] Could not find root admin for org %s: %v", org.ID, err)
		}
		mapOrganizations[i] = *mapOrganizationWithAdmin(&org, orgAdmin)
		memberCount, err := GetOrganizationMemberCount(server, org.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get member count for organization %s: %v", org.ID, err)
		}
		mapOrganizations[i].OrganizationCount = memberCount

		// Fetch and populate partner details if available
		if org.StakeholderID != nil && org.StakeholderType != nil {
			stakeholderType := *org.StakeholderType
			stakeholderID := int(*org.StakeholderID)

			stakeholder, err := usermetricsDB.GetPartnerEntityByID(stakeholderType, stakeholderID, server.TrovoWalletDB)
			if err != nil {
				log.Printf("[LIST_ORGANIZATIONS] Error fetching stakeholder (type: %s, id: %d): %v", stakeholderType, stakeholderID, err)
			} else {
				address, country, feeFixed, feePercent := GetPartnerDetails(stakeholder)
				mapOrganizations[i].Address = address
				mapOrganizations[i].Country = country
				mapOrganizations[i].FeeFixed = feeFixed
				mapOrganizations[i].FeePercent = feePercent
			}
		}
	}

	return mapOrganizations, total, nil
}

func GetAdminUserByEmail(email string, db *gorm.DB) (userInfo models.UserInfo, err error) {
	if err := db.Where("email = ?", email).First(&userInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userInfo, fmt.Errorf("admin user not found")
		}
		return userInfo, fmt.Errorf("error fetching admin user: %v", err)
	}
	return userInfo, nil
}

// get member count for an organization
func GetOrganizationMemberCount(server *serverModels.Server, organizationID string) (int64, error) {
	var count int64
	if err := server.AdminDB.Model(&models.OrganizationMember{}).Where("organization_id = ?", organizationID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetPartnerDetails extracts address, country, and fees from a partner entity
func GetPartnerDetails(entity interface{}) (string, string, string, string) {
	var address, country, feeFixed, feePercent string

	switch v := entity.(type) {
	case models.AssetManager:
		address = v.AssetManagerAddress
		country = v.AssetManagerCountry
		feeFixed = fmt.Sprintf("%f", v.FeeFixed)
		feePercent = fmt.Sprintf("%f", v.FeePercent)
	case models.AssetIssuingHouse:
		address = v.AssetIssuingHouseAddress
		country = v.AssetIssuingHouseCountry
		feeFixed = fmt.Sprintf("%f", v.FeeFixed)
		feePercent = fmt.Sprintf("%f", v.FeePercent)
	case models.ApprovedAssetCustodian:
		address = v.AssetCustodianAddress
		country = v.AssetCustodianCountry
		feeFixed = fmt.Sprintf("%f", v.FeeFixed)
		feePercent = fmt.Sprintf("%f", v.FeePercent)
	case models.LegalAndProfessionalsPartners:
		address = v.PartnerAddress
		country = v.PartnerCountry
		feeFixed = fmt.Sprintf("%f", v.FeeFixed)
		feePercent = fmt.Sprintf("%f", v.FeePercent)
	case models.RatingAgency:
		address = v.AgencyAddress
		country = v.AgencyCountry
		feeFixed = fmt.Sprintf("%f", v.FeeFixed)
		feePercent = fmt.Sprintf("%f", v.FeePercent)
	case models.Trustees:
		address = v.TrusteeAddress
		country = v.TrusteeCountry
		feeFixed = fmt.Sprintf("%f", v.FeeFixed)
		feePercent = fmt.Sprintf("%f", v.FeePercent)
	}

	return address, country, feeFixed, feePercent
}

// GetOrganizationRootAdmin returns the ROOT_SUPER_ADMIN member of an organization as TrovoAdminInfo
func GetOrganizationRootAdmin(server *serverModels.Server, organizationID string) (*models.TrovoAdminInfo, error) {
	var member models.OrganizationMember
	if err := server.AdminDB.Where("organization_id = ? AND role = ?", organizationID, "ROOT_SUPER_ADMIN").First(&member).Error; err != nil {
		return &models.TrovoAdminInfo{}, fmt.Errorf("root admin not found: %v", err)
	}
	return &models.TrovoAdminInfo{
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Email:     member.Email,
	}, nil
}

// mapOrganizationWithAdmin maps &models.Organization{} to models.OrganizationResponse using TrovoAdminInfo directly
func mapOrganizationWithAdmin(organization *models.Organization, adminInfo *models.TrovoAdminInfo) *models.OrganizationResponse {
	return &models.OrganizationResponse{
		ID:                organization.ID,
		Name:              organization.Name,
		Email:             organization.Email,
		Type:              organization.Type,
		Status:            organization.Status,
		CreatedAt:         organization.CreatedAt,
		UpdatedAt:         organization.UpdatedAt,
		OrganizationCount: 0,
		TrovoAdminInfo:    *adminInfo,
		StakeholderID:     organization.StakeholderID,
		StakeholderType:   organization.StakeholderType,
		Level:             organization.Level,
	}
}

// DeactivateOrganization deactivates an organization and all its members
func DeactivateOrganization(server *serverModels.Server, organizationID string) (*models.Organization, error) {
	// Check if organization exists
	var existingOrg models.Organization
	if err := server.AdminDB.Where("id = ?", organizationID).First(&existingOrg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to find organization: %v", err)
	}

	// Check if organization is already inactive
	if existingOrg.Status == models.OrganizationStatusInactive {
		return nil, fmt.Errorf("organization is already inactive")
	}

	// Start transaction
	tx := server.AdminDB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Update organization status to INACTIVE
	updateData := map[string]interface{}{
		"status":     models.OrganizationStatusInactive,
		"updated_at": time.Now(),
	}

	if err := tx.Model(&existingOrg).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update organization status: %v", err)
	}

	// Also deactivate all organization members
	if err := tx.Model(&models.OrganizationMember{}).
		Where("organization_id = ?", organizationID).
		Updates(map[string]interface{}{
			"status":     models.OrganizationStatusInactive,
			"updated_at": time.Now(),
		}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update organization members status: %v", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &existingOrg, nil
}

// ValidateOrganizationExists checks if an organization exists and returns it
func ValidateOrganizationExists(server *serverModels.Server, organizationID string) (*models.Organization, error) {
	var org models.Organization
	if err := server.AdminDB.Where("id = ?", organizationID).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to find organization: %v", err)
	}
	return &org, nil
}
