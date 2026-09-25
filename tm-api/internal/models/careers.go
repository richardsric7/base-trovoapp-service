package models

import (
	"encoding/json"
	"time"
)

// CareerRoleStatus represents the status of a career role
const (
	CareerRoleStatusDraft     = "draft"
	CareerRoleStatusPublished = "published"
	CareerRoleStatusDisabled  = "disabled"
)

// Valid location types
const (
	LocationRemote = "remote"
	LocationHybrid = "hybrid"
	LocationOnsite = "onsite"
)

// Valid work types
const (
	WorkTypeFulltime = "fulltime"
	WorkTypeParttime = "parttime"
	WorkTypeContract = "contract"
)

// CareerRoleSections represents the sections of a career role posting
// Fields can be either string or []string
type CareerRoleSections map[string]interface{}

// CareerRoleApplication represents the application information
// Can be an array of strings or an object with location, apply, subject fields
type CareerRoleApplication map[string]interface{}

// CareerRole represents a career/job posting
type CareerRole struct {
	ID                string          `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Heading           string          `gorm:"column:heading;not null" json:"heading"`
	CompanyOverview   string          `gorm:"column:company_overview;type:text" json:"company_overview"`
	Location          string          `gorm:"column:location;not null;default:'onsite'" json:"location"`
	YearsOfExperience string          `gorm:"column:years_of_experience" json:"years_of_experience"`
	WorkType          string          `gorm:"column:work_type;not null;default:'fulltime'" json:"work_type"`
	Sections          json.RawMessage `gorm:"column:sections;type:jsonb;not null" json:"sections"`
	Application       json.RawMessage `gorm:"column:application;type:jsonb;not null" json:"application"`
	Status            string          `gorm:"column:status;type:career_role_status;default:'draft'" json:"status"`
	PublishedAt       *time.Time      `gorm:"column:published_at" json:"published_at,omitempty"`
	CreatedAt         time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"column:updated_at" json:"updated_at"`
	CreatedBy         *string         `gorm:"column:created_by;type:uuid" json:"created_by,omitempty"`
	UpdatedBy         *string         `gorm:"column:updated_by;type:uuid" json:"updated_by,omitempty"`
}

// TableName overrides the default gorm table name
func (CareerRole) TableName() string {
	return "career_roles"
}

// CareerRoleRequest represents the request payload for create/update operations
type CareerRoleRequest struct {
	Heading           string                `json:"heading" binding:"required" example:"Head of People Operations (HPO)"`
	CompanyOverview   string                `json:"company_overview" example:"Trovotech Ltd is a SEC-regulated digital asset platform..."`
	Location          string                `json:"location" binding:"required" example:"remote"`
	YearsOfExperience string                `json:"years_of_experience" example:"5+ years"`
	WorkType          string                `json:"work_type" binding:"required" example:"fulltime"`
	Sections          CareerRoleSections    `json:"sections" example:"{\"The Role\":[\"Role description\"],\"What You'll Own\":[\"Responsibilities\"]}"`
	Application       CareerRoleApplication `json:"application" example:"{\"location\":\"Nigeria\",\"apply\":\"email@example.com\",\"subject\":\"Job Title\"}"`
}

// CareerRoleResponse represents the response for career role operations
type CareerRoleResponse struct {
	ID                string                `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Heading           string                `json:"heading" example:"Head of People Operations (HPO)"`
	CompanyOverview   string                `json:"company_overview" example:"Trovotech Ltd is a SEC-regulated digital asset platform..."`
	Location          string                `json:"location" example:"remote" enums:"remote,hybrid,onsite"`
	YearsOfExperience string                `json:"years_of_experience" example:"5+ years"`
	WorkType          string                `json:"work_type" example:"fulltime" enums:"fulltime,parttime,contract"`
	Sections          CareerRoleSections    `json:"sections"`
	Application       CareerRoleApplication `json:"application"`
	Status            string                `json:"status" example:"published" enums:"draft,published,disabled"`
	PublishedAt       *time.Time            `json:"published_at,omitempty" example:"2024-01-15T10:00:00Z"`
	CreatedAt         time.Time             `json:"created_at" example:"2024-01-15T10:00:00Z"`
	UpdatedAt         time.Time             `json:"updated_at" example:"2024-01-15T10:00:00Z"`
}

// UpdateStatusRequest represents the request payload for status update operations
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=published disabled" example:"disabled" enums:"published,disabled"`
}

// CareerRoleFilters represents filter options for admin career role list
type CareerRoleFilters struct {
	Status            string `form:"status" example:"published" enums:"draft,published,disabled"`
	Location          string `form:"location" example:"remote" enums:"remote,hybrid,onsite"`
	WorkType          string `form:"work_type" example:"fulltime" enums:"fulltime,parttime,contract"`
	YearsOfExperience string `form:"years_of_experience" example:"5+ years"`
	Search            string `form:"search" example:"backend engineer"` // Search in heading and company_overview
}
