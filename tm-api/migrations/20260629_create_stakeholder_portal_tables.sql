-- Stakeholder Portal workflow tables.
-- AdminDB owns these tables. Do not apply to production without the deployment checklist.

CREATE TABLE IF NOT EXISTS stakeholder_asset_assignments (
    id TEXT PRIMARY KEY,
    asset_id TEXT NOT NULL UNIQUE,
    asset_code TEXT,
    trustee_org_id TEXT REFERENCES organizations(id),
    trustee_stakeholder_id BIGINT,
    asset_manager_org_id TEXT REFERENCES organizations(id),
    asset_manager_stakeholder_id BIGINT,
    custodian_org_id TEXT REFERENCES organizations(id),
    custodian_stakeholder_id BIGINT,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    assigned_by TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_stakeholder_asset_assignments_asset_code ON stakeholder_asset_assignments(asset_code);
CREATE INDEX IF NOT EXISTS idx_stakeholder_asset_assignments_trustee_org_id ON stakeholder_asset_assignments(trustee_org_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_asset_assignments_trustee_stakeholder_id ON stakeholder_asset_assignments(trustee_stakeholder_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_asset_assignments_asset_manager_org_id ON stakeholder_asset_assignments(asset_manager_org_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_asset_assignments_asset_manager_stakeholder_id ON stakeholder_asset_assignments(asset_manager_stakeholder_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_asset_assignments_custodian_org_id ON stakeholder_asset_assignments(custodian_org_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_asset_assignments_custodian_stakeholder_id ON stakeholder_asset_assignments(custodian_stakeholder_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_asset_assignments_status ON stakeholder_asset_assignments(status);

CREATE TABLE IF NOT EXISTS stakeholder_audit_logs (
    id TEXT PRIMARY KEY,
    request_id TEXT,
    actor_member_id TEXT,
    actor_org_id TEXT,
    actor_role TEXT,
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    before_state JSONB NOT NULL DEFAULT '{}'::jsonb,
    after_state JSONB NOT NULL DEFAULT '{}'::jsonb,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_stakeholder_audit_logs_actor_org_id ON stakeholder_audit_logs(actor_org_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_audit_logs_entity ON stakeholder_audit_logs(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_audit_logs_action ON stakeholder_audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_stakeholder_audit_logs_created_at ON stakeholder_audit_logs(created_at);

CREATE TABLE IF NOT EXISTS stakeholder_notifications (
    id TEXT PRIMARY KEY,
    recipient_org_id TEXT NOT NULL REFERENCES organizations(id),
    sender_org_id TEXT REFERENCES organizations(id),
    type TEXT NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    related_entity_type TEXT,
    related_entity_id TEXT,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_stakeholder_notifications_recipient_org_id ON stakeholder_notifications(recipient_org_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_notifications_sender_org_id ON stakeholder_notifications(sender_org_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_notifications_type ON stakeholder_notifications(type);
CREATE INDEX IF NOT EXISTS idx_stakeholder_notifications_related_entity ON stakeholder_notifications(related_entity_type, related_entity_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_notifications_read_at ON stakeholder_notifications(read_at);
CREATE INDEX IF NOT EXISTS idx_stakeholder_notifications_created_at ON stakeholder_notifications(created_at);

CREATE TABLE IF NOT EXISTS stakeholder_authorization_challenges (
    id TEXT PRIMARY KEY,
    member_id TEXT NOT NULL REFERENCES organization_members(id),
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    stakeholder_type TEXT NOT NULL,
    dashboard_role TEXT NOT NULL,
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    auth_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'verified', 'expired', 'failed', 'used')),
    expires_at TIMESTAMPTZ NOT NULL,
    verified_at TIMESTAMPTZ,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_stakeholder_authorization_challenges_scope ON stakeholder_authorization_challenges(member_id, organization_id, action, entity_type, entity_id, status);
CREATE INDEX IF NOT EXISTS idx_stakeholder_authorization_challenges_expires_at ON stakeholder_authorization_challenges(expires_at);
CREATE INDEX IF NOT EXISTS idx_stakeholder_authorization_challenges_auth_id ON stakeholder_authorization_challenges(auth_id);

CREATE TABLE IF NOT EXISTS fund_release_requests (
    id TEXT PRIMARY KEY,
    asset_id TEXT NOT NULL,
    asset_code TEXT,
    requester_org_id TEXT NOT NULL REFERENCES organizations(id),
    requester_member_id TEXT NOT NULL REFERENCES organization_members(id),
    trustee_org_id TEXT NOT NULL REFERENCES organizations(id),
    custodian_org_id TEXT REFERENCES organizations(id),
    amount NUMERIC(30,8) NOT NULL CHECK (amount > 0),
    currency TEXT NOT NULL CHECK (char_length(currency) BETWEEN 3 AND 12),
    purpose TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'submitted' CHECK (status IN ('draft', 'submitted', 'trustee_approved', 'execution_pending', 'processing', 'completed', 'trustee_rejected', 'failed')),
    supporting_document_ids JSONB NOT NULL DEFAULT '[]'::jsonb,
    reviewed_by_member_id TEXT REFERENCES organization_members(id),
    reviewed_at TIMESTAMPTZ,
    rejection_reason TEXT,
    executed_by_member_id TEXT REFERENCES organization_members(id),
    executed_at TIMESTAMPTZ,
    execution_reference TEXT,
    failure_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_fund_release_requests_requester_org_id ON fund_release_requests(requester_org_id);
CREATE INDEX IF NOT EXISTS idx_fund_release_requests_trustee_org_id ON fund_release_requests(trustee_org_id);
CREATE INDEX IF NOT EXISTS idx_fund_release_requests_custodian_org_id ON fund_release_requests(custodian_org_id);
CREATE INDEX IF NOT EXISTS idx_fund_release_requests_asset_id ON fund_release_requests(asset_id);
CREATE INDEX IF NOT EXISTS idx_fund_release_requests_asset_code ON fund_release_requests(asset_code);
CREATE INDEX IF NOT EXISTS idx_fund_release_requests_status ON fund_release_requests(status);
CREATE INDEX IF NOT EXISTS idx_fund_release_requests_created_at ON fund_release_requests(created_at);

CREATE TABLE IF NOT EXISTS due_diligence_checklists (
    id TEXT PRIMARY KEY,
    asset_id TEXT NOT NULL,
    asset_code TEXT,
    trustee_org_id TEXT NOT NULL REFERENCES organizations(id),
    status TEXT NOT NULL DEFAULT 'in_review' CHECK (status IN ('draft', 'in_review', 'approved', 'rejected')),
    approved_by_member_id TEXT REFERENCES organization_members(id),
    approved_at TIMESTAMPTZ,
    rejected_by_member_id TEXT REFERENCES organization_members(id),
    rejected_at TIMESTAMPTZ,
    rejection_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(asset_id, trustee_org_id)
);

CREATE INDEX IF NOT EXISTS idx_due_diligence_checklists_asset_id ON due_diligence_checklists(asset_id);
CREATE INDEX IF NOT EXISTS idx_due_diligence_checklists_asset_code ON due_diligence_checklists(asset_code);
CREATE INDEX IF NOT EXISTS idx_due_diligence_checklists_trustee_org_id ON due_diligence_checklists(trustee_org_id);
CREATE INDEX IF NOT EXISTS idx_due_diligence_checklists_status ON due_diligence_checklists(status);

CREATE TABLE IF NOT EXISTS due_diligence_items (
    id TEXT PRIMARY KEY,
    checklist_id TEXT NOT NULL REFERENCES due_diligence_checklists(id) ON DELETE CASCADE,
    category TEXT NOT NULL,
    item TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'complete', 'failed', 'not_applicable')),
    notes TEXT,
    verified_by_member_id TEXT REFERENCES organization_members(id),
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_due_diligence_items_checklist_id ON due_diligence_items(checklist_id);
CREATE INDEX IF NOT EXISTS idx_due_diligence_items_status ON due_diligence_items(status);
CREATE INDEX IF NOT EXISTS idx_due_diligence_items_verified_by_member_id ON due_diligence_items(verified_by_member_id);

CREATE TABLE IF NOT EXISTS revenue_records (
    id TEXT PRIMARY KEY,
    asset_id TEXT NOT NULL,
    asset_code TEXT,
    manager_org_id TEXT NOT NULL REFERENCES organizations(id),
    period_start TIMESTAMPTZ,
    period_end TIMESTAMPTZ,
    source TEXT NOT NULL,
    amount NUMERIC(30,8) NOT NULL CHECK (amount > 0),
    currency TEXT NOT NULL CHECK (char_length(currency) BETWEEN 3 AND 12),
    status TEXT NOT NULL DEFAULT 'recorded' CHECK (status IN ('draft', 'recorded', 'submitted_for_distribution')),
    collected_at TIMESTAMPTZ,
    created_by_member_id TEXT NOT NULL REFERENCES organization_members(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_revenue_records_asset_id ON revenue_records(asset_id);
CREATE INDEX IF NOT EXISTS idx_revenue_records_asset_code ON revenue_records(asset_code);
CREATE INDEX IF NOT EXISTS idx_revenue_records_manager_org_id ON revenue_records(manager_org_id);
CREATE INDEX IF NOT EXISTS idx_revenue_records_status ON revenue_records(status);
CREATE INDEX IF NOT EXISTS idx_revenue_records_period ON revenue_records(period_start, period_end);

CREATE TABLE IF NOT EXISTS distributions (
    id TEXT PRIMARY KEY,
    asset_id TEXT NOT NULL,
    asset_code TEXT,
    proposed_by_org_id TEXT NOT NULL REFERENCES organizations(id),
    proposed_by_member_id TEXT NOT NULL REFERENCES organization_members(id),
    trustee_org_id TEXT NOT NULL REFERENCES organizations(id),
    amount NUMERIC(30,8) NOT NULL CHECK (amount > 0),
    currency TEXT NOT NULL CHECK (char_length(currency) BETWEEN 3 AND 12),
    source TEXT NOT NULL,
    scheduled_date TIMESTAMPTZ,
    status TEXT NOT NULL DEFAULT 'proposed' CHECK (status IN ('proposed', 'authorized', 'rejected', 'processing', 'completed', 'failed')),
    authorized_by_member_id TEXT REFERENCES organization_members(id),
    authorized_at TIMESTAMPTZ,
    rejection_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_distributions_asset_id ON distributions(asset_id);
CREATE INDEX IF NOT EXISTS idx_distributions_asset_code ON distributions(asset_code);
CREATE INDEX IF NOT EXISTS idx_distributions_proposed_by_org_id ON distributions(proposed_by_org_id);
CREATE INDEX IF NOT EXISTS idx_distributions_trustee_org_id ON distributions(trustee_org_id);
CREATE INDEX IF NOT EXISTS idx_distributions_status ON distributions(status);
CREATE INDEX IF NOT EXISTS idx_distributions_scheduled_date ON distributions(scheduled_date);

CREATE TABLE IF NOT EXISTS asset_valuations (
    id TEXT PRIMARY KEY,
    asset_id TEXT NOT NULL,
    asset_code TEXT,
    manager_org_id TEXT NOT NULL REFERENCES organizations(id),
    valuation NUMERIC(30,8) NOT NULL CHECK (valuation > 0),
    currency TEXT NOT NULL CHECK (char_length(currency) BETWEEN 3 AND 12),
    methodology TEXT NOT NULL,
    valuation_date TIMESTAMPTZ NOT NULL,
    report_document_id TEXT,
    status TEXT NOT NULL DEFAULT 'submitted' CHECK (status IN ('submitted', 'independent_requested', 'accepted', 'rejected')),
    created_by_member_id TEXT NOT NULL REFERENCES organization_members(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_asset_valuations_asset_id ON asset_valuations(asset_id);
CREATE INDEX IF NOT EXISTS idx_asset_valuations_asset_code ON asset_valuations(asset_code);
CREATE INDEX IF NOT EXISTS idx_asset_valuations_manager_org_id ON asset_valuations(manager_org_id);
CREATE INDEX IF NOT EXISTS idx_asset_valuations_status ON asset_valuations(status);
CREATE INDEX IF NOT EXISTS idx_asset_valuations_valuation_date ON asset_valuations(valuation_date);

CREATE TABLE IF NOT EXISTS segregated_accounts (
    id TEXT PRIMARY KEY,
    asset_id TEXT NOT NULL,
    asset_code TEXT,
    custodian_org_id TEXT NOT NULL REFERENCES organizations(id),
    account_type TEXT NOT NULL CHECK (account_type IN ('sale_proceeds', 'development_funds', 'revenue', 'reserve')),
    account_name TEXT NOT NULL,
    balance NUMERIC(30,8) NOT NULL DEFAULT 0,
    currency TEXT NOT NULL CHECK (char_length(currency) BETWEEN 3 AND 12),
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'closed')),
    bank_details JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_by_member_id TEXT REFERENCES organization_members(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_segregated_accounts_asset_id ON segregated_accounts(asset_id);
CREATE INDEX IF NOT EXISTS idx_segregated_accounts_asset_code ON segregated_accounts(asset_code);
CREATE INDEX IF NOT EXISTS idx_segregated_accounts_custodian_org_id ON segregated_accounts(custodian_org_id);
CREATE INDEX IF NOT EXISTS idx_segregated_accounts_status ON segregated_accounts(status);
CREATE INDEX IF NOT EXISTS idx_segregated_accounts_account_type ON segregated_accounts(account_type);

CREATE TABLE IF NOT EXISTS compliance_items (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL REFERENCES organizations(id),
    category TEXT NOT NULL,
    requirement TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'complete', 'overdue', 'waived')),
    due_date TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    completed_by_member_id TEXT REFERENCES organization_members(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_compliance_items_org_id ON compliance_items(org_id);
CREATE INDEX IF NOT EXISTS idx_compliance_items_status ON compliance_items(status);
CREATE INDEX IF NOT EXISTS idx_compliance_items_due_date ON compliance_items(due_date);

CREATE TABLE IF NOT EXISTS stakeholder_documents (
    id TEXT PRIMARY KEY,
    asset_id TEXT,
    asset_code TEXT,
    uploaded_by_org_id TEXT NOT NULL REFERENCES organizations(id),
    uploaded_by_member_id TEXT NOT NULL REFERENCES organization_members(id),
    category TEXT NOT NULL,
    title TEXT NOT NULL,
    file_url TEXT NOT NULL,
    mime_type TEXT,
    version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0),
    access_roles JSONB NOT NULL DEFAULT '[]'::jsonb,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'superseded', 'deleted')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_stakeholder_documents_asset_id ON stakeholder_documents(asset_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_documents_asset_code ON stakeholder_documents(asset_code);
CREATE INDEX IF NOT EXISTS idx_stakeholder_documents_uploaded_by_org_id ON stakeholder_documents(uploaded_by_org_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_documents_category ON stakeholder_documents(category);
CREATE INDEX IF NOT EXISTS idx_stakeholder_documents_status ON stakeholder_documents(status);

CREATE TABLE IF NOT EXISTS stakeholder_notification_preferences (
    id TEXT PRIMARY KEY,
    member_id TEXT NOT NULL UNIQUE REFERENCES organization_members(id),
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    email_enabled BOOLEAN NOT NULL DEFAULT true,
    in_app_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_stakeholder_notification_preferences_organization_id ON stakeholder_notification_preferences(organization_id);

CREATE TABLE IF NOT EXISTS stakeholder_asset_operations (
    id TEXT PRIMARY KEY,
    asset_id TEXT NOT NULL,
    asset_code TEXT,
    manager_org_id TEXT NOT NULL REFERENCES organizations(id),
    operational_status TEXT,
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_by_member_id TEXT NOT NULL REFERENCES organization_members(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(asset_id, manager_org_id)
);

CREATE INDEX IF NOT EXISTS idx_stakeholder_asset_operations_asset_id ON stakeholder_asset_operations(asset_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_asset_operations_asset_code ON stakeholder_asset_operations(asset_code);
CREATE INDEX IF NOT EXISTS idx_stakeholder_asset_operations_manager_org_id ON stakeholder_asset_operations(manager_org_id);

CREATE TABLE IF NOT EXISTS stakeholder_reports (
    id TEXT PRIMARY KEY,
    asset_id TEXT,
    asset_code TEXT,
    manager_org_id TEXT NOT NULL REFERENCES organizations(id),
    trustee_org_id TEXT REFERENCES organizations(id),
    report_type TEXT NOT NULL,
    title TEXT NOT NULL,
    file_url TEXT,
    status TEXT NOT NULL DEFAULT 'generated' CHECK (status IN ('generated', 'submitted')),
    generated_by_member_id TEXT NOT NULL REFERENCES organization_members(id),
    submitted_at TIMESTAMPTZ,
    submitted_by_member_id TEXT REFERENCES organization_members(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_stakeholder_reports_asset_id ON stakeholder_reports(asset_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_reports_asset_code ON stakeholder_reports(asset_code);
CREATE INDEX IF NOT EXISTS idx_stakeholder_reports_manager_org_id ON stakeholder_reports(manager_org_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_reports_trustee_org_id ON stakeholder_reports(trustee_org_id);
CREATE INDEX IF NOT EXISTS idx_stakeholder_reports_status ON stakeholder_reports(status);
