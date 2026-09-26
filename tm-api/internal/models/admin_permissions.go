package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string
type AdminStatus string
type UserStatus string

const (
	SuperAdmin             Role = "SUPER_ADMIN"
	OrganizationSuperAdmin Role = "ROOT_SUPER_ADMIN"
	OrganizationMemberRole Role = "MEMBER"
	EditLevelAdmin         Role = "EDIT_LEVEL_ADMIN"
	ViewOnlyAdmin          Role = "VIEW_ONLY_ADMIN"
	Suspended              Role = "SUSPENDED"
)

const (
	ActiveUser    UserStatus = "ACTIVE"
	SuspendedUser UserStatus = "SUSPENDED"
)

const (
	ActiveAdmin    AdminStatus = "ACTIVE"
	SuspendedAdmin AdminStatus = "SUSPENDED"
)
const (
	SUSPENDED   = "SUSPENDED"
	REACTIVATED = "REACTIVATED"
	ACTIVE      = "ACTIVE"
)

type SuspendAdminPayload struct {
	Email              string `json:"email" binding:"required"`
	SuspensionReasonID uint   `json:"suspension_reason_id" binding:"required"`
	SuspensionNote     string `json:"suspension_note" binding:"required"`
}

// RoleConfig model for configurable roles
//
//	type RoleConfig struct {
//		ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
//		RoleName    string `gorm:"type:varchar(100);not null;unique" json:"roleName"`
//		Description string `gorm:"type:varchar(255)" json:"description"`
//	}
type UserSuspensionReason struct {
	ID                 uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Reason             string `gorm:"type:varchar(255);not null" json:"reason"`
	Username           string `gorm:"not null"`
	Email              string
	ActionPerformedBy  string    `gorm:"not null"`
	SuspensionDateTime time.Time `gorm:"not null"`
	ActionType         string    `gorm:"not null"` // SUSPENDED, REACTIVATED
	SuspensionNote     string
}

type CreateRoleRequest struct {
	RoleName    string `json:"roleName" binding:"required"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	RoleName    string `json:"roleName"`
	Description string `json:"description"`
}

type CreateSuspensionReasonRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type UpdateSuspensionReasonRequest struct {
	Reason string `json:"reason"`
}

// RoleConfig represents the model for a role in the database
type RoleConfig struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	RoleName    string `gorm:"uniqueIndex;not null" json:"roleName"`
	Description string `json:"description"`
}

// UserSuspensionReason represents the model for a suspension reason in the database
type SuspensionReason struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Reason      string `gorm:"not null;uniqueIndex:idx_reason_target" json:"reason"`
	RoleName    string `json:"roleName"`
	Description string `json:"description"`
	Target      string `gorm:"not null;uniqueIndex:idx_reason_target" json:"target"` // Shared unique index with "reason"
}

//// UserSuspensionReason model for configurable suspension reasons
// type UserSuspensionReason struct {
//	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
//	Reason string `gorm:"type:varchar(255);not null" json:"reason"`
//}

// AdminChangeLog model to track changes made by admins
type AdminChangeLog struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AdminEmail      string    `gorm:"type:varchar(100);not null" json:"adminEmail"`
	ChangeType      string    `gorm:"type:varchar(50);not null"` // e.g., CREATE, UPDATE, DELETE
	Configuration   string    `gorm:"type:varchar(50);not null"` // e.g., Role or SuspensionReason
	ConfigurationID uint      `json:"configurationId"`
	ChangeDetails   string    `gorm:"type:text" json:"changeDetails"` // e.g., details of what was changed
	CreatedAt       time.Time `json:"createdAt"`
}

// Examples of SusPensionReasonID:
// 1 - Violation of terms of service
// 2 - Violation of community guidelines
// 3 - Violation of KYC/AML policy
// 4 - Violation of security policy
// 5 - Violation of privacy policy
// 6 - Violation of trading policy
// 7 - Violation of payment policy
// 8 - Violation of dispute resolution policy
// 9 - Violation of support policy

type UnsuspendAdminPayload struct {
	Email string `json:"email" binding:"required"`
}

type AdminSuspensionHistory struct {
	ID                 uint   `gorm:"primaryKey"`
	Username           string `gorm:"not null"`
	Email              string
	ActionPerformedBy  string    `gorm:"not null"`
	Reason             string    `gorm:"not null"`
	SuspensionDateTime time.Time `gorm:"not null"`
	ActionType         string    `gorm:"not null"` // SUSPENDED, REACTIVATED
	SuspensionNote     string
}

type AdminUser struct {
	ID           uint `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	FirstName    string
	LastName     string
	Status       string `gorm:"not null"`
	Role         Role   `gorm:"not null"`
	IsAdmin      bool   `gorm:"not null"`
	WalletUserID string `gorm:"size:100"`
}

type AdminPermission struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

type RolePermission struct {
	Role         Role `gorm:"primaryKey;autoIncrement:false"`
	PermissionID uint `gorm:"primaryKey;autoIncrement:false"`
}

// SuspendOrLiftUserPayload is the request body for the distinct
// suspend-user and lift-user-suspension endpoints. Unlike the legacy
// toggle (SuspendNormalUserPayload above), each direction is its own
// endpoint, and Reason is a plain mandatory free-text explanation rather
// than a lookup into the UserSuspensionReason table (which is never
// seeded - see SuspendOrReactivateUser's history).
type SuspendOrLiftUserPayload struct {
	Email  string `json:"email" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}

type UserSuspensionHistory struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username           string    `gorm:"type:varchar(255);not null" json:"username"`
	Email              string    `gorm:"type:varchar(255);not null" json:"email"`
	ActionPerformedBy  string    `gorm:"type:varchar(255);not null" json:"actionPerformedBy"`
	Reason             string    `gorm:"type:varchar(255);not null" json:"reason"`
	Note               string    `gorm:"type:varchar(255)" json:"note"`
	SuspensionDateTime time.Time `json:"suspensionDateTime"`
	ActionType         string    `gorm:"type:varchar(20);not null" json:"actionType"` // SUSPENDED or REACTIVATED
}

func HasPermission(db *gorm.DB, user *AdminUser, permissionName string) bool {
	var permission AdminPermission
	if err := db.First(&permission, "name = ?", permissionName).Error; err != nil {
		return false
	}

	var count int64
	db.Model(&RolePermission{}).Where("role = ? AND permission_id = ?", user.Role, permission.ID).Count(&count)
	return count > 0
}

// PaymentHistory holds payment information
type PaymentHistory struct {
	ID                    string
	TransactionType       string    `gorm:"index:idx_payment_history_unique_key,unique"`
	TransactionDate       time.Time `json:"transactionDate" gorm:"index:idx_payment_history_tx_time"`
	From                  *string   `json:"from" gorm:"size:150;index:idx_payment_history_from;null"` // trovoWallet alias and name
	FromAddress           string    `json:"fromAddress" gorm:"size:150;index:idx_payment_history_from_pk;not null"`
	To                    *string   `json:"to" gorm:"size:100;index:idx_payment_history_to;null"` // trovoWallet alias and name
	ToAddress             string    `json:"toAddress" gorm:"size:100;index:idx_payment_history_to_pk;not null;"`
	Memo                  *string   `json:"memo" gorm:"size:28;null"`
	ContractAddress       *string   `json:"contractAddress" gorm:"size:56;null;"`
	AssetCode             string    `json:"assetCode" gorm:"size:12;not null;"`
	Amount                string    `json:"amount" gorm:"index:idx_amount_ph"`
	TransactionID         string    `json:"transactionId" gorm:"size:70;not null;index:idx_payment_history_txid;index:idx_payment_history_unique_key,unique"`
	PT                    string    `json:"-" gorm:"size:70;not null;index:idx_payment_history_unique_key,unique;"`
	SourceAccountSequence string    `json:"-" gorm:"size:70;not null;index:idx_payment_history_unique_key,unique;"`
}

// PaymentHistoryJSON holds payment information in json format
type PaymentHistoryJSON struct {
	TransactionDate time.Time `json:"transactionDate"`
	TransactionType string    `json:"transactionType"`
	From            string    `json:"from"` // trovoWallet alias and name
	FromAddress     string    `json:"fromAddress"`
	To              string    `json:"to"` // trovoWallet alias and name
	ToAddress       string    `json:"toAddress"`
	Memo            string    `json:"memo"`
	ContractAddress string    `json:"contractAddress"`
	AssetCode       string    `json:"assetCode"`
	Amount          string    `json:"amount"`
	TransactionID   string    `json:"transactionId"`
}

// PaymentHistoryRequest holds the request parameters for fetching payment history
type PaymentHistoryRequest struct {
	Page            int    `json:"page"`
	PageSize        int    `json:"pageSize"`
	TransactionType string `json:"transactionType"`
	TransactionDate string `json:"transactionDate"`
	From            string `json:"from"`
	To              string `json:"to"`
	Memo            string `json:"memo"`
	ContractAddress string `json:"contractAddress"`
	AssetCode       string `json:"assetCode"`
	Search          string `json:"search"`
}

// PaymentHistoryResponse represents the paginated response for payment history
type PaymentHistoryResponse struct {
	Data     []PaymentHistoryJSON `json:"data"`
	Total    int                  `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"pageSize"`
}
