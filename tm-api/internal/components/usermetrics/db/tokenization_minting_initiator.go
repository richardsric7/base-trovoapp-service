package usermetrics

import (
	"admin-panel-dashboard/internal/models"

	"gorm.io/gorm"
)

// CreateMintingUser creates a new tokenization minting user (initiator or approver)
func CreateMintingUser(db *gorm.DB, username string, role string) (*models.MintingUserResponse, error) {
	var response *models.MintingUserResponse

	if role == "initiator" {
		newInitiator := &models.TokenizationMintingInitiator{
			Initiator: username,
		}
		if err := db.Create(newInitiator).Error; err != nil {
			return nil, err
		}
		response = &models.MintingUserResponse{
			ID:       newInitiator.ID,
			Username: newInitiator.Initiator,
			Role:     "initiator",
			//CreatedAt: newInitiator.CreatedAt,
			//UpdatedAt: newInitiator.UpdatedAt,
		}
	} else {
		newApprover := &models.TokenizationMintingApprover{
			Approver: username,
		}
		if err := db.Create(newApprover).Error; err != nil {
			return nil, err
		}
		response = &models.MintingUserResponse{
			ID:       newApprover.ID,
			Username: newApprover.Approver,
			Role:     "approver",
			//CreatedAt: newApprover.CreatedAt,
			//UpdatedAt: newApprover.UpdatedAt,
		}
	}

	return response, nil
}

// DeleteMintingUser deletes a tokenization minting user by ID and role
func DeleteMintingUser(db *gorm.DB, id int, role string) error {
	if role == "initiator" {
		return db.Delete(&models.TokenizationMintingInitiator{}, id).Error
	}
	return db.Delete(&models.TokenizationMintingApprover{}, id).Error
}

// GetAllMintingUsers retrieves all minting users with pagination and optional role filter
func GetAllMintingUsers(db *gorm.DB, role string, page, pageSize int) (*models.PaginatedMintingUserResponse, error) {
	var total int64
	var users []models.MintingUserResponse

	if role == "" || role == "initiator" {
		var initiators []models.TokenizationMintingInitiator
		if err := db.Model(&models.TokenizationMintingInitiator{}).Count(&total).Error; err != nil {
			return nil, err
		}
		offset := (page - 1) * pageSize
		if err := db.Offset(offset).Limit(pageSize).Find(&initiators).Error; err != nil {
			return nil, err
		}
		for _, initiator := range initiators {
			users = append(users, models.MintingUserResponse{
				ID:       initiator.ID,
				Username: initiator.Initiator,
				Role:     "initiator",
				//CreatedAt: initiator.CreatedAt,
				//UpdatedAt: initiator.UpdatedAt,
			})
		}
	}

	if role == "" || role == "approver" {
		var approvers []models.TokenizationMintingApprover
		var approverTotal int64
		if err := db.Model(&models.TokenizationMintingApprover{}).Count(&approverTotal).Error; err != nil {
			return nil, err
		}
		total += approverTotal
		offset := (page - 1) * pageSize
		if err := db.Offset(offset).Limit(pageSize).Find(&approvers).Error; err != nil {
			return nil, err
		}
		for _, approver := range approvers {
			users = append(users, models.MintingUserResponse{
				ID:       approver.ID,
				Username: approver.Approver,
				Role:     "approver",
				//CreatedAt: approver.CreatedAt,
				//UpdatedAt: approver.UpdatedAt,
			})
		}
	}

	response := &models.PaginatedMintingUserResponse{
		Data:       users,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return response, nil
}

// SearchMintingUsers searches minting users by query with pagination and optional role filter
func SearchMintingUsers(db *gorm.DB, query string, role string, page, pageSize int) (*models.PaginatedMintingUserResponse, error) {
	var total int64
	var users []models.MintingUserResponse

	if role == "" || role == "initiator" {
		var initiators []models.TokenizationMintingInitiator
		dbQuery := db.Model(&models.TokenizationMintingInitiator{})
		if query != "" {
			dbQuery = dbQuery.Where("initiator ILIKE ?", "%"+query+"%")
		}
		if err := dbQuery.Count(&total).Error; err != nil {
			return nil, err
		}
		offset := (page - 1) * pageSize
		if err := dbQuery.Offset(offset).Limit(pageSize).Find(&initiators).Error; err != nil {
			return nil, err
		}
		for _, initiator := range initiators {
			users = append(users, models.MintingUserResponse{
				ID:       initiator.ID,
				Username: initiator.Initiator,
				Role:     "initiator",
				//CreatedAt: initiator.CreatedAt,
				//UpdatedAt: initiator.UpdatedAt,
			})
		}
	}

	if role == "" || role == "approver" {
		var approvers []models.TokenizationMintingApprover
		var approverTotal int64
		dbQuery := db.Model(&models.TokenizationMintingApprover{})
		if query != "" {
			dbQuery = dbQuery.Where("approver ILIKE ?", "%"+query+"%")
		}
		if err := dbQuery.Count(&approverTotal).Error; err != nil {
			return nil, err
		}
		total += approverTotal
		offset := (page - 1) * pageSize
		if err := dbQuery.Offset(offset).Limit(pageSize).Find(&approvers).Error; err != nil {
			return nil, err
		}
		for _, approver := range approvers {
			users = append(users, models.MintingUserResponse{
				ID:       approver.ID,
				Username: approver.Approver,
				Role:     "approver",
				//CreatedAt: approver.CreatedAt,
				//UpdatedAt: approver.UpdatedAt,
			})
		}
	}

	response := &models.PaginatedMintingUserResponse{
		Data:       users,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return response, nil
}
