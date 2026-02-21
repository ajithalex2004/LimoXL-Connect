-- Create Outsource Companies Table for Limo Outsourcing
-- This table is separate from the garages table which handles maintenance garages

CREATE TABLE IF NOT EXISTS outsource_companies (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    name TEXT NOT NULL,
    contact_person TEXT,
    designation TEXT,
    email TEXT,
    contact_number TEXT,
    address TEXT,
    city TEXT,
    country TEXT,
    specialties TEXT[], -- Array of specialties (e.g., 'Luxury Vehicles', 'Airport Transfers')
    rating DECIMAL(3,2), -- Rating out of 5.00
    is_active BOOLEAN DEFAULT true,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_outsource_companies_active ON outsource_companies(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_outsource_companies_name ON outsource_companies(name);

-- Add comment to table
COMMENT ON TABLE outsource_companies IS 'Limo outsourcing partner companies - separate from maintenance garages';
