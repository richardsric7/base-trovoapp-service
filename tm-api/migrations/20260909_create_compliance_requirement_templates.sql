-- Templated KYC/KYB/compliance requirement system. See
-- PRD-kyc-compliance-templates.md. Schema is authoritatively applied via GORM
-- AutoMigrate (internal/db/main.go); this file documents it for parity with
-- the rest of backend/migrations/.

ALTER TABLE organizations
    ADD COLUMN IF NOT EXISTS level INTEGER;

CREATE TABLE IF NOT EXISTS compliance_requirement_templates (
    id TEXT PRIMARY KEY,
    org_type TEXT NOT NULL,
    level INTEGER NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (org_type, level)
);

CREATE TABLE IF NOT EXISTS compliance_requirement_template_items (
    id TEXT PRIMARY KEY,
    template_id TEXT NOT NULL REFERENCES compliance_requirement_templates(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    description TEXT,
    input_type TEXT NOT NULL,
    required BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_compliance_requirement_template_items_template_id
    ON compliance_requirement_template_items(template_id);

CREATE TABLE IF NOT EXISTS compliance_requirement_instances (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL REFERENCES organizations(id),
    category TEXT NOT NULL,
    requirement TEXT NOT NULL,
    description TEXT,
    required BOOLEAN NOT NULL DEFAULT true,
    template_item_id TEXT,
    template_version INTEGER,
    status TEXT NOT NULL DEFAULT 'pending',
    input_type TEXT NOT NULL DEFAULT 'text',
    submission_data JSONB,
    submitted_by_member_id TEXT,
    submitted_at TIMESTAMPTZ,
    rejection_reason TEXT,
    reviewed_by_admin_id TEXT,
    reviewed_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    completed_by_member_id TEXT,
    assigned_level_at_approval INTEGER,
    due_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_compliance_requirement_instances_org_id
    ON compliance_requirement_instances(org_id);
CREATE INDEX IF NOT EXISTS idx_compliance_requirement_instances_status
    ON compliance_requirement_instances(status);
CREATE INDEX IF NOT EXISTS idx_compliance_requirement_instances_template_item_id
    ON compliance_requirement_instances(template_item_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_compliance_requirement_org_template_item
    ON compliance_requirement_instances(org_id, template_item_id);
