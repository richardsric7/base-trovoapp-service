-- Add new fields to career_roles table
ALTER TABLE career_roles 
  ADD COLUMN location VARCHAR(50) NOT NULL DEFAULT 'onsite',
  ADD COLUMN years_of_experience VARCHAR(100),
  ADD COLUMN work_type VARCHAR(50) NOT NULL DEFAULT 'fulltime';

-- Add indexes for filtering
CREATE INDEX idx_career_roles_location ON career_roles(location);
CREATE INDEX idx_career_roles_work_type ON career_roles(work_type);
