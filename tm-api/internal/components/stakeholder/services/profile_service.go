package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileService struct {
	db    *gorm.DB
	audit *AuditService
}

func NewProfileService(db *gorm.DB, audit *AuditService) *ProfileService {
	return &ProfileService{db: db, audit: audit}
}

func (s *ProfileService) GetProfile(ctx context.Context, auth AuthContext) (*models.StakeholderProfileResponse, error) {
	var member coreModels.OrganizationMember
	if err := s.db.WithContext(ctx).First(&member, "id = ?", auth.MemberID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewHTTPError(http.StatusNotFound, "member not found")
		}
		return nil, err
	}
	var org coreModels.Organization
	if err := s.db.WithContext(ctx).First(&org, "id = ?", auth.OrganizationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewHTTPError(http.StatusNotFound, "organization not found")
		}
		return nil, err
	}
	stakeholderType := ""
	if org.StakeholderType != nil {
		stakeholderType = *org.StakeholderType
	}
	role, _ := models.DashboardRoleForStakeholderType(stakeholderType)
	return &models.StakeholderProfileResponse{
		Member: models.StakeholderProfileMember{
			ID:        member.ID,
			Email:     member.Email,
			Role:      member.Role,
			Status:    member.Status,
			FirstName: member.FirstName,
			LastName:  member.LastName,
		},
		Organization: models.StakeholderProfileOrganization{
			ID:     org.ID,
			Name:   org.Name,
			Email:  org.Email,
			Type:   org.Type,
			Status: org.Status,
		},
		Stakeholder: models.StakeholderProfileLinkage{
			ID:            org.StakeholderID,
			Type:          stakeholderType,
			DashboardRole: role,
		},
		Wallet: models.StakeholderProfileWallet{
			TrovoWalletUsername: member.TrovoWalletUsername,
			IsWalletLinked:      member.IsWalletLinked,
		},
	}, nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, auth AuthContext, req models.UpdateProfileRequest) (*models.StakeholderProfileResponse, error) {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{"updated_at": time.Now()}
		if req.FirstName != "" {
			updates["first_name"] = req.FirstName
		}
		if req.LastName != "" {
			updates["last_name"] = req.LastName
		}
		if len(updates) == 1 {
			return NewHTTPError(http.StatusBadRequest, "no profile fields to update")
		}
		if err := tx.Model(&coreModels.OrganizationMember{}).Where("id = ? AND organization_id = ?", auth.MemberID, auth.OrganizationID).Updates(updates).Error; err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionProfileUpdate, models.EntityProfile, auth.MemberID)
			event.AfterState = models.JSONMap(updates)
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.GetProfile(ctx, auth)
}

func (s *ProfileService) UpdateNotificationPreferences(ctx context.Context, auth AuthContext, req models.UpdateNotificationPreferencesRequest) (*models.StakeholderNotificationPreference, error) {
	var preference models.StakeholderNotificationPreference
	err := s.db.WithContext(ctx).Where("member_id = ?", auth.MemberID).First(&preference).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	now := time.Now()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		preference = models.StakeholderNotificationPreference{
			ID:             uuid.NewString(),
			MemberID:       auth.MemberID,
			OrganizationID: auth.OrganizationID,
			EmailEnabled:   true,
			InAppEnabled:   true,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
	}
	if req.EmailEnabled != nil {
		preference.EmailEnabled = *req.EmailEnabled
	}
	if req.InAppEnabled != nil {
		preference.InAppEnabled = *req.InAppEnabled
	}
	preference.UpdatedAt = now
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&preference).Error; err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionNotificationPreferencesUpdate, models.EntityProfile, auth.MemberID)
			event.AfterState = models.JSONMap{
				"email_enabled":  preference.EmailEnabled,
				"in_app_enabled": preference.InAppEnabled,
			}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &preference, nil
}
