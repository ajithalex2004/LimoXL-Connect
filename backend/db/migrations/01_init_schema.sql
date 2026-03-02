-- Enable PostGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Sequences for Custom IDs
CREATE SEQUENCE IF NOT EXISTS trip_id_seq START 260001;
CREATE SEQUENCE IF NOT EXISTS rfq_id_seq START 1001;

-- Companies Table
CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('DEMAND', 'SUPPLY', 'BOTH')),
    contact_info JSONB DEFAULT '{}',
    verified BOOLEAN DEFAULT FALSE,
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID REFERENCES companies(id),
    role VARCHAR(50) NOT NULL CHECK (role IN ('ADMIN', 'OPS', 'SUPPLIER_ADMIN', 'DISPATCHER')),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    is_super_admin BOOLEAN DEFAULT FALSE,
    password_change_required BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Vehicles Table
CREATE TABLE IF NOT EXISTS vehicles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID REFERENCES companies(id),
    plate_number VARCHAR(50) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('SEDAN', 'SUV', 'LUXURY_VAN', 'BUS')),
    capacity INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'OFFLINE' CHECK (status IN ('IDLE', 'ON_TRIP', 'OFFLINE', 'MAINTENANCE')),
    current_location GEOGRAPHY(POINT, 4326),
    last_heartbeat TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Drivers Table
CREATE TABLE IF NOT EXISTS drivers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID REFERENCES companies(id),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    license_number VARCHAR(50),
    current_vehicle_id UUID REFERENCES vehicles(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Trips Table
CREATE TABLE IF NOT EXISTS trips (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    requesting_company_id UUID REFERENCES companies(id),
    reference_no VARCHAR(100) UNIQUE NOT NULL,
    pickup_zone VARCHAR(100), -- Added for Partner App
    dropoff_zone VARCHAR(100), -- Added for Partner App
    pickup_location GEOGRAPHY(POINT, 4326) NOT NULL,
    dropoff_location GEOGRAPHY(POINT, 4326) NOT NULL,
    pickup_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'CREATED' CHECK (status IN ('CREATED', 'MARKETPLACE_SEARCH', 'OFFERED', 'ACCEPTED', 'DRIVER_ASSIGNED', 'EN_ROUTE', 'IN_TRIP', 'COMPLETED', 'CANCELLED', 'FAILED')),
    vehicle_type_requested VARCHAR(50) NOT NULL,
    assigned_vehicle_id UUID, -- Removed FK to avoid constraints issues on demo
    assigned_driver_id UUID, -- Removed FK
    fulfillment_company_id UUID REFERENCES companies(id),
    driver_link_token VARCHAR(255), -- Added New Column
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Idempotent Updates for Existing Tables (if CREATE was skipped)
ALTER TABLE trips ADD COLUMN IF NOT EXISTS pickup_zone VARCHAR(100);
ALTER TABLE trips ADD COLUMN IF NOT EXISTS dropoff_zone VARCHAR(100);
ALTER TABLE trips ADD COLUMN IF NOT EXISTS driver_link_token VARCHAR(255);
ALTER TABLE trips ADD COLUMN IF NOT EXISTS passenger_name VARCHAR(255);
ALTER TABLE trips ADD COLUMN IF NOT EXISTS passenger_phone VARCHAR(50);
ALTER TABLE trips ADD COLUMN IF NOT EXISTS price FLOAT DEFAULT 0.0;
ALTER TABLE trips ADD COLUMN IF NOT EXISTS payment_method VARCHAR(50);
ALTER TABLE trips ADD COLUMN IF NOT EXISTS supplier_name VARCHAR(255);

ALTER TABLE users ADD COLUMN IF NOT EXISTS company_id UUID REFERENCES companies(id);
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS name VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(50);
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_super_admin BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_change_required BOOLEAN DEFAULT FALSE;

ALTER TABLE trips ADD COLUMN IF NOT EXISTS service_type VARCHAR(50);
ALTER TABLE trips ADD COLUMN IF NOT EXISTS requested_vehicle_type VARCHAR(50);
ALTER TABLE trips ADD COLUMN IF NOT EXISTS requested_vehicle_class VARCHAR(50);
ALTER TABLE trips ADD COLUMN IF NOT EXISTS requested_vehicle_group VARCHAR(50);
ALTER TABLE trips DROP CONSTRAINT IF EXISTS trips_assigned_vehicle_id_fkey; -- Cleanup old strict FK if exists

-- Trip Offers Table
CREATE TABLE IF NOT EXISTS trip_offers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_id UUID REFERENCES trips(id),
    supplier_company_id UUID REFERENCES companies(id),
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'ACCEPTED', 'REJECTED', 'TIMEOUT')),
    score FLOAT DEFAULT 0.0,
    price FLOAT DEFAULT 0.0,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE trip_offers ADD COLUMN IF NOT EXISTS price FLOAT DEFAULT 0.0;
ALTER TABLE trip_offers ADD COLUMN IF NOT EXISTS notes TEXT;

-- Indexes for Geo-spatial queries
CREATE INDEX IF NOT EXISTS idx_vehicles_location ON vehicles USING GIST (current_location);
CREATE INDEX IF NOT EXISTS idx_trips_pickup_location ON trips USING GIST (pickup_location);

-- Seed Data (Idempotent)
INSERT INTO companies (id, name, type, verified) VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a00', 'Limo Operator Co', 'BOTH', true) ON CONFLICT (id) DO NOTHING;
INSERT INTO companies (name, type, verified) VALUES ('Demo Supplier Co', 'SUPPLY', true) ON CONFLICT DO NOTHING;
INSERT INTO companies (id, name, type, verified) VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Demo Partner Co', 'DEMAND', true) ON CONFLICT (id) DO NOTHING;

-- Seed Users
INSERT INTO users (id, company_id, role, email, password_hash, name, password_change_required)
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a99', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a00', 'ADMIN', 'admin@limoxlink.com', '$2a$10$.2E6s3gnLeOf3bEn.m2O2e1RWFcs9bhDrOQXm1uo9QN/Qf34l2PxO', 'System Admin', false)
ON CONFLICT (email) DO NOTHING;

-- Seed Trips for RFQ
-- Seed Trips for RFQ
INSERT INTO trips (requesting_company_id, reference_no, pickup_zone, dropoff_zone, pickup_location, dropoff_location, pickup_time, status, vehicle_type_requested, passenger_name)
VALUES 
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a00', 'RFQ-1004', 'Dubai Mall', 'Atlantis The Palm', ST_MakePoint(55.2708, 25.1972), ST_MakePoint(55.1166, 25.1304), NOW() + INTERVAL '4 hours', 'MARKETPLACE_SEARCH', 'Sedan', 'VIP Guest'),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a00', 'RFQ-1005', 'DXB Terminal 3', 'Abu Dhabi Corniche', ST_MakePoint(55.3644, 25.2532), ST_MakePoint(54.3773, 24.4672), NOW() + INTERVAL '6 hours', 'MARKETPLACE_SEARCH', 'SUV', 'Mr. Anderson'),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a00', 'RFQ-1006', 'Burj Khalifa', 'City Tour', ST_MakePoint(55.2744, 25.1972), ST_MakePoint(55.2744, 25.1972), NOW() + INTERVAL '24 hours', 'MARKETPLACE_SEARCH', 'Van', 'Tourist Group')
ON CONFLICT (reference_no) DO UPDATE SET 
    status = 'MARKETPLACE_SEARCH',
    requesting_company_id = EXCLUDED.requesting_company_id;

-- Invoices Table
CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_id UUID REFERENCES trips(id),
    supplier_company_id UUID REFERENCES companies(id),
    invoice_number VARCHAR(50) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'PAID', 'REJECTED')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Targeted Dispatch Enhancements
ALTER TABLE trips ADD COLUMN IF NOT EXISTS visibility VARCHAR(50) DEFAULT 'PUBLIC' CHECK (visibility IN ('PUBLIC', 'DIRECT'));
ALTER TABLE trips ADD COLUMN IF NOT EXISTS pickup_landmark VARCHAR(255);
ALTER TABLE trips ADD COLUMN IF NOT EXISTS dropoff_landmark VARCHAR(255);

CREATE TABLE IF NOT EXISTS trip_access (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_id UUID REFERENCES trips(id),
    company_id UUID REFERENCES companies(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(trip_id, company_id)
);

-- Vehicle Fields Update
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS vehicle_class VARCHAR(50); -- Standard, Premium, Luxury
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS vehicle_group VARCHAR(50); -- Sedan, SUV, Van
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS model VARCHAR(100);         -- Lexus ES350, etc.
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS capacity INT DEFAULT 4;    -- Missing in remote DB

-- Driver Fields Update
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS company_id UUID; -- Critical for partner isolation
-- Companies Table Updates
ALTER TABLE companies ADD COLUMN IF NOT EXISTS trade_license_no VARCHAR(100);
ALTER TABLE companies ADD COLUMN IF NOT EXISTS itc_permit_no VARCHAR(100);
ALTER TABLE companies ADD COLUMN IF NOT EXISTS vat_no VARCHAR(100);
ALTER TABLE companies ADD COLUMN IF NOT EXISTS sla_score FLOAT DEFAULT 100.0;
ALTER TABLE companies ADD COLUMN IF NOT EXISTS rating FLOAT DEFAULT 5.0;
ALTER TABLE companies ADD COLUMN IF NOT EXISTS verification_status VARCHAR(50) DEFAULT 'PENDING' CHECK (verification_status IN ('PENDING', 'APPROVED', 'REJECTED'));

-- Vehicles Table Updates
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS permit_expiry TIMESTAMP WITH TIME ZONE;
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS insurance_expiry TIMESTAMP WITH TIME ZONE;

-- Drivers Table Updates
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS license_expiry TIMESTAMP WITH TIME ZONE;
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS itc_permit_expiry TIMESTAMP WITH TIME ZONE;
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS visa_expiry TIMESTAMP WITH TIME ZONE;

-- Invoices Table Updates
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS platform_fee DECIMAL(10, 2) DEFAULT 0.0;
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS vat_amount DECIMAL(10, 2) DEFAULT 0.0;
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS net_payout DECIMAL(10, 2) DEFAULT 0.0;

-- Idempotent check for old invoices table reference (trip_invoices was used in code)
DO $$ 
BEGIN
    IF EXISTS (SELECT FROM pg_tables WHERE tablename = 'invoices') AND NOT EXISTS (SELECT FROM pg_tables WHERE tablename = 'trip_invoices') THEN
        ALTER TABLE invoices RENAME TO trip_invoices;
    END IF;
END $$;
