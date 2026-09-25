package models

import "time"

const (
	MintingRoleInitiator = "initiator"
	MintingRoleApprover  = "approver"
)

// TokenizationMintingInitiator represents a tokenization minting initiator
type TokenizationMintingInitiator struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	Initiator string `gorm:"not null;unique" json:"initiator"`
	// CreatedAt time.Time `json:"created_at"`
	//UpdatedAt time.Time `json:"updated_at"`
}

// TokenizationMintingApprover represents a tokenization minting approver
type TokenizationMintingApprover struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Approver string `gorm:"not null;unique" json:"approver"`
	// CreatedAt time.Time `json:"created_at"`
	//UpdatedAt time.Time `json:"updated_at"`
}

// CreateMintingUserRequest represents the request to create a new minting user
type CreateMintingUserRequest struct {
	Username string `json:"username" binding:"required" example:"john.doe"`
	Role     string `json:"role" binding:"required,oneof=initiator approver" example:"initiator" enums:"initiator,approver"`
}

// DeleteMintingUserRequest represents the request to delete a minting user
type DeleteMintingUserRequest struct {
	ID   int    `json:"id" binding:"required" example:"1"`
	Role string `json:"role" binding:"required,oneof=initiator approver" example:"initiator" enums:"initiator,approver"`
}

// SearchMintingUserRequest represents the request to search minting users
type SearchMintingUserRequest struct {
	Query    string `form:"query" example:"john"`
	Role     string `form:"role" binding:"omitempty,oneof=initiator approver" example:"initiator" enums:"initiator,approver"`
	Page     int    `form:"page" binding:"required,min=1" example:"1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100" example:"10"`
}

// MintingUserResponse represents the response for minting user operations
type MintingUserResponse struct {
	ID        int       `json:"id" example:"1"`
	Username  string    `json:"username" example:"john.doe"`
	Role      string    `json:"role" example:"initiator" enums:"initiator,approver"`
	CreatedAt time.Time `json:"created_at" example:"2024-03-20T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-03-20T10:00:00Z"`
}

// PaginatedMintingUserResponse represents a paginated response of minting users
type PaginatedMintingUserResponse struct {
	Data       []MintingUserResponse `json:"data"`
	Total      int64                 `json:"total" example:"100"`
	Page       int                   `json:"page" example:"1"`
	PageSize   int                   `json:"page_size" example:"10"`
	TotalPages int                   `json:"total_pages" example:"10"`
}
