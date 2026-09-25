package careers

import (
	"admin-panel-dashboard/internal/models"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CreateCareerRole creates a new career role as draft
func CreateCareerRole(db *gorm.DB, req models.CareerRoleRequest, userID string) (string, error) {
	// Convert structs to JSON
	sectionsJSON, err := json.Marshal(req.Sections)
	if err != nil {
		return "", err
	}
	applicationJSON, err := json.Marshal(req.Application)
	if err != nil {
		return "", err
	}

	role := models.CareerRole{
		Heading:           req.Heading,
		CompanyOverview:   req.CompanyOverview,
		Location:          req.Location,
		YearsOfExperience: req.YearsOfExperience,
		WorkType:          req.WorkType,
		Sections:          sectionsJSON,
		Application:       applicationJSON,
		Status:            models.CareerRoleStatusDraft,
		CreatedBy:         &userID,
		UpdatedBy:         &userID,
	}

	if err := db.Create(&role).Error; err != nil {
		return "", err
	}

	return role.ID, nil
}

// GetPublishedCareerRoles retrieves only published career roles (public endpoint)
func GetPublishedCareerRoles(db *gorm.DB, search string) ([]models.CareerRoleResponse, error) {
	var roles []models.CareerRole
	query := db.Where("status = ?", models.CareerRoleStatusPublished)

	// Apply search filter if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("heading ILIKE ? OR company_overview ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Order("created_at DESC").Find(&roles).Error; err != nil {
		return nil, err
	}

	return convertRolesToResponse(roles), nil
}

// GetPublishedCareerRoleByID retrieves a published career role by ID (public endpoint)
func GetPublishedCareerRoleByID(db *gorm.DB, id string) (*models.CareerRoleResponse, error) {
	var role models.CareerRole
	query := db.Where("id = ? AND status = ?", id, models.CareerRoleStatusPublished)

	if err := query.First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("career role not found")
		}
		return nil, err
	}

	return convertRoleToResponse(&role), nil
}

// GetCareerRoleByID retrieves a career role by ID (admin endpoint - includes all statuses)
func GetCareerRoleByID(db *gorm.DB, id string, includeDrafts bool) (*models.CareerRoleResponse, error) {
	var role models.CareerRole
	query := db.Where("id = ?", id)

	// If not including drafts, only return published roles
	if !includeDrafts {
		query = query.Where("status = ?", models.CareerRoleStatusPublished)
	}

	if err := query.First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("career role not found")
		}
		return nil, err
	}

	return convertRoleToResponse(&role), nil
}

// GetAllCareerRolesWithFilters retrieves all career roles with optional filters (admin endpoint)
func GetAllCareerRolesWithFilters(db *gorm.DB, filters models.CareerRoleFilters) ([]models.CareerRoleResponse, error) {
	var roles []models.CareerRole
	query := db

	// Apply filters
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Location != "" {
		query = query.Where("location = ?", filters.Location)
	}
	if filters.WorkType != "" {
		query = query.Where("work_type = ?", filters.WorkType)
	}
	if filters.YearsOfExperience != "" {
		query = query.Where("years_of_experience = ?", filters.YearsOfExperience)
	}
	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		query = query.Where("heading ILIKE ? OR company_overview ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Order("created_at DESC").Find(&roles).Error; err != nil {
		return nil, err
	}

	return convertRolesToResponse(roles), nil
}

// convertRoleToResponse converts a CareerRole model to CareerRoleResponse
func convertRoleToResponse(role *models.CareerRole) *models.CareerRoleResponse {
	var sections models.CareerRoleSections
	var application models.CareerRoleApplication
	json.Unmarshal(role.Sections, &sections)
	json.Unmarshal(role.Application, &application)

	return &models.CareerRoleResponse{
		ID:                role.ID,
		Heading:           role.Heading,
		CompanyOverview:   role.CompanyOverview,
		Location:          role.Location,
		YearsOfExperience: role.YearsOfExperience,
		WorkType:          role.WorkType,
		Sections:          sections,
		Application:       application,
		Status:            role.Status,
		PublishedAt:       role.PublishedAt,
		CreatedAt:         role.CreatedAt,
		UpdatedAt:         role.UpdatedAt,
	}
}

// convertRolesToResponse converts a slice of CareerRole models to CareerRoleResponse slice
func convertRolesToResponse(roles []models.CareerRole) []models.CareerRoleResponse {
	responses := make([]models.CareerRoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = *convertRoleToResponse(&role)
	}
	return responses
}

// UpdateCareerRole updates an existing career role
func UpdateCareerRole(db *gorm.DB, id string, req models.CareerRoleRequest, userID string) error {
	var existing models.CareerRole
	if err := db.First(&existing, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("career role not found")
		}
		return err
	}

	// Convert structs to JSON
	sectionsJSON, err := json.Marshal(req.Sections)
	if err != nil {
		return err
	}
	applicationJSON, err := json.Marshal(req.Application)
	if err != nil {
		return err
	}

	updateData := models.CareerRole{
		Heading:           req.Heading,
		CompanyOverview:   req.CompanyOverview,
		Location:          req.Location,
		YearsOfExperience: req.YearsOfExperience,
		WorkType:          req.WorkType,
		Sections:          sectionsJSON,
		Application:       applicationJSON,
		UpdatedBy:         &userID,
	}

	return db.Model(&models.CareerRole{}).Where("id = ?", id).Updates(updateData).Error
}

// DeleteCareerRole deletes a career role by ID
func DeleteCareerRole(db *gorm.DB, id string) error {
	var existing models.CareerRole
	if err := db.First(&existing, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("career role not found")
		}
		return err
	}

	return db.Delete(&existing).Error
}

// PublishCareerRole publishes a draft role
func PublishCareerRole(db *gorm.DB, id string) error {
	var existing models.CareerRole
	if err := db.First(&existing, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("career role not found")
		}
		return err
	}

	if existing.Status != models.CareerRoleStatusDraft {
		return errors.New("only draft roles can be published")
	}

	now := time.Now()
	return db.Model(&models.CareerRole{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       models.CareerRoleStatusPublished,
		"published_at": &now,
	}).Error
}

// UpdateCareerRoleStatus updates the status of a career role (published ↔ disabled)
func UpdateCareerRoleStatus(db *gorm.DB, id string, status string) error {
	var existing models.CareerRole
	if err := db.First(&existing, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("career role not found")
		}
		return err
	}

	// Validate status transition
	if existing.Status == models.CareerRoleStatusDraft {
		return errors.New("draft roles must be published first")
	}

	if status != models.CareerRoleStatusPublished && status != models.CareerRoleStatusDisabled {
		return errors.New("invalid status: must be 'published' or 'disabled'")
	}

	// If publishing a disabled role, ensure published_at is set
	updates := map[string]interface{}{"status": status}
	if status == models.CareerRoleStatusPublished && existing.PublishedAt == nil {
		now := time.Now()
		updates["published_at"] = &now
	}

	return db.Model(&models.CareerRole{}).Where("id = ?", id).Updates(updates).Error
}
