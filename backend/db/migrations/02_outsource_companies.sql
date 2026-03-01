-- Migration 02: Add outsource_companies table
CREATE TABLE IF NOT EXISTS outsource_companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    contact_person VARCHAR(255),
    designation VARCHAR(255),
    email VARCHAR(255),
    contact_number VARCHAR(100),
    address TEXT,
    city VARCHAR(100),
    country VARCHAR(100),
    trade_license_no VARCHAR(100),
    itc_permit_no VARCHAR(100),
    vat_no VARCHAR(100),
    rating DECIMAL(3,2),
    is_active BOOLEAN DEFAULT TRUE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);
