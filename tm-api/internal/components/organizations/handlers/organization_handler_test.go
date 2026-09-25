package handlers

import (
	"context"
	"testing"
	"time"

	"admin-panel-dashboard/internal/models"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSetupPasswordActivatesOrganizationForRootAdmin(t *testing.T) {
	db := newOrganizationTestDB(t)
	invite, member, organization := seedPasswordSetup(t, db, string(models.OrganizationSuperAdmin))

	if err := completePasswordSetup(db, &invite, &member, "hashed-password", time.Now()); err != nil {
		t.Fatalf("completePasswordSetup returned an error: %v", err)
	}

	if err := db.First(&organization, "id = ?", organization.ID).Error; err != nil {
		t.Fatal(err)
	}
	if organization.Status != models.OrganizationStatusActive {
		t.Fatalf("organization status = %q, want %q", organization.Status, models.OrganizationStatusActive)
	}

	if err := db.First(&member, "id = ?", member.ID).Error; err != nil {
		t.Fatal(err)
	}
	if member.Status != models.OrganizationStatusActive {
		t.Fatalf("member status = %q, want %q", member.Status, models.OrganizationStatusActive)
	}
}

func TestRevokeOrganizationMemberSessionsIncrementsVersion(t *testing.T) {
	database := newOrganizationTestDB(t)
	_, member, _ := seedPasswordSetup(t, database, string(models.OrganizationSuperAdmin))

	if err := revokeOrganizationMemberSessions(context.Background(), database, member.ID, member.OrganizationID); err != nil {
		t.Fatalf("revoke sessions: %v", err)
	}
	var updated models.OrganizationMember
	if err := database.First(&updated, "id = ?", member.ID).Error; err != nil {
		t.Fatalf("reload member: %v", err)
	}
	if updated.SessionVersion != 1 {
		t.Fatalf("session version = %d, want 1", updated.SessionVersion)
	}
	if err := revokeOrganizationMemberSessions(context.Background(), database, member.ID, "wrong-org"); err != gorm.ErrRecordNotFound {
		t.Fatalf("wrong organization error = %v, want record not found", err)
	}
}

func TestSetupPasswordDoesNotActivateOrganizationForRegularMember(t *testing.T) {
	db := newOrganizationTestDB(t)
	invite, member, organization := seedPasswordSetup(t, db, string(models.OrganizationMemberRole))

	if err := completePasswordSetup(db, &invite, &member, "hashed-password", time.Now()); err != nil {
		t.Fatalf("completePasswordSetup returned an error: %v", err)
	}

	if err := db.First(&organization, "id = ?", organization.ID).Error; err != nil {
		t.Fatal(err)
	}
	if organization.Status != models.OrganizationStatusPending {
		t.Fatalf("organization status = %q, want %q", organization.Status, models.OrganizationStatusPending)
	}
}

func newOrganizationTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+uuid.NewString()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Organization{}, &models.OrganizationInvite{}, &models.OrganizationMember{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func seedPasswordSetup(t *testing.T, db *gorm.DB, role string) (models.OrganizationInvite, models.OrganizationMember, models.Organization) {
	t.Helper()

	now := time.Now()
	organization := models.Organization{
		ID:        uuid.NewString(),
		Name:      "Test Organization",
		Email:     uuid.NewString() + "@example.com",
		Type:      models.OrganizationTypeAssetManager,
		Status:    models.OrganizationStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: "test",
	}
	invite := models.OrganizationInvite{
		ID:              uuid.NewString(),
		OrganizationID:  organization.ID,
		Email:           uuid.NewString() + "@example.com",
		Role:            role,
		Status:          models.InviteStatusPending,
		ExpiresAt:       now.Add(time.Hour),
		CreatedAt:       now,
		UpdatedAt:       now,
		InvitedBy:       "test",
		VerificationOTP: "123456",
		OTPExpiresAt:    now.Add(time.Hour),
	}
	member := models.OrganizationMember{
		ID:             uuid.NewString(),
		OrganizationID: organization.ID,
		Email:          invite.Email,
		Role:           role,
		Status:         models.OrganizationStatusPending,
		CreatedAt:      now,
		UpdatedAt:      now,
		CreatedBy:      "test",
		InviteID:       invite.ID,
	}

	if err := db.Create(&organization).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&invite).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}
	return invite, member, organization
}
