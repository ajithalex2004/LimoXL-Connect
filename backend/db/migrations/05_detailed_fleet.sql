-- Detailed Fleet Columns for Vehicles
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS chassis_no VARCHAR(100);
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS vin VARCHAR(100);
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS year_of_manufacture INT;
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS color VARCHAR(50);
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS registration_number VARCHAR(100);
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS plate_code VARCHAR(50);
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS plate_category VARCHAR(50);
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS emirate VARCHAR(50);
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS hierarchy VARCHAR(100);
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS vehicle_usage VARCHAR(100);

-- Detailed Fleet Columns for Drivers
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS first_name VARCHAR(100);
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS last_name VARCHAR(100);
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS date_of_birth DATE;
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS nationality VARCHAR(100);
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS emirates_id VARCHAR(50);
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS date_of_join DATE;
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS dallas_id VARCHAR(50);
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS communication_language VARCHAR(50);
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS hierarchy VARCHAR(100);
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS driver_type VARCHAR(50);

-- Rename contact_number to phone if needed or just use phone (Repo mapping)
DO $$ 
BEGIN
    IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'drivers' AND column_name = 'contact_number') THEN
        ALTER TABLE drivers RENAME COLUMN contact_number TO phone;
    END IF;
END $$;

-- Fleet Attachments Table
CREATE TABLE IF NOT EXISTS fleet_attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_id UUID NOT NULL,
    entity_type VARCHAR(20) NOT NULL CHECK (entity_type IN ('VEHICLE', 'DRIVER')),
    file_name VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_type VARCHAR(100),
    file_size BIGINT,
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_fleet_attachments_entity ON fleet_attachments(entity_id, entity_type);
