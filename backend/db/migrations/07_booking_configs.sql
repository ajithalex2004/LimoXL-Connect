-- Create Booking Configurations Table
CREATE TABLE IF NOT EXISTS booking_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(100) NOT NULL,
    booking_type VARCHAR(50) NOT NULL,
    request_type VARCHAR(50) NOT NULL,
    priority VARCHAR(50) NOT NULL,
    sort_order INTEGER DEFAULT 1000,
    vehicle_classes JSONB DEFAULT '[]',
    vehicle_groups JSONB DEFAULT '[]',
    vehicle_usages JSONB DEFAULT '[]',
    pickup_buffer INTEGER DEFAULT 60,
    auto_dispatch_buffer INTEGER DEFAULT 30,
    pricing_source VARCHAR(50) DEFAULT 'Zone Based',
    approval_workflow_required BOOLEAN DEFAULT false,
    epod_required BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(company_id, name)
);

CREATE INDEX idx_booking_configs_company ON booking_configs(company_id);
