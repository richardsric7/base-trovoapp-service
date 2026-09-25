ALTER TABLE fund_release_requests
    ADD COLUMN IF NOT EXISTS receiving_bank TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS receiving_account_name TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS receiving_account_number TEXT NOT NULL DEFAULT '';

