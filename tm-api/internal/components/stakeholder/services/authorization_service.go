package services

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"admin-panel-dashboard/internal/components/stakeholder/models"
	"admin-panel-dashboard/internal/trovosdk"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ServiceLinkAuthorizer interface {
	SendAuthorizationRequest(trovoUser, authDescription, deviceInfo, callbackUrl string, validityInMinutes int) (*trovosdk.TrovoWalletAuthorizationData, error)
	VerifyAuthorizationRequest(trovoUser, authID string) (*trovosdk.TrovoWalletAuthorizationRespomseData, error)
}

type AuthorizationService struct {
	db          *gorm.DB
	serviceLink ServiceLinkAuthorizer
	ttl         time.Duration
}

func NewAuthorizationService(db *gorm.DB, serviceLink ServiceLinkAuthorizer) *AuthorizationService {
	return &AuthorizationService{db: db, serviceLink: serviceLink, ttl: 10 * time.Minute}
}

func (s *AuthorizationService) CreateChallenge(ctx context.Context, auth AuthContext, req models.CreateAuthorizationRequest) (*models.AuthorizationChallengeResponse, error) {
	if !isValidChallengeTarget(req.Action, req.EntityType) {
		return nil, NewHTTPError(http.StatusBadRequest, "invalid action/entity_type combination")
	}
	if !auth.IsWalletLinked || auth.TrovoWalletUsername == nil || strings.TrimSpace(*auth.TrovoWalletUsername) == "" {
		return nil, NewHTTPError(http.StatusForbidden, "linked Trovo Wallet is required for this action")
	}
	cleanUsername := cleanTrovoUsername(*auth.TrovoWalletUsername)
	callbackURL := os.Getenv("LOGIN_CALLBACK_URL")
	if strings.TrimSpace(callbackURL) == "" {
		callbackURL = "http://localhost/stakeholder/authorization/callback"
	}
	description := fmt.Sprintf("Authorize %s for %s", req.Action, req.EntityType)
	authData, err := s.serviceLink.SendAuthorizationRequest(cleanUsername, description, "Trovo Manager Stakeholder Portal", callbackURL, int(s.ttl.Minutes()))
	if err != nil {
		return nil, NewHTTPError(http.StatusServiceUnavailable, "failed to create wallet authorization challenge")
	}
	now := time.Now()
	challenge := models.StakeholderAuthorizationChallenge{
		ID:              uuid.NewString(),
		MemberID:        auth.MemberID,
		OrganizationID:  auth.OrganizationID,
		StakeholderType: auth.StakeholderType,
		DashboardRole:   auth.DashboardRole,
		Action:          req.Action,
		EntityType:      req.EntityType,
		EntityID:        req.EntityID,
		AuthID:          authData.AuthID,
		Status:          models.AuthorizationStatusPending,
		ExpiresAt:       now.Add(s.ttl),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := s.db.WithContext(ctx).Create(&challenge).Error; err != nil {
		return nil, err
	}
	return &models.AuthorizationChallengeResponse{
		ID:          challenge.ID,
		AuthID:      challenge.AuthID,
		DynamicLink: authData.DynamicLink,
		QRCode:      authData.QRCode,
		Status:      challenge.Status,
		ExpiresAt:   challenge.ExpiresAt,
	}, nil
}

func (s *AuthorizationService) VerifyChallenge(ctx context.Context, auth AuthContext, challengeID string, authID string) (*models.AuthorizationChallengeResponse, error) {
	var challenge models.StakeholderAuthorizationChallenge
	err := s.db.WithContext(ctx).
		Where("id = ? AND member_id = ? AND organization_id = ?", challengeID, auth.MemberID, auth.OrganizationID).
		First(&challenge).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewHTTPError(http.StatusNotFound, "authorization challenge not found")
		}
		return nil, err
	}
	if challenge.Status != models.AuthorizationStatusPending {
		return nil, NewHTTPError(http.StatusConflict, "authorization challenge is not pending")
	}
	if time.Now().After(challenge.ExpiresAt) {
		now := time.Now()
		challenge.Status = models.AuthorizationStatusExpired
		challenge.UpdatedAt = now
		_ = s.db.WithContext(ctx).Save(&challenge).Error
		return nil, NewHTTPError(http.StatusForbidden, "authorization challenge has expired")
	}
	if strings.TrimSpace(authID) != "" && authID != challenge.AuthID {
		return nil, NewHTTPError(http.StatusBadRequest, "auth_id does not match challenge")
	}
	if !auth.IsWalletLinked || auth.TrovoWalletUsername == nil || strings.TrimSpace(*auth.TrovoWalletUsername) == "" {
		return nil, NewHTTPError(http.StatusForbidden, "linked Trovo Wallet is required for this action")
	}
	if _, err := s.serviceLink.VerifyAuthorizationRequest(cleanTrovoUsername(*auth.TrovoWalletUsername), challenge.AuthID); err != nil {
		now := time.Now()
		challenge.Status = models.AuthorizationStatusFailed
		challenge.UpdatedAt = now
		_ = s.db.WithContext(ctx).Save(&challenge).Error
		return nil, NewHTTPError(http.StatusForbidden, "wallet authorization verification failed")
	}
	now := time.Now()
	challenge.Status = models.AuthorizationStatusVerified
	challenge.VerifiedAt = &now
	challenge.UpdatedAt = now
	if err := s.db.WithContext(ctx).Save(&challenge).Error; err != nil {
		return nil, err
	}
	return &models.AuthorizationChallengeResponse{
		ID:        challenge.ID,
		AuthID:    challenge.AuthID,
		Status:    challenge.Status,
		ExpiresAt: challenge.ExpiresAt,
	}, nil
}

func (s *AuthorizationService) ConsumeVerifiedChallenge(ctx context.Context, tx *gorm.DB, auth AuthContext, action, entityType, entityID, challengeID string) error {
	if strings.TrimSpace(challengeID) == "" {
		return NewHTTPError(http.StatusForbidden, "authorization challenge is required")
	}
	if tx == nil {
		tx = s.db
	}
	var challenge models.StakeholderAuthorizationChallenge
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND member_id = ? AND organization_id = ? AND action = ? AND entity_type = ? AND entity_id = ?",
			challengeID, auth.MemberID, auth.OrganizationID, action, entityType, entityID).
		First(&challenge).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return NewHTTPError(http.StatusForbidden, "valid authorization challenge not found")
		}
		return err
	}
	if challenge.Status != models.AuthorizationStatusVerified {
		return NewHTTPError(http.StatusForbidden, "authorization challenge is not verified")
	}
	if time.Now().After(challenge.ExpiresAt) {
		challenge.Status = models.AuthorizationStatusExpired
		challenge.UpdatedAt = time.Now()
		_ = tx.WithContext(ctx).Save(&challenge).Error
		return NewHTTPError(http.StatusForbidden, "authorization challenge has expired")
	}
	now := time.Now()
	result := tx.WithContext(ctx).Model(&models.StakeholderAuthorizationChallenge{}).
		Where("id = ? AND status = ?", challenge.ID, models.AuthorizationStatusVerified).
		Updates(map[string]interface{}{
			"status":     models.AuthorizationStatusUsed,
			"used_at":    now,
			"updated_at": now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return NewHTTPError(http.StatusConflict, "authorization challenge was already used")
	}
	return nil
}

func isValidChallengeTarget(action, entityType string) bool {
	switch action {
	case models.ActionFundReleaseApprove:
		return entityType == models.EntityFundReleaseRequest
	case models.ActionFundReleaseExecute:
		return entityType == models.EntityFundReleaseRequest
	case models.ActionDistributionAuthorize:
		return entityType == models.EntityDistribution
	default:
		return false
	}
}

func cleanTrovoUsername(username string) string {
	username = strings.TrimSpace(strings.ToLower(username))
	return strings.Split(username, "@")[0]
}
