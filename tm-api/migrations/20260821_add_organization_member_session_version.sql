-- Incremented on logout to invalidate all previously issued organization JWTs.
ALTER TABLE organization_members
    ADD COLUMN IF NOT EXISTS session_version BIGINT NOT NULL DEFAULT 0;
