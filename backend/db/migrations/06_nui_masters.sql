-- Create NUI Masters Table
CREATE TABLE IF NOT EXISTS nui_masters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id),
    category VARCHAR(50) NOT NULL, -- TYPE, CLASS, USAGE
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(company_id, category, name)
);

CREATE INDEX idx_nui_masters_company_category ON nui_masters(company_id, category);
