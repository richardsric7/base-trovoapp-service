package server

import (
	"admin-panel-dashboard/internal/components/admin_invite"
	usermetrics "admin-panel-dashboard/internal/components/usermetrics/db"
	userServices "admin-panel-dashboard/internal/components/users/services"
	p2pErrors "admin-panel-dashboard/internal/errors"
	"admin-panel-dashboard/internal/mail"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverModels "admin-panel-dashboard/internal/server/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Changes the Status an admin
// @Description Changes the admin status by updating their status to 'SUPER_ADMIN', 'EDIT_LEVEL_ADMIN', 'VIEW_ONLY_ADMIN', in the database. Only super admins can perform this action.
// @Tags Admins
// @Accept json
// @Produce json
// @Param payload body ChangeRolePayload true "Admin suspension details"
// @Success 200 {object} map[string]string "message": "Admin status updated to suspended"
// @Failure 400 {object} map[string]string "error": "Invalid request payload" or "Only 'suspended' status is allowed" or "Admin is already suspended"
// @Failure 401 {object} map[string]string "error": "Unauthorized access" or "Only super admins can change admin status"
// @Failure 404 {object} map[string]string "error": "Admin user does not exist"
// @Failure 500 {object} map[string]string "error": "Failed to retrieve admin information" or "InternalServerError" or "Failed to update admin status"
// @Router /admin/modify/status [patch]
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
func ChangeAdminStatus(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve admin information from context
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID
		adminInfo, done, _ := GetUserFromContext(userID, s.TrovoWalletDB, c)
		if done {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve admin information"})
			return
		}

		// Check if the current admin is a super admin
		isSuperAdmin, err := IsSuperAdmin(adminInfo.Email, s.AdminDB)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			return
		}

		if !isSuperAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Only super admins can change admin status"})
			return
		}

		// Bind JSON payload
		var payload ChangeRolePayload
		if err = c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		// DB call to fetch
		if payload.Role != string(models.EditLevelAdmin) &&
			payload.Role != string(models.SuperAdmin) &&
			payload.Role != string(models.ViewOnlyAdmin) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin role"})
			return
		}

		// Retrieve the admin to be suspended
		var admin models.AdminUser
		if err = s.AdminDB.First(&admin, "email = ?", payload.Email).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Admin user does not exist"})
			return
		}

		// Update the admin's status to suspended
		admin.Role = models.Role(payload.Role)
		if err = s.AdminDB.Save(&admin).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin status"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Admin user has been suspended"})
	}
}

// func InviteAdminUser(c *gin.Context) {
//	//ad, _ := middleware.ExtractTokenMetadata(c.Request)
//	//userID := ad.UserID
//	//
//	//userInfo, err := userServices.GetUser(userID, walletDb)
//	//if err != nil {
//	//	log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)
//	//
//	//	var ex p2pErrors.GenericError
//	//	var ok bool
//	//
//	//	ex, ok = err.(p2pErrors.GenericError)
//	//	var statusCode int = 0
//	//	var response interface{}
//	//
//	//	if ok {
//	//		statusCode = ex.HTTPCode()
//	//		response = ex.JSONError()
//	//	} else {
//	//		statusCode = http.StatusBadRequest
//	//		response = gin.H{"error": err.Error()}
//	//	}
//
//	//	c.JSON(statusCode, response)
//	//	return
//	//}
//	// Retrieve data with pagination and filtering
//
//	payload, err := models.validateInviteAdminPayload(c)
//	if err != nil {
//		return
//	}
//
//	if models.checkExistingWallet(payload.Email) {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already has an existing wallet"})
//		return
//	}
//
//	username := generateUsername(payload.Email) // Assuming a function to generate a username
//	permission := "default_permission"          // Set default permission or get from payload
//
//	if err := models.storeAdminData(payload.Email, username, permission); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store admin data"})
//		return
//	}
//
//	if err := models.sendInvitationEmail(payload.Email); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send invitation email"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Invitation sent successfully"})
//
//	//req := models.UserProfile{
//	//	ID:       id,
//	//	Username: usernameFilter,
//	//	Email:    emailFilter,
//	//}
//	//userInfo, err := usermetrics.GetUserProfileInfo(req)
//	//if err != nil {
//	//	log.Println("[METRICS] error:", err)
//	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
//	//	return
//	//}
//	//response.JSON(c, http.StatusOK, "userInfp fetched successfully", userInfo, nil)
//}

type InviteAdminPayload struct {
	EmailOrUsername string `json:"email_or_username" binding:"required"`
	Role            string `json:"admin_role" binding:"required"`
}

// AddAdmin handles the invitation of a new admin user.
// @Summary      Add Admin
// @Description  Usable by SuperAdmin Only. Invite a new admin user by providing a trovo user's email and expected role accepted admin_roles are "SUPER_ADMIN" , "EDIT_LEVEL_ADMIN",  "VIEW_ONLY_ADMIN"
// @Tags         Admins
// @Accept       json
// @Produce      json
// @Param        payload  body      InviteAdminPayload  true  "Admin invitation payload"
// @Success      200      {object}  map[string]string   "Admin successfully created"
// @Failure      400      {object}  map[string]string   "Invalid payload or user does not exist"
// @Failure      500      {object}  map[string]string   "Internal server error"
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Router       /admin/invite [post]
func AddAdmin(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID
		adminInfo, done, _ := GetUserFromContext(userID, s.TrovoWalletDB, c)
		if done {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		var payload InviteAdminPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}

		// Validate EmailOrUsername field
		if strings.TrimSpace(payload.EmailOrUsername) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or username must be provided"})
			return
		}

		// Validate role
		if payload.Role != string(models.EditLevelAdmin) &&
			payload.Role != string(models.SuperAdmin) &&
			payload.Role != string(models.ViewOnlyAdmin) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin role"})
			return
		}

		isSuperAdmin, err := IsSuperAdmin(adminInfo.Email, s.AdminDB)
		if err != nil || !isSuperAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			return
		}

		var user models.User
		if err := s.TrovoWalletDB.
			Where("email = ? OR username = ?", payload.EmailOrUsername, payload.EmailOrUsername).
			First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		// Check for existing admin with this email or username
		var existingAdmin models.AdminUser
		if err := s.AdminDB.Where("email = ? OR username = ?", user.Email, user.Username).First(&existingAdmin).Error; err == nil {
			msg := fmt.Sprintf("This admin already exists with role %s", existingAdmin.Role)
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		// Create admin user
		admin := models.AdminUser{
			Username:     user.Username,
			Email:        user.Email,
			Role:         models.Role(payload.Role),
			Status:       string(models.ActiveAdmin),
			IsAdmin:      true,
			WalletUserID: user.ID,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
		}

		if err = s.AdminDB.Create(&admin).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin user"})
			return
		}

		// Send email notification
		go func() {
			emailBody := fmt.Sprintf(
				"Hi %v,\nYou have been added as an admin user. Please log in on %v with your Trovo app.\nWe're glad you are here.\nThe Trovo Team",
				user.FirstName, os.Getenv("TROVO_WALLET_BASE_URL"),
			)
			_, response, err := mail.SendEmailWithMailgunTemplate(admin, emailBody, "Admin Email Invite")
			if err != nil {
				log.Println("Error sending email:", err)
				return
			}
			log.Println("Email sent successfully:", response)
		}()

		c.JSON(http.StatusOK, gin.H{"message": "Admin user created successfully"})
	}
}

type EmailSend struct {
	FirstName string `json:"first_name" binding:"required"`
}

//// SuspendAdmin suspends an admin user, logs the suspension reason, and updates the admin status.
//// @Summary Suspend an admin
//// @Description Suspends an admin by updating their status to 'SUSPENDED' in the database. Only super admins can perform this action.
//// @Tags Admins
//// @Accept json
//// @Produce json
//// @Param payload body models.SuspendAdminPayload true "Admin suspension details"
//// @Success 200 {object} map[string]string "message": "Admin status updated to suspended"
//// @Failure 400 {object} map[string]string "error": "Invalid request payload" or "Only 'suspended' status is allowed" or "Admin is already suspended"
//// @Failure 401 {object} map[string]string "error": "Unauthorized access" or "Only super admins can change admin status"
//// @Failure 404 {object} map[string]string "error": "Admin user does not exist"
//// @Failure 500 {object} map[string]string "error": "Failed to retrieve admin information" or "InternalServerError" or "Failed to update admin status"
//// @Router /admin/suspend [patch]
//// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// func SuspendAdmin(s *serverModels.Server) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		ad, _ := middleware.ExtractTokenMetadata(c.Request)
//		userID := ad.UserID
//		superAdmin, done, _ := GetUserFromContext(userID, s.TrovoWalletDB, c)
//		if done {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve admin information"})
//			return
//		}
//
//		// Bind JSON payload
//		var payload models.SuspendAdminPayload
//		if err := c.ShouldBindJSON(&payload); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
//			return
//		}
//
//		// Validate the status to be 'suspended'
//		//if payload.Status != string(models.Suspended) {
//		//	c.JSON(http.StatusBadRequest, gin.H{"error": "Only 'suspended' status is allowed"})
//		//	return
//		//}
//
//		// Retrieve the admin to be suspended
//		var admin models.AdminUser
//		if err := s.AdminDB.First(&admin, "email = ?", payload.Email).Error; err != nil {
//			c.JSON(http.StatusNotFound, gin.H{"error": "Admin user does not exist"})
//			return
//		}
//
//		// Check if the admin is already suspended
//		if admin.Status == string(models.Suspended) {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Admin is already suspended"})
//			return
//		}
//
//		// Update the admin's status to suspended
//		admin.Status = string(models.Suspended)
//		if err := s.AdminDB.Save(&admin).Error; err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin status"})
//			return
//		}
//
//		// Log the suspension activity in the suspension history table
//		suspensionRecord := models.AdminSuspensionHistory{
//			Username:           admin.Username,
//			Email:              admin.Email,
//			ActionPerformedBy:  superAdmin.Email, // Admin performing the action
//			Reason:             payload.Reason,   // Suspension reason
//			SuspensionDateTime: time.Now(),
//			ActionType:         models.SUSPENDED, // Action type (SUSPENDED)
//		}
//
//		if err := s.AdminDB.Create(&suspensionRecord).Error; err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log suspension activity"})
//			return
//		}
//
//		c.JSON(http.StatusOK, gin.H{"message": "Admin user has been suspended"})
//	}
//}

// @Summary Suspend an admin
// @Description Suspends an admin by updating their status to 'SUSPENDED' in the database. Requires admin-level authorization.
// @Tags Admins
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param payload body models.SuspendAdminPayload true "Admin suspension details, including SuspensionReasonID and SuspensionNote"
// @Success 200 {object} map[string]string  "Admin user has been suspended"
// @Failure 400 {object} models.ErrorResponse "Invalid request payload or admin is already suspended"
// @Failure 401 {object} models.ErrorResponse "Unauthorized access"
// @Failure 404 {object} models.ErrorResponse "Admin user does not exist"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /admin/suspend [patch]
func SuspendAdmin(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID
		superAdmin, done, _ := GetUserFromContext(userID, s.TrovoWalletDB, c)
		if done {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve admin information"})
			return
		}

		// Bind JSON payload
		var payload models.SuspendAdminPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Retrieve the admin to be suspended
		var admin models.AdminUser
		if err := s.AdminDB.First(&admin, "email = ?", payload.Email).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Admin user does not exist"})
			return
		}

		// Check if the admin is already suspended
		if admin.Status == string(models.Suspended) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Admin is already suspended"})
			return
		}

		// DB call to suspend | the admin's status to suspended
		admin.Status = string(models.Suspended)
		if err := s.AdminDB.Save(&admin).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin status"})
			return
		}
		suspensionReason, err := usermetrics.GetSuspensionReasonByID(payload.SuspensionReasonID, s.AdminDB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve suspension reason"})
			return
		}

		// Log the suspension activity in the suspension history table
		suspensionRecord := models.AdminSuspensionHistory{
			Username:           admin.Username,
			Email:              admin.Email,
			ActionPerformedBy:  superAdmin.Email,
			Reason:             suspensionReason,
			SuspensionNote:     payload.SuspensionNote,
			SuspensionDateTime: time.Now(),
			ActionType:         models.SUSPENDED,
		}

		if err := s.AdminDB.Create(&suspensionRecord).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log suspension activity"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Admin user has been suspended"})
	}
}

// SuspendUser suspends a normal (Trovo Wallet) user. Distinct from
// LiftUserSuspension below (previously this was one toggle endpoint,
// SuspendOrReactivateUser - fragile against a stale read or a double
// click, and its reason was a lookup into the never-seeded
// UserSuspensionReason table, so every call to it errored). While
// suspended, app-backend blocks the user's wallet(s) from either side of
// any transaction (payments, swaps, P2P, shared-wallet access grants,
// tokenized-asset subscriptions/minting) - see User.EnsureNotSuspended in
// app-backend.
// @Summary Suspend a user
// @Description Suspends a normal user, logging the mandatory reason. Blocks the user's wallet(s) from any transaction until lifted.
// @Tags Users
// @Accept json
// @Produce json
// @Param payload body models.SuspendOrLiftUserPayload true "Email of the user to suspend, and the mandatory reason"
// @Success 200 {object} map[string]string "message": "User has been suspended"
// @Failure 400 {object} map[string]string "error": "Invalid request payload" or "User is already suspended"
// @Failure 401 {object} map[string]string "error": "Unauthorized access"
// @Failure 404 {object} map[string]string "error": "User does not exist"
// @Failure 500 {object} map[string]string "error": "Failed to update user status" or "Failed to log suspension activity"
// @Router /admin/users/suspend [patch]
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
func SuspendUser(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		adminUser, done, _ := GetUserFromContext(ad.UserID, s.TrovoWalletDB, c)
		if done {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve admin information"})
			return
		}

		var payload models.SuspendOrLiftUserPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		var user models.User
		if err := s.TrovoWalletDB.First(&user, "email = ?", payload.Email).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
			return
		}
		if user.Suspended == 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is already suspended"})
			return
		}

		// Map-based update (not a struct Save) so this write only ever
		// touches these two columns, regardless of what else the User
		// struct's field/tag set does or doesn't cover.
		if err := s.TrovoWalletDB.Model(&models.User{}).Where("email = ?", payload.Email).
			Updates(map[string]interface{}{"suspended": 1, "suspension_reason": payload.Reason}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
			return
		}

		userSuspensionRecord := models.UserSuspensionHistory{
			Username:           user.Username,
			Email:              user.Email,
			ActionPerformedBy:  adminUser.Email,
			Reason:             payload.Reason,
			SuspensionDateTime: time.Now(),
			ActionType:         models.SUSPENDED,
		}
		if err := s.AdminDB.Create(&userSuspensionRecord).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log suspension activity"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User has been suspended"})
	}
}

// LiftUserSuspension reactivates a suspended normal user, logging the
// mandatory reason for lifting the suspension. See SuspendUser above.
// @Summary Lift a user's suspension
// @Description Reactivates a suspended user, logging the mandatory reason for lifting the suspension.
// @Tags Users
// @Accept json
// @Produce json
// @Param payload body models.SuspendOrLiftUserPayload true "Email of the user to reactivate, and the mandatory reason for lifting the suspension"
// @Success 200 {object} map[string]string "message": "User's suspension has been lifted"
// @Failure 400 {object} map[string]string "error": "Invalid request payload" or "User is not currently suspended"
// @Failure 401 {object} map[string]string "error": "Unauthorized access"
// @Failure 404 {object} map[string]string "error": "User does not exist"
// @Failure 500 {object} map[string]string "error": "Failed to update user status" or "Failed to log suspension activity"
// @Router /admin/users/lift-suspension [patch]
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
func LiftUserSuspension(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		adminUser, done, _ := GetUserFromContext(ad.UserID, s.TrovoWalletDB, c)
		if done {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve admin information"})
			return
		}

		var payload models.SuspendOrLiftUserPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		var user models.User
		if err := s.TrovoWalletDB.First(&user, "email = ?", payload.Email).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
			return
		}
		if user.Suspended == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not currently suspended"})
			return
		}

		if err := s.TrovoWalletDB.Model(&models.User{}).Where("email = ?", payload.Email).
			Updates(map[string]interface{}{"suspended": 0, "suspension_reason": ""}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
			return
		}

		userSuspensionRecord := models.UserSuspensionHistory{
			Username:           user.Username,
			Email:              user.Email,
			ActionPerformedBy:  adminUser.Email,
			Reason:             payload.Reason,
			SuspensionDateTime: time.Now(),
			ActionType:         models.ACTIVE,
		}
		if err := s.AdminDB.Create(&userSuspensionRecord).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log suspension activity"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User's suspension has been lifted"})
	}
}

// GetUserSuspensionHistoryByEmail retrieves the logged suspend/lift
// history (each with its admin-supplied reason) for a normal user - the
// User-side counterpart to GetSuspensionHistoryByEmail (admin/staff)
// above.
// @Summary Get a user's suspension history
// @Description Retrieves suspend/lift activity (each with its logged reason) for a specific normal user by email.
// @Tags Users
// @Accept json
// @Produce json
// @Param email path string true "User email"
// @Success 200 {array} models.UserSuspensionHistory
// @Failure 400 {object} map[string]string "error": "Invalid or missing email"
// @Router /users/suspension-history/{email} [get]
func GetUserSuspensionHistoryByEmail(adminDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Param("email")
		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing email"})
			return
		}

		var history []models.UserSuspensionHistory
		if err := adminDB.Where("email = ?", email).Order("suspension_date_time desc").Find(&history).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve suspension history"})
			return
		}

		c.JSON(http.StatusOK, history)
	}
}

// UnsuspendAdmin unsuspends an admin user, logs the suspension reason, and updates the admin status.
// @Summary Suspend an admin
// @Description Suspends an admin by updating their status to 'SUSPENDED' in the database. Only super admins can perform this action.
// @Tags Admins
// @Accept json
// @Produce json
// @Param payload body models.UnsuspendAdminPayload true "Admin un-suspension details"
// @Success 200 {object} map[string]string "message": "Admin status updated to suspended"
// @Failure 400 {object} map[string]string "error": "Invalid request payload" or "Only 'suspended' status is allowed" or "Admin is already suspended"
// @Failure 401 {object} map[string]string "error": "Unauthorized access" or "Only super admins can change admin status"
// @Failure 404 {object} map[string]string "error": "Admin user does not exist"
// @Failure 500 {object} map[string]string "error": "Failed to retrieve admin information" or "InternalServerError" or "Failed to update admin status"
// @Router /admin/unsuspend [patch]
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
func UnsuspendAdmin(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {

		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID
		superAdmin, done, _ := GetUserFromContext(userID, s.TrovoWalletDB, c)
		if done {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve admin information"})
			return
		}

		// Bind JSON payload
		var payload models.UnsuspendAdminPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Validate the status to be 'suspended'
		// if payload.Status != string(models.Suspended) {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": "Only 'suspended' status is allowed"})
		//	return
		//}

		// Retrieve the admin to be suspended
		var admin models.AdminUser
		if err := s.AdminDB.First(&admin, "email = ?", payload.Email).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Admin user does not exist"})
			return
		}

		// Check if the admin is already suspended
		if admin.Status == string(models.ACTIVE) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Admin is already active"})
			return
		}

		// Update the admin's status to suspended
		admin.Status = string(models.ACTIVE)
		if err := s.AdminDB.Save(&admin).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin status"})
			return
		}

		// Log the suspension activity in the suspension history table
		unsuspensionRecord := models.AdminSuspensionHistory{
			Username:           admin.Username,
			Email:              admin.Email,
			ActionPerformedBy:  superAdmin.Email, // Admin performing the action
			Reason:             "REACTIVATED",    // Suspension reason
			SuspensionDateTime: time.Now(),
			ActionType:         models.REACTIVATED, // Action type (REACTIVATED)
		}

		if err := s.AdminDB.Create(&unsuspensionRecord).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log suspension activity"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Admin user has been REACTIVATED"})
	}
}

func GetUserFromContext(userID string, walletDB *gorm.DB, c *gin.Context) (models.UserInfo, bool, error) {

	adminInfo, err := userServices.GetUser(userID, walletDB)
	if err != nil {
		log.Println("[METRICS] error for user:", adminInfo.Username, "error: ", err)

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
			response = gin.H{"error": err.Error()}
		}

		c.JSON(statusCode, response)
		return models.UserInfo{}, true, nil
	}
	log.Println("checking user data: ", adminInfo, userID)
	adminInfo.ID = userID
	return adminInfo, false, err
}

// @Summary Get all suspension history
// @Description Retrieves all suspension activities from the history table.
// @Tags Admins
// @Accept json
// @Produce json
// @Success 200 {array} models.AdminSuspensionHistory
// @Failure 500 {object} map[string]string "error": "Internal server error"
// @Router /suspension/history [get]
func GetAllSuspensionHistory(adminDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var history []models.AdminSuspensionHistory
		if err := adminDB.Find(&history).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve suspension history"})
			return
		}

		c.JSON(http.StatusOK, history)
	}
}

// @Summary Get suspension history by email
// @Description Retrieves suspension activities of a specific admin by their email.
// @Tags Admins
// @Accept json
// @Produce json
// @Param email query string true "Admin email"
// @Success 200 {array} models.AdminSuspensionHistory
// @Failure 400 {object} map[string]string "error": "Invalid or missing email"
// @Failure 404 {object} map[string]string "error": "No suspension history found for this email"
// @Failure 500 {object} map[string]string "error": "Internal server error"
// @Router /suspension/history/{email} [get]
func GetSuspensionHistoryByEmail(adminDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Param("email")
		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing email"})
			return
		}

		var history []models.AdminSuspensionHistory
		if err := adminDB.Where("email = ?", email).Find(&history).Error; err != nil || len(history) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No suspension history found for this email"})
			return
		}

		c.JSON(http.StatusOK, history)
	}
}

func IsSuperAdmin(email string, db *gorm.DB) (bool, error) {
	var admin models.AdminUser
	if err := db.First(&admin, "email = ?", email).Error; err != nil {
		return false, err
	}
	if admin.Role == models.SuperAdmin {
		return true, nil
	}
	return false, nil
}

func IsAdmin(email string, db *gorm.DB) (bool, error) {
	var admin models.AdminUser
	if err := db.First(&admin, "email = ?", email).Error; err != nil {
		return false, err
	}
	return true, nil
}

type InvitationResponsePayload struct {
	Email  string `json:"email" binding:"required,email"`
	Accept bool   `json:"accept" binding:"required"`
}

func AcceptRejectInvitation(adminDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload InvitationResponsePayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}

		var admin models.AdminUser
		if err := adminDB.First(&admin, "email = ?", payload.Email).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Admin user not found"})
			return
		}

		if payload.Accept {
			admin.Role = models.EditLevelAdmin
			adminDB.Save(&admin)
		} else {
			adminDB.Delete(&admin)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Invitation response recorded"})
	}
}

type ChangeRolePayload struct {
	Email string `json:"email" binding:"required"`
	Role  string `json:"admin_role" binding:"required"`
}

// endpoint for admin user configs (rol0

// no need to accept email
//SUSPEND ENDPOINT / REMOVE SUSPENSION (separate)

// paginated
// GET /admin?page=2&pageSize=5
// ListAdminsHandler returns a handler function to list admin users with pagination.

// ListAdminsHandler returns a handler function to list admin users with pagination.
// @Summary      List Admins
// @Description  Retrieve a paginated list of admin users
// @Tags         Admins
// @Accept       json
// @Produce      json
// @Param        page      query    int     false  "Page number (must be >= 1)"          default(1)
// @Param        pageSize  query    int     false  "Number of admins per page (max 100)" default(10)
// @Success      200       {object} map[string]interface{}
// @Failure      400       {object} map[string]string
// @Failure      500       {object} map[string]string
// @Router       /admin/list [get]
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
func ListAdminsHandler(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		const (
			defaultPageSize = 10
			maxPageSize     = 100
		)

		// Get query parameters with defaults
		page := c.DefaultQuery("page", "1")
		pageSize := c.DefaultQuery("pageSize", fmt.Sprintf("%d", defaultPageSize))

		// Parse and validate query parameters
		pageInt, err := strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
			return
		}

		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil || pageSizeInt < 1 || pageSizeInt > maxPageSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
			return
		}

		limit := pageSizeInt
		offset := (pageInt - 1) * pageSizeInt

		// Fetch total count of admin users
		var total int64
		if err := s.AdminDB.Model(&models.AdminUser{}).Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch admin users count"})
			return
		}

		// Fetch paginated admins
		admins, err := admin_invite.FetchAdmins(limit, offset, s.AdminDB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch admin users"})
			return
		}

		totalPages := int((total + int64(pageSizeInt) - 1) / int64(pageSizeInt)) // Ceiling division for pages

		// Return response with pagination metadata
		c.JSON(http.StatusOK, gin.H{
			"admins": admins,
			"pagination": gin.H{
				"page":       pageInt,
				"pageSize":   pageSizeInt,
				"total":      total,
				"totalPages": totalPages,
			},
		})
	}
}

type RemoveAdminPayload struct {
	Email string `json:"email" binding:"required,email"`
}

// ONLY SUPER ADMIN
func RemoveAdmin(adminDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload RemoveAdminPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}

		if err := adminDB.Delete(&models.AdminUser{}, "email = ?", payload.Email).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove admin user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Admin removed successfully"})
	}
}

/*
log.Println("not a super admin")
*/
// test other untested endpoints
/*
up until now, look at how i am seeding my configurations at startup in the following methods.
But now, i want all my configurations to be totally configurable by admins by creating crud endpoints for them.
Note: i have to log the history of admin details that made a particular change to the configuration
Note that i want to be able to create, read, update and delete the following configurations:

CRUD OPERATION FOR Roles
const (
	SuperAdmin     Role = "SUPER_ADMIN"
	EditLevelAdmin Role = "EDIT_LEVEL_ADMIN"
	ViewOnlyAdmin  Role = "VIEW_ONLY_ADMIN"
	Suspended      Role = "SUSPENDED"
)
CRUD OPERATION FOR SuspensionReasons
type UserSuspensionReason struct {
	ID                 uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Reason             string `gorm:"type:varchar(255);not null" json:"reason"`
}


func SeedSuspensionReasons(db *gorm.DB) error {
	suspensionReasons := []models.UserSuspensionReason{
		{ID: 1, Reason: "Violation of terms of service"},
		{ID: 2, Reason: "Violation of community guidelines"},
		{ID: 3, Reason: "Violation of KYC/AML policy"},
		{ID: 4, Reason: "Violation of security policy"},
		{ID: 5, Reason: "Violation of privacy policy"},
		{ID: 6, Reason: "Violation of trading policy"},
		{ID: 7, Reason: "Violation of payment policy"},
		{ID: 8, Reason: "Violation of dispute resolution policy"},
		{ID: 9, Reason: "Violation of support policy"},
	}

	for _, reason := range suspensionReasons {
		var existingReason models.UserSuspensionReason
		// Check if the reason already exists in the database
		if err := db.Where("id = ?", reason.ID).First(&existingReason).Error; err != nil {
			// If not found, create a new reason
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err = db.Create(&reason).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
*/
