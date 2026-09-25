package db

import (
	stakeholderModels "admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"
	"fmt"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAdminAutoMigrateCreatesGORMOwnedStakeholderSchema(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatal(err)
	}
	models := []interface{}{
		&coreModels.Organization{},
		&coreModels.OrganizationMember{},
		&stakeholderModels.FundReleaseRequest{},
		&stakeholderModels.ComplianceItem{},
		&stakeholderModels.DueDiligenceChecklist{},
		&stakeholderModels.DueDiligenceItem{},
		&stakeholderModels.StakeholderDocument{},
		&stakeholderModels.StakeholderReport{},
		&stakeholderModels.ComplianceRequirementTemplate{},
		&stakeholderModels.ComplianceRequirementInstance{},
		&stakeholderModels.ComplianceItemDocument{},
		&stakeholderModels.DueDiligenceItemDocument{},
		&stakeholderModels.ComplianceRequirementTemplateItem{},
	}
	if err := database.AutoMigrate(models...); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
	if err := database.AutoMigrate(models...); err != nil {
		t.Fatalf("second AutoMigrate run failed: %v", err)
	}

	for _, model := range []interface{}{
		&stakeholderModels.ComplianceItemDocument{},
		&stakeholderModels.DueDiligenceItemDocument{},
		&stakeholderModels.ComplianceRequirementTemplate{},
		&stakeholderModels.ComplianceRequirementTemplateItem{},
		&stakeholderModels.ComplianceRequirementInstance{},
	} {
		if !database.Migrator().HasTable(model) {
			t.Fatalf("missing table for %T", model)
		}
	}
	if !database.Migrator().HasIndex(&stakeholderModels.ComplianceItemDocument{}, "idx_compliance_item_documents_document_id") {
		t.Fatal("missing compliance document lookup index")
	}
	if !database.Migrator().HasIndex(&stakeholderModels.DueDiligenceItemDocument{}, "idx_due_diligence_item_documents_document_id") {
		t.Fatal("missing due-diligence document lookup index")
	}
	if !database.Migrator().HasIndex(&stakeholderModels.ComplianceRequirementInstance{}, "idx_compliance_requirement_org_template_item") {
		t.Fatal("missing organization/template-item uniqueness index")
	}

	assertRequiredDefaultColumn(t, database, &stakeholderModels.FundReleaseRequest{}, "receiving_bank")
	assertRequiredDefaultColumn(t, database, &stakeholderModels.FundReleaseRequest{}, "receiving_account_name")
	assertRequiredDefaultColumn(t, database, &stakeholderModels.FundReleaseRequest{}, "receiving_account_number")
	assertRequiredDefaultColumn(t, database, &coreModels.OrganizationMember{}, "session_version")
	assertRequiredDefaultColumn(t, database, &stakeholderModels.ComplianceRequirementTemplateItem{}, "required")
	assertRequiredDefaultColumn(t, database, &stakeholderModels.ComplianceRequirementInstance{}, "required")
	if !database.Migrator().HasColumn(&coreModels.Organization{}, "level") {
		t.Fatal("organizations.level was not created by AutoMigrate")
	}
	if !database.Migrator().HasColumn(&stakeholderModels.StakeholderReport{}, "document_id") {
		t.Fatal("stakeholder_reports.document_id was not created by AutoMigrate")
	}

	compliance := stakeholderModels.ComplianceItem{ID: "compliance-1", OrgID: "custodian-1", Category: "regulatory", Requirement: "Evidence"}
	checklist := stakeholderModels.DueDiligenceChecklist{ID: "checklist-1", AssetID: "asset-1", TrusteeOrgID: "trustee-1"}
	dueItem := stakeholderModels.DueDiligenceItem{ID: "due-item-1", ChecklistID: checklist.ID, Category: "legal", Item: "Evidence"}
	document := stakeholderModels.StakeholderDocument{
		ID: "document-1", UploadedByOrgID: "custodian-1", UploadedByMemberID: "member-1",
		Category: "compliance", Title: "Evidence", FileURL: "https://example.invalid/document-1",
	}
	for _, model := range []interface{}{&compliance, &checklist, &dueItem, &document} {
		if err := database.Create(model).Error; err != nil {
			t.Fatalf("seed %T: %v", model, err)
		}
	}
	if err := database.Create(&stakeholderModels.ComplianceItemDocument{ComplianceItemID: compliance.ID, DocumentID: document.ID}).Error; err != nil {
		t.Fatalf("create compliance document link: %v", err)
	}
	if err := database.Create(&stakeholderModels.DueDiligenceItemDocument{DueDiligenceItemID: dueItem.ID, DocumentID: document.ID}).Error; err != nil {
		t.Fatalf("create due-diligence document link: %v", err)
	}
	if err := database.Delete(&document).Error; err != nil {
		t.Fatalf("delete linked document: %v", err)
	}
	assertLinkCount(t, database, &stakeholderModels.ComplianceItemDocument{}, 0)
	assertLinkCount(t, database, &stakeholderModels.DueDiligenceItemDocument{}, 0)

	document.ID = "document-2"
	document.FileURL = "https://example.invalid/document-2"
	if err := database.Create(&document).Error; err != nil {
		t.Fatalf("seed second document: %v", err)
	}
	if err := database.Create(&stakeholderModels.ComplianceItemDocument{ComplianceItemID: compliance.ID, DocumentID: document.ID}).Error; err != nil {
		t.Fatalf("create second compliance document link: %v", err)
	}
	if err := database.Create(&stakeholderModels.DueDiligenceItemDocument{DueDiligenceItemID: dueItem.ID, DocumentID: document.ID}).Error; err != nil {
		t.Fatalf("create second due-diligence document link: %v", err)
	}
	if err := database.Delete(&compliance).Error; err != nil {
		t.Fatalf("delete linked compliance item: %v", err)
	}
	if err := database.Delete(&dueItem).Error; err != nil {
		t.Fatalf("delete linked due-diligence item: %v", err)
	}
	assertLinkCount(t, database, &stakeholderModels.ComplianceItemDocument{}, 0)
	assertLinkCount(t, database, &stakeholderModels.DueDiligenceItemDocument{}, 0)
}

func assertRequiredDefaultColumn(t *testing.T, database *gorm.DB, model interface{}, name string) {
	t.Helper()
	columns, err := database.Migrator().ColumnTypes(model)
	if err != nil {
		t.Fatal(err)
	}
	for _, column := range columns {
		if column.Name() != name {
			continue
		}
		if nullable, ok := column.Nullable(); !ok || nullable {
			t.Fatalf("column %s nullable = %t, known = %t", name, nullable, ok)
		}
		if _, ok := column.DefaultValue(); !ok {
			t.Fatalf("column %s has no default", name)
		}
		return
	}
	t.Fatalf("column %s not found", name)
}

func assertLinkCount(t *testing.T, database *gorm.DB, model interface{}, want int64) {
	t.Helper()
	var count int64
	if err := database.Model(model).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != want {
		t.Fatalf("%T row count = %d, want %d", model, count, want)
	}
}

func TestAdminAutoMigrateCreatesStakeholderPortalTables(t *testing.T) {
	if os.Getenv("ADMIN_MIGRATION_INTEGRATION") != "1" {
		t.Skip("set ADMIN_MIGRATION_INTEGRATION=1 to run against a local PostgreSQL database")
	}

	dsn := os.Getenv("ADMIN_CONNECTION_STRING")
	if dsn == "" {
		t.Fatal("ADMIN_CONNECTION_STRING is required")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	schema := fmt.Sprintf("stakeholder_migration_test_%d", time.Now().UnixNano())
	tx := database.Begin()
	if tx.Error != nil {
		t.Fatal(tx.Error)
	}
	defer tx.Rollback()

	if err := tx.Exec(`CREATE SCHEMA "` + schema + `"`).Error; err != nil {
		t.Fatal(err)
	}
	if err := tx.Exec(`SET LOCAL search_path TO "` + schema + `"`).Error; err != nil {
		t.Fatal(err)
	}
	if err := migrateAdminSchema(tx); err != nil {
		t.Fatalf("migrateAdminSchema failed: %v", err)
	}
	if err := tx.Exec(`
		ALTER TABLE fund_release_requests
			ALTER COLUMN receiving_bank DROP NOT NULL,
			ALTER COLUMN receiving_bank DROP DEFAULT,
			ALTER COLUMN receiving_account_name DROP NOT NULL,
			ALTER COLUMN receiving_account_name DROP DEFAULT,
			ALTER COLUMN receiving_account_number DROP NOT NULL,
			ALTER COLUMN receiving_account_number DROP DEFAULT;
		INSERT INTO fund_release_requests (
			id, asset_id, requester_org_id, requester_member_id, trustee_org_id,
			amount, currency, purpose, status,
			receiving_bank, receiving_account_name, receiving_account_number
		) VALUES (
			'legacy-null-receiving-account', 'asset-1', 'manager-1', 'member-1', 'trustee-1',
			100, 'CNGN', 'Legacy request', 'submitted', NULL, NULL, NULL
		);

		ALTER TABLE compliance_item_documents
			DROP CONSTRAINT compliance_item_documents_compliance_item_id_fkey,
			DROP CONSTRAINT compliance_item_documents_document_id_fkey;
		ALTER TABLE due_diligence_item_documents
			DROP CONSTRAINT due_diligence_item_documents_due_diligence_item_id_fkey,
			DROP CONSTRAINT due_diligence_item_documents_document_id_fkey;
		DROP INDEX idx_compliance_item_documents_document_id;
		DROP INDEX idx_due_diligence_item_documents_document_id;

		INSERT INTO compliance_item_documents (compliance_item_id, document_id)
		VALUES ('missing-compliance', 'missing-document');
		INSERT INTO due_diligence_item_documents (due_diligence_item_id, document_id)
		VALUES ('missing-due-diligence-item', 'missing-document')
	`).Error; err != nil {
		t.Fatalf("prepare legacy partially migrated schema: %v", err)
	}
	if err := migrateAdminSchema(tx); err != nil {
		t.Fatalf("second migrateAdminSchema run failed: %v", err)
	}

	var tableCount int64
	if err := tx.Raw(`
		SELECT count(*)
		FROM information_schema.tables
		WHERE table_schema = ?
		  AND table_name IN (
			'stakeholder_asset_assignments',
			'stakeholder_audit_logs',
			'stakeholder_notifications',
			'stakeholder_authorization_challenges',
			'fund_release_requests',
			'due_diligence_checklists',
			'due_diligence_items',
			'revenue_records',
			'distributions',
			'asset_valuations',
			'segregated_accounts',
			'compliance_items',
			'stakeholder_documents',
			'compliance_item_documents',
			'due_diligence_item_documents',
			'compliance_requirement_templates',
			'compliance_requirement_template_items',
			'compliance_requirement_instances',
			'stakeholder_notification_preferences',
			'stakeholder_asset_operations',
			'stakeholder_reports',
			'stakeholder_structuring_statuses'
		  )
	`, schema).Scan(&tableCount).Error; err != nil {
		t.Fatal(err)
	}
	if tableCount != 22 {
		t.Fatalf("created %d stakeholder portal tables, want 22", tableCount)
	}

	type schemaColumn struct {
		ColumnName    string
		IsNullable    string
		ColumnDefault *string
	}
	var receivingColumns []schemaColumn
	if err := tx.Raw(`
		SELECT column_name, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_schema = ?
		  AND table_name = 'fund_release_requests'
		  AND column_name IN ('receiving_bank', 'receiving_account_name', 'receiving_account_number')
		  AND data_type = 'text'
		ORDER BY column_name
	`, schema).Scan(&receivingColumns).Error; err != nil {
		t.Fatal(err)
	}
	if len(receivingColumns) != 3 {
		t.Fatalf("created %d receiving-account columns, want 3", len(receivingColumns))
	}
	for _, column := range receivingColumns {
		if column.IsNullable != "NO" || column.ColumnDefault == nil {
			t.Fatalf("receiving-account column %s nullable=%s default=%v", column.ColumnName, column.IsNullable, column.ColumnDefault)
		}
	}

	var nullReceivingAccountCount int64
	if err := tx.Raw(`
		SELECT count(*)
		FROM fund_release_requests
		WHERE id = 'legacy-null-receiving-account'
		  AND (receiving_bank IS NULL OR receiving_account_name IS NULL OR receiving_account_number IS NULL)
	`).Scan(&nullReceivingAccountCount).Error; err != nil {
		t.Fatal(err)
	}
	if nullReceivingAccountCount != 0 {
		t.Fatal("legacy receiving-account NULL values were not backfilled")
	}

	var orphanedEvidenceLinkCount int64
	if err := tx.Raw(`
		SELECT
			(SELECT count(*) FROM compliance_item_documents WHERE compliance_item_id = 'missing-compliance') +
			(SELECT count(*) FROM due_diligence_item_documents WHERE due_diligence_item_id = 'missing-due-diligence-item')
	`).Scan(&orphanedEvidenceLinkCount).Error; err != nil {
		t.Fatal(err)
	}
	if orphanedEvidenceLinkCount != 0 {
		t.Fatalf("retained %d orphaned evidence links, want 0", orphanedEvidenceLinkCount)
	}

	var cascadeConstraintCount int64
	if err := tx.Raw(`
		SELECT count(*)
		FROM information_schema.table_constraints tc
		JOIN information_schema.referential_constraints rc
		  ON rc.constraint_schema = tc.constraint_schema
		 AND rc.constraint_name = tc.constraint_name
		WHERE tc.constraint_schema = ?
		  AND tc.table_name IN ('compliance_item_documents', 'due_diligence_item_documents')
		  AND tc.constraint_name IN (
			'compliance_item_documents_compliance_item_id_fkey',
			'compliance_item_documents_document_id_fkey',
			'due_diligence_item_documents_due_diligence_item_id_fkey',
			'due_diligence_item_documents_document_id_fkey'
		  )
		  AND tc.constraint_type = 'FOREIGN KEY'
		  AND rc.delete_rule = 'CASCADE'
		  AND rc.update_rule = 'NO ACTION'
	`, schema).Scan(&cascadeConstraintCount).Error; err != nil {
		t.Fatal(err)
	}
	if cascadeConstraintCount != 4 {
		t.Fatalf("created %d cascading evidence foreign keys, want 4", cascadeConstraintCount)
	}

	var documentIndexCount int64
	if err := tx.Raw(`
		SELECT count(*)
		FROM pg_indexes
		WHERE schemaname = ?
		  AND indexname IN (
			'idx_compliance_item_documents_document_id',
			'idx_due_diligence_item_documents_document_id'
		  )
	`, schema).Scan(&documentIndexCount).Error; err != nil {
		t.Fatal(err)
	}
	if documentIndexCount != 2 {
		t.Fatalf("created %d evidence document indexes, want 2", documentIndexCount)
	}

	var evidencePrimaryKeyColumnCount int64
	if err := tx.Raw(`
		SELECT count(*)
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
		  ON kcu.constraint_schema = tc.constraint_schema
		 AND kcu.constraint_name = tc.constraint_name
		WHERE tc.constraint_schema = ?
		  AND tc.table_name IN ('compliance_item_documents', 'due_diligence_item_documents')
		  AND tc.constraint_type = 'PRIMARY KEY'
	`, schema).Scan(&evidencePrimaryKeyColumnCount).Error; err != nil {
		t.Fatal(err)
	}
	if evidencePrimaryKeyColumnCount != 4 {
		t.Fatalf("created %d evidence primary-key columns, want 4", evidencePrimaryKeyColumnCount)
	}

	var evidenceTimestampCount int64
	if err := tx.Raw(`
		SELECT count(*)
		FROM information_schema.columns
		WHERE table_schema = ?
		  AND table_name IN ('compliance_item_documents', 'due_diligence_item_documents')
		  AND column_name = 'created_at'
		  AND is_nullable = 'NO'
		  AND column_default IS NOT NULL
	`, schema).Scan(&evidenceTimestampCount).Error; err != nil {
		t.Fatal(err)
	}
	if evidenceTimestampCount != 2 {
		t.Fatalf("created %d required evidence timestamps with defaults, want 2", evidenceTimestampCount)
	}

	var sessionVersionColumnCount int64
	if err := tx.Raw(`
		SELECT count(*)
		FROM information_schema.columns
		WHERE table_schema = ?
		  AND table_name = 'organization_members'
		  AND column_name = 'session_version'
		  AND data_type = 'bigint'
		  AND is_nullable = 'NO'
		  AND column_default IS NOT NULL
	`, schema).Scan(&sessionVersionColumnCount).Error; err != nil {
		t.Fatal(err)
	}
	if sessionVersionColumnCount != 1 {
		t.Fatalf("created %d required BIGINT session-version columns with defaults, want 1", sessionVersionColumnCount)
	}
}
