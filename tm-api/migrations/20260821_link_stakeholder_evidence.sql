-- Link managed stakeholder documents to custodian compliance requirements and
-- trustee due-diligence checklist items. AdminDB owns all affected tables.

ALTER TABLE compliance_items
    ADD COLUMN IF NOT EXISTS asset_id TEXT;

CREATE INDEX IF NOT EXISTS idx_compliance_items_asset_id
    ON compliance_items(asset_id);

CREATE TABLE IF NOT EXISTS compliance_item_documents (
    compliance_item_id TEXT NOT NULL REFERENCES compliance_items(id) ON DELETE CASCADE,
    document_id TEXT NOT NULL REFERENCES stakeholder_documents(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (compliance_item_id, document_id)
);

CREATE INDEX IF NOT EXISTS idx_compliance_item_documents_document_id
    ON compliance_item_documents(document_id);

CREATE TABLE IF NOT EXISTS due_diligence_item_documents (
    due_diligence_item_id TEXT NOT NULL REFERENCES due_diligence_items(id) ON DELETE CASCADE,
    document_id TEXT NOT NULL REFERENCES stakeholder_documents(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (due_diligence_item_id, document_id)
);

CREATE INDEX IF NOT EXISTS idx_due_diligence_item_documents_document_id
    ON due_diligence_item_documents(document_id);
