package handlers

import (
	"log"
	"time"

	vaultsignermodels "admin-panel-dashboard/internal/components/vaultsigner/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// writeAuditLog persists one VaultSignerAuditLog row (Section 4). This is a
// historical record only, written after the fact — a failure here is logged
// but never blocks or rolls back the operation it's recording.
func writeAuditLog(db *gorm.DB, row vaultsignermodels.VaultSignerAuditLog) {
	row.ID = uuid.NewString()
	row.ChangedAt = time.Now()
	if err := db.Create(&row).Error; err != nil {
		log.Println("[vaultsigner] failed to write audit log row:", err)
	}
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func intPtr(i int) *int {
	return &i
}
