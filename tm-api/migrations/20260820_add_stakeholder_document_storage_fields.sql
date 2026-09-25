BEGIN;

ALTER TABLE stakeholder_documents
    ADD COLUMN IF NOT EXISTS storage_object_id TEXT,
    ADD COLUMN IF NOT EXISTS original_filename TEXT,
    ADD COLUMN IF NOT EXISTS size_bytes BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS sha256 TEXT;

CREATE INDEX IF NOT EXISTS idx_stakeholder_documents_storage_object_id
    ON stakeholder_documents(storage_object_id);

COMMIT;
