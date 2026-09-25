package services

import (
	"context"
	"time"

	"admin-panel-dashboard/internal/components/stakeholder/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditService struct {
	db *gorm.DB
}

type AuditEvent struct {
	RequestID     string
	ActorMemberID string
	ActorOrgID    string
	ActorRole     string
	Action        string
	EntityType    string
	EntityID      string
	BeforeState   models.JSONMap
	AfterState    models.JSONMap
	Metadata      models.JSONMap
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{db: db}
}

func (s *AuditService) Record(ctx context.Context, tx *gorm.DB, event AuditEvent) error {
	if tx == nil {
		tx = s.db
	}
	if event.BeforeState == nil {
		event.BeforeState = models.JSONMap{}
	}
	if event.AfterState == nil {
		event.AfterState = models.JSONMap{}
	}
	if event.Metadata == nil {
		event.Metadata = models.JSONMap{}
	}
	log := models.StakeholderAuditLog{
		ID:            uuid.NewString(),
		RequestID:     event.RequestID,
		ActorMemberID: event.ActorMemberID,
		ActorOrgID:    event.ActorOrgID,
		ActorRole:     event.ActorRole,
		Action:        event.Action,
		EntityType:    event.EntityType,
		EntityID:      event.EntityID,
		BeforeState:   event.BeforeState,
		AfterState:    event.AfterState,
		Metadata:      event.Metadata,
		CreatedAt:     time.Now(),
	}
	return tx.WithContext(ctx).Create(&log).Error
}

func AuditEventFromAuth(auth AuthContext, action, entityType, entityID string) AuditEvent {
	actorRole := auth.DashboardRole
	if actorRole == "" {
		actorRole = auth.MemberRole
	}
	return AuditEvent{
		ActorMemberID: auth.MemberID,
		ActorOrgID:    auth.OrganizationID,
		ActorRole:     actorRole,
		Action:        action,
		EntityType:    entityType,
		EntityID:      entityID,
	}
}
