package usermetrics

import (
	usermetricsDB "admin-panel-dashboard/internal/components/usermetrics/db"
	"admin-panel-dashboard/internal/server/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// trovoAdminOnly answers 403 unless the caller is a Trovo admin (the JWT
// middleware also admits organization members).
func trovoAdminOnly(c *gin.Context) bool {
	if c.GetString("auth_type") != "trovo_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only Trovo admins can manage fee exemptions"})
		return false
	}
	return true
}

// @Summary List fee-exempt accounts
// @Description Accounts that pay no platform service fees (swap, payment, patron, account recovery, sub-wallet, tokenization application, closed group). The tokenization issuing profile is always exempt and is not listed.
// @ID GetFeeExemptUsers
// @Tags Fees
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} response.Data{data=[]models.FeeExemptUser}
// @Failure 403 {object} map[string]string
// @Router /fee/exempt-users [get]
func GetFeeExemptUsersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !trovoAdminOnly(c) {
			return
		}
		users, err := usermetricsDB.ListFeeExemptUsers(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response.JSON(c, http.StatusOK, "Fee-exempt accounts fetched successfully", users, nil)
	}
}

// AddFeeExemptUserRequest is the account to exempt from platform service fees.
type AddFeeExemptUserRequest struct {
	Username string `json:"username" binding:"required"`
	Reason   string `json:"reason" binding:"required"`
}

// @Summary Exempt an account from platform service fees
// @Description Adds a Trovo account (by username) to the fee-exempt list, or updates its reason. Takes effect on its next fee-charging request.
// @ID AddFeeExemptUser
// @Tags Fees
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param body body AddFeeExemptUserRequest true "Account and reason"
// @Success 200 {object} response.Data{data=models.FeeExemptUser}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /fee/exempt-users [post]
func AddFeeExemptUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !trovoAdminOnly(c) {
			return
		}
		var req AddFeeExemptUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		row, err := usermetricsDB.AddFeeExemptUser(db, req.Username, req.Reason, c.GetString("trovo_admin_email"))
		if err != nil {
			if errors.Is(err, usermetricsDB.ErrFeeExemptUserNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "no Trovo account has that username"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response.JSON(c, http.StatusOK, "Account exempted from service fees", row, nil)
	}
}

// @Summary Remove a fee exemption
// @Description The account pays platform service fees again from its next request.
// @ID RemoveFeeExemptUser
// @Tags Fees
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param username path string true "Username"
// @Success 200 {object} response.Data
// @Failure 403 {object} map[string]string
// @Router /fee/exempt-users/{username} [delete]
func RemoveFeeExemptUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !trovoAdminOnly(c) {
			return
		}
		if err := usermetricsDB.RemoveFeeExemptUser(db, c.Param("username")); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response.JSON(c, http.StatusOK, "Fee exemption removed", nil, nil)
	}
}
