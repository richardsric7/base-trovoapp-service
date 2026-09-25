package dto

import "admin-panel-dashboard/internal/models"

// CreateSuspensionReasonRequest represents the request payload for creating a suspension reason
type CreateSuspensionReasonRequest struct {
	ID     uint   `json:"id"`
	Reason string `json:"reason" binding:"required"`
}

// CreateSuspensionReasonResponse represents the response payload for a created suspension reason
type CreateSuspensionReasonResponse struct {
	ID     uint   `json:"id"`
	Reason string `json:"reason"`
}

// UpdateSuspensionReasonRequest represents the request payload for updating a suspension reason
type UpdateSuspensionReasonRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// UpdateSuspensionReasonResponse represents the response payload for an updated suspension reason
type UpdateSuspensionReasonResponse struct {
	ID     uint   `json:"id"`
	Reason string `json:"reason"`
}

// GetAllConfigurationsResponse represents the response payload for retrieving suspension reasons
type GetAllConfigurationsResponse struct {
	SuspensionReason []models.UserSuspensionReason
	Roles            []models.RoleConfig
}

// GetSuspensionReasonResponse represents the response payload for retrieving suspension reasons
type GetSuspensionReasonResponse struct {
	ID     uint   `json:"id"`
	Reason string `json:"reason"`
}

// DeleteSuspensionReasonResponse represents the response payload for deleting a suspension reason
type DeleteSuspensionReasonResponse struct {
	Message string `json:"message"`
}

type CreateRoleRequest struct {
	RoleName    string `json:"roleName" binding:"required"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	RoleName    string `json:"roleName"`
	Description string `json:"description"`
}
