package services

import (
	"context"
	"net/http"
	"time"

	"admin-panel-dashboard/internal/components/stakeholder/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

type NotificationInput struct {
	RecipientOrgID    string
	SenderOrgID       *string
	Type              string
	Title             string
	Message           string
	RelatedEntityType string
	RelatedEntityID   string
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

func (s *NotificationService) Create(ctx context.Context, tx *gorm.DB, input NotificationInput) error {
	if input.RecipientOrgID == "" {
		return nil
	}
	if tx == nil {
		tx = s.db
	}
	notification := models.StakeholderNotification{
		ID:                uuid.NewString(),
		RecipientOrgID:    input.RecipientOrgID,
		SenderOrgID:       input.SenderOrgID,
		Type:              input.Type,
		Title:             input.Title,
		Message:           input.Message,
		RelatedEntityType: input.RelatedEntityType,
		RelatedEntityID:   input.RelatedEntityID,
		CreatedAt:         time.Now(),
	}
	return tx.WithContext(ctx).Create(&notification).Error
}

func (s *NotificationService) List(ctx context.Context, auth AuthContext, page, limit int) ([]models.StakeholderNotification, int64, error) {
	query := s.db.WithContext(ctx).Model(&models.StakeholderNotification{}).
		Where("recipient_org_id = ?", auth.OrganizationID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var notifications []models.StakeholderNotification
	if err := query.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&notifications).Error; err != nil {
		return nil, 0, err
	}
	return notifications, total, nil
}

func (s *NotificationService) MarkRead(ctx context.Context, auth AuthContext, id string) (*models.StakeholderNotification, error) {
	var notification models.StakeholderNotification
	err := s.db.WithContext(ctx).Where("id = ? AND recipient_org_id = ?", id, auth.OrganizationID).First(&notification).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewHTTPError(http.StatusNotFound, "notification not found")
		}
		return nil, err
	}
	now := time.Now()
	if notification.ReadAt == nil {
		notification.ReadAt = &now
		if err := s.db.WithContext(ctx).Save(&notification).Error; err != nil {
			return nil, err
		}
	}
	return &notification, nil
}
