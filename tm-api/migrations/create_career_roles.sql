-- Create enum type for career role status
CREATE TYPE career_role_status AS ENUM ('draft', 'published', 'disabled');

-- Create career_roles table
CREATE TABLE IF NOT EXISTS career_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    heading VARCHAR(255) NOT NULL,
    company_overview TEXT,
    sections JSONB NOT NULL,
    application JSONB NOT NULL,
    status career_role_status NOT NULL DEFAULT 'draft',
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_by UUID
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_career_roles_status ON career_roles(status);
CREATE INDEX IF NOT EXISTS idx_career_roles_published_at ON career_roles(published_at) WHERE status = 'published';
CREATE INDEX IF NOT EXISTS idx_career_roles_created_at ON career_roles(created_at DESC);

-- Add comment to table
COMMENT ON TABLE career_roles IS 'Stores career/job posting information';
COMMENT ON COLUMN career_roles.status IS 'Status: draft (not published), published (visible), disabled (hidden)';
COMMENT ON COLUMN career_roles.sections IS 'JSONB containing structured sections: The Role, What You''ll Own, What We''re Looking For, What Success Looks Like';
COMMENT ON COLUMN career_roles.application IS 'JSONB containing application information: location, apply instructions, subject';
